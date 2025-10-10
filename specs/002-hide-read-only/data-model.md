# Data Model: Hide Read-Only Events by Default

## Overview

This feature requires **no backend data model changes**. All state management is client-side using React component state and browser localStorage.

## Client-Side State Model

### Filter State (Component State)

```typescript
interface FilterState {
  /**
   * Whether to hide read-only events (get, list, watch) from the timeline.
   * Default: true (read-only events hidden by default)
   */
  hideReadOnly: boolean;
}
```

**State Location**: React component state in `ui/pages/lifecycle/index.tsx`

**State Management**: useState hook with useLocalStorage hook for persistence

**Lifecycle**:
1. Component mounts → Read from localStorage (key: `lifecycle-hide-readonly`)
2. User toggles → Update component state + write to localStorage
3. Component unmounts → State persisted in localStorage
4. Component remounts (navigation/refresh) → Read persisted value

---

### localStorage Schema

**Key**: `lifecycle-hide-readonly`

**Value Format**: JSON-serialized boolean

```typescript
// Example stored values
localStorage.getItem('lifecycle-hide-readonly')
// Returns: "true"  (hide read-only events - default)
// Returns: "false" (show all events)
```

**Schema**:
```typescript
type StoredFilterPreference = 'true' | 'false';  // String representation of boolean
```

**Serialization**: `JSON.stringify(boolean)` → `JSON.parse(string)`

**Error Handling**:
- If localStorage unavailable (private browsing, quota exceeded): Default to `hideReadOnly = true`
- If corrupted value (non-boolean): Reset to `true` and log warning

---

## Event Data Model (Existing, Unchanged)

### LifecycleEvent (GraphQL Type)

**Source**: Generated from `gql/lifecycle.graphql:34-52`

```typescript
interface LifecycleEvent {
  /** Unique identifier for the audit event */
  id: string;

  /** Type of lifecycle event (CREATE, UPDATE, DELETE, GET) */
  type: EventType;

  /** ISO 8601 timestamp when the event occurred */
  timestamp: string;

  /** User or service account that triggered the event */
  user: string;

  /** Complete resource state at the time of this event (YAML as JSON) */
  resourceState: any;  // JSON scalar

  /** Diff showing changes from previous version (null for CREATE and DELETE events) */
  diff?: ResourceDiff;
}

enum EventType {
  CREATE = 'CREATE',
  UPDATE = 'UPDATE',
  DELETE = 'DELETE',
  GET = 'GET',
  // Note: 'LIST' and 'WATCH' may be added by backend in future
}
```

**No modifications to this structure.** Filtering operates on the `type` field only.

---

## Filtering Logic

### Event Classification

**Read-Only Events** (blacklist):
```typescript
const READ_ONLY_EVENTS = ['get', 'list', 'watch'] as const;
```

**Rationale**: Clarification spec.md:39 specifies blacklist approach—only these verbs are hidden by default.

**Case Handling**: Comparison is case-insensitive (`event.type.toLowerCase()`) to handle potential backend casing variations.

---

### Filter Function

```typescript
/**
 * Filters lifecycle events based on read-only status.
 *
 * @param events - Full list of lifecycle events from GraphQL query
 * @param hideReadOnly - Whether to hide read-only events (get, list, watch)
 * @returns Filtered event list (or original if hideReadOnly=false)
 */
function filterEvents(
  events: LifecycleEvent[],
  hideReadOnly: boolean
): LifecycleEvent[] {
  if (!hideReadOnly) {
    return events;  // Show all events
  }

  return events.filter(event => {
    const eventType = event.type.toLowerCase();
    return !READ_ONLY_EVENTS.includes(eventType);
  });
}
```

**Performance**: O(n) complexity for n events. Array.filter with includes check is optimized by V8/SpiderMonkey engines.

**Invariants**:
- Chronological order preserved (filter does not reorder)
- Original array unchanged (filter creates new array)
- No mutation of event objects

---

## Component Props Interfaces

### FilterToggle Component

