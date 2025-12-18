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
	assert.NotEmpty(t, token)
}

func TestJWTProvider_ValidateToken_Success(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:     "test-secret-key",
		Expiration: 24 * time.Hour,
	}
	provider := NewJWTProvider(cfg)

	user, _ := entity.NewUser("test@example.com", "password123", "John", "Doe", entity.UserRoleStudent)
	token, _ := provider.GenerateToken(user)

	claims, err := provider.ValidateToken(token)

	require.NoError(t, err)
	require.NotNil(t, claims)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Role, claims.Role)
}

func TestJWTProvider_ValidateToken_InvalidToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:     "test-secret-key",
		Expiration: 24 * time.Hour,
	}
	provider := NewJWTProvider(cfg)

	_, err := provider.ValidateToken("invalid.token.here")
	assert.Error(t, err)
}

func TestJWTProvider_ValidateToken_WrongSecret(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:     "test-secret-key",
		Expiration: 24 * time.Hour,
	}
	provider := NewJWTProvider(cfg)

	wrongProvider := NewJWTProvider(&config.JWTConfig{
		Secret:     "wrong-secret",
		Expiration: 24 * time.Hour,
	})

	user, _ := entity.NewUser("test@example.com", "password123", "John", "Doe", entity.UserRoleStudent)
	token, _ := wrongProvider.GenerateToken(user)

	_, err := provider.ValidateToken(token)
	assert.Error(t, err)
}
