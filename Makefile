.PHONY: help start build run test test-all test-workflow test-integration clean docker-build docker-up docker-down migrate migrate-up migrate-down migrate-status migrate-version migrate-force migrate-create fmt lint swagger build-all health db-shell db-backup db-restore db-tables db-reset migrate-status-docker docker-restart docker-logs dev deps

# Default target
help:
	@echo "🚀 Learning Management System - Available Commands"
	@echo ""
	@echo "Quick Start:"
	@echo "  make start              - 🚀 Start everything (Docker + open Swagger)"
	@echo "  make docker-down        - Stop all services"
	@echo ""
	@echo "Build & Run:"
	@echo "  make build              - Build API binary"
	@echo "  make build-all          - Build API and migration binaries"
	@echo "  make run                - Run application locally"
	@echo "  make dev                - Run with hot reload (requires air)"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-up          - Start services with docker-compose"
	@echo "  make docker-down        - Stop services"
	@echo "  make docker-restart     - Restart services"
	@echo "  make docker-logs        - View API logs"
	@echo ""
	@echo "Testing:"
	@echo "  make test               - Run unit tests (excludes DB adapters)"
	@echo "  make test-all           - Run all tests (includes DB adapters)"
	@echo "  make test-workflow      - Run full integration workflow"
	@echo "  make test-integration   - Full end-to-end test (requires Docker)"
	@echo ""
	@echo "Database Migrations:"
	@echo "  make migrate-up         - Apply all pending migrations"
	@echo "  make migrate-down       - Rollback last migration"
	@echo "  make migrate-status     - Show migration status"
	@echo "  make migrate-version    - Show current version"
	@echo "  make migrate-create     - Create new migration files"
	@echo "  make migrate-force      - Force version (emergency)"
	@echo ""
	@echo "Database Management:"
	@echo "  make db-shell           - Open PostgreSQL shell"
	@echo "  make db-tables          - List database tables"
	@echo "  make db-backup          - Create database backup"
	@echo "  make db-restore FILE=x  - Restore from backup file"
	@echo "  make db-reset           - ⚠️  Reset database (deletes all data)"
	@echo ""
	@echo "Development:"
	@echo "  make test               - Run tests with coverage"
	@echo "  make fmt                - Format code"
	@echo "  make lint               - Run linter"
	@echo "  make deps               - Install/update dependencies"
	@echo "  make swagger            - Generate Swagger docs"
	@echo "  make clean              - Remove build artifacts"
	@echo ""
	@echo "Monitoring:"
	@echo "  make health             - Check API health"
	@echo "  make migrate-status-docker - Check migration status in Docker"
	@echo ""
	@echo "📚 Documentation: See README.md and QUICK_REFERENCE.md"

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
	@echo "Running tests (excluding repository adapters)..."
	@pkgs=$$(go list ./... | grep -v src/internal/infrastructure/repository); \
	go test -v -race -coverprofile=coverage.out $$pkgs; \
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run all packages (including repository adapters)
test-all:
	@echo "Running all tests (including repository adapters)..."
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
	@docker-compose exec postgres psql -U postgres -d lms_db

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g src/cmd/api/main.go -o docs --parseDependency --parseInternal
	@echo "Swagger docs generated in docs/ directory"
	@echo "View at: http://localhost:8080/swagger/index.html"

# Database backup
db-backup:
	@echo "Creating database backup..."
	@docker-compose exec postgres pg_dump -U postgres lms_db > backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "Backup created: backup_$$(date +%Y%m%d_%H%M%S).sql"

# Database restore (usage: make db-restore FILE=backup.sql)
db-restore:
	@if [ -z "$(FILE)" ]; then \
		echo "Error: Please specify FILE parameter"; \
		echo "Usage: make db-restore FILE=backup.sql"; \
		exit 1; \
	fi
	@echo "Restoring database from $(FILE)..."
	@docker-compose exec -T postgres psql -U postgres lms_db < $(FILE)
	@echo "Database restored"

# Build both API and migration binaries
build-all:
	@echo "Building all binaries..."
	@mkdir -p bin
	@go build -o bin/api ./src/cmd/api
	@go build -o bin/migrate ./src/cmd/migrate
	@echo "✅ Binaries built:"
	@echo "   - bin/api"
	@echo "   - bin/migrate"

# Check application health
health:
	@echo "Checking application health..."
	@curl -s http://localhost:8080/api/health | jq . || echo "API not responding"

# View database tables
db-tables:
	@echo "Database tables:"
	@docker-compose exec postgres psql -U postgres -d lms_db -c "\dt"

# Check migration status via Docker
migrate-status-docker:
	@echo "Checking migration status in Docker..."
	@docker-compose exec api ./migrate status

# Reset database (WARNING: Deletes all data)
db-reset:
	@echo "⚠️  WARNING: This will delete ALL data!"
	@read -p "Are you sure? (yes/no): " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "Resetting database..."; \
		docker-compose down -v; \
		docker-compose up -d; \
		sleep 5; \
		echo "Database reset complete"; \
	fi

# Full integration workflow test
test-workflow:
	@echo "🧪 Running full integration workflow test..."
	@chmod +x test_workflow.sh
	@./test_workflow.sh

# End-to-end integration test (requires Docker)
test-integration:
	@echo "🧪 Running end-to-end integration tests..."
	@echo "Ensuring services are running..."
	@docker-compose ps | grep -q "api" || docker-compose up -d
	@sleep 3
	@chmod +x test_workflow.sh
	@./test_workflow.sh
	@echo ""
	@echo "✅ Integration test complete!"
		docker-compose up -d; \
		echo "✅ Database reset complete"; \
	else \
		echo "❌ Database reset cancelled"; \
	fi
