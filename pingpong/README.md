# Pingpong App

It is a simple Go service running in Kubernetes. It tracks the number of requests and stores the counter in a PostgreSQL database.
It is deployed in the `exercises` namespace and uses a StatefulSet-based PostgreSQL database for persistent storage.


## Endpoint

- **`/pingpong`** (Public via Ingress): Increments a counter stored in PostgreSQL and returns `pong N`.
- **`/pings`** (Internal Cluster IP): Returns the current counter value stored in PostgreSQL.


## Run with Docker

```bash
docker build -t pingpong-app:1.0 -f pingpong/Dockerfile .
docker run -p 8080:8080 pingpong-app:1.0
```

## Run in Kubernetes
```bash
k3d image import pingpong-app:1.0 -c k3s-default
kubectl apply -f manifests/
```

Now you can test it:
```bash
curl http://localhost:8082/pingpong
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
