{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "fullName": "Nguyen Van A",
    "phoneNumber": "0123456789",
    "userRole": "ROLE_BIDDER"
  },
  "message": "User fetched successfully"
}
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "fullName": "Nguyen Van A",
    "phoneNumber": "0123456789",
    "userRole": "ROLE_BIDDER"
  },
  "message": "User fetched successfully"
}
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "fullName": "Nguyen Van A",
    "phoneNumber": "0123456789",
    "userRole": "ROLE_BIDDER"
  },
  "message": "User fetched successfully"
}
**Authorization:** ROLE_BIDDER, ROLE_SELLER, ROLE_ADMIN
**Authorization:** ROLE_ADMIN, ROLE_SELLER
**Authorization:** ROLE_BIDDER
**Authorization:** ROLE_ADMIN
- `role` (optional): ROLE_BIDDER, ROLE_SELLER, ROLE_ADMIN
# ONLINE AUCTION API DOCUMENTATION

**Base URL:** `http://localhost:8080`

**API Gateway Port:** 8080

---

## 📋 Mục lục

1. [Authentication](#1-authentication-service)
2. [User Service](#2-user-service)
3. [Category Service](#3-category-service)
4. [Product Service](#4-product-service)
5. [Bidding Service](#5-bidding-service)
6. [Order Service](#6-order-service)
7. [Media Service](#7-media-service)
8. [Comment Service](#8-comment-service)
9. [Notification Service](#9-notification-service)

---

## 🔐 Headers Chung

Hầu hết các endpoint (trừ Auth Service) yêu cầu: X-User-Token chứa JWT Access Token trong header (không có "Bearer" prefix) và Content-Type là application/json.

```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Lưu ý:** 
- Token được trả về sau khi login/register thành công
- Token không cần prefix "Bearer"
- Các endpoint thuộc Auth Service không yêu cầu `X-User-Token`

---

## 1. Authentication Service

**Routing:** `GET/POST/PUT/DELETE http://localhost:8080/api/auth/*` → `http://localhost:8081/auth/*`

### 1.1. Đăng ký tài khoản

**Endpoint:** `POST http://localhost:8080/api/auth/register`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "fullName": "Nguyen Van A",
  "phoneNumber": "0123456789"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "data": null,
  "message": "Successfully registered user"
}
```

**Response Error (400):**
```json
{
  "success": false,
  "data": null,
  "message": "Fail to register user, email is already registered"
}
```

---

### 1.2. Xác thực OTP

**Endpoint:** `POST http://localhost:8080/api/auth/verify-otp`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "otpCode": "123456"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "data": null,
  "message": "OTP verified successfully"
}
```

---

### 1.3. Đăng nhập

**Endpoint:** `POST http://localhost:8080/api/auth/sign-in`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "accessToken": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response Error (400):**
```json
{
  "success": false,
  "accessToken": "",
  "refreshToken": ""
}
```

---

### 1.4. Validate JWT Token

**Endpoint:** `POST http://localhost:8080/api/auth/validate-jwt`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response Success (200):**
```json
{
  "valid": true
}
```

**Response Error (401):**
```json
{
  "valid": false
}
```

---

### 1.5. Đăng nhập bằng Google

**Endpoint:** `POST http://localhost:8080/api/auth/sign-in/google`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "token": "google_id_token_here"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "accessToken": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## 2. User Service

**Routing:** `GET/POST/PUT/DELETE http://localhost:8080/api/users/*` → `http://localhost:8084/api/users/*`

**Required Header:** `X-User-Token: <JWT_ACCESS_TOKEN>`

### 2.1. Lấy thông tin user đơn giản theo email

**Endpoint:** `GET http://localhost:8080/api/users/simple?email=user@example.com`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "fullName": "Nguyen Van A",
    "phoneNumber": "0123456789",
    "userRole": "ROLE_BIDDER"
  },
  "message": "User fetched successfully"
}
```

---

### 2.2. Lấy thông tin user đơn giản theo ID

**Endpoint:** `GET http://localhost:8080/api/users/{id}/simple`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "fullName": "Nguyen Van A",
    "phoneNumber": "0123456789",
    "userRole": "ROLE_BIDDER"
  },
  "message": "User fetched successfully"
}
```

---

### 2.3. Lấy profile của user hiện tại

**Endpoint:** `GET http://localhost:8080/api/users/profile/me`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Authorization:** ROLE_BIDDER, ROLE_SELLER, ROLE_ADMIN

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "fullName": "Nguyen Van A",
    "phoneNumber": "0123456789",
    "userRole": "ROLE_BIDDER",
    "isEmailVerified": true,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "message": "Profile retrieved"
}
```

---

### 2.4. Tìm kiếm users

**Endpoint:** `GET http://localhost:8080/api/users/search`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Authorization:** ADMIN, SELLER

**Query Parameters:**
- `keyword` (optional): Từ khóa tìm kiếm
- `role` (optional): ROLE_BIDDER, ROLE_SELLER, ROLE_ADMIN
- `page` (default: 0): Số trang
- `size` (default: 10): Số lượng kết quả mỗi trang

**Example:** `GET http://localhost:8080/api/users/search?keyword=nguyen&role=ROLE_BIDDER&page=0&size=10`
**Response Success (200):**
```json
{
  "content": [
    {
      "id": 1,
      "email": "user@example.com",
      "fullName": "Nguyen Van A",
      "userRole": "ROLE_BIDDER"
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

### 2.5. Yêu cầu nâng cấp lên Seller

**Endpoint:** `POST http://localhost:8080/api/users/upgrade-to-seller?reason=I want to sell`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Authorization:** ROLE_BIDDER

**Query Parameters:**
- `reason` (required): Lý do muốn nâng cấp

**Response Success (200):**
```json
"Upgrade request submitted"
```

---

### 2.6. Duyệt yêu cầu nâng cấp (Admin only)

**Endpoint:** `POST http://localhost:8080/api/users/{requestId}/approve`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Authorization:** ADMIN

**Response Success (200):**
```json
"User upgraded to SELLER"
```

---

### 2.7. Lấy danh sách yêu cầu nâng cấp

**Endpoint:** `GET http://localhost:8080/api/users`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Query Parameters:**
- `page` (default: 0): Số trang
- `size` (default: 10): Số lượng kết quả
- `sort` (default: createdAt): Trường sắp xếp
- `direction` (default: desc): Hướng sắp xếp (asc/desc)

**Response Success (200):**
```json
{
  "content": [
    {
      "id": 1,
      "userId": 5,
      "reason": "I want to sell products",
      "status": "PENDING",
      "createdAt": "2024-01-01T00:00:00Z"
    }
  ],
  "totalElements": 1,
  "totalPages": 1
}
```

---

### 2.8. Xác thực email (Internal)

**Endpoint:** `POST http://localhost:8080/api/users/verify-email`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "data": null,
  "message": "Email verified successfully"
}
```

---

### 2.9. Xóa user theo email (Internal)

**Endpoint:** `DELETE http://localhost:8080/api/users`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "data": null,
  "message": "User deleted successfully"
}
```

---

## 3. Category Service

**Routing:** `GET/POST/PUT/DELETE http://localhost:8080/api/categories/*` → `http://localhost:8082/*`

