package api

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	SecretKey = "supersecretkey"
)

//type RegoEngine struct {
//	engine *rego.Rego
//}

type JWTClaims struct {
	jwt.RegisteredClaims
}

func GenerateToken(username string) (string, error) {
	claims := JWTClaims{
		// We will keep claims encrypted
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "hello/world/server",
			Audience: jwt.ClaimStrings{
				"hello/world/" + username,
			},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}

//func VerifyToken(ctx context.Context, engine *RegoEngine, tokenString string) (bool, error) {
//	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return []byte(SecretKey), nil
//	})
//	if err != nil {
//		return false, err
//	}
//	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
//		// todo: use claims for something.
//		_ = claims
//
//		return true, nil
//	}
//	return false, nil
//}
