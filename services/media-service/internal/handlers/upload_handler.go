package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"media_service/internal/config"
	"media_service/internal/models"
	"media_service/internal/utils"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	s3Client *s3.Client
	cfg      *config.Config
}

func NewUploadHandler(s3Client *s3.Client, cfg *config.Config) *UploadHandler {
	return &UploadHandler{
		s3Client: s3Client,
		cfg:      cfg,
	}
}

// UploadSingleFile godoc
// @Summary Upload a single file
// @Description Upload a file to AWS S3
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param folder query string false "Folder path in S3 (default: uploads/)"
// @Success 200 {object} models.UploadResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /upload [post]
func (h *UploadHandler) UploadSingleFile(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Cần gửi file với key 'file'",
			Details: err.Error(),
		})
	}

	// Validate file size
	if fileHeader.Size > h.cfg.MaxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: fmt.Sprintf("File quá lớn, tối đa %dMB", h.cfg.MaxFileSize/(1024*1024)),
		})
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Không mở được file",
			Details: err.Error(),
		})
	}
	defer file.Close()

	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Lỗi đọc file",
			Details: err.Error(),
		})
	}

	// Get folder from query
	folder := c.Query("folder", "uploads/")
	if folder != "" && !strings.HasSuffix(folder, "/") {
		folder += "/"
	}

	// Generate unique filename
	filename := utils.GenerateUniqueFilename(fileHeader.Filename)
	key := folder + filename

	// Upload to S3
	uploader := manager.NewUploader(h.s3Client)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(h.cfg.AWSBucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(utils.GetContentType(fileHeader.Filename)),
		// ACL removed - use bucket policy instead for public access
	})

	if err != nil {
		slog.Error("Failed to upload file to S3", "error", err, "filename", fileHeader.Filename)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Upload lên S3 thất bại",
			Details: err.Error(),
		})
	}

	url := config.GetS3URL(h.cfg, key)

	slog.Info("File uploaded successfully", "filename", fileHeader.Filename, "key", key, "size", fileHeader.Size)

	return c.JSON(models.UploadResponse{
		Message:    "Upload thành công",
		URL:        url,
		Key:        key,
		Filename:   fileHeader.Filename,
		Size:       fileHeader.Size,
		UploadedAt: time.Now(),
	})
}

// UploadMultipleFiles godoc
// @Summary Upload multiple files
// @Description Upload multiple files to AWS S3
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Files to upload" multiple
// @Param folder query string false "Folder path in S3 (default: uploads/)"
// @Success 200 {object} models.MultipleUploadResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /upload/multiple [post]
func (h *UploadHandler) UploadMultipleFiles(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid multipart form",
			Details: err.Error(),
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Cần gửi ít nhất 1 file với key 'files'",
		})
	}

	// Limit number of files
	if len(files) > h.cfg.MaxFilesPerUpload {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: fmt.Sprintf("Tối đa %d files mỗi lần upload", h.cfg.MaxFilesPerUpload),
		})
	}

	folder := c.Query("folder", "uploads/")
	if folder != "" && !strings.HasSuffix(folder, "/") {
		folder += "/"
	}

	var uploadedFiles []models.UploadResponse
	var failedFiles []models.FailedUpload

	for _, fileHeader := range files {
		// Validate file size
		if fileHeader.Size > h.cfg.MaxFileSize {
			failedFiles = append(failedFiles, models.FailedUpload{
				Filename: fileHeader.Filename,
				Error:    fmt.Sprintf("File quá lớn (max %dMB)", h.cfg.MaxFileSize/(1024*1024)),
			})
			continue
		}

		file, err := fileHeader.Open()
		if err != nil {
			failedFiles = append(failedFiles, models.FailedUpload{
				Filename: fileHeader.Filename,
				Error:    "Không mở được file",
			})
			continue
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			failedFiles = append(failedFiles, models.FailedUpload{
				Filename: fileHeader.Filename,
				Error:    "Lỗi đọc file",
			})
			continue
		}

		filename := utils.GenerateUniqueFilename(fileHeader.Filename)
		key := folder + filename

		uploader := manager.NewUploader(h.s3Client)
		_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket:      aws.String(h.cfg.AWSBucketName),
			Key:         aws.String(key),
			Body:        bytes.NewReader(data),
			ContentType: aws.String(utils.GetContentType(fileHeader.Filename)),
			// ACL removed - use bucket policy instead for public access
		})

		if err != nil {
			slog.Error("Failed to upload file to S3", "error", err, "filename", fileHeader.Filename)
			failedFiles = append(failedFiles, models.FailedUpload{
				Filename: fileHeader.Filename,
				Error:    "Upload lên S3 thất bại: " + err.Error(),
			})
			continue
		}

		url := config.GetS3URL(h.cfg, key)
		uploadedFiles = append(uploadedFiles, models.UploadResponse{
			Message:    "Upload thành công",
			URL:        url,
			Key:        key,
			Filename:   fileHeader.Filename,
			Size:       fileHeader.Size,
			UploadedAt: time.Now(),
		})

		slog.Info("File uploaded successfully", "filename", fileHeader.Filename, "key", key)
	}

	return c.JSON(models.MultipleUploadResponse{
		Message:      fmt.Sprintf("Uploaded %d/%d files successfully", len(uploadedFiles), len(files)),
		Uploaded:     uploadedFiles,
		Failed:       failedFiles,
		Total:        len(files),
		SuccessCount: len(uploadedFiles),
		FailedCount:  len(failedFiles),
	})
}

