# Service Account Webhook Demo

This project is a demonstration of how to force pods with a certain service account onto certain nodes.

## Build Docker image

```
docker build -t jonwoodlief/webhook:latest .
```

## prep cluster

generate webhooktest namespace using webhooktest-ns.yaml. label a node with the special label from main.go

## generate tls certs

run `./ssl.sh` and copy/paste the CA_BUNDLE outputted into webhook.yaml where it says TODO

## deploy all resources to k8s

deploy all resources to kubernetes, verify it works by making sure that the nginx pods are always scheduled onto a worker with the assigned labels
