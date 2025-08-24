#!/bin/sh

set -e 
minikube start --driver=virtualbox

# We will install a custom ingress gateway
istioctl install --set profile=minimal

# We will use the new Kubernetes Gateway API
#kubectl get crd gateways.gateway.networking.k8s.io &> /dev/null ||
kubectl apply -f  "https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.3.0/standard-install.yaml"

# Enable istio-injection on default namespace
kubectl label namespace default istio-injection=enabled --overwrite

cd .. && make minikube && cd k8s

# Install application
kubectl apply -f drone-api.yaml

# Install kubernetes gateay and HTTProute
kubectl apply -f istio-tls-gw.yaml

cd certs && bash ./create_secret.sh && cd ..

# minikube tunnel
sleep 3
kubectl get pods -A
