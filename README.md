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
- OpenAPI 3.0 specification.
- Echo Web Framework.
- Golang.
