package main

import (
	"api_gateway/internal/config"
	"api_gateway/internal/handlers"
	"api_gateway/internal/logger"
	"api_gateway/internal/metrics"
	"api_gateway/internal/middleware"
	"api_gateway/internal/telemetry"
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "api_gateway/docs"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// @title API Gateway
// @version 1.0
// @description API Gateway cho hệ thống đấu giá trực tuyến - Proxy requests to internal services
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey X-User-Token
// @in header
// @name X-User-Token
// @description JWT token without Bearer prefix

func main() {
	ctx := context.Background()

	// Load config
	cfg := config.LoadConfig()

	// Initialize logger
	logger.InitLogger(cfg.OTelEnvironment)
	slog.Info("Starting API Gateway", "version", cfg.OTelServiceVersion, "env", cfg.OTelEnvironment)

	// Initialize OpenTelemetry
	otelShutdown, err := telemetry.InitOTel(ctx, telemetry.OTelConfig{
		ServiceName:    cfg.OTelServiceName,
		ServiceVersion: cfg.OTelServiceVersion,
		Environment:    cfg.OTelEnvironment,
		OTelEndpoint:   cfg.OTelEndpoint,
	})
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer func() {
		if err := otelShutdown(ctx); err != nil {
			slog.Error("Error shutting down OpenTelemetry", "error", err)
		}
	}()

	// Re-initialize logger with OTel bridge
	logger.InitLoggerWithOTel(cfg.OTelEnvironment)
	slog.Info("Logger with OpenTelemetry initialized")

	// Initialize metrics
	if err := metrics.InitMetrics(ctx); err != nil {
		log.Fatalf("Failed to initialize metrics: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(middleware.TracingMiddleware())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} | ${path}\n",
	}))
	
	// CORS middleware - CRITICAL: Must be enabled for frontend to work properly
	// This handles OPTIONS preflight requests from browsers
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-User-Token, X-Internal-JWT")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Max-Age", "86400") // 24 hours cache for preflight
		
		// Handle preflight OPTIONS request
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusOK)
		}
		
		return c.Next()
	})

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Initialize handlers
	proxyHandler := handlers.NewProxyHandler(cfg)

	// Health check
	app.Get("/health", proxyHandler.HealthCheck)

	// API routes with authentication
	api := app.Group("/api")

	// Auth service routes (no auth required for login/register)
	auth := api.Group("/auth")
	auth.All("/*", middleware.ProxyMiddlewareForAuthenService(cfg), proxyHandler.ProxyRequest(cfg.AuthServiceURL))

	// Protected routes - require authentication
	protected := api.Group("", middleware.AuthMiddleware(cfg))

	// Category service
	protected.All("/categories/*", middleware.ProxyMiddleware(cfg, cfg.CategoryServiceName), proxyHandler.ProxyRequest(cfg.CategoryServiceURL))

	// Product service
	protected.All("/products/*", middleware.ProxyMiddleware(cfg, cfg.ProductServiceName), proxyHandler.ProxyRequest(cfg.ProductServiceURL))

	// User service
	protected.All("/users/*", middleware.ProxyMiddleware(cfg, cfg.UserServiceName), proxyHandler.ProxyRequest(cfg.UserServiceURL))

	// Bidding service
	protected.All("/bids/*", middleware.ProxyMiddleware(cfg, cfg.BiddingServiceName), proxyHandler.ProxyRequest(cfg.BiddingServiceURL))

	// Order service
	protected.All("/orders/*", middleware.ProxyMiddleware(cfg, cfg.OrderServiceName), proxyHandler.ProxyRequest(cfg.OrderServiceURL))
	protected.All("/order-websocket/*", middleware.ProxyMiddleware(cfg, cfg.OrderServiceName), proxyHandler.OrderProxyWebSocket)

	// Payment service
	protected.All("/payments/*", middleware.ProxyMiddleware(cfg, cfg.PaymentServiceName), proxyHandler.ProxyRequest(cfg.PaymentServiceURL))

	// Notification service
	protected.All("/notifications/*", middleware.ProxyMiddleware(cfg, cfg.NotificationServiceName), proxyHandler.ProxyRequest(cfg.NotificationServiceURL))

	// Media service
	protected.All("/media/*", middleware.ProxyMiddleware(cfg, cfg.MediaServiceName), proxyHandler.ProxyRequest(cfg.MediaServiceURL))

	// Comment service
	protected.All("/comments/history/*", middleware.ProxyMiddleware(cfg, cfg.CommentServiceName), proxyHandler.ProxyRequest(cfg.CommentServiceURL))
	protected.All("/comments/websocket/*", middleware.ProxyMiddleware(cfg, cfg.CommentServiceName), proxyHandler.ProxyWebSocket)
	// Search service
	protected.All("/search/*", middleware.ProxyMiddleware(cfg, cfg.SearchServiceName), proxyHandler.ProxyRequest(cfg.ProductServiceURL))

	// Auto Bidding service
	protected.All("/auto-bidding/*", middleware.ProxyMiddleware(cfg, cfg.AutoBiddingServiceName), proxyHandler.ProxyRequest(cfg.BiddingServiceURL))

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start server
	go func() {
		slog.Info("API Gateway started", "port", cfg.Port, "swagger", "http://localhost:"+cfg.Port+"/swagger/")
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	slog.Info("Shutting down server...")

	// Graceful shutdown
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}

	slog.Info("Server stopped gracefully")
}
