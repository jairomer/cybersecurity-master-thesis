package api

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	SecretKey = "supersecretkey"
)

type JWTClaims struct {
	jwt.RegisteredClaims
}

func GenerateToken(username string) (string, error) {
	claims := JWTClaims{
		// We will keep claims encrypted
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "hello/world/server",
			Audience:  jwt.ClaimStrings{username},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}
