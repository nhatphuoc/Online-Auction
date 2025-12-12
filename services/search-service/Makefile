.PHONY: help build up down logs restart clean dev test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all Docker images
	docker-compose build

up: ## Start all services
	docker-compose up -d
	@echo "Services are starting..."
	@echo "PostgreSQL:     http://localhost:5432"
	@echo "Elasticsearch:  http://localhost:9200"
	@echo "Redis:          http://localhost:6379"
	@echo "Search Service: http://localhost:3000"
	@echo ""
	@echo "Health check:   curl http://localhost:3000/health"

down: ## Stop all services
	docker-compose down

logs: ## Show logs from all services
	docker-compose logs -f

logs-search: ## Show logs from search-service only
	docker-compose logs -f search-service

logs-es: ## Show logs from Elasticsearch only
	docker-compose logs -f elasticsearch

restart: ## Restart all services
	docker-compose restart

restart-search: ## Restart search-service only
	docker-compose restart search-service

clean: ## Stop and remove all containers, volumes, and images
	docker-compose down -v
	docker-compose rm -f

rebuild: ## Rebuild and restart all services
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

dev: ## Run in development mode (with logs)
	docker-compose up

test: ## Test the search service
	@echo "Testing health endpoint..."
	@curl -s http://localhost:3000/health | jq .
	@echo ""
	@echo "Testing search endpoint..."
	@curl -s "http://localhost:3000/api/search/products?query=test&page=1&page_size=10" | jq .

ps: ## Show running containers
	docker-compose ps

shell-search: ## Open shell in search-service container
	docker-compose exec search-service sh

shell-postgres: ## Open PostgreSQL shell
	docker-compose exec postgres psql -U postgres -d auction_db

shell-redis: ## Open Redis CLI
	docker-compose exec redis redis-cli

init-indices: ## Initialize Elasticsearch indices (run after first start)
	@echo "Elasticsearch indices will be created automatically on first API call"
	@curl -s http://localhost:3000/health

stats: ## Show container stats
	docker stats search-service search-elasticsearch search-redis search-postgres
