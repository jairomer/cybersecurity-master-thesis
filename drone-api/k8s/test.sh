#!/bin/bash
#
set -e
kubectl wait --for=condition=programmed gtw gateway -n drone-api-ingress
kubectl wait --for=condition=programmed gtw tls-gateway -n drone-api-ingress

export INGRESSDRONE_HOST=$(kubectl get gtw gateway -n drone-api-ingress -o jsonpath='{.status.addresses[0].value}')
export SECUREDRONE_HOST=$(kubectl get gtw tls-gateway -n drone-api-ingress -o jsonpath='{.status.addresses[0].value}')

echo "Ingress HTTP: $INGRESSDRONE_HOST"
echo "Ingress HTTPS: $SECUREDRONE_HOST"

# curl  -X POST -HHost:drone-api.com "http://$INGRESSDRONE_HOST/login" \
# 	-H "Content-Type: application/json"  -d '{"User":"officer-1","password":"changeme"}'
# curl  -X POST -HHost:drone-api.com "http://$INGRESSDRONE_HOST/battlefield"

#openssl s_client -connect $SECUREDRONE_HOST:443 -servername drone-api.com -showcerts < /dev/null
#openssl s_client -connect $SECUREDRONE_HOST:443 -servername drone-api.com </dev/null | \
#	openssl x509 -text -noout | grep -A1 "Subject Alternative Name"

curl -v --cacert certs/gen/ca/ca.crt \
	--resolve "api.drone.com:443:$SECUREDRONE_HOST" \
	-X POST  "https://api.drone.com:443/login" \
	-H "Content-Type: application/json" \
	-H "Host:officer.drone.com" \
	-d '{"User":"officer-1","password":"changeme"}'

#curl -s -k -X POST -HHost:drone-api.com "https://$SECUREDRONE_HOST:443/battlefield"

