package main

import (
	"log"
	"log/slog"
	"order_service/internal/config"
	"order_service/internal/handlers"
	"order_service/internal/middleware"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "order_service/docs"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/gofiber/websocket/v2"
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

	// Connect database
	db := config.ConnectDB(cfg)
	defer db.Close()

	// Init schema
	if err := config.InitSchema(db); err != nil {
		log.Fatalf("Lỗi khởi tạo schema: %v", err)
	}

	// Seed sample data (commented out - uncomment to seed database)
	// if err := config.SeedData(db); err != nil {
	// 	log.Fatalf("Lỗi seeding dữ liệu: %v", err)
	// }

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

	// CORS middleware - Must be enabled for frontend to work
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-User-Token, X-Internal-JWT, Upgrade, Connection, Sec-WebSocket-Key, Sec-WebSocket-Version")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Max-Age", "86400")

		// Handle preflight OPTIONS request
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusOK)
		}

		return c.Next()
	})

	// Middleware
	app.Use(recover.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} | ${path}\n",
	}))

	// Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Serve static files for testing
	app.Static("/", "./public")

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(db, cfg)
	likeHandler := handlers.NewLikeHandler(db, cfg)
	api := app.Group("")

	app.Get("/ws", websocket.New(orderHandler.HandleWebSocket))

	// Routes

	// WatchList routes (danh sách yêu thích) - /data/watchlist
	watchlist := api.Group("/watchlist", middleware.ExtractUserInfo(cfg))
	watchlist.Post("/", likeHandler.AddToWatchList)                   // Add product to watch list
	watchlist.Get("/", likeHandler.GetWatchList)                      // Get user's watch list
	watchlist.Delete("/:product_id", likeHandler.RemoveFromWatchList) // Remove from watch list
	watchlist.Get("/:product_id/check", likeHandler.CheckInWatchList) // Check if in watch list

	// Order routes
	orders := api.Group("")

	// -------------------------------------------------
	// PUBLIC endpoint (no middleware)
	// Called by auction service
	// -------------------------------------------------
	orders.Post("order/", orderHandler.CreateOrder)                                // Create order (no auth - called by auction service)

	// -------------------------------------------------
	// PROTECTED endpoints (middleware applies from here)
	// -------------------------------------------------
	orders.Use(middleware.ExtractUserInfo(cfg))

	orders.Get("order/", orderHandler.GetUserOrders)                               // Get user's orders
	orders.Get("order/:id", orderHandler.GetOrderByID)                             // Get order by ID
	orders.Post("order/:id/pay", orderHandler.PayOrder)                            // Pay for order
	orders.Post("order/:id/shipping-address", orderHandler.ProvideShippingAddress) // Provide shipping address
	orders.Post("order/:id/shipping-invoice", orderHandler.SendShippingInvoice)    // Send shipping invoice
	orders.Post("order/:id/confirm-delivery", orderHandler.ConfirmDelivery)        // Confirm delivery
	orders.Post("order/:id/cancel", orderHandler.CancelOrder)                      // Cancel order
	orders.Get("order/:id/messages", orderHandler.GetMessages)                     // Get chat history (REST API for initial load)
	orders.Post("order/:id/rate", orderHandler.RateOrder)                          // Rate order
	orders.Get("order/:id/rating", orderHandler.GetRating)                         // Get rating (public)
	// User rating routes
	api.Get("/users/:id/rating", middleware.ExtractUserInfo(cfg), orderHandler.GetUserRating) // Get user rating stats (public)

	// Admin routes
	admin := api.Group("/admin", middleware.ExtractUserInfo(cfg), middleware.AdminMiddleware())
	admin.Get("/orders", orderHandler.GetAllOrders) // Get all orders (admin only)

	// WebSocket endpoint for order chat

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start HTTP server
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
