package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/config"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
)

type TokenProvider interface {
	GenerateToken(user *entity.User) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type Claims struct {
	UserID string          `json:"user_id"`
	Email  string          `json:"email"`
	Role   entity.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type JWTProvider struct {
	secret     []byte
	expiration time.Duration
}

func NewJWTProvider(cfg *config.JWTConfig) TokenProvider {
	return &JWTProvider{
		secret:     []byte(cfg.Secret),
		expiration: cfg.Expiration,
	}
}

func (p *JWTProvider) GenerateToken(user *entity.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}

func (p *JWTProvider) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return p.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
