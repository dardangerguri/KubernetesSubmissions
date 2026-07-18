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

## GitOps & Automated Deployment Pipeline

This application is fully managed using a declarative GitOps model driven by **ArgoCD** and **Kustomize**. Manual builds and explicit `kubectl apply` commands are no longer required to deploy updates.

### The Pipeline Architecture

1. **Code Modification**: A developer pushes updates to application logic or Kubernetes configuration files.
2. **CI/CD Automation (GitHub Actions)**:
   - Triggers automatically on file modifications inside `log_output/`.
   - Compiles Go binaries and builds secure, cross-platform container images.
   - Pushes images to Docker Hub, tagged uniquely using the modern git execution state `${{ github.sha }}`.
   - Executes `kustomize edit set image` to programmatically overwrite targets within `manifests/kustomization.yaml`.
   - Commits and pushes the modified configuration back into the active tracking repository branch.
3. **Continuous Convergence (ArgoCD)**:
   - Monitors the state specified in the repository.
   - Detects drift or state modifications automatically.
   - Directs the internal cluster resources to sync and securely pull matching images without exposed ingress cluster ports.

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
