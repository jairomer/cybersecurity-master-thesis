#!/bin/bash
set -e

rm -Rf gen
mkdir -p gen/ca
mkdir gen/drone-api
mkdir gen/drone-cert
mkdir gen/pilot-cert
mkdir gen/officer-cert
mkdir gen/attacker-cert

##########################################################################
## CA
### Create Root certificate and private key to sign the certificates.
##########################################################################
export CA_KEY="gen/ca/ca.key"
export CA_CRT="gen/ca/ca.crt"

openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 \
	-subj '/O=UC3M/CN=drone.com' -keyout ${CA_KEY}\
	-out ${CA_CRT}

##########################################################################
## DRONE API
##########################################################################
export DRONE_API_KEY="gen/drone-api/api.drone.com.key"
export DRONE_API_CSR="gen/drone-api/api.drone.com.csr"
export DRONE_API_CRT="gen/drone-api/api.drone.com.crt"
export DRONE_API_CONFIG="openssl.cnf"

openssl req -out ${DRONE_API_CSR} -newkey rsa:2048 \
	-nodes -keyout ${DRONE_API_KEY} \
	-config ${DRONE_API_CONFIG}

openssl x509 -req -days 365 -CA ${CA_CRT} \
	-CAkey ${CA_KEY} -set_serial 0 \
	-in ${DRONE_API_CSR} -out  ${DRONE_API_CRT} \
	-extfile ${DRONE_API_CONFIG} -extensions req_ext

##########################################################################
# Generate a DRONE-CLI-CLIENT certificate and private key
##########################################################################
export DRONE_CLIENT_KEY="gen/drone-cert/cli.drone.api.key"
export DRONE_CLIENT_CSR="gen/drone-cert/cli.drone.api.csr"
export DRONE_CLIENT_CRT="gen/drone-cert/cli.drone.api.crt"

openssl req -out ${DRONE_CLIENT_CSR} -newkey rsa:2048 -nodes \
	-keyout ${DRONE_CLIENT_KEY} -subj "/CN=cli.drone.api/O=UC3M"

openssl x509 -req -sha256 -days 365 -CA ${CA_CRT} \
	-CAkey ${CA_KEY} -set_serial 1 -in ${DRONE_CLIENT_CSR} \
	-out ${DRONE_CLIENT_CRT}

##########################################################################
# Generate a PILOT-CLI-CLIENT certificate and private key
##########################################################################
export PILOT_CLIENT_KEY="gen/pilot-cert/pilot.drone.api.key"
export PILOT_CLIENT_CSR="gen/pilot-cert/pilot.drone.api.csr"
export PILOT_CLIENT_CRT="gen/pilot-cert/pilot.drone.api.crt"

openssl req -out ${PILOT_CLIENT_CSR} -newkey rsa:2048 -nodes \
	-keyout ${PILOT_CLIENT_KEY} -subj "/CN=pilot.drone.api/O=UC3M"

openssl x509 -req -sha256 -days 365 -CA ${CA_CRT} \
	-CAkey ${CA_KEY} -set_serial 2 -in ${PILOT_CLIENT_CSR} \
	-out ${PILOT_CLIENT_CRT}

##########################################################################
# Generate a OFFICER-CLI-CLIENT certificate and private key
##########################################################################
export OFFICER_CLIENT_KEY="gen/officer-cert/officer.drone.api.key"
export OFFICER_CLIENT_CSR="gen/officer-cert/officer.drone.api.csr"
export OFFICER_CLIENT_CRT="gen/officer-cert/officer.drone.api.crt"

openssl req -out ${OFFICER_CLIENT_CSR} -newkey rsa:2048 -nodes \
	-keyout ${OFFICER_CLIENT_KEY} -subj "/CN=officer.drone.api/O=UC3M"

openssl x509 -req -sha256 -days 365 -CA ${CA_CRT} \
	-CAkey ${CA_KEY} -set_serial 3 -in ${OFFICER_CLIENT_CSR} \
	-out ${OFFICER_CLIENT_CRT}

##########################################################################
# Generate a ATTACKER-CLI-CLIENT certificate and private key
# - It will not be signed by the internal CA.but by another one.
##########################################################################
export ATTACKER_CA_CERT="gen/attacker-cert/attacker.ca.crt"
export ATTACKER_CA_KEY="gen/attacker-cert/attacker.ca.key"
export ATTACKER_CLIENT_KEY="gen/attacker-cert/attacker.drone.api.key"
export ATTACKER_CLIENT_CSR="gen/attacker-cert/attacker.drone.api.csr"
export ATTACKER_CLIENT_CRT="gen/attacker-cert/attacker.drone.api.crt"

openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 \
	-subj '/O=ADVERSARIAL./CN=ca.com' \
	-keyout ${ATTACKER_CA_KEY} \
	-out ${ATTACKER_CA_CERT}

openssl req -out ${ATTACKER_CLIENT_CSR} -newkey rsa:2048 \
	-nodes -keyout ${ATTACKER_CLIENT_KEY} \
	-subj "/CN=client.attacker.cli/O=ADVERSARIAL"

openssl x509 -req -sha256 -days 365 -CA ${ATTACKER_CA_CERT} \
	-CAkey ${ATTACKER_CA_KEY} -set_serial 0 \
	-in ${ATTACKER_CLIENT_CSR} \
	-out ${ATTACKER_CLIENT_CRT}
