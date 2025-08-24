#!/bin/bash
#
set -e
#kubectl wait --for=condition=programmed gtw gateway -n drone-api-ingress
kubectl wait --for=condition=programmed gtw tls-gateway -n drone-api-ingress

#export INGRESSDRONE_HOST=$(kubectl get gtw gateway -n drone-api-ingress -o jsonpath='{.status.addresses[0].value}')
export SECUREDRONE_HOST=$(kubectl get gtw tls-gateway -n drone-api-ingress -o jsonpath='{.status.addresses[0].value}')


# curl  -X POST -HHost:drone-api.com "http://$INGRESSDRONE_HOST/login" \
# 	-H "Content-Type: application/json"  -d '{"User":"officer-1","password":"changeme"}'
# curl  -X POST -HHost:drone-api.com "http://$INGRESSDRONE_HOST/battlefield"

# openssl s_client -connect $SECUREDRONE_HOST:443 -servername drone-api.com -showcerts < /dev/null
#openssl s_client -connect $SECUREDRONE_HOST:443 -servername drone-api.com </dev/null | \
#	openssl x509 -text -noout | grep -A1 "Subject Alternative Name"

export OFFICER_CLIENT_CRT="certs/gen/officer-cert/officer.drone.api.crt"
export OFFICER_CLIENT_KEY="certs/gen/officer-cert/officer.drone.api.key"
export CA_CRT="certs/gen/ca/ca.crt"
export OFFICER_HOST="officer.drone.com"

#curl -v --cacert ${CA_CRT} \
#	--resolve "api.drone.com:443:$SECUREDRONE_HOST" \
#	-X POST  "https://api.drone.com:443/login" \
#	-H "Content-Type: application/json" \
#	-H "Host:${OFFICER_HOST}" \
#	-d '{"User":"officer-1","password":"changeme"}'

curl -v --cacert ${CA_CRT} \
    --cert ${OFFICER_CLIENT_CRT} \
	--key ${OFFICER_CLIENT_KEY} \
	--resolve "api.drone.com:443:$SECUREDRONE_HOST" \
	-X POST  "https://api.drone.com:443/login" \
	-H "Content-Type: application/json" \
	-H "Host:${OFFICER_HOST}" \
	-d '{"User":"officer-1","password":"changeme"}'

echo "Ingress HTTPS: $SECUREDRONE_HOST"
#curl -s -k -X POST -HHost:drone-api.com "https://$SECUREDRONE_HOST:443/battlefield"

