# Tasks: Resource Lifecycle Viewer

**Input**: Design documents from `/Users/strrl/playground/GitHub/kubernetes-auditing-dashboard/specs/001-lifecycle-for-resources/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → ✅ Tech stack: Go 1.25.1, Next.js 13.3.0, Ent, gqlgen, TanStack Query
   → ✅ Structure: Web app (backend gql/, frontend ui/)
2. Load optional design documents:
   → ✅ data-model.md: Service models (ResourceIdentifier, LifecycleEvent, ResourceDiff)
   → ✅ contracts/: lifecycle-query.graphql contract
   → ✅ research.md: Technical decisions (GraphQL, diff algorithm, URL routing)
3. Generate tasks by category:
   → Setup: Composite index, dependencies
   → Tests: GraphQL resolver tests, diff service tests
   → Core: GraphQL schema, resolver, diff service, types
   → Frontend: Dynamic route, UI components, GraphQL integration
   → Polish: Linting, quickstart validation
4. Apply task rules:
   → Different files = [P] for parallel
   → Same file = sequential
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001-T023)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → ✅ Contract has corresponding test
   → ✅ Service models defined
   → ✅ All GraphQL types and resolver implemented
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Backend**: `gql/`, `pkg/services/`, `ent/schema/` at repository root
- **Frontend**: `ui/pages/`, `ui/modules/`, `ui/gql/` at repository root

---

## Phase 3.1: Setup

- [x] **T001** Add composite index to AuditEvent schema in `ent/schema/auditevent.go` for lifecycle query optimization (apiGroup, apiVersion, resource, namespace, name, requestTimestamp)

- [x] **T002** [P] Ensure YAML parsing dependency `gopkg.in/yaml.v3` is available in go.mod (likely already present, verify)

---

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3

**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

- [x] **T003** [P] GraphQL resolver test for `resourceLifecycle` query in `gql/lifecycle_test.go` - test cases:
  - Query with valid namespaced resource returns events in DESC order
  - Query with cluster-scoped resource (no namespace) works
  - Query with non-existent resource returns empty array
  - Query validates required parameters (kind, version, name not empty)
  - Query with 100+ events returns all events correctly (pagination test per Constitution II)
  - Verify events maintain correct DESC order across large result sets
  - Test memory efficiency with 500+ events (should not OOM)

- [x] **T004** [P] Diff service test in `pkg/services/lifecycle/diff_test.go` - test cases:
  - Compute diff between two YAML states shows added/removed/modified fields
  - Handle identical states (no diff)
  - Handle malformed YAML gracefully
  - Deep nested structure diffs computed correctly

- [x] **T005** [P] ResourceIdentifier parsing test in `pkg/services/lifecycle/types_test.go` - test cases:
  - Parse GVK from URL format "apps-v1-Deployment"
  - Handle core resources with empty apiGroup
  - URL decode special characters in resource names
  - Validate required fields (kind, version, name)
  - Parse namespaced resources (namespace != "_cluster")
  - Parse cluster-scoped resources (namespace == "_cluster", set to empty string)

---

## Phase 3.3: Backend Implementation (ONLY after tests are failing)

### GraphQL Schema & Types

- [x] **T006** Define lifecycle query schema in `gql/lifecycle.graphql`:
  - Add `resourceLifecycle` query with parameters (apiGroup, version, kind, namespace, name)
  - Define `LifecycleEvent` type with id, type, timestamp, user, resourceState, diff
  - Define `EventType` enum (CREATE, UPDATE, DELETE)
  - Define `ResourceDiff` type with added, removed, modified
  - Define `DiffEntry` type with path, oldValue, newValue

- [x] **T007** Run `make generate` to generate GraphQL types and resolver stubs from schema

### Service Layer Models

- [x] **T008** Implement all service types in `pkg/services/lifecycle/types.go`:
  - ResourceIdentifier struct with APIGroup, Version, Kind, Namespace, Name fields
  - `ParseFromURL(gvk, namespace, name string)` method (converts "_cluster" to empty namespace)
  - `ToEntQuery()` method to convert to Ent query parameters
  - Validation logic for required fields
  - LifecycleEvent struct with ID, Type, Timestamp, User, ResourceState, Diff
  - ResourceDiff struct with Added, Removed, Modified maps
  - DiffEntry struct with Path, OldValue, NewValue
  - EventType enum mapping (create → CREATE, update/patch → UPDATE, delete → DELETE)

- [x] **T009** [P] Implement YAML diff algorithm in `pkg/services/lifecycle/diff.go`:
  - `ComputeDiff(prevYAML, currentYAML string) (*ResourceDiff, error)` function
  - Parse YAML to map[string]interface{} using gopkg.in/yaml.v3
  - Deep recursive comparison for nested structures
  - Return only changed fields (added, removed, modified with paths)
  - Handle malformed YAML with error logging
  - Return partial diff with error indicator when YAML parse fails

- [x] **T010** [P] Implement error handling and recovery in `pkg/services/lifecycle/errors.go`:
  - Define custom error types:
    * `MalformedYAMLError` with partial parse result
    * `ResourceNotFoundError` for missing audit events
    * `InvalidIdentifierError` for bad GVK parsing
  - Error context structure with ResourceIdentifier, PartialData, RecoveryHint
  - Logging strategy:
    * ERROR level: Malformed YAML that blocks diff computation
    * WARN level: Partial YAML parse with recoverable data
    * INFO level: Empty result sets (not an error)
  - Recovery strategies for malformed YAML, missing fields, parse failures

### GraphQL Resolver

- [x] **T011** Implement `resourceLifecycle` resolver in `gql/lifecycle.resolvers.go`:
  - Parse and validate input parameters (ResourceIdentifier)
  - Query AuditEvent entities filtered by apiGroup, apiVersion, resource, namespace, name
  - Order results by requestTimestamp DESC
  - Map AuditEvent to LifecycleEvent (derive event type from verb)
  - Compute diffs for UPDATE events (compare consecutive raw YAML)
  - Extract user from userAgent field
  - Parse raw YAML to resourceState
  - Return empty array if no events found
  - Log errors with context (resource identifier, error details)

---

## Phase 3.4: Frontend Implementation

### GraphQL Client

- [x] **T012** Define GraphQL query in `ui/gql/lifecycle.graphql`:
  - `GetResourceLifecycle` query with variables (apiGroup, version, kind, namespace, name)
  - Request all LifecycleEvent fields (id, type, timestamp, user, resourceState, diff)
  - Include diff subfields (added, removed, modified with path, oldValue, newValue)

- [x] **T013** Run `cd ui && npm run codegen` to generate TypeScript types from GraphQL schema

### UI Components

- [x] **T014** [P] Create TimestampDisplay component in `ui/modules/common/TimestampDisplay.tsx`:
  - Accept ISO 8601 timestamp prop
  - Format using `Intl.DateTimeFormat` with browser timezone
  - Display human-readable format (e.g., "Oct 6, 2025, 3:30 PM")

- [x] **T015** [P] Create EmptyState component in `ui/modules/lifecycle/EmptyState.tsx`:
  - Display "No audit event record" message
  - Simple centered layout with appropriate styling (Tailwind)

- [x] **T016** [P] Create DiffViewer component in `ui/modules/lifecycle/DiffViewer.tsx`:
  - Accept ResourceDiffGql prop
  - Render modified fields with path, oldValue → newValue
  - Display added fields (green highlight)
  - Display removed fields (red highlight)
  - Format YAML-like output for readability

- [x] **T017** [P] Create TimelineView component in `ui/modules/lifecycle/TimelineView.tsx`:
  - Accept array of LifecycleEventGql props
  - Render events in vertical timeline layout
  - Each event shows: type badge, timestamp (using TimestampDisplay), user
  - For UPDATE events, include DiffViewer component
  - For CREATE/DELETE events, show resourceState without diff
  - Use Tailwind for timeline styling

### Dynamic Route

- [x] **T018** Create dynamic route page in `ui/pages/lifecycle/[...params].tsx`:
  - Parse URL params: params = [gvk, namespace, name] (always 3 segments)
  - Convert "_cluster" namespace to empty string for cluster-scoped resources
  - Parse GVK format "apiGroup-version-kind" (handle core resources with empty apiGroup)
  - URL decode resource name for special characters
  - Use TanStack Query to fetch data with GetResourceLifecycle query
  - Pass variables (apiGroup, version, kind, namespace, name) to GraphQL
  - Handle loading state (spinner or skeleton)
  - Handle error state (display error message)
  - Handle empty data (render EmptyState component)
  - Render TimelineView component with lifecycle events

---

## Phase 3.5: Integration & Testing

- [x] **T019** Run `make generate` to ensure all Ent and GraphQL codegen is up-to-date

- [x] **T020** Verify composite index created:
  - Check `ent/auditevent/auditevent.go` for generated index code
  - If needed, run database migration or regenerate SQLite schema

- [x] **T021** Run backend tests: `go test ./gql/... ./pkg/services/lifecycle/...` - ensure all tests pass

- [x] **T022** Run frontend linting: `cd ui && npm run lint` - fix any linting errors (lifecycle page errors fixed; pre-existing errors in other files remain)

- [x] **T022.1** Add navigation link to lifecycle viewer in Recent Changes page:
  - Locate the resource name display in the Recent Changes table component
  - Construct lifecycle URL: `/lifecycle/{gvk}/{namespace}/{name}` (use `_cluster` for cluster-scoped)
  - Wrap resource name with Next.js `<Link>` component pointing to lifecycle URL
  - Ensure GVK format matches: `{apiGroup}-{version}-{Kind}` (e.g., `apps-v1-Deployment`)
  - Handle namespace encoding: use actual namespace or `_cluster` sentinel
  - Apply appropriate hover styling to indicate clickable link

---

## Phase 3.6: Manual Validation

- [ ] **T023** Execute quickstart test scenarios from `quickstart.md`:
  - **Scenario 1**: Namespaced Deployment lifecycle with multiple UPDATE events (URL: `/lifecycle/apps-v1-Deployment/default/webapp`)
  - **Scenario 2**: Deleted ConfigMap showing DELETE event
  - **Scenario 3**: Cluster-scoped Namespace (URL: `/lifecycle/v1-Namespace/_cluster/default`)
  - **Scenario 4**: Non-existent resource (empty state)
  - **Scenario 5**: Resource with special characters in name
  - **Scenario 6**: Rapid updates with distinct timestamps
  - **Scenario 7**: Click resource name link from Recent Changes page, verify navigation to lifecycle page
  - Verify all pass criteria for each scenario
  - Take screenshots of UI for PR documentation

---

## Dependencies

**Setup (T001-T002)**:
- T001 must complete before T011 (resolver needs index for query performance)
- T002 verification can run in parallel

**TDD Phase (T003-T005)**:
- All tests [P] independent, can run in parallel
- Must ALL complete and FAIL before T006-T011

**Backend Implementation (T006-T011)**:
- T006 (schema) → T007 (codegen) - sequential dependency
- T008 (all types in types.go) must complete after T007
- T009 (diff.go) and T010 (errors.go) can run in parallel after T008 completes
- T011 (resolver) depends on T008, T009, T010 (uses all service types)

**Frontend (T012-T018)**:
- T012 (query definition) → T013 (codegen) - sequential
- T014, T015, T016, T017 can run in parallel (different component files)
- T018 (dynamic route) depends on T013 (types), T014-T017 (components)

**Integration & Testing (T019-T022.1)**:
- T019 first (ensure codegen current)
- T020-T022 can run in parallel
- T022.1 (navigation links) depends on identifying Recent Changes component location

**Manual Validation (T023)**:
- Depends on all previous tasks completing (including T022.1)
- Final validation before PR

---

## Parallel Execution Examples

### Example 1: TDD Phase Tests (T003-T005)
```bash
# Run all tests in parallel (they will fail initially - expected for TDD)
go test ./gql/lifecycle_test.go &
go test ./pkg/services/lifecycle/diff_test.go &
go test ./pkg/services/lifecycle/types_test.go &
wait
```

### Example 2: Service Layer Implementation (T008-T010)
After T007 codegen completes:
```
# First, implement all types in T008
Task: "Implement all service types in pkg/services/lifecycle/types.go"

