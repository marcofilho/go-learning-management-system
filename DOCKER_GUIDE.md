# Docker Deployment Guide

## Overview

The Learning Management System uses Docker and Docker Compose for containerized deployment. This guide covers building, running, and managing the application in Docker.

## Architecture

```
┌─────────────────────────────────────────┐
│          Docker Network (lms-network)   │
│                                         │
│  ┌──────────────┐    ┌──────────────┐  │
│  │              │    │              │  │
│  │  PostgreSQL  │◄───┤   LMS API    │  │
│  │  Container   │    │  Container   │  │
│  │              │    │              │  │
│  └──────────────┘    └──────────────┘  │
│       :5432               :8080         │
└─────────────────────────────────────────┘
```

## Services

### 1. PostgreSQL Database
- **Image**: `postgres:15-alpine`
- **Container Name**: `lms-postgres`
- **Port**: `5432`
- **Volume**: `postgres_data` (persistent storage)
- **Health Check**: Monitors database readiness

### 2. LMS API
- **Image**: Built from `Dockerfile`
- **Container Name**: `lms-api`
- **Port**: `8080`
- **Dependencies**: Waits for PostgreSQL to be healthy
- **Auto-restart**: `unless-stopped`

## Quick Start

### 1. Prerequisites

```bash
# Check Docker installation
docker --version
docker-compose --version

# Install if needed
brew install docker docker-compose  # macOS
```

### 2. Start Everything

```bash
# Copy environment file
cp .env.example .env

# Edit .env with your settings
nano .env

# Start all services (recommended)
make start

# Or manually
docker-compose up -d
```

### 3. Verify Deployment

```bash
# Check service status
docker-compose ps

# View logs
docker-compose logs -f api

# Check health
curl http://localhost:8080/api/health
```

### 4. Access the Application

- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/api/health

## Docker Commands

### Service Management

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# Restart services
docker-compose restart

# View logs
docker-compose logs -f           # All services
docker-compose logs -f api       # API only
docker-compose logs -f postgres  # Database only

# Check status
docker-compose ps

# View resource usage
docker stats
```

### Database Management

```bash
# Access PostgreSQL shell
docker-compose exec postgres psql -U postgres -d lms_db

# Create database backup
docker-compose exec postgres pg_dump -U postgres lms_db > backup.sql

# Restore database
docker-compose exec -T postgres psql -U postgres lms_db < backup.sql

# View database logs
docker-compose logs postgres
```

### API Container Management

```bash
# Access API container shell
docker-compose exec api sh

# Run migrations manually
docker-compose exec api ./migrate up

# Check migration status
docker-compose exec api ./migrate status

# View API logs
docker-compose logs -f api
```

## Building Images

### Development Build

```bash
# Build without cache
docker-compose build --no-cache

# Build specific service
docker-compose build api
```

### Production Build

```bash
# Multi-stage build with optimizations
docker build -t lms-api:latest .

# Build with specific version
docker build -t lms-api:v1.0.0 .

# Build for specific platform
docker build --platform linux/amd64 -t lms-api:latest .
```

## Environment Configuration

### Required Environment Variables

```bash
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-secure-password
DB_NAME=lms_db
DATABASE_URL=postgresql://postgres:your-secure-password@postgres:5432/lms_db?sslmode=disable

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRATION_HOURS=24

# Webhook
WEBHOOK_SECRET=your-webhook-secret
```

### .env File Structure

```bash
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=lms_db
DATABASE_URL=postgresql://postgres:postgres@postgres:5432/lms_db?sslmode=disable

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# JWT Configuration
JWT_SECRET=your-secret-key
JWT_EXPIRATION_HOURS=24

# Webhook Configuration
WEBHOOK_SECRET=your-webhook-secret-change-in-production
```

## Dockerfile Explanation

```dockerfile
# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build binaries
COPY . .
RUN go build -o api ./src/cmd/api
RUN go build -o migrate ./src/cmd/migrate

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app

# Copy artifacts
COPY --from=builder /app/api .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./api"]
```

**Benefits**:
- ✅ Multi-stage build (smaller image)
- ✅ Cached dependencies layer
- ✅ Includes migration tool
- ✅ Includes migration files
- ✅ Minimal runtime image

## docker-compose.yml Explanation

### PostgreSQL Service

```yaml
postgres:
  image: postgres:15-alpine
  container_name: lms-postgres
  environment:
    POSTGRES_USER: ${DB_USER}
    POSTGRES_PASSWORD: ${DB_PASSWORD}
    POSTGRES_DB: ${DB_NAME}
  ports:
    - "${DB_PORT}:5432"
  volumes:
    - postgres_data:/var/lib/postgresql/data
  healthcheck:
    test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
    interval: 10s
    timeout: 5s
    retries: 5
  networks:
    - lms-network
