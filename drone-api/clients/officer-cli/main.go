package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	client "uc3m.es/officer-cli/api"
)

const (
	mtls_certificate = "todo"
	//droneapi         = "http://localhost:8000"
	droneapi = "http://10.100.242.82"
	host     = "drone-api.com"
)

func addHostHeader(host string) client.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Host = host
		return nil
	}
}

func addJwtHeader(token string) client.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Host = host
		req.Header.Add("Authorization", "Bearer "+token)
		return nil
	}
}

func main() {
	// TODO: Add  mtls certificate
	authToken := ""
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		},
	}
	hc := http.Client{
		Transport: tr,
	}

	user := client.LoginJSONRequestBody{
		User:     "officer-1",
		Password: "changeme",
	}

	log.Println("officer-cli")
	{
		log.Println("Login officer with default credentials")
		// Provision battlefield.
		c, err := client.NewClient(droneapi, client.WithHTTPClient(&hc))
		if err != nil {
			log.Fatal(err)
		}

		resp, err := c.Login(context.TODO(), user, addHostHeader(host))
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
		}
		log.Printf("Login with default credentials was successful.")
	}
	{
		c, err := client.NewClient(droneapi, client.WithHTTPClient(&hc))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Provisioning battlefield")
		provisioning := client.BattlefieldProvision{
			Credentials: []client.UserProvision{
				client.UserProvision{User: "officer-1", Password: "test12!", Role: client.Officer},
				client.UserProvision{User: "pilot-1", Password: "test12!", Role: client.Pilot},
				client.UserProvision{User: "pilot-2", Password: "test12!", Role: client.Pilot},
				client.UserProvision{User: "drone-1", Password: "test12!", Role: client.Drone},
				client.UserProvision{User: "drone-2", Password: "test12!", Role: client.Drone},
				client.UserProvision{User: "drone-3", Password: "test12!", Role: client.Drone},
				client.UserProvision{User: "drone-4", Password: "test12!", Role: client.Drone},
				client.UserProvision{User: "drone-5", Password: "test12!", Role: client.Drone},
				client.UserProvision{User: "drone-6", Password: "test12!", Role: client.Drone},
			},
			Pilots: []client.PilotProvisioning{
				client.PilotProvisioning{
					Id: "pilot-1",
					Drones: []client.DroneData{
						client.DroneData{
							Id: "drone-1",
							Location: client.Coordinate{
								Altitude:  0.0,
								Longitude: 0.0,
								Latitude:  0.0,
							},
							Target: client.Coordinate{
								Altitude:  0.0,
								Longitude: 0.0,
								Latitude:  0.0,
							},
						},
						client.DroneData{
							Id: "drone-2",
							Location: client.Coordinate{
								Altitude:  0.0,
								Longitude: 0.0,
								Latitude:  0.0,
							},
							Target: client.Coordinate{
								Altitude:  0.0,
								Longitude: 0.0,
								Latitude:  0.0,
							},
						},
						client.DroneData{
							Id: "drone-3",
							Location: client.Coordinate{
								Altitude:  0.0,
								Longitude: 0.0,
								Latitude:  0.0,
							},
							Target: client.Coordinate{
								Altitude:  0.0,
								Longitude: 0.0,
								Latitude:  0.0,
							},
						},
					},
				},
				client.PilotProvisioning{
					Id: "pilot-2",
					Drones: []client.DroneData{
						client.DroneData{
							Id: "drone-4",
							Location: client.Coordinate{
								Altitude:  0.0,
								Longitude: 0.0,
								Latitude:  0.0,
							},
							Target: client.Coordinate{
								Altitude:  0.0,
								Longitude: 0.0,
								Latitude:  0.0,
							},
						},
						client.DroneData{
							Id: "drone-5",
							Location: client.Coordinate{
								Altitude:  0.0,
								Longitude: 0.0,
								Latitude:  0.0,
							},
							Target: client.Coordinate{
								Altitude:  0.0,
								Longitude: 0.0,
								Latitude:  0.0,
							},
						},
						client.DroneData{
							Id: "drone-6",
							Location: client.Coordinate{
								Altitude:  0.0,
								Longitude: 0.0,
								Latitude:  0.0,
							},
							Target: client.Coordinate{
								Altitude:  0.0,
								Longitude: 0.0,
								Latitude:  0.0,
							},
						},
					},
				},
			},
		}

		resp, err := c.BattlefieldProvision(
			context.TODO(),
			provisioning,
			addJwtHeader(authToken),
			addHostHeader(host),
		)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Expected HTTP 200 but received %s", resp.Status)
		}

		provResp, err := client.ParseBattlefieldProvisionResponse(resp)
		if err != nil {
			log.Fatal(err)
		}
		if provResp == nil {
			log.Fatal("Empty provisioning result received.")
		}

		log.Println(string(provResp.Body))
	}
	{
		c, err := client.NewClient(droneapi, client.WithHTTPClient(&hc))
		if err != nil {
			log.Fatal(err)
		}

		for {
			fmt.Print("\033[H\033[2J")
			fmt.Printf("Monitoring battlefield as: %s", user.User)
			resp, err := c.GetBattlefieldData(
				context.TODO(),
				addJwtHeader(authToken),
				addHostHeader(host),
			)
			if err != nil {
				log.Println(err)
			} else {
				bdResp, err := client.ParseGetBattlefieldDataResponse(resp)
				if err != nil {
					log.Println(err)
				} else {
					log.Println(string(bdResp.Body))
				}
			}
			time.Sleep(2 * time.Second)
		}
	}

}
