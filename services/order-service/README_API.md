# Order Service API

Tài liệu này mô tả chi tiết các endpoint API của Order Service cho frontend sử dụng.

---

## 1. Tạo đơn hàng mới
- **POST** `/api/orders`
- **Input:**
```json
{
  "auction_id": 1,
  "winner_id": 6,
  "seller_id": 2,
  "final_price": 2500000
}
```
- **Output:**
```json
{
  "id": 1,
  "auction_id": 1,
  "winner_id": 6,
  "seller_id": 2,
  "final_price": 2500000,
  "status": "PENDING_PAYMENT",
  ...
}
```
- **Mục đích:** Tạo đơn hàng sau khi kết thúc đấu giá (gọi từ auction-service).

---

## 2. Lấy danh sách đơn hàng của user
- **GET** `/api/orders?role=buyer|seller&status=...`
- **Header:** `Authorization: Bearer <token>`
- **Output:**
```json
[
  {
    "id": 1,
    "auction_id": 1,
    "winner_id": 6,
    "seller_id": 2,
    "final_price": 2500000,
    "status": "PENDING_PAYMENT",
    ...
  },
  ...
]
```
- **Mục đích:** FE lấy danh sách đơn hàng của user (mua hoặc bán).

---

## 3. Lấy chi tiết đơn hàng
- **GET** `/api/orders/{id}`
- **Header:** `Authorization: Bearer <token>`
- **Output:**
```json
{
  "id": 1,
  "auction_id": 1,
  "winner_id": 6,
  "seller_id": 2,
  "final_price": 2500000,
  "status": "PENDING_PAYMENT",
  ...
}
```
- **Mục đích:** Xem chi tiết đơn hàng (chỉ buyer/seller liên quan mới xem được).

---

## 4. Thanh toán đơn hàng
- **POST** `/api/orders/{id}/pay`
- **Header:** `Authorization: Bearer <token>`
- **Input:**
```json
{
  "payment_method": "MOMO|ZALOPAY|VNPAY|STRIPE|PAYPAL",
  "payment_proof": "<url ảnh>"
}
```
- **Output:**
```json
{
  ...order object...
}
```
- **Mục đích:** Buyer thanh toán đơn hàng.

---

## 5. Cung cấp địa chỉ giao hàng
- **POST** `/api/orders/{id}/shipping-address`
- **Header:** `Authorization: Bearer <token>`
- **Input:**
```json
{
  "shipping_address": "123 Đường ABC, Quận 1, TP.HCM",
  "shipping_phone": "0901234567"
}
```
- **Output:**
```json
{
  ...order object...
}
```
- **Mục đích:** Buyer cung cấp địa chỉ nhận hàng.

---

## 6. Gửi hóa đơn vận chuyển (seller)
- **POST** `/api/orders/{id}/shipping-invoice`
- **Header:** `Authorization: Bearer <token>`
- **Input:**
```json
{
  "tracking_number": "VN123456789",
  "shipping_invoice": "<url file hóa đơn>"
}
```
- **Output:**
```json
{
  ...order object...
}
```
- **Mục đích:** Seller gửi mã vận đơn và hóa đơn vận chuyển.

---

## 7. Xác nhận đã nhận hàng (buyer)
- **POST** `/api/orders/{id}/confirm-delivery`
- **Header:** `Authorization: Bearer <token>`
- **Output:**
```json
{
  ...order object...
}
```
- **Mục đích:** Buyer xác nhận đã nhận hàng.

---

## 8. Hủy đơn hàng (seller)
- **POST** `/api/orders/{id}/cancel`
- **Header:** `Authorization: Bearer <token>`
- **Input:**
```json
{
  "cancel_reason": "Lý do hủy đơn"
}
```
- **Output:**
```json
{
  ...order object...
}
```
- **Mục đích:** Seller hủy đơn hàng trước khi hoàn thành.

---

## 9. Gửi tin nhắn chat trong đơn hàng
- **POST** `/api/orders/{id}/messages`
- **Header:** `Authorization: Bearer <token>`
- **Input:**
```json
{
  "message": "Nội dung tin nhắn"
}
```
- **Output:**
```json
{
  "id": 1,
  "order_id": 1,
  "sender_id": 6,
  "message": "Nội dung tin nhắn",
  "created_at": "2025-12-25T10:00:00Z"
}
```
- **Mục đích:** Buyer/seller chat với nhau trong đơn hàng.

---

## 10. Lấy danh sách tin nhắn chat
- **GET** `/api/orders/{id}/messages?limit=50&offset=0`
- **Header:** `Authorization: Bearer <token>`
- **Output:**
```json
[
  {
    "id": 1,
    "order_id": 1,
    "sender_id": 6,
    "message": "Nội dung tin nhắn",
    "created_at": "2025-12-25T10:00:00Z"
  },
  ...
]
```
- **Mục đích:** Lấy lịch sử chat của đơn hàng.

---

## 11. Đánh giá đơn hàng
- **POST** `/api/orders/{id}/rate`
- **Header:** `Authorization: Bearer <token>`
- **Input:**
```json
{
  "rating": 1, // 1: tốt, -1: xấu
  "comment": "Nhận xét về đối tác"
}
```
- **Output:**
```json
{
  "order_id": 1,
  "buyer_rating": 1,
  "buyer_comment": "Rất tốt!",
  "seller_rating": 1,
  "seller_comment": "Giao dịch tốt!",
  ...
}
```
- **Mục đích:** Buyer hoặc seller đánh giá đối tác sau giao dịch.

---

## 12. Lấy đánh giá đơn hàng
- **GET** `/api/orders/{id}/rating`
- **Output:**
```json
{
  "order_id": 1,
  "buyer_rating": 1,
  "buyer_comment": "Rất tốt!",
  "seller_rating": 1,
  "seller_comment": "Giao dịch tốt!",
  ...
}
```
- **Mục đích:** Xem đánh giá của đơn hàng.

---

## 13. Lấy thống kê rating của user
- **GET** `/api/users/{id}/rating`
- **Output:**
```json
{
  "user_id": 6,
  "total_number_reviews": 10,
  "total_number_good_reviews": 8
}
```
- **Mục đích:** Xem tổng số lượt đánh giá và lượt tốt của user.

---

## 14. Lấy tất cả đơn hàng (admin)
- **GET** `/api/admin/orders?status=...&limit=50&offset=0`
- **Header:** `Authorization: Bearer <token>`
- **Output:**
```json
[
  ...order object...
]
```
- **Mục đích:** Admin xem toàn bộ đơn hàng trong hệ thống.

---

## 15. Health check
- **GET** `/api/health`
- **Output:**
```json
{
  "status": "ok"
}
```
- **Mục đích:** Kiểm tra trạng thái service.

---

## 16. Swagger UI
- **GET** `/swagger/index.html`
- **Mục đích:** Xem tài liệu API trực quan.

---

## Lưu ý chung
- Các API cần xác thực đều yêu cầu header: `Authorization: Bearer <token>`
- Các trường output có thể có thêm các trường timestamp: `created_at`, `updated_at`, ...
- Các trường ...order object... là toàn bộ thông tin đơn hàng như mô tả ở trên.
