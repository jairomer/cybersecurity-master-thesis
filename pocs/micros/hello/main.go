package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Start a web server listening on a port
	// Create a goroutine for every request
	// Run the goroutine against a Handler
	server := &HelloServer{}
	fmt.Printf("Listening...\n")
	log.Fatal(http.ListenAndServe(":5000", server))
	fmt.Printf("Exiting\n")
}
