package api

import (
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/open-policy-agent/opa/v1/rego"
)

const (
	SecretKey = "supersecretkey"
)

type JWTClaims struct {
	jwt.RegisteredClaims
}

func populateAuthzData(ctx echo.Context, uri string, claims *JWTClaims) interface{} {
	log.Printf("Authorizing: %s accessing %s\n", claims.Audience[0], uri)
	return map[string]interface{}{
		"access_control": map[string]interface{}{
			"users": map[string]interface{}{
				"test1": map[string]interface{}{
					"acl": []string{"/hello/world"},
				},
				"test2": map[string]interface{}{
					"acl": []string{},
				},
			},
		},
		"jwt": map[string]interface{}{
			"aud": claims.Audience[0],
		},
		"uri": uri,
	}
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

func Authorize(ctx echo.Context) (bool, error) {
	uri := ctx.Request().URL.Path
	jwtStr := ctx.Request().Header.Get("Authorization")
	jwtStr = strings.ReplaceAll(jwtStr, "Bearer ", "")
	token, err := jwt.ParseWithClaims(jwtStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return false, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		data := populateAuthzData(ctx, uri, claims)

		query, err := rego.New(
			rego.Query("data.example.authz.allow"),
			rego.Load([]string{"./authz.rego"}, nil),
		).PrepareForEval(ctx.Request().Context())
		if err != nil {
			return false, err
		}

		results, err := query.Eval(ctx.Request().Context(), rego.EvalInput(data))
		if err != nil {
			return false, err
		}
		if len(results) > 0 && len(results[0].Expressions) > 0 {
			allowed, ok := results[0].Expressions[0].Value.(bool)
			if ok && allowed {
				log.Println("Authorized")
				return true, nil
			} else if !ok && allowed {
				log.Printf("Not Ok, but allowed")
			} else if ok && !allowed {
				log.Printf("Not allowed")
			}
		}
		return false, nil
	}
	return false, nil
}
