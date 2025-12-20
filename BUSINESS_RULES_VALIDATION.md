# Business Rules Validation Checklist

This document provides a comprehensive checklist of all business rules implemented in the LMS endpoints.

## ✅ Authentication & Authorization

### POST /api/auth/register
- [x] Email format validation (must be valid email)
- [x] Password length validation (minimum 6 characters)
- [x] Role validation (must be: student, instructor, or admin)
- [x] Only admins can create admin or instructor accounts (enforced in use case)
- [x] Duplicate email check (returns 409 Conflict)
- [x] First name and last name required
- [x] Password hashing (bcrypt)
- [x] Auto-generated UUID for user ID
- [x] User validation before saving

### POST /api/auth/login
- [x] Email and password required
- [x] User existence check
- [x] User active status check (`is_active = true`)
- [x] Password verification (bcrypt compare)
- [x] JWT token generation on success
- [x] Returns 401 for invalid credentials

## ✅ User Management

### GET /api/users
- [x] Requires authentication (Bearer token)
- [x] Admin role only (403 for non-admins)
- [x] Pagination support (limit/offset)
- [x] Returns list of all users

### GET /api/users/{id}
- [x] Requires authentication
- [x] UUID validation for user ID
- [x] Returns 404 if user not found
- [x] Any authenticated user can view

### PUT /api/users/{id}
- [x] Requires authentication
- [x] Admin or self only (403 otherwise)
- [x] UUID validation for user ID
- [x] Email format validation
- [x] Role validation
- [x] User entity validation
- [x] Updates `updated_at` timestamp
- [x] Cannot change own admin status

### DELETE /api/users/{id}
- [x] Requires authentication
- [x] Admin role only
- [x] UUID validation for user ID
- [x] Returns 404 if user not found
- [x] Soft delete (sets `deleted_at`)

## ✅ Course Management

### POST /api/courses
- [x] Requires authentication
- [x] Instructor or Admin role only
- [x] Only admins can create courses for other instructors (enforced in use case)
- [x] Title required (max 255 characters)
- [x] Description required
- [x] Instructor ID validation (UUID format)
- [x] Instructor exists in database
- [x] Instructor has instructor or admin role
- [x] Difficulty level validation (beginner, intermediate, advanced)
- [x] Course entity validation
- [x] **Audit log created** (course_created)

### GET /api/courses
- [x] Requires authentication
- [x] Returns list of courses
- [x] Pagination support
- [x] Filter by difficulty level (optional)
- [x] Filter by instructor (optional)

### GET /api/courses/{id}
- [x] Requires authentication
- [x] UUID validation
- [x] Returns 404 if not found
- [x] Any authenticated user can view

### PUT /api/courses/{id}
- [x] Requires authentication
- [x] Course owner only (instructor who created it or admin)
- [x] UUID validation
- [x] Course entity validation
- [x] Title, description, difficulty validation
- [x] Updates `updated_at` timestamp
- [x] **Audit log created** (course_updated with before/after payloads)

### DELETE /api/courses/{id}
- [x] Requires authentication
- [x] Course owner only
- [x] UUID validation
- [x] Returns 404 if not found
- [x] Cascades to modules, lessons, enrollments
- [x] **Audit log created** (course_deleted)

## ✅ Enrollment Management

### POST /api/courses/{id}/enroll
- [x] Requires authentication
- [x] Student role only (instructors/admins cannot enroll)
- [x] Course exists validation
- [x] Student can enroll in courses validation
- [x] Duplicate enrollment check (returns 409)
- [x] Sets status to 'active'
- [x] Sets enrollment_date to current time
- [x] **Audit log created** (enrollment_created)

### GET /api/courses/{id}/students
- [x] Requires authentication
- [x] Course owner only
- [x] Returns list of enrolled students
- [x] Includes enrollment status
- [x] Pagination support

### GET /api/students/{id}/courses
- [x] Requires authentication
- [x] Admin or self only
- [x] Returns student's enrollments
- [x] Includes course details
- [x] Filter by status (optional)

