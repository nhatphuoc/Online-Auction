package main

import (
	"category_service/internal/config"
	"category_service/internal/handlers"
	"category_service/internal/middleware"
	"category_service/scripts"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "category_service/docs"

	"github.com/gofiber/fiber/v2"
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

	// Load config
	cfg := config.LoadConfig()

	// Connect database
	db := config.ConnectDB(cfg)
	defer db.Close()

	// Init schema
	if err := config.InitSchema(db); err != nil {
		log.Fatalf("Lỗi khởi tạo schema: %v", err)
	}

	// Seed initial data
	if err := scripts.SeedInitialData(db); err != nil {
		log.Fatalf("Lỗi seed dữ liệu ban đầu: %v", err)
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
	
	// CORS middleware - IMPORTANT: Must be enabled for frontend to work
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-User-Token, X-Internal-JWT")
		c.Set("Access-Control-Allow-Credentials", "true")
		
		// Handle preflight OPTIONS request
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusOK)
		}
		
		return c.Next()
	})

	// Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Initialize handlers
	categoryHandler := handlers.NewCategoryHandler(db)

	// Category routes
	categories := app.Group("", middleware.ExtractUserInfo(cfg))
	categories.Get("/", categoryHandler.GetCategories)
	categories.Get("/parent/:parent_id", categoryHandler.GetCategoriesByParent)
	categories.Get("/:id", categoryHandler.GetCategoryByID)
	categories.Post("/", middleware.RequireAdminRole(), categoryHandler.CreateCategory)
	categories.Put("/:id", middleware.RequireAdminRole(), categoryHandler.UpdateCategory)
	categories.Delete("/:id", middleware.RequireAdminRole(), categoryHandler.DeleteCategory)

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