```

**Features**:
- Persistent data volume
- Health check monitoring
- Environment-based configuration
- Network isolation

### API Service

```yaml
api:
  build:
    context: .
    dockerfile: Dockerfile
  container_name: lms-api
  environment:
    DB_HOST: postgres
    DB_PORT: 5432
    DATABASE_URL: postgresql://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable
    SERVER_HOST: 0.0.0.0
    SERVER_PORT: 8080
    JWT_SECRET: ${JWT_SECRET}
    WEBHOOK_SECRET: ${WEBHOOK_SECRET}
  ports:
    - "${SERVER_PORT}:8080"
  depends_on:
    postgres:
      condition: service_healthy
  networks:
    - lms-network
  restart: unless-stopped
```

**Features**:
- Waits for database health check
- Auto-restart on failure
- Environment variable injection
- Network isolation

## Database Migrations in Docker

### Automatic Migration (Default)

The API automatically runs migrations on startup:

```go
// In main.go
if err := runMigrations(&cfg.Database); err != nil {
    log.Printf("⚠️  Migration warning: %v", err)
}
```

### Manual Migration

```bash
# Run migrations
docker-compose exec api ./migrate up

# Check status
docker-compose exec api ./migrate status

# Rollback
docker-compose exec api ./migrate down

# Force version
docker-compose exec api ./migrate force 1
```

## Data Persistence

### Volumes

```bash
# List volumes
docker volume ls | grep lms

# Inspect volume
docker volume inspect go-learning-management-system_postgres_data

# Backup volume
docker run --rm -v go-learning-management-system_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz -C /data .

# Restore volume
docker run --rm -v go-learning-management-system_postgres_data:/data -v $(pwd):/backup alpine tar xzf /backup/postgres_backup.tar.gz -C /data
```

### Clean Start

```bash
# Stop and remove everything
docker-compose down -v

# Remove all data (WARNING: DESTRUCTIVE)
docker volume rm go-learning-management-system_postgres_data

# Start fresh
docker-compose up -d
```

## Networking

### Internal Communication

```bash
# API connects to database via service name
DATABASE_URL=postgresql://postgres:password@postgres:5432/lms_db

# Not needed: localhost or 127.0.0.1
```

### External Access

```bash
# Access from host machine
DATABASE_URL=postgresql://postgres:password@localhost:5432/lms_db
API_URL=http://localhost:8080
```

### Network Inspection

```bash
# List networks
docker network ls

# Inspect network
docker network inspect go-learning-management-system_lms-network

# View connected containers
docker network inspect go-learning-management-system_lms-network | jq '.[0].Containers'
```

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker-compose logs api

# Check if port is in use
lsof -i :8080

# Remove and recreate
docker-compose down
docker-compose up -d
```

### Database Connection Issues

```bash
# Check database status
docker-compose ps postgres

# Test connection
docker-compose exec postgres psql -U postgres -d lms_db

# Check network
docker-compose exec api ping postgres
```

### Migration Failures

```bash
# Check migration status
docker-compose exec api ./migrate status

# View migration logs
docker-compose logs api | grep migration

# Force clean state (CAUTION)
docker-compose down -v
docker-compose up -d
```

### Performance Issues

```bash
# Check resource usage
docker stats

# Increase Docker resources
# Docker Desktop -> Settings -> Resources

# View container resource limits
docker inspect lms-api | jq '.[0].HostConfig.Memory'
```

## Production Deployment

### Security Hardening

```bash
# Use secrets instead of .env
docker-compose -f docker-compose.prod.yml up -d

# Run as non-root user
USER 1000:1000

# Read-only filesystem
read_only: true

# Drop capabilities
cap_drop:
  - ALL
cap_add:
  - NET_BIND_SERVICE
```

### Health Monitoring

```bash
# Enable health checks
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8080/api/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

### Resource Limits

```yaml
api:
  deploy:
    resources:
      limits:
        cpus: '1.0'
        memory: 512M
      reservations:
        cpus: '0.5'
        memory: 256M
```

### Logging

```yaml
api:
  logging:
    driver: "json-file"
    options:
      max-size: "10m"
      max-file: "3"
```

## CI/CD Integration

### GitHub Actions Example

```yaml
- name: Build Docker Image
  run: docker build -t lms-api:${{ github.sha }} .

- name: Push to Registry
  run: |
    docker tag lms-api:${{ github.sha }} registry.example.com/lms-api:latest
    docker push registry.example.com/lms-api:latest
```

### Automated Deployment

```bash
# Pull latest image
docker-compose pull

# Restart with new image
docker-compose up -d

# Run migrations
docker-compose exec api ./migrate up
```

## Monitoring

### Container Logs

```bash
# Real-time logs
docker-compose logs -f --tail=100

# Export logs
docker-compose logs > logs.txt

# Search logs
docker-compose logs | grep ERROR
```

### Health Checks

```bash
# API health
curl http://localhost:8080/api/health

# Database health
docker-compose exec postgres pg_isready -U postgres

# Container health
docker inspect --format='{{.State.Health.Status}}' lms-api
```

## References

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [PostgreSQL Docker Hub](https://hub.docker.com/_/postgres)
- [Alpine Linux](https://alpinelinux.org/)

---

**Last Updated**: December 17, 2025  
**Docker Version**: 20.10+  
**Docker Compose Version**: 2.0+
