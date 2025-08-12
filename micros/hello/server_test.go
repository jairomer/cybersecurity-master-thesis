package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// In order to test our server, we will need a Request to send in and we'll want
// to spy on what our handler writes to the ResponseWriter.

func TestGETHelloWorld(t *testing.T) {
	t.Run("returns 'Hello World'", func(t *testing.T) {
		// Create a new request
		request, _ := http.NewRequest(http.MethodGet, "/hello/world", nil)

		// Has a spy already made for use called ResponseRecorder
		response := httptest.NewRecorder()

		HelloWorldServer(response, request)

		got := response.Body.String()
		want := "hello world"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestGETHelloCountry(t *testing.T) {
	t.Run("returns 'Hello' for the country in path", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/hello/spain", nil)
		response := httptest.NewRecorder()

		HelloWorldServer(response, request)

		got := response.Body.String()
		want := "hello spain"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestSecret(t *testing.T) {
	secret := "tasty test"
	t.Run("puts a secret into the buffer", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/hello/secret", bytes.NewReader([]byte(secret)))
		response := httptest.NewRecorder()

		HelloWorldServer(response, request)

		got := response.Body.String()
		want := "Ok"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("recover a secret from the buffer", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/hello/secret", nil)
		response := httptest.NewRecorder()

		HelloWorldServer(response, request)

		got := response.Body.String()
		want := secret

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
