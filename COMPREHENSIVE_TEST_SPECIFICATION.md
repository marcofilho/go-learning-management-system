# Comprehensive Endpoint Test Specification

This document specifies all test cases for every endpoint, including success and all error scenarios with expected status codes.

## Test Execution Order

Tests must run sequentially following the system's natural flow:
1. Authentication (Register, Login)
2. User Management
3. Course Management
4. Module Management
5. Lesson Management
6. Enrollment
7. Webhook
8. Audit Logs
9. Authorization Edge Cases

---

## 1. Authentication Endpoints

### POST /api/auth/register

#### Success Cases (201 Created)
- ✅ Register student with valid data
- ✅ Register instructor with admin token
- ✅ Register admin (first admin can register, or with existing admin token)

#### Error Cases
- **400 Bad Request**: 
  - Invalid email format
  - Password too short (< 6 characters)
  - Missing required fields (email, password, first_name, last_name)
  - Invalid role value
  - Malformed JSON
- **403 Forbidden**:
  - Non-admin trying to register instructor/admin role
- **409 Conflict**:
  - Duplicate email registration

### POST /api/auth/login

#### Success Cases (200 OK)
- ✅ Login with valid credentials
- ✅ Login returns JWT token

#### Error Cases
- **400 Bad Request**:
  - Missing email
  - Missing password
  - Malformed JSON
- **401 Unauthorized**:
  - Invalid email
  - Invalid password
  - User not found
  - Inactive user (is_active = false)

---

## 2. User Management Endpoints

### GET /api/users

#### Success Cases (200 OK)
- ✅ Admin can list users
- ✅ Pagination works (limit, offset)
- ✅ Returns paginated response with metadata

#### Error Cases
- **401 Unauthorized**:
  - No authentication token
  - Invalid token
  - Expired token
- **403 Forbidden**:
  - Student trying to list users
  - Instructor trying to list users

### GET /api/users/{id}

#### Success Cases (200 OK)
- ✅ Admin can view any user
- ✅ User can view own profile
- ✅ Returns user data

#### Error Cases
- **400 Bad Request**:
  - Invalid UUID format
- **401 Unauthorized**:
  - No authentication token
- **404 Not Found**:
  - User doesn't exist
  - User is soft-deleted

### PUT /api/users/{id}

#### Success Cases (200 OK)
- ✅ Admin can update any user
- ✅ User can update own profile
- ✅ Returns updated user data

#### Error Cases
- **400 Bad Request**:
  - Invalid UUID format
  - Invalid email format
  - Invalid role value
  - Missing required fields
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Student trying to update other user
  - Instructor trying to update other user
- **404 Not Found**:
  - User doesn't exist

### DELETE /api/users/{id}

#### Success Cases (204 No Content / 200 OK)
- ✅ Admin can delete user
- ✅ Soft delete (user marked as deleted)

#### Error Cases
- **400 Bad Request**:
  - Invalid UUID format
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Non-admin trying to delete
- **404 Not Found**:
  - User doesn't exist

---

## 3. Course Management Endpoints

### POST /api/courses

#### Success Cases (201 Created)
- ✅ Instructor creates own course
- ✅ Admin creates course for any instructor
- ✅ Returns created course data

#### Error Cases
- **400 Bad Request**:
  - Missing title
  - Invalid difficulty_level
  - Invalid instructor_id format
  - Malformed JSON
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Student trying to create course
  - Instructor trying to create course for another instructor (without admin)
- **404 Not Found**:
  - Instructor doesn't exist

### GET /api/courses

#### Success Cases (200 OK)
- ✅ All authenticated users can list courses
- ✅ Pagination works
- ✅ Filtering works (instructor_id, difficulty_level, active_only)
- ✅ Returns paginated response

#### Error Cases
- **401 Unauthorized**:
  - No authentication token

### GET /api/courses/{id}

#### Success Cases (200 OK)
- ✅ Any authenticated user can view course
- ✅ Returns course data

#### Error Cases
- **400 Bad Request**:
  - Invalid UUID format
- **401 Unauthorized**:
  - No authentication token
- **404 Not Found**:
  - Course doesn't exist
  - Course is soft-deleted

### PUT /api/courses/{id}

#### Success Cases (200 OK)
- ✅ Admin can update any course
- ✅ Instructor can update own course
- ✅ Returns updated course data

#### Error Cases
- **400 Bad Request**:
  - Invalid UUID format
  - Invalid difficulty_level
  - Missing title
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Instructor trying to update other instructor's course
  - Student trying to update course
- **404 Not Found**:
  - Course doesn't exist

### DELETE /api/courses/{id}

#### Success Cases (204 No Content / 200 OK)
- ✅ Admin can delete any course
- ✅ Instructor can delete own course
- ✅ Soft delete

#### Error Cases
- **400 Bad Request**:
  - Invalid UUID format
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Instructor trying to delete other instructor's course
  - Student trying to delete course
- **404 Not Found**:
  - Course doesn't exist

---

