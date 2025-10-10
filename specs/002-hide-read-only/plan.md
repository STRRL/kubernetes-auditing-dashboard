# Implementation Plan: Hide Read-Only Events by Default

**Branch**: `002-hide-read-only` | **Date**: 2025-10-09 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-hide-read-only/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → ✅ Loaded successfully
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → ✅ All clarifications resolved in spec
3. Fill the Constitution Check section
   → ✅ Section filled below
4. Evaluate Constitution Check section
   → ✅ No violations detected
5. Execute Phase 0 → research.md
   → ✅ Complete
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, AGENTS.md
   → ✅ Complete
7. Re-evaluate Constitution Check section
   → ✅ PASS - No violations after design
8. Plan Phase 2 → Describe task generation approach
   → ✅ Described below
9. STOP - Ready for /tasks command
   → ✅ Ready
```

**IMPORTANT**: The /plan command STOPS at step 9. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary

This feature adds client-side filtering of read-only Kubernetes audit events (get, list, watch) on the lifecycle page. By default, only mutating events are displayed to reduce visual clutter. A toggle control in the page header allows users to show/hide read-only events, with the preference persisted in localStorage across browser sessions.

**Technical Approach**: Client-side filtering using React state and localStorage, with no backend changes required. The existing GraphQL `resourceLifecycle` query returns all events; filtering happens in the UI component before rendering.

## Technical Context

**Language/Version**: TypeScript 5.0.3, React 18.2.0, Next.js 13.3.0
**Primary Dependencies**: @tanstack/react-query 4.29.1, graphql-request 6.0.0, Tailwind CSS 3.3.1
**Storage**: Browser localStorage for user preference persistence
**Testing**: npm run lint (ESLint via next lint)
**Target Platform**: Modern browsers (ES2020+) supporting localStorage
**Project Type**: web (Next.js frontend + Go GraphQL backend)
**Performance Goals**: Instantaneous client-side filtering (<16ms for 60fps), no perceivable lag for typical event lists (100-1000 events)
**Constraints**: No backend changes, maintain existing GraphQL contract, preserve chronological ordering
**Scale/Scope**: Single lifecycle page modification, 3 new UI components/modules, localStorage key management

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**I. Auditing Data Integrity**
- [x] Webhook data persistence is atomic (acknowledge only after successful DB write) - N/A: No backend changes
- [x] Event timestamps, user identities, and resource changes are immutable - ✅ UI filtering does not modify audit data
- [x] No modifications to audit event data after ingestion - ✅ Client-side filtering preserves all data

**II. Test-Driven Development**
- [x] Tests written before implementation (TDD cycle enforced) - ✅ Will write TimelineView.test.tsx before modifying TimelineView component
- [x] Go tests use standard `testing` package in `*_test.go` files - N/A: No Go changes
- [x] GraphQL resolvers have pagination/query coverage in `gql/*_test.go` - N/A: No resolver changes
- [x] UI changes include lint checks (`npm run lint`) - ✅ Included in testing approach

**III. Code Generation & Schema-First Design**
- [x] Entity changes start with `ent/schema/*.go` modifications - N/A: No entity changes
- [x] GraphQL schemas defined before resolvers - N/A: No schema changes
- [x] `make generate` run before committing schema changes - N/A: No codegen changes
- [x] No manual edits to generated files - ✅ No generated file modifications

**IV. Conventional Commits & PR Discipline**
- [x] Commit messages follow Conventional Commits format - ✅ Will use `feat: add read-only event filtering to lifecycle page`
- [x] PRs include test results and rationale - ✅ Will include lint results and screenshots
- [x] UI changes include screenshots - ✅ Before/after screenshots of toggle and filtering
- [x] Feature branch workflow followed - ✅ Already on `002-hide-read-only` branch

**V. Observability & Debugging**
- [x] Critical operations emit structured logs - N/A: Client-side filtering, browser DevTools sufficient
- [x] Errors include context (user, resource, timestamp) - N/A: No error-prone operations introduced
- [x] UI displays query errors clearly - ✅ Existing error handling unchanged
- [x] Local dev uses `make dev` for full-stack visibility - ✅ Development workflow unchanged

## Project Structure

### Documentation (this feature)
```
specs/002-hide-read-only/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
│   └── lifecycle-filter.test.tsx  # Contract test for filtering logic
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
ui/
├── pages/
│   └── lifecycle/
│       └── index.tsx                    # MODIFY: Add toggle and filtering state
├── modules/
│   ├── lifecycle/
│   │   ├── TimelineView.tsx             # MODIFY: Accept filtered events prop
│   │   ├── TimelineView.test.tsx        # NEW: Test filtering logic
│   │   ├── FilterToggle.tsx             # NEW: Toggle component
│   │   ├── FilterToggle.test.tsx        # NEW: Toggle interaction tests
│   │   ├── EmptyState.tsx               # MODIFY: Add filtered empty state variant
│   │   └── DiffViewer.tsx               # UNCHANGED
│   └── hooks/
│       ├── useLocalStorage.ts           # NEW: localStorage persistence hook
│       └── useLocalStorage.test.ts      # NEW: Hook tests
└── __tests__/
    └── lifecycle-filter.test.tsx        # NEW: Integration test for filter workflow
```

**Structure Decision**: Web application structure with Next.js frontend (`ui/`) and Go GraphQL backend (unchanged). All modifications are client-side React components in the existing `ui/modules/lifecycle/` directory. New custom hook in `ui/modules/hooks/` for localStorage abstraction.

## Phase 0: Outline & Research

### Research Findings

**Decision 1: Client-Side vs Server-Side Filtering**
- **Decision**: Client-side filtering in React component
- **Rationale**:
  - Feature spec clarifies "Events should be filtered client-side" (spec.md:49)
  - No backend changes required, maintaining audit data integrity
  - Typical event lists (100-1000 items) filter instantly in modern browsers
  - GraphQL contract unchanged, reducing deployment complexity
- **Alternatives considered**:
  - Server-side filtering: Rejected due to constitutional preference for minimal backend changes when client-side suffices
  - GraphQL query parameter: Rejected to avoid schema modification for UI preference

**Decision 2: State Persistence Strategy**
- **Decision**: Browser localStorage with custom React hook
- **Rationale**:
  - Clarification spec.md:40 requires "persist indefinitely across browser sessions using localStorage"
  - Avoids backend user preference storage complexity
  - Standard React pattern with useLocalStorage hook
  - Graceful degradation if localStorage unavailable (defaults to hide read-only)
- **Alternatives considered**:
  - URL query parameter: Rejected, doesn't persist across navigation
  - Session storage: Rejected, doesn't persist across browser restarts
  - Cookie: Over-engineered for client-only preference

**Decision 3: Event Type Classification**
- **Decision**: Hardcoded blacklist array `['get', 'list', 'watch']`
- **Rationale**:
  - Clarification spec.md:39 specifies "Blacklist read-only; show rest"
  - Event types are stable Kubernetes audit verbs, unlikely to change
  - Simple array check faster than regex or complex logic
  - Extensible: future features can move to configuration if needed
- **Alternatives considered**:
  - Configuration file: Over-engineered for 3 static values
  - GraphQL enum metadata: Requires schema changes

**Decision 4: Toggle Component Library**
- **Decision**: Native HTML checkbox styled with Tailwind, or Radix UI Switch if interactive polish needed
- **Rationale**:
  - Project already uses Radix UI (@radix-ui/react-dialog in package.json)
  - Accessible by default (keyboard nav, screen readers)
  - Matches existing UI patterns in dashboard
- **Alternatives considered**:
  - Custom toggle: Reinventing accessibility is anti-constitutional (simplicity principle)
  - DaisyUI: Project migrated away from DaisyUI (commit ecef4b0), now uses shadcn-style Radix + Tailwind

**Decision 5: Empty State Messaging**
- **Decision**: Conditional EmptyState variant with guidance text
- **Rationale**:
  - Clarification spec.md:43 specifies exact message: "No events found. Toggle 'Hide read-only events' to see read-only operations."
  - Existing EmptyState.tsx component can accept custom message prop
  - Maintains consistent empty state styling
- **Alternatives considered**: Modal/toast notification: Too intrusive for informational message

## Phase 1: Design & Contracts

### Data Model

**No new entities or backend data model changes.** This feature is purely client-side state management.

**Client-Side State Model**:

```typescript
// Filter state (component-local)
interface FilterState {
  hideReadOnly: boolean;  // Default: true
}

// localStorage schema
{
  "lifecycle-hide-readonly": "true" | "false"  // String representation of boolean
}

// Event type (already defined in GraphQL codegen output)
interface LifecycleEvent {
  id: string;
  type: EventType;  // Enum: CREATE | UPDATE | DELETE | GET
  timestamp: string;
  user: string;
  resourceState: any;
  diff?: ResourceDiff;
}
```

**Filtering Logic**:
```typescript
const READ_ONLY_EVENTS = ['get', 'list', 'watch'] as const;

function filterEvents(events: LifecycleEvent[], hideReadOnly: boolean): LifecycleEvent[] {
  if (!hideReadOnly) return events;
  return events.filter(event =>
    !READ_ONLY_EVENTS.includes(event.type.toLowerCase())
  );
}
```

### API Contracts

**No GraphQL schema changes.** The existing `resourceLifecycle` query (gql/lifecycle.graphql:13-28) returns all events; filtering happens client-side.

**Component Contracts** (TypeScript interfaces):

```typescript
// ui/modules/lifecycle/FilterToggle.tsx
interface FilterToggleProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  label?: string;  // Default: "Hide read-only events"
}

// ui/modules/lifecycle/TimelineView.tsx (modified)
interface TimelineViewProps {
  events: LifecycleEvent[];  // Now expects pre-filtered events
}

// ui/modules/lifecycle/EmptyState.tsx (modified)
interface EmptyStateProps {
  variant?: 'no-events' | 'filtered-empty';  // NEW
  message?: string;  // Custom message override
}

// ui/modules/hooks/useLocalStorage.ts
function useLocalStorage<T>(
  key: string,
  initialValue: T
): [T, (value: T | ((prev: T) => T)) => void];
```

### Contract Tests

Contract test at `specs/002-hide-read-only/contracts/lifecycle-filter.test.tsx`:

```typescript
import { describe, it, expect } from '@jest/globals';

describe('Lifecycle Event Filtering Contract', () => {
  const READ_ONLY_EVENTS = ['get', 'list', 'watch'];

  it('should define read-only event types as get, list, watch', () => {
    expect(READ_ONLY_EVENTS).toEqual(['get', 'list', 'watch']);
  });

  it('should filter out read-only events when hideReadOnly is true', () => {
    const events = [
      { id: '1', type: 'CREATE', timestamp: '2023-01-01' },
      { id: '2', type: 'GET', timestamp: '2023-01-02' },
      { id: '3', type: 'UPDATE', timestamp: '2023-01-03' },
      { id: '4', type: 'LIST', timestamp: '2023-01-04' },
    ];

    const filtered = filterEvents(events, true);
    expect(filtered).toHaveLength(2);
    expect(filtered.map(e => e.type)).toEqual(['CREATE', 'UPDATE']);
  });

  it('should preserve chronological order after filtering', () => {
    const events = [
      { id: '1', type: 'UPDATE', timestamp: '2023-01-03' },
      { id: '2', type: 'GET', timestamp: '2023-01-02' },
      { id: '3', type: 'CREATE', timestamp: '2023-01-01' },
    ];

    const filtered = filterEvents(events, true);
    expect(filtered.map(e => e.id)).toEqual(['1', '3']);
  });

  it('should return all events when hideReadOnly is false', () => {
    const events = [
      { id: '1', type: 'CREATE', timestamp: '2023-01-01' },
      { id: '2', type: 'GET', timestamp: '2023-01-02' },
    ];

    const filtered = filterEvents(events, false);
    expect(filtered).toHaveLength(2);
  });
});
```

### User Story Test Scenarios

From spec.md:52-56, extracted integration test scenarios for `ui/__tests__/lifecycle-filter.test.tsx`:

**Scenario 1: Default filter state**
```typescript
it('should hide read-only events by default on page load', async () => {
  render(<LifecyclePage />);

  await waitFor(() => {
    expect(screen.queryByText(/GET/)).not.toBeInTheDocument();
    expect(screen.getByText(/CREATE|UPDATE|DELETE/)).toBeInTheDocument();
  });

  const toggle = screen.getByLabelText('Hide read-only events');
  expect(toggle).toBeChecked();
});
```

**Scenario 2: Toggle shows read-only events**
```typescript
it('should show read-only events when toggle is disabled', async () => {
  render(<LifecyclePage />);

  const toggle = screen.getByLabelText('Hide read-only events');
  fireEvent.click(toggle);

  await waitFor(() => {
    expect(screen.getByText(/GET|LIST|WATCH/)).toBeInTheDocument();
  });

  expect(toggle).not.toBeChecked();
});
```

**Scenario 3: Preference persistence**
```typescript
it('should persist toggle preference across page navigation', async () => {
  const { rerender } = render(<LifecyclePage />);

  const toggle = screen.getByLabelText('Hide read-only events');
  fireEvent.click(toggle);  // Disable filter

  // Simulate navigation by unmounting and remounting
  rerender(<div />);
  rerender(<LifecyclePage />);

  const newToggle = screen.getByLabelText('Hide read-only events');
  expect(newToggle).not.toBeChecked();  // Preference retained
});
```

**Scenario 4: Filtered empty state**
```typescript
it('should show guidance message when only read-only events exist', async () => {
  // Mock GraphQL to return only GET events
  mockGraphQL({ resourceLifecycle: [{ type: 'GET', ... }] });

  render(<LifecyclePage />);

  await waitFor(() => {
    expect(screen.getByText(
      /No events found. Toggle 'Hide read-only events' to see read-only operations./
    )).toBeInTheDocument();
  });
});
```

### Quickstart Test Plan

See `specs/002-hide-read-only/quickstart.md` for manual validation steps.

### Agent Context Update

Running incremental AGENTS.md update to preserve manual additions:

```bash
.specify/scripts/bash/update-agent-context.sh claude
```

This will add:
- Recent change: "002-hide-read-only: Client-side lifecycle event filtering with localStorage persistence"
- Technology: TypeScript React hooks (useLocalStorage pattern)
- UI testing reminder: Include lint checks and component interaction tests

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
1. Load `.specify/templates/tasks-template.md` as base structure
2. Generate tasks from Phase 1 deliverables:
   - Contract test implementation (1 task)
   - Custom hook implementation + tests (2 tasks, parallel)
   - Component modifications (3 tasks: FilterToggle, TimelineView update, EmptyState variant)
   - Integration (1 task: wire components in lifecycle/index.tsx)
   - Validation (2 tasks: lint, manual quickstart)

**Task Dependencies**:
- Contract tests can run immediately (no dependencies)
- useLocalStorage hook must complete before FilterToggle (hook dependency)
- FilterToggle and TimelineView modifications are parallel [P]
- EmptyState modification parallel with FilterToggle [P]
- Integration task depends on all component tasks
- Validation depends on integration

**Ordering Strategy**:
1. Write contract test (TDD: fails initially)
2. Implement useLocalStorage hook + tests [P] with FilterToggle component
3. Modify TimelineView to accept filtered events [P]
4. Modify EmptyState for filtered variant [P]
5. Integrate in lifecycle/index.tsx (wire toggle state + filtering logic)
6. Run npm run lint
7. Execute quickstart.md manual tests
8. Verify contract test now passes

**Estimated Output**: 12-15 numbered tasks in tasks.md

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)
**Phase 4**: Implementation (execute tasks.md following TDD workflow)
**Phase 5**: Validation (contract tests pass, lint clean, quickstart scenarios verified)

## Complexity Tracking
*No constitutional violations detected. This section is empty.*

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
- [x] Complexity deviations documented (none)

---
*Based on Constitution v1.0.0 - See `.specify/memory/constitution.md`*
