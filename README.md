# Service Account Webhook Demo

This project is a demonstration of how to force pods with a certain service account onto certain nodes.

## Build Docker image

    ```shell
    docker build -t jonwoodlief/webhook:latest .
    ```

## deploy all resources to k8s

deploy all resources to kubernetes, verify it works by making sure that the nginx pods are always scheduled onto a worker with the assigned labels