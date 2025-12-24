# Media Service API

## Base URL (qua API Gateway)
```
http://<api-gateway-host>/api/media
```

---

- **Endpoint:** `/api/media/upload`
- **Method:** POST
- **Headers:**
  - X-User-Token: <JWT> (bắt buộc)
  - Content-Type: multipart/form-data
- **Body:**
  - file: file cần upload (key: `file`)
  - folder: (query, optional) thư mục trên S3 (default: uploads/)
- **Response:**
```json
{
  "message": "Upload thành công",
  "url": "https://...",
  "key": "uploads/abc123.jpg",
  "filename": "abc.jpg",
  "size": 12345,
  "uploaded_at": "2025-12-24T10:00:00Z"
}
```

---

- **Endpoint:** `/api/media/upload/multiple`
- **Method:** POST
- **Headers:**
  - X-User-Token: <JWT> (bắt buộc)
  - Content-Type: multipart/form-data
- **Body:**
  - files: các file cần upload (key: `files`, type: file[], multiple)
  - folder: (query, optional) thư mục trên S3 (default: uploads/)
- **Response:**
```json
{
  "message": "Uploaded 2/2 files successfully",
  "uploaded": [
    { "message": "Upload thành công", "url": "...", "key": "...", "filename": "...", "size": 123, "uploaded_at": "..." }
  ],
  "failed": [
    { "filename": "...", "error": "..." }
  ],
  "total": 2,
  "success_count": 2,
  "failed_count": 0
}
```

---

- **Endpoint:** `/api/media/presign`
- **Method:** GET
- **Headers:**
  - X-User-Token: <JWT> (bắt buộc)
- **Query:**
  - filename: tên file muốn upload (bắt buộc)
  - folder: (optional) thư mục trên S3 (default: uploads/)
- **Response:**
```json
{
  "presigned_url": "https://s3.amazonaws.com/...",
  "image_url": "https://your-bucket.s3.region.amazonaws.com/uploads/abc123.jpg",
  "key": "uploads/abc123.jpg",
  "expires_in": 900
}
```
- **Cách sử dụng:**
  1. Gọi API này để lấy `presigned_url` và `image_url`
  2. Dùng `presigned_url` để upload file trực tiếp lên S3 (PUT request với file content)
  3. Sau khi upload thành công, lưu `image_url` vào database để hiển thị ảnh

---

- **Endpoint:** `/api/media/presign/multiple`
- **Method:** POST
- **Headers:**
  - X-User-Token: <JWT> (bắt buộc)
  - Content-Type: application/json
- **Body:**
  - JSON array tên file: ["a.jpg", "b.png"]
  - folder: (query, optional) thư mục trên S3 (default: uploads/)
- **Response:**
```json
{
  "presigned": [
    { 
      "filename": "a.jpg", 
      "presigned_url": "https://s3.amazonaws.com/...", 
      "image_url": "https://your-bucket.s3.region.amazonaws.com/uploads/a.jpg",
      "key": "uploads/a.jpg", 
      "expires_in": 900 
    },
    { 
      "filename": "b.png", 
      "presigned_url": "https://s3.amazonaws.com/...", 
      "image_url": "https://your-bucket.s3.region.amazonaws.com/uploads/b.png",
      "key": "uploads/b.png", 
      "expires_in": 900 
    }
  ]
}
```
- **Cách sử dụng:**
  1. Gọi API này với danh sách tên file
  2. Nhận về mảng chứa `presigned_url` và `image_url` cho từng file
  3. Upload từng file lên S3 bằng `presigned_url` tương ứng
  4. Lưu các `image_url` vào database

---

- **Endpoint:** `/api/media/health`
- **Method:** GET
- **Headers:**
  - X-User-Token: <JWT> (bắt buộc)
- **Response:**
```json
{
  "status": "ok",
  "service": "media-service",
  "bucket": "...",
  "region": "..."
}
```

---

## Lưu ý
- Các API upload file sử dụng multipart/form-data.
- Các API presign trả về URL để client upload trực tiếp lên S3.
- Tất cả các API đều đi qua API Gateway (trừ upload trực tiếp lên S3 bằng presigned URL).
- Tất cả các API (trừ upload trực tiếp lên S3) đều yêu cầu header `X-User-Token` hợp lệ.
- Response lỗi dạng: `{ "error": "..." }`
