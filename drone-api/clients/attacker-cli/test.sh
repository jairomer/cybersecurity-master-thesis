#!/bin/bash
export SECUREDRONE_HOST=$(kubectl get gtw tls-gateway -n drone-api-ingress -o jsonpath='{.status.addresses[0].value}')
./attacker-cli \
  -ca ../../k8s/certs/gen/ca/ca.crt \
  -officercert ../../k8s/certs/gen/officer-cert/officer.drone.api.crt \
  -officerkey ../../k8s/certs/gen/officer-cert/officer.drone.api.key \
  -pilotcert ../../k8s/certs/gen/pilot-cert/pilot.drone.api.crt \
  -pilotkey ../../k8s/certs/gen/pilot-cert/pilot.drone.api.key \
  -dronecert ../../k8s/certs/gen/drone-cert/cli.drone.api.crt \
  -dronekey ../../k8s/certs/gen/drone-cert/cli.drone.api.key \
  -apihost ${SECUREDRONE_HOST}

# -crack \