**Required Header:** `X-User-Token: <JWT_ACCESS_TOKEN>`

### 3.1. Tạo danh mục mới

**Endpoint:** `POST http://localhost:8080/api/categories`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Electronics",
  "slug": "electronics",
  "description": "Electronic devices and accessories",
  "parent_id": null,
  "display_order": 1
}
```

**Response Success (201):**
```json
{
  "id": 1,
  "name": "Electronics",
  "slug": "electronics",
  "description": "Electronic devices and accessories",
  "parent_id": null,
  "level": 1,
  "is_active": true,
  "display_order": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Response Error (400):**
```json
{
  "error": "Maximum category depth is 2 levels"
}
```

---

### 3.2. Lấy danh sách categories

**Endpoint:** `GET http://localhost:8080/api/categories`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Query Parameters:**
- `parent_id` (optional): Lọc theo parent category
- `level` (optional): Lọc theo level (1 hoặc 2)

**Response Success (200):**
```json
{
  "categories": [
    {
      "id": 1,
      "name": "Electronics",
      "slug": "electronics",
      "description": "Electronic devices",
      "parent_id": null,
      "level": 1,
      "is_active": true,
      "display_order": 1,
      "children": [
        {
          "id": 2,
          "name": "Laptops",
          "slug": "laptops",
          "parent_id": 1,
          "level": 2,
          "is_active": true,
          "display_order": 1
        }
      ]
    }
  ]
}
```

---

### 3.3. Lấy category theo ID

**Endpoint:** `GET http://localhost:8080/api/categories/{id}`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
{
  "id": 1,
  "name": "Electronics",
  "slug": "electronics",
  "description": "Electronic devices",
  "parent_id": null,
  "level": 1,
  "is_active": true,
  "display_order": 1,
  "children": [
    {
      "id": 2,
      "name": "Laptops",
      "slug": "laptops",
      "parent_id": 1,
      "level": 2
    }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

### 3.4. Cập nhật category

**Endpoint:** `PUT http://localhost:8080/api/categories/{id}`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Electronics & Gadgets",
  "slug": "electronics-gadgets",
  "description": "All electronic devices",
  "is_active": true,
  "display_order": 1,
  "parent_id": null
}
```

**Response Success (200):**
```json
{
  "id": 1,
  "name": "Electronics & Gadgets",
  "slug": "electronics-gadgets",
  "description": "All electronic devices",
  "parent_id": null,
  "level": 1,
  "is_active": true,
  "display_order": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-02T00:00:00Z"
}
```

---

### 3.5. Xóa category (Soft delete)

**Endpoint:** `DELETE http://localhost:8080/api/categories/{id}`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
{
  "message": "Category deleted successfully"
}
```

**Response Error (400):**
```json
{
  "error": "Cannot delete category with active children"
}
```

---

### 3.6. Lấy categories theo parent ID

**Endpoint:** `GET http://localhost:8080/api/categories/parent/{parent_id}`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
[
  {
    "id": 2,
    "name": "Laptops",
    "slug": "laptops",
    "parent_id": 1,
    "level": 2,
    "is_active": true,
    "display_order": 1
  }
]
```

---

## 4. Product Service

**Routing:** `GET/POST/PUT/DELETE http://localhost:8080/api/products/*` → `http://localhost:8083/api/products/*`

**Required Header:** `X-User-Token: <JWT_ACCESS_TOKEN>`

### 4.1. Tạo sản phẩm mới (Seller only)

**Endpoint:** `POST http://localhost:8080/api/products`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Authorization:** SELLER

**Request Body:**
```json
{
  "name": "iPhone 15 Pro Max",
  "thumbnailUrl": "https://s3.amazonaws.com/bucket/thumb.jpg",
  "images": [
    "https://s3.amazonaws.com/bucket/image1.jpg",
    "https://s3.amazonaws.com/bucket/image2.jpg",
    "https://s3.amazonaws.com/bucket/image3.jpg"
  ],
  "description": "Latest iPhone model with A17 chip",
  "categoryId": 5,
  "categoryName": "Smartphones",
  "parentCategoryId": 1,
  "parentCategoryName": "Electronics",
  "startingPrice": 20000000.0,
  "buyNowPrice": 30000000.0,
  "stepPrice": 500000.0,
  "endAt": "2024-01-17T10:00:00",
  "autoExtend": true
}
```

**Response Success (200):**
```json
{
  "id": 1,
  "name": "iPhone 15 Pro Max",
  "thumbnailUrl": "https://s3.amazonaws.com/bucket/thumb.jpg",
  "images": [
    "https://s3.amazonaws.com/bucket/image1.jpg",
    "https://s3.amazonaws.com/bucket/image2.jpg",
    "https://s3.amazonaws.com/bucket/image3.jpg"
  ],
  "description": "Latest iPhone model with A17 chip",
  "parentCategoryId": 1,
  "parentCategoryName": "Electronics",
  "categoryId": 5,
  "categoryName": "Smartphones",
  "startingPrice": 20000000.0,
  "currentPrice": 20000000.0,
  "buyNowPrice": 30000000.0,
  "stepPrice": 500000.0,
  "createdAt": "2024-01-10T09:00:00",
  "endAt": "2024-01-17T10:00:00",
  "autoExtend": true,
  "extendThresholdMinutes": 5,
  "extendDurationMinutes": 10,
  "sellerId": 10,
  "sellerInfo": {
    "userId": 10,
    "username": "seller_user",
    "avatarUrl": "https://avatar.com/johndoe.png"
  },
  "highestBidder": null
}
```

---

### 4.2. Lấy chi tiết sản phẩm

**Endpoint:** `GET http://localhost:8080/api/products/{productId}`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
{
  "id": 1,
  "name": "iPhone 15 Pro Max",
  "thumbnailUrl": "https://s3.amazonaws.com/bucket/thumb.jpg",
  "images": [
    "https://s3.amazonaws.com/bucket/image1.jpg",
    "https://s3.amazonaws.com/bucket/image2.jpg",
    "https://s3.amazonaws.com/bucket/image3.jpg"
  ],
  "description": "Latest iPhone model with A17 chip",
  "parentCategoryId": 1,
  "parentCategoryName": "Electronics",
  "categoryId": 5,
  "categoryName": "Smartphones",
  "startingPrice": 20000000.0,
  "currentPrice": 20000000.0,
  "buyNowPrice": 30000000.0,
  "stepPrice": 500000.0,
  "createdAt": "2024-01-10T09:00:00",
  "endAt": "2024-01-17T10:00:00",
  "autoExtend": true,
  "extendThresholdMinutes": 5,
  "extendDurationMinutes": 10,
  "sellerId": 10,
  "sellerInfo": {
    "userId": 10,
    "username": "seller_user",
    "avatarUrl": "https://avatar.com/johndoe.png"
  },
  "highestBidder": null
}
```

---

### 4.3. Cập nhật mô tả sản phẩm (Seller only)

**Endpoint:** `PATCH http://localhost:8080/api/products/{productId}/description`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Authorization:** SELLER (chỉ seller sở hữu sản phẩm)

**Request Body:**
```json
{
  "additionalDescription": "Additional info: Brand new, sealed box"
}
```

**Response Success (200):**
```json
{
  "id": 1,
  "name": "iPhone 15 Pro Max",
  "thumbnailUrl": "https://s3.amazonaws.com/bucket/thumb.jpg",
  "images": [
    "https://s3.amazonaws.com/bucket/image1.jpg",
    "https://s3.amazonaws.com/bucket/image2.jpg",
    "https://s3.amazonaws.com/bucket/image3.jpg"
  ],
  "description": "Latest iPhone model with A17 chip",
  "parentCategoryId": 1,
  "parentCategoryName": "Electronics",
  "categoryId": 5,
  "categoryName": "Smartphones",
  "startingPrice": 20000000.0,
  "currentPrice": 20000000.0,
  "buyNowPrice": 30000000.0,
  "stepPrice": 500000.0,
  "createdAt": "2024-01-10T09:00:00",
  "endAt": "2024-01-17T10:00:00",
  "autoExtend": true,
  "extendThresholdMinutes": 5,
  "extendDurationMinutes": 10,
  "sellerId": 10,
  "sellerInfo": {
    "userId": 10,
    "username": "seller_user",
    "avatarUrl": "https://avatar.com/johndoe.png"
  },
  "highestBidder": null
}
```

---

### 4.4. Lấy danh sách sản phẩm của seller

**Endpoint:** `GET http://localhost:8080/api/products/seller/{sellerId}`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
[
  {
    "id": 1,
    "name": "iPhone 15 Pro Max",
    "thumbnailUrl": "https://s3.amazonaws.com/bucket/thumb.jpg",
    "images": [
      "https://s3.amazonaws.com/bucket/image1.jpg",
      "https://s3.amazonaws.com/bucket/image2.jpg",
      "https://s3.amazonaws.com/bucket/image3.jpg"
    ],
    "description": "Latest iPhone model with A17 chip",
    "parentCategoryId": 1,
    "parentCategoryName": "Electronics",
    "categoryId": 5,
    "categoryName": "Smartphones",
    "startingPrice": 20000000.0,
    "currentPrice": 20000000.0,
    "buyNowPrice": 30000000.0,
    "stepPrice": 500000.0,
    "createdAt": "2024-01-10T09:00:00",
    "endAt": "2024-01-17T10:00:00",
    "autoExtend": true,
    "extendThresholdMinutes": 5,
    "extendDurationMinutes": 10,
    "sellerId": 10,
    "sellerInfo": {
      "userId": 10,
      "username": "seller_user",
      "avatarUrl": "https://avatar.com/johndoe.png"
    },
    "highestBidder": null
  }
]
```

---

### 4.5. Top 5 sản phẩm sắp kết thúc

**Endpoint:** `GET http://localhost:8080/api/products/top-ending`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "thumbnailUrl": "https://s3.amazonaws.com/bucket/image1.jpg",
      "name": "iPhone 15 Pro Max",
      "currentPrice": 25000000.0,
      "buyNowPrice": 30000000.0,
      "createdAt": "2024-01-10T08:00:00",
      "endAt": "2024-01-17T10:00:00",
      "bidCount": 15,
      "categoryParentId": 1,
      "categoryParentName": "Electronics",
      "categoryId": 5,
      "categoryName": "Smartphones"
    }
  ],
  "message": "Successfully fetching top5 ending-soon products"
}
```

---

### 4.6. Top 5 sản phẩm có nhiều lượt đấu giá nhất

**Endpoint:** `GET http://localhost:8080/api/products/top-most-bids`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "thumbnailUrl": "https://s3.amazonaws.com/bucket/image1.jpg",
      "name": "iPhone 15 Pro Max",
      "currentPrice": 25000000.0,
      "buyNowPrice": 30000000.0,
      "createdAt": "2024-01-10T08:00:00",
      "endAt": "2024-01-17T10:00:00",
      "bidCount": 15,
      "categoryParentId": 1,
      "categoryParentName": "Electronics",
      "categoryId": 5,
      "categoryName": "Smartphones"
    }
  ],
  "message": "Successfully fetching top5 most-bids products"
}
```

---

### 4.7. Top 5 sản phẩm giá cao nhất

**Endpoint:** `GET http://localhost:8080/api/products/top-highest-price`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "thumbnailUrl": "https://s3.amazonaws.com/bucket/image1.jpg",
      "name": "iPhone 15 Pro Max",
      "currentPrice": 25000000.0,
      "buyNowPrice": 30000000.0,
      "createdAt": "2024-01-10T08:00:00",
      "endAt": "2024-01-17T10:00:00",
      "bidCount": 15,
      "categoryParentId": 1,
      "categoryParentName": "Electronics",
      "categoryId": 5,
      "categoryName": "Smartphones"
    }
  ],
  "message": "Successfully fetching top5 highest-price products"
}
```

---

### 4.8. Tìm kiếm và lọc sản phẩm

**Endpoint:** `GET http://localhost:8080/api/products/search`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Query Parameters:**
- `query` (optional): Từ khóa tìm kiếm
- `parentCategoryId` (optional): Lọc theo danh mục cha
- `categoryId` (optional): Lọc theo danh mục con
- `page` (default: 0): Số trang
- `pageSize` (default: 10): Số lượng kết quả

