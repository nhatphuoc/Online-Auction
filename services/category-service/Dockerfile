# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o category-service cmd/main.go

# Runtime stage
FROM alpine:3.20

WORKDIR /app

# Install ca-certificates for HTTPS and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/category-service .

# Copy .env file (optional, can be overridden by environment variables)
COPY .env .env

# Expose port
EXPOSE 3000

# Run the application
CMD ["./category-service"]
