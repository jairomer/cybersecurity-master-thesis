package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// In order to test our server, we will need a Request to send in and we'll want
// to spy on what our handler writes to the ResponseWriter.

func TestGETPlayers(t *testing.T) {
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
