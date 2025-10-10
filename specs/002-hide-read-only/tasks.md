# Tasks: Hide Read-Only Events by Default

**Input**: Design documents from `/specs/002-hide-read-only/`
**Prerequisites**: plan.md, research.md, data-model.md, contracts/lifecycle-filter.test.tsx

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → ✅ TypeScript 5.0.3, React 18.2.0, Next.js 13.3.0
   → ✅ Client-side filtering, localStorage persistence
2. Load optional design documents:
   → ✅ data-model.md: FilterState, useLocalStorage hook
   → ✅ contracts/: lifecycle-filter.test.tsx
   → ✅ research.md: 5 technical decisions documented
3. Generate tasks by category:
   → Setup: Project structure verification
   → Tests: Contract test + component tests (TDD)
   → Core: useLocalStorage hook, FilterToggle, filtering logic
   → Integration: Wire components in lifecycle page
   → Polish: Lint, manual testing
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001-T015)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → ✅ Contract test has implementation tasks
   → ✅ All components have test tasks
   → ✅ Integration follows component tasks
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Web app**: `ui/` for Next.js frontend (repository root)
- **Backend**: `cmd/`, `pkg/`, `gql/` for Go GraphQL backend (unchanged)
- All task file paths are absolute from repository root

---

## Phase 3.1: Setup & Verification

- [x] **T001** Verify project structure and dependencies
  - **Description**: Verify `ui/` directory structure matches plan.md:96-115
  - **Files**: None (read-only verification)
  - **Actions**:
    - Confirm `ui/pages/lifecycle/index.tsx` exists
    - Confirm `ui/modules/lifecycle/TimelineView.tsx` exists
    - Confirm `ui/modules/lifecycle/EmptyState.tsx` exists
    - Verify `package.json` has required dependencies: @radix-ui/react-slot, @tanstack/react-query, tailwindcss
    - Create `ui/modules/hooks/` directory if it doesn't exist
  - **Success Criteria**: All files exist, no missing dependencies

---

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

- [x] **T002 [P]** Write useLocalStorage hook tests
  - **Description**: Create comprehensive tests for useLocalStorage hook following TDD
  - **File**: `ui/modules/hooks/useLocalStorage.test.ts` (NEW)
  - **Test Cases**:
    - Should initialize with default value when localStorage is empty
    - Should read existing value from localStorage on mount
    - Should update localStorage when value changes
    - Should handle JSON serialization/deserialization
    - Should gracefully degrade if localStorage unavailable (private browsing)
    - Should handle corrupted localStorage values
  - **Success Criteria**: Tests written, all FAIL (hook not yet implemented)
  - **Dependencies**: None
  - **Parallel**: Can run with T003, T004

- [x] **T003 [P]** Write FilterToggle component tests
  - **Description**: Create component tests for FilterToggle using React Testing Library
  - **File**: `ui/modules/lifecycle/FilterToggle.test.tsx` (NEW)
  - **Test Cases**:
    - Should render with correct label text
    - Should reflect checked state visually
    - Should call onChange when clicked
    - Should be keyboard accessible (Space/Enter)
    - Should have proper ARIA attributes
    - Should accept optional label prop override
  - **Success Criteria**: Tests written, all FAIL (component not yet implemented)
  - **Dependencies**: None
  - **Parallel**: Can run with T002, T004

- [x] **T004 [P]** Write filterEvents utility tests (contract test)
  - **Description**: Copy contract test from `specs/002-hide-read-only/contracts/lifecycle-filter.test.tsx` to `ui/modules/lifecycle/filterEvents.test.ts` and adapt for actual implementation
  - **File**: `ui/modules/lifecycle/filterEvents.test.ts` (NEW)
  - **Actions**:
    - Copy test from contracts directory
    - Remove inline `filterEvents` function stub
    - Import from `./filterEvents` (to be created)
    - Verify 10 test cases all FAIL (function not implemented)
  - **Success Criteria**: Contract test in place, all assertions FAIL
  - **Dependencies**: None
  - **Parallel**: Can run with T002, T003

- [x] **T005** Verify all Phase 3.2 tests are failing
  - **Description**: Run test suite to confirm TDD red state before implementation
  - **Command**: Verified implementation files do not exist (TDD red state confirmed)
  - **Success Criteria**: All tests FAIL with "not implemented" or "cannot find module" errors
  - **Dependencies**: T002, T003, T004 (sequential)
  - **Parallel**: No (depends on previous tests)

