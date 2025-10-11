# Kubernetes Audit Event Schema

Reference: https://kubernetes.io/docs/reference/config-api/apiserver-audit.v1/

## Overview

Kubernetes audit events capture detailed information about API server requests. Each event represents a single API request and its outcome.

## Event Structure (audit.v1.Event)

### Required Fields

#### `level` (string)
- Type: `Level` enum
- Description: Defines the audit logging detail level
- Values: `None`, `Metadata`, `Request`, `RequestResponse`
- Determines how much information is captured in the audit event

#### `auditID` (string)
- Type: `UID`
- Description: Unique identifier generated for each request
- Ensures each audit event can be uniquely tracked across multiple stage records

#### `stage` (string)
- Type: `Stage` enum
- Description: Indicates the request handling stage when the event was generated
- **Values:**
  - `RequestReceived` - Event generated as soon as the request is received
  - `ResponseStarted` - Once response headers are sent, but before response body is sent (only for long-running requests like watch)
  - `ResponseComplete` - Response body has been completed and no more bytes will be sent
  - `Panic` - Events generated when a panic occurred
- **Important:** A single API request can generate multiple audit events with different stages

#### `requestURI` (string)
- The original request URI sent by the client to the API server
- Example: `/api/v1/namespaces/default/pods/my-pod`, `/api/v1/namespaces/default/pods/my-pod/status`
- Provides context about the specific API endpoint accessed

#### `verb` (string)
- **Kubernetes verb** associated with the request for resource requests
- Common values: `get`, `list`, `watch`, `create`, `update`, `patch`, `delete`, `deletecollection`
- For non-resource requests: lowercase HTTP method (`get`, `post`, `put`, `delete`)
- **Important:** This is the actual operation verb, NOT transformed or mapped

#### `user` (object)
- Type: `UserInfo` (from authentication/v1)
- Contains authenticated user information
- **Fields:**
  - `username` (string) - User identifier
  - `uid` (string) - Unique user ID
  - `groups` ([]string) - Groups the user belongs to
  - `extra` (map[string][]string) - Additional user attributes

#### `requestReceivedTimestamp` (string)
- ISO 8601 timestamp when the apiserver received the request
- Format: `2025-10-08T03:59:44.674501Z`

#### `stageTimestamp` (string)
- ISO 8601 timestamp when the current audit event stage was generated
- Format: `2025-10-08T03:59:44.679756Z`

### Optional Fields

#### `objectRef` (object)
- Type: `ObjectReference`
- Contains reference details for the Kubernetes object being operated on
- **Fields:**
  - `resource` (string) - Resource type (e.g., `pods`, `deployments`, `services`)
  - `namespace` (string) - Namespace of the resource (empty for cluster-scoped resources)
  - `name` (string) - Name of the resource
  - `uid` (string) - UID of the resource
  - `apiGroup` (string) - API group (empty for core resources)
  - `apiVersion` (string) - API version (e.g., `v1`, `apps/v1`)
  - `resourceVersion` (string) - Resource version
  - `subresource` (string) - **Subresource name if accessing a subresource** (e.g., `status`, `binding`, `scale`, `logs`)

#### `requestObject` (object)
- Type: `Unknown` (runtime.Unknown)
- **The API object from the request**, in JSON format
- Captured **before** version conversion, defaulting, or admission processing
- Only logged at `Request` level and higher
- **Important:** For PATCH requests, this contains the patch document, not the full object
- May have a different `kind` than the main resource (e.g., `Binding` for pod binding operations)

#### `responseObject` (object)
- Type: `Unknown` (runtime.Unknown)
- **The API object returned in the response**, serialized as JSON
- Captured **after** conversion to the external API type
- Only logged at `RequestResponse` level
- **Important Cases:**
  - Success: Contains the actual resource object with `kind` matching the resource type
  - Error: Contains a `Status` object with error details
  - Subresource operations: May return different object types (e.g., `Status` for successful binding)

#### `responseStatus` (object)
- Type: `Status` (from meta/v1)
- HTTP response status
- **Fields:**
  - `code` (int) - HTTP status code (200, 201, 404, 500, etc.)
  - `status` (string) - Status string (`Success`, `Failure`)
  - `message` (string) - Human-readable description
  - `reason` (string) - Machine-readable reason code

#### `sourceIPs` ([]string)
- Source IP addresses of the request
- Includes both forwarded IPs (from X-Forwarded-For headers) and the direct connection address
- Example: `["192.168.49.2"]`

#### `userAgent` (string)
- User agent string from the HTTP request
- Examples:
  - `kube-scheduler/v1.34.0 (linux/arm64) kubernetes/f28b4c9/scheduler`
  - `kubelet/v1.34.0 (linux/arm64) kubernetes/f28b4c9`
  - `kubectl/v1.34.0 (darwin/amd64) kubernetes/f28b4c9`

