package api

import (
	"context"
	"crypto/sha512"
	"fmt"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/open-policy-agent/opa/v1/rego"
)

const (
	SALT = "saltysalt!"
)

type XFCC struct {
	Value string
	role  string
}

func (xfcc *XFCC) GetClientRole() *string {
	if xfcc.role != "" {
		return &xfcc.role
	}
	index := strings.Index(xfcc.Value, "officer.drone.api")
	if index != -1 {
		xfcc.role = "officer"
		return &xfcc.role
	}
	index = strings.Index(xfcc.Value, "pilot.drone.api")
	if index != -1 {
		xfcc.role = "pilot"
		return &xfcc.role
	}
	index = strings.Index(xfcc.Value, "cli.drone.api")
	if index != -1 {
		xfcc.role = "drone"
		return &xfcc.role
	}
	return nil
}

type UserDatabase struct {
	Officers map[string]string
	Pilots   map[string]string
	Drones   map[string]string
}

type AuthorizationDecision struct {
	user      string
	role      string
	uri       string
	operation string
	drones    []string
	allowed   bool
}

const (
	operationSetTarget      = "set-target"
	operationGetTarget      = "get-target"
	operationSetLocation    = "set-location"
	operationGetLocation    = "get-location"
	operationProvisioning   = "provisioning"
	operationGetBattelfield = "get-battlefield"
)

type Server struct {
	Users  UserDatabase
	Pilots []PilotProvisioning
}

func NewServer() Server {

	// Default admin user
	hash := sha512.New()
	aggr := "officer-1" + "changeme" + SALT
	hash.Write([]byte(aggr))

	s := Server{
		Users: UserDatabase{
			Officers: map[string]string{"officer-1": string(hash.Sum(nil))},
			Pilots:   map[string]string{},
			Drones:   map[string]string{},
		},
		Pilots: []PilotProvisioning{},
	}
	return s
}

func (s *Server) GetBattlefieldData(ctx context.Context, request GetBattlefieldDataRequestObject) (GetBattlefieldDataResponseObject, error) {
	authzDecision, err := s.authorized(ctx, operationGetBattelfield)
	if err != nil {
		log.Printf("Error running authorization at %s by user %s", authzDecision.uri, authzDecision.user)
		return GetBattlefieldData403Response{}, nil
	}

	if !authzDecision.allowed {
		log.Printf("Unauthorized access blocked at %s by user %s", authzDecision.uri, authzDecision.user)
		return GetBattlefieldData403Response{}, nil
	}
	dd := []DroneData{}
	for _, droneid := range authzDecision.drones {
		for i, _ := range s.Pilots {
			for _, drone := range s.Pilots[i].Drones {
				if drone.Id == droneid {
					dd = append(dd, drone)
				}
			}
		}
	}
	return GetBattlefieldData200JSONResponse{dd}, nil
}

func (s *Server) SetCurrentLocation(ctx context.Context, request SetCurrentLocationRequestObject) (SetCurrentLocationResponseObject, error) {
	authzDecision, err := s.authorized(ctx, operationSetLocation)
	if err != nil {
		log.Printf("Error running authorization at %s by user %s", authzDecision.uri, authzDecision.user)
		return SetCurrentLocation403Response{}, nil
	}
	if authzDecision.allowed {
		for i, _ := range s.Pilots {
			for j, _ := range s.Pilots[i].Drones {
				if s.Pilots[i].Drones[j].Id == authzDecision.user {
					s.Pilots[i].Drones[j].Location.Altitude = request.Body.Altitude
					s.Pilots[i].Drones[j].Location.Latitude = request.Body.Latitude
					s.Pilots[i].Drones[j].Location.Longitude = request.Body.Longitude
					return SetCurrentLocation200JSONResponse(s.Pilots[i].Drones[j]), nil
				}
			}
		}
	}
	return SetCurrentLocation403Response{}, nil
}

