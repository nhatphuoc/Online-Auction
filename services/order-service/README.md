# Order Service

Order Service quản lý quy trình thanh toán và giao hàng sau khi đấu giá kết thúc trong hệ thống Online Auction.

## Chức năng chính

### 1. Quản lý đơn hàng (Order Management)
- **Tạo đơn hàng**: Tự động tạo đơn hàng khi đấu giá kết thúc (được gọi từ auction-service)
- **Xem đơn hàng**: Người mua và người bán có thể xem chi tiết đơn hàng của mình
- **Danh sách đơn hàng**: Xem tất cả đơn hàng với phân trang

### 2. Quy trình thanh toán (Payment Flow)
Order service quản lý 4 bước thanh toán sau đấu giá:

1. **PENDING_PAYMENT**: Chờ người mua thanh toán
2. **PAYMENT_CONFIRMED**: Người mua đã thanh toán và gửi địa chỉ
3. **ADDRESS_PROVIDED**: Người mua đã cung cấp địa chỉ giao hàng
4. **INVOICE_SENT**: Người bán đã gửi hóa đơn vận chuyển
5. **DELIVERED**: Đã giao hàng thành công
6. **COMPLETED**: Giao dịch hoàn tất
7. **CANCELLED**: Đơn hàng bị hủy

### 3. Chat giữa người mua và người bán
- Gửi tin nhắn trong đơn hàng
- Xem lịch sử tin nhắn
- Real-time communication cho việc trao đổi thông tin giao hàng

### 4. Đánh giá (Rating)
- Người mua đánh giá người bán (+1 hoặc -1)
- Người bán đánh giá người mua (+1 hoặc -1)
- Chỉ có thể đánh giá sau khi giao hàng hoặc hoàn tất
- Mỗi bên chỉ đánh giá một lần

### 5. Hủy đơn hàng
- Cả người mua và người bán đều có thể hủy đơn
- Không thể hủy nếu đã hoàn tất
- Người bán có thể hủy và tự động -1 điểm người thắng

## Kiến trúc

```
order-service/
├── cmd/
│   └── main.go                 # Entry point
├── internal/
│   ├── config/                 # Configuration
│   │   ├── config.go
│   │   └── database.go
│   ├── handlers/               # HTTP handlers
│   │   └── order_handler.go
│   ├── middleware/             # Middleware
│   │   └── auth.go
│   ├── models/                 # Data models
│   │   └── order.go
│   ├── repository/             # Database layer
│   │   └── order_repository.go
│   ├── service/                # Business logic
│   │   └── order_service.go
│   └── utils/                  # Utilities
│       └── validator.go
└── docs/                       # Swagger documentation

```

## API Endpoints

### Orders
- `POST /api/orders` - Tạo đơn hàng mới (internal)
- `GET /api/orders` - Lấy danh sách đơn hàng của user (requires auth)
- `GET /api/orders/:id` - Xem chi tiết đơn hàng (requires auth)
- `PATCH /api/orders/:id/status` - Cập nhật trạng thái đơn hàng (requires auth)
- `POST /api/orders/:id/cancel` - Hủy đơn hàng (requires auth)

### Messages
- `POST /api/orders/:id/messages` - Gửi tin nhắn (requires auth)
- `GET /api/orders/:id/messages` - Xem tin nhắn (requires auth)

### Ratings
- `POST /api/orders/:id/rate` - Đánh giá đơn hàng (requires auth)
- `GET /api/orders/:id/rating` - Xem đánh giá

## Authentication

Order service **không xử lý authentication** trực tiếp. Thay vào đó:

- API Gateway sẽ validate JWT token từ auth-service
- API Gateway forward thông tin user qua headers:
  - `X-User-ID`: ID của user
  - `X-User-Email`: Email của user
  - `X-User-Role`: Role của user (BUYER, SELLER, ADMIN)
- Order service chỉ extract thông tin từ headers này

## Database Schema

### Table: orders
```sql
- id (BIGSERIAL PRIMARY KEY)
- auction_id (BIGINT) - ID của sản phẩm đấu giá
- winner_id (BIGINT) - Người thắng (buyer)
- seller_id (BIGINT) - Người bán
- final_price (DOUBLE PRECISION) - Giá cuối cùng
- status (VARCHAR) - Trạng thái đơn hàng
- payment_method (VARCHAR) - Phương thức thanh toán
- payment_proof (TEXT) - Ảnh chứng từ
- shipping_address (TEXT) - Địa chỉ giao hàng
- shipping_phone (VARCHAR) - SĐT nhận hàng
- tracking_number (VARCHAR) - Mã vận đơn
- shipping_invoice (TEXT) - Hóa đơn vận chuyển
- delivered_at (TIMESTAMP)
- completed_at (TIMESTAMP)
- cancelled_at (TIMESTAMP)
- cancel_reason (TEXT)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

### Table: order_messages
```sql
- id (BIGSERIAL PRIMARY KEY)
- order_id (BIGINT) - FK to orders
- sender_id (BIGINT) - Người gửi
- message (TEXT) - Nội dung tin nhắn
- created_at (TIMESTAMP)
```

### Table: order_ratings
```sql
- id (BIGSERIAL PRIMARY KEY)
- order_id (BIGINT UNIQUE) - FK to orders
- buyer_rating (INT) - Đánh giá của buyer (+1/-1)
- buyer_comment (TEXT)
- seller_rating (INT) - Đánh giá của seller (+1/-1)
- seller_comment (TEXT)
- buyer_rated_at (TIMESTAMP)
- seller_rated_at (TIMESTAMP)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

## Environment Variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=online_auction

# Server
PORT=3000
GRPC_PORT=50051

# OpenTelemetry (optional)
OTEL_ENDPOINT=localhost:4317
OTEL_SERVICE_NAME=order-service
OTEL_SERVICE_VERSION=1.0.0
OTEL_ENVIRONMENT=development
```

## Build & Run

### Development
```bash
# Install dependencies
go mod download

# Run service
go run cmd/main.go

# Generate Swagger docs
swag init -g cmd/main.go -o docs
```

### Docker
```bash
# Build image
docker build -t order-service .

# Run container
docker run -p 3000:3000 --env-file .env order-service
```

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Swagger Documentation

Sau khi chạy service, truy cập Swagger UI tại:
```
http://localhost:3000/swagger/
```

## Dependencies

- **Fiber v2**: Web framework
- **go-pg v10**: PostgreSQL ORM
- **validator v10**: Request validation
- **swag**: Swagger documentation generator

## Integration với các service khác

### Auction Service
- Khi đấu giá kết thúc, auction-service gọi `POST /api/orders` để tạo đơn hàng

### User Service
- Lấy thông tin user (buyer/seller) để hiển thị
- Cập nhật điểm đánh giá của user

### Notification Service
- Gửi email khi có thay đổi trạng thái đơn hàng
- Thông báo tin nhắn mới trong chat

### Payment Service
- Xử lý thanh toán (nếu có tích hợp payment gateway)

## Notes

- Order service là stateless, có thể scale horizontal
- Sử dụng PostgreSQL indexes để tối ưu query performance
- Tất cả API cần authentication đều yêu cầu headers từ API Gateway
- Status transition được validate nghiêm ngặt để đảm bảo flow đúng
