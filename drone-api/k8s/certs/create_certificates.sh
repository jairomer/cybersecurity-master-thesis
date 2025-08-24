#!/bin/bash
#
set -e

##########################################################################
## CA
##########################################################################
#
# Create Root certificate and private key to sign the certificates.
#
#openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 \
#	-subj '/O=UC3M.es./CN=ca.com' \
#	-keyout gen/ca/ca.key \
#	-out gen/ca/ca.crt
export CA_KEY="gen/ca/ca.key"
export CA_CRT="gen/ca/ca.crt"

openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 \
	-subj '/O=UC3M/CN=drone.com' -keyout ${CA_KEY}\
	-out ${CA_CRT}

##########################################################################
## DRONE API
##########################################################################
#
# Generate a certificate and private key for 'drone.api'
#
#openssl req -out gen/drone-api/drone.api.csr \
#	-newkey rsa:2048 -nodes -keyout gen/drone-api/drone.api.key \
#       	-subj "/CN=drone-api.com/O=UC3M" -config openssl.cnf

export DRONE_API_KEY="gen/drone-api/api.drone.com.key"
export DRONE_API_CSR="gen/drone-api/api.drone.com.csr"
export DRONE_API_CRT="gen/drone-api/api.drone.com.crt"

openssl req -out ${DRONE_API_CSR} -newkey rsa:2048 \
	-nodes -keyout ${DRONE_API_KEY} \
	-subj "/CN=api.drone.com/O=UC3M"

openssl x509 -req -sha256 -days 365 -CA ${CA_CRT} \
	-CAkey ${CA_KEY} -set_serial 0 \
	-in ${DRONE_API_CSR} -out  ${DRONE_API_CRT}

#openssl req -new -nodes -keyout gen/drone-api/drone.api.key \
#  -out gen/drone-api/drone.api.csr \
#  -config openssl.cnf

#openssl x509 -req -in gen/drone-api/drone.api.csr \
#  -CA gen/ca/ca.crt -CAkey gen/ca/ca.key -CAcreateserial \
#  -out gen/drone-api/drone.api.crt \
#  -days 365 -sha256 \
#  -extensions req_ext -extfile openssl.cnf

# openssl x509 -req -sha256 -days 365 -CA gen/ca/ca.crt \
# 	-CAkey gen/ca/ca.key -set_serial 0 \
# 	-in gen/drone-api/drone.api.csr \
# 	-out gen/drone-api/drone.api.crt \
#        	-extensions req_ext -extfile openssl.cnf

##########################################################################
# Generate a DRONE-CLI-CLIENT certificate and private key
##########################################################################
#openssl req -out gen/drone-cert/drone.cli.csr -newkey rsa:2048 \
#	-nodes -keyout gen/drone-cert/cli.drone.com.key \
#	-subj "/CN=cli.drone.com/O=UC3M"
#openssl x509 -req -sha256 -days 365 -CA gen/ca/ca.crt \
#	-CAkey gen/ca/ca.key -set_serial 1 \
#	-in gen/drone-cert/cli.drone.com.csr \
#	-out gen/drone-cert/cli.drone.com.crt

##########################################################################
# Generate a PILOT-CLI-CLIENT certificate and private key
##########################################################################
#openssl req -out gen/pilot-cert/pilot.cli.csr -newkey rsa:2048 \
#	-nodes -keyout gen/pilot-cert/pilot.cli.key \
#	-subj "/CN=client.pilot.cli/O=UC3M"
#
#openssl x509 -req -sha256 -days 365 -CA gen/ca/ca.crt \
#	-CAkey gen/ca/ca.key -set_serial 1 \
#	-in gen/pilot-cert/pilot.cli.csr \
#	-out gen/pilot-cert/pilot.cli.crt
#
##########################################################################
# Generate a OFFICER-CLI-CLIENT certificate and private key
##########################################################################
#openssl req -out gen/officer-cert/officer.cli.csr -newkey rsa:2048 \
#	-nodes -keyout gen/officer-cert/officer.cli.key \
#	-subj "/CN=client.officer.cli/O=UC3M"
#
#openssl x509 -req -sha256 -days 365 -CA gen/ca/ca.crt \
#	-CAkey gen/ca/ca.key -set_serial 1 \
#	-in gen/officer-cert/officer.cli.csr \
#	-out gen/officer-cert/officer.cli.crt

export OFFICER_CLIENT_KEY="gen/officer-cert/officer.drone.api.key"
export OFFICER_CLIENT_CSR="gen/officer-cert/officer.drone.api.csr"
export OFFICER_CLIENT_CRT="gen/officer-cert/officer.drone.api.crt"
export OFFICER_CONFIG="gen/officer-cert/openssl.cnf"

# openssl req -out ${OFFICER_CLIENT_CSR} -newkey rsa:2048 -nodes \
# 	-keyout ${OFFICER_CLIENT_KEY} -subj "/CN=officer.drone.api/O=UC3M"
# 
# openssl x509 -req -sha256 -days 365 -CA ${CA_CRT} \
# 	-CAkey ${CA_KEY} -set_serial 1 -in ${OFFICER_CLIENT_CSR} \
# 	-out ${OFFICER_CLIENT_CRT}

openssl req -out ${OFFICER_CLIENT_CSR} -newkey rsa:2048 -nodes \
	-keyout ${OFFICER_CLIENT_KEY} -subj "/CN=api.drone.com/O=UC3M"

openssl req -sha256 -days 365 -CA ${CA_CRT} \
	-CAkey ${CA_KEY} -set_serial 1 -in ${OFFICER_CLIENT_CSR} \
	-out ${OFFICER_CLIENT_CRT} -config ${OFFICER_CONFIG} \
	-copy_extensions=copy

##########################################################################
# Generate a ATTACKER-CLI-CLIENT certificate and private key
# - It will not be signed by the internal CA.but by another one.
##########################################################################
#openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 \
#	-subj '/O=ADVERSARIAL./CN=ca.com' \
#	-keyout gen/attacker-cert/ca.key \
#	-out gen/attacker-cert/ca.crt
#
#openssl req -out gen/attacker-cert/attacker.cli.csr -newkey rsa:2048 \
#	-nodes -keyout gen/attacker-cert/attacker.cli.key \
#	-subj "/CN=client.attacker.cli/O=ADVERSARIAL"
#
#openssl x509 -req -sha256 -days 365 -CA gen/attacker-cert/ca.crt \
#	-CAkey gen/attacker-cert/ca.key -set_serial 0 \
#	-in gen/attacker-cert/attacker.cli.csr \
#	-out gen/attacker-cert/attacker.cli.crt
#