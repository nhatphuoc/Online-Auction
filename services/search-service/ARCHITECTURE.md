# Search Service - Full-text Search với Elasticsearch + PostgreSQL + Redis Stream

## Tổng quan Kiến trúc

Search service được thiết kế theo mô hình **Event-Driven Architecture** với các components chính:

```
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│  PostgreSQL  │────────▶│ Redis Stream │────────▶│ Search Worker│
│  (Source of  │         │  (Message    │         │  (Consumer)  │
│   Truth)     │         │   Queue)     │         └──────┬───────┘
└──────────────┘         └──────────────┘                │
                                                          │
                                                          ▼
                                                  ┌──────────────┐
                                                  │Elasticsearch │
                                                  │  (Index +    │
                                                  │   Search)    │
                                                  └──────────────┘
                                                          ▲
                                                          │
                                                  ┌───────┴──────┐
                                                  │ Search API   │
                                                  │   (Fiber)    │
                                                  └──────────────┘
```

## Các Components

### 1. **PostgreSQL** - Source of Truth
- Lưu trữ dữ liệu sản phẩm và danh mục
- Schema: `product` và `category` với các trường như đã mô tả

### 2. **Redis Stream** - Message Queue
- Nhận events từ các services khác (product-service, category-service)
- Event types: `product.created`, `product.updated`, `product.deleted`, `category.created`, `category.updated`, `category.deleted`
- Consumer Group: `search_service_group`

### 3. **Worker** - Event Processor
- Lắng nghe Redis Stream
- Đọc dữ liệu từ PostgreSQL
- Index/Update/Delete documents trong Elasticsearch

### 4. **Elasticsearch** - Search Engine
- **Vietnamese Analyzer**: Sử dụng `asciifolding` filter để xử lý tiếng Việt không dấu
- **Indexing**: 
  - `products` index với fields: `name`, `name_no_accent`, `description`, `description_no_accent`, etc.
  - `categories` index
- **Search Features**:
  - Full-text search với multi_match query
  - Filter theo category, status, price range
  - Sort theo price, end_at, created_at
  - Function score để boost sản phẩm mới (trong N phút)

### 5. **Search API** - REST Endpoints
- `GET /api/search/products` - Tìm kiếm sản phẩm

## Luồng hoạt động

### Khi tạo/cập nhật sản phẩm:
1. Product Service ghi vào PostgreSQL
2. Sau commit → Publish event `product.updated` vào Redis Stream
3. Search Worker nhận event
4. Worker đọc dữ liệu mới nhất từ PostgreSQL (bao gồm cả category)
5. Worker index vào Elasticsearch với:
   - Convert `name` → `name_no_accent` (tiếng Việt không dấu)
   - Convert `description` → `description_no_accent`
   - Thêm thông tin category

### Khi xóa sản phẩm:
1. Product Service xóa trong PostgreSQL
2. Publish event `product.deleted`
3. Worker xóa document trong Elasticsearch

### Khi tìm kiếm:
1. Client gửi request đến `/api/search/products`
2. Search API build Elasticsearch query:
   - **multi_match** trên `name`, `name_no_accent`, `description`, `description_no_accent`
   - **filter** theo category_id, status, price range
   - **function_score** với gauss decay để boost sản phẩm mới
   - **sort** theo yêu cầu
3. Trả về kết quả với pagination

## Elasticsearch Mapping

### Product Index
```json
{
  "settings": {
    "analysis": {
      "analyzer": {
        "vietnamese_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "asciifolding"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "name": {"type": "text"},
      "name_no_accent": {"type": "text", "analyzer": "vietnamese_analyzer"},
      "description": {"type": "text"},
      "description_no_accent": {"type": "text", "analyzer": "vietnamese_analyzer"},
      "category_id": {"type": "long"},
      "status": {"type": "keyword"},
      "current_price": {"type": "double"},
      "end_at": {"type": "date"},
      "created_at": {"type": "date"}
    }
  }
}
```

## Search Query với Boost

```json
{
  "query": {
    "function_score": {
      "query": {
        "bool": {
          "must": [
            {
              "multi_match": {
                "query": "điện thoại",
                "fields": ["name^3", "name_no_accent^3", "description", "description_no_accent"],
                "fuzziness": "AUTO"
              }
            }
          ],
          "filter": [
            {"term": {"status": "ACTIVE"}},
            {"range": {"current_price": {"gte": 1000000, "lte": 10000000}}}
          ]
        }
      },
      "functions": [
        {
          "gauss": {
            "created_at": {
              "origin": "now",
              "scale": "60m",
              "decay": 0.5
            }
          },
          "weight": 2.0
        }
      ],
      "score_mode": "sum",
      "boost_mode": "multiply"
    }
  },
  "sort": [{"end_at": {"order": "desc"}}]
}
```

## API Endpoints

### Search Products
```
GET /api/search/products?query=điện thoại&category_id=1&status=ACTIVE&min_price=1000000&max_price=10000000&sort_by=price&sort_order=asc&page=1&page_size=20
```

**Query Parameters:**
- `query`: Từ khóa tìm kiếm (tên hoặc mô tả sản phẩm)
- `category_id`: Lọc theo danh mục
- `status`: Lọc theo trạng thái (ACTIVE, PENDING, FINISHED, REJECTED)
- `min_price`, `max_price`: Lọc theo khoảng giá
- `sort_by`: Sắp xếp theo (price, end_at, created_at)
- `sort_order`: Thứ tự (asc, desc)
- `page`, `page_size`: Phân trang

**Response:**
```json
{
  "products": [...],
  "total": 100,
  "page": 1,
  "page_size": 20,
  "total_pages": 5
}
```

## Environment Variables

```env
# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=auction_db

# Elasticsearch
ELASTICSEARCH_URL=http://localhost:9200
ELASTICSEARCH_INDEX_PRODUCT=products
ELASTICSEARCH_INDEX_CATEGORY=categories

# Redis Stream
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_STREAM_KEY=auction_events
REDIS_CONSUMER_GROUP=search_service_group
REDIS_CONSUMER_NAME=search_service_consumer_1

# Search Configuration
BOOST_MINUTES=60
BOOST_SCORE=2.0

# Server
PORT=3000
```

## Cách chạy

1. **Install dependencies:**
```bash
go mod download
```

2. **Start dependencies (Docker):**
```bash
docker-compose up -d postgres elasticsearch redis
```

3. **Run service:**
```bash
go run cmd/main.go
```

## Testing

### 1. Test Health Check
```bash
curl http://localhost:3000/health
```

### 2. Test Search
```bash
curl "http://localhost:3000/api/search/products?query=dien+thoai&page=1&page_size=10"
```

### 3. Test with Filters
```bash
curl "http://localhost:3000/api/search/products?query=laptop&category_id=1&status=ACTIVE&min_price=1000000&max_price=20000000&sort_by=price&sort_order=asc"
```

## Notes

- **Vietnamese Accent Removal**: Sử dụng custom function `RemoveVietnameseAccents()` để convert tiếng Việt có dấu sang không dấu
- **Boost Recent Products**: Sản phẩm mới đăng trong vòng N phút (config `BOOST_MINUTES`) sẽ được ưu tiên hiển thị cao hơn với score nhân `BOOST_SCORE`
- **Graceful Shutdown**: Service hỗ trợ graceful shutdown, đảm bảo worker hoàn thành xử lý event trước khi tắt
- **Error Handling**: Tất cả errors đều được log, events thất bại sẽ được acknowledge để không block queue
