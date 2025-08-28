# Design and Development of a Zero-Trust REST API in a Cloud Native Environment

**Contents**
- `/drone-api` - Main project implementation.
  - `api/` - REST API to be deployed on the cluster.
  - `clients/` - Clients to interact with the API.
    - `attacker-cli/` - Application specific security testing tool
    - `drone-cli/` - Drone REST API Client
    - `officer-cli/` - Officer REST API Client 
    - `pilot-cli/` - Pilot REST API Client 
  - `k8s/` - Assets to be deployed on the cluster
    - `certs/` - Certificates and certificate generation scripts
    -  install.sh - Installation script
  - `Makefile`  - Application development workflow
  - `openapi.yaml`  - OpenAPI specification of the REST API
- `/pocs` - Proofs of concept made along the way.

## **What is this?**

These are the technical assets produced for during the research and development of my master thesis in cybersecurity.

It is a toy REST API that validates Zero Trust principles leveraging the following technologies:
- Kubernetes API Gateway.
- Istio Service Mesh for inter-cluster secure communications and ingress controller.
- Mutual TLS for service authentication against the API Gateway.
- Rego for request policy management.
- Standard JWT for user authentication.
- XFCC - X-Forwarded-Client-Cert
- OpenAPI 3.0 specification.
- Echo Web Framework.
- Golang.

## **How to execute it?**

You will need minikube and istio installed on your system, as well as anything else that the installation and test scripts need to run.

**Install the test environment**
`cd /drone-api/k8s && ./install`

**Run the testbed**
After the syste has finished installed, you will need to open up access to the cluster by running `minikube tunnel` in a terminal.

Then you will need to open up 4 terminals, one for each of the clients, then go to their respective directories `/drone-api/clients/<client>`, then execute `./test.sh` to start simulating traffic, the attacker client will execute an attack battery to test the access control policy of the system.
- Initialize the officer client first so that the API is provisioned for the pilot and the drone!


