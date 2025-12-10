# Auto-Bidding Service

Service đấu giá tự động cho hệ thống Online Auction.

## 📋 Mô tả

Service này xử lý logic đấu giá tự động, cho phép người dùng đặt giá tối đa và hệ thống sẽ tự động đấu giá thay họ.

### Logic đấu giá tự động

Khi có một bid mới vào sản phẩm:

1. **Tìm tất cả auto-bid đang ACTIVE** của sản phẩm đó
2. **Sắp xếp theo max_amount giảm dần** (nếu bằng nhau thì người tạo trước win)
3. **Xử lý từng auto-bid:**
   - Những người có `max_amount < giá cao nhất hiện tại` → Bid hết `max_amount` của họ
   - Người có `max_amount` cao nhất → Bid cao hơn người thứ 2 đúng **một bước giá**

### Ví dụ

**Sản phẩm:** iPhone 11  
**Giá khởi điểm:** 10,000,000 VNĐ  
**Bước giá:** 100,000 VNĐ

| Bidder | Giá tối đa    | Giá vào sản phẩm | Người giữ giá |
|--------|---------------|------------------|---------------|
| #1     | 11,000,000    | 10,000,000       | #1            |
| #2     | 10,800,000    | 10,800,000       | #1            |
| #3     | 11,500,000    | 11,100,000       | #3            |
| #4     | 11,500,000    | 11,500,000       | #3            |
| #4     | 11,700,000    | 11,600,000       | #4            |

**Giải thích:**
- Bidder #1 đặt max 11tr → Bid 10tr (giá khởi điểm)
- Bidder #2 đặt max 10.8tr → Hệ thống tự động bid 10.8tr cho #2, nhưng #1 có max cao hơn nên tự động bid lại thành 10.9tr → #1 giữ giá
- Bidder #3 đặt max 11.5tr → #1 max 11tr < 11.5tr nên bid hết 11tr, #3 chỉ cần bid 11.1tr (cao hơn #1 một bước) → #3 giữ giá
- Bidder #4 cũng đặt max 11.5tr nhưng #3 đặt trước nên #3 win → #4 phải bid 11.5tr để vượt #3
- Bidder #4 tăng max lên 11.7tr → Bid 11.6tr (cao hơn #3 một bước) → #4 giữ giá

## 🚀 API Endpoints

### 1. Tạo Auto-Bid
```
POST /api/auto-bids
Headers: X-User-ID, X-User-Token
Body: {
  "product_id": 1,
  "max_amount": 15000000
}
```

### 2. Trigger Auto-Bidding (Internal)
```
POST /api/auto-bids/trigger
Headers: X-Internal-Key
Body: {
  "product_id": 1,
  "current_price": 11000000,
  "bid_increment": 100000,
  "new_bidder_id": 5,
  "new_bid_amount": 11000000
}
```

### 3. Lấy Auto-Bids của user
```
GET /api/auto-bids/my
Headers: X-User-ID
```

### 4. Lấy Auto-Bid theo ID
```
GET /api/auto-bids/:id
Headers: X-User-ID
```

### 5. Hủy Auto-Bid
```
POST /api/auto-bids/:id/cancel
Headers: X-User-ID
```

## 🗄️ Database Schema

```sql
CREATE TABLE auto_bids (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    bidder_id BIGINT NOT NULL,
    max_amount DOUBLE PRECISION NOT NULL,
    current_amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_auto_bids_product_id ON auto_bids(product_id);
CREATE INDEX idx_auto_bids_bidder_id ON auto_bids(bidder_id);
CREATE INDEX idx_auto_bids_product_status ON auto_bids(product_id, status);
CREATE INDEX idx_auto_bids_max_amount ON auto_bids(max_amount DESC);
```

## 🛠️ Tech Stack

- **Language:** Go 1.21
- **Framework:** Fiber v2
- **Database:** PostgreSQL (via go-pg)
- **Observability:** OpenTelemetry
- **API Docs:** Swagger

## 📦 Installation

```bash
# Install dependencies
go mod download

# Run service
go run cmd/main.go
```

## 🔧 Environment Variables

```env
DB_HOST=ep-morning-snow-a4t3v7lk-pooler.us-east-1.aws.neon.tech
DB_PORT=5432
DB_USER=neondb_owner
DB_PASSWORD=npg_5DwaV1nZgEor
DB_NAME=neondb
JWT_SECRET=your-super-secret-jwt-key-change-in-production
PORT=3002
BIDDING_SERVICE_URL=http://localhost:8082
PRODUCT_SERVICE_URL=http://localhost:8081
OTEL_ENDPOINT=localhost:4317
OTEL_SERVICE_NAME=auto-bidding-service
OTEL_SERVICE_VERSION=1.0.0
OTEL_ENVIRONMENT=development
```

## 🔄 Flow hoạt động

1. **User tạo auto-bid:**
   - Gọi API `POST /auto-bids` với `max_amount`
   - Service kiểm tra sản phẩm còn hoạt động
   - Kiểm tra `max_amount >= current_price`
   - Lưu auto-bid vào database với status `ACTIVE`

2. **Khi có bid mới từ bidding-service:**
   - Bidding-service gọi `POST /auto-bids/trigger`
   - Service lấy tất cả auto-bid ACTIVE của sản phẩm
   - Sắp xếp theo `max_amount DESC, created_at ASC`
   - Xử lý logic đấu giá tự động
   - Gọi bidding-service để thực hiện bid thực tế

3. **User xem auto-bids:**
   - Gọi API `GET /auto-bids/my`
   - Xem được tất cả auto-bid đã tạo và trạng thái

4. **User hủy auto-bid:**
   - Gọi API `POST /auto-bids/:id/cancel`
   - Cập nhật status thành `CANCELLED`

## 📊 Status của Auto-Bid

- `ACTIVE`: Đang hoạt động
- `WON`: Đã thắng đấu giá
- `OUTBID`: Bị đấu giá vượt quá (max_amount < giá hiện tại)
- `CANCELLED`: Đã hủy bởi user
- `EXPIRED`: Hết hạn (sản phẩm kết thúc đấu giá)

## 🔐 Security

- JWT Authentication qua header `X-User-Token`
- Internal API bảo vệ bằng `X-Internal-Key`
- Chỉ cho phép user thao tác trên auto-bid của chính mình

## 📝 Swagger Documentation

Truy cập: `http://localhost:3002/swagger/`

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## 🐳 Docker

```bash
# Build
docker build -t auto-bidding-service .

# Run
docker run -p 3002:3002 --env-file .env auto-bidding-service
```

## 📞 Contact

- **Service Port:** 3002
- **Health Check:** `GET /api/health`
- **Metrics:** `GET /metrics`
