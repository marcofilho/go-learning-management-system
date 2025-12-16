.PHONY: help build run test clean docker-build docker-up docker-down migrate fmt lint swagger

# Default target
help:
	@echo "Available targets:"
	@echo "  make build        - Build the application binary"
	@echo "  make run          - Run the application locally"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make swagger      - Generate Swagger documentation"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-up    - Start services with docker-compose"
	@echo "  make docker-down  - Stop services with docker-compose"
	@echo "  make fmt          - Format Go code"
	@echo "  make lint         - Run linter"
	@echo "  make migrate      - Run database migrations"

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/api ./src/cmd/api
	@echo "Build complete: bin/api"

# Run the application locally
run:
	@echo "Starting application..."
	@go run src/cmd/api/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t lms-api:latest .
	@echo "Docker image built: lms-api:latest"

# Start services with docker-compose
docker-up:
	@echo "Starting services..."
	@docker-compose up -d
	@echo "Services started. API available at http://localhost:8080"

# Stop services
docker-down:
	@echo "Stopping services..."
	@docker-compose down
	@echo "Services stopped"

# Restart services
docker-restart: docker-down docker-up

# View logs
docker-logs:
	@docker-compose logs -f api

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Run database migrations
migrate:
	@echo "Migrations run automatically on startup"
	@echo "To run manually, use: go run src/cmd/api/main.go"

# Development mode with hot reload (requires air)
dev:
	@echo "Starting development mode..."
	@air -c .air.toml || echo "air not installed. Run: go install github.com/cosmtrek/air@latest"

# Database shell
db-shell:
	@docker-compose exec postgres psql -U lms_user -d lms_db

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g src/cmd/api/main.go -o docs --parseDependency --parseInternal
	@echo "Swagger docs generated in docs/ directory"
	@echo "View at: http://localhost:8080/swagger/index.html"