**Example:** `GET http://localhost:8080/api/products/search?query=iphone&categoryId=5&page=0&pageSize=10`

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "content": [
      {
        "id": 1,
        "thumbnailUrl": "https://s3.amazonaws.com/bucket/image1.jpg",
        "name": "iPhone 15 Pro Max",
        "currentPrice": 25000000.0,
        "buyNowPrice": 30000000.0,
        "createdAt": "2024-01-10T08:00:00",
        "endAt": "2024-01-17T10:00:00",
        "bidCount": 15,
        "categoryParentId": 1,
        "categoryParentName": "Electronics",
        "categoryId": 5,
        "categoryName": "Smartphones"
      }
    ],
    "totalElements": 1,
    "totalPages": 1,
    "size": 10,
    "number": 0,
    "numberOfElements": 1,
    "first": true,
    "last": true,
    "empty": false
  },
  "message": "Query success"
}
```

---

### 4.9. Cập nhật category (Internal - Category Service)

**Endpoint:** `PUT http://localhost:8080/api/products/categories/{categoryId}`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Request Body:**
```json
{
  "categoryName": "Gaming Laptop",
  "parentCategoryId": 2,
  "parentCategoryName": "Tech Devices"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "updatedCount": 15
  },
  "message": "Category updated successfully"
}
```

