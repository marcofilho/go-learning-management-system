# Role-Based Access Control (RBAC) Implementation

## Overview

The API now implements comprehensive role-based access control with three user roles:
- **Admin**: Full access to all resources
- **Instructor**: Can create and manage their own courses, modules, and lessons
- **Student**: Can enroll in courses and view content they're enrolled in

## User Roles

### Admin
- Can manage all users (list, view, update, delete)
- Can create, view, update, and delete any course
- Can view all enrolled students in any course
- Can manage modules and lessons for any course
- Can view audit trail (all lesson versions)

### Instructor
- Can create courses
- Can manage (update/delete) only their own courses
- Can upload and modify content for their own courses
- Can create/update/delete modules and lessons for their own courses
- Can view enrolled students in their own courses
- Cannot enroll users or manage other instructors' courses

### Student
- Can self-enroll in courses
- Can view course content only for enrolled courses
- Can view their own enrollment history
- Cannot create or modify courses, modules, or lessons
- Cannot view audit trail

## Permission Matrix

| Action | Admin | Instructor | Student |
|--------|-------|------------|---------|
| Create/update/delete courses | ✓ | ✓ (own only) | ✗ |
| Upload/modify course content | ✓ | ✓ (own only) | ✗ |
| Enroll users | ✓ | ✗ | Self-enroll only |
| View course content | ✓ | ✓ (own only) | ✓ (enrolled only) |
| Manage lessons/versions | ✓ | ✓ (own only) | ✗ |
| View audit content | ✓ | ✗ | ✗ |
| Manage users | ✓ | ✗ | Self only |

## Middleware Components

### 1. AuthMiddleware
**Purpose**: Validates JWT token and extracts user claims  
**Location**: `src/internal/adapter/http/middleware/auth.go`

```go
middleware.AuthMiddleware(tokenProvider)
```

### 2. RequireRole
**Purpose**: Restricts access to specific roles  
**Location**: `src/internal/adapter/http/middleware/auth.go`

```go
middleware.RequireRole(entity.UserRoleAdmin, entity.UserRoleInstructor)
```

### 3. RequireCourseOwnership
**Purpose**: Ensures user is admin or course instructor  
**Location**: `src/internal/adapter/http/middleware/authorization.go`

```go
middleware.RequireCourseOwnership(courseRepo)
```

### 4. RequireEnrollment
**Purpose**: Ensures student is enrolled or user owns the course  
**Location**: `src/internal/adapter/http/middleware/authorization.go`

```go
middleware.RequireEnrollment(courseRepo, enrollmentRepo)
```

### 5. RequireModuleOwnership
**Purpose**: Ensures user is admin or owns the course that contains the module  
**Location**: `src/internal/adapter/http/middleware/authorization.go`

```go
middleware.RequireModuleOwnership(moduleRepo, courseRepo)
```

### 6. RequireAdminOrSelf
**Purpose**: Allows access to own resources or admin access to all  
**Location**: `src/internal/adapter/http/middleware/authorization.go`

```go
middleware.RequireAdminOrSelf()
```

## API Endpoints & Required Permissions

### Authentication (Public)
```
POST /api/auth/register - Public (creates student account)
POST /api/auth/login    - Public
```

### Users
```
GET    /api/users           - Admin only
GET    /api/users/{id}      - Authenticated (own profile or admin)
PUT    /api/users/{id}      - Admin or self
DELETE /api/users/{id}      - Admin only
```

### Courses
```
GET    /api/courses         - Authenticated
POST   /api/courses         - Admin or Instructor
GET    /api/courses/{id}    - Authenticated
PUT    /api/courses/{id}    - Admin or Course Owner
DELETE /api/courses/{id}    - Admin or Course Owner
```

### Enrollments
```
POST   /api/courses/{id}/enroll        - Student (self) or Admin
GET    /api/courses/{id}/students      - Admin or Course Owner
GET    /api/students/{id}/courses      - Admin or Self
```

### Modules
```
POST   /api/courses/{courseId}/modules  - Admin or Course Owner
GET    /api/courses/{courseId}/modules  - Enrolled Student, Course Owner, or Admin
PUT    /api/modules/{id}                - Admin or Course Owner
DELETE /api/modules/{id}                - Admin or Course Owner
```

### Lessons
```
POST   /api/modules/{moduleId}/lessons           - Admin or Course Owner
GET    /api/modules/{moduleId}/lessons           - Authenticated (enrollment checked in handler)
POST   /api/lessons/{lessonId}/version           - Admin or Instructor
DELETE /api/lessons/{lessonId}                   - Admin or Instructor
GET    /api/lessons/{lessonId}/all-versions      - Admin only (audit trail)
```

