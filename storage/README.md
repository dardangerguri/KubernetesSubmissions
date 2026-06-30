# Storage Setup

This folder contains the Persistent Volume Claim (PVC) required to request and manage persistent storage for the **Todo App** component in the cluster.

## What's inside?

- **`persistentvolumeclaim.yaml`**: A storage reservation ticket that requests 1GB of space using the cluster's default `local-path` storage class, ensuring that the Todo application's data persists even if its pod restarts.

## How to use it

Apply this manifest before deploying the applications that depend on it:

```bash
kubectl apply -f persistentvolumeclaim.yaml
```