#!/bin/bash

export DIR="./certificates"

# 1) root CA
openssl req -x509 -sha256 -days 365 \
  -nodes -newkey rsa:4096 \
  -subj "/O=Drone API CA/CN=drone-api-ca" \
  -keyout "$DIR/ca.key" -out "$DIR/ca.crt"

# 2) server cert for the ingress gateway (subject CN must match SNI / host)
openssl req -out "$DIR/server.csr" -newkey rsa:4096 -nodes \
  -keyout "$DIR/server.key" -subj "/CN=drone.api.xyz/O=uc3m"

openssl x509 -req -days 365 -sha256 \
  -in "$DIR/server.csr" -CA "$DIR/ca.crt" -CAkey "$DIR/ca.key" \
  -set_serial 01 -out "$DIR/server.crt"

# 3) Create client cert (signed by same CA)
entities=( "drone" "pilot" "officer" )
for entity in "${entities[@]}"
do
	openssl req -out "$DIR/$entity-client.csr" -newkey rsa:4096 -nodes \
  	  -keyout "$DIR/$entity-client.key" -subj "/CN=$entity.client.api/O=uc3m"

	openssl x509 -req -days 365 -sha256 \
  	  -in "$DIR/$entity-client.csr" -CA "$DIR/ca.crt" -CAkey "$DIR/ca.key" \
	  -set_serial 02 -out "$DIR/$entity-client.crt"

	echo "$entity certificate done: $entity-client.crt"
done

echo "Creating secrets..."

# TLS secret (server cert + key)
kubectl create -n istio-system secret tls server-credential \
	--cert="$DIR/server.crt" \
       	--key="$DIR/server.key"

# CA cert secret for client verification.
kubectl create -n istio-system secret generic clients-cacert \
       --from-file=cacert="$DIR/ca.crt"	

echo "Done"