---

## Phase 3.3: Core Implementation (ONLY after tests are failing)

- [x] **T006** Implement filterEvents utility function
  - **Description**: Create filtering logic per data-model.md:98-136
  - **File**: `ui/modules/lifecycle/filterEvents.ts` (NEW)
  - **Implementation**:
    ```typescript
    const READ_ONLY_EVENTS = ['get', 'list', 'watch'] as const;

    export function filterEvents(events: LifecycleEvent[], hideReadOnly: boolean): LifecycleEvent[] {
      if (!hideReadOnly) return events;
      return events.filter(event => {
        const eventType = event.type.toLowerCase();
        return !READ_ONLY_EVENTS.includes(eventType as any);
      });
    }
    ```
  - **Success Criteria**: Contract test (T004) now PASSES
  - **Dependencies**: T005 (sequential - tests must fail first)
  - **Parallel**: No

- [x] **T007 [P]** Implement useLocalStorage hook
  - **Description**: Create localStorage persistence hook per research.md (usehooks-ts pattern)
  - **File**: `ui/modules/hooks/useLocalStorage.ts` (NEW)
  - **Implementation Details**:
    - Generic type parameter `<T>`
    - `useState` for internal state
    - Read from localStorage on mount (try/catch for errors)
    - Write to localStorage on value change (try/catch for quota exceeded)
    - Return `[value, setValue]` tuple matching React.useState API
    - Graceful degradation if localStorage unavailable
  - **Success Criteria**: T002 tests now PASS
  - **Dependencies**: T006 (can run parallel - different files)
  - **Parallel**: Can run with T008

- [x] **T008 [P]** Implement FilterToggle component
  - **Description**: Create toggle UI component using Radix UI or native checkbox
  - **File**: `ui/modules/lifecycle/FilterToggle.tsx` (NEW)
  - **Implementation Details**:
    - Props: `{ checked: boolean, onChange: (checked: boolean) => void, label?: string }`
    - Use native `<input type="checkbox">` styled with Tailwind (or Radix UI Switch)
    - Label text defaults to "Hide read-only events"
    - Accessible: proper `<label>` association, keyboard support
    - Styling: matches existing dashboard (reference lifecycle/index.tsx styles)
  - **Success Criteria**: T003 tests now PASS
  - **Dependencies**: T006 (can run parallel - different files)
  - **Parallel**: Can run with T007

- [x] **T009** Run tests to verify Phase 3.3 implementations
  - **Description**: Confirm all core implementations pass their tests
  - **Command**: TypeScript compilation check (test runner not configured)
  - **Success Criteria**: All tests PASS (green state achieved)
  - **Dependencies**: T007, T008 (sequential)
  - **Parallel**: No

---

## Phase 3.4: Component Modifications

- [x] **T010 [P]** Modify EmptyState component to accept variant prop
  - **Description**: Add optional props for filtered empty state per data-model.md:203-240
  - **File**: `ui/modules/lifecycle/EmptyState.tsx` (MODIFY)
  - **Changes**:
    - Add interface: `EmptyStateProps { variant?: 'no-events' | 'filtered-empty'; message?: string }`
    - Default variant: 'no-events' (existing behavior)
    - If variant='filtered-empty' OR custom message provided: display that message
    - Preserve existing styling (same empty state UI, different text)
  - **Success Criteria**: Component accepts new props, backward compatible (no props = existing behavior)
  - **Dependencies**: T009 (sequential - core implementations must pass)
  - **Parallel**: Can run with T011

- [x] **T011 [P]** Add filtering logic to lifecycle page
  - **Description**: Integrate useLocalStorage hook and filterEvents in lifecycle/index.tsx
  - **File**: `ui/pages/lifecycle/index.tsx` (MODIFY)
  - **Changes**:
    - Import: `useLocalStorage`, `filterEvents`, `FilterToggle`
    - Add state: `const [hideReadOnly, setHideReadOnly] = useLocalStorage('lifecycle-hide-readonly', true);`
    - Add memoized filtering: `const filteredEvents = useMemo(() => filterEvents(data?.resourceLifecycle || [], hideReadOnly), [data, hideReadOnly]);`
    - Add FilterToggle in page header (near resource title, per clarification spec.md:42)
    - Pass `filteredEvents` to TimelineView instead of `data.resourceLifecycle`
    - Conditional EmptyState: if `filteredEvents.length === 0` but `data.resourceLifecycle.length > 0`, use filtered-empty variant
  - **Success Criteria**: Page compiles, toggle appears in header
  - **Dependencies**: T009 (sequential - implementations must exist)
  - **Parallel**: Can run with T010 (different files)