## 4. Module Management Endpoints

### POST /api/courses/{courseId}/modules

#### Success Cases (201 Created)
- ✅ Admin creates module
- ✅ Course instructor creates module
- ✅ Returns created module data

#### Error Cases
- **400 Bad Request**:
  - Invalid courseId UUID format
  - Missing title
  - Invalid order_index
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Student trying to create module
  - Instructor trying to create module in other instructor's course
- **404 Not Found**:
  - Course doesn't exist

### GET /api/courses/{courseId}/modules

#### Success Cases (200 OK)
- ✅ Admin can list modules
- ✅ Course instructor can list modules
- ✅ Enrolled student can list modules
- ✅ Returns paginated response

#### Error Cases
- **400 Bad Request**:
  - Invalid courseId UUID format
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Non-enrolled student trying to view modules
  - Instructor trying to view other instructor's course modules
- **404 Not Found**:
  - Course doesn't exist

### GET /api/modules/{id}

#### Success Cases (200 OK)
- ✅ Admin can view any module
- ✅ Course instructor can view own module
- ✅ Enrolled student can view module
- ✅ Returns module data

#### Error Cases
- **400 Bad Request**:
  - Invalid UUID format
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Non-enrolled student
  - Instructor trying to view other instructor's module
- **404 Not Found**:
  - Module doesn't exist

### PUT /api/modules/{id}

#### Success Cases (200 OK)
- ✅ Admin can update any module
- ✅ Course instructor can update own module
- ✅ Returns updated module data

#### Error Cases
- **400 Bad Request**:
  - Invalid UUID format
  - Missing title
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Instructor trying to update other instructor's module
  - Student trying to update module
- **404 Not Found**:
  - Module doesn't exist

### DELETE /api/modules/{id}

#### Success Cases (204 No Content / 200 OK)
- ✅ Admin can delete any module
- ✅ Course instructor can delete own module

#### Error Cases
- **400 Bad Request**:
  - Invalid UUID format
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Instructor trying to delete other instructor's module
  - Student trying to delete module
- **404 Not Found**:
  - Module doesn't exist

---

## 5. Lesson Management Endpoints

### POST /api/modules/{moduleId}/lessons

#### Success Cases (201 Created)
- ✅ Admin creates lesson
- ✅ Course instructor creates lesson
- ✅ Creates lesson thread + first version (version 1)
- ✅ Returns created lesson data

#### Error Cases
- **400 Bad Request**:
  - Invalid moduleId UUID format
  - Missing content and video_url
  - Invalid JSON
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Student trying to create lesson
  - Instructor trying to create lesson in other instructor's module
- **404 Not Found**:
  - Module doesn't exist

### GET /api/modules/{moduleId}/lessons

#### Success Cases (200 OK)
- ✅ Admin can list lessons
- ✅ Course instructor can list lessons
- ✅ Enrolled student can list lessons
- ✅ Returns latest versions only
- ✅ Pagination works

#### Error Cases
- **400 Bad Request**:
  - Invalid moduleId UUID format
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Non-enrolled student
  - Instructor trying to view other instructor's lessons
- **404 Not Found**:
  - Module doesn't exist

### POST /api/lessons/{lessonId}/version

#### Success Cases (201 Created)
- ✅ Admin creates new version
- ✅ Course instructor creates new version
- ✅ Auto-increments version number
- ✅ Returns created version data

#### Error Cases
- **400 Bad Request**:
  - Invalid lessonId UUID format
  - Missing content and video_url
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Student trying to create version
  - Instructor trying to create version for other instructor's lesson
- **404 Not Found**:
  - Lesson doesn't exist

### GET /api/lessons/{lessonId}/all-versions

#### Success Cases (200 OK)
- ✅ Admin can view all versions
- ✅ Returns all versions with pagination

#### Error Cases
- **400 Bad Request**:
  - Invalid lessonId UUID format
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Instructor trying to view all versions
  - Student trying to view all versions
- **404 Not Found**:
  - Lesson doesn't exist

### DELETE /api/lessons/{lessonId}

#### Success Cases (204 No Content / 200 OK)
- ✅ Admin can delete lesson
- ✅ Course instructor can delete own lesson
- ✅ Soft delete

#### Error Cases
- **400 Bad Request**:
  - Invalid lessonId UUID format
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Instructor trying to delete other instructor's lesson
  - Student trying to delete lesson
- **404 Not Found**:
  - Lesson doesn't exist

---

## 6. Enrollment Endpoints

### POST /api/courses/{id}/enroll

#### Success Cases (201 Created)
- ✅ Student self-enrolls
- ✅ Admin enrolls any student
- ✅ Returns enrollment data

#### Error Cases
- **400 Bad Request**:
  - Invalid course ID format
  - Invalid student ID format
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Instructor trying to enroll (unless admin)
  - Student trying to enroll another student (unless admin)
- **404 Not Found**:
  - Course doesn't exist
  - Student doesn't exist
- **409 Conflict**:
  - Already enrolled

### GET /api/students/{id}/courses

