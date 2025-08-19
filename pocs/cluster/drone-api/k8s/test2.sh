#!/bin/bash

export INGRESS_HOST=$(kubectl get gtw httpbin-gateway2 -o jsonpath='{.status.addresses[0].value}')
export INGRESS_PORT=$(kubectl get gtw httpbin-gateway2 -o jsonpath='{.spec.listeners[?(@.name=="http")].port}')

# Positive test:
curl -s -I  "http://$INGRESS_HOST:$INGRESS_PORT/headers"

# Negative test:
#
# Should fail because this path is not explicitly exposed
	curl -s -I "http://$INGRESS_HOST:$INGRESS_PORT/status/200"

