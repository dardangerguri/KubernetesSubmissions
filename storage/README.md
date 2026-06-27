# Storage Setup

This folder holds the backbone for how our apps share data. Since the `pingpong` app and the `log_output` app live in totally separate places, they use the files in here to talk to each other through a shared folder.

## What's inside?

* **`persistentvolume.yaml`**: This tells Kubernetes to look at a specific folder on your actual machine (`/tmp/k3d-shared-data`) and let the cluster nodes use it.
* **`persistentvolumeclaim.yaml`**: Think of this as a reservation ticket. It asks the cluster to grab 1GB of that local space so our application pods can securely plug into it.

## How to use it

You have to run this **before** you start up the apps, otherwise the containers will get stuck waiting for a place to save their files:

```bash
kubectl apply -f .
```