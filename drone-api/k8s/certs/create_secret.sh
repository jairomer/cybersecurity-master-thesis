#!/bin/bash
#
# set -e

# Create a secret for the ingress gateway
kubectl create -n istio-ingress secret tls drone-api-credential \
  --key=gen/drone-api/drone.api.key \
  --cert=gen/drone-api/drone.api.crt 

# Create a secret for the CA in order to do mTLS
kubectl create -n istio-ingress secret tls drone-api-credential-cacert \
  --key=gen/ca/ca.key \
  --cert=gen/ca/ca.crt 
