validate:
	openapi-spec-validator openapi.yaml

build:
	# Generate server stubs
	oapi-codegen --config="./api/codegen.yaml" openapi.yaml
	# cli clients
	oapi-codegen --config=./clients/drone-cli/codegen.yaml openapi.yaml
	oapi-codegen --config=./clients/pilot-cli/codegen.yaml openapi.yaml
	oapi-codegen --config=./clients/officer-cli/codegen.yaml openapi.yaml
