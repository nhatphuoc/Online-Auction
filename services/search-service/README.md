# Search Service - Docker Setup

## 🚀 Quick Start

### 1. Start all services với Docker Compose:

```bash
# Build và start tất cả services
make up

# Hoặc dùng docker-compose trực tiếp
docker-compose up -d
```

### 2. Kiểm tra services đã chạy chưa:

```bash
# Check status
docker-compose ps

# Check logs
make logs

# Hoặc xem log từng service
make logs-search
make logs-es
```

### 3. Test API:

```bash
# Health check
curl http://localhost:3000/health

# Search products
curl "http://localhost:3000/api/search/products?query=test&page=1&page_size=10"
```

## 📦 Available Services

| Service | Port | URL |
|---------|------|-----|
| Search Service | 3000 | http://localhost:3000 |
| PostgreSQL | 5432 | localhost:5432 |
| Elasticsearch | 9200 | http://localhost:9200 |
| Redis | 6379 | localhost:6379 |

## 🛠️ Makefile Commands

```bash
make help           # Show all available commands
make build          # Build Docker images
make up             # Start all services
make down           # Stop all services
make logs           # Show logs from all services
make logs-search    # Show logs from search-service only
make restart        # Restart all services
make clean          # Stop and remove all containers and volumes
make rebuild        # Rebuild from scratch
make dev            # Run with logs visible
make test           # Test the API endpoints
make ps             # Show running containers
make stats          # Show container resource usage
```

## 🔧 Development

### Run without Docker (local development):

```bash
# Install dependencies
go mod download

# Start only infrastructure services
docker-compose up -d postgres elasticsearch redis

# Run application locally
go run cmd/main.go
```

### Access service shells:

```bash
# Search service shell
make shell-search

# PostgreSQL shell
make shell-postgres

# Redis CLI
make shell-redis
```

## 📝 Environment Variables

All environment variables are configured in `docker-compose.yml`. To override:

1. Create `.env.local` file
2. Add your custom variables
3. Update `docker-compose.yml` to use `env_file`

```yaml
search-service:
  env_file:
    - .env
    - .env.local
```

## 🔍 Troubleshooting

### Elasticsearch không start:

```bash
# Increase vm.max_map_count on Linux
sudo sysctl -w vm.max_map_count=262144

# Hoặc add vào /etc/sysctl.conf
vm.max_map_count=262144
```

### Search service không connect được database:

```bash
# Check logs
make logs-search

# Restart service
make restart-search

# Check database connectivity
make shell-postgres
```

### Clean start (xóa hết data):

```bash
# Stop và xóa tất cả volumes
make clean

# Start lại
make up
```

## 📊 Monitoring

### Check resource usage:

```bash
make stats
```

### View Elasticsearch indices:

```bash
curl http://localhost:9200/_cat/indices?v
```

### View Redis keys:

```bash
make shell-redis
# In Redis CLI:
KEYS *
```

## 🔄 Update & Rebuild

```bash
# After code changes
make rebuild

# Or manual steps
docker-compose build search-service
docker-compose up -d search-service
```

## 🧪 Testing

```bash
# Run automated tests
make test

# Manual testing
curl http://localhost:3000/health
curl "http://localhost:3000/api/search/products?query=laptop&page=1"
```

## 🌐 Production Deployment

For production, update `docker-compose.yml`:

1. Remove exposed ports for internal services
2. Add resource limits
3. Use Docker secrets for sensitive data
4. Enable Elasticsearch security
5. Add Redis password
6. Use production-grade PostgreSQL setup

Example:
```yaml
search-service:
  deploy:
    resources:
      limits:
        cpus: '1'
        memory: 512M
      reservations:
        cpus: '0.5'
        memory: 256M
```
