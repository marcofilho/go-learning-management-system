package handler

import (
	"context"
	"net/http"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/middleware"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
)

// GetUserIDFromContext extracts the user ID from the request context
func GetUserIDFromContext(r *http.Request) (string, bool) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		return "", false
	}
	return claims.UserID, true
}

// GetUserFromContext extracts the full user claims from the request context
func GetUserFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(middleware.UserContextKey).(*auth.Claims)
	return claims, ok
}

// IsAdmin checks if the current user is an admin
func IsAdmin(r *http.Request) bool {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		return false
	}
	return claims.Role == entity.UserRoleAdmin
}

// IsInstructor checks if the current user is an instructor
func IsInstructor(r *http.Request) bool {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		return false
	}
	return claims.Role == entity.UserRoleInstructor || claims.Role == entity.UserRoleAdmin
}
