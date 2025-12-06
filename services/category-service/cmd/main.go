package main

import (
	"category_service/internal/config"
	"category_service/internal/handlers"
	"category_service/internal/logger"
	"category_service/internal/metrics"
	"category_service/internal/middleware"
	"category_service/internal/telemetry"
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "category_service/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// @title Category Service API
// @version 1.0
// @description API Backend cho hệ thống quản lý danh mục và sản phẩm
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	ctx := context.Background()

	// Load config
	cfg := config.LoadConfig()

	// Initialize logger (without OTel first)
	logger.InitLogger(cfg.OTelEnvironment)
	slog.Info("Starting category_service API", "version", cfg.OTelServiceVersion, "env", cfg.OTelEnvironment)

	// Initialize OpenTelemetry
	otelShutdown, err := telemetry.InitOTel(ctx, telemetry.OTelConfig{
		ServiceName:    cfg.OTelServiceName,
		ServiceVersion: cfg.OTelServiceVersion,
		Environment:    cfg.OTelEnvironment,
		OTelEndpoint:   cfg.OTelEndpoint,
	})
	if err != nil {
		log.Fatalf("Lỗi khởi tạo OpenTelemetry: %v", err)
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
		log.Fatalf("Lỗi khởi tạo metrics: %v", err)
	}

	// Connect database
	db := config.ConnectDB(cfg)
	defer db.Close()

	// Init schema
	if err := config.InitSchema(db); err != nil {
		log.Fatalf("Lỗi khởi tạo schema: %v", err)
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
	app.Use(middleware.TracingMiddleware()) // OpenTelemetry tracing
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} | ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH",
	}))

	// Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(db)
	productHandler := handlers.NewProductHandler(db)

	// Routes
	api := app.Group("/api")

	// Category routes
	categories := api.Group("/categories")
	categories.Get("/", categoryHandler.GetCategories)
	categories.Get("/:id", categoryHandler.GetCategoryByID)
	categories.Get("/parent/:parent_id", categoryHandler.GetCategoriesByParent)
	categories.Post("/", middleware.AuthMiddleware(cfg), categoryHandler.CreateCategory)
	categories.Put("/:id", middleware.AuthMiddleware(cfg), categoryHandler.UpdateCategory)
	categories.Delete("/:id", middleware.AuthMiddleware(cfg), categoryHandler.DeleteCategory)

	// Product routes
	products := api.Group("/products")
	products.Get("/", productHandler.GetProductsByCategory)
	products.Get("/:id", productHandler.GetProductByID)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start HTTP server trong goroutine
	go func() {
		slog.Info("HTTP Server started", "port", cfg.Port, "swagger", "http://localhost:"+cfg.Port+"/swagger/")
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Lỗi chạy HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	slog.Info("Shutting down servers...")

	// Graceful shutdown HTTP server
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		slog.Error("HTTP server shutdown error", "error", err)
	}

	slog.Info("Servers stopped gracefully")
}
