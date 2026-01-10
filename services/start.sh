#!/bin/bash
set -e

echo "Starting services..."

GO_SERVICES=(
  api-gateway
  auto-bidding-service
  category-service
  media-service
  order-service
  comment-service
)

JAVA_SERVICES=(
  auth-service
  bidding-service
  notification-service
  product-service
  user-service
)

PIDS=()

cleanup() {
  echo ""
  echo "Stopping all services..."

  for pid in "${PIDS[@]}"; do
    if ps -p "$pid" > /dev/null 2>&1; then
      echo "Stopping process group $pid"
      kill -TERM -- -"$pid" 2>/dev/null || true
    fi
  done

  echo "All services stopped."
  exit 0
}

trap cleanup SIGINT SIGTERM

echo "Starting Go services..."
for service in "${GO_SERVICES[@]}"; do
  echo "Starting $service"
  (
    cd "$service" || exit
    go run cmd/main.go
  ) &
  PIDS+=($!)
done

echo "Starting Spring Boot services..."
for service in "${JAVA_SERVICES[@]}"; do
  echo "Starting $service"
  (
    cd "$service" || exit
    ./mvnw spring-boot:run
  ) &
  PIDS+=($!)
done

echo "All services started."
wait
