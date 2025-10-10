# Quickstart: Hide Read-Only Events by Default

## Overview

This quickstart guide validates the read-only event filtering feature through manual testing. Follow these steps **after implementation** to verify all acceptance criteria from the specification.

**Prerequisites**:
- Kubernetes cluster with audit webhook configured (see `script/kube-apiserver-config/`)
- Application running (`make dev` from repository root)
- Browser with DevTools (Chrome/Firefox recommended)
- Some lifecycle events already ingested (both mutating and read-only)

**Estimated Time**: 10 minutes

---

## Test Scenario 1: Default Filter State

**Acceptance Criterion** (spec.md:53):
> **Given** a user navigates to a resource lifecycle page, **When** the page loads, **Then** all events except read-only operations (get, list, watch) are displayed, and the "Hide read-only events" toggle is enabled by default.

### Steps

1. **Navigate to lifecycle page**:
   ```
   http://localhost:3000/lifecycle?group=apps&version=v1&kind=Deployment&namespace=default&name=nginx
   ```
   *(Adjust resource identifiers to match available resources in your cluster)*

2. **Verify toggle state**:
   - [ ] Toggle labeled "Hide read-only events" is visible in page header (near resource title)
   - [ ] Toggle is in the "enabled/checked" state by default

3. **Verify filtered events**:
   - [ ] Timeline displays CREATE, UPDATE, DELETE, PATCH events
   - [ ] Timeline does NOT display GET, LIST, WATCH events
   - [ ] Events are in chronological order (newest first)

4. **Inspect localStorage** (Browser DevTools → Application → Local Storage):
   ```javascript
   localStorage.getItem('lifecycle-hide-readonly')
   // Expected: "true"
   ```

**Expected Result**: ✅ Read-only events hidden by default, toggle enabled

---

## Test Scenario 2: Toggle Shows All Events

**Acceptance Criterion** (spec.md:54):
> **Given** the lifecycle page is displaying filtered events (read-only hidden), **When** the user disables the "Hide read-only events" toggle, **Then** all events including get/list/watch operations become visible.

### Steps

1. **Start from Test Scenario 1** (toggle enabled, read-only events hidden)

2. **Click the toggle** to disable filtering

3. **Verify toggle state**:
   - [ ] Toggle is now in the "disabled/unchecked" state

4. **Verify all events visible**:
   - [ ] Timeline now includes GET, LIST, WATCH events
   - [ ] Previous CREATE, UPDATE, DELETE events still present
   - [ ] Total event count increased (more events visible)
   - [ ] Chronological order preserved

5. **Inspect localStorage**:
   ```javascript
   localStorage.getItem('lifecycle-hide-readonly')
   // Expected: "false"
   ```

**Expected Result**: ✅ All events (including read-only) now visible

---

## Test Scenario 3: Toggle Hides Events Again

**Acceptance Criterion** (spec.md:55):
> **Given** the lifecycle page is showing all events (read-only visible), **When** the user enables the "Hide read-only events" toggle, **Then** the page returns to showing only non-read-only events.

### Steps

1. **Start from Test Scenario 2** (toggle disabled, all events visible)

2. **Click the toggle** to re-enable filtering

3. **Verify toggle state**:
   - [ ] Toggle is now in the "enabled/checked" state again

4. **Verify filtered events**:
   - [ ] GET, LIST, WATCH events are hidden again
   - [ ] Only CREATE, UPDATE, DELETE, PATCH events remain
   - [ ] Timeline reverted to initial filtered state

5. **Inspect localStorage**:
   ```javascript
   localStorage.getItem('lifecycle-hide-readonly')
   // Expected: "true"
   ```

**Expected Result**: ✅ Read-only events hidden again, matching initial state

---

## Test Scenario 4: Preference Persistence Across Navigation

**Acceptance Criterion** (spec.md:56):
> **Given** a user has set their preference for showing/hiding read-only events, **When** they navigate to a different resource lifecycle page or return later in a new browser session, **Then** the preference is retained and applied automatically.

### Steps

1. **Set preference to show all events**:
   - Disable the "Hide read-only events" toggle (from Test Scenario 2)
   - Verify localStorage: `lifecycle-hide-readonly` = `"false"`

2. **Navigate to different resource**:
   ```
   http://localhost:3000/lifecycle?group=&version=v1&kind=ConfigMap&namespace=kube-system&name=coredns
   ```
   *(Any different resource from Step 1)*

3. **Verify preference persisted**:
   - [ ] Page loads with toggle in "disabled/unchecked" state
   - [ ] All events (including read-only) are visible immediately
   - [ ] No momentary flash of filtered view

4. **Refresh the page** (F5 or Cmd+R)

5. **Verify preference after refresh**:
   - [ ] Toggle still in "disabled/unchecked" state
   - [ ] All events still visible

6. **Close browser completely**, then reopen

7. **Navigate to lifecycle page** (either resource)

8. **Verify preference after browser restart**:
   - [ ] Toggle still in "disabled/unchecked" state
   - [ ] All events still visible
   - [ ] localStorage value: `lifecycle-hide-readonly` = `"false"`

**Expected Result**: ✅ Preference persists across navigation, refresh, and browser restart

---

## Test Scenario 5: Filtered Empty State

**Acceptance Criterion** (spec.md:59):
> What happens when a resource has only read-only events (no other operations)? The page should display the message "No events found. Toggle 'Hide read-only events' to see read-only operations."

### Steps

