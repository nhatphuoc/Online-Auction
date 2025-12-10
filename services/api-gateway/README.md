# API Gateway

API Gateway cho hệ thống đấu giá trực tuyến. Gateway nhận requests từ UI và proxy đến các internal services.

## Kiến trúc

```
UI -> API Gateway -> Internal Services
       |
       +-> Auth Service (validate JWT)
       +-> Category Service
       +-> Product Service  
       +-> User Service
       +-> Bidding Service
       +-> Order Service
       +-> Payment Service
       +-> Notification Service
       +-> Media Service
```

## Flow xác thực

1. **UI gửi request đến Gateway** với header:
   ```
   X-User-Token: <jwt_token_without_bearer>
   ```

2. **Gateway validate token** qua Auth Service:
   - Gọi `POST /auth/validate-jwt`
   - Headers: `X-User-Token`, `X-Api-Gateway`, `X-Auth-Internal-Service`

3. **Gateway proxy request** đến internal service với headers:
   ```
   X-User-ID: <user_id>
   X-User-Email: <email>
   X-User-Role: <role>
   X-User-Token: <jwt_token>
   X-Api-Gateway: <gateway_secret>
   X-Auth-Internal-Service: <internal_secret>
   ```

## Headers

### From UI to Gateway
- `X-User-Token`: JWT token (không có Bearer prefix)

### From Gateway to Auth Service (validate)
- `X-User-Token`: JWT token
- `X-Api-Gateway`: Secret để xác thực gateway
- `X-Auth-Internal-Service`: Secret cho internal services

### From Gateway to Internal Services
- `X-User-ID`: User ID từ JWT
- `X-User-Email`: Email từ JWT
- `X-User-Role`: Role từ JWT
- `X-User-Token`: JWT token gốc
- `X-Api-Gateway`: Secret để xác thực gateway
- `X-Auth-Internal-Service`: Secret cho internal services

## Endpoints

### Public (không cần auth)
- `GET /health` - Health check
- `GET /swagger/*` - API documentation
- `POST /api/auth/*` - Auth endpoints (login, register, etc.)

### Protected (cần auth)
- `GET/POST/PUT/DELETE /api/categories/*` - Category service
- `GET/POST/PUT/DELETE /api/products/*` - Product service
- `GET/POST/PUT/DELETE /api/users/*` - User service
- `GET/POST/PUT/DELETE /api/bidding/*` - Bidding service
- `GET/POST/PUT/DELETE /api/orders/*` - Order service
- `GET/POST/PUT/DELETE /api/payments/*` - Payment service
- `GET/POST/PUT/DELETE /api/notifications/*` - Notification service
- `GET/POST /api/media/*` - Media service

## Configuration (.env)

```env
# Server
PORT=8080

# Security
API_GATEWAY_SECRET=your-gateway-secret
AUTH_INTERNAL_SECRET=your-internal-secret

# Service URLs
AUTH_SERVICE_URL=http://localhost:8081
CATEGORY_SERVICE_URL=http://localhost:8082
PRODUCT_SERVICE_URL=http://localhost:8083
USER_SERVICE_URL=http://localhost:8084
BIDDING_SERVICE_URL=http://localhost:8085
ORDER_SERVICE_URL=http://localhost:8086
PAYMENT_SERVICE_URL=http://localhost:8087
NOTIFICATION_SERVICE_URL=http://localhost:8088
MEDIA_SERVICE_URL=http://localhost:8089

# OpenTelemetry
OTEL_ENDPOINT=localhost:4317
OTEL_SERVICE_NAME=api-gateway
OTEL_SERVICE_VERSION=1.0.0
OTEL_ENVIRONMENT=development
```

## Run

```bash
# Development
go run cmd/main.go

# Build
go build -o bin/api-gateway cmd/main.go

# Run binary
./bin/api-gateway
```

## Middleware

### AuthMiddleware
- Validate JWT token via Auth Service
- Set user info in context (userID, email, role)

### ProxyMiddleware
- Add required headers for internal services
- Forward user information

### TracingMiddleware
- OpenTelemetry distributed tracing

## Example Request Flow

### 1. User Login (Public)
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```

Response:
```json
{
  "token": "eyJhbGc...",
  "user": {...}
}
```

### 2. Get Categories (Protected)
```bash
curl -X GET http://localhost:8080/api/categories \
  -H "X-User-Token: eyJhbGc..."
```

Flow:
1. Gateway receives request
2. AuthMiddleware validates token with Auth Service
3. ProxyMiddleware adds headers
4. Request proxied to Category Service
5. Response returned to UI

## Security

- Tất cả internal services phải verify headers:
  - `X-Api-Gateway`: Verify request từ gateway
  - `X-Auth-Internal-Service`: Verify internal authentication
  
- JWT token được validate tại Auth Service
- Gateway không decode/validate JWT trực tiếp
- User info được lấy từ Auth Service response

## Observability

- **Logging**: Structured logging với slog
- **Metrics**: OpenTelemetry metrics
- **Tracing**: Distributed tracing qua tất cả services

## Technologies

- **Fiber v2**: Web framework
- **OpenTelemetry**: Observability (tracing, metrics, logging)
- **Go**: Programming language
