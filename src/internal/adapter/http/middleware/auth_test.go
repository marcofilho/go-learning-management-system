package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTokenProvider for testing
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

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		setupMock      func(*MockTokenProvider)
		expectedStatus int
		checkContext   bool
	}{
		{
			name:       "valid token",
			authHeader: "Bearer valid-token",
			setupMock: func(m *MockTokenProvider) {
				claims := &auth.Claims{
					UserID: "user-123",
					Email:  "test@example.com",
					Role:   entity.UserRoleStudent,
				}
				m.On("ValidateToken", "valid-token").Return(claims, nil)
			},
			expectedStatus: http.StatusOK,
			checkContext:   true,
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			setupMock:      func(m *MockTokenProvider) {},
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
		{
			name:           "invalid header format - no Bearer",
			authHeader:     "invalid-token",
			setupMock:      func(m *MockTokenProvider) {},
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
		{
			name:           "invalid header format - only Bearer",
			authHeader:     "Bearer",
			setupMock:      func(m *MockTokenProvider) {},
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
		{
			name:       "invalid token",
			authHeader: "Bearer invalid-token",
			setupMock: func(m *MockTokenProvider) {
				m.On("ValidateToken", "invalid-token").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
			checkContext:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := new(MockTokenProvider)
			tt.setupMock(mockProvider)

			// Create test handler
			var contextChecked bool
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkContext {
					claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
					assert.True(t, ok, "Claims should be in context")
					assert.NotNil(t, claims, "Claims should not be nil")
					contextChecked = true
				}
				w.WriteHeader(http.StatusOK)
			})

			// Apply middleware
			middleware := AuthMiddleware(mockProvider)
			wrappedHandler := middleware(handler)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			// Execute
			wrappedHandler.ServeHTTP(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.checkContext {
				assert.True(t, contextChecked, "Context should have been checked")
			}
			mockProvider.AssertExpectations(t)
		})
	}
}

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name           string
		requiredRoles  []entity.UserRole
		userRole       entity.UserRole
		hasContext     bool
		expectedStatus int
	}{
		{
			name:           "admin accessing admin route",
			requiredRoles:  []entity.UserRole{entity.UserRoleAdmin},
			userRole:       entity.UserRoleAdmin,
			hasContext:     true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "instructor accessing instructor route",
			requiredRoles:  []entity.UserRole{entity.UserRoleInstructor},
			userRole:       entity.UserRoleInstructor,
			hasContext:     true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "student accessing student route",
			requiredRoles:  []entity.UserRole{entity.UserRoleStudent},
			userRole:       entity.UserRoleStudent,
			hasContext:     true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "admin accessing instructor or admin route",
			requiredRoles:  []entity.UserRole{entity.UserRoleInstructor, entity.UserRoleAdmin},
			userRole:       entity.UserRoleAdmin,
			hasContext:     true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "instructor accessing instructor or admin route",
			requiredRoles:  []entity.UserRole{entity.UserRoleInstructor, entity.UserRoleAdmin},
			userRole:       entity.UserRoleInstructor,
			hasContext:     true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "student accessing admin route - forbidden",
			requiredRoles:  []entity.UserRole{entity.UserRoleAdmin},
			userRole:       entity.UserRoleStudent,
			hasContext:     true,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "student accessing instructor route - forbidden",
			requiredRoles:  []entity.UserRole{entity.UserRoleInstructor},
			userRole:       entity.UserRoleStudent,
			hasContext:     true,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "no context - unauthorized",
			requiredRoles:  []entity.UserRole{entity.UserRoleAdmin},
			userRole:       entity.UserRoleAdmin,
			hasContext:     false,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Apply middleware
			middleware := RequireRole(tt.requiredRoles...)
			wrappedHandler := middleware(handler)

			// Create request with or without context
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.hasContext {
				claims := &auth.Claims{
					UserID: "user-123",
					Email:  "test@example.com",
					Role:   tt.userRole,
				}
				ctx := context.WithValue(req.Context(), UserContextKey, claims)
				req = req.WithContext(ctx)
			}
			rr := httptest.NewRecorder()

			// Execute
			wrappedHandler.ServeHTTP(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestRequireRole_MultipleRoles(t *testing.T) {
	roles := []entity.UserRole{entity.UserRoleAdmin, entity.UserRoleInstructor, entity.UserRoleStudent}

	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware := RequireRole(roles...)
			wrappedHandler := middleware(handler)

			claims := &auth.Claims{
				UserID: "user-123",
				Email:  "test@example.com",
				Role:   role,
			}
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := context.WithValue(req.Context(), UserContextKey, claims)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
		})
	}
}
