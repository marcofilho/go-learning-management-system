# Test Suite Documentation Index

## 🎯 Start Here

Choose based on your need:

### ⚡ I want to run tests RIGHT NOW
→ [QUICK_TEST_GUIDE.md](QUICK_TEST_GUIDE.md)

**Command:**
```bash
make start && sleep 10 && make test-workflow
```

---

### 📖 I want to understand the full test suite
→ [TEST_WORKFLOW_GUIDE.md](TEST_WORKFLOW_GUIDE.md)

**Covers:**
- 5 setup options
- Complete test flow
- Troubleshooting guide
- CI/CD integration

---

### 📋 I want setup details and features
→ [TEST_SETUP_SUMMARY.md](TEST_SETUP_SUMMARY.md)

**Covers:**
- What was created
- Test coverage breakdown
- Script features
- Next steps

---

## 📁 Files Created

| File | Size | Purpose |
|------|------|---------|
| [test_workflow.sh](test_workflow.sh) | 19KB | Main executable test script (556 lines) |
| [TEST_WORKFLOW_GUIDE.md](TEST_WORKFLOW_GUIDE.md) | 15KB | Comprehensive guide (400+ lines) |
| [QUICK_TEST_GUIDE.md](QUICK_TEST_GUIDE.md) | 8KB | Quick reference |
| [TEST_SETUP_SUMMARY.md](TEST_SETUP_SUMMARY.md) | 12KB | Setup documentation |
| Makefile | Updated | Added `test-workflow` and `test-integration` targets |

---

## 🧪 Test Summary

**Total Tests: 27** across **8 phases**

```
Phase 1: Authentication (6)      → Registration & Login
Phase 2: Courses (4)             → CRUD operations
Phase 3: Modules (2)             → Create & List
Phase 4: Lessons (4)             → Versioning
Phase 5: Enrollment (3)          → Self-enroll & List
Phase 6: Webhook (2)             → Certification
Phase 7: Audit Logs (3)          → Logging & Access
Phase 8: Authorization (3)       → Permission Tests
```

---

## 🚀 Quick Commands

### Execute Tests
```bash
# One command - recommended
make start && sleep 10 && make test-workflow

# Just run tests (if services already running)
make test-workflow

# Full end-to-end
make test-integration
```

### Alternative Run Methods
```bash
# Direct execution
./test_workflow.sh

# With environment variables
API_URL=http://localhost:8080/api ./test_workflow.sh

# With verbose output
bash -x test_workflow.sh
```

### Supporting Commands
```bash
make test          # 243 unit tests
make test-all      # Include DB adapters
make docker-up     # Start services
make docker-down   # Stop services
make db-reset      # Reset database
```

---

## 📊 What Gets Tested

### Authentication
- ✅ User registration (all roles)
- ✅ User login
- ✅ JWT token generation
- ✅ Token validation

### Course Management
- ✅ Create courses
- ✅ List courses
- ✅ Get course details
- ✅ Update courses
- ✅ Ownership verification
- ✅ Soft deletes

### Content Management
- ✅ Module creation
- ✅ Lesson versioning
- ✅ Auto-incrementing versions
- ✅ Version listing

### Enrollment
- ✅ Self-enrollment
- ✅ Duplicate prevention
- ✅ Student course listing
- ✅ Instructor student listing

### Webhooks
- ✅ Certification processing
- ✅ Status updates
- ✅ Business rule validation

### Audit Trail
- ✅ Action logging
- ✅ Before/after states
- ✅ User attribution
- ✅ Access control

### Authorization
- ✅ Role-based access
- ✅ Ownership verification
- ✅ Permission boundaries
- ✅ 403 Forbidden responses

---

## 🔍 Example Test Flow

```bash
$ make start && sleep 10 && make test-workflow

# Output shows:
✓ ========================================
✓ LMS API Integration Test Suite
✓ ========================================

Testing API at: http://localhost:8080/api
Database: postgres@localhost:5432/lms_db

✓ Pre-Flight Checks
✓ API is running
✓ Database is accessible

✓ 1. Authentication Tests
✓ Admin registered successfully
✓ Instructor registered successfully
✓ Student registered successfully
✓ Admin login successful
✓ Instructor login successful
✓ Student login successful

✓ 2. Course Management Tests
✓ Course created successfully
✓ Courses listed successfully
✓ Course retrieved successfully
✓ Course updated successfully

... (more phases) ...

✓ Test Summary
Total Tests: 27
✓ Passed: 27
✗ Failed: 0

✓ All tests passed!
```