func (s *Server) SetTargetLocation(ctx context.Context, request SetTargetLocationRequestObject) (SetTargetLocationResponseObject, error) {
	authzDecision, err := s.authorized(ctx, operationSetTarget)
	if err != nil {
		log.Printf("Error running authorization at %s by user %s", authzDecision.uri, authzDecision.user)
		return SetTargetLocation403Response{}, nil
	}
	if authzDecision.allowed {
		for _, droneid := range authzDecision.drones {
			if droneid == request.Droneid {
				for i, _ := range s.Pilots {
					if s.Pilots[i].Id == authzDecision.user {
						// user of this endpoint is always a pilot.
						for j, _ := range s.Pilots[i].Drones {
							if s.Pilots[i].Drones[j].Id == request.Droneid {
								s.Pilots[i].Drones[j].Target.Altitude = request.Body.Altitude
								s.Pilots[i].Drones[j].Target.Latitude = request.Body.Latitude
								s.Pilots[i].Drones[j].Target.Longitude = request.Body.Longitude
								return SetTargetLocation200JSONResponse(s.Pilots[i].Drones[j]), nil
							}
						}
					}
				}
			}
		}
	}
	return SetTargetLocation403Response{}, nil
}

func (s *Server) BattlefieldProvision(ctx context.Context, request BattlefieldProvisionRequestObject) (BattlefieldProvisionResponseObject, error) {
	jwtStr := ctx.Value("jwt").(string)
	if jwtStr == "" {
		log.Fatal("JWT received nil or emtpy from context")
	}

	token, err := jwt.ParseWithClaims(jwtStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		log.Printf("Error on battlefield provision: %s\n", err.Error())
		return BattlefieldProvision403Response{}, nil
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		user := claims.ID
		// This endpoint should only be accessed by officers.
		if _, ok := s.Users.Officers[user]; !ok {
			// Not an officer.
			log.Printf("Unauthorized officer attempted a new battlefield provision: %s\n", user)
			return BattlefieldProvision403Response{}, nil
		}
	} else {
		return BattlefieldProvision403Response{}, nil
	}

	// Add in credentials
	for _, cred := range request.Body.Credentials {
		hash := sha512.New()
		aggr := cred.User + cred.Password + SALT
		hash.Write([]byte(aggr))

		if cred.Role == "officer" {
			s.Users.Officers[cred.User] = string(hash.Sum(nil))
		} else if cred.Role == "pilot" {
			s.Users.Pilots[cred.User] = string(hash.Sum(nil))
		} else if cred.Role == "drone" {
			s.Users.Drones[cred.User] = string(hash.Sum(nil))
		} else {
			log.Printf("Error while provisioning battlefield: Unknown role specified for new user '%s' of role '%s'\n ", cred.User, cred.Role)
			return BattlefieldProvision400Response{}, nil
		}
	}

	// Add in the drones and the pilots.
	s.Pilots = append(s.Pilots, request.Body.Pilots...)
	dd := []DroneData{}
	for _, pilot := range s.Pilots {
		for _, drone := range pilot.Drones {
			inBattlefield := false
			for _, droneInBattlefield := range dd {
				if droneInBattlefield.Id == drone.Id {
					// Repeated
					inBattlefield = true
					break
				}
			}
			if !inBattlefield {
				dd = append(dd, drone)
			}
		}
	}
	return BattlefieldProvision200JSONResponse{dd}, nil
}

func (s *Server) authenticate(user, password string) bool {
	hash := sha512.New()
	aggr := user + password + SALT
	hash.Write([]byte(aggr))
	hsh, ok := s.Users.Officers[user]
	if !ok {
		hsh, ok = s.Users.Pilots[user]
	}
	if !ok {
		hsh, ok = s.Users.Drones[user]
	}
	if !ok {
		log.Printf("Failed to autenticate: user '%s' not found\n", user)
	}
	if hsh == string(hash.Sum(nil)) {
		log.Printf("User %s has been authenticated successfully.\n", user)
		return true
	}
	log.Printf("Failed to authenticate:  user '%s' presented invalid credentials\n", user)
	return false
}

