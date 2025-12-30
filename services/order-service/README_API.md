# ORDER SERVICE API DOCUMENTATION

**Base URL (qua API Gateway):** `http://localhost:8080/api/orders`
**Direct URL:** `http://localhost:8086`

**Note:** Tất cả requests phải đi qua API Gateway tại port 8080.

---

## 🔐 Authentication

Tất cả endpoints yêu cầu JWT token trong header:
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

Token được trả về sau khi login thành công qua Auth Service.

---

## 📋 Order Lifecycle

```
1. PENDING_PAYMENT    -> Người mua cần thanh toán
2. PAID               -> Đã thanh toán, chờ địa chỉ
3. ADDRESS_PROVIDED   -> Đã có địa chỉ, chờ seller gửi hàng
4. SHIPPING           -> Đang vận chuyển
5. DELIVERED          -> Đã giao hàng
6. COMPLETED          -> Hoàn thành (sau khi đánh giá)
7. CANCELLED          -> Đã hủy
```

---

## ❤️ WATCH LIST API (Danh sách yêu thích)

### 1. Add Product to Watch List

**POST** `http://localhost:8080/api/orders/data/watchlist`

**Description:** Thêm sản phẩm vào danh sách yêu thích

**Headers:**
```
Content-Type: application/json
X-User-Token: <JWT_TOKEN>
```

**Request Body:**
```json
{
  "product_id": 123
}
```

**Success Response (201 Created):**
```json
{
  "message": "Product added to watch list successfully",
  "data": {
    "id": 1,
    "user_id": 456,
    "product_id": 123,
    "created_at": "2025-12-31T10:30:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request:** Invalid product_id
- **401 Unauthorized:** Missing or invalid token
- **409 Conflict:** Product already in watch list
- **500 Internal Server Error:** Database error

---

### 2. Get User's Watch List

**GET** `http://localhost:8080/api/orders/data/watchlist?page=1&limit=20`

**Description:** Lấy danh sách sản phẩm yêu thích của user