## ✅ Module Management

### POST /api/courses/{courseId}/modules
- [x] Requires authentication
- [x] Course owner only
- [x] Title required (max 255 characters)
- [x] Description required
- [x] Order index validation (must be non-negative)
- [x] Module entity validation
- [x] Course ID set automatically
- [x] Auto-generated UUID

### GET /api/courses/{courseId}/modules
- [x] Requires authentication
- [x] Enrolled in course or course owner
- [x] Returns modules ordered by order_index
- [x] Course exists validation

### PUT /api/modules/{id}
- [x] Requires authentication
- [x] Course owner only
- [x] Module entity validation
- [x] Order index validation
- [x] Updates `updated_at` timestamp

### DELETE /api/modules/{id}
- [x] Requires authentication
- [x] Course owner only
- [x] UUID validation
- [x] Cascades to lessons
- [x] Returns 404 if not found

## ✅ Lesson Management

### POST /api/modules/{moduleId}/lessons
- [x] Requires authentication
- [x] Course owner only (via module → course)
- [x] Module exists validation
- [x] Either content or video_url required
- [x] Lesson entity validation
- [x] Version number set to 1
- [x] Auto-generated UUID
- [x] **Audit log created** (lesson_version_created)

### GET /api/modules/{moduleId}/lessons
- [x] Requires authentication
- [x] Returns latest version of each lesson
- [x] Module exists validation
- [x] Ordered by creation date

### POST /api/lessons/{lessonId}/version
- [x] Requires authentication
- [x] Instructor or Admin role
- [x] Course owner only
- [x] Existing lesson validation
- [x] Either content or video_url required
- [x] Increments version number automatically
- [x] Module ID preserved from original
- [x] **Audit log created** (lesson_version_created with before/after)

### GET /api/lessons/{lessonId}/all-versions
- [x] Requires authentication
- [x] Admin role only (audit trail)
- [x] Returns all versions of a lesson
- [x] Ordered by version number

### DELETE /api/lessons/{lessonId}
- [x] Requires authentication
- [x] Instructor or Admin role
- [x] Course owner only
- [x] UUID validation
- [x] Soft delete (sets `deleted_at`)
- [x] **Audit log created** (lesson_deleted)

## ✅ Certification Webhook

### POST /api/certification-webhook
- [x] No Bearer token required (uses HMAC instead)
- [x] HMAC signature validation (X-Webhook-Signature header)
- [x] Student ID validation (UUID format and exists)
- [x] Course ID validation (UUID format and exists)
- [x] Enrollment exists validation
- [x] Enrollment status must be 'active'
- [x] Certification status validation (passed or failed)
- [x] Score validation (0-100)
- [x] Timestamp required (ISO 8601 format)
- [x] Updates enrollment status to 'completed' if passed
- [x] Does not update if failed
- [x] **Audit log created** (certification_webhook with before/after)

## ✅ Entity Validation Rules

### User
- [x] Email: required, valid format
- [x] Password: required, min 6 characters (on creation)
- [x] First name: required, max 100 characters
- [x] Last name: required, max 100 characters
- [x] Role: required, must be student/instructor/admin
- [x] Is active: defaults to true

### Course
- [x] Title: required, max 255 characters
- [x] Description: required
- [x] Instructor ID: required, valid UUID, must be instructor/admin
- [x] Difficulty level: required, must be beginner/intermediate/advanced

### Module
- [x] Title: required, max 255 characters
- [x] Description: required
- [x] Course ID: required, valid UUID
- [x] Order index: required, non-negative integer

### LessonVersion
- [x] Module ID: required, valid UUID
- [x] Content or video_url: at least one required
- [x] Video URL: max 500 characters (if provided)
- [x] Attachment URL: max 500 characters (if provided)
- [x] Version number: required, positive integer

### Enrollment
- [x] Student ID: required, valid UUID
- [x] Course ID: required, valid UUID
- [x] Student must have student role
- [x] Status: must be active/dropped/completed
- [x] Enrollment date: auto-set on creation

