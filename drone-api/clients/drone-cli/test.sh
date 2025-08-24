#!/bin/bash

./drone-cli \
	-ca ../../k8s/certs/gen/ca/ca.crt \
	-clientcert ../../k8s/certs/gen/drone-cert/cli.drone.api.crt \
	-clientkey ../../k8s/certs/gen/drone-cert/cli.drone.api.key \
	-apihost 10.100.11.31