## Helper Functions

### GetUserIDFromContext
Extracts the authenticated user's ID from the request context.

```go
userID, ok := GetUserIDFromContext(r)
if !ok {
    // Handle unauthorized access
}
```

### GetUserFromContext
Extracts the full user claims (ID, email, role) from the context.

```go
claims, ok := GetUserFromContext(r.Context())
if ok {
    // Access claims.UserID, claims.Email, claims.Role
}
```

### IsAdmin
Checks if the current user has admin role.

```go
if IsAdmin(r) {
    // Admin-specific logic
}
```

### IsInstructor
Checks if the current user is an instructor or admin.

```go
if IsInstructor(r) {
    // Instructor logic
}
```

## JWT Token Structure

The JWT token includes the following claims:

```json
{
  "user_id": "uuid-string",
  "email": "user@example.com",
  "role": "admin|instructor|student",
  "exp": 1234567890,
  "iat": 1234567890,
  "nbf": 1234567890
}
```

## Testing RBAC

### 1. Register Users with Different Roles

**Create Admin** (requires database access):
```sql
UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';
```

**Create Instructor**:
```sql
UPDATE users SET role = 'instructor' WHERE email = 'instructor@example.com';
```

**Create Student** (default):
```bash
POST /api/auth/register
{
  "email": "student@example.com",
  "password": "password",
  "first_name": "John",
  "last_name": "Doe"
}
```

### 2. Login and Get Token

```bash
POST /api/auth/login
{
  "email": "admin@example.com",
  "password": "password"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 3. Use Token in Requests

```bash
GET /api/courses
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 4. Test Permission Scenarios

#### Scenario 1: Student tries to create course (should fail)
```bash
POST /api/courses
Authorization: Bearer <student-token>

Expected: 403 Forbidden
```

#### Scenario 2: Instructor creates course (should succeed)
```bash
POST /api/courses
Authorization: Bearer <instructor-token>
{
  "title": "My Course",
  "description": "Course description"
}

Expected: 201 Created
```

#### Scenario 3: Instructor tries to update another instructor's course (should fail)
```bash
PUT /api/courses/{other-instructor-course-id}
Authorization: Bearer <instructor-token>

Expected: 403 Forbidden
```

#### Scenario 4: Student tries to view unenrolled course content (should fail)
```bash
GET /api/courses/{course-id}/modules
Authorization: Bearer <student-token>

Expected: 403 Forbidden (if not enrolled)
```

#### Scenario 5: Admin can do anything (should succeed)
```bash
DELETE /api/courses/{any-course-id}
Authorization: Bearer <admin-token>

Expected: 204 No Content
```

## Error Responses

### 401 Unauthorized
Missing or invalid JWT token.

```json
{
  "error": "Unauthorized",
  "message": "Invalid token"
}
```

### 403 Forbidden
Valid token but insufficient permissions.

```json
{
  "error": "Insufficient permissions",
  "message": "Only the course instructor can perform this action"
}
```

### 404 Not Found
Resource doesn't exist (also prevents information disclosure).

```json
{
  "error": "Not found",
  "message": "Course not found"
}
```

## Security Best Practices

1. **Always use HTTPS in production** to protect JWT tokens in transit
2. **Store JWT securely** on the client side (httpOnly cookies recommended)
3. **Set appropriate token expiration** (current default: 24 hours)
4. **Rotate JWT secrets** regularly in production
5. **Never log or expose JWT tokens** in error messages
6. **Validate input** at both middleware and handler levels
7. **Use rate limiting** to prevent brute force attacks
8. **Implement refresh tokens** for long-lived sessions (future enhancement)

## Migration from Old System

The system previously had basic authentication but no authorization. The changes include:

1. **Added role field** to User entity (default: 'student')
2. **Created authorization middleware** for ownership and enrollment checks
3. **Updated all routes** with appropriate middleware chains
4. **Updated handlers** to use new context helper functions
5. **All existing users** will have 'student' role by default

### Database Migration

Run the application to auto-apply migrations:
```bash
make start
```

Or manually:
```bash
go run src/cmd/migrate/main.go
```

## Future Enhancements

1. **Granular Permissions**: Replace role-based with permission-based (e.g., can_create_course, can_delete_lesson)
2. **Course Moderators**: Allow instructors to assign TAs or co-instructors
3. **API Rate Limiting**: Per-role rate limits
4. **Audit Logging**: Track all administrative actions
5. **Two-Factor Authentication**: Enhanced security for admin accounts
6. **OAuth/SSO Integration**: Enterprise authentication support
7. **Session Management**: Active session tracking and revocation
8. **IP Whitelisting**: Restrict admin access by IP range
