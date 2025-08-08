package main

import (
	"fmt"
	"net/http"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func HelloWorldServer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Processing request at %s\n", r.URL)
	fmt.Fprintf(w, "hello world")
}
