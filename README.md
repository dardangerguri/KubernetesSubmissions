# KubernetesSubmissions

## Exercises

### Chapter 2

- [1.1.](https://github.com/dardangerguri/KubernetesSubmissions/tree/1.1/log_output)
- [1.2.](https://github.com/dardangerguri/KubernetesSubmissions/tree/1.2/the_project)
- [1.3.](https://github.com/dardangerguri/KubernetesSubmissions/tree/1.3/log_output)
- [1.4.](https://github.com/dardangerguri/KubernetesSubmissions/tree/1.4/the_project)
- [1.5.](https://github.com/dardangerguri/KubernetesSubmissions/tree/1.5/the_project)
- [1.6.](https://github.com/dardangerguri/KubernetesSubmissions/tree/1.6/the_project)
- [1.7.](https://github.com/dardangerguri/KubernetesSubmissions/tree/1.7/log_output)
- [1.8.](https://github.com/dardangerguri/KubernetesSubmissions/tree/1.8/the_project)
- [1.9.](https://github.com/dardangerguri/KubernetesSubmissions/tree/1.9/pingpong)
- [1.10.](https://github.com/dardangerguri/KubernetesSubmissions/tree/1.10/log_output)
- [1.11.](https://github.com/dardangerguri/KubernetesSubmissions/tree/1.11)
- [1.12.](https://github.com/dardangerguri/KubernetesSubmissions/tree/1.12/the_project)
- [1.13.](https://github.com/dardangerguri/KubernetesSubmissions/tree/1.13/the_project)


### Chapter 3

- [2.1.](https://github.com/dardangerguri/KubernetesSubmissions/tree/2.1)
- [2.2.](https://github.com/dardangerguri/KubernetesSubmissions/tree/2.2/the_project)
- [2.3.](https://github.com/dardangerguri/KubernetesSubmissions/tree/2.3)
- [2.4.](https://github.com/dardangerguri/KubernetesSubmissions/tree/2.4)
- [2.5.](https://github.com/dardangerguri/KubernetesSubmissions/tree/2.5/log_output)
- [2.6.](https://github.com/dardangerguri/KubernetesSubmissions/tree/2.6/the_project)
- [2.7.](https://github.com/dardangerguri/KubernetesSubmissions/tree/2.7)
- [2.8.](https://github.com/dardangerguri/KubernetesSubmissions/tree/2.8/the_project)
- [2.9.](https://github.com/dardangerguri/KubernetesSubmissions/tree/2.9/the_project)
- [2.10.](https://github.com/dardangerguri/KubernetesSubmissions/tree/2.10/the_project)


### Chapter 4

- [3.1.](https://github.com/dardangerguri/KubernetesSubmissions/tree/3.1/pingpong)
- [3.2.](https://github.com/dardangerguri/KubernetesSubmissions/tree/3.2)
- [3.3.](https://github.com/dardangerguri/KubernetesSubmissions/tree/3.3)
- [3.4.](https://github.com/dardangerguri/KubernetesSubmissions/tree/3.4)
- [3.5.](https://github.com/dardangerguri/KubernetesSubmissions/tree/3.5)
- [3.6.](https://github.com/dardangerguri/KubernetesSubmissions/tree/3.6)
- [3.7.](https://github.com/dardangerguri/KubernetesSubmissions/tree/3.7)
- [3.8.](https://github.com/dardangerguri/KubernetesSubmissions/tree/3.8)
- [3.9.](https://github.com/dardangerguri/KubernetesSubmissions/tree/3.9)
- [3.10.](https://github.com/dardangerguri/KubernetesSubmissions/tree/3.10)
- [3.11.](https://github.com/dardangerguri/KubernetesSubmissions/tree/3.11)
- [3.12.](https://github.com/dardangerguri/KubernetesSubmissions/tree/3.12)


### Chapter 5

- [4.1.](https://github.com/dardangerguri/KubernetesSubmissions/tree/4.1)
- [4.2.](https://github.com/dardangerguri/KubernetesSubmissions/tree/4.2)
- [4.3.](https://github.com/dardangerguri/KubernetesSubmissions/tree/4.3)
- [4.4.](https://github.com/dardangerguri/KubernetesSubmissions/tree/4.4)
- [4.5.](https://github.com/dardangerguri/KubernetesSubmissions/tree/4.5)
- [4.6.](https://github.com/dardangerguri/KubernetesSubmissions/tree/4.6)
- [4.7.](https://github.com/dardangerguri/KubernetesSubmissions/tree/4.7)
- [4.8.](https://github.com/dardangerguri/KubernetesSubmissions/tree/4.8)
- [4.9.](https://github.com/dardangerguri/KubernetesSubmissions/tree/4.9)
- [4.10.](https://github.com/dardangerguri/KubernetesSubmissions/tree/4.10)

## Exercise 3.9: DBaaS vs DIY Database Comparison

A pros and cons of using Database-as-a-Service (DBaaS) like Google Cloud SQL versus a Do-It-Yourself (DIY) database running inside GKE (PostgreSQL StatefulSet).

### 1. Initialization (Work and Costs)
**DBaaS:**
- Very low work. It gets initialized with few clicks or a single command. The infrastructure and base engine are automatically provisioned.
- Higher base costs. It covers also a premium managed service wrapper on top of the compute and storage.

**DIY:**
- High work. Requires writing complex Kubernetes manifests (Statefulset, Services, PersistentVolumeClaims) and configuring storage classes.
- Low costs. It shares the compute resources of existing GKE cluster nodes, adding minimal costs for the raw persistent disk storage.

### 2. Maintenance
**DBaaS:**
- **Pros:** Hands-off. Cloud providers handle OS patching, database updates, version upgrades during maintenance windows. High availability replication can be turned with a toggle.
- **Cons:** Less control over specific low-level configuration parameters, extensions, and OS optimizations.

**DIY:**
- **Pros:** Control over database configurations, extensions, and OS optimizations.
- **Cons:** High work needed. You are responsible for updating the database version, patching security vulnerabilities.

### 3. Backup Methods and Ease of Usage
**DBaaS:**
- Easy to use. There are automated daily backups and point-in-time recovery are built-in features that can be enabled via cloud console. Restoring the database is simple.

**DIY:**
- Complex to use. Backups must be manually scripted and scheduled. Managing secure, automated retention policies and verifying snapshot integrity requires engineering efforts.
