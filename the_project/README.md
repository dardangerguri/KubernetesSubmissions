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

## Access in Kubernetes (Port Forward)
```bash
kubectl port-forward deployment/todo-app-dep 3000:3000
```
After port-forwarding, open:
```bash
http://localhost:3000
```

