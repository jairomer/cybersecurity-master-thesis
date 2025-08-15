package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"uc3m.es/hello/server/pkg/api"
)

func main() {
	server := api.NewServer()
	e := echo.New()
	api.RegisterHandlers(e, api.NewStrictHandler(
		server,
		// Middlewares
		[]api.StrictMiddlewareFunc{},
	))
	if err := e.Start("127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}
