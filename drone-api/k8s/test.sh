#!/bin/bash
#
set -e
kubectl wait --for=condition=programmed gtw gateway -n drone-api-ingress

export INGRESSDRONE_HOST=$(kubectl get gtw gateway -n drone-api-ingress -o jsonpath='{.status.addresses[0].value}')
#export INGRESS_HOST=$(kubectl get gateways.gateway.networking.k8s.io gateway -n istio-ingress -ojsonpath='{.status.addresses[0].value}')
export SECURE_INGRESS_PORT=$(kubectl get gtw gateway -n drone-api-ingress -o jsonpath='{.spec.listeners[?(@.name=="https")].port}')

echo "Ingress Host: $INGRESSDRONE_HOST"
echo "Ingress Port: $INGRESSDRONE_PORT"

curl -v  -X POST -HHost:drone-api.com "http://$INGRESSDRONE_HOST/login" \
	-H "Content-Type: application/json"  -d '{"User":"officer-1","password":"changeme"}'

curl -s -X POST -HHost:drone-api.com "http://$INGRESSDRONE_HOST/battlefield"
# Another service
# curl -s -I -HHost:httpbin.example.com "http://$INGRESS_HOST/get"

