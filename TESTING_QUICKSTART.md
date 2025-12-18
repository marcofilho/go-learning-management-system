# Unit Testing - Quick Start Guide

## 🚀 Running Tests

### All Tests
```bash
go test ./... -v
```

### Entity Tests (Currently Implemented)
```bash
go test ./src/internal/domain/entity -v
```

### With Coverage
```bash
go test ./... -cover
```

## ✅ Current Test Status

### Working Tests
- **User Entity**: 3 tests passing ✅
  - TestNewUser
  - TestUser_ValidatePassword  
  - TestUser_Roles

### Test Output
```
=== RUN   TestNewUser
--- PASS: TestNewUser (0.07s)
=== RUN   TestUser_ValidatePassword
--- PASS: TestUser_ValidatePassword (0.14s)
=== RUN   TestUser_Roles
--- PASS: TestUser_Roles (0.13s)
PASS
ok  ...entity  0.618s
```

## 📚 Documentation

1. **[UNIT_TESTING_GUIDE.md](./UNIT_TESTING_GUIDE.md)** - Comprehensive testing guide
2. **[TESTING_IMPLEMENTATION_SUMMARY.md](./TESTING_IMPLEMENTATION_SUMMARY.md)** - Implementation details

## 🎯 What's Next

To continue implementing tests:

### 1. Entity Tests
```bash
# Create test file
touch src/internal/domain/entity/course_test.go

# Run tests
go test ./src/internal/domain/entity -v
```

### 2. Use Case Tests
```bash
# Create mocks and tests
touch src/usecase/user_usecase_test.go

# Run tests
go test ./src/usecase -v
```

### 3. Handler Tests
```bash
# Create HTTP tests
touch src/internal/adapter/http/handler/user_handler_test.go

# Run tests
go test ./src/internal/adapter/http/handler -v
```

## 🛠️ Test Template

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

## 📦 Dependencies

Already installed:
- `github.com/stretchr/testify` ✅
- `github.com/google/uuid` ✅

## 💡 Tips

1. **Run tests frequently** during development
2. **Use `-v` flag** for verbose output
3. **Check coverage** with `-cover` flag
4. **Test both success and error cases**
5. **Keep tests simple and focused**

## 🔍 Debugging Failed Tests

```bash
# Run specific test
go test ./src/internal/domain/entity -run TestNewUser -v

# Show more detail
go test -v -failfast ./...

# Run with race detector
go test -race ./...
```

## 🎓 Learn More

- Read [UNIT_TESTING_GUIDE.md](./UNIT_TESTING_GUIDE.md) for comprehensive guide
- Check [Go Testing Docs](https://pkg.go.dev/testing)
- Review [Testify Documentation](https://github.com/stretchr/testify)

---

**Status**: Foundation Complete ✅  
**Coverage**: ~10% (Goal: 70%+)  
**Next**: Implement remaining entity tests