**Headers:**
```
X-User-Token: <JWT_TOKEN>
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)

**Success Response (200 OK):**
```json
{
  "message": "Watch list fetched successfully",
  "data": [
    {
      "id": 1,
      "user_id": 456,
      "product_id": 123,
      "created_at": "2025-12-31T10:30:00Z"
    },
    {
      "id": 2,
      "user_id": 456,
      "product_id": 124,
      "created_at": "2025-12-30T15:20:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 2
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Missing or invalid token
- **500 Internal Server Error:** Database error

---

### 3. Remove Product from Watch List

**DELETE** `http://localhost:8080/api/orders/data/watchlist/{product_id}`

**Description:** Xóa sản phẩm khỏi danh sách yêu thích

**Headers:**
```
X-User-Token: <JWT_TOKEN>
```

**Path Parameters:**
- `product_id`: ID của sản phẩm cần xóa

**Success Response (200 OK):**
```json
{
  "message": "Product removed from watch list successfully"
}
```

**Error Responses:**
- **400 Bad Request:** Invalid product_id
- **401 Unauthorized:** Missing or invalid token
- **404 Not Found:** Product not found in watch list
- **500 Internal Server Error:** Database error

---

### 4. Check if Product is in Watch List

**GET** `http://localhost:8080/api/orders/data/watchlist/{product_id}/check`

**Description:** Kiểm tra xem sản phẩm có trong danh sách yêu thích không

**Headers:**
```
X-User-Token: <JWT_TOKEN>
```

**Path Parameters:**
- `product_id`: ID của sản phẩm cần kiểm tra

**Success Response (200 OK):**
```json
{
  "is_in_watchlist": true,
  "product_id": 123
}
```

**Error Responses:**
- **400 Bad Request:** Invalid product_id
- **401 Unauthorized:** Missing or invalid token
- **500 Internal Server Error:** Database error

---

## 📚 ORDER API Endpoints

### 1. Create Order (Internal - After Auction Ends)

**POST** `http://localhost:8080/api/orders`

**Authorization:** Internal service only

**Headers:**
```
Content-Type: application/json
X-User-Token: <JWT_TOKEN>
```

**Request Body:**
```json
{
  "auction_id": 1,
  "winner_id": 5,
  "seller_id": 3,
  "final_price": 25000000
}
```

**Response (201):**
```json
{
  "id": 1,
  "auction_id": 1,
  "winner_id": 5,
  "seller_id": 3,
  "final_price": 25000000,
  "status": "PENDING_PAYMENT",
  "payment_method": "",
  "payment_proof": "",
  "shipping_address": "",
  "shipping_phone": "",
  "tracking_number": "",
  "shipping_invoice": "",
  "paid_at": null,
  "delivered_at": null,
  "completed_at": null,
  "cancelled_at": null,
  "cancel_reason": "",
  "created_at": "2025-12-30T10:00:00Z",
  "updated_at": "2025-12-30T10:00:00Z"
}
```

---

### 2. Get Order By ID

**GET** `http://localhost:8080/api/orders/{id}`

**Authorization:** ROLE_BIDDER, ROLE_SELLER (chỉ buyer hoặc seller của đơn hàng)

**Headers:**
```
X-User-Token: <JWT_TOKEN>
```

**Response (200):**
```json
{
  "id": 1,
  "auction_id": 1,
  "winner_id": 5,
  "seller_id": 3,
  "final_price": 25000000,
  "status": "PAID",
  "payment_method": "MOMO",
  "payment_proof": "https://s3.amazonaws.com/proof.jpg",
  "shipping_address": "123 Nguyễn Huệ, Q1, TP.HCM",
  "shipping_phone": "0901234567",
  "tracking_number": "",
  "shipping_invoice": "",
  "paid_at": "2025-12-30T10:30:00Z",
  "delivered_at": null,
  "completed_at": null,
  "cancelled_at": null,
  "cancel_reason": "",
  "created_at": "2025-12-30T10:00:00Z",
  "updated_at": "2025-12-30T10:30:00Z",
  "buyer_name": "Nguyễn Văn A",
  "seller_name": "Trần Thị B",
  "rating": {
    "id": 1,
    "order_id": 1,
    "buyer_rating": null,
    "buyer_comment": "",
    "seller_rating": null,
    "seller_comment": "",
    "buyer_rated_at": null,
    "seller_rated_at": null,
    "created_at": "2025-12-30T10:00:00Z",
    "updated_at": "2025-12-30T10:00:00Z"
  }
}
```

**Response (403):**
```json
{
  "error": "You don't have permission to view this order"
}
```

---

### 3. Get User Orders

**GET** `http://localhost:8080/api/orders?role=buyer&status=COMPLETED`

**Authorization:** ROLE_BIDDER, ROLE_SELLER

**Headers:**
```
X-User-Token: <JWT_TOKEN>
```

**Query Parameters:**
- `role` (optional): `buyer` hoặc `seller` - lọc đơn hàng theo vai trò
- `status` (optional): Lọc theo trạng thái (PENDING_PAYMENT, PAID, SHIPPING, COMPLETED, CANCELLED)

**Examples:**
- Lấy tất cả đơn mua: `GET http://localhost:8080/api/orders?role=buyer`
- Lấy đơn bán đã hoàn thành: `GET http://localhost:8080/api/orders?role=seller&status=COMPLETED`
- Lấy tất cả đơn hàng: `GET http://localhost:8080/api/orders`

**Response (200):**
```json
[
  {
    "id": 1,
    "auction_id": 1,
    "winner_id": 5,
    "seller_id": 3,
    "final_price": 25000000,
    "status": "COMPLETED",
    "payment_method": "MOMO",
    "shipping_address": "123 Nguyễn Huệ, Q1, TP.HCM",
    "shipping_phone": "0901234567",
    "tracking_number": "VN123456789",
    "paid_at": "2025-12-30T10:30:00Z",
    "delivered_at": "2025-12-31T14:00:00Z",
    "completed_at": "2025-12-31T15:00:00Z",
    "created_at": "2025-12-30T10:00:00Z",
    "updated_at": "2025-12-31T15:00:00Z",
    "buyer_name": "Nguyễn Văn A",
    "seller_name": "Trần Thị B"
  }
]
```

---

### 4. Pay for Order (Buyer)

**POST** `http://localhost:8080/api/orders/{id}/pay`

**Authorization:** ROLE_BIDDER (chỉ buyer của đơn hàng)

**Headers:**
```
Content-Type: application/json
X-User-Token: <JWT_TOKEN>
```

**Request Body:**
```json
{
  "payment_method": "MOMO",
  "payment_proof": "https://s3.amazonaws.com/payment-proof.jpg"
}
```

**Available payment methods:** MOMO, ZALOPAY, VNPAY, STRIPE, PAYPAL

**Response (200):**
```json
{
  "id": 1,
  "status": "PAID",
  "payment_method": "MOMO",
  "payment_proof": "https://s3.amazonaws.com/payment-proof.jpg",
  "paid_at": "2025-12-30T10:30:00Z",
  "updated_at": "2025-12-30T10:30:00Z"
}
```

---

### 5. Provide Shipping Address (Buyer)

**POST** `http://localhost:8080/api/orders/{id}/shipping-address`

**Authorization:** ROLE_BIDDER (chỉ buyer của đơn hàng)

**Headers:**
```
Content-Type: application/json
X-User-Token: <JWT_TOKEN>
```

**Request Body:**
```json
{
  "shipping_address": "123 Nguyễn Huệ, Quận 1, TP. Hồ Chí Minh",
  "shipping_phone": "0901234567"
}
```

**Response (200):**
```json
{
  "id": 1,
  "status": "ADDRESS_PROVIDED",
  "shipping_address": "123 Nguyễn Huệ, Quận 1, TP. Hồ Chí Minh",
  "shipping_phone": "0901234567",
  "updated_at": "2025-12-30T11:00:00Z"
}
```

---

### 6. Provide Tracking Number (Seller)

**POST** `http://localhost:8080/api/orders/{id}/tracking`

**Authorization:** ROLE_SELLER (chỉ seller của đơn hàng)

**Headers:**
```
Content-Type: application/json
X-User-Token: <JWT_TOKEN>
```

**Request Body:**
```json
{
  "tracking_number": "VN123456789",
  "shipping_invoice": "https://s3.amazonaws.com/invoice.pdf"
}
```

**Response (200):**
```json
{
  "id": 1,
  "status": "SHIPPING",
  "tracking_number": "VN123456789",
  "shipping_invoice": "https://s3.amazonaws.com/invoice.pdf",
  "updated_at": "2025-12-30T12:00:00Z"
}
```

---

### 7. Confirm Delivery (Buyer)

**POST** `http://localhost:8080/api/orders/{id}/confirm-delivery`

**Authorization:** ROLE_BIDDER (chỉ buyer của đơn hàng)

**Headers:**
```
X-User-Token: <JWT_TOKEN>
```

**Response (200):**
```json
{
  "id": 1,
  "status": "DELIVERED",
  "delivered_at": "2025-12-31T14:00:00Z",
  "updated_at": "2025-12-31T14:00:00Z"
}
```

---

### 8. Rate Seller (Buyer)

**POST** `http://localhost:8080/api/orders/{id}/rate-seller`

**Authorization:** ROLE_BIDDER (chỉ buyer của đơn hàng, sau khi DELIVERED)

**Headers:**
```
Content-Type: application/json
X-User-Token: <JWT_TOKEN>
```

**Request Body:**
```json
{
  "rating": 1,
  "comment": "Người bán rất nhiệt tình, giao hàng nhanh!"
}
```

**Rating values:** `1` (positive) hoặc `-1` (negative)

**Response (200):**
```json
{
  "id": 1,
  "order_id": 1,
  "buyer_rating": 1,
  "buyer_comment": "Người bán rất nhiệt tình, giao hàng nhanh!",
  "buyer_rated_at": "2025-12-31T15:00:00Z",
  "seller_rating": null,
  "seller_comment": "",
  "seller_rated_at": null,
  "updated_at": "2025-12-31T15:00:00Z"
}
```

---

### 9. Rate Buyer (Seller)

**POST** `http://localhost:8080/api/orders/{id}/rate-buyer`

**Authorization:** ROLE_SELLER (chỉ seller của đơn hàng, sau khi DELIVERED)

**Headers:**
```
Content-Type: application/json
X-User-Token: <JWT_TOKEN>
```

**Request Body:**
```json
{
  "rating": 1,
  "comment": "Người mua thanh toán nhanh, dễ giao tiếp!"
}
```

**Response (200):**
```json
{
  "id": 1,
  "order_id": 1,
  "buyer_rating": 1,
  "buyer_comment": "Người bán rất nhiệt tình!",
  "buyer_rated_at": "2025-12-31T15:00:00Z",
  "seller_rating": 1,
  "seller_comment": "Người mua thanh toán nhanh, dễ giao tiếp!",
  "seller_rated_at": "2025-12-31T15:30:00Z",
  "updated_at": "2025-12-31T15:30:00Z"
}
```

---

### 10. Cancel Order (Seller)

**POST** `http://localhost:8080/api/orders/{id}/cancel`

**Authorization:** ROLE_SELLER (chỉ seller của đơn hàng)

**Headers:**
```
Content-Type: application/json
X-User-Token: <JWT_TOKEN>
```

**Request Body:**
```json
{
  "reason": "Người thắng không thanh toán trong 24h"
}
```

**Response (200):**
```json
{
  "id": 1,
  "status": "CANCELLED",
  "cancel_reason": "Người thắng không thanh toán trong 24h",
  "cancelled_at": "2025-12-30T18:00:00Z",
  "updated_at": "2025-12-30T18:00:00Z"
}
```

**Note:** Khi hủy đơn, seller tự động rate buyer -1 với comment là lý do hủy.

---

### 11. Get Chat History (Messages)

**GET** `http://localhost:8080/api/orders/data/product/{id}/messages`

**Authorization:** ROLE_BIDDER, ROLE_SELLER (buyer hoặc seller của đơn hàng)

**Headers:**
```
X-User-Token: <JWT_TOKEN>
```

**Query Parameters:**
- `limit` (optional, default: 50, max: 100): Số lượng messages mỗi trang
- `offset` (optional, default: 0): Vị trí bắt đầu (dùng cho pagination)

**Example:** 
```
GET http://localhost:8080/api/orders/data/product/1/messages?limit=20&offset=0
```

**Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "order_id": 1,
      "sender_id": 5,
      "message": "Xin chào, khi nào bạn giao hàng?",
      "created_at": "2025-12-30T13:00:00Z"
    },
    {
      "id": 2,
      "order_id": 1,
      "sender_id": 3,
      "message": "Mình sẽ gửi hàng ngày mai bạn nhé!",
      "created_at": "2025-12-30T13:05:00Z"
    }
  ],
  "pagination": {
    "total": 2,
    "limit": 50,
    "offset": 0
  }
}
```

**Note:** 
- Messages được sắp xếp theo thứ tự thời gian (cũ nhất đến mới nhất)
- Endpoint này dùng để load lịch sử chat ban đầu
- Để nhận real-time messages, sử dụng WebSocket connection

---

### 12. WebSocket Connection for Real-time Chat

**WebSocket URL:** `ws://localhost:8086/ws`

