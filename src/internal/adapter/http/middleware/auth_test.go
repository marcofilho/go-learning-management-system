package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTokenProvider struct {
	mock.Mock
}

func (m *MockTokenProvider) GenerateToken(user *entity.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockTokenProvider) ValidateToken(tokenString string) (*auth.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Claims), args.Error(1)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	mockProvider := new(MockTokenProvider)
	claims := &auth.Claims{
		UserID: uuid.New(),
		Email:  "test@example.com",
		Role:   entity.UserRoleStudent,
	}
	mockProvider.On("ValidateToken", "valid-token").Return(claims, nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, ok := r.Context().Value(UserContextKey).(*auth.Claims)
		assert.True(t, ok)
		assert.NotNil(t, c)
		w.WriteHeader(http.StatusOK)
	})

	middleware := AuthMiddleware(mockProvider)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockProvider.AssertExpectations(t)
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	mockProvider := new(MockTokenProvider)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := AuthMiddleware(mockProvider)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestRequireRole_Admin(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequireRole(entity.UserRoleAdmin)
	wrappedHandler := middleware(handler)

	claims := &auth.Claims{
		UserID: uuid.New(),
		Role:   entity.UserRoleAdmin,
	}
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, claims)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRequireRole_Forbidden(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequireRole(entity.UserRoleAdmin)
	wrappedHandler := middleware(handler)

	claims := &auth.Claims{
		UserID: uuid.New(),
		Role:   entity.UserRoleStudent,
	}
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, claims)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)
}
