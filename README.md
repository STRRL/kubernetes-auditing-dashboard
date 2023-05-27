# kubernetes-auditing-dashboard

Simple Dashboard for Kubernetes Auditing Events

## Start Development Environment with minikube

Bootstrap a local Kubernetes cluster with minikube, and configure the auditing webhook.

```bash
minikube start \
  --mount --mount-string ./script/kube-apiserver-config:/kube-apiserver-config \
  --wait apiserver && \
  minikube ssh sudo bash /kube-apiserver-config/patch-kube-apiserver-manifest.sh
```
Then start the application:

```bash
make dev
```