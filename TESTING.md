# Testing Guide

Complete guide for running and writing tests for the Learning Management System.

---

## Quick Start

### Run All Tests
```bash
go test ./... -v
```

### Run Tests with Coverage
```bash
go test ./... -cover
```

### Run Specific Package Tests
```bash
# Entity tests
go test ./src/internal/domain/entity -v

# Use case tests
go test ./src/usecase -v

# Handler tests
go test ./src/internal/adapter/http/handler -v

# Middleware tests
go test ./src/internal/adapter/http/middleware -v
```

### Run Specific Test
```bash
go test ./src/internal/domain/entity -run TestNewUser -v
```

---

## Test Status

**Current Status**: ✅ **243 tests passing**

- ✅ Unit tests for all entities
- ✅ Handler tests for all endpoints
- ✅ Use case tests for business logic
- ✅ Middleware tests for authorization
- ✅ Test coverage: 76%

---

## Test Structure

### 1. Domain Entity Tests

**Location**: `src/internal/domain/entity/*_test.go`

**Coverage**:
- User entity (creation, validation, roles, authentication)
- Course entity (creation, validation, ownership, permissions)
- Enrollment entity (creation, status transitions, validation)
- Module entity (creation, validation, ordering)
- Lesson entity (versioning, validation)
- Audit log entity (logging, validation)

**Example**:
```go
func TestNewUser(t *testing.T) {
    user, err := entity.NewUser("test@example.com", "password123", "Test", "User", entity.UserRoleStudent)
    require.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
}
```

---

### 2. Use Case Tests

**Location**: `src/usecase/*_test.go`

**Coverage**:
- User use case (registration, login, CRUD)
- Course use case (creation, updates, ownership)
- Enrollment use case (enrollment, status updates)
- Module use case (creation, ordering)
- Lesson use case (versioning, content management)
- Certification webhook use case

**Mocking**: Uses `testify/mock` for repository dependencies

**Example**:
```go
func TestUserUseCase_Register(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockTokenProvider := new(MockTokenProvider)
    useCase := NewUserUseCase(mockRepo, mockTokenProvider)
    
    // Setup mocks
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    // Test
    user, err := useCase.Register(ctx, ...)
    assert.NoError(t, err)
}
```

---

### 3. HTTP Handler Tests

**Location**: `src/internal/adapter/http/handler/*_test.go`

**Coverage**:
- All API endpoints
- Request/response validation
- HTTP status codes
- Error handling
- Authorization checks

**Example**:
```go
func TestUserHandler_Register(t *testing.T) {
    mockUseCase := new(MockUserUseCase)
    handler := NewUserHandler(mockUseCase)
    
    req := httptest.NewRequest("POST", "/api/auth/register", ...)
    w := httptest.NewRecorder()
    
    handler.Register(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

---

### 4. Middleware Tests

**Location**: `src/internal/adapter/http/middleware/*_test.go`

**Coverage**:
- Authentication middleware
- Authorization middleware (roles, ownership, enrollment)
- Error handling
- Context propagation

---

## Test Framework

**Framework**: Go's built-in `testing` package  
**Assertions**: `testify/assert` and `testify/require`  
**Mocking**: `testify/mock`

### Dependencies

Already installed in `go.mod`:
- `github.com/stretchr/testify` ✅
- `github.com/google/uuid` ✅

---

## Test Template

```go
package entity

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestYourFunction(t *testing.T) {
    // Arrange
    input := "test data"
    
    // Act
    result, err := YourFunction(input)
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

---

## Best Practices

1. **Run tests frequently** during development
2. **Use `-v` flag** for verbose output
3. **Check coverage** with `-cover` flag
4. **Test both success and error cases**
5. **Keep tests simple and focused**
6. **Use table-driven tests** for multiple scenarios
7. **Mock external dependencies** (databases, APIs)

---

## Advanced Testing

### Race Condition Detection
```bash
go test -race ./...
```

### Coverage Report
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Fail Fast
```bash
go test -failfast ./...
```

### Verbose with Specific Test
```bash
go test -v -run TestSpecificTest ./...
```

---

## Test Coverage Goals

- **Target**: 80%+ coverage
- **Current**: 76%
- **Focus Areas**: Edge cases, error handling, authorization paths

---

## Troubleshooting

### Tests Failing
1. Check test output for specific error messages
2. Run with `-v` flag for detailed output
3. Check database/test data setup
4. Verify mock expectations

### Slow Tests
1. Use in-memory databases for unit tests
2. Parallelize tests with `t.Parallel()` where safe
3. Use test fixtures instead of creating data in each test

---

**Status**: Comprehensive test suite ✅  
**Coverage**: 76%  
**Tests Passing**: 243/243 ✅