1. **Create a resource that only has read-only events**:
   ```bash
   # In your Kubernetes cluster, create a ConfigMap that's never modified
   kubectl create configmap test-readonly -n default --from-literal=key=value
   # Then only read it (don't update/delete)
   kubectl get configmap test-readonly -n default -o yaml
   ```

2. **Navigate to its lifecycle page**:
   ```
   http://localhost:3000/lifecycle?group=&version=v1&kind=ConfigMap&namespace=default&name=test-readonly
   ```

3. **Verify toggle is enabled** (hiding read-only events)

4. **Verify empty state message**:
   - [ ] Empty state displayed (no timeline events)
   - [ ] Message reads: **"No events found. Toggle 'Hide read-only events' to see read-only operations."**
   - [ ] Message styling matches existing empty state component

5. **Disable the toggle** to show read-only events

6. **Verify read-only events now visible**:
   - [ ] Timeline displays GET events for the ConfigMap
   - [ ] Empty state message is gone

**Expected Result**: ✅ Helpful guidance message when all events are filtered out

---

## Test Scenario 6: Truly Empty State

**Edge Case** (spec.md:61):
> What happens if there are no events at all? The existing empty state should be displayed.

### Steps

1. **Navigate to a resource with NO audit events** (newly created in a cluster without audit webhook):
   - If audit webhook is always enabled, skip this test (not applicable)
   - Alternatively, fabricate by navigating to invalid resource identifier

2. **Verify toggle state**:
   - [ ] Toggle is enabled (hiding read-only events by default)

3. **Verify empty state**:
   - [ ] Standard empty state message displayed
   - [ ] NO mention of toggling filter (because no events exist at all)

**Expected Result**: ✅ Standard empty state shown when truly no events (not filtered empty)

---

## Test Scenario 7: Keyboard Accessibility

**Accessibility Requirement** (Constitutional Principle II: Quality standards):

### Steps

1. **Navigate to lifecycle page**

2. **Use Tab key** to focus on toggle

3. **Verify focus indicator**:
   - [ ] Toggle has visible focus outline/ring

4. **Press Space or Enter** to toggle state

5. **Verify toggle responds**:
   - [ ] State changes (enabled ↔ disabled)
   - [ ] Events filter/unfilter accordingly

6. **Repeat toggle via keyboard**:
   - [ ] Consistently responds to Space/Enter

**Expected Result**: ✅ Toggle is fully keyboard accessible

---

## Test Scenario 8: Screen Reader Compatibility

**Accessibility Requirement**:

### Steps (macOS with VoiceOver)

1. **Enable VoiceOver** (Cmd+F5)

2. **Navigate to lifecycle page**

3. **Tab to toggle control**

4. **Verify VoiceOver announcement**:
   - [ ] Announces control type (e.g., "checkbox" or "switch")
   - [ ] Announces label: "Hide read-only events"
   - [ ] Announces state (e.g., "checked" or "unchecked")

5. **Toggle state via keyboard** (Space)

6. **Verify state change announced**:
   - [ ] VoiceOver announces new state

**Expected Result**: ✅ Screen reader users can understand and operate toggle

---

## Lint Validation

**Constitutional Requirement II**: UI changes include lint checks

### Steps

1. **Run lint**:
   ```bash
   cd ui
   npm run lint
   ```

2. **Verify clean output**:
   - [ ] No ESLint errors
   - [ ] No ESLint warnings in modified files:
     - `ui/pages/lifecycle/index.tsx`
     - `ui/modules/lifecycle/TimelineView.tsx`
     - `ui/modules/lifecycle/FilterToggle.tsx`
     - `ui/modules/lifecycle/EmptyState.tsx`
     - `ui/modules/hooks/useLocalStorage.ts`

**Expected Result**: ✅ All lint checks pass

---

## Performance Validation

**Performance Goal** (plan.md:46): Instantaneous filtering (<16ms for 60fps)

### Steps

1. **Navigate to lifecycle page with 100+ events**

2. **Open Browser DevTools → Performance tab**

3. **Start recording**

4. **Toggle filter state**

5. **Stop recording**

6. **Analyze timeline**:
   - [ ] Toggle interaction → re-render completes in <16ms
   - [ ] No noticeable lag or jank

7. **Test with larger dataset** (if available):
   - [ ] 500 events: Still instant (<50ms)
   - [ ] 1000 events: Still acceptable (<100ms)

**Expected Result**: ✅ Filtering is instantaneous, no perceivable delay

---

## Regression Check: Existing Functionality

**Verify no regressions** in lifecycle page features:

- [ ] Timeline diff viewer still expands/collapses changes
- [ ] Timestamps display correctly
- [ ] User field shows audit user
- [ ] Resource state JSON is viewable
- [ ] Page header resource title/namespace unchanged
- [ ] GraphQL error handling unchanged (test by stopping backend)

**Expected Result**: ✅ No existing functionality broken

---

## Cleanup

After validation:

```bash
# Remove test ConfigMap if created
kubectl delete configmap test-readonly -n default

# Reset localStorage to default (if needed)
localStorage.setItem('lifecycle-hide-readonly', 'true')
```

---

## Acceptance Sign-Off

**Feature Ready for Merge** when:

- [x] All 8 test scenarios pass
- [x] Lint validation clean
- [x] Performance within targets
- [x] No regressions detected
- [x] Accessibility verified (keyboard + screen reader)

**Tested By**: ___________
**Date**: ___________
**Browser/Version**: ___________
**Notes**: ___________

---

*Quickstart guide version: 1.0*
