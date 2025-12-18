# Go Learning Management System - Clean Architecture

## Project Structure

```
migrations/                     # Database migrations (golang-migrate)
├── 000001_initial_schema.up.sql
└── 000001_initial_schema.down.sql
src/
├── cmd/
│   ├── api/                    # Application entry point
│   │   ├── main.go             # Main application
│   │   └── routes.go           # HTTP route configuration
│   └── migrate/                # Migration CLI tool
│       └── main.go             # Migration commands
├── internal/
│   ├── domain/                 # Domain layer (entities & interfaces)
│   │   ├── entity/             # Domain entities
│   │   │   ├── user.go         # User entity with roles
│   │   │   ├── course.go       # Course entity with status
│   │   │   ├── enrollment.go   # Enrollment tracking
│   │   │   └── errors.go       # Domain errors
│   │   └── repository/         # Repository interfaces (DIP)
│   │       ├── user_repository.go
│   │       ├── course_repository.go
│   │       └── enrollment_repository.go
│   ├── infrastructure/         # Infrastructure layer
│   │   ├── auth/               # Authentication
│   │   │   └── jwt_provider.go # JWT token provider
│   │   ├── database/           # Database connection
│   │   │   └── postgres.go     # PostgreSQL connection pool
│   │   └── repository/         # Repository implementations
│   │       ├── user_repository_postgres.go
│   │       ├── course_repository_postgres.go
│   │       └── enrollment_repository_postgres.go
│   ├── adapter/http/           # HTTP adapter layer
│   │   ├── handler/            # HTTP handlers
│   │   │   ├── response.go     # Response helpers & DTO mappers
│   │   │   ├── user_handler.go # User endpoints
│   │   │   ├── course_handler.go # Course endpoints
│   │   │   └── enrollment_handler.go # Enrollment endpoints
│   │   ├── middleware/         # HTTP middleware
│   │   │   ├── auth.go         # JWT authentication & authorization
│   │   │   ├── logger.go       # Request logging
│   │   │   └── common.go       # CORS & content-type
│   │   └── dto/                # Data Transfer Objects
│   │       ├── auth.go         # Auth DTOs (Register, Login)
│   │       ├── user.go         # User DTOs
│   │       ├── course.go       # Course DTOs
│   │       ├── enrollment.go   # Enrollment DTOs
│   │       └── common.go       # Common DTOs
│   └── config/                 # Configuration
│       └── config.go           # App config loading
└── usecase/                    # Business logic layer
    ├── user_usecase.go         # User business logic
    ├── course_usecase.go       # Course business logic
    └── enrollment_usecase.go   # Enrollment business logic
```

## Architecture Principles

### Clean Architecture Layers
1. **Domain Layer** (`internal/domain/`)
   - Core business entities
   - Repository interfaces (Dependency Inversion)
   - Domain errors
   - No dependencies on other layers

2. **Use Case Layer** (`src/usecase/`)
   - Business logic
   - Orchestrates domain entities
   - Depends only on domain layer

3. **Infrastructure Layer** (`internal/infrastructure/`)
   - Implements repository interfaces
   - Database connections
   - Authentication (JWT)
   - External services

4. **Adapter Layer** (`internal/adapter/http/`)
   - HTTP handlers
   - Middleware
   - DTOs (Data Transfer Objects)
   - Converts between HTTP and domain

### SOLID Principles Implemented

**Single Responsibility Principle (SRP)**
- Each use case handles one domain
- Each handler manages one resource type
- Each repository handles one entity

**Open/Closed Principle (OCP)**
- Repository interfaces allow switching implementations
- Middleware chain is extensible
- Use cases depend on abstractions

**Liskov Substitution Principle (LSP)**
- Repository implementations are interchangeable
- TokenProvider interface allows different auth mechanisms

**Interface Segregation Principle (ISP)**
- Separate repository interfaces per entity
- Focused middleware functions

**Dependency Inversion Principle (DIP)**
- Use cases depend on repository interfaces, not implementations
- Handlers depend on use cases, not repositories
- Infrastructure implements domain interfaces

## Key Features

### Authentication & Authorization
- **JWT-based authentication**
- **Role-based access control**: Student, Instructor, Admin
- **Middleware protection** for routes
- **Token validation** on each protected request

