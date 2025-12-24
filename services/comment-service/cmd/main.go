package main

import (
	"comment_service/internal/config"
	"comment_service/internal/handlers"
	"comment_service/internal/middleware"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "comment_service/docs"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/gofiber/websocket/v2"
)

// @title comment_service API
// @version 1.0
// @description API Backend cho hệ thống bình luận sản phẩm
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
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

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "*",
	// 	// Cho phép tất cả các Header tùy chỉnh từ trình duyệt
	// 	AllowHeaders: "*",
	// 	// Thêm OPTIONS vào danh sách phương thức để trình duyệt có thể thực hiện Preflight check
	// 	AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	// 	// Tùy chọn: Cho phép trình duyệt lưu kết quả Preflight trong một khoảng thời gian (giây)
	// 	MaxAge: 86400,
	// }))

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
	commentHandler := handlers.NewCommentHandler(db)

	// Routes
	api := app.Group("")

	// WebSocket endpoint for comments
	app.Get("/ws", websocket.New(commentHandler.HandleWebSocket))

	// Comment routes (protected by auth middleware)
	comments := api.Group("", middleware.ExtractUserInfo(cfg))
	comments.Get("/products/:productId", commentHandler.GetProductComments)

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
