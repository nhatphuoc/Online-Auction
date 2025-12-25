# Bidding Service API Documentation

This document describes the endpoints for the Bidding Service, as accessed via the API Gateway. All requests must include the `X-User-Token` header for authentication. The APIs are documented and accessible via Swagger UI.

**Base URL via API Gateway:**
```
http://localhost:8080/api/bids
```

**Internal Service URL:**
```
http://localhost:8085/bids
```

---

## Endpoints

### 1. Search Bidding History
- **Endpoint:** `POST /api/bids/search`
- **Description:** Search for bidding history by product or user.
- **Headers:**
  - `X-User-Token: <token>`
- **Request Body:**
```json
{
  "productId": "string", // optional
  "userId": "string",    // optional
  "page": 0,              // integer, optional
  "size": 10              // integer, optional
}
```
- **Response:**
```json
{
  "success": true,
  "data": [
    {
      "bidId": "string",
      "productId": "string",
      "userId": "string",
      "amount": 100000,
      "timestamp": "2025-12-26T12:00:00Z"
    }
  ],
  "total": 1
}
```

---

### 2. Place a Bid
- **Endpoint:** `POST /api/bids`
- **Description:** Place a new bid on a product.
- **Headers:**
  - `X-User-Token: <token>`
- **Request Body:**
```json
{
  "productId": "string",
  "amount": 100000
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "bidId": "string",
    "productId": "string",
    "userId": "string",
    "amount": 100000,
    "timestamp": "2025-12-26T12:00:00Z"
  }
}
```

---

### 3. Get Bidding History for a Product
- **Endpoint:** `GET /api/bids/product/{productId}`
- **Description:** Get all bids for a specific product.
- **Headers:**
  - `X-User-Token: <token>`
- **Response:**
```json
{
  "success": true,
  "data": [
    {
      "bidId": "string",
      "productId": "string",
      "userId": "string",
      "amount": 100000,
      "timestamp": "2025-12-26T12:00:00Z"
    }
  ]
}
```

---

### 4. Get User's Bidding History
- **Endpoint:** `GET /api/bids/user/{userId}`
- **Description:** Get all bids placed by a specific user.
- **Headers:**
  - `X-User-Token: <token>`
- **Response:**
```json
{
  "success": true,
  "data": [
    {
      "bidId": "string",
      "productId": "string",
      "userId": "string",
      "amount": 100000,
      "timestamp": "2025-12-26T12:00:00Z"
    }
  ]
}
```

---

## Notes
- All endpoints require the `X-User-Token` header.
- All APIs are documented and testable via Swagger UI.
- All requests and responses are in JSON format.
- The API Gateway handles routing and authentication.

---

## Example: Call Search API via API Gateway
```
curl -X POST "http://localhost:8080/api/bids/search" \
  -H "Content-Type: application/json" \
  -H "X-User-Token: <token>" \
  -d '{"productId": "12345", "page": 0, "size": 10}'
```

---

## Contact
For more information, contact the backend team.
