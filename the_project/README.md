# Todo App

It is a simple web server for the DevOps with Kubernetes course that displays an hourly cached picture along with your tasks. Now updated into a microservices architecture with a dedicated persistent database!

## What it does

- Starts a frontend web server, a separate backend server, and a reliable database.
- **Frontend (`todo-app`):** Serves the HTML, talks to the backend, and handles the image caching logic. Updated with modern flexbox and defensive word-wrapping rules to keep buttons and task cards justified properly on the page.
- **Backend (`todo-backend`):** Exposes endpoints to fetch, save, and update your tasks into the database. GET requests are sorted explicitly by ID to preserve chronological list ordering.
- **Database (`todo-postgres`):** A PostgreSQL database managed by a StatefulSet to store your tasks securely. Updated to include a tracking column for completed tasks.
- **Secure Configurations:** Uses Kubernetes ConfigMaps and Secrets to securely manage database credentials (host, port, user, database name, and password) without hardcoding them in the manifest files.
- **Persistent Cache:** Keeps the image identical for 10 minutes. If the container crashes or restarts, it reuses the cached image from the persistent volume instead of hitting the external API again.
- **Todo Interface:** Features a user input field with a strict 140-character maximum limit, a submit action button, a **"Mark done"** toggle button, and renders tasks dynamically from the backend service.
- **Automated Wikipedia Bot (`wikipedia-todo-job`):** A Kubernetes CronJob that triggers once every hour. It automatically fetches a random Wikipedia article URL and issues a `POST` request to the backend to add a "Read <URL>" reminder task.

---

## Automated Database Backups (`postgres-backup-job`)
A secure Kubernetes CronJob that triggers once every 24 hours. It automatically dumps the PostgreSQL database contents and uploads the backup file securely to a Google Cloud Storage (GCS) bucket.

- **Security:** Uses GKE **Workload Identity Federation** to safely authenticate with Google Cloud without storing any hardcoded service account JSON keys inside the cluster or GitHub.
- **Tools:** Uses a lightweight `google/cloud-sdk:alpine` image to install native Postgres client utilities (`pg_dump`) on the fly, keeping the backup image minimal and up to date.
- **Storage Bucket:** Backups are saved with a timestamp format (`backup-YYYYMMDDHHMMSS.sql`) inside the `todo-db-backups-8802feef` GCS bucket.

### Testing the Backup Job Manually
To force an immediate database backup without waiting for midnight:
```bash
# Trigger a one-time manual execution from the CronJob
kubectl create job --from=cronjob/postgres-backup-job test-backup-run -n project

# View the execution and upload progress logs
kubectl logs -l job-name=test-backup-run -n project

# Verify the backup file exists in the cloud storage bucket
gcloud storage ls gs://todo-db-backups-8802feef/

# Clean up the manual test run
kubectl delete job test-backup-run -n project

## Run with Docker

```bash
docker build -t todo-app:1.0 ./frontend
docker build -t todo-backend:1.0 ./backend
```

## Run in Kubernetes (GKE with GitHub Actions & Kustomize)
This project uses **GitHub Actions** for automated CI/CD. Every push builds the backend and frontend Docker images, pushes them to **Google Artifact Registry**, and deploys them to the GKE cluster using **Kustomize**.


- Pushes to the `main` branch are deployed to the `project` namespace.
- Pushes to any other branch automatically create (if needed) and deploy to a namespace with the same name as the branch, providing an isolated preview environment.


To trigger a deployment:
```bash
git add .
git commit -m "Your deployment message"
git push origin main
```

## Access
If Ingress port is enabled:
```bash
http://<INGRESS-IP>
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
The production deployment runs inside the `project` namespace.
```bash
kubectl get all -n project
kubectl get all -n test-preview
```

## Resource Limits (No Cluster Crashing)

To stop the app from hogging all the cluster resources or crashing the nodes, everything has strict CPU and memory limits now:

- **Frontend & Backend:** Requests `50m` CPU / `64Mi` RAM, can burst up to `150m` CPU / `128Mi` RAM max.
- **Postgres DB:** Requests `100m` CPU / `128Mi` RAM, capped at `300m` CPU / `256Mi` RAM so database queries run smooth.

## GKE Cloud Logging & Monitoring

The cluster uses Google Cloud's native **Logs Explorer** to track what's happening in real-time. No need to spam `kubectl logs` all day—everything streams right to the GCP console.

To find the backend logs instantly, run this query in the Logs Explorer dashboard:

```query
resource.type="k8s_container"
resource.labels.namespace_name="project"
resource.labels.container_name="todo-backend"
```

## Prometheus Monitoring

To verify the monitoring setup, the Alertmanager StatefulSet was scaled to 3 replicas to test aggregation metrics:
```bash
kubectl scale statefulset prometheus-test-alertmanager --replicas=3 -n monitoring
```
The following PromQL query successfully filters and counts the 3 running StatefulSet pods:
```bash
sum(kube_pod_info{namespace="monitoring", created_by_kind="StatefulSet"})
```

## Health Checks

The backend exposes two internal endpoints:

- `GET /healthz` - Returns healthy only when the application and database are available.
- `POST /break` - Simulates an application failure. Kubernetes detects the failed liveness probe and automatically restarts the pod.
- `PUT /todos/:id` - Updates a specific task's done status in the database to toggle its completed state.

## Broadcaster Microservice & Message Broker (NATS)

The application has been upgraded with an event-driven architecture using **NATS** to handle live notifications whenever todos are added or completed.

- **Message Broker (`nats`):** A lightweight message queue deployed as a ClusterIP service inside the cluster.
- **Broadcaster (`todo-broadcaster`):** A dedicated microservice written in Go that listens to the `todos` subject on NATS and forwards the event message payload securely to an external chat webhook.
- **Horizontal Scaling & Deduplication:** Scaled to **6 replicas** to meet high availability requirements. It utilizes NATS **Queue Groups** (`QueueSubscribe` under the `todos-group` channel) ensuring that even with 6 active replicas, **only one instance** receives and broadcasts any given message, completely preventing duplicate notifications.

### Verifying the Broadcaster Logs & Replicas

```bash
kubectl get pods -l app=todo-broadcaster -n project

kubectl logs -l app=todo-broadcaster -n project --tail=20 -f
```

## GitOps Deployment

This project uses **Argo CD** for GitOps-based continuous deployment.

- **Target Namespace:** `project`
- **Sync Behavior:** Automatically tracks changes pushed to the `main` branch.
- **Application Directory:** `the_project/manifests`

### Manual Verification
To check the status of the GitOps synchronization directly from the cluster:
```bash
kubectl get application -n argocd todo-app-gitops
```