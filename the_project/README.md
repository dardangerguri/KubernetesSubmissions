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
kubectl create job --from=cronjob/postgres-backup-job test-backup-run -n production

# View the execution and upload progress logs
kubectl logs -l job-name=test-backup-run -n production

# Verify the backup file exists in the cloud storage bucket
gcloud storage ls gs://todo-db-backups-8802feef/

# Clean up the manual test run
kubectl delete job test-backup-run -n production
```

## Run with Docker

```bash
docker build -t todo-app:1.0 ./frontend
docker build -t todo-backend:1.0 ./backend
```

## Run in Kubernetes (GKE with GitHub Actions & Kustomize)
> ⚠️ **GitOps Architecture Update (Exercise 4.10):** > To separate application source code from structural deployment manifests, all Kubernetes resource definitions and Kustomize environment overlays (`/the_project/manifests`) have been migrated to a dedicated configuration repository: [todo-app-gitops-config](https://github.com/dardangerguri/todo-app-gitops-config).

This project uses **GitHub Actions** for automated CI/CD. Every push builds the backend, frontend, and broadcaster Docker images, pushes them to **Google Artifact Registry**, and updates the base Kustomize specifications.

The infrastructure automatically reconciles across isolated environments using **Argo CD** based on the following Git triggers:

- **Staging Namespace (`staging`):** Reconciles automatically on every push to the `main` branch. Built for rapid iteration, Kustomize patches the environment to clear the broadcaster's `WEBHOOK_URL` (`value: ""`) so it strictly prints to standard container logs. The database backup CronJob is also pruned entirely from this environment.
- **Production Namespace (`production`):** Reconciles strictly when a version tag matching the pattern `v*` (e.g., `v1.0.0`) is pushed. This is the fully featured deployment layer where external webhook notifications and hourly database backups are kept active.

To trigger a staging deployment:
```bash
git add .
git commit -m "Your deployment message"
git push origin main
```
To trigger a production release:
```bash
git tag v1.0.0
git push origin v1.0.0
```

## Access
```bash
# For Staging
kubectl port-forward service/todo-app-svc 8080:80 -n staging

# For Production
kubectl port-forward service/todo-app-svc 8081:80 -n production
```

## Logging
The backend includes structural input checks and logging explicitly tracked for observability stack monitoring

- Enforces a strict 140-character maximum limit on the backend. Requests exceeding this limit receive a `400 Bad Request`.
- Formats structured logs to standard output (`stdout`) for successful posts as well as input rejections.

Testing the logs:
```bash
kubectl port-forward service/todo-backend-svc 8082:80 -n production

curl -X POST http://localhost:8082/todos \
  -H "Content-Type: application/json" \
  -d '{"text": "LONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONGLONG."}'

kubectl logs -l app=todo-backend -n production --tail=20
```

## Testing the Bot
```bash
# Trigger bot job manually in production
kubectl create job --from=cronjob/wikipedia-todo-job test-wikipedia-job -n production
kubectl logs -l job-name=test-wikipedia-job -n production
kubectl delete job test-wikipedia-job -n production
```

## Namespace Separation
The multi-environment setup separates environments neatly across their own designated namespaces.
```bash
kubectl get all -n staging

kubectl get all -n production
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
resource.labels.namespace_name="production"
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
# Check Staging Replicas and Logs
kubectl get pods -l app=todo-broadcaster -n staging
kubectl logs -l app=todo-broadcaster -n staging --tail=20 -f

# Check Production Replicas and Logs
kubectl get pods -l app=todo-broadcaster -n production
kubectl logs -l app=todo-broadcaster -n production --tail=20 -f
```

## GitOps Deployment
> ⚠️ **GitOps Architecture Update (Exercise 4.10):** > To separate application source code from structural deployment manifests, all Kubernetes resource definitions and Kustomize environment overlays (`/the_project/manifests`) have been migrated to a dedicated configuration repository: [todo-app-gitops-config](https://github.com/dardangerguri/todo-app-gitops-config).

This project uses **Argo CD** for automated GitOps continuous state synchronization.

- **Staging Monitor:** Tracks changes pushed directly to the `main` branch.
- **Production Monitor:** Tracks releases triggered via version tags (`v*`).

### Manual Verification
To check the running sync status of both environments directly from the cluster:
```bash
kubectl get application -n argocd todo-app-staging
kubectl get application -n argocd todo-app-production
```
