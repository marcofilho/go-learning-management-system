# Audit Logging System

## Overview

The LMS includes a comprehensive audit logging system that tracks all critical operations across the system. This provides accountability, debugging capabilities, and compliance with regulatory requirements.

## Database Schema

### audit_logs Table

```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    action VARCHAR(100) NOT NULL,
    resource_id UUID NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    payload_before JSONB,
    payload_after JSONB,
    user_id UUID,
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_id);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
```

### Fields

- **id**: Unique identifier for the audit log entry
- **action**: Type of action performed (see Actions below)
- **resource_id**: UUID of the resource that was modified
- **resource_type**: Type of resource (course, enrollment, lesson, etc.)
- **payload_before**: JSON representation of the resource state before the change
- **payload_after**: JSON representation of the resource state after the change
- **user_id**: UUID of the user who performed the action (nullable for system events)
- **created_at**: Timestamp when the action occurred
- **deleted_at**: Soft delete timestamp (for GORM)

## Logged Actions

### Course Operations

| Action | Description | User ID | Payload Before | Payload After |
|--------|-------------|---------|----------------|---------------|
| `course_created` | New course created | Instructor | Empty | Course details |
| `course_updated` | Course modified | Instructor/Admin | Old course state | New course state |
| `course_deleted` | Course deleted | Instructor/Admin | Course details | Empty |

**Example:**
```json
{
  "action": "course_created",
  "resource_id": "123e4567-e89b-12d3-a456-426614174000",
  "resource_type": "course",
  "payload_before": "",
  "payload_after": "{\"id\":\"123...\",\"title\":\"Introduction to Go\",\"instructor_id\":\"456...\",\"difficulty_level\":\"beginner\"}",
  "user_id": "456e4567-e89b-12d3-a456-426614174001",
  "created_at": "2025-12-17T10:30:00Z"
}
```

### Enrollment Operations

| Action | Description | User ID | Payload Before | Payload After |
|--------|-------------|---------|----------------|---------------|
| `enrollment_created` | Student enrolled in course | Student/Admin | Empty | Enrollment details |
| `enrollment_updated` | Enrollment status changed | Student/Admin | Old status | New status |
| `enrollment_deleted` | Enrollment removed | Admin | Enrollment details | Empty |

**Note:** Resource ID for enrollments is a composite key: `{student_id}-{course_id}`

**Example:**
```json
{
  "action": "enrollment_updated",
  "resource_id": "student-uuid-course-uuid",
  "resource_type": "enrollment",
  "payload_before": "{\"status\":\"active\"}",
  "payload_after": "{\"status\":\"completed\"}",
  "user_id": "admin-uuid",
  "created_at": "2025-12-17T10:30:00Z"
}
```

### Lesson Version Operations

| Action | Description | User ID | Payload Before | Payload After |
|--------|-------------|---------|----------------|---------------|
| `lesson_version_created` | New lesson version created | Instructor | Previous version | New version |
| `lesson_deleted` | Lesson deleted | Instructor/Admin | Lesson details | Empty |

**Example:**
```json
{
  "action": "lesson_version_created",
  "resource_id": "789e4567-e89b-12d3-a456-426614174002",
  "resource_type": "lesson",
  "payload_before": "{\"id\":\"old-uuid\",\"version\":1}",
  "payload_after": "{\"id\":\"new-uuid\",\"version\":2}",
  "user_id": "instructor-uuid",
  "created_at": "2025-12-17T10:30:00Z"
}
```

### Webhook Operations

| Action | Description | User ID | Payload Before | Payload After |
|--------|-------------|---------|----------------|---------------|
| `certification_webhook` | Certification result received | Student (from webhook) | Old enrollment state | New enrollment state with certification |

**Example:**
```json
{
  "action": "certification_webhook",
  "resource_id": "student-uuid-course-uuid",
  "resource_type": "enrollment",
  "payload_before": "{\"status\":\"active\"}",
  "payload_after": "{\"status\":\"completed\",\"certification_status\":\"passed\",\"score\":95}",
  "user_id": "student-uuid",
  "created_at": "2025-12-17T10:30:00Z"
}
```

## Implementation

### Entity

