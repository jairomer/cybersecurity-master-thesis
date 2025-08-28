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

// Custom middleware to log XFCC header
func LogXFCC(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		xfcc := c.Request().Header.Get("X-Forwarded-Client-Cert")
		log.Println("XFCC Header:", xfcc)
		v := api.XFCC{Value: xfcc}
		ctx := context.WithValue(c.Request().Context(), "xfcc", &v)
		req := c.Request().WithContext(ctx)
		c.SetRequest(req)
		// Call the next handler
		return next(c)
	}
}

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

	e.Use(LogXFCC)
	e.Use(echojwt.WithConfig(jwtconfig))

	api.RegisterHandlers(e, api.NewStrictHandler(
		&server,
		// Middlewares
		[]api.StrictMiddlewareFunc{},
	))
	if err := e.Start("0.0.0.0:8000"); err != nil {
		log.Fatal(err)
	}
}
