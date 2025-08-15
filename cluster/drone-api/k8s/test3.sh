#!/bin/bash
#
set -e
export INGRESS_HOST=$(kubectl get gtw tls-gateway -n istio-ingress -o jsonpath='{.status.addresses[0].value}')
export SECURE_INGRESS_PORT=$(kubectl get gtw tls-gateway -n istio-ingress -o jsonpath='{.spec.listeners[?(@.name=="https")].port}')

kubectl wait --for=condition=programmed gtw tls-gateway -n istio-ingress

curl -v -HHost:httpbin.example.com \
	--resolve "httpbin.example.com:$SECURE_INGRESS_PORT:$INGRESS_HOST" \
  	--cacert ./certs/example_certs1/example.com.crt \
	"https://httpbin.example.com:$SECURE_INGRESS_PORT/status/418"
# curl -v -HHost:httpbin.example.com \
#   	--cacert ./certs/example_certs1/example.com.crt \
# 	"https://$INGRESS_HOST:/$SECURE_INGRESS_PORT/status/418"

curl -v -HHost:httpbin.example.com --resolve "httpbin.example.com:$SECURE_INGRESS_PORT:$INGRESS_HOST" \
  --cacert certs/example_certs2/example.com.crt "https://httpbin.example.com:$SECURE_INGRESS_PORT/status/418"

