# Todo App

It is a simple web server for the DevOps with Kubernetes course.


## What it does

- Starts a web server on a configurable port
- Outputs "Server started in port NNNN" on startup
- Responds with "Todo App is running!" at the root endpoint

## Example Output

```bash
Server started in port 3000
```

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
docker build -t todo-app:1.0 .
docker run -e PORT=3000 todo-app:1.0
```

## Run in Kubernetes

```bash
kubectl apply -f manifests/deployment.yaml
kubectl get pods
kubectl logs -f <pod-name>
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

After deploying to the cluster, it can be accessed at:
```bash
http://localhost:8081
```
