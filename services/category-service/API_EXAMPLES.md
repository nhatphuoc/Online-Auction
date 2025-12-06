# API Examples - Category Service

## Base URL
```
http://localhost:3000/api
```

## 1. Lấy tất cả danh mục (cấu trúc phân cấp)

### Request
```bash
curl -X GET "http://localhost:3000/api/categories"
```

### Response
```json
{
  "categories": [
    {
      "id": 1,
      "name": "Điện tử",
      "slug": "dien-tu",
      "description": "Các sản phẩm điện tử",
      "level": 1,
      "is_active": true,
      "display_order": 1,
      "created_at": "2025-12-06T10:00:00Z",
      "updated_at": "2025-12-06T10:00:00Z",
      "children": [
        {
          "id": 4,
          "name": "Điện thoại di động",
          "slug": "dien-thoai-di-dong",
          "description": "Điện thoại thông minh, điện thoại cơ bản",
          "parent_id": 1,
          "level": 2,
          "is_active": true,
          "display_order": 1,
          "created_at": "2025-12-06T10:00:00Z",
          "updated_at": "2025-12-06T10:00:00Z"
        },
        {
          "id": 5,
          "name": "Máy tính xách tay",
          "slug": "may-tinh-xach-tay",
          "description": "Laptop, notebook",
          "parent_id": 1,
          "level": 2,
          "is_active": true,
          "display_order": 2,
          "created_at": "2025-12-06T10:00:00Z",
          "updated_at": "2025-12-06T10:00:00Z"
        }
      ]
    },
    {
      "id": 2,
      "name": "Thời trang",
      "slug": "thoi-trang",
      "description": "Các sản phẩm thời trang",
      "level": 1,
      "is_active": true,
      "display_order": 2,
      "created_at": "2025-12-06T10:00:00Z",
      "updated_at": "2025-12-06T10:00:00Z",
      "children": [
        {
          "id": 6,
          "name": "Giày",
          "slug": "giay",
          "description": "Giày thể thao, giày tây, giày sneaker",
          "parent_id": 2,
          "level": 2,
          "is_active": true,
          "display_order": 1,
          "created_at": "2025-12-06T10:00:00Z",
          "updated_at": "2025-12-06T10:00:00Z"
        },
        {
          "id": 7,
          "name": "Đồng hồ",
          "slug": "dong-ho",
          "description": "Đồng hồ đeo tay nam, nữ",
          "parent_id": 2,
          "level": 2,
          "is_active": true,
          "display_order": 2,
          "created_at": "2025-12-06T10:00:00Z",
          "updated_at": "2025-12-06T10:00:00Z"
        }
      ]
    }
  ]
}
```

## 2. Lấy danh mục theo cấp (level)

### Lấy danh mục cấp 1
```bash
curl -X GET "http://localhost:3000/api/categories?level=1"
```

### Lấy danh mục cấp 2
```bash
curl -X GET "http://localhost:3000/api/categories?level=2"
```

## 3. Lấy danh mục con của "Điện tử"

### Request
```bash
curl -X GET "http://localhost:3000/api/categories/parent/1"
```

### Response
```json
[
  {
    "id": 4,
    "name": "Điện thoại di động",
    "slug": "dien-thoai-di-dong",
    "description": "Điện thoại thông minh, điện thoại cơ bản",
    "parent_id": 1,
    "level": 2,
    "is_active": true,
    "display_order": 1,
    "created_at": "2025-12-06T10:00:00Z",
    "updated_at": "2025-12-06T10:00:00Z"
  },
  {
    "id": 5,
    "name": "Máy tính xách tay",
    "slug": "may-tinh-xach-tay",
    "description": "Laptop, notebook",
    "parent_id": 1,
    "level": 2,
    "is_active": true,
    "display_order": 2,
    "created_at": "2025-12-06T10:00:00Z",
    "updated_at": "2025-12-06T10:00:00Z"
  }
]
```

## 4. Tạo danh mục mới

### Request (cần JWT token)
```bash
curl -X POST "http://localhost:3000/api/categories" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Sách",
    "slug": "sach",
    "description": "Sách các loại",
    "parent_id": null,
    "display_order": 4
  }'
```

### Response
```json
{
  "id": 8,
  "name": "Sách",
  "slug": "sach",
  "description": "Sách các loại",
  "level": 1,
  "is_active": true,
  "display_order": 4,
  "created_at": "2025-12-06T10:00:00Z",
  "updated_at": "2025-12-06T10:00:00Z"
}
```

## 5. Tạo danh mục con

