
# Implementation Plan: Resource Lifecycle Viewer

**Branch**: `001-lifecycle-for-resources` | **Date**: 2025-10-06 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/Users/strrl/playground/GitHub/kubernetes-auditing-dashboard/specs/001-lifecycle-for-resources/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from file system structure or context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code, or `AGENTS.md` for all other agents).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Create a resource lifecycle viewer that displays the complete history of Kubernetes resources through a URL-based interface. Users can view all create, update (with YAML diffs), and delete events for any resource by navigating to `/lifecycle/{gvk}/{namespace}/{name}` (namespaced) or `/lifecycle/{gvk}/_cluster/{name}` (cluster-scoped). The feature leverages existing audit event data to provide timeline visualization with browser-timezone formatted timestamps and diff-only YAML views. Navigation links from the Recent Changes page enable quick access to any resource's lifecycle history.

## Technical Context
**Language/Version**: Go 1.25.1 (backend), TypeScript/Next.js 13.3.0 (frontend)
**Primary Dependencies**: Ent v0.14.5 (ORM), gqlgen v0.17.81 (GraphQL), Gin v1.11.0 (HTTP), React 18.2.0, TanStack Query v4.29.1
**Storage**: SQLite (existing audit events database)
**Testing**: Go standard `testing` package, Next.js/ESLint for frontend
**Target Platform**: Web application (localhost:3000 development, Kubernetes cluster deployment)
**Project Type**: Web (Next.js frontend + Go GraphQL backend)
**Performance Goals**: Sub-second query response for lifecycle queries, smooth diff rendering for resources up to 10k lines
**Constraints**: Read-only feature (no data modification), must work with existing audit event schema, browser timezone handling for timestamps
**Scale/Scope**: Support resources with 100+ lifecycle events, handle GVK parsing for all Kubernetes resource types

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**I. Auditing Data Integrity**
- [x] Webhook data persistence is atomic (acknowledge only after successful DB write) - N/A: read-only feature
- [x] Event timestamps, user identities, and resource changes are immutable - YES: only reads existing audit data
- [x] No modifications to audit event data after ingestion - YES: purely a query/display feature

**II. Test-Driven Development**
- [x] Tests written before implementation (TDD cycle enforced) - YES: will write GraphQL query tests first
- [x] Go tests use standard `testing` package in `*_test.go` files - YES: follows existing patterns
- [x] GraphQL resolvers have pagination/query coverage in `gql/*_test.go` - YES: lifecycle query needs pagination tests
- [x] UI changes include lint checks (`npm run lint`) - YES: new Next.js page requires linting

**III. Code Generation & Schema-First Design**
- [x] Entity changes start with `ent/schema/*.go` modifications - N/A: no new entities, uses existing Event schema
- [x] GraphQL schemas defined before resolvers - YES: define lifecycle query schema first
- [x] `make generate` run before committing schema changes - YES: required for GraphQL codegen
- [x] No manual edits to generated files - YES: only edit schema and resolvers

**IV. Conventional Commits & PR Discipline**
- [x] Commit messages follow Conventional Commits format - YES: will use `feat: add resource lifecycle viewer`
- [x] PRs include test results and rationale - YES: will include test output and screenshots
- [x] UI changes include screenshots - YES: timeline and diff views need visual validation
- [x] Feature branch workflow followed - YES: already on branch `001-lifecycle-for-resources`

**V. Observability & Debugging**
- [x] Critical operations emit structured logs - YES: GraphQL resolver will log query params and errors
- [x] Errors include context (user, resource, timestamp) - YES: error responses include resource identifier
- [x] UI displays query errors clearly - YES: empty state and error boundaries
- [x] Local dev uses `make dev` for full-stack visibility - YES: standard development workflow

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
gql/
├── schema.graphql                # Add lifecycle query schema
├── lifecycle.resolvers.go        # New resolver for lifecycle queries
└── lifecycle_test.go             # GraphQL query tests (TDD)

ui/
├── pages/
│   └── lifecycle/
│       └── [...params].tsx       # Dynamic route: /lifecycle/{gvk}/{namespace}/{name} (use _cluster for cluster-scoped)
├── modules/
│   ├── lifecycle/
│   │   ├── TimelineView.tsx      # Event timeline component
│   │   ├── DiffViewer.tsx        # YAML diff component
│   │   └── EmptyState.tsx        # "No audit event record" component
│   └── common/
│       └── TimestampDisplay.tsx  # Browser-timezone timestamp formatter
└── gql/
    └── lifecycle.graphql         # GraphQL query definitions

pkg/services/
└── lifecycle/
    └── diff.go                   # YAML diff computation logic
```

**Structure Decision**: Web application structure using existing Next.js frontend (`ui/`) and Go GraphQL backend (`gql/`). The feature adds:
- GraphQL schema extension and resolver in `gql/`
- Next.js dynamic route in `ui/pages/lifecycle/`
- Reusable UI components in `ui/modules/lifecycle/`
- Diff computation service in `pkg/services/lifecycle/`

No new Ent entities required; queries use existing Event entity from audit webhook ingestion.

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `.specify/scripts/bash/update-agent-context.sh claude`
     **IMPORTANT**: Execute it exactly as specified above. Do not add or remove any arguments.
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs (contracts, data model, quickstart)
- GraphQL contract test → lifecycle query test in `gql/lifecycle_test.go` [P]
- Service layer tests → diff computation tests in `pkg/services/lifecycle/diff_test.go` [P]
- UI component tests → linting and manual validation via quickstart
- Implementation tasks follow TDD: make failing tests pass

**Ordering Strategy**:
1. **Setup Phase**: Add composite index to AuditEvent schema
2. **TDD Phase**: Write all tests first (GraphQL resolver test, diff service tests) [P]
3. **Backend Implementation**: GraphQL schema → resolver → diff service
4. **Frontend Implementation**: Dynamic route → UI components → GraphQL integration
5. **Polish Phase**: Linting, quickstart validation, screenshots

**Estimated Output**: 20-25 numbered, ordered tasks in tasks.md

**Key Parallel Opportunities**:
- GraphQL resolver test + diff service test can run in parallel (different files)
- UI component creation can parallelize (TimelineView, DiffViewer, EmptyState in separate files)
- GraphQL schema definition + TypeScript type generation can run sequentially (codegen dependency)

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [x] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented (none required)

**Artifacts Generated**:
- [x] research.md - Technical decisions and alternatives analysis
- [x] data-model.md - Entity definitions and GraphQL schema design
- [x] contracts/lifecycle-query.graphql - GraphQL API contract specification
- [x] quickstart.md - Manual test scenarios and validation steps
- [x] CLAUDE.md - Updated agent context with new tech stack info

---
*Based on Constitution v1.0.0 - See `.specify/memory/constitution.md`*