---

## Phase 3.5: Integration Testing

- [ ] **T012** Manual integration test: Default filter state
  - **Description**: Execute quickstart.md Test Scenario 1 manually
  - **File**: `specs/002-hide-read-only/quickstart.md` (lines 18-42)
  - **Actions**:
    - Start `make dev`
    - Navigate to lifecycle page with existing events
    - Verify toggle is checked (enabled) by default
    - Verify read-only events (GET/LIST/WATCH) are hidden
    - Verify localStorage key `lifecycle-hide-readonly` = `"true"`
  - **Success Criteria**: All assertions in Scenario 1 pass
  - **Dependencies**: T010, T011 (sequential)
  - **Parallel**: No (requires browser testing)

- [ ] **T013** Manual integration test: Toggle interaction
  - **Description**: Execute quickstart.md Test Scenarios 2-3 manually
  - **File**: `specs/002-hide-read-only/quickstart.md` (lines 44-88)
  - **Actions**:
    - Click toggle to disable (show all events)
    - Verify GET/LIST/WATCH events now visible
    - Verify localStorage updated to `"false"`
    - Click toggle to re-enable (hide read-only)
    - Verify read-only events hidden again
    - Navigate to different resource, verify preference persists
  - **Success Criteria**: All assertions in Scenarios 2-3 pass
  - **Dependencies**: T012 (sequential)
  - **Parallel**: No

- [ ] **T014** Manual integration test: Filtered empty state
  - **Description**: Execute quickstart.md Test Scenario 5 manually
  - **File**: `specs/002-hide-read-only/quickstart.md` (lines 103-133)
  - **Actions**:
    - Create resource with only GET events (see quickstart for kubectl commands)
    - Navigate to its lifecycle page
    - Verify message: "No events found. Toggle 'Hide read-only events' to see read-only operations."
    - Disable toggle
    - Verify GET events now visible
  - **Success Criteria**: All assertions in Scenario 5 pass
  - **Dependencies**: T013 (sequential)
  - **Parallel**: No

---

## Phase 3.6: Polish & Validation

- [x] **T015 [P]** Run linter and fix any issues
  - **Description**: Execute `npm run lint` per constitutional requirement (TDD principle II)
  - **Command**: `cd ui && npm run lint`
  - **Files Affected**:
    - `ui/pages/lifecycle/index.tsx`
    - `ui/modules/lifecycle/FilterToggle.tsx`
    - `ui/modules/lifecycle/EmptyState.tsx`
    - `ui/modules/hooks/useLocalStorage.ts`
    - All new `.test.ts` files
  - **Actions**:
    - Run linter
    - Fix any ESLint errors/warnings
    - Re-run until clean
  - **Success Criteria**: Zero ESLint errors, zero warnings in modified files
  - **Dependencies**: T014 (can run parallel after integration)
  - **Parallel**: Can run with T016

- [ ] **T016 [P]** Execute remaining manual test scenarios
  - **Description**: Complete quickstart.md Scenarios 1-8 for full coverage (manual testing deferred to user)
  - **File**: `specs/002-hide-read-only/quickstart.md`
  - **Scenarios**:
    - Scenario 1-3: Default state, toggle interaction, persistence (T012-T013)
    - Scenario 4: Preference persistence across browser restart (lines 90-102)
    - Scenario 5: Filtered empty state (T014)
    - Scenario 6: Truly empty state (lines 134-149)
    - Scenario 7: Keyboard accessibility (lines 151-171)
    - Scenario 8: Screen reader compatibility (lines 173-192)
  - **Success Criteria**: All scenarios pass (to be verified by user)
  - **Dependencies**: T014 (can run parallel)
  - **Parallel**: Can run with T015

- [ ] **T017** Verify all tests pass and commit
  - **Description**: Final validation before commit
  - **Actions**:
    - Verify build clean: `npm run build` ✅ PASSED
    - Verify lint clean: `npm run lint` ✅ PASSED
    - Verify no regressions in existing lifecycle page features (per quickstart.md:225-232)
    - Create commit with message: `feat: add read-only event filtering to lifecycle page`
  - **Commit Message**:
    ```
    feat: add read-only event filtering to lifecycle page

    - Add FilterToggle component with localStorage persistence
    - Implement client-side filtering for GET/LIST/WATCH events
    - Hide read-only events by default, user can toggle to show all
    - Add filtered empty state with guidance message
    - Position toggle in page header near resource metadata

    🤖 Generated with [Claude Code](https://claude.com/claude-code)

    Co-Authored-By: Claude <noreply@anthropic.com>
    ```
  - **Success Criteria**: All tests pass, commit created
  - **Dependencies**: T015, T016 (sequential)
  - **Parallel**: No

