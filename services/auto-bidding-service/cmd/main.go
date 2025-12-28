package main

import (
	"auto-bidding-service/internal/client"
	"auto-bidding-service/internal/config"
	"auto-bidding-service/internal/handlers"
	"auto-bidding-service/internal/logger"
	"auto-bidding-service/internal/metrics"
	"auto-bidding-service/internal/middleware"
	"auto-bidding-service/internal/repository"
	"auto-bidding-service/internal/service"
	"auto-bidding-service/internal/telemetry"
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.OTelEnvironment)
	slog.Info("Starting auto-bidding-service API")

	otelShutdown, err := telemetry.InitOTel(ctx, telemetry.OTelConfig{
		ServiceName:    cfg.OTelServiceName,
		ServiceVersion: cfg.OTelServiceVersion,
		Environment:    cfg.OTelEnvironment,
		OTelEndpoint:   cfg.OTelEndpoint,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer otelShutdown(ctx)

	logger.InitLoggerWithOTel(cfg.OTelEnvironment)
	metrics.InitMetrics(ctx)

	db := config.ConnectDB(cfg)
	defer db.Close()
	config.InitSchema(db)

	autoBidRepo := repository.NewAutoBidRepository(db)
	biddingClient := client.NewBiddingServiceClient(os.Getenv("BIDDING_SERVICE_URL"))
	productClient := client.NewProductServiceClient(os.Getenv("PRODUCT_SERVICE_URL"))
	autoBidService := service.NewAutoBidService(autoBidRepo, biddingClient, productClient)
	autoBidHandler := handlers.NewAutoBidHandler(autoBidService)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(fiberlogger.New())
	
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
	
	app.Use(middleware.TracingMiddleware())

	api := app.Group("/api")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	autoBids := api.Group("/auto-bids")
	autoBids.Post("/", middleware.AuthMiddleware(cfg), autoBidHandler.CreateAutoBid)
	autoBids.Post("/trigger", autoBidHandler.TriggerAutoBidding)
	autoBids.Get("/my", middleware.AuthMiddleware(cfg), autoBidHandler.GetMyAutoBids)
	autoBids.Get("/:id", middleware.AuthMiddleware(cfg), autoBidHandler.GetAutoBidByID)
	autoBids.Post("/:id/cancel", middleware.AuthMiddleware(cfg), autoBidHandler.CancelAutoBid)
	app.Get("/swagger/*", swagger.HandlerDefault)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}

	slog.Info("Starting server", "port", port)
	app.Listen(":" + port)
}
