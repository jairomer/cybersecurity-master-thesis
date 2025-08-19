package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type HelloServer struct {
	secret string
}

func SayHello() string {
	if to == "world" {
		return "hello world"
	} else if to == "spain" {
		return "hello spain"
	} else {
		return ""
	}
}

func (s *HelloServer) StoreSecret(w http.ResponseWriter, r *http.Request) {
	_, err := r.Body.Read([]byte(s.secret))
	if err != nil {
		fmt.Fprintf(w, "KO")
		return
	}
	fmt.Fprintf(w, "OK")
}

func (s *HelloServer) RecoverSecret(w http.ResponseWriter, r *http.Request) string {
	return s.secret
}

func (s *HelloServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := strings.TrimPrefix(r.URL.Path, "/hello/")

	switch r.Method {
	case http.MethodGet:
		if route == "secret" {
			s.RecoverSecret(w, r)
			return
		} else {
			r.Context().Value()
			SayHello(w, r)
			return
		}
	case http.MethodPut:
		s.StoreSecret(w, r)
		return
	}

	fmt.Printf("Processing request at %s\n", r.URL)
	fmt.Fprint(w, RouteTo(route, r))
}
