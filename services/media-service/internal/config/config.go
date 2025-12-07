package config

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		// Không panic, vì có thể dùng environment variables
	}
}

type Config struct {
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSBucketName      string
	Port               string
	MaxFileSize        int64
	MaxFilesPerUpload  int
}

func LoadConfig() *Config {
	return &Config{
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSRegion:          getEnv("AWS_REGION", "ap-southeast-1"),
		AWSBucketName:      getEnv("AWS_BUCKET_NAME", ""),
		Port:               getEnv("PORT", "3000"),
		MaxFileSize:        50 * 1024 * 1024, // 50MB
		MaxFilesPerUpload:  10,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) Validate() error {
	if c.AWSAccessKeyID == "" {
		return ErrMissingAWSAccessKey
	}
	if c.AWSSecretAccessKey == "" {
		return ErrMissingAWSSecretKey
	}
	if c.AWSRegion == "" {
		return ErrMissingAWSRegion
	}
	if c.AWSBucketName == "" {
		return ErrMissingAWSBucket
	}
	return nil
}
