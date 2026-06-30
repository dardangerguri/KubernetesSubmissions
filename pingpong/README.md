# Pingpong App

It is a simple Go application that tracks request counts and exposes them both publicly to users and internally to other cluster applications over HTTP.


## Endpoint

- **`/pingpong`** (Public via Ingress): Increments and returns `pong N` where N increases on each request.
- **`/pings`** (Internal Cluster IP): Returns just the raw count integer `N` for the `log-output` service to consume.

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
