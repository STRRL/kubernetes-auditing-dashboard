# kubernetes-auditing-dashboard

Simple Dashboard for Kubernetes Auditing Events

## Start Development Environment with minikube

```bash
minikube start \
  --mount --mount-string ./script/kube-apiserver-config:/kube-apiserver-config \
  --wait apiserver && \
  minikube ssh sudo bash /kube-apiserver-config/patch-kube-apiserver-manifest.sh
```
