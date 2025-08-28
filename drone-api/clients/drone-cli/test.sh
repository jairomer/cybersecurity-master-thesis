#!/bin/bash
export SECUREDRONE_HOST=$(kubectl get gtw tls-gateway -n drone-api-ingress -o jsonpath='{.status.addresses[0].value}')
./drone-cli \
	-ca ../../k8s/certs/gen/ca/ca.crt \
	-clientcert ../../k8s/certs/gen/drone-cert/cli.drone.api.crt \
	-clientkey ../../k8s/certs/gen/drone-cert/cli.drone.api.key \
	-apihost ${SECUREDRONE_HOST} \
	-droneid "drone-1" \
	-password "test12!"
