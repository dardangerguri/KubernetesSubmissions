# Log Output App

It is a simple Go application for the DevOps with Kubernetes course.


## What it does

- Generates a random UUID on startup
- Stores it in memory
- Prints it every 5 seconds with a timestamp

## Example Output
```bash
2020-03-30T12:15:17.705Z: 8523ecb1-c716-4cb6-a044-b9e83bb98e43
2020-03-30T12:15:22.705Z: 8523ecb1-c716-4cb6-a044-b9e83bb98e43
```

## Run Locally

```bash
go run main.go
```

## Run with Docker

```bash
docker build -t log-output:1.0 .
docker run log-output:1.0
```

## Run in Kubernetes

```bash
kubectl apply -f manifests/deployment.yaml
kubectl logs -f <pod-name>
```
