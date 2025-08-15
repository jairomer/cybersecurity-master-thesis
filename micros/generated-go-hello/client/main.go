package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	client "uc3m.es/hello/client/api"
)

func main() {
	hc := http.Client{}
	{
		c, err := client.NewClient("http://localhost:8000", client.WithHTTPClient(&hc))
		if err != nil {
			log.Fatal(err)
		}

		resp, err := c.GetHelloWorld(context.TODO(), &client.GetHelloWorldParams{})
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Expected HTTP 200 but received %d", resp.StatusCode)
		}

		hello_world, err := client.ParseGetHelloWorldResponse(resp)
		if err != nil {
			log.Fatal(err)
		}
		if hello_world != nil {
			fmt.Printf("%s %s\n", hello_world.JSON200.Hello, hello_world.JSON200.Country)
		}
	}
	{
		c, err := client.NewClient("http://127.0.0.1:8000", client.WithHTTPClient(&hc))
		if err != nil {
			log.Fatal(err)
		}
		country := "Spain"
		params := client.GetHelloWorldParams{Country: &country}
		resp, err := c.GetHelloWorld(context.TODO(), &params)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Expected HTTP 200 but received %d", resp.StatusCode)
		}

		hello_world, err := client.ParseGetHelloWorldResponse(resp)
		if err != nil {
			log.Fatal(err)
		}
		if hello_world != nil {
			fmt.Printf("%s %s\n", hello_world.JSON200.Hello, hello_world.JSON200.Country)
		}
	}
}
