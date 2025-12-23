package config

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../shared/.env")
	if err != nil {
		// Không panic, có thể dùng environment variables
	}
}

type Config struct {
	Port                 string
	APIGatewayPrivateKey string
	AuthInternalSecret   string

	AuthServiceURL             string
	CategoryServiceURL         string
	ProductServiceURL          string
	UserServiceURL             string
	BiddingServiceURL          string
	OrderServiceURL            string
	PaymentServiceURL          string
	NotificationServiceURL     string
	MediaServiceURL            string
	SearchServiceURL           string
	CommentServiceURL          string
	AutoBiddingServiceURL      string
	CommentServiceWebSocketURL string

	OTelEndpoint            string
	OTelServiceName         string
	OTelServiceVersion      string
	OTelEnvironment         string
	JWTPublicKeyAuthService string
	JWTPrivateKey           string

	APIGatewayName          string
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
		Port:                       getEnv("API_GATEWAY_PORT", "8080"),
		APIGatewayPrivateKey:       getEnv("API_GATEWAY_SECRET", "api-gateway-secret"),
		AuthInternalSecret:         getEnv("X_AUTH_INTERNAL_KEY", "internal-auth-secret"),
		AuthServiceURL:             getEnv("AUTH_SERVICE_URL", "http://localhost:8081/auth"),
		CategoryServiceURL:         getEnv("CATEGORY_SERVICE_URL", "http://localhost:8082"),
		ProductServiceURL:          getEnv("PRODUCT_SERVICE_URL", "http://localhost:8083"),
		UserServiceURL:             getEnv("USER_SERVICE_URL", "http://localhost:8084"),
		BiddingServiceURL:          getEnv("BIDDING_SERVICE_URL", "http://localhost:8085"),
		OrderServiceURL:            getEnv("ORDER_SERVICE_URL", "http://localhost:8086"),
		PaymentServiceURL:          getEnv("PAYMENT_SERVICE_URL", "http://localhost:8087"),
		NotificationServiceURL:     getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8088"),
		MediaServiceURL:            getEnv("MEDIA_SERVICE_URL", "http://localhost:8089"),
		SearchServiceURL:           getEnv("SEARCH_SERVICE_URL", "http://localhost:8090"),
		CommentServiceURL:          getEnv("COMMENT_SERVICE_URL", "http://localhost:8091"),
		AutoBiddingServiceURL:      getEnv("AUTO_BIDDING_SERVICE_URL", "http://localhost:8092"),
		CommentServiceWebSocketURL: getEnv("COMMENT_SERVICE_WEBSOCKET_URL", "ws://localhost:8091/ws"),

		APIGatewayName:          getEnv("API_GATEWAY_NAME", "api-gateway"),
		CategoryServiceName:     getEnv("CATEGORY_SERVICE_NAME", "category-service"),
		ProductServiceName:      getEnv("PRODUCT_SERVICE_NAME", "product-service"),
		UserServiceName:         getEnv("USER_SERVICE_NAME", "user-service"),
		BiddingServiceName:      getEnv("BIDDING_SERVICE_NAME", "bidding-service"),
		OrderServiceName:        getEnv("ORDER_SERVICE_NAME", "order-service"),
		PaymentServiceName:      getEnv("PAYMENT_SERVICE_NAME", "payment-service"),
		NotificationServiceName: getEnv("NOTIFICATION_SERVICE_NAME", "notification-service"),
		MediaServiceName:        getEnv("MEDIA_SERVICE_NAME", "media-service"),
		SearchServiceName:       getEnv("SEARCH_SERVICE_NAME", "search-service"),
		CommentServiceName:      getEnv("COMMENT_SERVICE_NAME", "comment-service"),
		AutoBiddingServiceName:  getEnv("AUTO_BIDDING_SERVICE_NAME", "auto-bidding-service"),

		OTelEndpoint:       getEnv("OTEL_ENDPOINT", "localhost:4317"),
		OTelServiceName:    getEnv("OTEL_SERVICE_NAME", "api-gateway"),
		OTelServiceVersion: getEnv("OTEL_SERVICE_VERSION", "1.0.0"),
		OTelEnvironment:    getEnv("OTEL_ENVIRONMENT", "development"),

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
