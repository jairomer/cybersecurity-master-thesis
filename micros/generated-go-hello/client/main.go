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

	authToken := ""
	{
		c, err := client.NewClient("http://localhost:8000", client.WithHTTPClient(&hc))
		if err != nil {
			log.Fatal(err)
		}

		user := client.LoginJSONRequestBody{User: "test1", Password: "test"}
		resp, err := c.Login(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode == 200 {
			loginRes, err := client.ParseLoginResponse(resp)
			if err != nil {
				log.Fatal(err)
			}
			authToken = loginRes.JSON200.Token
		} else {
			log.Fatalf("Authentication failed: %s", resp.Status)
			return
		}
		log.Printf("Login in as %s with Authentication: Bearer %s", "test1", authToken)
	}
	{
		c, err := client.NewClient("http://localhost:8000", client.WithHTTPClient(&hc))
		if err != nil {
			log.Fatal(err)
		}

		resp, err := c.GetHelloWorld(
			context.TODO(),
			&client.GetHelloWorldParams{},
			func(ctx context.Context, req *http.Request) error {
				req.Header.Add("Authorization", "Bearer "+authToken)
				return nil
			})

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
		resp, err := c.GetHelloWorld(context.TODO(), &params,
			func(ctx context.Context, req *http.Request) error {
				req.Header.Add("Authorization", "Bearer "+authToken)
				return nil
			})
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
		c, err := client.NewClient("http://localhost:8000", client.WithHTTPClient(&hc))
		if err != nil {
			log.Fatal(err)
		}

		user := client.LoginJSONRequestBody{User: "test2", Password: "test"}
		resp, err := c.Login(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode == 200 {
			loginRes, err := client.ParseLoginResponse(resp)
			if err != nil {
				log.Fatal(err)
			}
			authToken = loginRes.JSON200.Token
		} else {
			log.Fatalf("Authentication failed: %s", resp.Status)
			return
		}
		log.Printf("Login in as %s with Authentication: Bearer %s", "test2", authToken)
	}
	{
		c, err := client.NewClient("http://localhost:8000", client.WithHTTPClient(&hc))
		if err != nil {
			log.Fatal(err)
		}

		resp, err := c.GetHelloWorld(
			context.TODO(),
			&client.GetHelloWorldParams{},
			func(ctx context.Context, req *http.Request) error {
				req.Header.Add("Authorization", "Bearer "+authToken)
				return nil
			})

		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusForbidden {
			log.Fatalf("Expected HTTP 403 but received %d", resp.StatusCode)
		}
		log.Print("Request Forbidden")
	}
	{
		c, err := client.NewClient("http://127.0.0.1:8000", client.WithHTTPClient(&hc))
		if err != nil {
			log.Fatal(err)
		}
		country := "Spain"
		params := client.GetHelloWorldParams{Country: &country}
		resp, err := c.GetHelloWorld(context.TODO(), &params,
			func(ctx context.Context, req *http.Request) error {
				req.Header.Add("Authorization", "Bearer "+authToken)
				return nil
			})
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusForbidden {
			log.Fatalf("Expected HTTP 403 but received %d", resp.StatusCode)
		}
		log.Print("Request Forbidden")
	}
}
