.PHONY: help build run test clean swagger docker-build docker-run

# Variables
APP_NAME=category-service
DOCKER_IMAGE=$(APP_NAME):latest
PORT=3000

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) cmd/main.go
	@echo "Build completed!"

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	@go run cmd/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "Clean completed!"

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@swag init -g cmd/main.go -o docs
	@echo "Swagger docs generated!"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built!"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -p $(PORT):$(PORT) --env-file .env $(DOCKER_IMAGE)

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy
	@echo "Go modules tidied!"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted!"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run
	@echo "Linting completed!"

all: clean swagger build ## Clean, generate swagger and build

dev: swagger run ## Generate swagger and run in development mode

seed: ## Run seed data script
	@echo "Running seed data..."
	@go run scripts/seed.go
	@echo "Seed data completed!"