#### Success Cases (200 OK)
- ✅ Admin can view any student's courses
- ✅ Student can view own courses
- ✅ Returns paginated response

#### Error Cases
- **400 Bad Request**:
  - Invalid student ID format
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Student trying to view other student's courses
  - Instructor trying to view student courses
- **404 Not Found**:
  - Student doesn't exist

### GET /api/courses/{id}/students

#### Success Cases (200 OK)
- ✅ Admin can view course students
- ✅ Course instructor can view own course students
- ✅ Returns paginated response

#### Error Cases
- **400 Bad Request**:
  - Invalid course ID format
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Instructor trying to view other instructor's course students
  - Student trying to view course students
- **404 Not Found**:
  - Course doesn't exist

### PUT /api/enrollments/{courseId}/status

#### Success Cases (200 OK)
- ✅ Student updates own enrollment status
- ✅ Admin updates any enrollment status
- ✅ Returns updated enrollment data

#### Error Cases
- **400 Bad Request**:
  - Invalid course ID format
  - Invalid status value
  - Missing status
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Student trying to update other student's enrollment
  - Instructor trying to update enrollment
- **404 Not Found**:
  - Enrollment doesn't exist

### DELETE /api/enrollments/{courseId}

#### Success Cases (200 OK / 204 No Content)
- ✅ Student drops own enrollment
- ✅ Admin drops any enrollment
- ✅ Returns updated enrollment data

#### Error Cases
- **400 Bad Request**:
  - Invalid course ID format
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Student trying to drop other student's enrollment
  - Instructor trying to drop enrollment
- **404 Not Found**:
  - Enrollment doesn't exist

---

## 7. Webhook Endpoint

### POST /api/certification-webhook

#### Success Cases (200 OK)
- ✅ Valid HMAC signature
- ✅ Student and course exist
- ✅ Enrollment is active
- ✅ Updates enrollment status to completed
- ✅ Creates audit log

#### Error Cases
- **400 Bad Request**:
  - Missing student_id
  - Missing course_id
  - Missing certification_status
  - Invalid UUID formats
  - Invalid certification_status value
- **401 Unauthorized**:
  - Invalid HMAC signature
  - Missing signature header
- **404 Not Found**:
  - Student doesn't exist
  - Course doesn't exist
  - Enrollment doesn't exist
- **409 Conflict**:
  - Enrollment is not active

---

## 8. Audit Logs Endpoint

### GET /api/audit-logs

#### Success Cases (200 OK)
- ✅ Admin can view audit logs
- ✅ Pagination works
- ✅ Filtering works (action, resource_type, user_id)
- ✅ Returns paginated response

#### Error Cases
- **401 Unauthorized**:
  - No authentication token
- **403 Forbidden**:
  - Instructor trying to view audit logs
  - Student trying to view audit logs

---

## 9. Authorization Edge Cases

### Cross-Role Access Tests
- ✅ Student cannot access instructor endpoints
- ✅ Instructor cannot access admin endpoints
- ✅ Instructor cannot modify other instructors' resources
- ✅ Student cannot access other students' resources
- ✅ Admin can access all endpoints

### Token Validation
- ✅ Expired tokens rejected (401)
- ✅ Invalid token format rejected (401)
- ✅ Missing Bearer prefix rejected (401)
- ✅ Tampered tokens rejected (401)

### Ownership Validation
- ✅ Instructor can only modify own courses
- ✅ Instructor can only modify own modules
- ✅ Instructor can only modify own lessons
- ✅ Student can only modify own enrollment

---

## Status Code Summary

### Expected Status Codes by Scenario

| Scenario | Status Code | Description |
|----------|-------------|-------------|
| Success (Create) | 201 | Resource created |
| Success (Get/Update) | 200 | Operation successful |
| Success (Delete) | 204/200 | Resource deleted |
| Bad Request | 400 | Invalid input/validation error |
| Unauthorized | 401 | Missing/invalid authentication |
| Forbidden | 403 | Insufficient permissions |
| Not Found | 404 | Resource doesn't exist |
| Conflict | 409 | Resource conflict (duplicate) |
| Internal Error | 500 | Server error |

---

## Test Execution Strategy

1. **Setup Phase**: Create test users (admin, instructor, student)
2. **Success Path**: Test all success scenarios in natural flow
3. **Error Path**: Test all error scenarios for each endpoint
4. **Authorization**: Test all role combinations
5. **Edge Cases**: Test boundary conditions and edge cases
6. **Cleanup**: Clean up test data

---

## Expected Test Count

Approximate test cases:
- Authentication: ~15 test cases
- User Management: ~25 test cases
- Course Management: ~30 test cases
- Module Management: ~25 test cases
- Lesson Management: ~25 test cases
- Enrollment: ~25 test cases
- Webhook: ~10 test cases
- Audit Logs: ~5 test cases
- Authorization Edge Cases: ~15 test cases

**Total: ~175 test cases**

---

**This specification serves as the complete reference for endpoint testing.**

