package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// GetContentType returns the MIME type based on file extension
func GetContentType(filename string) string {
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
	case ".bmp":
		return "image/bmp"
	case ".ico":
		return "image/x-icon"
	case ".mp4":
		return "video/mp4"
	case ".mov":
		return "video/quicktime"
	case ".avi":
		return "video/x-msvideo"
	case ".wmv":
		return "video/x-ms-wmv"
	case ".flv":
		return "video/x-flv"
	case ".webm":
		return "video/webm"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".txt":
		return "text/plain"
	case ".csv":
		return "text/csv"
	case ".zip":
		return "application/zip"
	case ".rar":
		return "application/x-rar-compressed"
	case ".7z":
		return "application/x-7z-compressed"
	default:
		return "application/octet-stream"
	}
}

// GenerateUniqueFilename generates a unique filename with timestamp
func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	nameWithoutExt := strings.TrimSuffix(originalFilename, ext)
	// Clean filename: remove special characters
	nameWithoutExt = cleanFilename(nameWithoutExt)
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%s_%d%s", nameWithoutExt, timestamp, ext)
}

// cleanFilename removes special characters from filename
func cleanFilename(name string) string {
	// Replace spaces with underscore
	name = strings.ReplaceAll(name, " ", "_")
	// Remove or replace other special characters if needed
	return name
}

// ValidateFileExtension checks if file extension is allowed
func ValidateFileExtension(filename string, allowedExtensions []string) bool {
	if len(allowedExtensions) == 0 {
		return true // No restriction
	}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowed := range allowedExtensions {
		if ext == strings.ToLower(allowed) {
			return true
		}
	}
	return false
}
