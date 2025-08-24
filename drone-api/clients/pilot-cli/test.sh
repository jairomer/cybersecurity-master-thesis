#!/bin/bash

./pilot-cli \
	-ca ../../k8s/certs/gen/ca/ca.crt \
	-clientcert ../../k8s/certs/gen/officer-cert/officer.drone.api.crt \
	-clientkey ../../k8s/certs/gen/officer-cert/officer.drone.api.key \
	-apihost 10.100.11.31 \
	-pilotid "pilot-1" \
	-password "test12!"
	
