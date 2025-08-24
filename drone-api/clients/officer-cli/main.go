package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	client "uc3m.es/officer-cli/api"
)

const (
	host     = "api.drone.com"
	droneapi = "https://api.drone.com:443"
)

func addHostHeader(host string) client.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Host = "officer.drone.com"
		return nil
	}
}

func addJwtHeader(token string) client.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", "Bearer "+token)
		return nil
	}
}

func main() {
	caPath := flag.String("ca", "certs/ca.crt", "Path to CA certificate")
	certPath := flag.String("clientcert", "certs/cert.crt", "Path to client certificate")
	keyPath := flag.String("clientkey", "certs/cert.key", "Path to client key")
	apihost := flag.String("apihost", "10.101.92.59", "IP for drone API gateway")
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

	authToken := ""
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
