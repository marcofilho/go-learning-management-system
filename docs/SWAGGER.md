# Swagger API Documentation

## Overview

The LMS API includes interactive Swagger documentation for testing and exploring all endpoints.

## Accessing Swagger UI

1. **Start the application**:
   ```bash
   make start
   # or
   docker-compose up -d
   ```

2. **Open your browser** and navigate to:
   ```
   http://localhost:8080/swagger/index.html
   ```

## Using Swagger UI

### Authentication

Most endpoints require authentication:

1. **Register a new user** (POST /api/auth/register):
   ```json
   {
     "email": "test@example.com",
     "password": "password123",
     "first_name": "Test",
     "last_name": "User",
     "role": "student"
   }
   ```

2. **Login to get JWT token** (POST /api/auth/login):
   ```json
   {
     "email": "test@example.com",
     "password": "password123"
   }
   ```

3. **Authorize**:
   - Click the "Authorize" button at the top
   - Enter: `Bearer <your-jwt-token>`
   - Click "Authorize"
   - Now you can test protected endpoints

### Testing Endpoints

1. Expand an endpoint
2. Click "Try it out"
3. Fill in parameters/body
4. Click "Execute"
5. View response

## Generating Documentation

### Regenerate Swagger Docs

```bash
make swagger
# or
swag init -g src/cmd/api/main.go -o docs
```

### Swagger Annotations

Annotations are placed directly above handler function signatures:

```go
// CreateCourse godoc
// @Summary Create a new course
// @Description Create a course (Admin or Instructor only)
// @Tags Courses
// @Accept json
// @Produce json
// @Param course body dto.CreateCourseRequest true "Course data"
// @Success 201 {object} dto.SuccessResponse{data=dto.CourseDTO}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /courses [post]
func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
    // Handler implementation
}
```

## API Information

- **Title**: Learning Management System API
- **Version**: 1.0
- **Base URL**: http://localhost:8080/api
- **Authentication**: JWT Bearer Token

## Documentation Features

- ✅ All endpoints documented
- ✅ Request/response schemas
- ✅ Authentication requirements
- ✅ Role-based access control info
- ✅ Example payloads
- ✅ Interactive testing

---

**Access**: http://localhost:8080/swagger/index.html  
**Generation**: `make swagger`

