# Certification Provider Webhook

## Overview

The certification webhook allows external certification providers to notify the LMS when students complete exams, automatically updating enrollment status.

## Endpoint

```
POST /api/certification-webhook
```

**Note:** This endpoint does NOT require authentication (no Bearer token needed) but requires HMAC signature validation.

## Request

### Headers

- `Content-Type: application/json` (required)
- `X-Webhook-Signature: <hmac_signature>` (required)

### Payload

```json
{
  "student_id": "uuid-of-student",
  "course_id": "uuid-of-course",
  "certification_status": "passed|failed",
  "score": 0-100,
  "timestamp": "2025-12-17T10:30:00Z"
}
```

**Field Descriptions:**

- `student_id` (string, required): UUID of the student
- `course_id` (string, required): UUID of the course
- `certification_status` (string, required): Must be either "passed" or "failed"
- `score` (integer, required): Score between 0 and 100
- `timestamp` (string, required): ISO 8601 formatted timestamp

## Response

### Success (200 OK)

```json
{
  "message": "Certification webhook processed successfully"
}
```

### Errors

#### 400 Bad Request

```json
{
  "error": "Missing required fields"
}
```

Possible error messages:
- "Failed to read request body"
- "Invalid JSON payload"
- "Missing required fields"
- "certification_status must be 'passed' or 'failed'"
- "score must be between 0 and 100"

#### 401 Unauthorized

```json
{
  "error": "Invalid webhook signature"
}
```

Possible error messages:
- "Missing webhook signature"
- "Invalid webhook signature"

#### 404 Not Found

```json
{
  "error": "student not found"
}
```

Possible error messages:
- "student not found"
- "course not found"
- "enrollment not found"

#### 500 Internal Server Error

```json
{
  "error": "error message"
}
```

## Security: HMAC Signature

The webhook uses HMAC-SHA256 for authentication. The signature is computed as:

```
HMAC-SHA256(webhook_secret, request_body)
```

### Configuration

Set the webhook secret via environment variable:

```bash
export WEBHOOK_SECRET="your-secret-key-here"
```

**Default:** `default-webhook-secret-change-in-production` (change this in production!)

### Generating Signature (Examples)

#### Bash/OpenSSL

```bash
echo -n '{"student_id":"...","course_id":"...","certification_status":"passed","score":95,"timestamp":"2025-12-17T10:30:00Z"}' | \
  openssl dgst -sha256 -hmac "your-webhook-secret" | \
  awk '{print $2}'
```

#### Python

```python
import hmac
import hashlib
import json

payload = {
    "student_id": "uuid",
    "course_id": "uuid",
    "certification_status": "passed",
    "score": 95,
    "timestamp": "2025-12-17T10:30:00Z"
}

secret = "your-webhook-secret"
message = json.dumps(payload, separators=(',', ':'))

signature = hmac.new(
    secret.encode('utf-8'),
    message.encode('utf-8'),
    hashlib.sha256
).hexdigest()

print(f"X-Webhook-Signature: {signature}")
```

#### Node.js

```javascript
const crypto = require('crypto');

const payload = {
  student_id: "uuid",
  course_id: "uuid",
  certification_status: "passed",
  score: 95,
  timestamp: "2025-12-17T10:30:00Z"
};

const secret = "your-webhook-secret";
const message = JSON.stringify(payload);

const signature = crypto
  .createHmac('sha256', secret)
  .update(message)
  .digest('hex');

console.log(`X-Webhook-Signature: ${signature}`);
```

## Behavior

1. **Validates HMAC signature** - Returns 401 if signature is missing or invalid
2. **Validates payload format** - Returns 400 if JSON is malformed or fields are invalid
3. **Validates student exists** - Returns 404 if student UUID not found
4. **Validates course exists** - Returns 404 if course UUID not found
5. **Validates enrollment exists** - Returns 404 if no enrollment found for student-course pair
6. **Checks enrollment status** - Only processes enrollments with `status = active`
7. **Updates enrollment** - If `certification_status = "passed"`, sets `enrollment.status = completed`
8. **Creates audit log** - Records the webhook event with student, course, and score details
9. **Returns success** - Returns 200 OK

### Important Notes

- Only enrollments with `status = "active"` are updated
- If `certification_status = "failed"`, the enrollment status is NOT changed
- If `certification_status = "passed"`, enrollment is marked as completed
- All webhook events are logged in the audit_log table for tracking

## Testing

Use the provided test script:

```bash
# 1. Start the server
go run ./src/cmd/api

# 2. Update test_webhook.sh with valid UUIDs
# 3. Run the test
./test_webhook.sh
```

Or use curl directly:

```bash
PAYLOAD='{"student_id":"uuid","course_id":"uuid","certification_status":"passed","score":95,"timestamp":"2025-12-17T10:30:00Z"}'
SIGNATURE=$(echo -n "$PAYLOAD" | openssl dgst -sha256 -hmac "default-webhook-secret-change-in-production" | awk '{print $2}')

curl -X POST http://localhost:8080/api/certification-webhook \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Signature: $SIGNATURE" \
  -d "$PAYLOAD"
```

## Database Schema

### Audit Log Table

The webhook creates entries in the `audit_logs` table:

```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    action VARCHAR(100) NOT NULL,
    entity_id UUID NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    details TEXT,
    user_id UUID,
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);
```

Example audit log entry:

```json
{
  "id": "audit-uuid",
  "action": "certification_webhook",
  "entity_id": "student-uuid-course-uuid",
  "entity_type": "enrollment",
  "details": "Certification webhook processed: status=passed, score=95, student=uuid (John Doe), course=uuid (Introduction to Go)",
  "user_id": "student-uuid",
  "created_at": "2025-12-17T10:30:00Z"
}
```

## Architecture

```
External Provider
       |
       | POST with HMAC signature
       v
Certification Webhook Handler
       |
       | Validates signature, payload
       v
Certification Webhook Use Case
       |
       +-- Validates student exists (UserRepository)
       +-- Validates course exists (CourseRepository)
       +-- Gets enrollment (EnrollmentRepository)
       +-- Updates enrollment status if passed
       +-- Creates audit log (AuditLogRepository)
       v
   Returns 200 OK
```

## Files Created

- `src/internal/domain/entity/audit_log.go` - Audit log entity
- `src/internal/domain/repository/audit_log_repository.go` - Repository interface
- `src/internal/infrastructure/repository/audit_log_repository_postgres.go` - PostgreSQL implementation
- `src/usecase/certification_webhook_usecase.go` - Business logic
- `src/internal/adapter/http/handler/certification_webhook_handler.go` - HTTP handler
- `test_webhook.sh` - Test script
- `CERTIFICATION_WEBHOOK.md` - This documentation

## Configuration

Environment variables:

- `WEBHOOK_SECRET` - Shared secret for HMAC validation (required in production)

## Security Considerations

1. **Always change the default webhook secret in production**
2. **Use HTTPS in production** to protect the payload in transit
3. **Validate the signature before processing** any webhook data
4. **Log all webhook attempts** for security auditing
5. **Consider rate limiting** to prevent abuse
6. **Implement webhook retry logic** on the provider side with exponential backoff
7. **Monitor failed webhook attempts** for potential security issues