// Place new handler methods at the end of the file to avoid breaking package structure
// GetPresignedURL godoc
// @Summary Get presigned URL for single file upload
// @Description Get a presigned URL to upload a file directly to S3
// @Tags media
// @Accept json
// @Produce json
// @Param filename query string true "Tên file muốn upload"
// @Param folder query string false "Folder path in S3 (default: uploads/)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /presign [get]
func (h *UploadHandler) GetPresignedURL(c *fiber.Ctx) error {
	filename := c.Query("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Thiếu tên file (filename)",
		})
	}
	folder := c.Query("folder", "uploads/")
	if folder != "" && !strings.HasSuffix(folder, "/") {
		folder += "/"
	}
	key := folder + utils.GenerateUniqueFilename(filename)

	presignClient := s3.NewPresignClient(h.s3Client)
	presignParams := &s3.PutObjectInput{
		Bucket:      aws.String(h.cfg.AWSBucketName),
		Key:         aws.String(key),
		ContentType: aws.String(utils.GetContentType(filename)),
		// ACL removed - use bucket policy instead for public access
	}
	presignDuration := 15 * time.Minute
	presigned, err := presignClient.PresignPutObject(context.TODO(), presignParams, func(opts *s3.PresignOptions) {
		opts.Expires = presignDuration
	})
	if err != nil {
		slog.Error("Failed to generate presigned URL", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Không tạo được presigned URL",
			Details: err.Error(),
		})
	}

	// Generate the final image URL
	imageURL := config.GetS3URL(h.cfg, key)

	return c.JSON(fiber.Map{
		"presigned_url": presigned.URL,
		"image_url":     imageURL,
		"key":           key,
		"expires_in":    int(presignDuration.Seconds()),
	})
}

// GetPresignedURLs godoc
// @Summary Get presigned URLs for multiple files
// @Description Get presigned URLs to upload multiple files directly to S3
// @Tags media
// @Accept json
// @Produce json
// @Param filenames body []string true "Danh sách tên file muốn upload"
// @Param folder query string false "Folder path in S3 (default: uploads/)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /presign/multiple [post]
func (h *UploadHandler) GetPresignedURLs(c *fiber.Ctx) error {
	var filenames []string
	if err := c.BodyParser(&filenames); err != nil || len(filenames) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Body phải là mảng tên file ([]string)",
		})
	}
	folder := c.Query("folder", "uploads/")
	if folder != "" && !strings.HasSuffix(folder, "/") {
		folder += "/"
	}
	presignClient := s3.NewPresignClient(h.s3Client)
	presignDuration := 15 * time.Minute
	result := make([]map[string]interface{}, 0, len(filenames))
	for _, filename := range filenames {
		key := folder + utils.GenerateUniqueFilename(filename)
		presignParams := &s3.PutObjectInput{
			Bucket:      aws.String(h.cfg.AWSBucketName),
			Key:         aws.String(key),
			ContentType: aws.String(utils.GetContentType(filename)),
			// ACL removed - use bucket policy instead for public access
		}
		presigned, err := presignClient.PresignPutObject(context.TODO(), presignParams, func(opts *s3.PresignOptions) {
			opts.Expires = presignDuration
		})
		if err != nil {
			slog.Error("Failed to generate presigned URL", "error", err, "filename", filename)
			result = append(result, map[string]interface{}{
				"filename": filename,
				"error":    err.Error(),
			})
			continue
		}

		// Generate the final image URL
		imageURL := config.GetS3URL(h.cfg, key)

		result = append(result, map[string]interface{}{
			"filename":      filename,
			"presigned_url": presigned.URL,
			"image_url":     imageURL,
			"key":           key,
			"expires_in":    int(presignDuration.Seconds()),
		})
	}
	return c.JSON(fiber.Map{
		"presigned": result,
	})
}
