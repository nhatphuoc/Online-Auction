package config

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file")
	}
}

type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	JWTSecret          string
	Port               string
	GRPCPort           string
	OTelEndpoint       string
	OTelServiceName    string
	OTelServiceVersion string
	OTelEnvironment    string
}

func LoadConfig() *Config {

	return &Config{
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", "neondb"),
		JWTSecret:          getEnv("JWT_SECRET", "secret"),
		Port:               getEnv("PORT", "3000"),
		GRPCPort:           getEnv("GRPC_PORT", "50051"),
		OTelEndpoint:       getEnv("OTEL_ENDPOINT", "localhost:4317"),
		OTelServiceName:    getEnv("OTEL_SERVICE_NAME", "final4-api"),
		OTelServiceVersion: getEnv("OTEL_SERVICE_VERSION", "1.0.0"),
		OTelEnvironment:    getEnv("OTEL_ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
