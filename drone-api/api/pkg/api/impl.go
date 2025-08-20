package api

import (
	"context"
	"crypto/sha512"
	"log"

	"github.com/golang-jwt/jwt/v4"
	"github.com/open-policy-agent/opa/v1/rego"
)

const (
	SALT = "saltysalt!"
)

type UserDatabase struct {
	Officers map[string]string
	Pilots   map[string]string
	Drones   map[string]string
}

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
	userid, authz_drones, err := s.authorized(ctx)
	if err != nil || userid == "" {
		log.Printf("Error while processing user: %s", err.Error())
		return GetBattlefieldData403Response{}, nil
	}
	if len(authz_drones) == 0 {
		log.Printf("Unauthorized attempt at getting battlefield data by user %s\n", userid)
		return GetBattlefieldData403Response{}, nil
	}
	// Get the data for these drones.
	dd := []DroneData{}
	for _, droneid := range authz_drones {
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

func (s *Server) SetTargetLocation(ctx context.Context, request SetTargetLocationRequestObject) (SetTargetLocationResponseObject, error) {
	userid, authz_drones, err := s.authorized(ctx)

	if err != nil || userid == "" {
		log.Printf("Unauthorized attempt at setting location for drone: %s\n", err.Error())
		return SetTargetLocation403Response{}, nil
	}
	if len(authz_drones) == 0 {
		log.Printf("Unauthorized attempt at setting location for drone '%s' by user '%s'\n", request.Droneid, userid)
		return SetTargetLocation403Response{}, nil
	}

	for _, droneid := range authz_drones {
		if droneid == request.Droneid {
			for i, _ := range s.Pilots {
				if s.Pilots[i].Id == userid {
					// user is a defined pilot
					for j, _ := range s.Pilots[i].Drones {
						if s.Pilots[i].Drones[j].Id == request.Droneid {
							s.Pilots[i].Drones[j].Location.Altitude = request.Body.Altitude
							s.Pilots[i].Drones[j].Location.Latitude = request.Body.Latitude
							s.Pilots[i].Drones[j].Location.Longitude = request.Body.Longitude
							return SetTargetLocation200JSONResponse(s.Pilots[i].Drones[j]), nil
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
			log.Fatal(err)
			return Login401Response{}, err
		}
		return Login200JSONResponse{Token: token}, nil
	}
	log.Printf("User %s has failed authentication\n", evaluatedUser)
	return Login401Response{}, nil
}

func (s *Server) populateAuthzData(ctx context.Context, uri string, claims *JWTClaims) interface{} {
	log.Printf("Authorizing: %s accessing %s\n", claims.ID, uri)
	return `
		{
	  "battlefield": {
	    "officers" : [
		    { "id": "officer-1" }
	    ],
	    "pilots": [
		{
		  "id": "pilot-1",
		  "drones": [ "drone-1", "drone-2", "drone-3" ]
		},
		{
		  "id": "pilot-2",
		  "drones": [ "drone-4", "drone-5", "drone-6"]
		}
	    ]
	  },
	  "request": {
	    "user": {
		"id": "officer-1",
		"role": "officer"
	    }
	  }
	}`
	//return map[string]interface{}{
	//	"battlefield": map[string]interface{}{

	//	},
	//}
	//	return map[string]interface{}{
	//		"access_control": map[string]interface{}{
	//			"users": map[string]interface{}{
	//				"test1": map[string]interface{}{
	//					"acl": []string{"/hello/world"},
	//				},
	//				"test2": map[string]interface{}{
	//					"acl": []string{},
	//				},
	//			},
	//		},
	//		"jwt": map[string]interface{}{
	//			"aud": claims.ID,
	//		},
	//		"uri": uri,
	//	}
}

func (s *Server) authorized(ctx context.Context) (string, []string, error) {
	jwtStr := ctx.Value("jwt").(string)
	uri := ctx.Value("uri").(string)

	if jwtStr == "" {
		log.Fatal("JWT received nil or emtpy from context")
	}
	if uri == "" {
		log.Fatal("URI received nil or emtpy from context")
	}

	token, err := jwt.ParseWithClaims(jwtStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return "", []string{}, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		userid := claims.ID
		data := s.populateAuthzData(ctx, uri, claims)

		query, err := rego.New(
			rego.Query("data.battlefield.authz"),
			// TODO: Use environment variable instead.
			rego.Load([]string{"./pkg/api/auth/access_policy.rego"}, nil),
		).PrepareForEval(ctx)

		if err != nil {
			return userid, []string{}, err
		}

		results, err := query.Eval(ctx, rego.EvalInput(data))
		if err != nil {
			return userid, []string{}, err
		}
		if len(results) > 0 && len(results[0].Expressions) > 0 {
			authz, ok := results[0].Expressions[0].Value.(map[string]interface{})
			allowed := authz["allow"].(bool)
			if ok && allowed {
				log.Println("Authorized")
				return userid, authz["drones"].([]string), nil
			} else if !ok && allowed {
				log.Printf("Not Ok, but allowed")
			} else if ok && !allowed {
				log.Printf("Not allowed")
			}
		}
		return userid, []string{}, nil
	}
	return "", []string{}, nil
}
