# Quickstart: Resource Lifecycle Viewer

## Prerequisites

- Minikube cluster running with audit webhook configured (see main README)
- Application running via `make dev` (backend on GraphQL, frontend on localhost:3000)
- Some Kubernetes resources created/modified to generate audit events

## Setup Test Data

### Step 1: Create Test Resources

```bash
# Create a test namespace
kubectl create namespace lifecycle-test

# Create a ConfigMap (cluster-scoped resource for comparison)
kubectl create configmap test-config \
  -n lifecycle-test \
  --from-literal=key1=value1

# Create a Deployment
kubectl create deployment nginx \
  -n lifecycle-test \
  --image=nginx:latest \
  --replicas=1
```

### Step 2: Update Test Resources

```bash
# Update ConfigMap (to generate UPDATE event)
kubectl patch configmap test-config \
  -n lifecycle-test \
  --type='json' \
  -p='[{"op": "add", "path": "/data/key2", "value": "value2"}]'

# Scale Deployment (to generate UPDATE event with visible diff)
kubectl scale deployment nginx \
  -n lifecycle-test \
  --replicas=3

# Update Deployment image (another UPDATE event)
kubectl set image deployment/nginx \
  -n lifecycle-test \
  nginx=nginx:1.25
```

### Step 3: Delete a Resource (Optional)

```bash
# Delete ConfigMap to generate DELETE event
kubectl delete configmap test-config -n lifecycle-test
```

Wait ~10 seconds for audit events to be ingested into the database.

## Test Scenarios

### Scenario 1: View Deployment Lifecycle (Namespaced Resource)

**Action**: Navigate to lifecycle page
```
http://localhost:3000/lifecycle/apps-v1-Deployment/lifecycle-test/nginx
```

**Expected Result**:
- Page loads successfully with timeline view
- Events displayed in reverse chronological order:
  1. UPDATE: Image change to nginx:1.25 (newest)
  2. UPDATE: Replica count change from 1 to 3
  3. CREATE: Initial deployment creation (oldest)
- Each UPDATE event shows a diff view with only changed fields:
  - Image change diff shows: `spec.template.spec.containers[0].image: nginx:latest → nginx:1.25`
  - Replica diff shows: `spec.replicas: 1 → 3`
- Timestamps are formatted in browser's local timezone (e.g., "Oct 6, 2025, 3:30 PM")
- User information displayed for each event (e.g., "kubectl/v1.28.0")

**Pass Criteria**:
- [x] Page renders without errors
- [x] All 3 events visible in timeline
- [x] Events ordered newest first (image change at top)
- [x] Diff view highlights only changed YAML fields
- [x] Timestamps readable in local timezone
- [x] No console errors in browser DevTools

### Scenario 2: View ConfigMap Lifecycle (Deleted Resource)

**Action**: Navigate to deleted ConfigMap lifecycle
```
http://localhost:3000/lifecycle/core-v1-ConfigMap/lifecycle-test/test-config
```

**Expected Result**:
- Page loads successfully
- Timeline shows:
  1. DELETE: Deletion event with final state (newest)
  2. UPDATE: Patch adding key2 field
  3. CREATE: Initial creation (oldest)
- DELETE event displays:
  - Event type badge shows "DELETE"
  - Final resource state visible (last known YAML before deletion)
  - No diff shown for DELETE event
- UPDATE event shows diff:
  - Added: `data.key2: "value2"`

**Pass Criteria**:
- [x] Deleted resource history still accessible
- [x] DELETE event clearly marked
- [x] Final resource state displayed
- [x] All 3 events visible
- [x] No diff shown for DELETE event

### Scenario 3: View Namespace Lifecycle (Cluster-Scoped Resource)

**Action**: Navigate to namespace lifecycle (no namespace in URL)
```
http://localhost:3000/lifecycle/core-v1-Namespace/lifecycle-test
```

**Expected Result**:
- Page loads successfully
- URL parsing correctly identifies cluster-scoped resource (2 path segments after /lifecycle/)
- Timeline shows:
  1. CREATE: Namespace creation event
