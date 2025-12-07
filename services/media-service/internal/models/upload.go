package models

import "time"

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

// MultipleUploadResponse represents response for multiple file uploads
type MultipleUploadResponse struct {
	Message      string           `json:"message" example:"Uploaded 8/10 files successfully"`
	Uploaded     []UploadResponse `json:"uploaded"`
	Failed       []FailedUpload   `json:"failed"`
	Total        int              `json:"total" example:"10"`
	SuccessCount int              `json:"success_count" example:"8"`
	FailedCount  int              `json:"failed_count" example:"2"`
}

// FailedUpload represents a failed file upload
type FailedUpload struct {
	Filename string `json:"filename" example:"large_file.jpg"`
	Error    string `json:"error" example:"File quá lớn (max 50MB)"`
}
