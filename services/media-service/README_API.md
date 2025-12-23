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
  "url": "https://...",
  "key": "uploads/abc123.jpg",
  "expires_in": 900
}
```

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
    { "filename": "a.jpg", "url": "...", "key": "...", "expires_in": 900 },
    { "filename": "b.png", "url": "...", "key": "...", "expires_in": 900 }
  ]
}
```

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
