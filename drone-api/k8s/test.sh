#!/bin/bash
#
set -e
export INGRESS_HOST=$(kubectl get gtw tls-gateway -n istio-ingress -o jsonpath='{.status.addresses[0].value}')
export SECURE_INGRESS_PORT=$(kubectl get gtw tls-gateway -n istio-ingress -o jsonpath='{.spec.listeners[?(@.name=="https")].port}')

kubectl wait --for=condition=programmed gtw tls-gateway -n istio-ingress
