package main

import (
	"log"
	"log/slog"
	"media_service/internal/config"
	"media_service/internal/handlers"
	"media_service/internal/middleware"
	"os"

	_ "media_service/docs"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// @title Media Service API
// @version 1.0
// @description API for uploading and managing media files (images, videos) to AWS S3
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api
// @schemes http https

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize S3 client
	s3Client, err := config.InitS3Client(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: int(cfg.MaxFileSize),
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
	
	// CORS middleware - Must handle OPTIONS properly
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-User-Token, X-Internal-JWT")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Max-Age", "86400")
		
		// Handle preflight OPTIONS request
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusOK)
		}
		
		return c.Next()
	})

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Initialize handlers
	uploadHandler := handlers.NewUploadHandler(s3Client, cfg)

	// API routes
	api := app.Group("")

	// Upload endpoints
	// api.Post("/upload", uploadHandler.UploadSingleFile)
	// api.Post("/upload/multiple", uploadHandler.UploadMultipleFiles)
	// Presigned URL endpoints
	api.Get("/presign", middleware.ExtractUserInfo(cfg), uploadHandler.GetPresignedURL)
	api.Post("/presign/multiple", middleware.ExtractUserInfo(cfg), uploadHandler.GetPresignedURLs)
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "media-service",
			"bucket":  cfg.AWSBucketName,
			"region":  cfg.AWSRegion,
		})
	})

	// Start server
	slog.Info("Starting media service", "port", cfg.Port, "bucket", cfg.AWSBucketName, "region", cfg.AWSRegion)
	slog.Info("Swagger documentation", "url", "http://localhost:"+cfg.Port+"/swagger/")

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
