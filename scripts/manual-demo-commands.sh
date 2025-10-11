#!/bin/bash

# Manual Demo Commands
# Run these step by step to generate audit events

# Create namespace
kubectl create namespace lifecycle-test

# Create deployment
kubectl create deployment demo-app --image=nginx:1.24 --replicas=3 -n lifecycle-test

# Wait and scale
sleep 5
kubectl scale deployment demo-app --replicas=5 -n lifecycle-test

# Update image
sleep 3
kubectl set image deployment/demo-app nginx=nginx:1.25 -n lifecycle-test

# Patch deployment
sleep 5
kubectl patch deployment demo-app -n lifecycle-test -p '{"metadata":{"labels":{"updated":"true"}}}'

# Scale down
sleep 3
kubectl scale deployment demo-app --replicas=2 -n lifecycle-test

# Create ConfigMap
sleep 3
kubectl create configmap demo-config \
  --from-literal=app=demo \
  --from-literal=env=test \
  -n lifecycle-test

# Update ConfigMap
sleep 2
kubectl patch configmap demo-config -n lifecycle-test \
  -p '{"data":{"new-key":"new-value"}}'

# Create Service
sleep 2
kubectl expose deployment demo-app \
  --port=80 \
  --target-port=80 \
  --name=demo-service \
  -n lifecycle-test

# Delete a pod to trigger recreation
sleep 3
kubectl delete pod -n lifecycle-test -l app=demo-app --field-selector=status.phase=Running | head -1

# View all resources
sleep 5
kubectl get all,configmap -n lifecycle-test