---

### 4.10. Đổi tên parent category (Internal - Category Service)

**Endpoint:** `PUT http://localhost:8080/api/products/parent-categories/{parentCategoryId}/rename`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Request Body:**
```json
{
  "parentCategoryName": "Tech Devices"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "data": {
    "updatedCount": 25
  },
  "message": "Parent category renamed successfully"
}
```

---

## 5. Bidding Service

**Routing:** `GET/POST http://localhost:8080/api/bids/*` → `http://localhost:8085/*`

**Required Header:** `X-User-Token: <JWT_ACCESS_TOKEN>`

### 5.1. Đặt giá thầu

**Endpoint:** `POST http://localhost:8080/api/bids`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Authorization:** ROLE_BIDDER, ROLE_SELLER

**Request Body:**
```json
{
  "productId": 1,
  "amount": 25500000,
  "requestId": "unique-request-id-123"
}
```

**Response Success (200):**
```json
{
  "success": true,
  "message": "Bid placed successfully",
  "data": {
    "newHighest": 25500000.0,
    "previousHighestBidder": 10
  }
}
```

**Response Error (400):**
```json
{
  {
    "success": false,
    "message": "Bid amount too low",
    "data": null
  }
}
```

---

### 5.2. Tìm kiếm lịch sử đấu giá

**Endpoint:** `GET http://localhost:8080/api/bids/search`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Authorization:** ROLE_ADMIN, ROLE_SELLER, ROLE_BIDDER

**Query Parameters:**
- `productId` (optional): Lọc theo sản phẩm
- `bidderId` (optional): Lọc theo người đấu giá
- `status` (optional): SUCCESS, FAILED
- `requestId` (optional): Lọc theo request ID
- `from` (optional): Thời gian bắt đầu (ISO 8601)
- `to` (optional): Thời gian kết thúc (ISO 8601)
- `page` (default: 0): Số trang
- `size` (default: 10): Số lượng kết quả

**Example:** `GET http://localhost:8080/api/bids/search?productId=1&status=SUCCESS&page=0&size=10`

**Response Success (200):**
```json
{
  "content": [
    {
      "id": 1,
      "productId": 1,
      "bidderId": 5,
      "amount": 25500000,
      "status": "SUCCESS",
      "requestId": "unique-request-id-123",
      "createdAt": "2024-01-15T10:30:00Z"
    }
  ],
  "totalElements": 1,
  "totalPages": 1,
  "size": 10,
  "number": 0
}
```

---

## 6. Order Service

**Routing:** `GET/POST/PUT http://localhost:8080/api/orders/data/*` → `http://localhost:8086/product/*`

**Required Header:** `X-User-Token: <JWT_ACCESS_TOKEN>`

