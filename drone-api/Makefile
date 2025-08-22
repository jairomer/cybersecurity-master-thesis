validate:
	openapi-spec-validator openapi.yaml

gen:
	# Generate server stubs
	oapi-codegen --config="./api/codegen.yaml" openapi.yaml
	# cli clients
	oapi-codegen --config=./clients/drone-cli/codegen.yaml openapi.yaml
	oapi-codegen --config=./clients/pilot-cli/codegen.yaml openapi.yaml
	oapi-codegen --config=./clients/officer-cli/codegen.yaml openapi.yaml
	oapi-codegen --config=./clients/attacker-cli/codegen.yaml openapi.yaml

build: 
	cd api && go mod tidy && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build . 
	cd api && podman build -t drone-api .
	cd clients && cd drone-cli && go mod tidy && go build . 
	cd clients && cd pilot-cli && go mod tidy && go build . 
	cd clients && cd officer-cli && go mod tidy && go build . 
	cd clients && cd attacker-cli && go mod tidy && go build . 