- No UPDATE or DELETE events (namespace hasn't been modified)

**Pass Criteria**:
- [x] Cluster-scoped resource URL works (no namespace segment)
- [x] CREATE event displayed
- [x] No errors for resources with single event
- [x] Empty diff section for CREATE event

### Scenario 4: Non-Existent Resource (Empty State)

**Action**: Navigate to a resource that doesn't exist
```
http://localhost:3000/lifecycle/apps-v1-Deployment/default/nonexistent-app
```

**Expected Result**:
- Page loads successfully (not 404)
- Empty state component displayed with message: "No audit event record"
- No timeline or events shown
- No error messages in UI
- URL remains unchanged (user can bookmark and retry later)

**Pass Criteria**:
- [x] Empty state message displayed
- [x] No JavaScript errors
- [x] Page remains usable
- [x] "No audit event record" message visible

### Scenario 5: Resource with Special Characters

**Action**: Create and view a resource with special characters in name
```bash
# Create deployment with special name
kubectl create deployment test-app-v2.1 \
  -n lifecycle-test \
  --image=nginx:latest
```

Navigate to:
```
http://localhost:3000/lifecycle/apps-v1-Deployment/lifecycle-test/test-app-v2.1
```

**Expected Result**:
- URL encoding handled correctly (periods and hyphens)
- Resource lifecycle loads successfully
- CREATE event visible

**Pass Criteria**:
- [x] Special characters in URL work correctly
- [x] Resource identifier parsed accurately
- [x] Events displayed normally

### Scenario 6: Rapid Updates (Multiple Events Close Together)

**Action**: Make rapid changes to a resource
```bash
# Rapid scaling changes
kubectl scale deployment nginx -n lifecycle-test --replicas=5
sleep 1
kubectl scale deployment nginx -n lifecycle-test --replicas=7
sleep 1
kubectl scale deployment nginx -n lifecycle-test --replicas=10
```

Navigate to deployment lifecycle:
```
http://localhost:3000/lifecycle/apps-v1-Deployment/lifecycle-test/nginx
```

**Expected Result**:
- All 3 rapid UPDATE events visible
- Each event has distinct timestamp (even if milliseconds apart)
- Diffs show progression: 3→5, 5→7, 7→10
- No events merged or lost

**Pass Criteria**:
- [x] All rapid updates captured
- [x] Timestamps distinguish close events
- [x] Diffs computed correctly for consecutive changes

## GraphQL API Testing

### Direct GraphQL Query Test

**Action**: Open GraphQL Playground or use curl
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query { resourceLifecycle(apiGroup: \"apps\", version: \"v1\", kind: \"Deployment\", namespace: \"lifecycle-test\", name: \"nginx\") { id type timestamp user resourceState diff { modified { path oldValue newValue } } } }"
  }'
```

**Expected Result**:
```json
{
  "data": {
    "resourceLifecycle": [
      {
        "id": "...",
        "type": "UPDATE",
        "timestamp": "2025-10-06T...",
        "user": "kubectl/...",
        "resourceState": { ... },
        "diff": {
          "modified": [
            { "path": "spec.replicas", "oldValue": 7, "newValue": 10 }
          ]
        }
      },
      ...
    ]
  }
}
```

**Pass Criteria**:
- [x] GraphQL query executes successfully
- [x] Returns array of lifecycle events
- [x] Events include all required fields
- [x] Diff structure matches contract

## Performance Validation

### Step 1: Create Resource with Many Updates

```bash
# Script to generate 50 updates
for i in {1..50}; do
  kubectl scale deployment nginx -n lifecycle-test --replicas=$i
  sleep 0.5
done
```

### Step 2: Load Lifecycle Page

Navigate to:
```
http://localhost:3000/lifecycle/apps-v1-Deployment/lifecycle-test/nginx
```

**Expected Result**:
- Page loads in <2 seconds
- All 50+ events rendered in timeline
- Scrolling is smooth
- Diff computation completes without UI freeze
- Browser DevTools shows no performance warnings

**Pass Criteria**:
- [x] Page load time <2s
- [x] All events rendered
- [x] UI remains responsive
- [x] No memory leaks (check DevTools Performance tab)

## Cleanup

```bash
# Remove test namespace and all resources
kubectl delete namespace lifecycle-test
```

## Success Criteria Summary

**Must Pass**:
- ✅ All 6 test scenarios execute successfully
- ✅ GraphQL API returns correct data structure
- ✅ UI displays lifecycle events in reverse chronological order
- ✅ Diff views show only changed YAML fields
- ✅ Timestamps formatted in browser timezone
- ✅ Empty state handled gracefully
- ✅ No JavaScript errors in console
- ✅ Linting passes: `cd ui && npm run lint`

**Performance**:
- ✅ Query response time <500ms for <100 events
- ✅ Page load time <2s for resources with 50+ events
- ✅ Diff rendering completes without UI blocking

**Edge Cases**:
- ✅ Deleted resources show full history
- ✅ Cluster-scoped resources work without namespace
- ✅ Special characters in resource names handled
- ✅ Rapid updates all captured with distinct timestamps

---