**Order Status Flow:**
```
PENDING_PAYMENT → PAID → ADDRESS_PROVIDED → SHIPPING → DELIVERED → COMPLETED
                                                             ↓
                                                        CANCELLED (có thể cancel bất kỳ lúc nào trước COMPLETED)
```

### 6.1. Tạo đơn hàng (Internal - sau khi auction kết thúc)

**Endpoint:** `POST http://localhost:8080/api/orders/data/product/`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Request Body:**
```json
{
  "auction_id": 1,
  "winner_id": 5,
  "seller_id": 10,
  "final_price": 26000000
}
```

**Response Success (201):**
```json
{
  "id": 1,
  "auction_id": 1,
  "winner_id": 5,
  "seller_id": 10,
  "final_price": 26000000,
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
  "created_at": "2024-01-17T10:00:00Z",
  "updated_at": "2024-01-17T10:00:00Z"
}
```

---

### 6.2. Lấy chi tiết đơn hàng

**Endpoint:** `GET http://localhost:8080/api/orders/data/product/{id}`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Authorization:** Chỉ buyer hoặc seller của đơn hàng

**Response Success (200):**
```json
{
  "id": 1,
  "auction_id": 1,
  "winner_id": 5,
  "seller_id": 10,
  "final_price": 26000000,
  "status": "PAID",
  "payment_method": "MOMO",
  "payment_proof": "https://s3.amazonaws.com/proof.jpg",
  "shipping_address": "",
  "shipping_phone": "",
  "tracking_number": "",
  "shipping_invoice": "",
  "paid_at": "2024-01-17T11:00:00Z",
  "delivered_at": null,
  "completed_at": null,
  "cancelled_at": null,
  "cancel_reason": "",
  "rating": {
    "id": 1,
    "order_id": 1,
    "buyer_rating": null,
    "buyer_comment": "",
    "seller_rating": null,
    "seller_comment": "",
    "buyer_rated_at": null,
    "seller_rated_at": null,
    "created_at": "2024-01-17T10:00:00Z",
    "updated_at": "2024-01-17T10:00:00Z"
  },
  "created_at": "2024-01-17T10:00:00Z",
  "updated_at": "2024-01-17T11:00:00Z"
}
```

**Response Error (403):**
```json
{
  "error": "Access denied"
}
```

**Response Error (404):**
```json
{
  "error": "Order not found"
}
```

---

### 6.3. Lấy danh sách đơn hàng của user

**Endpoint:** `GET http://localhost:8080/api/orders/data/product/`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Query Parameters:**
- `role` (optional): `buyer`, `seller`, `ROLE_BIDDER`, `ROLE_SELLER` - Lọc theo vai trò
- `status` (optional): `PENDING_PAYMENT`, `PAID`, `ADDRESS_PROVIDED`, `SHIPPING`, `DELIVERED`, `COMPLETED`, `CANCELLED`

**Example:** `GET http://localhost:8080/api/orders/data/product/?role=buyer&status=COMPLETED`

**Response Success (200):**
```json
[
  {
    "id": 1,
    "auction_id": 1,
    "winner_id": 5,
    "seller_id": 10,
    "final_price": 26000000,
    "status": "COMPLETED",
    "payment_method": "MOMO",
    "payment_proof": "https://s3.amazonaws.com/proof.jpg",
    "shipping_address": "123 Nguyen Hue, District 1, HCM",
    "shipping_phone": "0901234567",
    "tracking_number": "VN123456789",
    "shipping_invoice": "https://s3.amazonaws.com/invoice.jpg",
    "paid_at": "2024-01-17T11:00:00Z",
    "delivered_at": "2024-01-20T15:00:00Z",
    "completed_at": "2024-01-20T16:00:00Z",
    "cancelled_at": null,
    "cancel_reason": "",
    "rating": {
      "id": 1,
      "order_id": 1,
      "buyer_rating": 1,
      "buyer_comment": "Great seller!",
      "seller_rating": 1,
      "seller_comment": "Good buyer!",
      "buyer_rated_at": "2024-01-20T16:00:00Z",
      "seller_rated_at": "2024-01-20T16:05:00Z",
      "created_at": "2024-01-17T10:00:00Z",
      "updated_at": "2024-01-20T16:05:00Z"
    },
    "created_at": "2024-01-17T10:00:00Z",
    "updated_at": "2024-01-20T16:05:00Z"
  }
]
```

---

### 6.4. Thanh toán đơn hàng (Buyer)

**Endpoint:** `POST http://localhost:8080/api/orders/data/product/{id}/pay`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Authorization:** Chỉ buyer của đơn hàng, order phải ở trạng thái `PENDING_PAYMENT`

**Request Body:**
```json
{
  "payment_method": "MOMO",
  "payment_proof": "https://s3.amazonaws.com/proof.jpg"
}
```

**Payment Methods:** `MOMO`, `ZALOPAY`, `VNPAY`, `STRIPE`, `PAYPAL`

**Response Success (200):**
```json
{
  "id": 1,
  "auction_id": 1,
  "winner_id": 5,
  "seller_id": 10,
  "final_price": 26000000,
  "status": "PAID",
  "payment_method": "MOMO",
  "payment_proof": "https://s3.amazonaws.com/proof.jpg",
  "paid_at": "2024-01-17T11:00:00Z",
  "created_at": "2024-01-17T10:00:00Z",
  "updated_at": "2024-01-17T11:00:00Z"
}
```

**Response Error (400):**
```json
{
  "error": "Cannot pay order with status: PAID"
}
```

**Response Error (403):**
```json
{
  "error": "Only buyer can pay for order"
}
```

---

### 6.5. Cung cấp địa chỉ giao hàng (Buyer)

**Endpoint:** `POST http://localhost:8080/api/orders/data/product/{id}/shipping-address`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Authorization:** Chỉ buyer, order phải ở trạng thái `PAID`

