# Research: Resource Lifecycle Viewer

## Decision 1: GraphQL Query Design

**Decision**: Create a new `resourceLifecycle` query that accepts GVK (group/version/kind), namespace, and resource name parameters, returning a list of AuditEvent objects ordered by timestamp.

**Rationale**:
- Existing AuditEvent entity already contains all necessary fields (namespace, name, apiGroup, apiVersion, resource, verb, requestTimestamp, raw payload)
- GraphQL provides type-safe querying and integrates with existing UI data fetching patterns (TanStack Query)
- Filtering by GVK + namespace + name can leverage existing indexes on the AuditEvent table

**Alternatives Considered**:
- REST endpoint: Rejected because project uses GraphQL throughout, maintaining consistency is important
- Direct database access from UI: Rejected, violates separation of concerns and existing architecture

## Decision 2: URL Route Pattern

**Decision**: Use Next.js dynamic catch-all route `/lifecycle/[...params].tsx` to handle both patterns:
- Namespaced: `/lifecycle/{apiGroup}-{version}-{kind}/{namespace}/{name}`
- Cluster-scoped: `/lifecycle/{apiGroup}-{version}-{kind}/{name}`

**Rationale**:
- Catch-all route `[...params]` captures variable-length path segments
- Can parse params array length to determine if namespace is present
- Aligns with Kubernetes resource identifier conventions
- Handles URL encoding for special characters in resource names

**Alternatives Considered**:
- Separate routes for namespaced vs cluster-scoped: Rejected, duplicates logic
- Query parameters for namespace: Rejected, less RESTful and harder to bookmark
- Single path with optional namespace marker: Rejected, ambiguous parsing

## Decision 3: YAML Diff Algorithm

**Decision**: Implement diff computation in Go backend (`pkg/services/lifecycle/diff.go`) using YAML parsing and field-level comparison. Return only changed fields as structured diff data to the UI.

**Rationale**:
- Go has robust YAML libraries (gopkg.in/yaml.v3) for parsing audit event payloads
- Computing diffs server-side reduces client bundle size and improves performance
- Structured diff data (not text) allows UI to render with semantic highlighting
- Backend can cache diff results for frequently accessed resources

**Alternatives Considered**:
- Client-side diff with JavaScript: Rejected, increases bundle size and complexity
- Text-based diff (unified format): Rejected, loses semantic structure for YAML
- Full object comparison: Rejected, requirement specifies "diff-only" view

## Decision 4: Timestamp Formatting

**Decision**: Send ISO 8601 timestamps from backend, format to browser timezone using JavaScript `Intl.DateTimeFormat` API in a reusable `TimestampDisplay` component.

**Rationale**:
- ISO 8601 is standard for timestamp transmission in APIs
- `Intl.DateTimeFormat` automatically handles browser timezone detection
- Centralized component ensures consistent formatting across the feature
- Existing audit events already use ISO timestamps (requestTimestamp field)

**Alternatives Considered**:
- Server-side timezone formatting: Rejected, cannot detect user's browser timezone
- Relative time only ("2 hours ago"): Rejected, requirement specifies human-readable timestamps
- Client-side timezone selection: Rejected, adds unnecessary complexity for MVP

## Decision 5: Event Type Detection

**Decision**: Derive lifecycle event type (create/update/delete) from the `verb` field in AuditEvent:
- `create` → CREATE event
- `update`, `patch` → UPDATE event
- `delete` → DELETE event

**Rationale**:
- AuditEvent schema already captures verb from Kubernetes audit logs
- Kubernetes audit verbs map directly to lifecycle event types
- No additional database fields or indexing required
- Verb field is already indexed for query performance

**Alternatives Considered**:
- Add explicit event_type field: Rejected, duplicates existing verb information
- Parse from raw JSON payload: Rejected, verb field already extracted for performance
- Infer from object state changes: Rejected, unnecessary complexity

## Decision 6: Empty State Handling

**Decision**: Return empty array from GraphQL query when no events found, UI renders `EmptyState` component with "No audit event record" message.

**Rationale**:
- Follows GraphQL best practice (null for errors, empty array for no results)
- Allows UI to distinguish between loading, error, and no-data states
- Consistent with TanStack Query patterns used in existing UI code
- Simple conditional rendering in React component

**Alternatives Considered**:
- Return null with error: Rejected, "no data" is not an error condition
- HTTP 404 response: Rejected, using GraphQL not REST
- Redirect to home page: Rejected, user should see explicit feedback

## Technology Stack Summary

**Backend**:
- GraphQL resolver in `gql/lifecycle.resolvers.go`
- Diff service in `pkg/services/lifecycle/diff.go` using `gopkg.in/yaml.v3`
- Existing Ent AuditEvent queries with filtering

**Frontend**:
- Next.js dynamic route: `ui/pages/lifecycle/[...params].tsx`
- React components: `TimelineView`, `DiffViewer`, `EmptyState`, `TimestampDisplay`
- TanStack Query for GraphQL data fetching
- TypeScript with generated GraphQL types

**Testing**:
- Go unit tests: `gql/lifecycle_test.go` for resolver logic
- Go service tests: `pkg/services/lifecycle/diff_test.go` for diff algorithm
- Frontend linting: `npm run lint`
- Manual UI testing via quickstart scenarios

---
