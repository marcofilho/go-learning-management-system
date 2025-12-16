package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/dto"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(tokenProvider auth.TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, entity.ErrUnauthorized, "Missing authorization header")
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondWithError(w, http.StatusUnauthorized, entity.ErrUnauthorized, "Invalid authorization header format")
				return
			}
			claims, err := tokenProvider.ValidateToken(parts[1])
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, entity.ErrUnauthorized, "Invalid token")
				return
			}
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...entity.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
			if !ok {
				respondWithError(w, http.StatusUnauthorized, entity.ErrUnauthorized, "User not authenticated")
				return
			}
			hasRole := false
			for _, role := range roles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}
			if !hasRole {
				respondWithError(w, http.StatusForbidden, entity.ErrInsufficientPermissions, "Insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func respondWithError(w http.ResponseWriter, code int, err error, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(dto.NewErrorResponse(err, message))
}
