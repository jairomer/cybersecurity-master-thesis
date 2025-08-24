package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	client "uc3m.es/attacker-cli/api"
)

const (
	droneapi = "https://api.drone.com:443"
)

//	func addJwtHeader(token string) client.RequestEditorFn {
//		return func(ctx context.Context, req *http.Request) error {
//			req.Header.Add("Authorization", "Bearer "+token)
//			return nil
//		}
//	}
func addHostHeader(host string) client.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Host = host
		return nil
	}
}

func getClient(caPath, certPath, keyPath, apihost *string) *client.Client {
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	hc := http.Client{}
	if caPath != nil && certPath != nil && keyPath != nil {
		cert, err := tls.LoadX509KeyPair(*certPath, *keyPath)
		if err != nil {
			log.Printf("Failed to load client cert/key pair: %w", err)
		}

		caCert, err := os.ReadFile(*caPath)
		if err != nil {
			log.Printf("Failed to load CA file: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			log.Printf("Failed to append CA cert")
		}
		hc.Transport = &http.Transport{
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
	} else {
		hc.Transport = &http.Transport{
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
	}
	c, err := client.NewClient(droneapi, client.WithHTTPClient(&hc))
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func officerAttackBattery(caPath, certPath, keyPath, apihost *string) {
	host := "officer.drone.com"
	{
		// Scenario 1: Attempt connection to officer endpoints without certificate.
		// 	Expected: Failure to connect.
		c := getClient(nil, nil, nil, apihost)
		log.Println("Attept to provision battlefield without certificate nor authentication...")
		provisioning := client.BattlefieldProvision{
			Credentials: []client.UserProvision{},
		}
		if _, err := c.BattlefieldProvision(context.TODO(), provisioning, addHostHeader(host)); err == nil {
			log.Printf("Attacker was able to access battlefield provision endpoint without credentials nor certificate.")
		} else {
			log.Printf("Attacker blocked: %s", err.Error())
		}
	}
	// Scenario 2: Officer certificate has been disclosed.
	// 	Expected: mTLS successful, receive 401 to officer endpoints.
	{
		// Scenario 2: Officer certificate has been disclosed.
		// 	Expected: mTLS successful, receive 401 to officer endpoints.
		c := getClient(caPath, certPath, keyPath, apihost)
		log.Println("Attept to provision battlefield without authentication...")
		provisioning := client.BattlefieldProvision{
			Credentials: []client.UserProvision{},
		}
		if resp, err := c.BattlefieldProvision(context.TODO(), provisioning, addHostHeader(host)); err == nil {
			if resp.StatusCode == 200 {
				log.Printf("Attacker was able to access battlefield provision endpoint without credentials.")
			} else {
				log.Printf("Attacker blocked: %s", resp.Status)
			}
		} else {
			log.Printf("Attacker blocked: %s", err.Error())
		}
	}

	// Scenario 2A: Attempt to access drone endpoints with officer certificate.
	//	Expected: mTLS failure

	// Scenario 2B: Attempt to access pilot endpoints with officer certificate.
	//	Expected: mTLS failure
}

func pilotAttackBattery(caPath, certPath, keyPath, apihost *string) {
	host := "pilot.drone.com"
	{
		// Scenario 1: Attempt connection to pilot endpoints without  certificate.
		// 	Expected: Failure to connect.
		c := getClient(nil, nil, nil, apihost)
		log.Println("Attacker takeover of a drone without certificate nor authentication...")
		_, err := c.SetTargetLocation(
			context.TODO(),
			"drone-1",
			client.Coordinate{Altitude: 0, Longitude: 0, Latitude: 0},
			addHostHeader(host),
		)
		if err == nil {
			log.Printf("Attacker was able to access pilot endpoint without credentials nor certificate.")
		} else {
			log.Printf("Attacker blocked: %s", err.Error())
		}
	}
	{
		// Scenario 2: Pilot certificate has been disclosed.
		// 	Expected: mTLS successful, receive 401 to pilot endpoints.
		c := getClient(caPath, certPath, keyPath, apihost)
		log.Println("Attacker takeover of a drone without authentication...")
		resp, err := c.SetTargetLocation(
			context.TODO(),
			"drone-1",
			client.Coordinate{Altitude: 0, Longitude: 0, Latitude: 0},
			addHostHeader(host),
		)
		if err == nil {
			if resp.StatusCode == 200 {
				log.Printf("Attacker was able to execute drone takeover without credentials.")
			} else {
				log.Printf("Attacker blocked: %s", resp.Status)
			}
		} else {
			log.Printf("Attacker blocked: %s", err.Error())
		}
	}

	// Scenario 2A: Attempt to access drone endpoints with pilot certificate.
	// 	Expected: mTLS fails

	// Scenario 2B: Attempt to access officer endpoints with pilot certificate.
	// 	Expected: mTLS fails

	// Scenario 3: Pilot credentials have been disclosed.
	// 	Expected:
	// 		- Pilot data integrity and confidentiality has been compromissed!
	// 		- Data integrity and confidentiality for other pilots and drones of other pilots is secured.
}

func droneAttackBattery(caPath, certPath, keyPath, apihost *string) {
	// Scenario 1: Attempt connection to Drone endpoints without  certificate.
	host := "cli.drone.com"
	{
		// Scenario 1: Attempt connection to pilot endpoints without  certificate.
		// 	Expected: Failure to connect.
		c := getClient(nil, nil, nil, apihost)
		log.Println("Attacker attempt to spoof a drone location without certificate nor authentication...")
		_, err := c.SetCurrentLocation(
			context.TODO(),
			"drone-1",
			client.Coordinate{Altitude: 0, Longitude: 0, Latitude: 0},
			addHostHeader(host),
		)
		if err == nil {
			log.Println("Attacker was able to access drone endpoint without credentials nor certificate.")
		} else {
			log.Printf("Attacker blocked: %s", err.Error())
		}
	}
	{
		{
			// Scenario 2: Drone certificate has been disclosed.
			// 	Expected: mTLS successful, receive 401 to pilot endpoints.
			c := getClient(caPath, certPath, keyPath, apihost)
			log.Println("Attacker attempt to spoof a drone location without authentication...")
			resp, err := c.SetCurrentLocation(
				context.TODO(),
				"drone-1",
				client.Coordinate{Altitude: 0, Longitude: 0, Latitude: 0},
				addHostHeader(host),
			)
			if err == nil {
				if resp.StatusCode == 200 {
					log.Println("Attacker was able to execute drone takeover without credentials.")
				} else {
					log.Printf("Attacker blocked: %s\n", resp.Status)
				}
			} else {
				log.Printf("Attacker blocked: %s\n", err.Error())
			}
		}
		{
			// Scenario 2A: Attempt to access officer endpoints with drone certificate.
			// 	Expected: mTLS fails.
			c := getClient(caPath, certPath, keyPath, apihost)
			log.Println("Attept to provision battlefield with a drone certificate without authentication...")
			provisioning := client.BattlefieldProvision{
				Credentials: []client.UserProvision{},
			}
			if _, err := c.BattlefieldProvision(context.TODO(), provisioning, addHostHeader(host)); err == nil {
				log.Printf("Attacker was able to access battlefield provision endpoint without officer certificate using a drone certificate.")
			} else {
				log.Printf("Attacker blocked: %s", err.Error())
			}
		}
		{
			// Scenario 2B: Attempt to access pilot endpoints with a drone certificate.
			// 	Expected: mTLS fails.
			c := getClient(caPath, certPath, keyPath, apihost)
			log.Println("Attacker takeover of a drone with a drone certificate without authentication...")
			_, err := c.SetTargetLocation(
				context.TODO(),
				"drone-1",
				client.Coordinate{Altitude: 0, Longitude: 0, Latitude: 0},
				addHostHeader(host),
			)
			if err == nil {
				log.Println("Attacker was able to access a pilot specific endpoint with a drone certificate.")
			} else {
				log.Printf("Attacker blocked: %s", err.Error())
			}
		}
	}
}

func login_cracking(caPath, certPath, keyPath, apihost *string) {
	host := "cli.drone.com"
	{
		// Scenario 1: Attempt connection to pilot endpoints without certificate.
		// 	Expected: Failure to connect.
		c := getClient(nil, nil, nil, apihost)
		log.Println("Attacker attempt to login as a drone location without certificate nor authentication...")
		_, err := c.Login(context.TODO(), client.UserLogin{User: "test", Password: "test"}, addHostHeader(host))
		if err == nil {
			log.Println("Attacker was able to access drone endpoint without credentials nor certificate.")
		} else {
			log.Printf("Attacker blocked: %s", err.Error())
		}
	}
	{
		// Scenario 2: After successfully disclosing a drone certificate, a brute force attack is attempted.
		// 	Expected: 401.
		c := getClient(caPath, certPath, keyPath, apihost)
		log.Println("Attacker attempt to login as a drone location without certificate nor authentication...")

		for {
			// This is just for traffic generation.
			// At this point, a successful attack will most likely be possible if the attacker creates a good Dictionary
			// of if the defender uses weak and predictable passwords and users.
			resp, err := c.Login(context.TODO(), client.UserLogin{User: "test", Password: "test"}, addHostHeader(host))
			if err == nil {
				if resp.StatusCode == 200 {
					log.Printf("Attacker was able to Login.")
				} else {
					log.Printf(".")
					//log.Printf("Attacker blocked: %s", resp.Status)
				}
			} else {
				//log.Printf("Attacker blocked: %s", err.Error())
			}
		}
	}
}

func main() {
	caPath := flag.String("ca", "certs/ca.crt", "Path to CA certificate")

	officerCertPath := flag.String("officercert", "certs/cert.crt", "Path to officer certificate")
	officerKeyPath := flag.String("officerkey", "certs/cert.key", "Path to officer key")

	pilotCertPath := flag.String("pilotcert", "certs/cert.crt", "Path to pilot certificate")
	pilotKeyPath := flag.String("pilotkey", "certs/cert.key", "Path to pilot key")

	droneCertPath := flag.String("dronecert", "certs/cert.crt", "Path to drone certificate")
	droneKeyPath := flag.String("dronekey", "certs/cert.key", "Path to drone key")

	apihost := flag.String("apihost", "10.101.92.59", "IP for drone API gateway")

	loginCracking := flag.String("crack", "", "Add to enable login bruteforce.")
	flag.Parse()

	log.Println("===============================")
	log.Println("= Attacking Officer Endpoints =")
	log.Println("===============================")
	officerAttackBattery(
		caPath,
		officerCertPath,
		officerKeyPath,
		apihost,
	)

	log.Println("=============================")
	log.Println("= Attacking Pilot Endpoints =")
	log.Println("=============================")
	pilotAttackBattery(
		caPath,
		pilotCertPath,
		pilotKeyPath,
		apihost,
	)

	log.Println("=============================")
	log.Println("= Attacking Drone Endpoints =")
	log.Println("=============================")
	droneAttackBattery(
		caPath,
		droneCertPath,
		droneKeyPath,
		apihost,
	)

	if loginCracking != nil {
		log.Println("=============================")
		log.Println("= Attacking Login Endpoints =")
		log.Println("=============================")
		login_cracking(
			caPath,
			droneCertPath,
			droneKeyPath,
			apihost,
		)
	}
}
