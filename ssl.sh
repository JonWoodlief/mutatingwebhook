#!/bin/bash

openssl genrsa -out ca.key 2048

openssl req -new -x509 -days 365 -key ca.key \
  -subj "/C=AU/CN=webhook-service.webhooktest.svc"\
  -out ca.crt

openssl req -newkey rsa:2048 -nodes -keyout server.key \
  -subj "/C=AU/CN=webhook-service.webhooktest.svc" \
  -out server.csr

openssl x509 -req \
  -extfile <(printf "subjectAltName=DNS:webhook-service.webhooktest.svc") \
  -days 365 \
  -in server.csr \
  -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out server.crt

echo
echo ">> Generating kube secrets..."
kubectl create secret tls -n webhooktest webhook-tls \
  --cert=server.crt \
  --key=server.key \
  --dry-run=client -o yaml \
  > ./kubeyamls/webhook.tls.secret.yaml

echo
echo ">> MutatingWebhookConfiguration caBundle:"
cat ca.crt | base64

rm ca.crt ca.key ca.srl server.crt server.csr server.key
