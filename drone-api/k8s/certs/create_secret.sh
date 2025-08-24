#!/bin/bash
#
# set -e

export DRONE_API_KEY="gen/drone-api/api.drone.com.key"
export DRONE_API_CSR="gen/drone-api/api.drone.com.csr"
export DRONE_API_CRT="gen/drone-api/api.drone.com.crt"
export K8S_NS="drone-api-ingress"
export CA_KEY="gen/ca/ca.key"
export CA_CRT="gen/ca/ca.crt"

# Create a secret for the ingress gateway
kubectl delete secret -n "${K8S_NS}" drone-api-credential --ignore-not-found
kubectl create -n "${K8S_NS}" secret tls drone-api-credential \
  --key="${DRONE_API_KEY}" --cert="${DRONE_API_CRT}"

# Create a secret for the CA in order to do mTLS
kubectl delete secret -n "${K8S_NS}" drone-api-credential-cacert --ignore-not-found
kubectl create -n "${K8S_NS}" secret tls drone-api-credential-cacert \
  --key="${CA_KEY}" --cert="${CA_CRT}" 
