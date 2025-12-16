# Docker Configuration Guide

This guide explains the Docker setup for the Learning Management System API, including different configurations for development, testing, and production environments.

## Overview

The project uses Docker multi-stage builds and Docker Compose for containerization, with separate configurations optimized for different environments.

## Files Structure

```
.
├── Dockerfile              # Production-optimized build
├── Dockerfile.dev          # Development with hot reload
├── docker-compose.yml      # Base configuration
├── docker-compose.dev.yml  # Development overrides
├── docker-compose.prod.yml # Production overrides
├── .dockerignore          # Files excluded from build
└── .air.toml              # Hot reload configuration
```

## Dockerfile Configurations

### Production Dockerfile

**File:** `Dockerfile`

**Features:**
- Multi-stage build (builder + runtime)
- Go 1.24 Alpine base
- Optimized layer caching
- Static binary compilation
- Non-root user (appuser)
- Health check included
- Final image size: ~15MB

**Build optimizations:**
```dockerfile
# Separate go.mod/go.sum for better caching
COPY go.mod go.sum ./
RUN go mod download

# Static binary with stripped symbols
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"'
```

**Security features:**
- Runs as non-root user (UID 1000)
- Minimal Alpine Linux base
- No shell access needed
- CA certificates included for HTTPS

### Development Dockerfile

**File:** `Dockerfile.dev`

**Features:**
- Hot reload with Air
- Development tools included
- Volume mounting support
- Quick iteration cycles

**Usage:**
```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

## Docker Compose Configurations

### Base Configuration

**File:** `docker-compose.yml`

**Services:**
1. **PostgreSQL 15 Alpine**
   - Data persistence with named volume
   - Health check for startup coordination
   - UTF-8 encoding
   - Network isolation

2. **API Service**
   - Built from Dockerfile
   - Environment variables from .env
   - Health check endpoint
   - Depends on healthy PostgreSQL
   - Automatic restart

**Features:**
- Default values for all environment variables
- Health checks on both services
- Named volumes for data persistence
- Isolated network (lms-network)

### Development Overrides

**File:** `docker-compose.dev.yml`

**Additional features:**
- Hot reload enabled (Air)
- Source code volume mounting
- Direct PostgreSQL port access (5432)
- Debug logging enabled
- Optional pgAdmin service

**Usage:**
```bash
# Start development environment
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# With pgAdmin (uncomment in file first)
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```

**Development workflow:**
1. Code changes are detected by Air
2. Application automatically rebuilds
3. Server restarts with new changes
4. Zero downtime for development

### Production Overrides

**File:** `docker-compose.prod.yml`

**Production features:**
- Resource limits (CPU/Memory)
- Stricter health checks
- JSON file logging with rotation
- No exposed database ports
- Restart on failure (max 3 attempts)
- Release mode (no debug logs)

**Resource limits:**
- **API Service:**
  - CPU: 0.5-1 core
  - Memory: 256MB-512MB
- **PostgreSQL:**
  - CPU: 1-2 cores
  - Memory: 1GB-2GB

**Usage:**
```bash
# Start production environment
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

## Environment Variables

The services use environment variables from `.env` file with sensible defaults:

### Database Configuration
```env
DB_HOST=localhost          # Override to 'postgres' in docker-compose
DB_PORT=5432
DB_USER=postgres          # Default: postgres
DB_PASSWORD=postgres      # Change in production!
DB_NAME=lms_db           # Default: lms_db
```

### Server Configuration
```env
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

### JWT Configuration
```env
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION_HOURS=24
```

**⚠️ Production Notes:**
- Always change default passwords
- Use strong JWT secrets (32+ characters)
- Store secrets in secure vault (not in .env)

## Make Commands

Convenient commands for Docker operations:

```bash
# Standard operations
make docker-build      # Build production image
make docker-up         # Start services (standard mode)
make docker-down       # Stop services
make docker-logs       # View API logs
make docker-logs-all   # View all service logs

# Environment-specific
make docker-dev        # Start with hot reload
make docker-prod       # Start production mode

# Maintenance
make docker-rebuild    # Rebuild and restart
make docker-clean      # Stop and remove volumes
make docker-restart    # Quick restart
```

## Health Checks

### API Health Check
- **Endpoint:** `/api/health`
- **Interval:** 30 seconds
- **Timeout:** 3 seconds
- **Retries:** 3
- **Start period:** 5-10 seconds

### PostgreSQL Health Check
- **Command:** `pg_isready`
- **Interval:** 10 seconds
- **Timeout:** 5 seconds
- **Retries:** 5
- **Start period:** 10 seconds

**Check service health:**
```bash
docker-compose ps
docker inspect lms-api | grep Health -A 10
```

## Volume Management

### Named Volumes
- `postgres_data` - PostgreSQL data persistence

**Commands:**
```bash
# List volumes
docker volume ls

# Inspect volume
docker volume inspect go-learning-management-system_postgres_data

# Backup database
docker-compose exec postgres pg_dump -U postgres lms_db > backup.sql

