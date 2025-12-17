package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
)

// OptionalAuthMiddleware attempts to authenticate the user if a token is present,
// but allows the request to proceed even if no token is provided.
// This is useful for endpoints that are public but have enhanced functionality for authenticated users.
func OptionalAuthMiddleware(tokenProvider auth.TokenProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			
			// If no auth header, proceed without authentication
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			// If auth header exists, try to validate it
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				// Invalid format, proceed without authentication (don't block)
				next.ServeHTTP(w, r)
				return
			}

			// Try to validate token
			claims, err := tokenProvider.ValidateToken(parts[1])
			if err != nil {
				// Invalid token, proceed without authentication (don't block)
				next.ServeHTTP(w, r)
				return
			}

			// Valid token, add claims to context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