---

## 🛠️ Troubleshooting

### API Not Running
```bash
curl http://localhost:8080/api/health
make run  # Start API
```

### Database Issues
```bash
docker-compose ps
make db-reset  # Nuclear option
```

### Script Problems
```bash
chmod +x test_workflow.sh
bash test_workflow.sh
```

---

## 📚 Related Documentation

For understanding the system:

- [REQUIREMENTS_EVALUATION.md](REQUIREMENTS_EVALUATION.md) - Requirements overview
- [COMPLETE_COMPLIANCE_VERIFICATION.md](COMPLETE_COMPLIANCE_VERIFICATION.md) - Detailed compliance
- [RBAC.md](RBAC.md) - Authorization matrix
- [BUSINESS_RULES_VALIDATION.md](BUSINESS_RULES_VALIDATION.md) - Business logic
- [AUDIT_LOGGING.md](AUDIT_LOGGING.md) - Audit trail details

---

## ✨ Features

**Pre-Flight Checks**
- ✓ API health check
- ✓ Database connectivity
- ✓ Automatic error exit

**Test Management**
- ✓ Sequential execution
- ✓ Token extraction
- ✓ JSON parsing
- ✓ UUID extraction

**Reporting**
- ✓ Color-coded output
- ✓ Test statistics
- ✓ Pass/fail summary
- ✓ Exit codes

**Error Handling**
- ✓ Graceful failures
- ✓ Detailed error messages
- ✓ Validation feedback

---

## 🎓 Learning Path

1. **Get Started** → QUICK_TEST_GUIDE.md
2. **Run Tests** → `make start && sleep 10 && make test-workflow`
3. **Review Results** → Check output
4. **Understand System** → REQUIREMENTS_EVALUATION.md
5. **Deep Dive** → COMPLETE_COMPLIANCE_VERIFICATION.md
6. **Authorization** → RBAC.md

---

## 📞 Support

### Common Issues

**Tests fail immediately?**
- Check API is running: `curl http://localhost:8080/api/health`
- Check database: `docker-compose ps`
- View logs: `docker-compose logs`

**Tests run slow?**
- Check API logs: `docker-compose logs api`
- Check DB logs: `docker-compose logs postgres`
- Restart services: `make docker-restart`

**Script errors?**
- Make executable: `chmod +x test_workflow.sh`
- Try with explicit shell: `bash test_workflow.sh`
- Check permissions: `ls -la test_workflow.sh`

---

## 🚀 Next Steps

After tests pass:

1. ✅ Verify all 27 tests pass
2. 📊 Review test output
3. 🔍 Check API logs for warnings
4. 📚 Review COMPLETE_COMPLIANCE_VERIFICATION.md
5. 🚀 Deploy to staging
6. 📈 Set up CI/CD pipeline
7. 🔄 Run tests periodically

---

## 📋 Checklist

Before running production:

- [ ] All 27 tests passing
- [ ] No errors in API logs
- [ ] No errors in database logs
- [ ] Requirements verified ([COMPLETE_COMPLIANCE_VERIFICATION.md](COMPLETE_COMPLIANCE_VERIFICATION.md))
- [ ] Authorization tested ([RBAC.md](RBAC.md))
- [ ] Business rules validated ([BUSINESS_RULES_VALIDATION.md](BUSINESS_RULES_VALIDATION.md))
- [ ] Audit logging working ([AUDIT_LOGGING.md](AUDIT_LOGGING.md))

---

## 🎯 Quick Access

| Need | Link | Command |
|------|------|---------|
| Quick start | [QUICK_TEST_GUIDE.md](QUICK_TEST_GUIDE.md) | `make test-workflow` |
| Full guide | [TEST_WORKFLOW_GUIDE.md](TEST_WORKFLOW_GUIDE.md) | Read docs |
| Requirements | [REQUIREMENTS_EVALUATION.md](REQUIREMENTS_EVALUATION.md) | Review |
| Compliance | [COMPLETE_COMPLIANCE_VERIFICATION.md](COMPLETE_COMPLIANCE_VERIFICATION.md) | Review |
| Auth rules | [RBAC.md](RBAC.md) | Review |
| Business rules | [BUSINESS_RULES_VALIDATION.md](BUSINESS_RULES_VALIDATION.md) | Review |

---

**Status: ✅ Ready to Test**

Run this now:
```bash
make start && sleep 10 && make test-workflow
```
