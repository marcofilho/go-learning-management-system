# Quick Test Execution Guide

## Fastest Way to Test Everything

### 1. One-Command Full Setup & Test

```bash
make start && sleep 10 && make test-workflow
```

This will:
- ✅ Start Docker (PostgreSQL + API)
- ✅ Run migrations automatically
- ✅ Wait for services to be ready
- ✅ Execute full integration test suite

### 2. If Services Already Running

```bash
make test-workflow
```

This will:
- ✅ Skip the Docker startup
- ✅ Run all 26 integration tests
- ✅ Display results with color coding

---

## Test Categories

### Run Unit Tests Only

```bash
make test
```

**Runs**: 243 unit tests (excludes DB repository adapters)

### Run All Tests (Including DB)

```bash
make test-all
```

**Runs**: Full test suite including PostgreSQL adapter tests

### Run Integration Workflow

```bash
make test-workflow
```

**Runs**: Full sequential workflow (8 phases, 26+ tests)

### Run Full End-to-End

```bash
make test-integration
```

**Runs**: Ensures Docker is up, then runs workflow tests

---

## Troubleshooting Quick Fixes

### API Not Responding

```bash
# Check if running
make docker-logs

# Restart
make docker-restart
```

### Database Connection Issues

```bash
# Reset database (WARNING: clears all data)
make db-reset
```

### Test Script Fails to Run

```bash
# Make executable
chmod +x test_workflow.sh

# Run directly
bash test_workflow.sh
```

### Debug Mode

```bash
# Run with verbose output
bash -x test_workflow.sh
```

---

## What Gets Tested

| Phase | Tests | Validates |
|-------|-------|-----------|
| Authentication | 6 | Registration, login for all roles |
| Courses | 4 | CRUD operations with ownership |
| Modules | 2 | Create and list |
| Lessons | 4 | Versioning, auto-increment |
| Enrollment | 3 | Self-enroll, listing |
| Webhook | 2 | Certification processing |
| Audit Logs | 3 | Logging and access control |
| Authorization | 3 | Permission boundaries |
| **Total** | **27** | **Complete workflow** |

---

## Expected Output Summary

```
✓ All tests passed!

Total Tests: 27
Passed: 27
Failed: 0
```

---

## Test Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Integration Test Suite                   │
└─────────────────────────────────────────────────────────────┘
                             ↓
                    ┌─────────────────┐
                    │ Pre-Flight Check│
                    │ API + Database  │
                    └────────┬────────┘
                             ↓
┌─────────────────────────────────────────────────────────────┐
│ 1. Authentication          → Register & Login (6 tests)     │
│ 2. Course Management       → Create, List, Update (4 tests) │
│ 3. Module Management       → Create, List (2 tests)         │
│ 4. Lesson Management       → Versioning (4 tests)           │
│ 5. Enrollment              → Enroll & List (3 tests)        │
│ 6. Webhook                 → Certification (2 tests)        │
│ 7. Audit Logging           → Verify Logs (3 tests)          │
│ 8. Authorization           → Permission Tests (3 tests)     │
└─────────────────────────────────────────────────────────────┘
                             ↓
                    ┌─────────────────┐
                    │ Test Summary    │
                    │ Report Results  │
                    └─────────────────┘
```

---

## Environment Variables

```bash
# API Configuration
API_URL=http://localhost:8080/api

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=lms_db

# Run with custom values
DB_HOST=myhost.com DB_PORT=5433 ./test_workflow.sh
```

---

## Manual Testing with cURL

If you want to test specific endpoints:

```bash
# Get health
curl http://localhost:8080/api/health

# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email":"test@example.com",
    "password":"Pass123",
    "first_name":"Test",
    "last_name":"User",
    "role":"student"
  }'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email":"test@example.com",
    "password":"Pass123"
  }'
```

---

## Next Steps

✅ **All tests passing?**

1. Review [REQUIREMENTS_EVALUATION.md](REQUIREMENTS_EVALUATION.md) - Requirement compliance
2. Review [COMPLETE_COMPLIANCE_VERIFICATION.md](COMPLETE_COMPLIANCE_VERIFICATION.md) - Detailed verification
3. Deploy to production or staging
4. Set up continuous integration
5. Monitor with periodic test runs

---

## Documentation

- 📚 [TEST_WORKFLOW_GUIDE.md](TEST_WORKFLOW_GUIDE.md) - Full guide
- 📋 [REQUIREMENTS_EVALUATION.md](REQUIREMENTS_EVALUATION.md) - Requirements overview
- ✅ [COMPLETE_COMPLIANCE_VERIFICATION.md](COMPLETE_COMPLIANCE_VERIFICATION.md) - Detailed verification
- 🔐 [RBAC.md](RBAC.md) - Authorization matrix
- 📝 [BUSINESS_RULES_VALIDATION.md](BUSINESS_RULES_VALIDATION.md) - Business logic

---

## Commands Summary

```bash
# Setup & Start
make start                  # Start everything
make docker-up              # Start services only
make docker-down            # Stop services

# Testing
make test                   # Unit tests
make test-all               # All tests
make test-workflow          # Integration workflow
make test-integration       # Full E2E

# Development
make run                    # Run API locally
make dev                    # Dev mode with hot reload
make build                  # Build binary

# Database
make migrate-up             # Apply migrations
make db-reset               # Reset database
make db-shell               # Connect to DB
```

---

**Status: ✅ Ready to Test!**
