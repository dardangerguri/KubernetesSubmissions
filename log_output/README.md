# Log Output App

It is a multi-container Go application for the DevOps with Kubernetes course. Split into decoupled components running inside a single Pod.

## What it does

The application is split into two distinct components that communicate via a shared volume (`emptyDir`):
1. **Writer**: Generates a random UUID on startup and appends a line with a timestamp and the UUID every 5 seconds to a shared file.
2. **Reader**: An HTTP server that reads the shared log file and displays its contents to the user.

## Example Output
```bash
2020-03-30T12:15:17.705Z: 8523ecb1-c716-4cb6-a044-b9e83bb98e43
2020-03-30T12:15:22.705Z: 8523ecb1-c716-4cb6-a044-b9e83bb98e43
```

## Build and Push to Cluster

#### Build Docker images
```bash
docker build -t log-writer:1.0 -f writer/Dockerfile .
docker build -t log-reader:1.0 -f reader/Dockerfile .
```
#### Import images into the k3d cluster nodes
```bash
k3d image import log-writer:1.0 log-reader:1.0 -c k3s-default
```

## Run in Kubernetes

```bash
kubectl apply -f manifests/
```

## Access in Kubernetes (Ingress)
The application is exposed using a Kubernetes Ingress.

After deploying to the cluster, it can be accessed at:
```bash
http://localhost:8081
```
