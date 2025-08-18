package api

import (
	"context"
	"log"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (Server) GetHelloWorld(ctx context.Context, request GetHelloWorldRequestObject) (GetHelloWorldResponseObject, error) {

	res := HelloWorldResponse{Hello: "Hello", Country: "World!"}
	if request.Params.Country != nil {
		res.Country = *request.Params.Country
	}

	return GetHelloWorld200JSONResponse(res), nil
}

func authenticate(user, password string) bool {
	testUsers := map[string]string{
		"test1": "test",
		"test2": "test",
	}
	return testUsers[user] == password
}

func (Server) Login(ctx context.Context, request LoginRequestObject) (LoginResponseObject, error) {

	if request.Body == nil || request.Body.Password == "" || request.Body.User == "" {
		log.Printf("Error decoding request on /login")
		return Login403Response{}, nil
	}
	evaluatedUser := request.Body.User
	evaluatedPassword := request.Body.Password
	if authenticate(evaluatedUser, evaluatedPassword) {
		log.Printf("User %s has been authenticated\n", evaluatedUser)
		evaluatedPassword = "" // clear it asap
		token, err := GenerateToken(request.Body.User)
		if err != nil {
			log.Fatal(err)
			return Login403Response{}, err
		}
		return Login200JSONResponse{Token: token}, nil
	}
	log.Printf("User %s has failed authentication\n", evaluatedUser)
	return Login403Response{}, nil
}
