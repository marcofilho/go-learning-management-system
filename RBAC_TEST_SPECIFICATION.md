# Role-Based Access Control (RBAC) Test Specification

## Overview
This document provides comprehensive test specifications for validating roles and permissions across all endpoints in the Learning Management System.

## Test Roles Configuration

### Test Users
```
- Admin: admin@test.com / password123 (Role: admin)
- Instructor 1: instructor@test.com / password123 (Role: instructor)
- Instructor 2: instructor2@test.com / password123 (Role: instructor)
- Student 1: student@test.com / password123 (Role: student)
- Student 2: student2@test.com / password123 (Role: student)
```

---

## 1. Authentication Endpoints

### POST /api/auth/register
**Purpose**: Create new user account

| Test Case | Input | Expected Result | Reason |
|-----------|-------|-----------------|---------|
| Valid registration | Email, password, names, role | 201 Created | Anyone can register |
| Duplicate email | Existing email | 400 Bad Request | Email must be unique |
| Invalid role | role: "superuser" | 400 Bad Request | Role must be valid enum |
| Weak password | password: "123" | 400 Bad Request | Password requirements |

### POST /api/auth/login
**Purpose**: Authenticate and receive JWT token

| Test Case | Input | Expected Result | Reason |
|-----------|-------|-----------------|---------|
| Valid credentials | Correct email/password | 200 OK + token | Authentication success |
| Invalid credentials | Wrong password | 401 Unauthorized | Authentication failure |
| Non-existent email | fake@test.com | 401 Unauthorized | User not found |

---

## 2. User Management Endpoints

### GET /api/users
**Purpose**: List all users (Admin only)

| Role | Expected Status | Reason |
|------|-----------------|---------|
| Admin | 200 OK | Admin can list all users |
| Instructor | 403 Forbidden | Only admin can list users |
| Student | 403 Forbidden | Only admin can list users |
| Unauthenticated | 401 Unauthorized | Authentication required |

**Validation**: Response should contain user list with pagination

### GET /api/users/:id
**Purpose**: Get user details

| Requester | Target User | Expected Status | Reason |
|-----------|-------------|-----------------|---------|
| Admin | Any user | 200 OK | Admin can view any user |
| User | Self | 200 OK | Users can view own profile |
| User | Other user | 403 Forbidden | Users cannot view other profiles |
| Unauthenticated | Any | 401 Unauthorized | Authentication required |

### PUT /api/users/:id
**Purpose**: Update user information

| Requester | Target User | Expected Status | Reason |
|-----------|-------------|-----------------|---------|
| Admin | Any user | 200 OK | Admin can update any user |
| User | Self | 200 OK | Users can update own profile |
| Student 1 | Student 2 | 403 Forbidden | Cannot update other users |
| Instructor 1 | Student 1 | 403 Forbidden | Cannot update other users |

**Additional Validations**:
- Non-admin cannot change role
- Email must remain unique
- UUID validation on user_id parameter

### DELETE /api/users/:id
**Purpose**: Soft delete user

| Role | Expected Status | Reason |
|------|-----------------|---------|
| Admin | 204 No Content | Only admin can delete users |
| Instructor | 403 Forbidden | Only admin can delete |
| Student | 403 Forbidden | Only admin can delete |

---

## 3. Course Management Endpoints

### POST /api/courses
**Purpose**: Create new course

| Role | With Valid instructor_id | Expected Status | Reason |
|------|-------------------------|-----------------|---------|
| Instructor | Self | 201 Created | Instructors create own courses |
| Instructor | Other instructor | 400 Bad Request | Cannot create for others |
| Admin | Any instructor | 201 Created | Admin can create for any instructor |
| Student | Any | 403 Forbidden | Students cannot create courses |

**Validations**:
- instructor_id must exist
- instructor_id must have role 'instructor'
- Title and description required
- difficulty_level must be: beginner, intermediate, advanced