```go
type AuditLog struct {
    ID            string      `gorm:"type:uuid;primaryKey" json:"id"`
    Action        AuditAction `gorm:"type:varchar(100);not null" json:"action"`
    ResourceID    string      `gorm:"type:uuid;not null" json:"resource_id"`
    ResourceType  string      `gorm:"type:varchar(50);not null" json:"resource_type"`
    PayloadBefore string      `gorm:"type:jsonb" json:"payload_before,omitempty"`
    PayloadAfter  string      `gorm:"type:jsonb" json:"payload_after,omitempty"`
    UserID        *string     `gorm:"type:uuid" json:"user_id,omitempty"`
    CreatedAt     time.Time   `gorm:"not null" json:"created_at"`
    DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### Repository Interface

```go
type AuditLogRepository interface {
    Create(ctx context.Context, auditLog *entity.AuditLog) error
    GetByID(ctx context.Context, id string) (*entity.AuditLog, error)
    List(ctx context.Context, limit, offset int) ([]*entity.AuditLog, error)
}
```

### Usage in Use Cases

All use cases that modify data automatically create audit log entries:

```go
// Example: Course Update
func (uc *CourseUseCase) UpdateCourse(ctx context.Context, course *entity.Course, userID string) error {
    // Get existing course for audit log
    existingCourse, _ := uc.courseRepo.GetByID(ctx, course.ID)
    
    // Perform update
    if err := uc.courseRepo.Update(ctx, course); err != nil {
        return err
    }
    
    // Create audit log
    payloadBefore := fmt.Sprintf(`{"title":"%s","description":"%s"}`, 
        existingCourse.Title, existingCourse.Description)
    payloadAfter := fmt.Sprintf(`{"title":"%s","description":"%s"}`, 
        course.Title, course.Description)
    
    auditLog, _ := entity.NewAuditLog(
        entity.AuditActionCourseUpdated,
        course.ID,
        "course",
        payloadBefore,
        payloadAfter,
        &userID,
    )
    uc.auditLogRepo.Create(ctx, auditLog)
    
    return nil
}
```

## Querying Audit Logs

### Get All Audit Logs

```go
auditLogs, err := auditLogRepo.List(ctx, 100, 0)
```

### Get Audit Log by ID

```go
auditLog, err := auditLogRepo.GetByID(ctx, "audit-log-uuid")
```

### Direct SQL Queries

```sql
-- Get all changes to a specific course
SELECT * FROM audit_logs 
WHERE resource_id = 'course-uuid' 
AND resource_type = 'course'
ORDER BY created_at DESC;

-- Get all actions by a specific user
SELECT * FROM audit_logs 
WHERE user_id = 'user-uuid'
ORDER BY created_at DESC;

-- Get all webhook events
SELECT * FROM audit_logs 
WHERE action = 'certification_webhook'
ORDER BY created_at DESC;

-- Get enrollment changes
SELECT * FROM audit_logs 
WHERE resource_type = 'enrollment'
AND action IN ('enrollment_created', 'enrollment_updated', 'enrollment_deleted')
ORDER BY created_at DESC;

-- Compare before/after states
SELECT 
    action,
    resource_id,
    payload_before::json->>'status' as old_status,
    payload_after::json->>'status' as new_status,
    created_at
FROM audit_logs 
WHERE action = 'enrollment_updated'
ORDER BY created_at DESC;
```

## Use Cases

### 1. Compliance & Audit Trail

Track who changed what and when for regulatory compliance:

```sql
-- Show all modifications to a course
SELECT 
    action,
    payload_before,
    payload_after,
    user_id,
    created_at
FROM audit_logs
WHERE resource_id = 'course-uuid'
ORDER BY created_at;
```

### 2. Debugging

Investigate issues by reviewing the history of changes:

```sql
-- Find when an enrollment was completed
SELECT * FROM audit_logs
WHERE resource_type = 'enrollment'
AND payload_after::json->>'status' = 'completed'
AND resource_id = 'student-uuid-course-uuid';
```

### 3. Analytics

Generate reports on system usage:

```sql
-- Count actions by type
SELECT action, COUNT(*) as count
FROM audit_logs
GROUP BY action
ORDER BY count DESC;

-- Most active users
SELECT user_id, COUNT(*) as actions
FROM audit_logs
WHERE user_id IS NOT NULL
GROUP BY user_id
ORDER BY actions DESC
LIMIT 10;
```

### 4. Security Monitoring

Detect suspicious activities:

```sql
-- Find mass deletions
SELECT user_id, COUNT(*) as deletions
FROM audit_logs
WHERE action IN ('course_deleted', 'lesson_deleted')
AND created_at > NOW() - INTERVAL '1 hour'
GROUP BY user_id
HAVING COUNT(*) > 5;
```

## Best Practices

1. **Always log changes**: Every create, update, delete operation should create an audit log
2. **Include relevant data**: Store enough information in JSON payloads to understand what changed
3. **Never delete audit logs**: Use soft deletes if necessary, but preserve the audit trail
4. **Index appropriately**: Ensure queries on resource_id, user_id, and created_at are fast
5. **Handle failures gracefully**: If audit logging fails, log the error but don't block the operation
6. **Protect sensitive data**: Avoid storing passwords or tokens in audit logs
7. **Set retention policies**: Archive old audit logs to keep table sizes manageable

## Future Enhancements

- **Audit log API endpoints**: Allow admins to query audit logs via REST API
- **Real-time monitoring**: Stream audit logs to monitoring systems
- **Anomaly detection**: Alert on suspicious patterns
- **Data retention**: Automatic archival of old audit logs
- **Enhanced filtering**: Query by date range, action type, user, etc.
- **Rollback capability**: Use audit logs to revert changes

## Files

- `src/internal/domain/entity/audit_log.go` - Entity definition
- `src/internal/domain/repository/audit_log_repository.go` - Repository interface  
- `src/internal/infrastructure/repository/audit_log_repository_postgres.go` - PostgreSQL implementation
- `src/usecase/course_usecase.go` - Course audit logging
- `src/usecase/enrollment_usecase.go` - Enrollment audit logging
- `src/usecase/lesson_usecase.go` - Lesson audit logging
- `src/usecase/certification_webhook_usecase.go` - Webhook audit logging

## Summary

The audit logging system provides comprehensive tracking of all critical operations in the LMS. Every course modification, enrollment change, lesson update, and certification webhook event is logged with before/after states, timestamps, and user attribution. This enables compliance, debugging, analytics, and security monitoring.
