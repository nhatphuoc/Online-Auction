# API Gateway Documentation

## Tổng quan flow

- **Client gửi request** tới API Gateway (qua các endpoint /api/...)
- **API Gateway xác thực JWT** (header: `X-User-Token`) bằng public key của Auth Service
- Nếu hợp lệ, API Gateway **proxy request** tới các internal service
- Khi proxy, API Gateway **ký JWT nội bộ** bằng private key của chính nó, gửi qua header `X-Internal-JWT` cho các service nội bộ

---

## Header Input (Client → API Gateway)
- `X-User-Token`: JWT của user (do Auth Service cấp)
- `Authorization`: (tuỳ chọn, nếu dùng Bearer)
- `Content-Type`: application/json

## Header Output (API Gateway → Internal Service)
- `X-User-ID`: ID của user
- `X-User-Email`: Email user
- `X-User-Role`: Role user
- `X-User-Token`: JWT của user
- `X-Api-Gateway`: Secret của API Gateway
- `X-Auth-Internal-Service`: Secret nội bộ
- `X-Internal-JWT`: JWT nội bộ do API Gateway ký (RS256, chứa iss, aud, exp, sub)

---

## Endpoint

### Public
- `GET /health`: Kiểm tra trạng thái API Gateway
- `GET /swagger/*`: Swagger UI

### Auth Service
- `POST /api/auth/register`: Đăng ký
- `POST /api/auth/sign-in`: Đăng nhập
- `POST /api/auth/verify-otp`: Xác thực OTP
- `POST /api/auth/validate-jwt`: Kiểm tra JWT
- `POST /api/auth/sign-in/google`: Đăng nhập Google

### Protected (yêu cầu JWT hợp lệ)
- `ALL /api/categories/*`: Category Service
- `ALL /api/products/*`: Product Service
- `ALL /api/users/*`: User Service
- `ALL /api/bidding/*`: Bidding Service
- `ALL /api/orders/*`: Order Service
- `ALL /api/payments/*`: Payment Service
- `ALL /api/notifications/*`: Notification Service
- `ALL /api/media/*`: Media Service

---

## JWT nội bộ (X-Internal-JWT)
- **Algorithm**: RS256
- **Payload**:
  - `iss`: "api-gateway"
  - `aud`: "internal-service"
  - `exp`: Thời gian hết hạn (5 phút)
  - `sub`: userID

---

## Response mẫu
```json
{
  "success": true,
  "data": ...,
  "message": "..."
}
```

---

## Lưu ý
- Các service nội bộ cần xác thực JWT nội bộ qua header `X-Internal-JWT`.
- Private key của API Gateway cần bảo mật tuyệt đối.
