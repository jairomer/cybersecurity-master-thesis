package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
	"os"
	"time"

	client "uc3m.es/pilot-cli/api"
)

const (
	droneapi = "https://api.drone.com:443"
)

func addHostHeader() client.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Host = "pilot.drone.com"
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

	caPath := flag.String("ca", "certs/ca.crt", "Path to CA certificate")
	certPath := flag.String("clientcert", "certs/cert.crt", "Path to client certificate")
	keyPath := flag.String("clientkey", "certs/cert.key", "Path to client key")
	apihost := flag.String("apihost", "10.101.92.59", "IP for drone API gateway")
	pilotid := flag.String("pilotid", "pilot-1", "Provisioned id for the pilot")
	password := flag.String("password", "test12!", "Provisioned password for the pilot")
	flag.Parse()

	cert, err := tls.LoadX509KeyPair(*certPath, *keyPath)
	if err != nil {
		log.Fatal("Failed to load client cert/key pair: %w", err)
	}

	caCert, err := os.ReadFile(*caPath)
	if err != nil {
		log.Fatal("Failed to load CA file: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatal("Failed to append CA cert")
	}

	pilotState := client.PilotProvisioning{
		Id:     *pilotid,
		Drones: []client.DroneData{},
	}

	dialer := &net.Dialer{Timeout: 10 * time.Second}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName:   "api.drone.com",
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
			MinVersion:   tls.VersionTLS12,
		},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// This is for resolving the TLS CN to a particular IP.
			if apihost != nil {
				hst := *apihost
				if addr == "api.drone.com:443" {
					// log.Printf("Dialling address: %s\n", addr)
					return dialer.DialContext(ctx, network, net.JoinHostPort(hst, "443"))
				}
			}
			return dialer.DialContext(ctx, network, addr)
		},
	}
	hc := http.Client{
		Transport: tr,
	}

	user := client.LoginJSONRequestBody{
		User:     *pilotid,
		Password: *password,
	}

	authToken := ""
	{
		log.Println("Login pilot with default credentials")
		c, err := client.NewClient(droneapi, client.WithHTTPClient(&hc))
		if err != nil {
			log.Fatal(err)
		}
		resp, err := c.Login(context.TODO(), user, addHostHeader())
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
				addHostHeader(),
			)
			if err != nil {
				log.Println(err)
			} else {
				bdResp, err := client.ParseGetBattlefieldDataResponse(resp)
				if err != nil {
					log.Println(err)
				} else {
					if bdResp.StatusCode() != 200 {
						log.Println(bdResp.Status())
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
										addHostHeader(),
									)
								}
							}
						}
					}
				}
			}
			time.Sleep(2 * time.Second)
		}
	}
}
