package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/middleware"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
)

func TestGetUserIDFromContext(t *testing.T) {
	t.Run("should return user id when user is in context", func(t *testing.T) {
		userID := uuid.New()
		ctx := context.WithValue(context.Background(), middleware.UserContextKey, &auth.Claims{UserID: userID})
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

		retrievedUserID, ok := GetUserIDFromContext(req)
		assert.True(t, ok)
		assert.Equal(t, userID, retrievedUserID)
	})

	t.Run("should return empty string when user is not in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		retrievedUserID, ok := GetUserIDFromContext(req)
		assert.False(t, ok)
		assert.Equal(t, uuid.Nil, retrievedUserID)
	})
}

func TestGetUserFromContext(t *testing.T) {
	t.Run("should return user when user is in context", func(t *testing.T) {
		claims := &auth.Claims{UserID: uuid.New()}
		ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)

		retrievedClaims, ok := GetUserFromContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, claims, retrievedClaims)
	})

	t.Run("should return nil when user is not in context", func(t *testing.T) {
		ctx := context.Background()

		retrievedClaims, ok := GetUserFromContext(ctx)
		assert.False(t, ok)
		assert.Nil(t, retrievedClaims)
	})
}

func TestIsAdmin(t *testing.T) {
	t.Run("should return true when user is admin", func(t *testing.T) {
		claims := &auth.Claims{Role: "admin"}
		ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

		isAdmin := IsAdmin(req)
		assert.True(t, isAdmin)
	})

	t.Run("should return false when user is not admin", func(t *testing.T) {
		claims := &auth.Claims{Role: "student"}
		ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

		isAdmin := IsAdmin(req)
		assert.False(t, isAdmin)
	})

	t.Run("should return false when user is not in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		isAdmin := IsAdmin(req)
		assert.False(t, isAdmin)
	})
}

func TestIsInstructor(t *testing.T) {
	t.Run("should return true when user is instructor", func(t *testing.T) {
		claims := &auth.Claims{Role: "instructor"}
		ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

		isInstructor := IsInstructor(req)
		assert.True(t, isInstructor)
	})

	t.Run("should return false when user is not instructor", func(t *testing.T) {
		claims := &auth.Claims{Role: "student"}
		ctx := context.WithValue(context.Background(), middleware.UserContextKey, claims)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

		isInstructor := IsInstructor(req)
		assert.False(t, isInstructor)
	})

	t.Run("should return false when user is not in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		isInstructor := IsInstructor(req)
		assert.False(t, isInstructor)
	})
}
