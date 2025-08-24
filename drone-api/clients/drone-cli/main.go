package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	client "uc3m.es/drone-cli/api"
)

const (
	mtls_certificate = "todo"
	apihost          = "10.101.92.59"
	host             = "cli.drone.com"
	droneapi         = "https://api.drone.com:443"
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

func moveToTarget(dd *client.DroneData) {
	if dd.Location.Altitude < dd.Target.Altitude {
		dd.Location.Altitude += 50
	} else if dd.Location.Altitude > dd.Target.Altitude {
		dd.Location.Altitude -= 50
	}

	if dd.Location.Longitude < dd.Target.Longitude {
		dd.Location.Longitude += 50
	} else if dd.Location.Longitude > dd.Target.Longitude {
		dd.Location.Longitude -= 50
	}

	if dd.Location.Latitude < dd.Target.Longitude {
		dd.Location.Latitude += 50
	} else if dd.Location.Latitude > dd.Target.Longitude {
		dd.Location.Latitude -= 50
	}
}

func inTargetLocation(dd *client.DroneData) bool {
	return dd.Target.Altitude == dd.Location.Altitude &&
		dd.Target.Latitude == dd.Location.Latitude &&
		dd.Target.Longitude == dd.Location.Longitude
}

func main() {
	droneid := "drone-1"

	fmt.Println("drone-cli")
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // TODO: Remove this
			ServerName:         host,
		},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// This is for resolving the TLS CN to a particular IP.
			if addr == "api.drone.com:443" {
				//log.Printf("Dialling address: %s\n", addr)
				return dialer.DialContext(ctx, network, net.JoinHostPort(apihost, "443"))
			}
			return dialer.DialContext(ctx, network, addr)
		},
	}
	hc := http.Client{
		Transport: tr,
	}
	// TODO: Add  mtls certificate
	authToken := ""
	user := client.LoginJSONRequestBody{
		User:     droneid,
		Password: "test12!",
	}

	droneState := client.DroneData{
		Id:       droneid,
		Location: client.Coordinate{},
		Target:   client.Coordinate{},
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
					droneDataArr := client.BattlefieldData{}
					if err := json.Unmarshal(bdResp.Body, &droneDataArr); err != nil {
						log.Printf("Invalid response from API: %s\n", err.Error())
					} else {
						if len(droneDataArr.Drones) == 0 {
							log.Println("Unauthorized access or undefined drone")
						} else {
							// In this client we are expecting only one drone.
							droneState.Target = droneDataArr.Drones[0].Target
						}
					}
				}
			}

			if !inTargetLocation(&droneState) {
				moveToTarget(&droneState)
				c.SetCurrentLocation(
					context.TODO(),
					droneState.Id,
					droneState.Location,
					addJwtHeader(authToken),
					addHostHeader(host),
				)
			}

			time.Sleep(2 * time.Second)
		}
	}

}
