package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	UpdatedAt int64  `json:"updated_at"`
	Type      string `json:"type"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId, role string, updatedAt int64) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expired := os.Getenv("JWT_EXPIRED")
	if expired == "" {
		expired = "15m"
	}

	duration, err := time.ParseDuration(expired)
	if err != nil {
		return "", errors.New("Gagal parse JWT_EXPIRED")
	}

	claims := &Claims{
		UserID:    userId,
		Role:      role,
		UpdatedAt: updatedAt,
		Type:      "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID, role string, updatedAt int64, rememberMe bool) (string, error) {
	secret := os.Getenv("REFRESH_TOKEN_SECRET")

	duration := 24 * time.Hour
	if rememberMe {
		duration = 7 * 24 * time.Hour
	}

	claims := &Claims{
		UserID:    userID,
		Role:      role,
		UpdatedAt: updatedAt,
		Type:      "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateEmailVerificationToken(userID string) (string, error) {
	secret := os.Getenv("EMAIL_VERIFY_SECRET")

	claims := &Claims{
		UserID: userID,
		Type:   "email_verification",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateResetPasswordToken(userID string, updatedAt int64) (string, error) {
	secret := os.Getenv("RESET_PASSWORD_SECRET")

	claims := &Claims{
		UserID:    userID,
		UpdatedAt: updatedAt,
		Type:      "reset_password",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokens, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokens, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
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