func (s *Server) Login(ctx context.Context, request LoginRequestObject) (LoginResponseObject, error) {

	if request.Body == nil || request.Body.Password == "" || request.Body.User == "" {
		log.Printf("Error decoding request on /login")
		return Login401Response{}, nil
	}
	evaluatedUser := request.Body.User
	if s.authenticate(evaluatedUser, request.Body.Password) {
		log.Printf("User %s has been authenticated\n", evaluatedUser)
		request.Body.Password = "" // clear it asap
		token, err := GenerateToken(request.Body.User)
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			return Login401Response{}, err
		}
		log.Printf("Login successful")
		return Login200JSONResponse{Token: token}, nil
	}
	log.Printf("User %s has failed authentication\n", evaluatedUser)
	return Login401Response{}, nil
}

func (s *Server) getRole(user string) (string, error) {
	_, ok := s.Users.Officers[user]
	if ok {
		return "officer", nil
	}
	_, ok = s.Users.Pilots[user]
	if ok {
		return "pilot", nil
	}
	_, ok = s.Users.Drones[user]
	if ok {
		return "drone", nil
	}
	return "", fmt.Errorf("unknown role for user %s", user)
}

func (s *Server) populateAuthzData(ctx context.Context, authzDecision *AuthorizationDecision) interface{} {
	log.Printf("Authorizing: %s accessing %s\n", authzDecision.user, authzDecision.uri)
	var officers []string
	for k, _ := range s.Users.Officers {
		officers = append(officers, k)
	}
	var pilots []map[string]interface{}
	for i, _ := range s.Pilots {
		drones := []string{}
		for _, drone := range s.Pilots[i].Drones {
			drones = append(drones, drone.Id)
		}
		pilots = append(pilots, map[string]interface{}{"id": s.Pilots[i].Id, "drones": drones})
	}
	return map[string]interface{}{
		"battlefield": map[string]interface{}{
			"officers": officers,
			"pilots":   pilots,
		},
		"request": map[string]interface{}{
			"user": map[string]string{
				"id":        authzDecision.user,
				"role":      authzDecision.role,
				"operation": authzDecision.operation,
			},
		},
	}
}

func (s *Server) authorized(ctx context.Context, operation string) (*AuthorizationDecision, error) {
	authzDecision := AuthorizationDecision{}
	authzDecision.allowed = false
	authzDecision.operation = operation
	jwtStr := ctx.Value("jwt").(string)
	if jwtStr == "" {
		log.Fatal("JWT received nil or emtpy from context")
	}

	authzDecision.uri = ctx.Value("uri").(string)
	if authzDecision.uri == "" {
		log.Fatal("URI received nil or emtpy from context")
	}

	xfcc := ctx.Value("xfcc").(*XFCC)

	token, err := jwt.ParseWithClaims(jwtStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		authzDecision.user = claims.ID
		if authzDecision.role, err = s.getRole(claims.ID); err != nil {
			return nil, err
		}
		clientRole := xfcc.GetClientRole()
		if authzDecision.role != *clientRole {
			log.Printf("User authenticated as '%s' but using certificate for '%s' detected, unauthorized.\n", authzDecision.role, *clientRole)
			authzDecision.allowed = false
			return &authzDecision, nil
		}
		data := s.populateAuthzData(ctx, &authzDecision)
		query, err := rego.New(
			rego.Query("data.battlefield.authz"),
			rego.Load([]string{"./access_policy.rego"}, nil),
		).PrepareForEval(ctx)

		if err != nil {
			log.Printf("Error parsing authorization policy: %s\n", err.Error())
			return nil, err
		}

		results, err := query.Eval(ctx, rego.EvalInput(data))
		if err != nil {
			log.Printf("Error evaluating policy query: %s\n", err.Error())
			return nil, err
		}
		authz := results[0].Expressions[0].Value.(map[string]interface{})

		log.Printf("Authorization result for user %s: %s", authzDecision.user, authz)

		authzDecision.allowed = authz["allow"].(bool)
		drones := authz["drones"].(map[string]interface{})
		log.Println(drones)
		for k, _ := range drones {
			authzDecision.drones = append(authzDecision.drones, k)
		}
	}
	return &authzDecision, nil
}