**Request Body:**
```json
{
  "shipping_address": "123 Nguyen Hue, District 1, Ho Chi Minh City, Vietnam",
  "shipping_phone": "0901234567"
}
```

**Response Success (200):**
```json
{
  "id": 1,
  "auction_id": 1,
  "winner_id": 5,
  "seller_id": 10,
  "final_price": 26000000,
  "status": "ADDRESS_PROVIDED",
  "payment_method": "MOMO",
  "shipping_address": "123 Nguyen Hue, District 1, Ho Chi Minh City, Vietnam",
  "shipping_phone": "0901234567",
  "paid_at": "2024-01-17T11:00:00Z",
  "created_at": "2024-01-17T10:00:00Z",
  "updated_at": "2024-01-17T12:00:00Z"
}
```

**Response Error (400):**
```json
{
  "error": "Cannot provide address for order with status: PENDING_PAYMENT"
}
```

---

### 6.6. Gửi thông tin vận chuyển (Seller)

**Endpoint:** `POST http://localhost:8080/api/orders/data/product/{id}/shipping-invoice`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Authorization:** Chỉ seller, order phải ở trạng thái `ADDRESS_PROVIDED`

**Request Body:**
```json
{
  "tracking_number": "VN123456789",
  "shipping_invoice": "https://s3.amazonaws.com/invoice.jpg"
}
```

**Response Success (200):**
```json
{
  "id": 1,
  "auction_id": 1,
  "winner_id": 5,
  "seller_id": 10,
  "final_price": 26000000,
  "status": "SHIPPING",
  "tracking_number": "VN123456789",
  "shipping_invoice": "https://s3.amazonaws.com/invoice.jpg",
  "created_at": "2024-01-17T10:00:00Z",
  "updated_at": "2024-01-18T09:00:00Z"
}
```

**Response Error (403):**
```json
{
  "error": "Only seller can send shipping invoice"
}
```

---

### 6.7. Xác nhận đã nhận hàng (Buyer)

**Endpoint:** `POST http://localhost:8080/api/orders/data/product/{id}/confirm-delivery`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Authorization:** Chỉ buyer, order phải ở trạng thái `SHIPPING`

**Response Success (200):**
```json
{
  "id": 1,
  "auction_id": 1,
  "winner_id": 5,
  "seller_id": 10,
  "final_price": 26000000,
  "status": "DELIVERED",
  "delivered_at": "2024-01-20T15:00:00Z",
  "created_at": "2024-01-17T10:00:00Z",
  "updated_at": "2024-01-20T15:00:00Z"
}
```

**Response Error (400):**
```json
{
  "error": "Cannot confirm delivery for order with status: PAID"
}
```

---

### 6.8. Hủy đơn hàng (Seller)

**Endpoint:** `POST http://localhost:8080/api/orders/data/product/{id}/cancel`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Authorization:** Chỉ seller, order không được ở trạng thái `COMPLETED` hoặc `CANCELLED`

**Request Body:**
```json
{
  "cancel_reason": "Product is no longer available"
}
```

**Response Success (200):**
```json
{
  "id": 1,
  "auction_id": 1,
  "winner_id": 5,
  "seller_id": 10,
  "final_price": 26000000,
  "status": "CANCELLED",
  "cancel_reason": "Product is no longer available",
  "cancelled_at": "2024-01-18T10:00:00Z",
  "created_at": "2024-01-17T10:00:00Z",
  "updated_at": "2024-01-18T10:00:00Z"
}
```

**Lưu ý:** Khi seller hủy đơn, buyer sẽ tự động nhận rating -1 từ hệ thống

---

### 6.9. Gửi tin nhắn trong order

**Endpoint:** `POST http://localhost:8080/api/orders/data/product/{id}/messages`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Authorization:** Buyer hoặc seller của order

**Request Body:**
```json
{
  "message": "When will you ship the product?"
}
```

**Response Success (201):**
```json
{
  "id": 1,
  "order_id": 1,
  "sender_id": 5,
  "message": "When will you ship the product?",
  "created_at": "2024-01-17T12:00:00Z"
}
```

---

### 6.10. Lấy danh sách tin nhắn trong order

**Endpoint:** `GET http://localhost:8080/api/orders/data/product/{id}/messages`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Authorization:** Buyer hoặc seller của order

**Query Parameters:**
- `limit` (default: 50): Số lượng tin nhắn
- `offset` (default: 0): Vị trí bắt đầu

**Example:** `GET http://localhost:8080/api/orders/data/product/1/messages?limit=20&offset=0`

**Response Success (200):**
```json
[
  {
    "id": 1,
    "order_id": 1,
    "sender_id": 5,
    "message": "When will you ship the product?",
    "created_at": "2024-01-17T12:00:00Z"
  },
  {
    "id": 2,
    "order_id": 1,
    "sender_id": 10,
    "message": "I will ship it tomorrow",
    "created_at": "2024-01-17T13:00:00Z"
  }
]
```

---

### 6.11. Đánh giá đơn hàng

**Endpoint:** `POST http://localhost:8080/api/orders/data/product/{id}/rate`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Authorization:** Buyer hoặc seller của order

**Request Body:**
```json
{
  "rating": 1,
  "comment": "Great transaction!"
}
```

**Rating Values:**
- `1`: Đánh giá tích cực (good review)
- `-1`: Đánh giá tiêu cực (bad review)

**Response Success (200):**
```json
{
  "id": 1,
  "order_id": 1,
  "buyer_rating": 1,
  "buyer_comment": "Great transaction!",
  "seller_rating": null,
  "seller_comment": "",
  "buyer_rated_at": "2024-01-20T16:00:00Z",
  "seller_rated_at": null,
  "created_at": "2024-01-17T10:00:00Z",
  "updated_at": "2024-01-20T16:00:00Z"
}
```

**Lưu ý:** 
- Nếu cả buyer và seller đều đã đánh giá và order ở trạng thái `DELIVERED`, order sẽ tự động chuyển sang `COMPLETED`
- Rating sẽ cập nhật vào thống kê rating của user (total_number_reviews và total_number_good_reviews)

