build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .
	podman build -t hello .

test:
	go test .

clean:
	go clean .
	podman image rm --force hello

run:
	podman run  --network host -it localhost/hello:latest