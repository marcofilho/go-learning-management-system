# Quick RBAC Testing Guide

This is a quick reference for testing the newly implemented Role-Based Access Control system.

## Setup Test Users

### 1. Start the Application
```bash
make start
```

### 2. Register Test Users

**Register a Student (default role):**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@test.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Student"
  }'
```

**Register an Instructor:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "instructor@test.com",
    "password": "password123",
    "first_name": "Jane",
    "last_name": "Instructor"
  }'
```

**Register an Admin:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "password123",
    "first_name": "Admin",
    "last_name": "User"
  }'
```

### 3. Update Roles in Database

Connect to the database and update roles:
```bash
make db-shell

# In PostgreSQL shell:
UPDATE users SET role = 'instructor' WHERE email = 'instructor@test.com';
UPDATE users SET role = 'admin' WHERE email = 'admin@test.com';
\q
```

### 4. Login and Get Tokens

**Login as Student:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@test.com",
    "password": "password123"
  }'
```
Save the returned token as `$STUDENT_TOKEN`

**Login as Instructor:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "instructor@test.com",
    "password": "password123"
  }'
```
Save the returned token as `$INSTRUCTOR_TOKEN`

**Login as Admin:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "password123"
  }'
```
Save the returned token as `$ADMIN_TOKEN`

## Test Scenarios

### Scenario 1: Student Cannot Create Course ❌

```bash
curl -X POST http://localhost:8080/api/courses \
  -H "Authorization: Bearer $STUDENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Course",
    "description": "Test course"
  }'

# Expected: 403 Forbidden
```

### Scenario 2: Instructor Can Create Course ✅

```bash
curl -X POST http://localhost:8080/api/courses \
  -H "Authorization: Bearer $INSTRUCTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Instructor Course",
    "description": "A course by instructor",
    "difficulty_level": "beginner"
  }'

# Expected: 201 Created
# Save the returned course ID as $COURSE_ID
```

### Scenario 3: Instructor Cannot Modify Another's Course ❌

Create a second course with admin, then try to modify it with instructor:

```bash
# Admin creates a course
curl -X POST http://localhost:8080/api/courses \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Admin Course",
    "description": "A course by admin"
  }'

# Save the returned course ID as $ADMIN_COURSE_ID

# Instructor tries to update it
curl -X PUT http://localhost:8080/api/courses/$ADMIN_COURSE_ID \
  -H "Authorization: Bearer $INSTRUCTOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Modified Title"
  }'

# Expected: 403 Forbidden
```

### Scenario 4: Student Can Self-Enroll ✅

```bash
curl -X POST http://localhost:8080/api/courses/$COURSE_ID/enroll \
  -H "Authorization: Bearer $STUDENT_TOKEN"

# Expected: 201 Created
```

### Scenario 5: Student Cannot View Unenrolled Course Content ❌

```bash
# Before enrollment
curl -X GET http://localhost:8080/api/courses/$ADMIN_COURSE_ID/modules \
  -H "Authorization: Bearer $STUDENT_TOKEN"

# Expected: 403 Forbidden
```

### Scenario 6: Student Can View Enrolled Course Content ✅

```bash
# After enrollment in $COURSE_ID
curl -X GET http://localhost:8080/api/courses/$COURSE_ID/modules \
  -H "Authorization: Bearer $STUDENT_TOKEN"

# Expected: 200 OK (empty array if no modules yet)
```

### Scenario 7: Admin Can Do Everything ✅

```bash
# Admin can view all users
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Admin can delete any course
curl -X DELETE http://localhost:8080/api/courses/$COURSE_ID \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Admin can view audit trail
curl -X GET http://localhost:8080/api/lessons/$LESSON_ID/all-versions \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Expected: 200 OK or 204 No Content
```

### Scenario 8: Instructor Cannot View Audit Trail ❌

```bash
curl -X GET http://localhost:8080/api/lessons/$LESSON_ID/all-versions \
  -H "Authorization: Bearer $INSTRUCTOR_TOKEN"

# Expected: 403 Forbidden
```

### Scenario 9: User Can Only Update Own Profile ✅/❌

```bash
# Get user IDs from login response

# User updates own profile (works)
curl -X PUT http://localhost:8080/api/users/$OWN_USER_ID \
  -H "Authorization: Bearer $STUDENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Updated Name"
  }'

# Expected: 200 OK

# User tries to update another user (fails)
curl -X PUT http://localhost:8080/api/users/$OTHER_USER_ID \
  -H "Authorization: Bearer $STUDENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Hacker"
  }'

# Expected: 403 Forbidden
```

### Scenario 10: Instructor Can View Own Course Students ✅

```bash
curl -X GET http://localhost:8080/api/courses/$COURSE_ID/students \
  -H "Authorization: Bearer $INSTRUCTOR_TOKEN"

# Expected: 200 OK with list of enrolled students
```

## Using Swagger UI

The easiest way to test is through Swagger UI at `http://localhost:8080/swagger/index.html`:

1. **Authorize** - Click the green "Authorize" button at the top
2. **Enter Token** - Paste your JWT token (without "Bearer " prefix)
3. **Click Authorize** - Now all requests will include your token
4. **Test Endpoints** - Try different endpoints and see the responses

You can switch between different users by authorizing with different tokens.

## Expected Behaviors Summary

| Action | Student | Instructor | Admin |
|--------|---------|------------|-------|
| Create course | ❌ 403 | ✅ 201 | ✅ 201 |
| Update own course | ❌ 403 | ✅ 200 | ✅ 200 |
| Update others' course | ❌ 403 | ❌ 403 | ✅ 200 |
| View course list | ✅ 200 | ✅ 200 | ✅ 200 |
| Enroll in course | ✅ 201 | ❌ 403 | ✅ 201 |
| View enrolled content | ✅ 200 | N/A | ✅ 200 |
| View unenrolled content | ❌ 403 | ❌ 403 | ✅ 200 |
| View audit trail | ❌ 403 | ❌ 403 | ✅ 200 |
| List all users | ❌ 403 | ❌ 403 | ✅ 200 |
| Update own profile | ✅ 200 | ✅ 200 | ✅ 200 |
| Update others' profile | ❌ 403 | ❌ 403 | ✅ 200 |

## Troubleshooting

**401 Unauthorized:**
- Check if token is included in Authorization header
- Verify token format: `Bearer <token>`
- Token might be expired (default 24h)

**403 Forbidden:**
- This is expected! RBAC is working correctly
- Verify you're using the correct role's token
- Check RBAC.md for required permissions

**404 Not Found:**
- Resource doesn't exist
- Or you don't have permission to know it exists (security feature)

## Next Steps

For more detailed information, see:
- [RBAC.md](RBAC.md) - Complete RBAC documentation
- [README.md](README.md) - General project documentation
- [SWAGGER_GUIDE.md](SWAGGER_GUIDE.md) - Swagger usage guide
