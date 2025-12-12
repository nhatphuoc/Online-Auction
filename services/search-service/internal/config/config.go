package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		// Try loading from parent directory for development
		_ = godotenv.Load("../.env")
	}
}

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string

	// Elasticsearch
	ElasticsearchURL          string
	ElasticsearchIndexProduct string
	ElasticsearchIndexCategory string

	// Redis Stream
	RedisHost          string
	RedisPort          string
	RedisPassword      string
	RedisDB            int
	RedisStreamKey     string
	RedisConsumerGroup string
	RedisConsumerName  string

	// Search Configuration
	BoostMinutes int
	BoostScore   float64
}

func LoadConfig() *Config {
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	boostMinutes, _ := strconv.Atoi(getEnv("BOOST_MINUTES", "60"))
	boostScore, _ := strconv.ParseFloat(getEnv("BOOST_SCORE", "2.0"), 64)

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "neondb"),
		Port:       getEnv("PORT", "3000"),

		ElasticsearchURL:           getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),
		ElasticsearchIndexProduct:  getEnv("ELASTICSEARCH_INDEX_PRODUCT", "products"),
		ElasticsearchIndexCategory: getEnv("ELASTICSEARCH_INDEX_CATEGORY", "categories"),

		RedisHost:          getEnv("REDIS_HOST", "localhost"),
		RedisPort:          getEnv("REDIS_PORT", "6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            redisDB,
		RedisStreamKey:     getEnv("REDIS_STREAM_KEY", "auction_events"),
		RedisConsumerGroup: getEnv("REDIS_CONSUMER_GROUP", "search_service_group"),
		RedisConsumerName:  getEnv("REDIS_CONSUMER_NAME", "search_service_consumer_1"),

		BoostMinutes: boostMinutes,
		BoostScore:   boostScore,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