### GET /api/courses
**Purpose**: List all courses

| Role | Expected Status | Content |
|------|-----------------|---------|
| Admin | 200 OK | All courses |
| Instructor | 200 OK | All courses |
| Student | 200 OK | All courses |
| Unauthenticated | 401 Unauthorized | N/A |

### GET /api/courses/:id
**Purpose**: Get course details

| Requester | Expected Status | Reason |
|-----------|-----------------|---------|
| Any authenticated user | 200 OK | Course details are public to authenticated users |
| Unauthenticated | 401 Unauthorized | Authentication required |

### PUT /api/courses/:id
**Purpose**: Update course

| Requester | Course Owner | Expected Status | Reason |
|-----------|--------------|-----------------|---------|
| Admin | Any | 200 OK | Admin can update any course |
| Instructor | Own course | 200 OK | Owner can update |
| Instructor 1 | Instructor 2's course | 403 Forbidden | Cannot update others' courses |
| Student | Any | 403 Forbidden | Students cannot update courses |

**Validations**:
- Cannot change instructor_id to invalid instructor
- All course fields must remain valid

### DELETE /api/courses/:id
**Purpose**: Soft delete course

| Requester | Course Owner | Expected Status | Reason |
|-----------|--------------|-----------------|---------|
| Admin | Any | 204 No Content | Admin can delete any course |
| Instructor | Own course | 204 No Content | Owner can delete own course |
| Instructor 1 | Instructor 2's course | 403 Forbidden | Cannot delete others' courses |
| Student | Any | 403 Forbidden | Students cannot delete courses |

---

## 4. Enrollment Endpoints

### POST /api/courses/:id/enroll
**Purpose**: Enroll current user in course

| Role | Expected Status | Reason |
|------|-----------------|---------|
| Student | 201 Created | Students can enroll |
| Instructor | 400 Bad Request | Instructors cannot enroll |
| Admin | 400 Bad Request | Admins cannot enroll |

**Validations**:
- Course must exist
- Student cannot enroll twice in same course
- Enrollment status defaults to 'active'

### GET /api/courses/:id/students
**Purpose**: List students enrolled in course

| Requester | Course Owner | Expected Status | Reason |
|-----------|--------------|-----------------|---------|
| Admin | Any | 200 OK | Admin can view all enrollments |
| Instructor | Own course | 200 OK | Owner can view enrolled students |
| Instructor | Other's course | 403 Forbidden | Cannot view other instructors' students |
| Student | Any | 403 Forbidden | Students cannot view enrollment list |

### GET /api/students/:id/courses
**Purpose**: List courses student is enrolled in

| Requester | Target Student | Expected Status | Reason |
|-----------|----------------|-----------------|---------|
| Admin | Any student | 200 OK | Admin can view any student's courses |
| Student | Self | 200 OK | Students can view own courses |
| Student 1 | Student 2 | 403 Forbidden | Cannot view other students' courses |
| Instructor | Any student | 403 Forbidden | Instructors cannot view student courses |

---

## 5. Module Management Endpoints

### POST /api/courses/:courseId/modules
**Purpose**: Create module in course

| Requester | Course Owner | Expected Status | Reason |
|-----------|--------------|-----------------|---------|
| Admin | Any | 201 Created | Admin can create modules |
| Instructor | Own course | 201 Created | Owner can create modules |
| Instructor | Other's course | 403 Forbidden | Cannot modify other courses |
| Student | Any | 403 Forbidden | Students cannot create modules |

**Validations**:
- Course must exist
- order_index must be positive
- Title required

### GET /api/courses/:courseId/modules
**Purpose**: List modules in course

| Requester | Enrollment Status | Expected Status | Reason |
|-----------|-------------------|-----------------|---------|
| Admin | N/A | 200 OK | Admin views all |
| Instructor | Own course | 200 OK | Owner views modules |
| Student | Enrolled | 200 OK | Enrolled students view modules |
| Student | Not enrolled | 403 Forbidden | Must be enrolled to view content |

