<!--
Sync Impact Report:
- Version change: [INITIAL] → 1.0.0
- Constitution creation: First version based on project description and existing AGENTS.md
- Modified principles: N/A (initial creation)
- Added sections:
  * Core Principles (5 principles)
  * Development Workflow
  * Security & Data Handling
  * Governance
- Templates requiring updates:
  ✅ constitution.md (this file)
  ⚠ plan-template.md (needs Constitution Check section examples updated)
  ⚠ spec-template.md (alignment verified, no changes needed)
  ⚠ tasks-template.md (alignment verified, no changes needed)
- Follow-up TODOs: None
-->

# Kubernetes Auditing Dashboard Constitution

## Core Principles

### I. Auditing Data Integrity

The system MUST accurately capture and preserve Kubernetes audit events without modification or loss. All webhook data received from the Kubernetes API server must be persisted atomically before acknowledgment. Event timestamps, user identities, and resource changes must remain immutable once recorded. Audit logs are the single source of truth for cluster activity.

**Rationale**: Audit integrity is non-negotiable for debugging, compliance, and security investigations. Any data corruption or loss undermines the entire purpose of the dashboard.

### II. Test-Driven Development (NON-NEGOTIABLE)

All features MUST follow the TDD cycle: write tests first, verify they fail, then implement. Go backend tests use the standard `testing` package and reside alongside implementation files (`*_test.go`). GraphQL resolvers require pagination and query coverage (`gql/*_test.go`). UI changes need visual validation and linting (`npm run lint`). Integration tests validate end-to-end flows.

**Rationale**: The repository already demonstrates TDD discipline with pagination tests and explicit test guidelines. This prevents regressions in audit data processing and ensures query correctness.

### III. Code Generation & Schema-First Design

Entity definitions live in `ent/schema/*.go` as declarative schemas. GraphQL schemas define the contract before resolvers. Run `make generate` to regenerate Ent code, GraphQL types, and typed documents before committing. Never manually edit generated files. Schema changes require corresponding migration and test updates.

**Rationale**: Ent and GraphQL codegen ensure type safety across the stack. Manual edits to generated code create maintenance nightmares and drift from schema definitions.

### IV. Conventional Commits & PR Discipline

Follow Conventional Commits (`feat:`, `fix:`, `chore:`) with imperative subjects under 72 characters. PRs MUST include:

- Description of changes and rationale
- Test results (`go test ./...`, `npm run lint`)
- UI screenshots when frontend changes
- References to related issues

Commits go to feature branches; merge only after CI/local checks pass.

**Rationale**: Consistent history enables automated changelog generation and clear communication. The repository already uses this pattern (`feat: update recent changes page`, `chore: setup AGNETS.md`).

### V. Observability & Debugging

All critical operations (webhook ingestion, database writes, GraphQL queries) MUST emit structured logs. Errors include context (user, resource, timestamp). The UI displays query errors clearly. Local development uses `make dev` for full-stack visibility (Next.js, GraphQL codegen watch, Go backend logs).

**Rationale**: Debugging audit pipelines requires clear visibility into data flow. Silent failures in webhook handlers or GraphQL resolvers are unacceptable.

## Development Workflow

### Project Structure

- **Backend**: `cmd/kubernetes-auditing-dashboard` (entry point), `pkg/services` (shared logic), `ent/schema` (entities), `gql` (GraphQL resolvers/tests)
- **Frontend**: `ui/pages` (Next.js routes), `ui/modules` (shared UI), `ui/styles` (Tailwind), `ui/public` (static assets)
- **Kubernetes Config**: `script/kube-apiserver-config` (Minikube audit webhook setup)

### Build & Test Commands

- `make dev`: Full-stack development loop (Next.js dev server + GraphQL codegen watch + Go backend)
- `make generate`: Regenerate Ent code, run `go generate ./...`, tidy modules, refresh GraphQL documents
- `go test ./...`: Run all Go unit and integration tests
- `cd ui && npm run lint`: Apply Next.js/ESLint checks

### Code Style

- **Go**: Run `gofmt` before committing. Exported symbols PascalCase, internal helpers camelCase. No package-level underscores.
- **TypeScript/React**: Follow bundled ESLint config. Components PascalCase, hooks `useCamelCase`. Use Tailwind utilities from existing styles.
- **Comments**: Never use end-of-line comments. Never use Chinese in code comments.

## Security & Data Handling

### Kubernetes Configuration

`script/kube-apiserver-config` is the single source of truth for Minikube audit setup. Update both the manifest patch and README instructions together. Never commit real auditing data from production clusters. Redact or refresh `data.db` before pushing.

### Secrets Management

Store sensitive connection strings, API tokens, and credentials in environment variables or Kubernetes secrets. Never hardcode them in Go services or Next.js modules. The `.env` pattern is acceptable for local development but must be gitignored.

### Data Retention

Audit events are retained indefinitely by default for compliance and historical analysis. Deletion policies (if implemented) MUST be explicit, logged, and require user confirmation. Bulk deletion without audit trails is prohibited.

## Governance

### Amendment Process

Constitution changes require:

1. Documented rationale for the amendment
2. Review of impact on existing features and workflows
3. Update of dependent templates (plan, spec, tasks)
4. Semantic version bump (MAJOR for breaking changes, MINOR for new principles, PATCH for clarifications)
5. Commit message: `docs: amend constitution to vX.Y.Z (description)`

### Compliance Review

All PRs MUST verify compliance with constitutional principles. Violations require explicit justification in the Complexity Tracking section of the implementation plan. Unjustified violations block merge.

### Versioning Policy

- **MAJOR** (X.0.0): Backward-incompatible governance changes, principle removals, or redefinitions that invalidate existing implementations
- **MINOR** (0.X.0): New principles added, existing principles materially expanded, new mandatory sections
- **PATCH** (0.0.X): Clarifications, wording improvements, typo fixes, non-semantic refinements

### Runtime Guidance

For agent-specific implementation details, testing strategies, and technology-specific best practices, see `AGENTS.md` in the repository root. The constitution defines "what" and "why"; agent guidance files define "how."

**Version**: 1.0.0 | **Ratified**: 2025-10-06 | **Last Amended**: 2025-10-06
