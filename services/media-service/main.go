package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	_ "media_service/docs"
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

var s3Client *s3.Client
var bucketName string
var awsRegion string

// UploadResponse represents the response after successful upload
type UploadResponse struct {
	Message    string    `json:"message" example:"Upload thành công"`
	URL        string    `json:"url" example:"https://bucket.s3.region.amazonaws.com/uploads/image.jpg"`
	Key        string    `json:"key" example:"uploads/image.jpg"`
	Filename   string    `json:"filename" example:"image.jpg"`
	Size       int64     `json:"size" example:"1048576"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Details string `json:"details,omitempty" example:"Missing file parameter"`
}

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Không tìm thấy .env, sẽ dùng biến môi trường hệ thống")
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion = os.Getenv("AWS_REGION")
	bucketName = os.Getenv("AWS_BUCKET_NAME")

	if accessKey == "" || secretKey == "" || awsRegion == "" || bucketName == "" {
		log.Fatal("Thiếu cấu hình AWS trong .env")
	}

	// Tạo config AWS SDK v2
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		log.Fatal("Không load được AWS config:", err)
	}

	// Tạo S3 client (v2)
	s3Client = s3.NewFromConfig(cfg)

	// Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024, // 50MB max upload size
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(ErrorResponse{
				Error: err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} | ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API routes
	api := app.Group("/api")

	// Upload endpoints
	api.Post("/upload", uploadHandler)
	api.Post("/upload/multiple", uploadMultipleHandler)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "media-service",
			"bucket":  bucketName,
			"region":  awsRegion,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Media service chạy trên :%s", port)
	log.Printf("Swagger documentation: http://localhost:%s/swagger/", port)
	log.Fatal(app.Listen(":" + port))
}

// uploadHandler godoc
// @Summary Upload a single file
// @Description Upload a file to AWS S3
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param folder query string false "Folder path in S3 (default: uploads/)"
// @Success 200 {object} UploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /upload [post]
func uploadHandler(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Cần gửi file với key 'file'",
			Details: err.Error(),
		})
	}

	// Validate file size (max 50MB)
	if fileHeader.Size > 50*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "File quá lớn, tối đa 50MB",
		})
	}

	// Mở file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{
			Error:   "Không mở được file",
			Details: err.Error(),
		})
	}
	defer file.Close()

	// Đọc toàn bộ nội dung
	data, err := io.ReadAll(file)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{
			Error:   "Lỗi đọc file",
			Details: err.Error(),
		})
	}

	// Lấy folder từ query (tùy chọn)
	folder := c.Query("folder", "uploads/")
	if folder != "" && !strings.HasSuffix(folder, "/") {
		folder += "/"
	}

	// Generate unique filename
	filename := generateUniqueFilename(fileHeader.Filename)
	key := folder + filename

	// Upload to S3
	uploader := manager.NewUploader(s3Client)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(getContentType(fileHeader.Filename)),
	})

	if err != nil {
		return c.Status(500).JSON(ErrorResponse{
			Error:   "Upload lên S3 thất bại",
			Details: err.Error(),
		})
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, awsRegion, key)

	return c.JSON(UploadResponse{
		Message:    "Upload thành công",
		URL:        url,
		Key:        key,
		Filename:   fileHeader.Filename,
		Size:       fileHeader.Size,
		UploadedAt: time.Now(),
	})
}

// uploadMultipleHandler godoc
// @Summary Upload multiple files
// @Description Upload multiple files to AWS S3
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Files to upload" multiple
// @Param folder query string false "Folder path in S3 (default: uploads/)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /upload/multiple [post]
func uploadMultipleHandler(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid multipart form",
			Details: err.Error(),
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Cần gửi ít nhất 1 file với key 'files'",
		})
	}

	// Limit number of files
	if len(files) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Tối đa 10 files mỗi lần upload",
		})
	}

	folder := c.Query("folder", "uploads/")
	if folder != "" && !strings.HasSuffix(folder, "/") {
		folder += "/"
	}

	var uploadedFiles []UploadResponse
	var failedFiles []map[string]string

	for _, fileHeader := range files {
		// Validate file size
		if fileHeader.Size > 50*1024*1024 {
			failedFiles = append(failedFiles, map[string]string{
				"filename": fileHeader.Filename,
				"error":    "File quá lớn (max 50MB)",
			})
			continue
		}

		file, err := fileHeader.Open()
		if err != nil {
			failedFiles = append(failedFiles, map[string]string{
				"filename": fileHeader.Filename,
				"error":    "Không mở được file",
			})
			continue
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			failedFiles = append(failedFiles, map[string]string{
				"filename": fileHeader.Filename,
				"error":    "Lỗi đọc file",
			})
			continue
		}

		filename := generateUniqueFilename(fileHeader.Filename)
		key := folder + filename

		uploader := manager.NewUploader(s3Client)
		_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(key),
			Body:        bytes.NewReader(data),
			ContentType: aws.String(getContentType(fileHeader.Filename)),
			ACL:         "public-read",
		})

		if err != nil {
			failedFiles = append(failedFiles, map[string]string{
				"filename": fileHeader.Filename,
				"error":    "Upload lên S3 thất bại: " + err.Error(),
			})
			continue
		}

		url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, awsRegion, key)
		uploadedFiles = append(uploadedFiles, UploadResponse{
			Message:    "Upload thành công",
			URL:        url,
			Key:        key,
			Filename:   fileHeader.Filename,
			Size:       fileHeader.Size,
			UploadedAt: time.Now(),
		})
	}

	return c.JSON(fiber.Map{
		"message":       fmt.Sprintf("Uploaded %d/%d files successfully", len(uploadedFiles), len(files)),
		"uploaded":      uploadedFiles,
		"failed":        failedFiles,
		"total":         len(files),
		"success_count": len(uploadedFiles),
		"failed_count":  len(failedFiles),
	})
}

// Helper functions
func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".avi":
		return "video/x-msvideo"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}

func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	nameWithoutExt := strings.TrimSuffix(originalFilename, ext)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%s_%d%s", nameWithoutExt, timestamp, ext)
}