**Connection Parameters:**
- `orderId`: ID của order
- `X-User-Token`: JWT Access Token
- `X-Internal-JWT`: Internal JWT token (lấy từ API Gateway endpoint `/api/order-websocket/`)

**Cách connect:**

1. **Lấy WebSocket info từ API Gateway:**
```
GET http://localhost:8080/api/order-websocket/
Headers:
  X-User-Token: <JWT_TOKEN>

Response:
{
  "order_service_websocket_url": "ws://localhost:8086/ws",
  "internal_jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

2. **Connect WebSocket:**
```javascript
const ws = new WebSocket(
  'ws://localhost:8086/ws?orderId=1&X-User-Token=<JWT>&X-Internal-JWT=<internal_jwt>'
);
```

**WebSocket Message Types:**

**Send Message:**
```json
{
  "type": "message",
  "content": "Xin chào!"
}
```

**Send Typing Indicator:**
```json
{
  "type": "typing"
}
```

**Receive Message:**
```json
{
  "type": "message",
  "order_id": 1,
  "data": {
    "id": 3,
    "order_id": 1,
    "sender_id": 5,
    "message": "Xin chào!",
    "created_at": "2025-12-30T14:00:00Z"
  }
}
```

**Receive Typing Indicator:**
```json
{
  "type": "typing",
  "order_id": 1,
  "data": {
    "userId": 5
  }
}
```

**Best Practice:**
- Sử dụng REST API `/messages` để load lịch sử chat khi mở trang
- Sử dụng WebSocket để nhận và gửi messages real-time
- Kết nối WebSocket khi user đang xem order detail page
- Disconnect WebSocket khi rời khỏi trang

---

### 13. Rate Order (Rate Seller)

## 🔄 Workflow Example

### Complete Order Flow (Buyer Perspective):

1. **Auction ends** → Order created automatically với status `PENDING_PAYMENT`
2. **Buyer pays** → `POST /orders/{id}/pay` → Status: `PAID`
3. **Buyer provides address** → `POST /orders/{id}/shipping-address` → Status: `ADDRESS_PROVIDED`
4. **Seller ships** → `POST /orders/{id}/tracking` → Status: `SHIPPING`
5. **Buyer confirms** → `POST /orders/{id}/confirm-delivery` → Status: `DELIVERED`
6. **Buyer rates seller** → `POST /orders/{id}/rate-seller` → Status: `COMPLETED`
7. **Seller rates buyer** → `POST /orders/{id}/rate-buyer`

### Cancel Flow:

- **Seller cancels** → `POST /orders/{id}/cancel` → Status: `CANCELLED`, auto rate buyer -1

---

## 🎯 User Roles

- **ROLE_BIDDER**: Người mua (winner of auction)
  - Can: pay, provide address, confirm delivery, rate seller
  - Can view: own orders as buyer

- **ROLE_SELLER**: Người bán
  - Can: provide tracking, rate buyer, cancel order
  - Can view: own orders as seller

- **ROLE_ADMIN**: Quản trị viên
  - Can view: all orders (future feature)

---

## ⚠️ Error Responses

**400 Bad Request:**
```json
{
  "error": "Invalid request body"
}
```

**401 Unauthorized:**
```json
{
  "error": "Missing user information"
}
```

**403 Forbidden:**
```json
{
  "error": "You don't have permission to perform this action"
}
```

**404 Not Found:**
```json
{
  "error": "Order not found"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Internal server error"
}
```

---

## 📝 Changelog

### v2.0 - Rating & User Info Improvements

**Fixed:**
- ✅ **Rating API** - `POST /orders/{id}/rate` now automatically creates rating record if not exists (UPSERT logic)
- ✅ **Cancel Order** - `POST /orders/{id}/cancel` now creates rating record before applying negative rating
- ✅ **User Names** - All order responses now include `buyer_name` and `seller_name` from JOIN with users table

**Response Changes:**
All order responses now include:
```json
{
  "buyer_name": "Nguyễn Văn A",
  "seller_name": "Trần Thị B"
}
```

**Endpoints affected:**
- `GET /orders/{id}` - Single order with user names
- `GET /orders` - List orders with user names
- `POST /orders/{id}/rate` - No longer returns 404 "Rating record not found"
- `POST /orders/{id}/cancel` - Properly creates rating record before negative rating