---

### 6.12. Lấy thông tin rating của order

**Endpoint:** `GET http://localhost:8080/api/orders/data/product/{id}/rating`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
{
  "id": 1,
  "order_id": 1,
  "buyer_rating": 1,
  "buyer_comment": "Great seller!",
  "seller_rating": 1,
  "seller_comment": "Good buyer!",
  "buyer_rated_at": "2024-01-20T16:00:00Z",
  "seller_rated_at": "2024-01-20T16:05:00Z",
  "created_at": "2024-01-17T10:00:00Z",
  "updated_at": "2024-01-20T16:05:00Z"
}
```

---

### 6.13. Lấy thống kê rating của user

**Endpoint:** `GET http://localhost:8080/api/orders/data/users/{id}/rating`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
{
  "user_id": 5,
  "total_number_reviews": 25,
  "total_number_good_reviews": 22,
  "rating_percentage": 88.0
}
```

**Response Error (404):**
```json
{
  "error": "User not found"
}
```

---

### 6.14. Lấy tất cả orders (Admin only)

**Endpoint:** `GET http://localhost:8080/api/orders/data/admin/orders`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Authorization:** ROLE_ADMIN only

**Query Parameters:**
- `status` (optional): Filter by status
- `limit` (default: 50): Số lượng orders
- `offset` (default: 0): Vị trí bắt đầu

**Example:** `GET http://localhost:8080/api/orders/data/admin/orders?status=COMPLETED&limit=20`

**Response Success (200):**
```json
[
  {
    "id": 1,
    "auction_id": 1,
    "winner_id": 5,
    "seller_id": 10,
    "final_price": 26000000,
    "status": "COMPLETED",
    "rating": {
      "id": 1,
      "order_id": 1,
      "buyer_rating": 1,
      "seller_rating": 1
    },
    "created_at": "2024-01-17T10:00:00Z",
    "updated_at": "2024-01-20T16:05:00Z"
  }
]
```

**Response Error (403):**
```json
{
  "error": "Admin access required"
}
```

---

### 6.15. WebSocket kết nối chat trong order

