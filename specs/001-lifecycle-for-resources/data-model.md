# Data Model: Resource Lifecycle Viewer

## Existing Entities (No Changes Required)

### AuditEvent
**Source**: `ent/schema/auditevent.go` (already exists)

**Fields**:
- `id`: Int (auto-increment primary key)
- `raw`: Text (complete audit event JSON, immutable)
- `level`: String (audit level, immutable)
- `auditID`: String (unique audit identifier, immutable)
- `verb`: String (Kubernetes verb: create/update/patch/delete, immutable)
- `userAgent`: String (client user agent, immutable)
- `requestTimestamp`: Time (when request occurred, immutable, indexed)
- `stageTimestamp`: Time (audit stage timestamp, immutable, indexed)
- `namespace`: String (resource namespace, default empty for cluster-scoped)
- `name`: String (resource name)
- `apiVersion`: String (version only, e.g., "v1", "v1beta1" - NOT "apps/v1")
- `apiGroup`: String (API group like "apps", empty for core resources)
- `resource`: String (resource type like "deployments")
- `subResource`: String (sub-resource if applicable)
- `stage`: String (audit stage, immutable)

**Indexes**:
- Composite: (level, verb)
- Single: verb, auditID, userAgent, requestTimestamp, stageTimestamp

**Usage for Lifecycle Feature**:
- Query by: apiGroup, apiVersion, resource, namespace (optional), name
- Order by: requestTimestamp DESC (newest first)
- Event type derived from: verb field
- Diff source: raw field (parse YAML for comparison)

**Validation Rules**:
- All fields except namespace/name are required (NotEmpty)
- All fields are immutable after creation
- Timestamps must be valid RFC3339 format

## New Service Models (Backend Only)

### ResourceIdentifier
**Location**: `pkg/services/lifecycle/types.go`

**Purpose**: Parse and validate resource identifiers from URL paths

**Fields**:
- `APIGroup`: string (e.g., "apps", "" for core)
- `Version`: string (e.g., "v1")
- `Kind`: string (e.g., "Deployment")
- `Namespace`: string (optional, empty for cluster-scoped)
- `Name`: string

**Methods**:
- `ParseFromURL(gvk, namespace, name string) (*ResourceIdentifier, error)`: Parse URL segments
- `ToEntQuery() (apiGroup, apiVersion, resource, namespace, name string)`: Convert to Ent query params

**Validation**:
- Kind must not be empty
- Name must not be empty
- Version must not be empty
- APIGroup can be empty (core resources)
- Namespace validation based on resource scope (TBD: may need cluster-scoped resource list)

### LifecycleEvent
**Location**: `pkg/services/lifecycle/types.go`

**Purpose**: Represent processed lifecycle events for API responses

**Fields**:
- `ID`: int (AuditEvent ID)
- `Type`: EventType enum (CREATE, UPDATE, DELETE)
- `Timestamp`: time.Time
- `User`: string (extracted from userAgent or raw payload)
- `ResourceState`: map[string]interface{} (parsed YAML from raw field)
- `Diff`: *ResourceDiff (optional, only for UPDATE events)

**Derived From**: AuditEvent entity
- Type: derived from `verb` field
- Timestamp: uses `requestTimestamp`
- ResourceState: parsed from `raw` field YAML

### ResourceDiff
**Location**: `pkg/services/lifecycle/types.go`

**Purpose**: Represent changes between resource versions

**Fields**:
- `Added`: map[string]interface{} (fields added in this update)
- `Removed`: map[string]interface{} (fields removed in this update)
- `Modified`: map[string]DiffEntry (fields changed with old/new values)

**Subtype: DiffEntry**:
- `OldValue`: interface{}
- `NewValue`: interface{}
- `Path`: string (JSON path to the changed field, e.g., "spec.replicas")

**Computation**:
- Compare consecutive AuditEvent `raw` fields as YAML
- Deep diff recursive algorithm for nested structures
- Only include changed portions (per FR-009)

## GraphQL Schema Extensions

