# Product Service API Documentation

This document describes the Product Service endpoints as exposed via the API Gateway. All requests must include the `X-User-Token` header for authentication. The APIs are documented in Swagger and are accessible through the API Gateway at `http://localhost:8080` (which proxies to the Product Service at `http://localhost:8083`).

## Base URL

```
http://localhost:8080/api/products
```

## Common Headers

- `X-User-Token`: string (required)

---

## Endpoints

### 1. Get Top Most Bids
- **Endpoint:** `GET /api/products/top-most-bids`
- **Description:** Get products with the highest number of bids.
- **Request:**
  - Headers: `X-User-Token`
- **Response:**
  - Status: 200 OK
  - Body: `List<ProductDTO>`

### 2. Get Product By ID
- **Endpoint:** `GET /api/products/{id}`
- **Description:** Get product details by product ID.
- **Request:**
  - Path: `id` (Long)
  - Headers: `X-User-Token`
- **Response:**
  - Status: 200 OK
  - Body: `ProductDTO`

### 3. Create Product
- **Endpoint:** `POST /api/products`
- **Description:** Create a new product.
- **Request:**
  - Headers: `X-User-Token`
  - Body: `ProductCreateRequest`
- **Response:**
  - Status: 201 Created
  - Body: `ProductDTO`

### 4. Update Product
- **Endpoint:** `PUT /api/products/{id}`
- **Description:** Update an existing product.
- **Request:**
  - Path: `id` (Long)
  - Headers: `X-User-Token`
  - Body: `ProductUpdateRequest`
- **Response:**
  - Status: 200 OK
  - Body: `ProductDTO`

### 5. Delete Product
- **Endpoint:** `DELETE /api/products/{id}`
- **Description:** Delete a product by ID.
- **Request:**
  - Path: `id` (Long)
  - Headers: `X-User-Token`
- **Response:**
  - Status: 200 OK
  - Body: `ApiResponse`

### 6. List Products (with filter)
- **Endpoint:** `GET /api/products`
- **Description:** List products with optional filters (category, search, etc.).
- **Request:**
  - Query params: `categoryId`, `search`, `page`, `size`, ...
  - Headers: `X-User-Token`
- **Response:**
  - Status: 200 OK
  - Body: `Page<ProductDTO>`

### 7. Batch Update Products
- **Endpoint:** `PUT /api/products/batch-update`
- **Description:** Batch update multiple products.
- **Request:**
  - Headers: `X-User-Token`
  - Body: `List<ProductUpdateRequest>`
- **Response:**
  - Status: 200 OK
  - Body: `BatchUpdateResult`

---

## Swagger UI

- Access API documentation at: `http://localhost:8080/swagger-ui.html`

## Notes
- All endpoints require `X-User-Token` header.
- All requests are routed through the API Gateway.
- For detailed request/response models, refer to Swagger UI.

---

## Example cURL Request

```
curl -H "X-User-Token: <token>" http://localhost:8080/api/products/top-most-bids
```

---

## Contact
For further information, contact the backend team.
