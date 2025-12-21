package config

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		// Không panic, có thể dùng environment variables
	}
}

type Config struct {
	Port                    string
	APIGatewaySecret        string
	AuthInternalSecret      string
	AuthServiceURL          string
	CategoryServiceURL      string
	ProductServiceURL       string
	UserServiceURL          string
	BiddingServiceURL       string
	OrderServiceURL         string
	PaymentServiceURL       string
	NotificationServiceURL  string
	MediaServiceURL         string
	OTelEndpoint            string
	OTelServiceName         string
	OTelServiceVersion      string
	OTelEnvironment         string
	JWTPublicKeyAuthService string
	JWTPrivateKey           string
	CategoryServiceName     string
	ProductServiceName      string
	UserServiceName         string
	BiddingServiceName      string
	OrderServiceName        string
	PaymentServiceName      string
	NotificationServiceName string
	MediaServiceName        string
	SearchServiceName       string
	CommentServiceName      string
	AutoBiddingServiceName  string
}

func LoadConfig() *Config {
	return &Config{
		Port:                    getEnv("PORT", "8080"),
		APIGatewaySecret:        getEnv("API_GATEWAY_SECRET", "api-gateway-secret"),
		AuthInternalSecret:      getEnv("AUTH_INTERNAL_SECRET", "internal-auth-secret"),
		AuthServiceURL:          getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		CategoryServiceURL:      getEnv("CATEGORY_SERVICE_URL", "http://localhost:8082"),
		ProductServiceURL:       getEnv("PRODUCT_SERVICE_URL", "http://localhost:8083"),
		UserServiceURL:          getEnv("USER_SERVICE_URL", "http://localhost:8084"),
		BiddingServiceURL:       getEnv("BIDDING_SERVICE_URL", "http://localhost:8085"),
		OrderServiceURL:         getEnv("ORDER_SERVICE_URL", "http://localhost:8086"),
		PaymentServiceURL:       getEnv("PAYMENT_SERVICE_URL", "http://localhost:8087"),
		NotificationServiceURL:  getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8088"),
		MediaServiceURL:         getEnv("MEDIA_SERVICE_URL", "http://localhost:8089"),
		OTelEndpoint:            getEnv("OTEL_ENDPOINT", "localhost:4317"),
		OTelServiceName:         getEnv("OTEL_SERVICE_NAME", "api-gateway"),
		OTelServiceVersion:      getEnv("OTEL_SERVICE_VERSION", "1.0.0"),
		OTelEnvironment:         getEnv("OTEL_ENVIRONMENT", "development"),
		JWTPublicKeyAuthService: getEnv("JWT_PUBLIC_KEY_AUTH_SERVICE", ""),
		JWTPrivateKey:           getEnv("JWT_PRIVATE_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
