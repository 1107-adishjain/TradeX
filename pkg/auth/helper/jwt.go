package helper

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrJWTSecretMissing = errors.New("JWT_SECRET is not set")

type CustomClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	Type  string `json:"type"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(email, role string) (string, error) {
	return generateToken(email, role, "access", 15*time.Minute)
}

func GenerateRefreshToken(email, role string) (string, error) {
	return generateToken(email, role, "refresh", 7*24*time.Hour)
}

func generateToken(email, role, tokenType string, ttl time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", ErrJWTSecretMissing
	}

	now := time.Now().UTC()
	claims := CustomClaims{
		Email: email,
		Role:  role,
		Type:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
