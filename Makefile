.PHONY: help build run test clean docker-build docker-up docker-down docker-dev docker-prod docker-logs docker-rebuild migrate fmt lint swagger

# Default target
help:
	@echo "Available targets:"
	@echo "  make build         - Build the application binary"
	@echo "  make run           - Run the application locally"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make swagger       - Generate Swagger documentation"
	@echo ""
	@echo "Docker commands:"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start services (standard mode)"
	@echo "  make docker-dev    - Start services (development mode with hot reload)"
	@echo "  make docker-prod   - Start services (production mode)"
	@echo "  make docker-down   - Stop services"
	@echo "  make docker-logs   - View API logs"
	@echo "  make docker-rebuild- Rebuild and restart services"
	@echo ""
	@echo "Development:"
	@echo "  make fmt           - Format Go code"
	@echo "  make lint          - Run linter"
	@echo "  make dev           - Run with hot reload (requires air)"
	@echo "  make migrate       - Run database migrations"

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

# Start services with docker-compose (standard mode)
docker-up:
	@echo "Starting services (standard mode)..."
	@docker-compose up -d
	@echo "Services started. API available at http://localhost:8080"
	@echo "Swagger UI: http://localhost:8080/swagger/index.html"

# Start services in development mode with hot reload
docker-dev:
	@echo "Starting services (development mode with hot reload)..."
	@docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
	@echo "Development mode - Changes will auto-reload"

# Start services in production mode
docker-prod:
	@echo "Starting services (production mode)..."
	@docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
	@echo "Production services started"

# Stop services
docker-down:
	@echo "Stopping services..."
	@docker-compose down
	@echo "Services stopped"

# Stop services and remove volumes
docker-clean:
	@echo "Stopping services and removing volumes..."
	@docker-compose down -v
	@echo "Services stopped and volumes removed"

# Rebuild and restart services
docker-rebuild:
	@echo "Rebuilding services..."
	@docker-compose build --no-cache
	@docker-compose up -d
	@echo "Services rebuilt and restarted"

# Restart services
docker-restart: docker-down docker-up

# View logs
docker-logs:
	@docker-compose logs -f api

# View all logs
docker-logs-all:
	@docker-compose logs -f

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