### GET /api/modules/:id
**Purpose**: Get module details

| Requester | Enrollment Status | Expected Status | Reason |
|-----------|-------------------|-----------------|---------|
| Admin | N/A | 200 OK | Admin views all |
| Instructor | Own course | 200 OK | Owner views module |
| Student | Enrolled in course | 200 OK | Enrolled students view |
| Student | Not enrolled | 403 Forbidden | Must be enrolled |

### PUT /api/modules/:id
**Purpose**: Update module

| Requester | Course Owner | Expected Status | Reason |
|-----------|--------------|-----------------|---------|
| Admin | Any | 200 OK | Admin updates all |
| Instructor | Own course | 200 OK | Owner updates module |
| Instructor | Other's course | 403 Forbidden | Cannot modify others |
| Student | Any | 403 Forbidden | Students cannot update |

### DELETE /api/modules/:id
**Purpose**: Soft delete module

| Requester | Course Owner | Expected Status | Reason |
|-----------|--------------|-----------------|---------|
| Admin | Any | 204 No Content | Admin deletes all |
| Instructor | Own course | 204 No Content | Owner deletes module |
| Instructor | Other's course | 403 Forbidden | Cannot modify others |
| Student | Any | 403 Forbidden | Students cannot delete |

---

## 6. Lesson Management Endpoints

### POST /api/modules/:moduleId/lessons
**Purpose**: Create lesson in module (creates version 1)

| Requester | Course Owner | Expected Status | Reason |
|-----------|--------------|-----------------|---------|
| Admin | Any | 201 Created | Admin creates lessons |
| Instructor | Own course | 201 Created | Owner creates lessons |
| Instructor | Other's course | 403 Forbidden | Cannot modify others |
| Student | Any | 403 Forbidden | Students cannot create |

**Validations**:
- Module must exist
- Content or video_url required (at least one)
- First version created as v1

### GET /api/modules/:moduleId/lessons
**Purpose**: List lessons in module (latest versions only)

| Requester | Enrollment Status | Expected Status | Reason |
|-----------|-------------------|-----------------|---------|
| Admin | N/A | 200 OK | Admin views all |
| Instructor | Own course | 200 OK | Owner views lessons |
| Student | Enrolled | 200 OK | Enrolled students view |
| Student | Not enrolled | 403 Forbidden | Must be enrolled |

### GET /api/lessons/:id
**Purpose**: Get latest lesson version

| Requester | Enrollment Status | Expected Status | Reason |
|-----------|-------------------|-----------------|---------|
| Admin | N/A | 200 OK | Admin views all |
| Instructor | Own course | 200 OK | Owner views lesson |
| Student | Enrolled | 200 OK | Enrolled students view |
| Student | Not enrolled | 403 Forbidden | Must be enrolled |

### POST /api/lessons/:lessonId/version
**Purpose**: Create new lesson version

| Requester | Course Owner | Expected Status | Reason |
|-----------|--------------|-----------------|---------|
| Admin | Any | 201 Created | Admin creates versions |
| Instructor | Own course | 201 Created | Owner creates versions |
| Instructor | Other's course | 403 Forbidden | Cannot modify others |
| Student | Any | 403 Forbidden | Students cannot create versions |

**Validations**:
- Lesson must exist
- Version number auto-increments
- Previous version remains in audit trail

### GET /api/lessons/:id/all-versions
**Purpose**: Get all lesson versions (audit trail)

| Role | Expected Status | Reason |
|------|-----------------|---------|
| Admin | 200 OK | Admin can view audit trail |
| Instructor | 403 Forbidden | Cannot view audit trail |
| Student | 403 Forbidden | Cannot view audit trail |

**Purpose**: Only admin should see version history for auditing

