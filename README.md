# Go Learning Management System

A Golang REST API for a Learning Management System, where students enroll in courses, instructors manage course content (modules and lessons with versioning), and an external certification provider sends webhook updates.

## Features

- **User Management**: Authentication and authorization with JWT
- **Role-Based Access Control (RBAC)**: Comprehensive permission system for Admin, Instructor, and Student roles
  - Admins: Full system access
  - Instructors: Manage their own courses, modules, and lessons
  - Students: Self-enroll and access enrolled course content
  - See [RBAC.md](RBAC.md) for complete documentation
- **Course Management**: Create and manage courses with difficulty levels
- **Module & Lesson System**: Organize content in modules with versioned lessons
- **Enrollment System**: Students enroll in courses with status tracking
- **Ownership Controls**: Instructors can only modify their own content
- **Swagger Documentation**: Interactive API documentation and testing
- **Clean Architecture**: Organized code structure with domain-driven design

## Tech Stack

- **Go 1.21+**: Programming language
- **PostgreSQL**: Database
- **GORM**: ORM for database operations
- **Gorilla Mux**: HTTP router
- **JWT**: Token-based authentication
- **Docker**: Containerization

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make (optional)

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/yourusername/go-learning-management-system.git
cd go-learning-management-system
```

2. Copy the environment file:
```bash
cp .env.example .env
```

3. Start the entire project (recommended):
```bash
make start
```
This will:
- Start Docker containers (PostgreSQL + API)
- Wait for services to be ready
- Run database migrations automatically
- Open Swagger UI in your browser

**Alternative:** Start services without opening browser:
```bash
make docker-up
# or
docker-compose up -d
```

4. The API will be available at `http://localhost:8080`
   - Swagger documentation: `http://localhost:8080/swagger/index.html`
   - Health check: `http://localhost:8080/api/health`

5. Stop the services:
```bash
make docker-down
# or
docker-compose down
```

### Local Development

1. Ensure PostgreSQL is running locally

2. Copy and configure environment variables:
```bash
cp .env.example .env
# Edit .env with your local database credentials
```

3. Install dependencies:
```bash
go mod download
```

4. Run the application:
```bash
make run
# or
go run ./src/cmd/api
```

## Available Make Commands

```bash
make start          # 🚀 Start the entire project (Docker + open Swagger)
make help           # Show all available commands
make build          # Build the application
make run            # Run the application locally
make test           # Run tests with coverage
make swagger        # Generate Swagger documentation
make fmt            # Format code
make lint           # Run linter
make docker-build   # Build Docker image
make docker-up      # Start services with Docker Compose
make docker-down    # Stop Docker services
make docker-logs    # View API logs
make dev            # Run with hot reload (requires air)
make db-shell       # Open PostgreSQL shell
```

## API Documentation

Interactive Swagger documentation is available at `/swagger/index.html` when the server is running.

Access it at: `http://localhost:8080/swagger/index.html`

The documentation includes:
- All available endpoints with request/response examples
- Authentication requirements (JWT Bearer token)
- Role-based access control information
- Request parameter descriptions
- Interactive API testing interface

### Authentication & Authorization

All protected endpoints require a JWT token in the `Authorization` header:
```
Authorization: Bearer <your-jwt-token>
```

See [RBAC.md](RBAC.md) for detailed information about:
- User roles and permissions
- Endpoint access requirements
- Testing different permission scenarios
- Security best practices

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user

### Courses
- `GET /api/courses` - List courses (with filters)
- `GET /api/courses/:id` - Get course details
- `POST /api/courses` - Create course (instructor only)
- `PUT /api/courses/:id` - Update course (instructor only)
- `DELETE /api/courses/:id` - Delete course (instructor only)

### Modules
- `GET /api/courses/:courseId/modules` - List course modules
- `GET /api/modules/:id` - Get module details
- `POST /api/courses/:courseId/modules` - Create module (instructor only)
- `PUT /api/modules/:id` - Update module (instructor only)
- `DELETE /api/modules/:id` - Delete module (instructor only)

### Lessons
- `GET /api/modules/:moduleId/lessons` - Get latest lessons for module
- `GET /api/lessons/:id/all-versions` - Get all versions of a lesson (admin only)
- `POST /api/modules/:moduleId/lessons` - Create lesson (instructor only)
- `POST /api/lessons/:lessonId/version` - Create new lesson version (instructor only)
- `DELETE /api/lessons/:id` - Delete lesson (instructor only)

### Enrollments
- `POST /api/courses/:id/enroll` - Enroll in course
- `GET /api/students/:id/courses` - Get student's courses
- `GET /api/courses/:id/students` - Get course's students
- `PUT /api/enrollments/:courseId/status` - Update enrollment status
- `DELETE /api/enrollments/:courseId` - Drop course

## Project Structure

```
.
├── src/
│   ├── cmd/
│   │   └── api/          # Application entry point
│   ├── internal/
│   │   ├── adapter/
│   │   │   └── http/     # HTTP handlers and DTOs
│   │   ├── domain/
│   │   │   ├── entity/   # Domain entities
│   │   │   └── repository/ # Repository interfaces
│   │   └── infrastructure/
│   │       ├── database/ # Database setup
│   │       └── repository/ # Repository implementations
│   ├── middleware/        # HTTP middleware
│   └── usecase/          # Business logic
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

## Environment Variables

See [.env.example](.env.example) for all available configuration options:

- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `SERVER_HOST`: API server host (default: 0.0.0.0)
- `SERVER_PORT`: API server port (default: 8080)
- `JWT_SECRET`: Secret key for JWT tokens
- `JWT_EXPIRATION_HOURS`: Token expiration time (default: 24)

## Database Migrations

The project uses **golang-migrate** for versioned database migrations. Migrations are stored in the `migrations/` directory and run automatically on application startup.

### Migration Commands

```bash
make migrate-up       # Apply all pending migrations
make migrate-down     # Rollback last migration
make migrate-status   # Check migration status
make migrate-version  # Show current version
make migrate-create   # Create new migration files
```

For detailed migration documentation, see [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md).

### Schema Version

Current schema version: **1** (initial_schema)

**Tables**: users, courses, modules, lessons, lesson_versions, course_enrollments, audit_logs

## Development

### Running Tests

Unit tests are implemented for all layers of the application. By default, repository adapter (Postgres) tests are excluded to keep CI fast and deterministic; run the integration sweep when you explicitly need it.

```bash
# Run fast suite (excludes repository adapters)
make test
# equivalent
go test -v -cover $(go list ./... | grep -v src/internal/infrastructure/repository)

# Run everything (including repository adapters)
make test-all
go test -v -cover ./...

# Run entity tests only
go test ./src/internal/domain/entity -v

# Run with coverage report (fast suite)
go test -coverprofile=coverage.out $(go list ./... | grep -v src/internal/infrastructure/repository)
go tool cover -html=coverage.out
```

**Testing Documentation**:
- [TESTING.md](TESTING.md) - Complete testing guide

**Current Test Status**: ✅ **243 tests passing** with 76% coverage

### Code Formatting

```bash
make fmt
# or
go fmt ./...
```

### Linting

```bash
make lint
# or
golangci-lint run
```

## License

MIT License
