#!/bin/bash
# This script patches the kube-apiserver manifest, enable the auditing features
# install yq
export VERSION=v4.33.2 BINARY=yq_linux_amd64 && curl -L -o - https://github.com/mikefarah/yq/releases/download/${VERSION}/${BINARY}.tar.gz |\
  tar xz && mv ${BINARY} /usr/bin/yq

systemctl stop kubelet

yq -i '.spec.containers[0].command |= . + ["--audit-policy-file=/kube-apiserver-config/audit-policy.yaml", "--audit-webhook-config-file=/kube-apiserver-config/audit-webhook.kubeconfig", "--audit-webhook-initial-backoff=1s", "--audit-webhook-batch-max-wait=10s"]' /etc/kubernetes/manifests/kube-apiserver.yaml
yq -i '.spec.volumes |= . + [{"hostPath": {"path":"/kube-apiserver-config"}, "name": "kube-apiserver-config"}]' /etc/kubernetes/manifests/kube-apiserver.yaml
yq -i '.spec.containers[0].volumeMounts |= . + [{"mountPath": "/kube-apiserver-config", "name": "kube-apiserver-config", "readOnly": true}]' /etc/kubernetes/manifests/kube-apiserver.yaml

systemctl start kubelet
