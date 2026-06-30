# Todo App

It is a simple web server for the DevOps with Kubernetes course that displays an hourly cached picture along with your tasks. Now updated into a microservices architecture!


## What it does

- Starts a frontend web server and a separate backend server.
- **Frontend (`todo-app`):** Serves the HTML, talks to the backend, and handles the image caching logic.
- **Backend (`todo-backend`):** Exposes `GET /todos` and `POST /todos` endpoints to save your tasks in memory.
- **Persistent Cache:** Keeps the image identical for 10 minutes. If the container crashes or restarts, it reuses the cached image from the persistent volume instead of hitting the external API again.
- **Todo Interface:** Features a user input field with a strict 140-character maximum limit, a submit action button, and renders tasks dynamically from the backend service.

## Run with Docker

```bash
docker build -t todo-app:1.0 ./frontend
docker build -t todo-backend:1.0 ./backend
```

## Run in Kubernetes

```bash
docker build -t todo-app:1.0 ./frontend
k3d image import todo-app:1.0 -c k3s-default
docker build -t todo-backend:1.0 ./backend
k3d image import todo-backend:1.0 -c k3s-default
kubectl apply -f storage/persistentvolumeclaim.yaml
kubectl apply -f manifests/
```

## Access
If Ingress port is mapped:
```bash
http://localhost:8081
```
or, use port-forwarding directly to the frontend service:
```bash
kubectl port-forward service/todo-app-svc 8080:80
````
then got to
```bash
http://localhost:8080
```
