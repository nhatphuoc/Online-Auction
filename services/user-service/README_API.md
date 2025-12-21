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
- **Header:** `Authorization: Bearer <token>`
- **Response:**
```json
{
  "content": [ ... ],
  "totalElements": 100,
  ...
}
```