### Query Extension
**File**: `gql/auditevents.graphql` or new `gql/lifecycle.graphql`

```graphql
type Query {
  resourceLifecycle(
    apiGroup: String!
    version: String!
    kind: String!
    namespace: String
    name: String!
  ): [LifecycleEventGQL!]!
}
```

### New GraphQL Types
```graphql
type LifecycleEventGQL {
  id: ID!
  type: EventTypeEnum!
  timestamp: Time!
  user: String!
  resourceState: JSON!
  diff: ResourceDiffGQL
}

enum EventTypeEnum {
  CREATE
  UPDATE
  DELETE
}

type ResourceDiffGQL {
  added: JSON
  removed: JSON
  modified: [DiffEntryGQL!]!
}

type DiffEntryGQL {
  path: String!
  oldValue: JSON!
  newValue: JSON!
}
```

**Mapping**:
- `LifecycleEventGQL` ← `LifecycleEvent` service model
- `ResourceDiffGQL` ← `ResourceDiff` service model
- `JSON` scalar for dynamic YAML content (already defined in project)
- `Time` scalar for timestamps (already defined in `gql/time.graphql`)

## Frontend Data Flow

### GraphQL Operation
**File**: `ui/gql/lifecycle.graphql`

```graphql
query GetResourceLifecycle(
  $apiGroup: String!
  $version: String!
  $kind: String!
  $namespace: String
  $name: String!
) {
  resourceLifecycle(
    apiGroup: $apiGroup
    version: $version
    kind: $kind
    namespace: $namespace
    name: $name
  ) {
    id
    type
    timestamp
    user
    resourceState
    diff {
      added
      removed
      modified {
        path
        oldValue
        newValue
      }
    }
  }
}
```

### TypeScript Types (Generated)
**File**: `ui/gql/graphql.ts` (auto-generated by codegen)

```typescript
export type LifecycleEventGql = {
  id: string;
  type: EventTypeEnum;
  timestamp: string; // ISO 8601
  user: string;
  resourceState: any; // JSON
  diff?: ResourceDiffGql;
};

export enum EventTypeEnum {
  Create = 'CREATE',
  Update = 'UPDATE',
  Delete = 'DELETE'
}

export type ResourceDiffGql = {
  added?: any;
  removed?: any;
  modified: Array<DiffEntryGql>;
};

export type DiffEntryGql = {
  path: string;
  oldValue: any;
  newValue: any;
};
```

## State Transitions

### Resource Lifecycle States
```
[Not Exists]
    ↓ CREATE event
[Exists]
    ↓ UPDATE/PATCH events (0 or more)
[Exists (Modified)]
    ↓ DELETE event (optional)
[Deleted]
```

**Edge Cases**:
- Resource with only CREATE: Display single event, no diffs
- Resource with CREATE + DELETE (no updates): Show both events, deletion timestamp
- Multiple rapid UPDATEs: Each gets separate diff, ordered by timestamp
- PATCH vs UPDATE: Both treated as UPDATE event type for display purposes

## Query Performance Considerations

### Index Usage
Lifecycle query will use composite filtering:
```sql
SELECT * FROM audit_events
WHERE api_group = ?
  AND api_version = ?
  AND resource = ?
  AND namespace = ?
  AND name = ?
ORDER BY request_timestamp DESC
```

**Optimization Needed**: Add composite index for lifecycle queries:
```go
index.Fields("apiGroup", "apiVersion", "resource", "namespace", "name", "requestTimestamp")
```

### Diff Caching Strategy
- Compute diffs on-demand in resolver
- Future optimization: Cache computed diffs in memory with LRU eviction
- For MVP: No caching, acceptable for <100 events per resource

## Validation & Error Handling

### Input Validation
- Empty GVK components: Return GraphQL error "Invalid resource identifier"
- Special characters in name: URL-decode before query
- Non-existent resource: Return empty array (not error)

### Data Integrity
- Malformed YAML in `raw` field: Log error, show partial state with warning
- Missing verb field: Skip event, log warning
- Invalid timestamp: Use stageTimestamp as fallback

---
