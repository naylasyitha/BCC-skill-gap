package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userId, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expired := os.Getenv("JWT_EXPIRED")

	duration, err := time.ParseDuration(expired)
	if err != nil {
		return "", errors.New("gagal parse JWT_EXPIRED")
	}

	claims := &Claims{
		UserID: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokens string) (*Claims, error) {
	secret := os.Getenv("JWt_SECRET")

	token, err := jwt.ParseWithClaims(tokens, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		err := t.Method.(*jwt.SigningMethodHMAC)
		if err != nil {
			return nil, errors.New("signing method tidak valid")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, errors.New("Token tidak valid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("Token tidak valid")
	}

	return claims, nil
}