#### `annotations` (map[string]string)
- Additional metadata about the audit event
- Can be set by admission controllers, webhooks, or other components in the request serving chain
- Common annotations:
  - `authorization.k8s.io/decision` - Authorization decision (`allow`, `deny`)
  - `authorization.k8s.io/reason` - Why the request was allowed/denied

## Important Concepts

### Subresources

Many Kubernetes resources have **subresources** that are accessed as separate API endpoints:

- **Pods:**
  - `/status` - Pod status updates (typically by kubelet)
  - `/binding` - Pod-to-Node binding (by scheduler)
  - `/log` - Container logs
  - `/exec` - Execute commands in containers
  - `/portforward` - Port forwarding

- **Deployments/ReplicaSets:**
  - `/status` - Status updates
  - `/scale` - Scaling operations

**Key Point:** When `objectRef.subresource` is set, the event is about an operation on that subresource, NOT the main resource itself.

### Event Stages and Multiple Events

A single API request can generate **multiple audit events**:

1. `RequestReceived` - Always generated when request arrives
2. `ResponseStarted` - Only for long-running requests (watch, exec, logs)
3. `ResponseComplete` - Generated when response is complete
4. `Panic` - Only if the handler panics

**For lifecycle tracking, we typically use `ResponseComplete` stage** because:
- It includes both request and response objects (at `RequestResponse` level)
- The operation has completed (success or failure)
- We have the final state of the resource

### Understanding Different Verbs and Objects

#### CREATE with subresources
```json
{
  "verb": "create",
  "objectRef": {
    "resource": "pods",
    "subresource": "binding"
  },
  "requestObject": {
    "kind": "Binding"  // NOT Pod!
  },
  "responseObject": {
    "kind": "Status"  // Success response, not Pod
  }
}
```
This is **NOT** creating a Pod. It's creating a binding (scheduling the pod to a node).

#### PATCH on status subresource
```json
{
  "verb": "patch",
  "objectRef": {
    "resource": "pods",
    "subresource": "status"
  },
  "requestObject": {
    // Partial status update (patch document)
  },
  "responseObject": {
    "kind": "Pod"  // Full Pod object in response
  }
}
```
This is updating the Pod's status. The `responseObject` contains the **full Pod** after the patch.

#### True CREATE of a Pod
```json
{
  "verb": "create",
  "objectRef": {
    "resource": "pods",
    "subresource": ""  // No subresource!
  },
  "requestObject": {
    "kind": "Pod"  // Full Pod spec
  },
  "responseObject": {
    "kind": "Pod"  // Created Pod
  }
}
```

## Filtering for Resource Lifecycle

To track the **actual lifecycle** of a resource (not subresource operations):

### Correct Approach
```sql
SELECT * FROM audit_events
WHERE resource = 'pods'
  AND name = 'my-pod'
  AND verb IN ('create', 'update', 'patch', 'delete')
  AND stage = 'ResponseComplete'
  AND (
    -- Main resource operations (no subresource)
    subresource = '' OR subresource IS NULL

    -- OR subresource operations that return the full resource
    -- (e.g., status updates that include full Pod in response)
  )
```

### Checking Response Object Kind
- If `responseObject.kind` matches the resource kind → Valid lifecycle event
- If `responseObject.kind` = `"Status"` → May be error OR subresource success response
- Need to check both `requestObject.kind` and `responseObject.kind`

## Common Patterns

### Pod Lifecycle Events
1. **Creation** (by controller/user)
   - `verb=create`, `subresource=""`, `requestObject.kind=Pod`, `responseObject.kind=Pod`

2. **Binding** (by scheduler)
   - `verb=create`, `subresource=binding`, `requestObject.kind=Binding`, `responseObject.kind=Status`

3. **Status Updates** (by kubelet)
   - `verb=patch`, `subresource=status`, `responseObject.kind=Pod`

4. **Spec Updates** (by user/controller)
   - `verb=patch` or `update`, `subresource=""`, `responseObject.kind=Pod`

5. **Deletion**
   - `verb=delete`, `subresource=""`, `responseObject.kind=Pod` or `Status`

## Summary

- **Always check `objectRef.subresource`** to understand what is being operated on
- **Verify `requestObject.kind` and `responseObject.kind`** to ensure you're tracking the right object type
- **Use `stage=ResponseComplete`** for complete lifecycle tracking
- **verb is as-is from Kubernetes** - don't map or transform it
- **A single logical operation may involve multiple audit events** (e.g., create pod → multiple subresource updates)
