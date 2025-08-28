#!/bin/sh

set -e
minikube start --driver=virtualbox
minikube addons enable metrics-server

# We will install a custom ingress gateway
istioctl install --set profile=minimal -f topology.yml

# We will use the new Kubernetes Gateway API
#kubectl get crd gateways.gateway.networking.k8s.io &> /dev/null ||
kubectl apply -f "https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.3.0/standard-install.yaml"

# Enable istio-injection on default namespace
kubectl label namespace default istio-injection=enabled --overwrite

# Load the REST API app image into the cluster.
cd .. && make minikube && cd k8s


# Install application
kubectl apply -f drone-api.yaml

# Install kubernetes gateay and HTTProute
kubectl apply -f istio-tls-gw.yaml

# Create the certificates and secret
cd certs && ./clean.sh && ./create_certificates.sh && ./create_secret.sh && cd ..

sleep 3
kubectl get pods -A
