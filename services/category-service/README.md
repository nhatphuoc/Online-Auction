# Category Service

Service quản lý danh mục sản phẩm và hiển thị danh sách sản phẩm theo danh mục cho hệ thống đấu giá trực tuyến.

## Tính năng

### Quản lý danh mục (Category)
- **CRUD danh mục**: Tạo, đọc, cập nhật, xóa danh mục
- **Cấu trúc phân cấp 2 cấp**: 
  - Cấp 1: Danh mục cha (VD: Điện tử, Thời trang)
  - Cấp 2: Danh mục con (VD: Điện thoại di động, Máy tính xách tay)
- **Lấy danh mục theo cấp**: Lấy tất cả danh mục cấp 1 hoặc cấp 2
- **Lấy danh mục con**: Lấy tất cả danh mục con của một danh mục cha
- **Hiển thị cây danh mục**: Hiển thị danh mục theo cấu trúc phân cấp

### Quản lý sản phẩm (Product)
- **Lấy sản phẩm theo danh mục**: Lấy danh sách sản phẩm của một danh mục (bao gồm cả danh mục con)
- **Phân trang**: Hỗ trợ phân trang cho danh sách sản phẩm
- **Lọc theo trạng thái**: Lọc sản phẩm theo trạng thái (ACTIVE, PENDING, FINISHED, REJECTED)
- **Sắp xếp**: Sắp xếp theo giá, thời gian tạo, thời gian kết thúc

## API Endpoints

### Category APIs

#### Lấy tất cả danh mục (có cấu trúc phân cấp)
```
GET /api/categories
Query params:
  - parent_id: ID của danh mục cha (để lọc danh mục con)
  - level: Cấp độ danh mục (1 hoặc 2)
```

#### Lấy danh mục theo ID
```
GET /api/categories/:id
```

#### Lấy danh mục con của một danh mục cha
```
GET /api/categories/parent/:parent_id
```

#### Tạo danh mục mới (cần authentication)
```
POST /api/categories
Body:
{
  "name": "Điện tử",
  "slug": "dien-tu",
  "description": "Các sản phẩm điện tử",
  "parent_id": null,
  "display_order": 1
}
```

#### Cập nhật danh mục (cần authentication)
```
PUT /api/categories/:id
Body:
{
  "name": "Điện tử",
  "description": "Mô tả mới",
  "is_active": true
}
```

#### Xóa danh mục (soft delete, cần authentication)
```
DELETE /api/categories/:id
```

### Product APIs

#### Lấy sản phẩm theo danh mục
```
GET /api/products?category_id=1&page=1&page_size=20&status=ACTIVE&sort_by=created_at&sort_order=desc
Query params:
  - category_id: ID danh mục (bắt buộc)
  - page: Trang hiện tại (default: 1)
  - page_size: Số sản phẩm mỗi trang (default: 20, max: 100)
  - status: Trạng thái sản phẩm (ACTIVE, PENDING, FINISHED, REJECTED)
  - sort_by: Trường sắp xếp (created_at, current_price, end_at, name)
  - sort_order: Thứ tự sắp xếp (asc, desc)
```

#### Lấy chi tiết sản phẩm
```
GET /api/products/:id
```

## Cấu trúc Database

### Table: categories
```sql
- id: BIGSERIAL PRIMARY KEY
- name: VARCHAR(255) NOT NULL
- slug: VARCHAR(255) NOT NULL UNIQUE
- description: TEXT
- parent_id: BIGINT (FOREIGN KEY -> categories.id)
- level: INT NOT NULL DEFAULT 1
- is_active: BOOLEAN NOT NULL DEFAULT true
- display_order: INT DEFAULT 0
- created_at: TIMESTAMP DEFAULT CURRENT_TIMESTAMP
- updated_at: TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

### Table: products
```sql
- id: BIGSERIAL PRIMARY KEY
- name: VARCHAR(255) NOT NULL
- description: TEXT
- category_id: BIGINT NOT NULL (FOREIGN KEY -> categories.id)
- seller_id: BIGINT NOT NULL
- starting_price: DOUBLE PRECISION NOT NULL
- current_price: DOUBLE PRECISION
- buy_now_price: DOUBLE PRECISION
- step_price: DOUBLE PRECISION NOT NULL
- status: VARCHAR(255) NOT NULL (CHECK: ACTIVE, FINISHED, PENDING, REJECTED)
- thumbnail_url: TEXT
- auto_extend: BOOLEAN NOT NULL DEFAULT false
- end_at: TIMESTAMP NOT NULL
- created_at: TIMESTAMP DEFAULT CURRENT_TIMESTAMP
```

### Table: product_images
```sql
- product_id: BIGINT NOT NULL (FOREIGN KEY -> products.id)
- image_url: VARCHAR(255) NOT NULL
```

## Ví dụ cấu trúc danh mục

```
Điện tử (level 1)
├── Điện thoại di động (level 2)
└── Máy tính xách tay (level 2)

Thời trang (level 1)
├── Giày (level 2)
└── Đồng hồ (level 2)
```

## Cài đặt và chạy

### Prerequisites
- Go 1.25.4+
- PostgreSQL database

### Cấu hình
Tạo file `.env` với nội dung:
```
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=your-db-user
DB_PASSWORD=your-db-password
DB_NAME=your-db-name
JWT_SECRET=your-jwt-secret
PORT=3000
OTEL_ENDPOINT=localhost:4317
OTEL_SERVICE_NAME=category-service-api
OTEL_SERVICE_VERSION=1.0.0
OTEL_ENVIRONMENT=development
```

### Chạy service
```bash
# Cài đặt dependencies
go mod download

# Chạy service
go run cmd/main.go
```

Service sẽ chạy tại `http://localhost:3000`

### Swagger Documentation
Truy cập `http://localhost:3000/swagger/` để xem API documentation

## Authentication
Các endpoint tạo, cập nhật, xóa danh mục yêu cầu JWT token trong header:
```
Authorization: Bearer <your-jwt-token>
```

## Technologies
- **Fiber v2**: Web framework
- **go-pg**: PostgreSQL ORM
- **Swagger**: API documentation
- **OpenTelemetry**: Observability (tracing, metrics, logging)
- **JWT**: Authentication
