# Log Output App

It is a multi-container Go application for the DevOps with Kubernetes course. Split into decoupled components running inside a single Pod.

## What it does

The application consists of two distinct components that run concurrently in a single pod:
1. **Writer**: Generates a random UUID on startup and appends a line with a timestamp and the UUID every 5 seconds to an internal Pod volume file (`/shared/log.txt`).
2. **Reader**: An HTTP server that reads the local log file, makes an **internal HTTP GET request** to the `pingpong-app-svc` cluster service to fetch the live pong counter, aggregates them, and outputs the unified state to the user over port `8080`.
3. **ConfigMap Integration**: Mounts a text file (`/config/information.txt`) containing static content and injects an environment variable (`MESSAGE`) directly from a cluster ConfigMap to customize the response.

## Example Output
```bash
2026-06-27T22:46:12.733096425Z: f35da301-bc6c-4eb6-b813-3fa6870a775e
2026-06-27T22:46:17.738108845Z: f35da301-bc6c-4eb6-b813-3fa6870a775e.Ping / Pongs: 3
```

## How to Configure the App

You don’t need to touch the Go code to change what the application prints! Everything is handled dynamically using a Kubernetes ConfigMap (`manifests/configmap.yaml`).

Inside that file, you can tweak two things:

* **`MESSAGE`** (Environment Variable): Change this value to update the text printed on the second line of the web response. It defaults to `"hello world"`.
* **`information.txt`** (File Mount): This injects static text into the container filesystem at `/config/information.txt`. The app reads this file live, and it currently defaults to `"this text is from file"`.

If you update either of these settings in your YAML, just re-apply the config and give the deployment a quick kick to load the changes instantly:

```bash
kubectl apply -f manifests/configmap.yaml
kubectl rollout restart deployment log-output-dep -n exercises
```

## Build and Push to Cluster

#### Build Docker images for GKE (amd64)
```bash
docker build --platform linux/amd64 -t dardangerguri/log-writer:1.0-amd64 -f log_output/writer/Dockerfile log_output/
docker build --platform linux/amd64 -t dardangerguri/log-reader:1.0-amd64 -f log_output/reader/Dockerfile log_output/
```

#### Push to Docker Hub
```bash
docker push dardangerguri/log-writer:1.0-amd64
docker push dardangerguri/log-reader:1.0-amd64
```

## Run in Kubernetes
```bash
kubectl apply -f manifests/
```

Now you can test it:
```bash
curl http://136.68.44.240/
```

## Namespace Separation
This application is deployed inside the isolated `exercises` namespace.
```bash
kubectl get all -n exercises
```
## Readiness Probe

The reader container exposes a `/ready` endpoint for Kubernetes.

- Returns **200 OK** when the Pingpong application can be reached.
- Returns **503 Service Unavailable** while the Pingpong application is unavailable.

This ensures the application only receives traffic after its dependency is ready.
