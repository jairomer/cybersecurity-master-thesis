package main

import (
	"log"
	"net/http"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"uc3m.es/hello/server/pkg/api"
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
		},
		SigningKey: []byte(api.SecretKey),
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		// This is to be executed after the jwt is verified to be valid.
		return func(c echo.Context) error {
			if c.Request().URL.Path == "/login" {
				// Avoid authorization if /login
				return next(c)
			}
			authorized, err := api.Authorize(c)
			if err != nil {
				log.Println(err)
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
			if !authorized {
				log.Printf("Unauthorized request on %s\f", c.Request().RequestURI)
				return echo.NewHTTPError(http.StatusForbidden)
			}
			// Authorized
			return next(c)
		}
	})
	e.Use(echojwt.WithConfig(jwtconfig))

	api.RegisterHandlers(e, api.NewStrictHandler(
		server,
		// Middlewares
		[]api.StrictMiddlewareFunc{},
	))
	if err := e.Start("127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}
