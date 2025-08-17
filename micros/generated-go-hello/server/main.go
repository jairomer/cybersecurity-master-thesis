package main

import (
	"log"

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
			log.Println("Request authenticated, authorizing...")
			// TODO: Verify rego rules here
		},
		SigningKey: []byte(api.SecretKey),
	}

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
