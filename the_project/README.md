# Todo App

It is a simple web server for the DevOps with Kubernetes course that displays an hourly cached picture along with your tasks. Now updated into a microservices architecture with a dedicated persistent database!


## What it does

- Starts a frontend web server, a separate backend server, and a reliable database.
- **Frontend (`todo-app`):** Serves the HTML, talks to the backend, and handles the image caching logic.
- **Backend (`todo-backend`):** Exposes `GET /todos` and `POST /todos` endpoints to fetch and save your tasks into the database.
- **Database (`todo-postgres`):** A PostgreSQL database managed by a StatefulSet to store your tasks securely.
- **Secure Configurations:** Uses Kubernetes ConfigMaps and Secrets to securely manage database credentials (host, port, user, database name, and password) without hardcoding them in the manifest files.
- **Persistent Cache:** Keeps the image identical for 10 minutes. If the container crashes or restarts, it reuses the cached image from the persistent volume instead of hitting the external API again.
- **Todo Interface:** Features a user input field with a strict 140-character maximum limit, a submit action button, and renders tasks dynamically from the backend service.
- **Automated Wikipedia Bot (`wikipedia-todo-job`):** A Kubernetes CronJob that triggers once every hour. It automatically fetches a random Wikipedia article URL and issues a `POST` request to the backend to add a "Read <URL>" reminder task.

## Run with Docker

```bash
docker build -t todo-app:1.0 ./frontend
docker build -t todo-backend:1.0 ./backend
```

## Run in Kubernetes

```bash
docker build -t todo-app:1.0 ./frontend
docker build -t todo-backend:1.0 ./backend
k3d image import todo-app:1.0 -c k3s-default
k3d image import todo-backend:1.0 -c k3s-default
kubectl apply -f manifests/db_config.yaml
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
kubectl port-forward service/todo-app-svc 8080:80 -n project
````
then got to
```bash
http://localhost:8080
```

## Logging
The backend includes structural input checks and logging explicitly tracked for observability stack monitoring

- Enforces a strict 140-character maximum limit on the backend. Requests exceeding this limit receive a `400 Bad Request`.
- Formats structured logs to standard output (`stdout`) for successful posts as well as input rejections.

Testing the logs:
```bash
kubectl port-forward service/todo-backend-svc 8082:80 -n project

curl -X POST http://localhost:8082/todos \
  -H "Content-Type: application/json" \
  -d '{"text": "LONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONG."}'

kubectl logs -l app=todo-backend -n project --tail=20
```

## Testing the Bot
```bash
kubectl create job --from=cronjob/wikipedia-todo-job test-wikipedia-job -n project
kubectl logs -l job-name=test-wikipedia-job -n project
kubectl delete job test-wikipedia-job -n project
```

## Namespace Separation
This application is deployed inside the isolated `project` namespace.
```bash
kubectl get all -n project
```
