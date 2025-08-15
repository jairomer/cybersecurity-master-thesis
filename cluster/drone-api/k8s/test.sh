#!/bin/bash

export INGRESS_HOST=$(kubectl get gtw httpbin-gateway -o jsonpath='{.status.addresses[0].value}')
export INGRESS_PORT=$(kubectl get gtw httpbin-gateway -o jsonpath='{.spec.listeners[?(@.name=="http")].port}')

# Positive test:
curl -s -I -H Host:httpbin.example.com  "http://$INGRESS_HOST:$INGRESS_PORT/status/200"

# Negative test:
#
# Should fail because this path is not explicitly exposed
curl -s -I -HHost:httpbin.example.com "http://$INGRESS_HOST:$INGRESS_PORT/headers"

# Should fail because the host header was not specified.
curl -s -I  "http://$INGRESS_HOST:$INGRESS_PORT/status/200"
