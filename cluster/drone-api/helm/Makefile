.PHONY: render

lint:
	helm lint .

render:
	helm template . > render/drone-api.yml && vim render/drone-api.yml

install:
	helm install drone-api --atomic .

upgrade:
	helm upgrade drone-api --atomic .

uninstall:
	helm uninstall drone-api 

add-istio:
	helm repo add istio https://istio-release.storage.googleapis.com/charts
	helm repo update

install-istio:
	helm install istio-base istio/base --atomic -n istio-system \
		--set defaultRevision=default --create-namespace --wait
	helm install istiod istio/istiod -n istio-system --wait
	kubectl create namespace istio-ingress
	helm install istio-ingress istio/gateway -n istio-ingress --wait

uninstall-istio:
	helm uninstall istio-ingress -n istio-ingress --wait || true
	kubectl delete namespace istio-system || true
	helm uninstall istiod -n istio-system --wait || true
	helm uninstall istio-base -n istio-system --wait || true

