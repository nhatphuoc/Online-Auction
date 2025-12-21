# Auth Service API Documentation

Dành cho Frontend: Danh sách các endpoint của Auth Service

Base URL: `/auth`

---

## 1. Đăng ký tài khoản
**POST** `/auth/register`
- **Request Body:**
```json
{
  "fullName": "string",
  "email": "string",
  "password": "string",
  "birthDay": "yyyy-MM-dd",
  "emailVerified": true/false
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

---

## 2. Xác thực OTP
**POST** `/auth/verify-otp`
- **Request Body:**
```json
{
  "email": "string",
  "otpCode": "string"
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

---

## 3. Đăng nhập
**POST** `/auth/sign-in`
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
  "accessToken": "string",
  "refreshToken": "string"
}
```

---

## 4. Xác thực JWT
**POST** `/auth/validate-jwt`
- **Request Body:**
```json
{
  "token": "string"
}
```
- **Response:**
```json
{
  "valid": true/false
}
```

---

## 5. Đăng nhập với Google
**POST** `/auth/sign-in/google`
- **Request Body:**
```json
{
  "idToken": "string"
}
```
- **Response:**
```json
{
  "success": true/false,
  "accessToken": "string",
  "refreshToken": "string"
}
```

---

## Response chuẩn
- Tất cả các response đều có thể trả về dạng:
```json
{
  "success": true/false,
  "message": "string",
  "data": ...
}
```

- Các trường hợp đặc biệt sẽ trả về các trường như ví dụ trên.

---

Nếu cần chi tiết hơn về từng trường, hãy xem lại các file DTO trong thư mục `src/main/java/com/Online_Auction/auth_service/dto/request` và `dto/response`.