### User Management
- User registration with email validation
- Login with JWT token generation
- Role assignment (student/instructor/admin)
- User CRUD operations
- Admin-only user management

### Course Management
- Instructors can create courses
- Course status: Draft, Published, Archived
- Price and duration management
- Public course listing (published only)
- Instructor-specific course management

### Enrollment System
- Students enroll in published courses
- Progress tracking (0-100%)
- Auto-completion at 100%
- Enrollment status: Active, Completed, Dropped, Suspended
- Last access tracking

## API Endpoints

### Public Endpoints
```
GET  /api/v1/health            # Health check
POST /api/v1/auth/register     # User registration
POST /api/v1/auth/login        # User login
GET  /api/v1/courses/published # List published courses
GET  /api/v1/courses/{id}      # Get course details
```

### Authenticated Endpoints

**Users**
```
GET    /api/v1/users           # List users (admin)
GET    /api/v1/users/{id}      # Get user
PUT    /api/v1/users/{id}      # Update user
DELETE /api/v1/users/{id}      # Delete user (admin)
```

**Courses**
```
GET    /api/v1/courses               # List all courses
POST   /api/v1/courses               # Create course (instructor/admin)
PUT    /api/v1/courses/{id}          # Update course
DELETE /api/v1/courses/{id}          # Delete course
POST   /api/v1/courses/{id}/publish  # Publish course
```

**Enrollments**
```
POST /api/v1/enrollments                  # Enroll in course
GET  /api/v1/enrollments/my               # My enrollments
GET  /api/v1/enrollments/{id}             # Get enrollment
PUT  /api/v1/enrollments/{id}/progress    # Update progress
POST /api/v1/enrollments/{id}/drop        # Drop enrollment
GET  /api/v1/enrollments/course/{course_id} # Course enrollments
```

## Environment Variables

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=lms_db
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION_HOURS=24
```

## Running the Application

1. **Start PostgreSQL**:
   ```bash
   docker-compose up -d postgres
   ```

2. **Run migrations**:
   ```bash
   make migrate-up
   ```

3. **Start the API**:
   ```bash
   ./bin/api
   # or
   make run
   ```

4. **Health check**:
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

## Build Commands

```bash
# Build the application
go build -o bin/api ./src/cmd/api

# Build using Make
make build

# Run the application
make run

# Run database migrations
make migrate-up
make migrate-down

# Development with hot reload
make dev
```

## Dependencies

- **gorilla/mux v1.8.1**: HTTP routing
- **pgx/v5 v5.5.1**: PostgreSQL driver
- **golang-jwt/jwt/v5 v5.2.0**: JWT authentication
- **golang.org/x/crypto**: Password hashing (bcrypt)
- **google/uuid v1.5.0**: UUID generation

## Database Schema

### Users
- id (UUID)
- email (unique)
- password_hash
- first_name, last_name
- role (student/instructor/admin)
- is_active (boolean)
- created_at, updated_at

### Courses
- id (UUID)
- title, description
- instructor_id (FK to users)
- status (draft/published/archived)
- price, duration
- is_published (boolean)
- created_at, updated_at

### Enrollments
- id (UUID)
- user_id (FK to users)
- course_id (FK to courses)
- status (active/completed/dropped/suspended)
- progress (0-100)
- enrolled_at, completed_at, last_access_at
- created_at, updated_at

## Security Features

1. **Password Hashing**: bcrypt with default cost
2. **JWT Tokens**: HS256 signing algorithm
3. **Role-based Access Control**: Admin, Instructor, Student
4. **Authorization Middleware**: Protects sensitive endpoints
5. **Input Validation**: DTOs with validation tags

## Project Status

✅ **Complete and Successfully Built**
- All layers implemented
- Clean architecture maintained
- SOLID principles applied
- Database integration ready
- JWT authentication configured
- Role-based authorization
- Full CRUD operations
- Project builds without errors

## Next Steps (Optional Enhancements)

1. Add request validation using validator package
2. Implement API documentation (Swagger/OpenAPI)
3. Add unit tests and integration tests
4. Implement rate limiting
5. Add Redis for caching
6. Implement file uploads for course materials
7. Add email notifications
8. Create admin dashboard
9. Implement course categories and tags
10. Add search and filtering capabilities
