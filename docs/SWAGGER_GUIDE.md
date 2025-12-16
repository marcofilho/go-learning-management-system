# Swagger API Documentation Guide

## Overview

The Learning Management System API now includes interactive Swagger documentation that allows you to:
- View all available endpoints
- Test endpoints directly from your browser
- See request/response schemas
- Understand authentication requirements

## Accessing Swagger UI

1. **Start the application**:
   ```bash
   # Using Docker Compose (recommended)
   make docker-up
   
   # OR run locally
   make run
   ```

2. **Open your browser** and navigate to:
   ```
   http://localhost:8080/swagger/index.html
   ```

## Using Swagger UI

### 1. Authentication

Most endpoints require authentication. To test protected endpoints:

1. **Register a new user** (if you don't have an account):
   - Expand the `POST /api/auth/register` endpoint
   - Click "Try it out"
   - Fill in the request body:
     ```json
     {
       "email": "test@example.com",
       "password": "password123",
       "first_name": "Test",
       "last_name": "User",
       "role": "student"
     }
     ```
   - Click "Execute"
   - Note the response

2. **Login to get JWT token**:
   - Expand the `POST /api/auth/login` endpoint
   - Click "Try it out"
   - Fill in credentials:
     ```json
     {
       "email": "test@example.com",
       "password": "password123"
     }
     ```
   - Click "Execute"
   - **Copy the token** from the response

3. **Authorize Swagger**:
   - Click the **"Authorize"** button (🔒) at the top right
   - In the "Value" field, enter: `Bearer YOUR_TOKEN_HERE`
     - Example: `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
   - Click "Authorize"
   - Click "Close"

Now all your requests will include the authentication header!

### 2. Testing Endpoints

#### Example: Creating a Course

1. Navigate to the **Courses** section
2. Expand `POST /api/courses`
3. Click "Try it out"
4. Modify the request body:
   ```json
   {
     "title": "Introduction to Programming",
     "description": "Learn the basics of programming",
     "difficulty_level": "beginner"
   }
   ```
5. Click "Execute"
6. View the response below

#### Example: Listing Courses with Filters

1. Expand `GET /api/courses`
2. Click "Try it out"
3. Add query parameters:
   - `difficulty_level`: beginner
   - `active_only`: true
   - `limit`: 10
4. Click "Execute"

### 3. API Tags

Endpoints are organized by functionality:

- **Authentication**: Register and login
- **Users**: User management
- **Courses**: Course CRUD operations
- **Modules**: Module management within courses
- **Lessons**: Lesson versioning and management
- **Enrollments**: Student enrollment in courses

## Regenerating Documentation

If you modify the API code and add/change endpoints:

```bash
# Regenerate Swagger docs
make swagger

# Rebuild and restart the application
make docker-up
# or
make run
```

## Swagger Annotations

The documentation is generated from code annotations. Example:

```go
// CreateCourse godoc
// @Summary Create a new course
// @Description Create a new course (instructor only)
// @Tags Courses
// @Accept json
// @Produce json
// @Param request body dto.CreateCourseRequest true "Course details"
// @Success 201 {object} dto.SuccessResponse{data=dto.CourseDTO}
// @Failure 400 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /courses [post]
func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

## Common Issues

### Issue: "Authorize" button not working
**Solution**: Make sure you include the `Bearer ` prefix before your token.

### Issue: Getting 401 Unauthorized
**Solution**: 
1. Check if your token has expired (default: 24 hours)
2. Login again to get a fresh token
3. Make sure you clicked "Authorize" after entering the token

### Issue: Swagger UI not loading
**Solution**:
1. Check if the server is running: `curl http://localhost:8080/api/health`
2. Verify the port is correct (default: 8080)
3. Check logs for errors: `make docker-logs`

## API Testing Workflow

Typical workflow for testing the complete API:

1. **Register a user** (role: student)
2. **Login** and get token
3. **Authorize** in Swagger
4. **Register an instructor** (role: instructor)
5. **Login as instructor**
6. **Create a course** (as instructor)
7. **Add modules** to the course
8. **Add lessons** to modules
9. **Switch back to student** token
10. **Enroll in the course**
11. **View enrolled courses**
12. **Access course modules and lessons**

## Additional Resources

- [Swagger Official Documentation](https://swagger.io/docs/)
- [Go Swag Library](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)

## Security Notes

⚠️ **Important**: 
- Never commit your actual JWT secret to version control
- Use strong passwords in production
- Keep your authentication tokens secure
- Tokens expire after 24 hours by default

## Support

For issues or questions:
- Check the application logs
- Review the README.md for setup instructions
- Ensure all environment variables are configured correctly
