# User Service API Documentation

Base URL: `/api/users`

---

## 1. Lấy thông tin user đơn giản theo email
**GET** `/api/users/simple?email={email}`
- **Response:**
```json
{
  "success": true/false,
  "message": "string",
  "data": {
    "id": 1,
    "email": "string",
    "fullName": "string",
    "userRole": "BIDDER|SELLER|ADMIN"
  }
}
```

## 2. Lấy thông tin user đơn giản theo id
**GET** `/api/users/{id}/simple`
- **Response:** như trên

## 3. Đăng ký user
**POST** `/api/users`
- **Request Body:**
```json
{
  "fullName": "string",
  "email": "string",
  "password": "string",
  "birthDay": "yyyy-MM-dd"
}
```
- **Response:**
```json
{
  "success": true/false,
  "message": "string",
  "data": null
}
```

## 4. Xác thực email
**POST** `/api/users/verify-email`
- **Request Body:**
```json
{
  "email": "string"
}
```
- **Response:** như trên

## 5. Xoá user theo email
**DELETE** `/api/users`
- **Request Body:**
```json
{
  "email": "string"
}
```
- **Response:** như trên

## 6. Xác thực đăng nhập
**POST** `/api/users/authenticate`
- **Request Body:**
```json
{
  "email": "string",
  "password": "string"
}
```
- **Response:**
```json
{
  "success": true/false,
  "message": "string",
  "data": {
    "id": 1,
    "email": "string",
    "fullName": "string",
    "userRole": "BIDDER|SELLER|ADMIN"
  }
}
```

## 7. Lấy profile user hiện tại
**GET** `/api/users/profile/me`
- **Header:** `Authorization: Bearer <token>`
- **Response:**
```json
{
  "success": true/false,
  "message": "string",
  "data": {
    "id": 1,
    "email": "string",
    "fullName": "string",
    "birthDay": "yyyy-MM-dd",
    "userRole": "BIDDER|SELLER|ADMIN"
  }
}
```

## 8. Tìm kiếm user (admin, seller)
**GET** `/api/users/search?keyword=...&role=...&page=0&size=10`
- **Header:** `X-User-Token: <token>`
- **Authorization:** ROLE_ADMIN, ROLE_SELLER
- **Response:**
```json
{
  "content": [ ... ],
  "totalElements": 100,
  ...
}
```

---

## 9. Yêu cầu nâng cấp lên Seller
**POST** `/api/users/upgrade-to-seller?reason={reason}`
- **Header:** `X-User-Token: <token>`
- **Authorization:** ROLE_BIDDER
- **Query Parameters:**
  - `reason` (required): Lý do muốn nâng cấp
- **Response:**
```json
"Upgrade request submitted"
```

---

## 10. Lấy danh sách yêu cầu nâng cấp (Admin)
**GET** `/api/users/upgrade-requests?status=PENDING&page=0&size=10&sort=createdAt&direction=desc`
- **Header:** `X-User-Token: <token>`
- **Authorization:** ROLE_ADMIN
- **Query Parameters:**
  - `status` (optional): PENDING, APPROVED, REJECTED
  - `page` (default: 0): Số trang
  - `size` (default: 10): Số lượng kết quả
  - `sort` (default: createdAt): Trường sắp xếp
  - `direction` (default: desc): Hướng sắp xếp (asc/desc)
- **Response:**
```json
{
  "content": [
    {
      "id": 1,
      "user": {
        "id": 5,
        "email": "bidder@example.com",
        "fullName": "Nguyen Van A",
        "role": "ROLE_BIDDER"
      },
      "status": "PENDING",
      "reason": "I want to sell my handmade products",
      "rejectionReason": null,
      "createdAt": "2024-01-15T10:00:00Z",
      "reviewedAt": null,
      "reviewedByAdminId": null
    }
  ],
  "pageable": {
    "pageNumber": 0,
    "pageSize": 10
  },
  "totalElements": 1,
  "totalPages": 1
}
```

---

## 11. Duyệt yêu cầu nâng cấp (Admin)
**POST** `/api/users/{requestId}/approve`
- **Header:** `X-User-Token: <token>`
- **Authorization:** ROLE_ADMIN
- **Response:**
```json
"User upgraded to SELLER"
```

---

## 12. Từ chối yêu cầu nâng cấp (Admin)
**POST** `/api/users/{requestId}/reject?rejectReason={reason}`
- **Header:** `X-User-Token: <token>`
- **Authorization:** ROLE_ADMIN
- **Query Parameters:**
  - `rejectReason` (optional): Lý do từ chối
- **Response:**
```json
"Upgrade request rejected"
```

---

## 13. Lấy tất cả yêu cầu nâng cấp
**GET** `/api/users?page=0&size=10&sort=createdAt&direction=desc`
- **Header:** `X-User-Token: <token>`
- **Query Parameters:**
  - `page` (default: 0): Số trang
  - `size` (default: 10): Số lượng kết quả
  - `sort` (default: createdAt): Trường sắp xếp
  - `direction` (default: desc): Hướng sắp xếp
- **Response:** Similar to endpoint 10

---

## Database Schema - user_upgrade_requests

```sql
CREATE TABLE user_upgrade_requests (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    status VARCHAR(20) NOT NULL, -- PENDING, APPROVED, REJECTED
    reason TEXT,
    rejection_reason TEXT,
    created_at TIMESTAMP NOT NULL,
    reviewed_at TIMESTAMP,
    reviewed_by_admin_id BIGINT
);
```