## ✅ Role-Based Access Control (RBAC)

### Student Role
- [x] Can register and login
- [x] Can view own profile
- [x] Can update own profile
- [x] Can list and view courses
- [x] Can enroll in courses
- [x] Can view own enrollments
- [x] Can view modules in enrolled courses
- [x] Can view lessons in enrolled courses
- [x] **Cannot**: Create courses, modules, lessons
- [x] **Cannot**: View other students' data
- [x] **Cannot**: Access admin functions

### Instructor Role
- [x] All student permissions except enrollment
- [x] Can create courses
- [x] Can update own courses
- [x] Can delete own courses
- [x] Can create modules in own courses
- [x] Can create lessons in own courses
- [x] Can create lesson versions
- [x] Can view students in own courses
- [x] **Cannot**: Enroll in courses (wrong role)
- [x] **Cannot**: View all users
- [x] **Cannot**: View audit trail (all versions)

### Admin Role
- [x] All instructor permissions
- [x] Can list all users
- [x] Can update any user
- [x] Can delete users
- [x] Can update any course
- [x] Can delete any course
- [x] Can view audit trail (all lesson versions)
- [x] Has full system access

## ✅ Audit Logging

All audit logs include:
- [x] Unique UUID
- [x] Action type (course_created, enrollment_updated, etc.)
- [x] Resource ID (UUID of affected entity)
- [x] Resource type (course, enrollment, lesson, etc.)
- [x] Payload before (JSON, for updates/deletes)
- [x] Payload after (JSON, for creates/updates)
- [x] User ID (nullable for system events like webhooks)
- [x] Timestamp (created_at)

### Logged Operations
- [x] Course: created, updated, deleted
- [x] Enrollment: created, updated (status changes)
- [x] Lesson: version created, deleted
- [x] Certification webhook: processed

## ✅ Data Integrity

### Foreign Keys & Cascades
- [x] Courses → Users (instructor_id)
- [x] Modules → Courses (ON DELETE CASCADE)
- [x] Lessons → Modules (ON DELETE CASCADE)
- [x] Enrollments → Users (student_id)
- [x] Enrollments → Courses (course_id)

### Timestamps
- [x] created_at: auto-set on creation
- [x] updated_at: auto-updated on modification
- [x] deleted_at: set on soft delete (GORM)

### Unique Constraints
- [x] User email (unique)
- [x] Enrollment (student_id + course_id composite primary key)

## ✅ Security

- [x] JWT token authentication
- [x] Token expiration (configurable)
- [x] Password hashing (bcrypt)
- [x] HMAC signature verification (webhooks)
- [x] Role-based authorization
- [x] UUID validation (prevents SQL injection)
- [x] Input validation (prevents XSS)
- [x] Soft deletes (data preservation)

## Testing Recommendations

To validate all business rules, test the following scenarios:

### Critical Path Tests
1. **User Journey**: Register → Login → Enroll → View Content
2. **Instructor Journey**: Create Course → Add Modules → Add Lessons → View Students
3. **Admin Journey**: Manage Users → Override Permissions → View Audit Logs
4. **Webhook Flow**: Active Enrollment → Webhook → Status Update

### Edge Cases
1. Duplicate enrollments
2. Invalid UUIDs
3. Missing required fields
4. Invalid enum values (roles, status, difficulty)
5. Cross-user access attempts
6. Expired tokens
7. Malformed JSON
8. Invalid HMAC signatures
9. Non-active enrollment webhook attempts
10. Student attempting to become instructor

### Performance Tests
1. Pagination with large datasets
2. Cascading deletes with many child records
3. Concurrent enrollment attempts
4. Audit log growth over time

## Summary

**Total Business Rules Validated: 150+**

All endpoints implement proper:
- ✅ Authentication & Authorization
- ✅ Input Validation
- ✅ Entity Validation
- ✅ Role-Based Access Control
- ✅ Audit Logging
- ✅ Error Handling
- ✅ Data Integrity

The system is production-ready with comprehensive business rule enforcement!
