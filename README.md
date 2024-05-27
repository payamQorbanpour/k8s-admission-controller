# Kubernetes Admission Webhook

This repository contains a Kubernetes admission webhook example written in Go. The webhook validates admission requests and allows or denies the creation of Kubernetes resources based on custom logic.

## Contents

- `cmd/main.go`: The Go code for the webhook server.
- `example/`: Kubernetes manifests for the webhook server.
- `openssl.cnf`: OpenSSL configuration file for generating self-signed certificates.


## Getting Started

### Prerequisites

- Kubernetes cluster([kind](https://kind.sigs.k8s.io/))
- `kubectl` configured to interact with your cluster
- OpenSSL
- Docker
- Golang

### Generate TLS Certificates

Use the provided `openssl.cnf` to generate self-signed certificates:

```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout webhook-server.key -out webhook-server.crt -config openssl.cnf
```

Create a Kubernetes secret with the generated certificate and key:

```bash
kubectl create secret tls webhook-server-tls - cert=webhook-server.crt - key=webhook-server.key -n webhook
```

First, base64 encode the certificate:
```bash
cat webhook-server.crt | base64 | tr -d '\n'
```

Replace <base64-encoded-certificate> with the actual base64-encoded certificate in `example/validatingwebhookconfiguration.yaml` file:
```yaml
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
...
      caBundle: <base64-encoded-certificate>
...
```

### Build the project
Build the Docker image and push it to some registery.

```bash
docker buildx build --network host -t <webhook-server-image> . --push
```

Add the image to the deployment `example/deployment.yaml` file:
```yaml
apiVersion: apps/v1
kind: Deployment
...
 containers:
      - image: <webhook-server-image>
...
```

### Apply resources
```bash
kubectl apply -f example/service.yaml
kubectl apply -f example/deployment.yaml
kubectl apply -f example/validatingwebhookconfiguration.yaml
```