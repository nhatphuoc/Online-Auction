# Bidding Service API Documentation

This document describes the API endpoints for the Bidding Service as accessed via the API Gateway. All endpoints require the `X-User-Token` header for authentication. The APIs are documented and accessible via Swagger UI.

- **API Gateway Base URL:** `http://localhost:8080`
- **Bidding Service Base Path:** `/api/bids`
- **Internal Service URL:** `http://localhost:8085/bids`

## Common Headers
- `X-User-Token: <JWT Token>`
- `Content-Type: application/json`

---

## Endpoints

### 1. Search Bidding History
- **Endpoint:** `GET /api/bids/search`
- **Description:** Search for bidding history based on criteria.
- **Query Parameters:**
  - `productId` (string, optional)
  - `userId` (string, optional)
  - `page` (integer, optional, default: 0)
  - `size` (integer, optional, default: 10)
- **Example:**
  `GET /api/bids/search?productId=abc123&userId=user1&page=0&size=10`
- **Response:**
```json
{
  "success": true,
  "data": {
    "content": [
      {
        "id": "string",
        "productId": "string",
        "userId": "string",
        "bidAmount": 1000,
        "bidTime": "2024-01-01T12:00:00Z"
      }
    ],
    "totalElements": 1,
    "totalPages": 1,
    "page": 0,
    "size": 10
  },
  "message": "Search successful"
}
```

### 2. Place a Bid
- **Endpoint:** `POST /api/bids`
- **Description:** Place a new bid on a product.
- **Request Body:**
```json
{
  "productId": "string",
  "bidAmount": 1000
}
```
- **Response (Success):**
```json
{
  "success": true,
  "data": {
    "productId": "string",
    "bidAmount": 1000,
    "userId": "string",
    "bidTime": "2024-01-01T12:00:00Z"
  },
  "message": "Bid placed successfully"
}
```
- **Response (Failure):**
```json
{
  "success": false,
  "data": {
    "reason": "Bid amount too low"
  },
  "message": "Bid failed"
}
```

### 3. Get Bidding History for a Product
- **Endpoint:** `GET /api/bids/product/{productId}`
- **Description:** Get all bidding history for a specific product.
- **Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "string",
      "productId": "string",
      "userId": "string",
      "bidAmount": 1000,
      "bidTime": "2024-01-01T12:00:00Z"
    }
  ],
  "message": "Fetched successfully"
}
```

---

## Notes
- All endpoints are protected and require a valid `X-User-Token`.
- Use Swagger UI for interactive API documentation and testing.
- All requests and responses are in JSON format.

---

## Swagger UI
- Access via: `http://localhost:8080/swagger-ui.html` (API Gateway)

---

## Contact
For further information, contact the backend team.
