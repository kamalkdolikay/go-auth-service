package handlers

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"auth/config"
)

var (
	jwtSecret       []byte
	jwtExpiresMinutes int
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// InitJWT loads config once
func InitJWT() {
	secret := config.GetEnv("JWT_SECRET", "fallback-secret-insecure")
	if len(secret) < 32 {
		panic("JWT_SECRET must be at least 32 characters")
	}
	jwtSecret = []byte(secret)

	hours, _ := strconv.Atoi(config.GetEnv("JWT_EXPIRES_MINUTES", "15"))
	jwtExpiresMinutes = hours
}

// GenerateJWT creates signed token
func GenerateJWT(userID int, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(jwtExpiresMinutes))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseJWT verifies token from cookie
func ParseJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}