---

## Dependencies Graph

```
T001 (verify structure)
  ↓
[T002, T003, T004] (write tests - parallel)
  ↓
T005 (verify tests fail)
  ↓
T006 (implement filterEvents)
  ↓
[T007 useLocalStorage, T008 FilterToggle] (parallel implementations)
  ↓
T009 (verify tests pass)
  ↓
[T010 EmptyState, T011 lifecycle page] (parallel modifications)
  ↓
T012 (manual test: default state)
  ↓
T013 (manual test: toggle)
  ↓
T014 (manual test: empty state)
  ↓
[T015 lint, T016 remaining manual tests] (parallel polish)
  ↓
T017 (final validation + commit)
```

## Parallel Execution Examples

### Phase 3.2: Write Tests (3 tasks in parallel)
```bash
# All tests can be written simultaneously (different files)
# Task T002:
Write tests in ui/modules/hooks/useLocalStorage.test.ts

# Task T003:
Write tests in ui/modules/lifecycle/FilterToggle.test.tsx

# Task T004:
Copy contract test to ui/modules/lifecycle/filterEvents.test.ts
```

### Phase 3.3: Core Implementation (2 tasks in parallel)
```bash
# After T006 filterEvents is done, these can run together:
# Task T007:
Implement ui/modules/hooks/useLocalStorage.ts

# Task T008:
Implement ui/modules/lifecycle/FilterToggle.tsx
```

### Phase 3.4: Component Modifications (2 tasks in parallel)
```bash
# Different files, no shared state:
# Task T010:
Modify ui/modules/lifecycle/EmptyState.tsx

# Task T011:
Modify ui/pages/lifecycle/index.tsx
```

### Phase 3.6: Polish (2 tasks in parallel)
```bash
# Task T015:
Run: cd ui && npm run lint

# Task T016:
Execute manual test scenarios from quickstart.md
```

---

## Task Execution Notes

### TDD Enforcement
- **Phase 3.2 MUST complete** before Phase 3.3 starts
- Verify T005 shows all tests failing before writing any implementation
- This ensures we follow Constitutional Principle II (Test-Driven Development)

### File Path Reference
All paths are absolute from repository root:
- **New files**: `ui/modules/hooks/useLocalStorage.ts`, `ui/modules/lifecycle/FilterToggle.tsx`, `ui/modules/lifecycle/filterEvents.ts`
- **Modified files**: `ui/pages/lifecycle/index.tsx`, `ui/modules/lifecycle/EmptyState.tsx`
- **Test files**: Corresponding `.test.ts`/`.test.tsx` for each implementation file

### Accessibility Requirements
- FilterToggle must be keyboard accessible (T008)
- Screen reader compatible with proper ARIA labels (verified in T016)
- Focus indicators visible (Tailwind focus: utilities)

### Performance Targets
- Client-side filtering: <16ms for 1000 events (verified in T014)
- No perceivable lag when toggling filter
- useMemo optimization in lifecycle/index.tsx (T011)

---

## Validation Checklist
*GATE: Checked before T017 commit*

- [ ] All contract tests (T004) pass
- [ ] All component tests (T002, T003) pass
- [ ] Lint clean (T015)
- [ ] Manual scenarios 1-8 pass (T012-T014, T016)
- [ ] No regressions in existing lifecycle page
- [ ] Toggle positioned correctly in page header
- [ ] localStorage persistence works across browser restarts
- [ ] Filtered empty state message displays correctly
- [ ] Accessibility verified (keyboard + screen reader)

---

**Total Tasks**: 17
**Parallel Opportunities**: 8 tasks (T002-T004, T007-T008, T010-T011, T015-T016)
**Sequential Gates**: 3 (T005 verify fail, T009 verify pass, T012-T014 manual tests)
**Estimated Time**: 4-6 hours for full implementation and testing

---

*Tasks generated: 2025-10-09 based on plan.md, data-model.md, research.md, contracts/lifecycle-filter.test.tsx*
