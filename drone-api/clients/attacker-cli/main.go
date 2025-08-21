package main

import (
	"context"
	"net/http"

	client "uc3m.es/attacker-cli/api"
)

const (
	mtls_officer_certificate = "todo"
	mtls_pilot_certificate   = "todo"
	mtls_drone_certificate   = "todo"
	droneapi                 = "http://localhost:8000"
)

func addJwtHeader(token string) client.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Bearer "+token)
		return nil
	}
}

func officerAttackBattery() {
	// Scenario 1: Attempt connection to officer endpoints without certificate.
	// 	Expected: Failure to connect.

	// Scenario 2: Officer certificate has been disclosed.
	// 	Expected: mTLS successful, receive 401 to officer endpoints.

	// Scenario 2A: Attempt to access drone endpoints with officer certificate.
	//	Expected: mTLS failure

	// Scenario 2B: Attempt to access pilot endpoints with officer certificate.
	//	Expected: mTLS failure
}

func pilotAttackBattery() {
	// Scenario 1: Attempt connection to pilot endpoints without  certificate.
	// 	Expected: Failure to connect.

	// Scenario 2: Pilot certificate has been disclosed.
	// 	Expected: mTLS successful, receive 401 to pilot endpoints.

	// Scenario 2A: Attempt to access drone endpoints with pilot certificate.
	// 	Expected: mTLS fails

	// Scenario 2B: Attempt to access officer endpoints with pilot certificate.
	// 	Expected: mTLS fails

	// Scenario 3: Pilot credentials have been disclosed.
	// 	Expected:
	// 		- Pilot data integrity and confidentiality has been compromissed!
	// 		- Data integrity and confidentiality for other pilots and drones of other pilots is secured.
}

func droneAttackBattery() {
	// Scenario 1: Attempt connection to Drone endpoints without  certificate.

	// Scenario 2: Drone certificate has been disclosed.
	// 	Expected: mTLS successful, receive 401 to pilot endpoints.

	// Scenario 2A: Attempt to access officer endpoints with drone certificate.
	// 	Expected: mTLS fails.

	// Scenario 2B: Attempt to access pilot endpoints with drone certificates.
	// 	Expected: mTLS fails.

	// Scenario 3: Drone credentials have been disclosed.
	// 	Expected:
	// 		- Drone data integrity and confidentiality has been compromissed
	//		- Data integrity and confidentiality for other drones, pilots or officers is maintained.
}

func main() {
	/*
		   log.Println("attacker-cli")

		hc := http.Client{}
		// TODO: Add  mtls certificate
		authToken := ""
		user := client.LoginJSONRequestBody{
			User:     "officer-1",
			Password: "changeme",
		}
		   	{
		   		log.Println("Login officer with default credentials")
		   		// Provision battlefield.
		   		c, err := client.NewClient(droneapi, client.WithHTTPClient(&hc))
		   		if err != nil {
		   			log.Fatal(err)
		   		}

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
	*/
}
