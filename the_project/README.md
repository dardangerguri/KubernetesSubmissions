# Todo App

It is a simple web server for the DevOps with Kubernetes course that displays an hourly cached picture along with your tasks.


## What it does

- Starts a web server on a configurable port
- Fetches a random 1200px image from Lorem Picsum and caches it locally.
- **Persistent Cache:** Keeps the image identical for 10 minutes. If the container crashes or restarts, it reuses the cached image from the persistent volume instead of hitting the external API again.
- **Todo Interface:** Features a user input field with a strict 140-character maximum limit, a submit action button, and a pre-seeded task list.

## Run Locally

```bash
go run main.go
```
## or with Custom Port

```bash
PORT=8080 go run main.go
```

## Run with Docker

```bash
docker build -t todo-app:1.0 .
docker run -e PORT=8080 todo-app:1.0
```

## Run in Kubernetes

```bash
kubectl apply -f storage/
docker build -t todo-app:1.0 .
k3d image import todo-app:1.0 -c k3s-default
kubectl apply -f manifests/
```

## Access
```bash
http://localhost:8081
```
