#!/bin/bash
export SECUREDRONE_HOST=$(kubectl get gtw tls-gateway -n drone-api-ingress -o jsonpath='{.status.addresses[0].value}')
./pilot-cli \
	-ca ../../k8s/certs/gen/ca/ca.crt \
	-clientcert ../../k8s/certs/gen/pilot-cert/pilot.drone.api.crt \
	-clientkey ../../k8s/certs/gen/pilot-cert/pilot.drone.api.key \
	-apihost ${SECUREDRONE_HOST} \
	-pilotid "pilot-1" \
	-password "test12!"
	
