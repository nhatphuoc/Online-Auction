package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"search-service/internal/config"
	"search-service/internal/elasticsearch"
	"search-service/internal/handlers"
	"search-service/internal/models"
	"search-service/internal/repository"
	"search-service/internal/stream"
	"search-service/internal/worker"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.LoadConfig()
	log.Println("Configuration loaded")

	db := config.ConnectDB(cfg)
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	esClient, err := elasticsearch.NewClient(cfg.ElasticsearchURL)
	if err != nil {
		log.Fatalf("Elasticsearch error: %v", err)
	}

	productIndexExists, _ := elasticsearch.IndexExists(ctx, esClient, cfg.ElasticsearchIndexProduct)
	if !productIndexExists {
		if err := elasticsearch.CreateProductIndex(ctx, esClient, cfg.ElasticsearchIndexProduct); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	categoryIndexExists, _ := elasticsearch.IndexExists(ctx, esClient, cfg.ElasticsearchIndexCategory)
	if !categoryIndexExists {
		if err := elasticsearch.CreateCategoryIndex(ctx, esClient, cfg.ElasticsearchIndexCategory); err != nil {
			log.Printf("Warning: %v", err)
		}
	}

	indexer := elasticsearch.NewIndexer(esClient, cfg.ElasticsearchIndexProduct, cfg.ElasticsearchIndexCategory)
	searcher := elasticsearch.NewSearcher(esClient, cfg.ElasticsearchIndexProduct, cfg.BoostMinutes, cfg.BoostScore)

	syncWorker := worker.NewSyncWorker(productRepo, categoryRepo, indexer)

	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	consumer, err := stream.NewConsumer(redisAddr, cfg.RedisPassword, cfg.RedisDB, cfg.RedisStreamKey, cfg.RedisConsumerGroup, cfg.RedisConsumerName)
	if err != nil {
		log.Fatalf("Redis error: %v", err)
	}
	defer consumer.Close()

	go func() {
		log.Println("Starting worker...")
		if err := consumer.ReadMessages(ctx, func(event *models.Event) error {
			return syncWorker.HandleEvent(ctx, event)
		}); err != nil {
			log.Printf("Worker error: %v", err)
		}
	}()

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	productHandler := handlers.NewProductHandler(searcher)
	api := app.Group("/api")
	api.Get("/search/products", productHandler.SearchProducts)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down...")
		cancel()
		app.Shutdown()
	}()

	port := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Running on port %s", cfg.Port)
	if err := app.Listen(port); err != nil {
		log.Fatal(err)
	}
}
