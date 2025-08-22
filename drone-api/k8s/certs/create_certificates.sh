#!/bin/bash
#
set -e

##########################################################################
## CA
##########################################################################
#
# Create Root certificate and private key to sign the certificates.
#
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 \
	-subj '/O=UC3M.es./CN=ca.com' \
	-keyout gen/ca/ca.key \
	-out gen/ca/ca.crt

##########################################################################
## DRONE API
##########################################################################
#
# Generate a certificate and private key for 'drone.api'
#
openssl req -out gen/drone-api/drone.api.csr \
	-newkey rsa:2048 -nodes -keyout gen/drone-api/drone.api.key \
       	-subj "/CN=drone.api.com/O=UC3M"

openssl x509 -req -sha256 -days 365 -CA gen/ca/ca.crt \
	-CAkey gen/ca/ca.key -set_serial 0 \
	-in gen/drone-api/drone.api.csr \
	-out gen/drone-api/drone.api.crt

##########################################################################
# Generate a DRONE-CLI-CLIENT certificate and private key
##########################################################################
openssl req -out gen/drone-cert/drone.cli.csr -newkey rsa:2048 \
	-nodes -keyout gen/drone-cert/drone.cli.key \
	-subj "/CN=client.drone.cli/O=UC3M"

openssl x509 -req -sha256 -days 365 -CA gen/ca/ca.crt \
	-CAkey gen/ca/ca.key -set_serial 1 \
	-in gen/drone-cert/drone.cli.csr \
	-out gen/drone-cert/drone.cli.crt

##########################################################################
# Generate a PILOT-CLI-CLIENT certificate and private key
##########################################################################
openssl req -out gen/pilot-cert/pilot.cli.csr -newkey rsa:2048 \
	-nodes -keyout gen/pilot-cert/pilot.cli.key \
	-subj "/CN=client.pilot.cli/O=UC3M"

openssl x509 -req -sha256 -days 365 -CA gen/ca/ca.crt \
	-CAkey gen/ca/ca.key -set_serial 1 \
	-in gen/pilot-cert/pilot.cli.csr \
	-out gen/pilot-cert/pilot.cli.crt

##########################################################################
# Generate a OFFICER-CLI-CLIENT certificate and private key
##########################################################################
openssl req -out gen/officer-cert/officer.cli.csr -newkey rsa:2048 \
	-nodes -keyout gen/officer-cert/officer.cli.key \
	-subj "/CN=client.officer.cli/O=UC3M"

openssl x509 -req -sha256 -days 365 -CA gen/ca/ca.crt \
	-CAkey gen/ca/ca.key -set_serial 1 \
	-in gen/officer-cert/officer.cli.csr \
	-out gen/officer-cert/officer.cli.crt

##########################################################################
# Generate a ATTACKER-CLI-CLIENT certificate and private key
# - It will not be signed by the internal CA.but by another one.
##########################################################################
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 \
	-subj '/O=ADVERSARIAL./CN=ca.com' \
	-keyout gen/attacker-cert/ca.key \
	-out gen/attacker-cert/ca.crt

openssl req -out gen/attacker-cert/attacker.cli.csr -newkey rsa:2048 \
	-nodes -keyout gen/attacker-cert/attacker.cli.key \
	-subj "/CN=client.attacker.cli/O=ADVERSARIAL"

openssl x509 -req -sha256 -days 365 -CA gen/attacker-cert/ca.crt \
	-CAkey gen/attacker-cert/ca.key -set_serial 0 \
	-in gen/attacker-cert/attacker.cli.csr \
	-out gen/attacker-cert/attacker.cli.crt
