# Quick Reference Guide

## 🚀 Getting Started

### Start the Application
```bash
make start                  # Complete startup with Swagger UI
docker-compose up -d        # Start without opening browser
make run                    # Run locally (requires PostgreSQL)
```

### Stop the Application
```bash
docker-compose down         # Stop all services
docker-compose down -v      # Stop and remove volumes (clean slate)
```

## 📊 Database Migrations

### Common Commands
```bash
make migrate-up             # Apply all pending migrations
make migrate-down           # Rollback last migration
make migrate-status         # Check migration status
make migrate-version        # Show current version
make migrate-create         # Create new migration files
```

### Manual Migration Tool
```bash
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/lms_db?sslmode=disable"
go run src/cmd/migrate/main.go up        # Apply migrations
go run src/cmd/migrate/main.go status    # Check status
go run src/cmd/migrate/main.go down      # Rollback
```

## 🐳 Docker Commands

### Service Management
```bash
docker-compose ps           # Check service status
docker-compose logs -f      # View all logs
docker-compose logs -f api  # View API logs only
docker-compose restart      # Restart all services
```

### Database Access
```bash
docker-compose exec postgres psql -U postgres -d lms_db
```

### Container Shell Access
```bash
docker-compose exec api sh              # Access API container
docker-compose exec api ./migrate up    # Run migrations in container
```

## 📚 API Documentation

### Access Points
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Swagger JSON**: http://localhost:8080/swagger/doc.json
- **Health Check**: http://localhost:8080/api/health

### Regenerate Swagger
```bash
make swagger
# or
swag init -g src/cmd/api/main.go -o docs --parseDependency --parseInternal
```

## 🔧 Development

### Build
```bash
make build                  # Build binaries
go build -o bin/api src/cmd/api/*.go
go build -o bin/migrate src/cmd/migrate/main.go
```

### Test
```bash
make test                   # Run all tests
go test -v ./...           # Verbose test output
go test -cover ./...       # With coverage
```

### Format & Lint
```bash
make fmt                    # Format code
make lint                   # Run linter
go mod tidy                 # Clean up dependencies
```

## 🔐 Environment Variables

### Required Variables
```bash
# Database
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/lms_db?sslmode=disable
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=lms_db

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION_HOURS=24

# Webhook
WEBHOOK_SECRET=your-webhook-secret
```

## 🧪 Testing API Endpoints

### Register User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123",
    "first_name": "Admin",
    "last_name": "User",
    "role": "admin"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

### Access Protected Endpoint
```bash
TOKEN="your-jwt-token"
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN"
```

### Health Check
```bash
curl http://localhost:8080/api/health
```

## 📖 Documentation Files

### Getting Started
- [README.md](README.md) - Main documentation
- [DOCKER_GUIDE.md](DOCKER_GUIDE.md) - Docker deployment

### Migration System
- [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) - Complete migration guide
- [MIGRATION_IMPLEMENTATION_SUMMARY.md](MIGRATION_IMPLEMENTATION_SUMMARY.md) - Technical details

### Features
- [RBAC.md](RBAC.md) - Role-based access control
- [AUDIT_LOGGING.md](AUDIT_LOGGING.md) - Audit logging system
- [CERTIFICATION_WEBHOOK.md](CERTIFICATION_WEBHOOK.md) - Webhook integration
- [BUSINESS_RULES_VALIDATION.md](BUSINESS_RULES_VALIDATION.md) - Business rules

### Testing & Compliance
- [RBAC_TEST_SPECIFICATION.md](RBAC_TEST_SPECIFICATION.md) - Test specifications
- [REQUIREMENTS_COMPLIANCE_REPORT.md](REQUIREMENTS_COMPLIANCE_REPORT.md) - 100% compliance

### Architecture
- [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - Project architecture
- [docs/SWAGGER_IMPLEMENTATION.md](docs/SWAGGER_IMPLEMENTATION.md) - API docs

## 🐛 Troubleshooting

### Port Already in Use
```bash
lsof -i :8080               # Check what's using port 8080
lsof -i :5432               # Check what's using port 5432
docker-compose down         # Stop services
```

### Database Connection Issues
```bash
docker-compose ps postgres                              # Check if running
docker-compose logs postgres                            # View logs
docker-compose exec postgres pg_isready -U postgres     # Test connection
```

### Migration Failures
```bash
make migrate-status         # Check current status
docker-compose down -v      # Clean slate (WARNING: deletes data)
docker-compose up -d        # Restart fresh
```

### Build Errors
```bash
go mod tidy                 # Fix dependencies
go clean -cache             # Clear Go cache
make clean                  # Clean build artifacts
```

### Docker Build Issues
```bash
docker-compose build --no-cache     # Rebuild without cache
docker system prune -a              # Clean Docker system
```

## 🎯 Common Tasks

### Reset Database
```bash
docker-compose down -v
docker-compose up -d
# Migrations run automatically
```

### Update Dependencies
```bash
go get -u ./...
go mod tidy
```

### Create New Migration
```bash
make migrate-create
# Enter migration name: add_new_feature
# Edit: migrations/000002_add_new_feature.up.sql
# Edit: migrations/000002_add_new_feature.down.sql
make migrate-up
```

### View All Routes
```bash
grep -r "@Router" src/internal/adapter/http/handler/ | grep -v Binary
```

### Check Code Coverage
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📊 Monitoring

### Check Service Health
```bash
curl http://localhost:8080/api/health
docker-compose ps
docker stats
```

### View Logs
```bash
docker-compose logs -f --tail=100           # Last 100 lines, follow
docker-compose logs api | grep ERROR        # Search for errors
docker-compose logs --since 10m             # Last 10 minutes
```

### Database Queries
```bash
docker-compose exec postgres psql -U postgres -d lms_db -c "SELECT * FROM schema_migrations;"
docker-compose exec postgres psql -U postgres -d lms_db -c "\dt"
```

## 🔗 Useful Links

- **Local API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/api/health
- **PostgreSQL**: localhost:5432

## 💡 Tips

1. **Use `make start`** for quick development setup
2. **Check `make help`** for all available commands
3. **Migrations run automatically** on application startup
4. **Use Swagger UI** for interactive API testing
5. **Check logs first** when troubleshooting issues
6. **Clean slate**: `docker-compose down -v && docker-compose up -d`

---

**Last Updated**: December 17, 2025  
**Quick Help**: Run `make help` for all commands
