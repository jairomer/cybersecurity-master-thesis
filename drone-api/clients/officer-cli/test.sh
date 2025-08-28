#!/bin/bash
set -x
export SECUREDRONE_HOST=$(kubectl get gtw tls-gateway -n drone-api-ingress -o jsonpath='{.status.addresses[0].value}')
./officer-cli \
	-ca ../../k8s/certs/gen/ca/ca.crt \
	-clientcert ../../k8s/certs/gen/officer-cert/officer.drone.api.crt \
	-clientkey ../../k8s/certs/gen/officer-cert/officer.drone.api.key \
	-apihost ${SECUREDRONE_HOST} \
	-officerid "officer-1" \
	-password "changeme" # Changed to test12! after first exec.