### DELETE /api/lessons/:id
**Purpose**: Soft delete lesson (all versions)

| Requester | Course Owner | Expected Status | Reason |
|-----------|--------------|-----------------|---------|
| Admin | Any | 204 No Content | Admin deletes lessons |
| Instructor | Own course | 204 No Content | Owner deletes lessons |
| Instructor | Other's course | 403 Forbidden | Cannot modify others |
| Student | Any | 403 Forbidden | Students cannot delete |

---

## 7. Certification Webhook

### POST /api/certification-webhook
**Purpose**: Receive certification results from external provider

**Authentication**: HMAC-SHA256 signature validation (NO Bearer token)

| Scenario | Signature | Expected Status | Reason |
|----------|-----------|-----------------|---------|
| Valid payload + signature | Correct HMAC | 200 OK | Webhook processed |
| Valid payload + invalid signature | Wrong HMAC | 401 Unauthorized | Signature mismatch |
| Missing signature header | None | 401 Unauthorized | Signature required |
| With Bearer token | Valid JWT | Should fail | Webhook uses HMAC only |

**Payload Requirements**:
```json
{
  "student_id": "uuid",
  "course_id": "uuid",
  "certification_status": "passed|failed",
  "score": 0-100,
  "timestamp": "ISO8601"
}
```

**Validations**:
- student_id must exist
- course_id must exist
- Student must be enrolled with status 'active'
- certification_status must be 'passed' or 'failed'
- score must be 0-100
- If 'passed', enrollment status updated to 'completed'
- Audit log created with webhook action

**HMAC Calculation**:
```
secret = process.env.WEBHOOK_SECRET || "default-webhook-secret-change-in-production"
signature = HMAC-SHA256(secret, payload_json)
header = hex(signature)
```

---

## 8. Authorization Logic Summary

### Permission Patterns

**1. Admin Override**
- Admin has full access to all endpoints (except enrollment which is student-only)
- Admin can perform CRUD operations on any resource
- Admin can view audit trails and version history

**2. Resource Ownership**
- Instructors can only modify their own courses and associated resources
- Users can only update their own profile
- Students can only view their own enrollment data

**3. Enrollment Requirements**
- Students must be enrolled to view course content (modules/lessons)
- Enrollment is verified by checking `enrollments` table
- Enrollment status must be 'active' for content access

**4. Role Restrictions**
- Students: Enroll only, cannot create content
- Instructors: Create and manage own courses, cannot enroll
- Admin: Full system access, cannot enroll (business rule)

---

## 9. Test Execution Checklist

### Manual Testing Steps

#### Setup
- [ ] Start PostgreSQL database
- [ ] Run migrations
- [ ] Start API server
- [ ] Prepare HTTP client (Postman/curl)

#### Authentication Tests
- [ ] Register all test users (admin, 2 instructors, 2 students)
- [ ] Login each user and save JWT tokens
- [ ] Test invalid credentials rejection
- [ ] Test missing token returns 401

#### User Permissions
- [ ] Admin can list all users
- [ ] Instructor cannot list users (403)
- [ ] Student cannot list users (403)
- [ ] Users can update own profile
- [ ] Users cannot update other profiles
- [ ] Only admin can delete users

#### Course Permissions
- [ ] Instructor creates course (201)
- [ ] Instructor cannot create course for other instructor
- [ ] Student cannot create course (403)
- [ ] Instructor updates own course (200)
- [ ] Instructor cannot update other's course (403)
- [ ] Instructor deletes own course (204)
- [ ] Instructor cannot delete other's course (403)

#### Enrollment Permissions
- [ ] Student enrolls in course (201)
- [ ] Student cannot enroll twice (400)
- [ ] Instructor cannot enroll (400)
- [ ] Course owner can view enrolled students
- [ ] Non-owner cannot view enrolled students (403)
- [ ] Student can view own courses
- [ ] Student cannot view other student's courses (403)