**Endpoint:** `GET http://localhost:8080/api/order-websocket/`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
{
  "order_service_websocket_url": "ws://localhost:8086/ws",
  "internal_jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Cách sử dụng WebSocket:**

1. Gọi endpoint trên để lấy `order_service_websocket_url` và `internal_jwt`
2. Kết nối WebSocket với URL:
   ```
   ws://localhost:8086/ws?orderId=1&X-User-Token=<JWT_ACCESS_TOKEN>&X-Internal-JWT=<internal_jwt>
   ```

3. Gửi tin nhắn:
   ```json
   {
     "type": "message",
     "content": "Hello, seller!"
   }
   ```

4. Gửi typing indicator:
   ```json
   {
     "type": "typing"
   }
   ```

5. Nhận tin nhắn mới:
   ```json
   {
     "type": "message",
     "orderId": 1,
     "data": {
       "id": 5,
       "order_id": 1,
       "sender_id": 10,
       "message": "Hi, buyer!",
       "created_at": "2024-01-17T14:00:00Z"
     }
   }
   ```

6. Nhận typing indicator:
   ```json
   {
     "type": "typing",
     "orderId": 1,
     "data": {
       "userId": 10
     }
   }
   ```

---

## 7. Media Service

**Routing:** `GET/POST http://localhost:8080/api/media/*` → `http://localhost:8089/*`

**Required Header:** `X-User-Token: <JWT_ACCESS_TOKEN>`

### 7.1. Upload file đơn

**Endpoint:** `POST http://localhost:8080/api/media/upload`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: multipart/form-data
```

**Request Body (Form Data):**
- `file`: File cần upload
- `folder` (query param, optional): Thư mục đích (default: "uploads/")

**Example:** `POST http://localhost:8080/api/media/upload?folder=products/`

**Response Success (200):**
```json
{
  "message": "Upload thành công",
  "url": "https://wnc-s3.s3.ap-southeast-1.amazonaws.com/products/20240101-uuid-image.jpg",
  "key": "products/20240101-uuid-image.jpg",
  "filename": "image.jpg",
  "size": 1048576,
  "uploaded_at": "2024-01-01T00:00:00Z"
}
```

**Response Error (400):**
```json
{
  "error": "File quá lớn, tối đa 10MB"
}
```

---

### 7.2. Upload nhiều file

**Endpoint:** `POST http://localhost:8080/api/media/upload/multiple`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: multipart/form-data
```

**Request Body (Form Data):**
- `files`: Danh sách file cần upload (multiple)
- `folder` (query param, optional): Thư mục đích

**Example:** `POST http://localhost:8080/api/media/upload/multiple?folder=products/`

**Response Success (200):**
```json
{
  "message": "Uploaded 3/3 files successfully",
  "uploaded": [
    {
      "message": "Upload thành công",
      "url": "https://wnc-s3.s3.ap-southeast-1.amazonaws.com/products/image1.jpg",
      "key": "products/image1.jpg",
      "filename": "image1.jpg",
      "size": 1048576,
      "uploaded_at": "2024-01-01T00:00:00Z"
    }
  ],
  "failed": [],
  "total": 3,
  "success_count": 3,
  "failed_count": 0
}
```

---

### 7.3. Lấy Presigned URL cho upload trực tiếp

**Endpoint:** `GET http://localhost:8080/api/media/presign`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Query Parameters:**
- `filename` (required): Tên file muốn upload
- `folder` (optional): Thư mục đích

**Example:** `GET http://localhost:8080/api/media/presign?filename=product.jpg&folder=products/`

**Response Success (200):**
```json
{
  "presigned_url": "https://wnc-s3.s3.ap-southeast-1.amazonaws.com/products/20240101-uuid-product.jpg?X-Amz-Algorithm=...",
  "image_url": "https://wnc-s3.s3.ap-southeast-1.amazonaws.com/products/20240101-uuid-product.jpg",
  "key": "products/20240101-uuid-product.jpg",
  "expires_in": 900
}
```

**Cách sử dụng:**
1. Client gọi endpoint này để lấy `presigned_url`
2. Client upload file trực tiếp đến `presigned_url` bằng PUT request
3. Sau khi upload thành công, sử dụng `image_url` để lưu vào database

---

## 8. Comment Service

**Routing:** `GET/POST http://localhost:8080/api/comments/*` → `http://localhost:8091/*`

**Required Header:** `X-User-Token: <JWT_ACCESS_TOKEN>`

### 8.1. Lấy lịch sử bình luận của sản phẩm

**Endpoint:** `GET http://localhost:8080/api/comments/history/products/{productId}`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Query Parameters:**
- `limit` (default: 50): Số lượng bình luận
- `offset` (default: 0): Vị trí bắt đầu

**Example:** `GET http://localhost:8080/api/comments/history/products/1?limit=50&offset=0`

**Response Success (200):**
```json
[
  {
    "id": 1,
    "product_id": 1,
    "sender_id": 5,
    "content": "Sản phẩm còn bảo hành không?",
    "created_at": "2024-01-15T10:30:00Z"
  },
  {
    "id": 2,
    "product_id": 1,
    "sender_id": 10,
    "content": "Còn bảo hành 12 tháng bạn nhé",
    "created_at": "2024-01-15T10:35:00Z"
  }
]
```

---

### 8.2. WebSocket kết nối chat real-time

**Endpoint:** `GET http://localhost:8080/api/comments/websocket/*`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
```

**Response Success (200):**
```json
{
  "comment_service_websocket_url": "ws://localhost:8091/ws",
  "internal_jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Cách sử dụng WebSocket:**

1. **Lấy thông tin kết nối:**
   ```
   GET http://localhost:8080/api/comments/websocket/*
   ```

2. **Kết nối WebSocket:**
   ```
   ws://localhost:8091/ws?productId=1&X-User-Token=<JWT>&X-Internal-JWT=<internal_jwt>
   ```

3. **Gửi tin nhắn:**
   ```json
   {
     "type": "comment",
     "content": "Hello, is this product still available?"
   }
   ```

4. **Nhận tin nhắn:**
   ```json
   {
     "type": "new_comment",
     "data": {
       "id": 1,
       "product_id": 1,
       "sender_id": 5,
       "content": "Hello, is this product still available?",
       "created_at": "2024-01-15T10:30:00Z"
     }
   }
   ```

---

## 9. Notification Service

**Routing:** `POST http://localhost:8080/api/notifications/*` → `http://localhost:8088/api/notify/*`

**Required Header:** `X-User-Token: <JWT_ACCESS_TOKEN>`

### 9.1. Gửi email thông báo

**Endpoint:** `POST http://localhost:8080/api/notifications/email`

**Headers:**
```
X-User-Token: <JWT_ACCESS_TOKEN>
Content-Type: application/json
```

**Request Body:**
```json
{
  "to": "user@example.com",
  "subject": "Bid Notification",
  "body": "Your bid has been placed successfully"
}
```

**Response Success (200):**
```json
{
  "message": "Email sent successfully"
}
```

---

## 📝 Lưu ý quan trọng

### Authentication Flow

1. **Đăng ký:** `POST /api/auth/register` → Nhận OTP qua email
2. **Xác thực OTP:** `POST /api/auth/verify-otp` → Kích hoạt tài khoản
3. **Đăng nhập:** `POST /api/auth/sign-in` → Nhận `accessToken` và `refreshToken`
4. **Sử dụng API:** Gửi `accessToken` qua header `X-User-Token` cho các request tiếp theo

### Status Code Summary

- `200 OK`: Request thành công
- `201 Created`: Tạo resource thành công
- `400 Bad Request`: Dữ liệu request không hợp lệ
- `401 Unauthorized`: Chưa đăng nhập hoặc token không hợp lệ
- `403 Forbidden`: Không có quyền truy cập
- `404 Not Found`: Resource không tồn tại
- `500 Internal Server Error`: Lỗi server

### Role-Based Access Control

- **ROLE_BIDDER**: Người đấu giá (có thể đặt giá, xem sản phẩm)
- **ROLE_SELLER**: Người bán (có thể tạo sản phẩm, quản lý sản phẩm của mình)
- **ROLE_ADMIN**: Quản trị viên (toàn quyền truy cập)

### WebSocket Connections

Để kết nối WebSocket:
1. Lấy thông tin kết nối từ API Gateway
2. Sử dụng `internal_jwt` và `X-User-Token` khi kết nối
3. Gửi/nhận message theo format JSON

### API Gateway Routing Rules

- **Auth Service**: `/api/auth/*` → Không cần token
- **Protected Services**: `/api/{service}/*` → Yêu cầu `X-User-Token`
- **WebSocket**: API Gateway trả về URL và JWT để kết nối trực tiếp

---

## 🔗 Service URLs (Internal)

Các URL này chỉ dùng trong môi trường development và không được expose ra ngoài:

- API Gateway: `http://localhost:8080`
- Auth Service: `http://localhost:8081`
- Category Service: `http://localhost:8082`
- Product Service: `http://localhost:8083`
- User Service: `http://localhost:8084`
- Bidding Service: `http://localhost:8085`
- Order Service: `http://localhost:8086`
- Notification Service: `http://localhost:8088`
- Media Service: `http://localhost:8089`
- Search Service: `http://localhost:8090`
- Comment Service: `http://localhost:8091`

**Tất cả requests từ client phải đi qua API Gateway (port 8080).**

---

**Tài liệu này được tạo bởi Senior Backend Engineer - Online Auction System**

**Version:** 1.0  
**Last Updated:** December 27, 2025
