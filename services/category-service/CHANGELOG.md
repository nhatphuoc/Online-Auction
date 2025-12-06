# Thay đổi Category Service

## Tổng quan
Đã chuyển đổi hoàn toàn từ comment-service sang category-service với đầy đủ chức năng quản lý danh mục và hiển thị sản phẩm theo danh mục.

## Các thay đổi chính

### 1. Models (internal/models/)
- ✅ **Tạo mới `category.go`**: Model cho danh mục với cấu trúc phân cấp 2 cấp
  - Category model với parent_id, level, display_order
  - CreateCategoryRequest, UpdateCategoryRequest
  - CategoryResponse, CategoryTreeResponse
  
- ✅ **Tạo mới `product.go`**: Model cho sản phẩm
  - Product model với đầy đủ thông tin đấu giá
  - ProductImage model cho nhiều ảnh
  - ProductListResponse với pagination
  - ProductQueryParams cho filtering

- ❌ **Xóa `message.go`**: Model cũ của comment service

### 2. Handlers (internal/handlers/)
- ✅ **Tạo mới `category_handler.go`**: Handler đầy đủ cho CRUD categories
  - CreateCategory (POST /categories) - Tạo danh mục mới
  - GetCategories (GET /categories) - Lấy danh sách với filter
  - GetCategoryByID (GET /categories/:id) - Chi tiết danh mục
  - GetCategoriesByParent (GET /categories/parent/:parent_id) - Lấy danh mục con
  - UpdateCategory (PUT /categories/:id) - Cập nhật danh mục
  - DeleteCategory (DELETE /categories/:id) - Xóa danh mục (soft delete)
  - Helper functions: buildCategoryTree, toCategoryResponse

- ✅ **Tạo mới `product_handler.go`**: Handler cho products
  - GetProductsByCategory (GET /products?category_id=X) - Lấy sản phẩm theo danh mục
    - Hỗ trợ pagination (page, page_size)
    - Hỗ trợ filtering (status)
    - Hỗ trợ sorting (sort_by, sort_order)
    - Tự động bao gồm sản phẩm của danh mục con
  - GetProductByID (GET /products/:id) - Chi tiết sản phẩm
  - Helper function: getCategoryWithChildren

- ❌ **Xóa `chat_handler.go`**: Handler cũ của comment service

### 3. Database Schema (internal/config/database.go)
- ✅ **Cập nhật InitSchema()**: Tạo bảng mới
  - Table `categories`: Danh mục với parent-child relationship
    - Indexes: parent_id, level
    - Constraint: Foreign key tự reference
  - Table `products`: Sản phẩm đấu giá
    - Indexes: category_id, status
    - Constraint: Status check, foreign key to categories
  - Table `product_images`: Ảnh sản phẩm
    - Index: product_id
    - Constraint: Foreign key to products with CASCADE delete

- ❌ **Xóa**: Bảng comments và related tables

### 4. Main Application (cmd/main.go)
- ✅ **Cập nhật routes**:
  ```
  /api/categories
    GET / - Lấy tất cả danh mục
    GET /:id - Chi tiết danh mục
    GET /parent/:parent_id - Danh mục con
    POST / (auth) - Tạo danh mục
    PUT /:id (auth) - Cập nhật
    DELETE /:id (auth) - Xóa
  
  /api/products
    GET / - Lấy sản phẩm theo danh mục
    GET /:id - Chi tiết sản phẩm
  ```

- ❌ **Xóa**: WebSocket endpoints, comment routes

- ✅ **Cập nhật Swagger**: Đổi title và description phù hợp

### 5. Dependencies (go.mod)
- ❌ **Xóa**: `github.com/gofiber/websocket/v2`
- ❌ **Xóa**: `github.com/fasthttp/websocket`
- ✅ **Giữ lại**: Các dependencies cần thiết (fiber, go-pg, otel, jwt, validator)

