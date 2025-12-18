package auth
package auth

import (
	"testing"
	"time"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/config"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTProvider_GenerateToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:     "test-secret-key",
		Expiration: 24 * time.Hour,
	}
	provider := NewJWTProvider(cfg)

	user, _ := entity.NewUser("test@example.com", "password123", "John", "Doe", entity.UserRoleStudent)

	token, err := provider.GenerateToken(user)

	require.NoError(t, err)






























































































































}	return user	user, _ := entity.NewUser(email, "password123", "Test", "User", role)func mustCreateUser(email string, role entity.UserRole) *entity.User {}	}		})			assert.Equal(t, user.Role, claims.Role)			assert.Equal(t, user.Email, claims.Email)			assert.Equal(t, user.ID, claims.UserID)			// Verify claims			require.NotNil(t, claims)			require.NoError(t, err)			claims, err := provider.ValidateToken(token)			// Validate token			require.NotEmpty(t, token)			require.NoError(t, err)			token, err := provider.GenerateToken(user)			// Generate token		t.Run(string(user.Role), func(t *testing.T) {	for _, user := range users {	}		mustCreateUser("admin@example.com", entity.UserRoleAdmin),		mustCreateUser("instructor@example.com", entity.UserRoleInstructor),		mustCreateUser("student@example.com", entity.UserRoleStudent),	users := []*entity.User{	provider := NewJWTProvider(cfg)	}		Expiration: 24 * time.Hour,		Secret:     "test-secret-key",	cfg := &config.JWTConfig{func TestJWTProvider_RoundTrip(t *testing.T) {}	}		})			}				}					tt.check(t, claims)				if tt.check != nil {				require.NotNil(t, claims)				require.NoError(t, err)			} else {				assert.Error(t, err)			if tt.wantErr {			claims, err := provider.ValidateToken(token)			token := tt.setup()		t.Run(tt.name, func(t *testing.T) {	for _, tt := range tests {	}		},			wantErr: true,			},				return token				token, _ := expiredProvider.GenerateToken(user)				})					Expiration: -1 * time.Hour, // Already expired					Secret:     "test-secret-key",				expiredProvider := NewJWTProvider(&config.JWTConfig{			setup: func() string {			name: "expired token",		{		},			wantErr: true,			},				return token				token, _ := wrongProvider.GenerateToken(user)				})					Expiration: 24 * time.Hour,					Secret:     "wrong-secret",				wrongProvider := NewJWTProvider(&config.JWTConfig{			setup: func() string {			name: "token with wrong secret",		{		},			wantErr: true,			},				return ""			setup: func() string {			name: "empty token",		{		},			wantErr: true,			},				return "invalid.token.here"			setup: func() string {			name: "invalid token",		{		},			},				assert.Equal(t, user.Role, claims.Role)				assert.Equal(t, user.Email, claims.Email)				assert.Equal(t, user.ID, claims.UserID)			check: func(t *testing.T, claims *Claims) {			wantErr: false,			},				return token				token, _ := provider.GenerateToken(user)			setup: func() string {			name: "valid token",		{	}{		check   func(*testing.T, *Claims)		wantErr bool		setup   func() string		name    string	tests := []struct {	user, _ := entity.NewUser("test@example.com", "password123", "John", "Doe", entity.UserRoleStudent)	provider := NewJWTProvider(cfg)	}		Expiration: 24 * time.Hour,		Secret:     "test-secret-key",	cfg := &config.JWTConfig{func TestJWTProvider_ValidateToken(t *testing.T) {}	assert.NotEmpty(t, token)