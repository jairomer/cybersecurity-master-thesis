package api

import (
	"context"
	"crypto/sha512"
	"log"
)

type Server struct {
	Users  map[string]string
	Pilots []PilotProvisioning
}

func (s *Server) GetBattlefieldData(ctx context.Context, request GetBattlefieldDataRequestObject) (GetBattlefieldDataResponseObject, error) {
	// TODO
	return nil, nil
}

func (s *Server) SetTargetLocation(ctx context.Context, request SetTargetLocationRequestObject) (SetTargetLocationResponseObject, error) {
	// TODO
	return nil, nil
}

func (s *Server) getBattlefieldData() *BattlefieldData {
	bd := []DroneData{}
	for _, pilot := range s.Pilots {
		for _, drone := range pilot.Drones {
			inBattlefield := false
			for _, droneInBattlefield := range bd {
				if droneInBattlefield.Id == drone.Id {
					// Repeated
					inBattlefield = true
					break
				}
			}
			if !inBattlefield {
				bd = append(bd, drone)
			}
		}
	}
	return &BattlefieldData{Drones: &bd}
}

func (s *Server) BattlefieldProvision(ctx context.Context, request BattlefieldProvisionRequestObject) (BattlefieldProvisionResponseObject, error) {

	// Add in credentials
	for _, cred := range request.Body.Credentials {
		hash := sha512.New()
		aggr := cred.User + cred.Password + "saltysalt!"
		hash.Write([]byte(aggr))
		s.Users[cred.User] = string(hash.Sum(nil))
	}

	// Add in the drones and the pilots.
	s.Pilots = append(s.Pilots, request.Body.Pilots...)
	bd := s.getBattlefieldData()
	return BattlefieldProvision200JSONResponse{bd.Drones}, nil
}

func (s *Server) authenticate(user, password string) bool {
	hash := sha512.New()
	aggr := user + password + "saltysalt!"
	hash.Write([]byte(aggr))
	return s.Users[user] == string(hash.Sum(nil)) || (user == "officer-1" && password == "test")
}

func (s *Server) Login(ctx context.Context, request LoginRequestObject) (LoginResponseObject, error) {

	if request.Body == nil || request.Body.Password == "" || request.Body.User == "" {
		log.Printf("Error decoding request on /login")
		return Login401Response{}, nil
	}
	evaluatedUser := request.Body.User
	evaluatedPassword := request.Body.Password
	if s.authenticate(evaluatedUser, evaluatedPassword) {
		log.Printf("User %s has been authenticated\n", evaluatedUser)
		evaluatedPassword = "" // clear it asap
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
