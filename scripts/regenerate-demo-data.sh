#!/bin/bash

set -e

echo "=== Kubernetes Auditing Dashboard - Demo Data Generation ==="
echo ""

# Step 1: Backup and clean database
echo "Step 1: Cleaning database..."
cd /Users/strrl/playground/GitHub/kubernetes-auditing-dashboard
if [ -f data.db ]; then
    cp data.db data.db.backup.$(date +%s)
    rm data.db
    echo "✓ Database backed up and removed"
else
    echo "✓ No existing database found"
fi

# Step 2: Rebuild application
echo ""
echo "Step 2: Building application..."
go build -o kubernetes-auditing-dashboard ./cmd/kubernetes-auditing-dashboard/
echo "✓ Application built"

# Step 3: Start application in background
echo ""
echo "Step 3: Starting application..."
./kubernetes-auditing-dashboard > app.log 2>&1 &
APP_PID=$!
echo "✓ Application started (PID: $APP_PID)"
echo "  Waiting 3 seconds for initialization..."
sleep 3

# Step 4: Ensure namespace exists
echo ""
echo "Step 4: Preparing Kubernetes namespace..."
kubectl create namespace lifecycle-test --dry-run=client -o yaml | kubectl apply -f -
echo "✓ Namespace ready"

# Step 5: Generate demo data
echo ""
echo "Step 5: Generating demo data..."

# Clean up any existing demo resources
echo "  Cleaning up existing demo resources..."
kubectl delete deployment demo-app -n lifecycle-test --ignore-not-found=true
sleep 2

# Create initial deployment (nginx:1.24)
echo "  Creating deployment (nginx:1.24)..."
kubectl create deployment demo-app --image=nginx:1.24 --replicas=3 -n lifecycle-test
sleep 5

# Scale deployment
echo "  Scaling deployment to 5 replicas..."
kubectl scale deployment demo-app --replicas=5 -n lifecycle-test
sleep 3

# Update image
echo "  Updating image to nginx:1.25..."
kubectl set image deployment/demo-app nginx=nginx:1.25 -n lifecycle-test
sleep 5

# Patch deployment (add label)
echo "  Patching deployment (adding label)..."
kubectl patch deployment demo-app -n lifecycle-test -p '{"metadata":{"labels":{"updated":"true"}}}'
sleep 2

# Scale down
echo "  Scaling down to 2 replicas..."
kubectl scale deployment demo-app --replicas=2 -n lifecycle-test
sleep 3

# Create a ConfigMap
echo "  Creating ConfigMap..."
kubectl create configmap demo-config --from-literal=key1=value1 --from-literal=key2=value2 -n lifecycle-test
sleep 2

# Update ConfigMap
echo "  Updating ConfigMap..."
kubectl patch configmap demo-config -n lifecycle-test -p '{"data":{"key3":"value3"}}'
sleep 2

# Create a Service
echo "  Creating Service..."
kubectl expose deployment demo-app --port=80 --target-port=80 --name=demo-service -n lifecycle-test
sleep 2

# Delete one pod to trigger recreation
echo "  Deleting one pod to trigger recreation..."
POD_NAME=$(kubectl get pods -n lifecycle-test -l app=demo-app -o jsonpath='{.items[0].metadata.name}')
kubectl delete pod $POD_NAME -n lifecycle-test
sleep 5

echo ""
echo "=== Demo Data Generation Complete ==="
echo ""
echo "Generated resources:"
kubectl get all,configmap -n lifecycle-test
echo ""
echo "Application is running at: http://localhost:23333"
echo "Frontend should be at: http://localhost:3000"
echo ""
echo "To stop the application: kill $APP_PID"
echo "To view logs: tail -f app.log"