# Then T009 and T010 can run in parallel
Task: "Implement YAML diff algorithm in pkg/services/lifecycle/diff.go"
Task: "Implement error handling and recovery in pkg/services/lifecycle/errors.go"
```

### Example 3: UI Components (T014-T017)
After T013 codegen, create components in parallel:
```
Task: "Create TimestampDisplay component in ui/modules/common/TimestampDisplay.tsx"
Task: "Create EmptyState component in ui/modules/lifecycle/EmptyState.tsx"
Task: "Create DiffViewer component in ui/modules/lifecycle/DiffViewer.tsx"
Task: "Create TimelineView component in ui/modules/lifecycle/TimelineView.tsx"
```

---

## Notes

- **[P] tasks** = different files, no dependencies
- **TDD Critical**: T003-T005 must be written first and must fail before implementation
- **Codegen steps**: T007 and T013 regenerate types, run after schema changes
- **Commit strategy**: Commit after each logical group (setup, tests, implementation phase)
- **Constitution compliance**: All tasks follow TDD principle, schema-first design, and structured logging

---

## Validation Checklist

*GATE: Checked before considering tasks complete*

- [x] Contract (lifecycle-query.graphql) has corresponding test (T003)
- [x] Service models (ResourceIdentifier, LifecycleEvent, ResourceDiff) have implementation tasks (T008-T010)
- [x] All tests come before implementation (T003-T005 before T006-T011)
- [x] Parallel tasks are truly independent (different files, verified)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task

---
