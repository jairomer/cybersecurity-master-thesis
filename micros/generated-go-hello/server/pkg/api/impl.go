package api

import (
	"context"
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
