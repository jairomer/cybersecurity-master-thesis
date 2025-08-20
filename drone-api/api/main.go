package main

import (
	"context"
	"log"
	"strings"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"uc3m.es/drone-api/pkg/api"
)

func main() {
	server := api.NewServer()
	e := echo.New()
	e.Use(middleware.Logger())

	jwtconfig := echojwt.Config{
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == "/login" {
				log.Println("Accessing /login endpoint.")
				return true
			}
			return false
		},
		SuccessHandler: func(c echo.Context) {
			log.Printf("JWT Authentication successful on: %s\n", c.Request().URL)
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, "uri", c.Request().URL.Path)
			ctx = context.WithValue(ctx, "jwt", strings.ReplaceAll(c.Request().Header.Get("Authorization"), "Bearer ", ""))
			c.SetRequest(c.Request().WithContext(ctx))
		},
		SigningKey: []byte(api.SecretKey),
	}

	e.Use(echojwt.WithConfig(jwtconfig))

	api.RegisterHandlers(e, api.NewStrictHandler(
		&server,
		// Middlewares
		[]api.StrictMiddlewareFunc{},
	))
	if err := e.Start("127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}