#### Module Permissions
- [ ] Course owner creates module (201)
- [ ] Non-owner cannot create module (403)
- [ ] Enrolled student can view modules
- [ ] Non-enrolled student cannot view modules (403)
- [ ] Course owner updates/deletes modules
- [ ] Non-owner cannot update/delete (403)

#### Lesson Permissions
- [ ] Course owner creates lesson (201)
- [ ] Non-owner cannot create lesson (403)
- [ ] Enrolled student views lesson content
- [ ] Non-enrolled cannot view content (403)
- [ ] Owner creates new version (v2, v3...)
- [ ] Student cannot create versions (403)
- [ ] Only admin views all versions (audit trail)
- [ ] Instructor cannot view all versions (403)

#### Webhook Security
- [ ] Valid HMAC signature processes webhook (200)
- [ ] Invalid signature rejects webhook (401)
- [ ] Missing signature rejects webhook (401)
- [ ] Bearer token does not work on webhook
- [ ] Enrollment status updated on 'passed'
- [ ] Audit log created for webhook

---

## 10. Expected Database State After Tests

### Users Table
```
5 users:
- 1 admin
- 2 instructors  
- 2 students
```

### Courses Table
```
Multiple courses created by instructors
Some may be soft-deleted (DeletedAt != null)
```

### Enrollments Table
```
Students enrolled in various courses
Some with status: active
Some with status: completed (via webhook)
```

### Audit_logs Table
```
Entries for:
- Course created/updated/deleted
- Enrollment created/updated
- Lesson versions created
- Certification webhooks processed
```

---

## 11. Common Error Codes

| Status | Meaning | When It Occurs |
|--------|---------|----------------|
| 200 OK | Success | GET requests successful |
| 201 Created | Resource created | POST successful |
| 204 No Content | Deleted | DELETE successful |
| 400 Bad Request | Invalid input | Validation failure |
| 401 Unauthorized | Not authenticated | Missing/invalid JWT or HMAC |
| 403 Forbidden | No permission | Role/ownership check failed |
| 404 Not Found | Resource missing | Invalid ID in URL |
| 500 Internal Error | Server error | Unexpected failure |

---

## 12. Security Validations

### UUID Validation
All ID parameters must be valid UUIDs (google/uuid format)

### JWT Token Validation
- Token must be present in `Authorization: Bearer <token>` header
- Token must not be expired
- Token must have valid signature
- User ID from token must exist in database
- User must not be soft-deleted

### HMAC Signature Validation (Webhook Only)
- Signature must be in `X-Webhook-Signature` header
- Must match HMAC-SHA256(secret, body)
- Secret from env: `WEBHOOK_SECRET`

### Role Validation
- Role must be: student, instructor, or admin
- Role cannot be changed by non-admin users
- Instructor_id must have instructor role

### Enrollment Validation
- Student can only have one active enrollment per course
- Must be enrolled to access course content
- Status must be: active, completed, or dropped

---

## Testing Recommendations

1. **Use Postman Collections**: Create collection with all endpoints and test cases
2. **Environment Variables**: Store tokens and IDs as variables
3. **Test Order**: Run in sequence (create resources before testing permissions)
4. **Clean State**: Reset database between test runs
5. **Error Messages**: Verify error messages are specific and helpful
6. **Audit Logs**: Check that all mutations create audit entries
7. **Soft Deletes**: Verify DeletedAt is set, not hard deleted
8. **Performance**: Test with multiple enrollments and nested queries

---

## Summary

✅ **Total Test Scenarios: 150+**

This specification covers comprehensive RBAC testing including:
- Authentication and authorization
- Role-based permissions (Admin, Instructor, Student)
- Resource ownership validation
- Enrollment requirements for content access
- Webhook security with HMAC validation
- Audit logging verification
- Error handling and validation

All permissions are enforced at the HTTP handler layer through middleware and explicit checks in use cases.