### Request (cần JWT token)
```bash
curl -X POST "http://localhost:3000/api/categories" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Tai nghe",
    "slug": "tai-nghe",
    "description": "Tai nghe không dây, có dây",
    "parent_id": 1,
    "display_order": 4
  }'
```

## 6. Cập nhật danh mục

### Request (cần JWT token)
```bash
curl -X PUT "http://localhost:3000/api/categories/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Điện tử & Công nghệ",
    "description": "Các sản phẩm điện tử và công nghệ hiện đại"
  }'
```

## 7. Xóa danh mục

### Request (cần JWT token)
```bash
curl -X DELETE "http://localhost:3000/api/categories/8" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Response
```json
{
  "message": "Category deleted successfully"
}
```

## 8. Lấy sản phẩm theo danh mục

### Request - Lấy sản phẩm của danh mục "Điện tử" (ID: 1)
```bash
curl -X GET "http://localhost:3000/api/products?category_id=1&page=1&page_size=20&status=ACTIVE&sort_by=created_at&sort_order=desc"
```

### Response
```json
{
  "products": [
    {
      "id": 1,
      "name": "iPhone 15 Pro Max",
      "description": "Điện thoại iPhone 15 Pro Max 256GB",
      "category_id": 4,
      "seller_id": 123,
      "starting_price": 25000000,
      "current_price": 26000000,
      "buy_now_price": 30000000,
      "step_price": 500000,
      "status": "ACTIVE",
      "thumbnail_url": "https://example.com/iphone15.jpg",
      "auto_extend": false,
      "end_at": "2025-12-10T18:00:00Z",
      "created_at": "2025-12-06T10:00:00Z",
      "category": {
        "id": 4,
        "name": "Điện thoại di động",
        "slug": "dien-thoai-di-dong",
        "description": "Điện thoại thông minh, điện thoại cơ bản",
        "parent_id": 1,
        "level": 2,
        "is_active": true,
        "display_order": 1,
        "created_at": "2025-12-06T10:00:00Z",
        "updated_at": "2025-12-06T10:00:00Z"
      },
      "images": [
        {
          "product_id": 1,
          "image_url": "https://example.com/iphone15-1.jpg"
        },
        {
          "product_id": 1,
          "image_url": "https://example.com/iphone15-2.jpg"
        }
      ]
    }
  ],
  "total": 50,
  "page": 1,
  "page_size": 20,
  "total_pages": 3
}
```

### Các query parameters khả dụng:
- `category_id` (bắt buộc): ID của danh mục
- `page` (mặc định: 1): Số trang
- `page_size` (mặc định: 20, max: 100): Số sản phẩm mỗi trang
- `status`: Lọc theo trạng thái (ACTIVE, PENDING, FINISHED, REJECTED)
- `sort_by`: Sắp xếp theo trường (created_at, current_price, end_at, name)
- `sort_order`: Thứ tự sắp xếp (asc, desc)

## 9. Lấy chi tiết sản phẩm

### Request
```bash
curl -X GET "http://localhost:3000/api/products/1"
```

### Response
```json
{
  "id": 1,
  "name": "iPhone 15 Pro Max",
  "description": "Điện thoại iPhone 15 Pro Max 256GB, màu Titan tự nhiên",
  "category_id": 4,
  "seller_id": 123,
  "starting_price": 25000000,
  "current_price": 26000000,
  "buy_now_price": 30000000,
  "step_price": 500000,
  "status": "ACTIVE",
  "thumbnail_url": "https://example.com/iphone15.jpg",
  "auto_extend": false,
  "end_at": "2025-12-10T18:00:00Z",
  "created_at": "2025-12-06T10:00:00Z",
  "category": {
    "id": 4,
    "name": "Điện thoại di động",
    "slug": "dien-thoai-di-dong",
    "description": "Điện thoại thông minh, điện thoại cơ bản",
    "parent_id": 1,
    "level": 2,
    "is_active": true,
    "display_order": 1,
    "created_at": "2025-12-06T10:00:00Z",
    "updated_at": "2025-12-06T10:00:00Z"
  },
  "images": [
    {
      "product_id": 1,
      "image_url": "https://example.com/iphone15-1.jpg"
    },
    {
      "product_id": 1,
      "image_url": "https://example.com/iphone15-2.jpg"
    }
  ]
}
```

## 10. Health Check

### Request
```bash
curl -X GET "http://localhost:3000/health"
```

### Response
```json
{
  "status": "ok"
}
```

## Lưu ý về Authentication

Các endpoint cần authentication (tạo, cập nhật, xóa danh mục) yêu cầu JWT token trong header:

```
Authorization: Bearer <your-jwt-token>
```

JWT token có thể lấy từ auth-service sau khi đăng nhập thành công.
