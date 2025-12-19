# Integration Test Workflow Guide

This document explains how to run the complete integration test suite for the LMS API.

## Overview

The `test_workflow.sh` script tests the entire LMS API workflow sequentially, covering:

1. **Authentication** - User registration and login (admin, instructor, student)
2. **Course Management** - Create, list, retrieve, and update courses
3. **Module Management** - Create and list modules within courses
4. **Lesson Management** - Create lessons (v1) and additional versions, list all versions
5. **Enrollment** - Student self-enrollment, view enrollments
6. **Webhook Processing** - Certification webhook testing
7. **Audit Logging** - View and verify audit trail
8. **Authorization** - Test permission boundaries

## Quick Start

### Option 1: Using Docker (Recommended)

```bash
# Start all services (API + Database)
make start

# In another terminal, run the integration tests
make test-workflow
```

### Option 2: Manual Local Setup

```bash
# 1. Start PostgreSQL database
docker-compose up -d postgres

# 2. Wait for database to be ready
sleep 3

# 3. Run migrations
make migrate-up

# 4. Start the API (in another terminal)
make run

# 5. Run the integration tests (in a third terminal)
./test_workflow.sh
```

### Option 3: Full End-to-End Test

```bash
# Ensures Docker is running and executes all tests
make test-integration
```

## Configuration

The test script accepts environment variables for configuration:

```bash
# API endpoint (default: http://localhost:8080/api)
export API_URL="http://localhost:8080/api"

# Database configuration
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_NAME="lms_db"

# Run tests
./test_workflow.sh
```

## Test Execution Flow

### Phase 1: Pre-Flight Checks
- ✅ Verify API is running
- ✅ Verify database connectivity

### Phase 2: Authentication (Tests 1-6)
- ✅ Register admin user
- ✅ Register instructor user
- ✅ Register student user
- ✅ Admin login
- ✅ Instructor login
- ✅ Student login

### Phase 3: Course Management (Tests 7-9)
- ✅ Create course as instructor
- ✅ List courses
- ✅ Get course details
- ✅ Update course

### Phase 4: Module Management (Tests 10-11)
- ✅ Create module
- ✅ List modules

### Phase 5: Lesson Management (Tests 12-15)
- ✅ Create lesson (version 1)
- ✅ List lessons in module
- ✅ Create lesson version 2
- ✅ Get all lesson versions (admin only)

### Phase 6: Enrollment (Tests 16-18)
- ✅ Student self-enroll
- ✅ List student courses
- ✅ View course enrolled students

### Phase 7: Webhook (Tests 19-20)
- ✅ Process certification webhook (passed status)

### Phase 8: Audit Logging (Tests 21-23)
- ✅ Admin view audit logs
- ✅ Verify course creation logged
- ✅ Check audit trail entries

### Phase 9: Authorization (Tests 24-26)
- ✅ Student cannot create course (403)
- ✅ Instructor cannot view audit logs (403)
- ✅ Instructor cannot view all lesson versions (403)

## Expected Output

```
✅ ========================================
✅ LMS API Integration Test Suite
✅ ========================================

Testing API at: http://localhost:8080/api
Database: postgres@localhost:5432/lms_db

========================================
Pre-Flight Checks
========================================

TEST: Checking if API is running
✓ API is running on http://localhost:8080/api

TEST: Checking database connectivity
✓ Database is accessible

========================================
1. Authentication Tests
========================================

TEST: Register admin user
✓ Admin registered successfully

... (more tests) ...

========================================
Test Summary
========================================

Total Tests: 26
Passed: 26
Failed: 0

✓ All tests passed!
```

## Troubleshooting

### API Not Running

```bash
# Check if API is accessible
curl http://localhost:8080/api/health

# If not, start it
make run

# Or with Docker
docker-compose up -d
```

### Database Connection Failed

```bash
# Check PostgreSQL status
docker-compose ps

# View database logs
docker-compose logs postgres

# Recreate database
docker-compose down -v
docker-compose up -d
```

### Test Hangs or Times Out

```bash
# Check API logs
docker-compose logs api

# Manually test API endpoint
curl -X GET http://localhost:8080/api/health

# Increase timeout in test script (edit test_workflow.sh)
# Modify timeout in api_call function if needed
```

### Permission Denied on Script

```bash
# Make script executable
chmod +x test_workflow.sh

# Run with explicit shell
bash test_workflow.sh
```

## Manual Test Examples

If you prefer to test manually with `curl`:

### 1. Register User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email":"admin@example.com",
    "password":"SecurePass123",
    "first_name":"Admin",
    "last_name":"User",
    "role":"admin"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email":"admin@example.com",
    "password":"SecurePass123"
  }'
```

### 3. Create Course (with token from login)
```bash
TOKEN="your_token_here"

curl -X POST http://localhost:8080/api/courses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title":"Go Programming",
    "description":"Learn Go",
    "difficulty_level":"beginner"
  }'
```

### 4. List Courses
```bash
curl -X GET http://localhost:8080/api/courses \
  -H "Authorization: Bearer $TOKEN"
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: lms_db
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.24
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run migrations
        run: make migrate-up
        env:
          DB_HOST: localhost
          DB_USER: postgres
          DB_PASSWORD: postgres
          DB_NAME: lms_db
      
      - name: Start API
        run: make run &
        env:
          DB_HOST: localhost
          DB_USER: postgres
          DB_PASSWORD: postgres
          DB_NAME: lms_db
      
      - name: Wait for API
        run: sleep 5
      
      - name: Run integration tests
        run: make test-workflow
        env:
          API_URL: http://localhost:8080/api
```

## Test Script Features

### Color-Coded Output
- 🟢 **Green**: Successful tests
- 🔴 **Red**: Failed tests
- 🟡 **Yellow**: Test descriptions
- 🔵 **Blue**: Headers and information

### Test Statistics
- Track total tests run
- Count passed/failed tests
- Display summary at the end

### Graceful Error Handling
- Stops execution on critical errors (pre-flight checks)
- Continues on non-critical test failures
- Provides detailed error messages

### Flexible API Calls
- Supports GET, POST, PUT, DELETE
- Handles authentication tokens
- Parses JSON responses
- Extracts UUIDs from responses

## Performance Considerations

- **Total Runtime**: ~10-15 seconds (with local API + database)
- **Network Delay**: Add 5-10 seconds if testing remote API
- **Database Seeding**: ~1-2 seconds

## Next Steps

After successful integration tests:

1. **Deploy to Staging** - Use the same test script to validate staging environment
2. **Monitor Production** - Run tests periodically to ensure system health
3. **Extend Tests** - Add more edge cases and failure scenarios
4. **Load Testing** - Use wrk, Apache Bench, or similar tools for performance testing

## Support

For issues or questions:

1. Check [COMPLETE_COMPLIANCE_VERIFICATION.md](COMPLETE_COMPLIANCE_VERIFICATION.md) for requirement details
2. Review API logs: `docker-compose logs api`
3. Check database state: `make db-shell`
4. Review test script: `cat test_workflow.sh`

## Additional Resources

- [API Documentation](docs/SWAGGER_GUIDE.md)
- [RBAC Documentation](RBAC.md)
- [Business Rules](BUSINESS_RULES_VALIDATION.md)
- [Requirements Compliance](REQUIREMENTS_COMPLIANCE_REPORT.md)
