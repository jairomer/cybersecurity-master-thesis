package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	client "uc3m.es/pilot-cli/api"
)

const (
	mtls_certificate = "todo"
	droneapi         = "http://10.100.242.82"
	host             = "drone-api.com"
)

func addHostHeader(host string) client.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Host = host
		return nil
	}
}

func addJwtHeader(token string) client.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Bearer "+token)
		return nil
	}
}

func moveDroneToTarget(dd *client.DroneData) {
	dd.Target.Altitude = dd.Target.Altitude + float64(rand.IntN(100))
	dd.Target.Longitude = dd.Target.Longitude + float64(rand.IntN(100))
	dd.Target.Latitude = dd.Target.Latitude + float64(rand.IntN(100))
}

func main() {
	fmt.Println("pilot-cli")

	pilotId := "pilot-1"

	pilotState := client.PilotProvisioning{
		Id:     pilotId,
		Drones: []client.DroneData{},
	}
	hc := http.Client{}
	// TODO: Add  mtls certificate
	authToken := ""
	user := client.LoginJSONRequestBody{
		User:     pilotId,
		Password: "test12!",
	}

	{
		log.Println("Login pilot with default credentials")
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
					battlefieldData := client.BattlefieldData{}
					if json.Unmarshal(bdResp.Body, &battlefieldData) != nil {
						log.Println("Invalid response from API")
					} else {
						pilotState.Drones = battlefieldData.Drones
						if len(pilotState.Drones) == 0 {
							log.Println("No drones assigned")
						} else {
							if rand.Float32() > 0.5 {
								// Toss a coint, randomly give orders to one random drone.
								droneToMove := rand.IntN(len(pilotState.Drones))
								moveDroneToTarget(&pilotState.Drones[droneToMove])
								c.SetTargetLocation(
									context.TODO(),
									pilotState.Drones[droneToMove].Id,
									pilotState.Drones[droneToMove].Target,
									addJwtHeader(authToken),
									addHostHeader(host),
								)
							}
						}
					}
				}
			}
			time.Sleep(2 * time.Second)
		}
	}
}
