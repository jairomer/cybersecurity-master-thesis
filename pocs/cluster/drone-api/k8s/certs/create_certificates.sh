#!/bin/bash
#
set -e

## FIRST SET
#
# Create Root certificate and private key to sign the certificates.
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 \
	-subj '/O=example Inc./CN=example.com' \
	-keyout example_certs1/example.com.key \
	-out example_certs1/example.com.crt

# Generate a certificate and private key for 'httpbin.example.com'
openssl req -out example_certs1/httpbin.example.com.csr \
	-newkey rsa:2048 -nodes -keyout example_certs1/httpbin.example.com.key \
       	-subj "/CN=httpbin.example.com/O=httpbin organization"

openssl x509 -req -sha256 -days 365 -CA example_certs1/example.com.crt \
	-CAkey example_certs1/example.com.key -set_serial 0 \
	-in example_certs1/httpbin.example.com.csr \
	-out example_certs1/httpbin.example.com.crt

## SECOND SET
#
# Create another set of the same.
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 \
	-subj '/O=example Inc./CN=example.com' \
	-keyout example_certs2/example.com.key \
       	-out example_certs2/example.com.crt

openssl req -out example_certs2/httpbin.example.com.csr -newkey rsa:2048 \
       	-nodes -keyout example_certs2/httpbin.example.com.key \
	-subj "/CN=httpbin.example.com/O=httpbin organization"

openssl x509 -req -sha256 -days 365 -CA example_certs2/example.com.crt \
       -CAkey example_certs2/example.com.key -set_serial 0 \
       -in example_certs2/httpbin.example.com.csr \
       -out example_certs2/httpbin.example.com.crt

# Generate a certificate and a private key for 'helloworld.example.com'
openssl req -out example_certs1/helloworld.example.com.csr -newkey rsa:2048 \
       -nodes -keyout example_certs1/helloworld.example.com.key \
       -subj "/CN=helloworld.example.com/O=helloworld organization"

openssl x509 -req -sha256 -days 365 -CA example_certs1/example.com.crt \
	-CAkey example_certs1/example.com.key -set_serial 1 \
	-in example_certs1/helloworld.example.com.csr \
	-out example_certs1/helloworld.example.com.crt

# Generate a CLIENT certificate and private key
openssl req -out example_certs1/client.example.com.csr -newkey rsa:2048 \
	-nodes -keyout example_certs1/client.example.com.key \
	-subj "/CN=client.example.com/O=client organization"
openssl x509 -req -sha256 -days 365 -CA example_certs1/example.com.crt \
	-CAkey example_certs1/example.com.key -set_serial 1 \
	-in example_certs1/client.example.com.csr \
	-out example_certs1/client.example.com.crt

# Create a secret for the ingress gateway
kubectl create -n istio-ingress secret tls httpbin-credential \
  --key=example_certs1/httpbin.example.com.key \
  --cert=example_certs1/httpbin.example.com.crt