```typescript
interface FilterToggleProps {
  /**
   * Current toggle state.
   * true = Hide read-only events (toggle enabled/checked)
   * false = Show all events (toggle disabled/unchecked)
   */
  checked: boolean;

  /**
   * Callback when toggle state changes.
   * @param checked - New toggle state
   */
  onChange: (checked: boolean) => void;

  /**
   * Optional label text override.
   * Default: "Hide read-only events"
   */
  label?: string;
}
```

**Usage Example**:
```tsx
<FilterToggle
  checked={hideReadOnly}
  onChange={setHideReadOnly}
  label="Hide read-only events"
/>
```

---

### TimelineView Component (Modified)

**Current Props** (from `ui/modules/lifecycle/TimelineView.tsx`):
```typescript
interface TimelineViewProps {
  events: LifecycleEvent[];
}
```

**No prop interface changes required.** TimelineView already accepts `events` array; it will now receive pre-filtered events instead of raw GraphQL result.

**Behavioral Change**:
- **Before**: Receives all events from `data.resourceLifecycle`
- **After**: Receives filtered events from `filterEvents(data.resourceLifecycle, hideReadOnly)`

---

### EmptyState Component (Modified)

**Current Props** (from `ui/modules/lifecycle/EmptyState.tsx`):
```typescript
interface EmptyStateProps {
  // Currently no props (shows hardcoded message)
}
```

**New Props**:
```typescript
interface EmptyStateProps {
  /**
   * Variant of empty state to display.
   * - 'no-events': No events exist at all (default message)
   * - 'filtered-empty': Events exist but all filtered out (show guidance)
   */
  variant?: 'no-events' | 'filtered-empty';

  /**
   * Custom message override (takes precedence over variant).
   * If provided, displays this message instead of default.
   */
  message?: string;
}
```

**Usage Examples**:
```tsx
// All events filtered out
<EmptyState
  variant="filtered-empty"
  message="No events found. Toggle 'Hide read-only events' to see read-only operations."
/>

// Truly no events (existing behavior)
<EmptyState variant="no-events" />
```

---

## State Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│ Component Mount (ui/pages/lifecycle/index.tsx)             │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ useLocalStorage('lifecycle-hide-readonly', true)            │
│ → Reads from localStorage                                   │
│ → Returns: [hideReadOnly, setHideReadOnly]                  │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ GraphQL Query: resourceLifecycle                            │
│ → Returns: LifecycleEvent[]                                 │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ useMemo(() => filterEvents(events, hideReadOnly))           │
│ → Computes filtered events (memoized)                       │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ Render Logic:                                               │
│ - FilterToggle (in page header)                             │
│ - TimelineView (receives filteredEvents)                    │
│ - EmptyState (if filteredEvents.length === 0)               │
└─────────────────────────────────────────────────────────────┘

User clicks toggle
         │
         ▼
┌─────────────────────────────────────────────────────────────┐
│ setHideReadOnly(!hideReadOnly)                              │
│ → Updates component state                                   │
│ → Writes to localStorage                                    │
│ → Triggers re-render (useMemo recalculates filteredEvents)  │
└─────────────────────────────────────────────────────────────┘
```

---

## Data Integrity Guarantees

**Constitutional Principle I: Auditing Data Integrity**

✅ **Webhook data persistence**: Not affected (no backend changes)

✅ **Event immutability**: Filtering creates a new array; original `data.resourceLifecycle` unchanged

✅ **No modifications after ingestion**: Client-side filter only affects presentation; GraphQL response unchanged

**Verification**:
```typescript
// Before filtering
console.log(data.resourceLifecycle.length);  // e.g., 100

// After filtering
const filtered = filterEvents(data.resourceLifecycle, true);
console.log(data.resourceLifecycle.length);  // Still 100 (original untouched)
console.log(filtered.length);                 // e.g., 75 (read-only events excluded)
```

---

## No Backend Schema Changes

**GraphQL Schema**: No modifications to `gql/lifecycle.graphql`

**Ent Entities**: No changes to `ent/schema/*.go`

**Database**: No migrations required

**Rationale**: Feature is purely presentational. All audit data continues to be stored and queryable; filtering is a client-side view preference.

---

*Data model complete: 2025-10-09*
