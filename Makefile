.PHONY: help start build run test clean docker-build docker-up docker-down migrate migrate-up migrate-down migrate-status migrate-version migrate-force migrate-create fmt lint swagger

# Default target
help:
	@echo "Available targets:"
	@echo "  make start         - 🚀 Start the entire project (Docker + open Swagger)"
	@echo "  make build         - Build the application binary"
	@echo "  make run           - Run the application locally"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make swagger       - Generate Swagger documentation"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start services with docker-compose"
	@echo "  make docker-down   - Stop services with docker-compose"
	@echo "  make fmt           - Format Go code"
	@echo "  make lint          - Run linter"
	@echo ""
	@echo "Database Migration Commands:"
	@echo "  make migrate-up      - Apply all pending migrations"
	@echo "  make migrate-down    - Rollback the last migration"
	@echo "  make migrate-status  - Show migration status"
	@echo "  make migrate-version - Show current migration version"
	@echo "  make migrate-create  - Create a new migration file"
	@echo "  make migrate-force   - Force set migration version (emergency use)"

# Start the entire project
start:
	@echo "🚀 Starting Learning Management System..."
	@echo ""
	@if [ ! -f .env ]; then \
		echo "⚙️  Creating .env file from .env.example..."; \
		cp .env.example .env; \
		echo "✅ .env file created!"; \
		echo ""; \
	fi
	@echo "📦 Starting Docker containers..."
	@docker-compose up -d
	@echo ""
	@echo "⏳ Waiting for services to be ready..."
	@sleep 5
	@echo ""
	@echo "🔍 Checking API health..."
	@until curl -s http://localhost:8080/api/health > /dev/null 2>&1; do \
		echo "   Waiting for API to be ready..."; \
		sleep 2; \
	done
	@echo ""
	@echo "✅ All services are running!"
	@echo ""
	@echo "📚 API Documentation: http://localhost:8080/swagger/index.html"
	@echo "🏥 Health Check:      http://localhost:8080/api/health"
	@echo ""
	@echo "Opening Swagger UI in your browser..."
	@sleep 1
	@open http://localhost:8080/swagger/index.html 2>/dev/null || xdg-open http://localhost:8080/swagger/index.html 2>/dev/null || echo "Please open http://localhost:8080/swagger/index.html in your browser"
	@echo ""
	@echo "💡 Useful commands:"
	@echo "   make docker-logs  - View API logs"
	@echo "   make docker-down  - Stop all services"
	@echo ""

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
migrate-up:
	@echo "Running migrations..."
	@go run src/cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back last migration..."
	@go run src/cmd/migrate/main.go down

migrate-status:
	@echo "Checking migration status..."
	@go run src/cmd/migrate/main.go status

migrate-version:
	@echo "Getting migration version..."
	@go run src/cmd/migrate/main.go version

migrate-force:
	@echo "Force migration version (use with caution)..."
	@read -p "Enter version number: " version; \
	go run src/cmd/migrate/main.go force $$version

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name; \
	echo "Migration files created in migrations/"

# Legacy migrate command (for backwards compatibility)
migrate: migrate-up

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
