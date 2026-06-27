# Pingpong App

It is a simple Go application that responds to GET requests and tracks request counts persistently using a shared storage volume.


## Endpoint

- `/pingpong` → returns `pong N` where N increases on each request (persisted to `/shared/pongs.txt`)

## Run Locally

```bash
go run main.go
```
## or with Custom Port

```bash
PORT=3000 go run main.go
```

## Run with Docker

```bash
docker build -t pingpong-app:1.0 -f pingpong/Dockerfile .
docker run -p 8080:8080 -v /tmp/k3d-shared-data:/shared pingpong-app:1.0
```

## Run in Kubernetes
```bash
k3d image import pingpong-app:1.0 -c k3s-default
kubectl apply -f storage/
kubectl apply -f pingpong/manifests/
```

## Access in Kubernetes (Port-Forward)
The application is exposed using a Kubernetes Ingress.

```bash
kubectl port-forward deployment/pingpong-app-dep 8082:8080
```
Now you can test it:
```bash
curl http://localhost:8082/pingpong
```
