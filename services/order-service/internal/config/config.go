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
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	PublicKeys  map[string]string
	ServiceURLs map[string]string

	JWTPublicKeyAuthService string
	JWTPrivateKey           string

	APIGatewayName          string
	AuthServiceName         string
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

	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSBucketName      string
	MaxFileSize        int64
	MaxFilesPerUpload  int
}

func LoadConfig() *Config {
	return &Config{
		Port: getEnv("ORDER_SERVICE_PORT", "8086"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "neondb"),

		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSRegion:          getEnv("AWS_REGION", "ap-southeast-1"),
		AWSBucketName:      getEnv("AWS_BUCKET_NAME", ""),
		MaxFileSize:        50 * 1024 * 1024, // 50MB
		MaxFilesPerUpload:  10,

		PublicKeys: map[string]string{
			"api-gateway":          getEnv("JWT_PUBLIC_KEY_API_GATEWAY", ""),
			"auth-service":         getEnv("JWT_PUBLIC_KEY_AUTH_SERVICE", ""),
			"category-service":     getEnv("JWT_PUBLIC_KEY_CATEGORY_SERVICE", ""),
			"product-service":      getEnv("JWT_PUBLIC_KEY_PRODUCT_SERVICE", ""),
			"user-service":         getEnv("JWT_PUBLIC_KEY_USER_SERVICE", ""),
			"bidding-service":      getEnv("JWT_PUBLIC_KEY_BIDDING_SERVICE", ""),
			"order-service":        getEnv("JWT_PUBLIC_KEY_ORDER_SERVICE", ""),
			"payment-service":      getEnv("JWT_PUBLIC_KEY_PAYMENT_SERVICE", ""),
			"notification-service": getEnv("JWT_PUBLIC_KEY_NOTIFICATION_SERVICE", ""),
			"media-service":        getEnv("JWT_PUBLIC_KEY_MEDIA_SERVICE", ""),
			"search-service":       getEnv("JWT_PUBLIC_KEY_SEARCH_SERVICE", ""),
			"comment-service":      getEnv("JWT_PUBLIC_KEY_COMMENT_SERVICE", ""),
			"auto-bidding-service": getEnv("JWT_PUBLIC_KEY_AUTO_BIDDING_SERVICE", ""),
		},

		ServiceURLs: map[string]string{
			"auth-service":         getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
			"category-service":     getEnv("CATEGORY_SERVICE_URL", "http://localhost:8082"),
			"product-service":      getEnv("PRODUCT_SERVICE_URL", "http://localhost:8083"),
			"user-service":         getEnv("USER_SERVICE_URL", "http://localhost:8084"),
			"bidding-service":      getEnv("BIDDING_SERVICE_URL", "http://localhost:8085"),
			"order-service":        getEnv("ORDER_SERVICE_URL", "http://localhost:8086"),
			"payment-service":      getEnv("PAYMENT_SERVICE_URL", "http://localhost:8087"),
			"notification-service": getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8088"),
			"media-service":        getEnv("MEDIA_SERVICE_URL", "http://localhost:8089"),
			"comment-service":      getEnv("COMMENT_SERVICE_URL", "http://localhost:8090"),
			"auto-bidding-service": getEnv("AUTO_BIDDING_SERVICE_URL", "http://localhost:8091"),
		},

		APIGatewayName:          getEnv("API_GATEWAY_NAME", "api-gateway"),
		AuthServiceName:         getEnv("AUTH_SERVICE_NAME", "auth-service"),
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