### 6. Configuration
- ✅ **Cập nhật .env**: OTEL_SERVICE_NAME = "category-service-api"
- ✅ **Dockerfile**: Build từ cmd/main.go, expose port 3000
- ✅ **docker-compose.yml**: Service mới với đúng tên và config

### 7. Documentation & Scripts
- ✅ **README.md**: Mô tả đầy đủ về category service
  - Tính năng
  - API endpoints
  - Database schema
  - Ví dụ cấu trúc danh mục
  - Hướng dẫn cài đặt và chạy

- ✅ **API_EXAMPLES.md**: Ví dụ chi tiết cho tất cả endpoints
  - cURL commands
  - Request/Response samples
  - Query parameters
  - Authentication guide

- ✅ **Makefile**: Commands tiện lợi
  - build, run, test, clean
  - swagger, docker-build, docker-run
  - tidy, fmt, lint
  - seed (chạy seed data)

- ✅ **scripts/seed.go**: Script seed dữ liệu mẫu
  - 3 categories cấp 1: Điện tử, Thời trang, Gia dụng
  - 6 categories cấp 2: Điện thoại, Laptop, Tablet, Giày, Đồng hồ, Túi xách

- ✅ **.gitignore**: Ignore files phù hợp với Go project

### 8. Xóa các files/folders không cần thiết
- ❌ **Xóa `public/`**: Static files của comment service
- ❌ **Xóa `internal/models/message.go`**: Model cũ
- ❌ **Xóa `internal/handlers/chat_handler.go`**: Handler cũ

## Tính năng chính

### Quản lý Danh mục
1. **Cấu trúc phân cấp 2 cấp**:
   - Level 1: Danh mục cha (VD: Điện tử, Thời trang)
   - Level 2: Danh mục con (VD: Điện thoại, Giày)
   
2. **CRUD đầy đủ**:
   - Tạo, đọc, cập nhật, xóa (soft delete)
   - Validation đầy vào
   - Không cho xóa danh mục có children
   
3. **Query linh hoạt**:
   - Lấy theo level
   - Lấy theo parent_id
   - Hiển thị dạng tree hoặc flat list

### Quản lý Sản phẩm
1. **Lấy sản phẩm theo danh mục**:
   - Tự động bao gồm sản phẩm của danh mục con
   - VD: Lấy danh mục "Điện tử" → bao gồm cả "Điện thoại", "Laptop"

2. **Pagination**:
   - page, page_size
   - total, total_pages

3. **Filtering & Sorting**:
   - Filter theo status (ACTIVE, PENDING, FINISHED, REJECTED)
   - Sort theo created_at, current_price, end_at, name
   - Sort order: asc/desc

4. **Relations**:
   - Include category information
   - Include product images

## API Structure

```
/api
├── /categories
│   ├── GET / (list all with tree structure)
│   ├── GET /:id (get one)
│   ├── GET /parent/:parent_id (get children)
│   ├── POST / (create - auth required)
│   ├── PUT /:id (update - auth required)
│   └── DELETE /:id (soft delete - auth required)
└── /products
    ├── GET / (list by category with pagination)
    └── GET /:id (get one)
```

## Database Tables

### categories
- Hierarchical structure với parent_id
- 2 levels maximum
- Soft delete với is_active flag
- Display order cho sorting

### products
- Full auction information
- Foreign key to categories
- Status constraint
- Indexes cho performance

### product_images
- One-to-many với products
- Cascade delete khi xóa product

## Testing

### Chạy service
```bash
make run
# hoặc
go run cmd/main.go
```

### Seed dữ liệu
```bash
make seed
# hoặc
go run scripts/seed.go
```

### Build
```bash
make build
```

### Swagger docs
```bash
make swagger
```

Sau khi chạy, truy cập: http://localhost:3000/swagger/

## Lưu ý
- Service yêu cầu PostgreSQL database
- Các endpoint tạo/sửa/xóa category cần JWT authentication
- Product endpoints chỉ READ, không có CREATE/UPDATE/DELETE (dành cho product-service)
- OpenTelemetry được tích hợp cho observability
