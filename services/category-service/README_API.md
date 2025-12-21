# Category Service API

## Base URL (qua API Gateway)
```
http://<api-gateway-host>/api/categories
```

## 1. Lấy tất cả danh mục (GET)
- **Endpoint:** `/api/categories/`
- **Method:** GET
- **Headers:**
  - X-User-Token: <token> (bắt buộc)
- **Request:**
  - Query params (optional):
    - `parent_id`: int (lọc theo parent)
    - `level`: int (lọc theo level)
- **Response:**
```json
{
  "categories": [
    {
      "id": 1,
      "name": "string",
      "slug": "string",
      "description": "string",
      "parent_id": null,
      "level": 1,
      "is_active": true,
      "display_order": 0,
      "created_at": "2025-12-22T10:00:00Z",
      "updated_at": "2025-12-22T10:00:00Z",
      "children": [ ... ]
    }
  ]
}
```

## 2. Lấy danh mục theo parent (GET)
- **Endpoint:** `/api/categories/parent/{parent_id}`
- **Method:** GET
- **Headers:**
  - X-User-Token: <token>
- **Response:**
```json
[
  {
    "id": 2,
    "name": "string",
    "slug": "string",
    "description": "string",
    "parent_id": 1,
    "level": 2,
    "is_active": true,
    "display_order": 0,
    "created_at": "2025-12-22T10:00:00Z",
    "updated_at": "2025-12-22T10:00:00Z",
    "children": []
  }
]
```

## 3. Lấy chi tiết danh mục (GET)
- **Endpoint:** `/api/categories/{id}`
- **Method:** GET
- **Headers:**
  - X-User-Token: <token>
- **Response:**
```json
{
  "id": 1,
  "name": "string",
  "slug": "string",
  "description": "string",
  "parent_id": null,
  "level": 1,
  "is_active": true,
  "display_order": 0,
  "created_at": "2025-12-22T10:00:00Z",
  "updated_at": "2025-12-22T10:00:00Z",
  "children": [ ... ]
}
```

## 4. Tạo danh mục (POST)
- **Endpoint:** `/api/categories/`
- **Method:** POST
- **Headers:**
  - X-User-Token: <token>
- **Request:**
```json
{
  "name": "string",
  "slug": "string",
  "description": "string",
  "parent_id": 1,
  "display_order": 0
}
```
- **Response:**
```json
{
  "id": 3,
  "name": "string",
  "slug": "string",
  "description": "string",
  "parent_id": 1,
  "level": 2,
  "is_active": true,
  "display_order": 0,
  "created_at": "2025-12-22T10:00:00Z",
  "updated_at": "2025-12-22T10:00:00Z"
}
```

## 5. Cập nhật danh mục (PUT)
- **Endpoint:** `/api/categories/{id}`
- **Method:** PUT
- **Headers:**
  - X-User-Token: <token>
- **Request:**
```json
{
  "name": "string",
  "slug": "string",
  "description": "string",
  "parent_id": 1,
  "is_active": true,
  "display_order": 0
}
```
- **Response:**
```json
{
  "id": 3,
  "name": "string",
  "slug": "string",
  "description": "string",
  "parent_id": 1,
  "level": 2,
  "is_active": true,
  "display_order": 0,
  "created_at": "2025-12-22T10:00:00Z",
  "updated_at": "2025-12-22T10:00:00Z"
}
```

## 6. Xóa danh mục (DELETE)
- **Endpoint:** `/api/categories/{id}`
- **Method:** DELETE
- **Headers:**
  - X-User-Token: <token>
- **Response:**
```json
{
  "message": "Category deleted successfully"
}
```

## Lưu ý
- Tất cả request đều phải có header X-User-Token.
- Response lỗi dạng: `{ "error": "..." }`
