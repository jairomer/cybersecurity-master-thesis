#!/bin/sh

# We will install a custom ingress gateway
istioctl install --set profile=minimal

# We will use the new Kubernetes Gateway API 
#kubectl get crd gateways.gateway.networking.k8s.io &> /dev/null || 
kubectl apply -f "https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.3.0/standard-install.yaml"

# Enable istio-injection on default namespace
kubectl label namespace default istio-injection=enabled --overwrite

# Apply manifests for istio ingress gateway.
kubectl apply -f istio-ingress.yaml 

minikube tunnel
