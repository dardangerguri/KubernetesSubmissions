# Pingpong App

It is a simple Go service running in Kubernetes. It tracks the number of requests and stores the counter in a PostgreSQL database.
It is deployed in the `exercises` namespace and uses a StatefulSet-based PostgreSQL database for persistent storage.


## Endpoint

- **Public Traffic (`/pingpong`)**: Users access the application via `/pingpong`. The Kubernetes Gateway API uses a `URLRewrite` filter to rewrite the `/pingpong` prefix to `/` before forwarding the request to the application.
- **Internal Traffic (`/pings`)**: Returns the current counter value stored in PostgreSQL without incrementing it. This endpoint is intended for internal cluster communication.


## Run with Docker Locally

```bash
docker build -t pingpong-app:1.0 -f pingpong/Dockerfile .
docker run -p 8080:8080 pingpong-app:1.0
```

## Cloud Deployment (GKE)
The application is containerized and deployed to Google Kubernetes Engine (GKE).

1. Build and Push to Docker Hub
```bash
docker build --platform linux/amd64 -t dardangerguri/pingpong-app:2.1 -f pingpong/Dockerfile pingpong/
docker push dardangerguri/pingpong-app:2.1
```
2. Deploy to GKE Cluster
Ensure your kubectl context is pointed to your active GKE cluster, then apply the manifests:
```bash
kubectl create namespace exercises --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f pingpong/manifests/
```
3. Verify Public Traffic & Routing Rewrites
Traffic routing is handled using the modern Kubernetes Gateway API (`gke-l7-global-external-managed`) connected to a Google Cloud Layer 7 Load Balancer. Fetch the assigned public IP address from the Gateway:
```bash
kubectl get gateway shared-gateway -n exercises
```
Test the live public endpoint using the ADDRESS IP:
```bash
curl http://<GATEWAY-ADDRESS-IP>/pingpong
```

## Namespace Separation
This application is deployed inside the isolated `exercises` namespace.
```bash
kubectl get all -n exercises
```

## Postgres
- PostgreSQL is deployed as a StatefulSet with 1 replica.
- It uses its own PersistentVolumeClaim defined with the StatefulSet (volumeClaimTemplates)
- The database stores the Pingpong counter
- Data persists across pod restarts and rescheduling

## Readiness Probe

The application exposes a `/ready` endpoint for Kubernetes.

- Returns **200 OK** when the application has a working connection to PostgreSQL.
- Returns **503 Service Unavailable** while the database is unavailable.

This allows Kubernetes to keep the container running while delaying traffic until the database is ready.

## Canary Deployments (Argo Rollouts)

The application utilizes **Argo Rollouts** instead of standard Kubernetes Deployments to manage safe, automated canary releases.

- **Analysis Template**: Monitors total CPU usage rate (`container_cpu_usage_seconds_total`) across the entire `exercises` namespace during updates.
- **Automated Rollback**: If the total namespace CPU usage exceeds the threshold checked via Prometheus (`prometheus-test-server.monitoring`), the rollout automatically aborts and reverts to the previous stable version.

### Manage the Rollout
```bash
kubectl get analysisruns -n exercises
```