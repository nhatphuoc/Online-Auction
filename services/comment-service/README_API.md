# Base URL (qua API Gateway)
```
http://<api-gateway-host>/api/comments
```

# Comment Service API

## Tổng quan
Comment Service cung cấp các API cho hệ thống bình luận sản phẩm, bao gồm WebSocket cho realtime comment và REST API cho lịch sử bình luận. Tất cả các REST API đều truy cập qua API Gateway (mặc định: http://localhost:8080).

---

## 1. Lấy thông tin truy cập WebSocket qua API Gateway
- **Endpoint:** `/api/comments/websocket/`
- **Method:** GET
- **Headers:**
  - `X-User-Token: <JWT>` (bắt buộc)
- **Request:**
  - Không cần body, chỉ cần header và truy cập đúng URL
- **Response:**
```json
{
  "comment_service_websocket_url": "ws://localhost:8091/ws",
  "internal_jwt": "<internal-jwt-token>"
}
```

---

## 2. Kết nối WebSocket tới Comment Service (trực tiếp, không qua Gateway)
- **URL:** Sử dụng `comment_service_websocket_url` nhận được ở bước trên
- **Query Params:**
  - `productId`: ID sản phẩm
  - `X-User-Token`: JWT của user
  - `X-Internal-JWT`: internal JWT nhận từ gateway
- **Ví dụ:**
```
ws://localhost:8091/ws?productId=1&X-User-Token=<JWT>&X-Internal-JWT=<internal-jwt-token>
```
- **Gửi bình luận:**
```json
{
  "type": "comment",
  "product_id": 1,
  "content": "Nội dung bình luận"
}
```
- **Gửi typing indicator:**
```json
{
  "type": "typing",
  "product_id": 1
}
```
- **Nhận bình luận:**
```json
{
  "type": "comment",
  "product_id": 1,
  "data": {
    "product_id": 1,
    "sender_id": 123,
    "content": "Nội dung bình luận",
    "created_at": "2025-12-23T10:00:00Z"
  }
}
```
- **Nhận typing indicator:**
```json
{
  "type": "typing",
  "product_id": 1,
  "data": {
    "userId": 123
  }
}
```

---

## 3. Lấy lịch sử bình luận (qua API Gateway)
- **Endpoint:** `/api/comments/history/products/:productId`
- **Method:** GET
- **Headers:**
  - `X-User-Token: <JWT>` (bắt buộc)
- **Response:**
```json
[
  {
    "product_id": 1,
    "sender_id": 123,
    "content": "Nội dung bình luận",
    "created_at": "2025-12-23T10:00:00Z"
  },
  ...
]
```

---

## 4. Health Check (có thể gọi trực tiếp hoặc qua Gateway)
- **Endpoint:** `/health`
- **Method:** GET
- **Response:**
```json
{
  "status": "ok"
}
```

---

## Lưu ý
- Tất cả các REST API đều phải gọi qua API Gateway (http://localhost:8080).
- Chỉ endpoint WebSocket là kết nối trực tiếp tới comment service.
- Luôn truyền đúng các query param khi kết nối WebSocket.
- Response lỗi dạng: `{ "error": "..." }`

---

## Quy trình FE tích hợp
1. Gọi `GET /comments/websocket/` qua API Gateway để lấy URL và internal JWT.
2. Kết nối WebSocket tới comment service với các query param.
3. Gửi/nhận bình luận realtime qua WebSocket.
4. Lấy lịch sử bình luận qua REST API (qua Gateway).
