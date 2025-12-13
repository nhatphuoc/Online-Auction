package main

import (
	"log"
	"log/slog"
	"order_service/internal/config"
	"order_service/internal/handlers"
	"order_service/internal/middleware"
	"order_service/internal/repository"
	"order_service/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "order_service/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// @title Order Service API
// @version 1.0
// @description API Backend cho hệ thống quản lý đơn hàng sau đấu giá - Order Management Service
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

	// Load config
	cfg := config.LoadConfig()

	// Initialize logger (without OTel first)
	slog.Info("Starting order_service API", "version", cfg.OTelServiceVersion, "env", cfg.OTelEnvironment)

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
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Auth middleware (extracts user info from headers set by API Gateway)
	authMiddleware := middleware.AuthMiddleware()

	// Routes
	api := app.Group("/api")

	// Order routes
	orders := api.Group("/orders")
	orders.Post("/", orderHandler.CreateOrder)                                    // Create order (internal use, no auth)
	orders.Get("/", authMiddleware, orderHandler.GetUserOrders)                   // Get user's orders (requires auth)
	orders.Get("/:id", authMiddleware, orderHandler.GetOrderByID)                 // Get order by ID (requires auth)
	orders.Patch("/:id/status", authMiddleware, orderHandler.UpdateOrderStatus)   // Update order status (requires auth)
	orders.Post("/:id/cancel", authMiddleware, orderHandler.CancelOrder)          // Cancel order (requires auth)
	orders.Post("/:id/messages", authMiddleware, orderHandler.SendMessage)        // Send message (requires auth)
	orders.Get("/:id/messages", authMiddleware, orderHandler.GetMessages)         // Get messages (requires auth)
	orders.Post("/:id/rate", authMiddleware, orderHandler.RateOrder)              // Rate order (requires auth)
	orders.Get("/:id/rating", orderHandler.GetRating)                             // Get rating (public)

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