# Restore database
docker-compose exec -T postgres psql -U postgres lms_db < backup.sql

# Remove volumes (⚠️ destroys data)
docker-compose down -v
```

## Networking

### Network Configuration
- **Name:** lms-network
- **Driver:** bridge
- **Isolation:** Container-only communication

**Access patterns:**
- API → PostgreSQL: Via service name `postgres`
- Host → API: Via `localhost:8080`
- Host → PostgreSQL (dev only): Via `localhost:5432`

**Inspect network:**
```bash
docker network inspect go-learning-management-system_lms-network
```

## Troubleshooting

### Container won't start

**Check logs:**
```bash
docker-compose logs api
docker-compose logs postgres
```

**Check health:**
```bash
docker-compose ps
```

**Restart services:**
```bash
docker-compose restart
```

### Database connection issues

**Verify PostgreSQL is healthy:**
```bash
docker-compose exec postgres pg_isready -U postgres
```

**Check environment variables:**
```bash
docker-compose exec api env | grep DB_
```

**Test connection:**
```bash
docker-compose exec api wget -qO- http://localhost:8080/api/health
```

### Port already in use

**Find process using port:**
```bash
# macOS/Linux
lsof -i :8080
lsof -i :5432

# Windows
netstat -ano | findstr :8080
```

**Change port in .env:**
```env
SERVER_PORT=8081
DB_PORT=5433
```

### Slow builds

**Enable BuildKit:**
```bash
export DOCKER_BUILDKIT=1
docker-compose build
```

**Use cache:**
```bash
docker-compose build --pull
```

**Clean build (if needed):**
```bash
docker-compose build --no-cache
```

### Out of disk space

**Remove unused resources:**
```bash
# Remove stopped containers
docker container prune

# Remove unused images
docker image prune -a

# Remove unused volumes (⚠️ careful)
docker volume prune

# Complete cleanup
docker system prune -a --volumes
```

## Performance Optimization

### Build Performance
1. Use `.dockerignore` to exclude unnecessary files
2. Order Dockerfile commands by change frequency
3. Enable BuildKit for parallel builds
4. Use layer caching effectively

### Runtime Performance
1. Use resource limits in production
2. Enable health checks for automatic recovery
3. Use named volumes instead of bind mounts in production
4. Configure appropriate restart policies

### Development Performance
1. Use volume mounts for hot reload
2. Exclude vendor and bin from mounts
3. Use docker-compose.dev.yml for development-specific settings

## Security Best Practices

### Image Security
- ✅ Non-root user
- ✅ Minimal base image (Alpine)
- ✅ No unnecessary packages
- ✅ Static binary (no runtime dependencies)
- ✅ Health checks enabled

### Runtime Security
- ✅ Network isolation
- ✅ No exposed database ports in production
- ✅ Environment variable validation
- ✅ Resource limits
- ✅ Read-only filesystem (where possible)

### Secrets Management
- ❌ Don't commit .env to git
- ✅ Use Docker secrets in swarm mode
- ✅ Use external secret managers (Vault, AWS Secrets Manager)
- ✅ Rotate credentials regularly

## Production Deployment

### Pre-deployment Checklist
- [ ] Change default passwords
- [ ] Use strong JWT secret
- [ ] Review resource limits
- [ ] Configure logging
- [ ] Set up monitoring
- [ ] Configure backups
- [ ] Test health checks
- [ ] Review security settings

### Deployment Commands
```bash
# Pull latest code
git pull origin main

# Build production image
docker-compose -f docker-compose.yml -f docker-compose.prod.yml build

# Start services
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Verify health
docker-compose ps
curl http://localhost:8080/api/health

# View logs
docker-compose logs -f api
```

### Rolling Updates
```bash
# Build new version
docker-compose build api

# Stop old version
docker-compose stop api

# Start new version
docker-compose up -d api

# Verify
docker-compose ps
curl http://localhost:8080/api/health
```

## Monitoring

### Container Metrics
```bash
# Resource usage
docker stats

# Inspect container
docker inspect lms-api

# Process list
docker-compose exec api ps aux
```

### Logs
```bash
# Follow logs
docker-compose logs -f api

# Last 100 lines
docker-compose logs --tail=100 api

# Since timestamp
docker-compose logs --since="2024-01-01T00:00:00" api
```

## Advanced Usage

### Multi-stage Development
```bash
# Development with hot reload
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

# Testing
docker-compose -f docker-compose.yml -f docker-compose.test.yml run --rm api go test ./...

# Production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### Database Operations
```bash
# Access PostgreSQL shell
make db-shell
# or
docker-compose exec postgres psql -U postgres -d lms_db

# Run migrations
docker-compose exec api ./api migrate

# Backup
docker-compose exec postgres pg_dump -U postgres lms_db > backup-$(date +%Y%m%d).sql

# Restore
docker-compose exec -T postgres psql -U postgres lms_db < backup.sql
```

## References

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose File Reference](https://docs.docker.com/compose/compose-file/)
- [Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [Dockerfile Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [Air (Hot Reload)](https://github.com/cosmtrek/air)
