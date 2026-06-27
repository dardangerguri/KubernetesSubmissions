# Pingpong App

It is a simple Go application that responds to GET requests and counts requests in memory.


## Endpoint

- `/pingpong` → returns `pong N` where N increases on each request

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
docker build -t pingpong:1.0 .
docker run -p 8080:8080 pingpong:1.0
```

## Run in Kubernetes
```bash
kubectl apply -f manifests/
kubectl get pods
kubectl get svc
kubectl get ingress
```

## Access in Kubernetes (Ingress)
The application is exposed using a Kubernetes Ingress.

```bash
http://localhost:8081/pingpong
```
