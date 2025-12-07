# Media Service

Service quản lý upload và lưu trữ media files (images, videos, documents) lên AWS S3 cho hệ thống đấu giá trực tuyến.

## Tính năng

### Upload Media
- **Upload đơn file**: Upload một file lên S3
- **Upload nhiều files**: Upload tối đa 10 files cùng lúc
- **Tự động generate unique filename**: Tránh trùng lặp tên file
- **Validate file size**: Giới hạn max 50MB/file
- **Hỗ trợ nhiều loại file**: Images, videos, documents

### Supported File Types
- **Images**: jpg, jpeg, png, gif, webp, svg, bmp, ico
- **Videos**: mp4, mov, avi, wmv, flv, webm
- **Audio**: mp3, wav
- **Documents**: pdf, doc, docx, xls, xlsx, ppt, pptx, txt, csv
- **Archives**: zip, rar, 7z

## API Endpoints

### Upload single file
```
POST /api/upload
Content-Type: multipart/form-data

Parameters:
- file: File to upload (required)
- folder: Folder path in S3 (optional, default: "uploads/")
```

**Example with cURL:**
```bash
curl -X POST "http://localhost:3000/api/upload?folder=products/" \
  -F "file=@image.jpg"
```

**Response:**
```json
{
  "message": "Upload thành công",
  "url": "https://bucket.s3.region.amazonaws.com/products/image_1234567890.jpg",
  "key": "products/image_1234567890.jpg",
  "filename": "image.jpg",
  "size": 1048576,
  "uploaded_at": "2025-12-07T10:00:00Z"
}
```

### Upload multiple files
```
POST /api/upload/multiple
Content-Type: multipart/form-data

Parameters:
- files: Files to upload (required, max 10 files)
- folder: Folder path in S3 (optional, default: "uploads/")
```

**Example with cURL:**
```bash
curl -X POST "http://localhost:3000/api/upload/multiple?folder=gallery/" \
  -F "files=@image1.jpg" \
  -F "files=@image2.jpg" \
  -F "files=@image3.jpg"
```

**Response:**
```json
{
  "message": "Uploaded 2/3 files successfully",
  "uploaded": [
    {
      "message": "Upload thành công",
      "url": "https://bucket.s3.region.amazonaws.com/gallery/image1_1234567890.jpg",
      "key": "gallery/image1_1234567890.jpg",
      "filename": "image1.jpg",
      "size": 1048576,
      "uploaded_at": "2025-12-07T10:00:00Z"
    },
    {
      "message": "Upload thành công",
      "url": "https://bucket.s3.region.amazonaws.com/gallery/image2_1234567891.jpg",
      "key": "gallery/image2_1234567891.jpg",
      "filename": "image2.jpg",
      "size": 2097152,
      "uploaded_at": "2025-12-07T10:00:01Z"
    }
  ],
  "failed": [
    {
      "filename": "image3.jpg",
      "error": "File quá lớn (max 50MB)"
    }
  ],
  "total": 3,
  "success_count": 2,
  "failed_count": 1
}
```

### Health check
```
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "service": "media-service",
  "bucket": "my-bucket",
  "region": "ap-southeast-1"
}
```

## Configuration

### Environment Variables
Create a `.env` file with the following variables:

```env
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
AWS_REGION=ap-southeast-1
AWS_BUCKET_NAME=your_bucket_name
PORT=3000
```

### AWS S3 Setup
1. Create an S3 bucket
2. Configure bucket permissions for public read access (if needed)
3. Create IAM user with S3 write permissions
4. Get Access Key ID and Secret Access Key

## Installation & Running

### Prerequisites
- Go 1.25.4+
- AWS S3 account
- AWS credentials with S3 permissions

### Run locally
```bash
# Install dependencies
go mod download

# Run service
go run cmd/main.go
```

Service will start on `http://localhost:3000`

### Swagger Documentation
Access API documentation at: `http://localhost:3000/swagger/`

### Build
```bash
go build -o bin/media-service cmd/main.go
```

### Docker
```bash
# Build image
docker build -t media-service .

# Run container
docker run -p 3000:3000 --env-file .env media-service
```

## Project Structure

```
media-service/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   ├── config.go        # Configuration management
│   │   ├── errors.go        # Custom errors
│   │   └── s3.go            # S3 client initialization
│   ├── handlers/
│   │   └── upload_handler.go # Upload handlers
│   ├── models/
│   │   └── upload.go        # Data models
│   └── utils/
│       └── file.go          # File utilities
├── docs/                    # Swagger documentation
├── .env                     # Environment variables
├── Dockerfile              # Docker configuration
├── go.mod                  # Go module dependencies
└── README.md               # This file
```

## Features

### File Validation
- **Size validation**: Max 50MB per file
- **Multiple files limit**: Max 10 files per upload
- **Content-Type detection**: Automatic MIME type detection

### Unique Filenames
- Automatically appends timestamp to avoid filename conflicts
- Original filename is preserved in response
- Format: `{original_name}_{timestamp}.{extension}`

### Error Handling
- Graceful error handling for failed uploads
- Detailed error messages
- Continues uploading other files on multi-upload failure

### Logging
- Structured JSON logging with slog
- Request/response logging
- Upload success/failure tracking

## Technologies
- **Fiber v2**: Web framework
- **AWS SDK Go v2**: S3 client
- **Swagger**: API documentation
- **slog**: Structured logging

## Security Considerations
- AWS credentials should be kept secret
- Use IAM roles with minimal required permissions
- Consider using pre-signed URLs for sensitive files
- Implement rate limiting for production
- Add authentication/authorization middleware if needed

## Limitations
- Max file size: 50MB
- Max files per multi-upload: 10
- Files are stored with public-read ACL by default

## Future Improvements
- [ ] Image resizing/optimization
- [ ] Video thumbnail generation
- [ ] File type restrictions by configuration
- [ ] Virus scanning before upload
- [ ] CDN integration
- [ ] Delete file endpoint
- [ ] List files endpoint
- [ ] Private file upload with pre-signed URLs
