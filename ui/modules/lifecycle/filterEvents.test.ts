/**
 * Tests for filterEvents utility function
 * Adapted from contract test in specs/002-hide-read-only/contracts/lifecycle-filter.test.tsx
 * Following TDD: These tests MUST FAIL initially (function not yet implemented)
 */

import { describe, it, expect } from '@jest/globals';
import { filterEvents, READ_ONLY_EVENTS } from './filterEvents';

/**
 * Type definitions (matching GraphQL codegen output)
 */
type EventType = 'CREATE' | 'UPDATE' | 'DELETE' | 'GET' | 'LIST' | 'WATCH' | 'PATCH';

interface LifecycleEvent {
  id: string;
  type: EventType;
  timestamp: string;
  user: string;
  resourceState: any;
  diff?: any;
}

describe('Lifecycle Event Filtering', () => {
  it('should define read-only event types as get, list, watch', () => {
    expect(READ_ONLY_EVENTS).toEqual(['get', 'list', 'watch']);
    expect(READ_ONLY_EVENTS).toHaveLength(3);
  });

  it('should filter out read-only events when hideReadOnly is true', () => {
    const events: LifecycleEvent[] = [
      { id: '1', type: 'CREATE', timestamp: '2023-01-01T10:00:00Z', user: 'user1', resourceState: {} },
      { id: '2', type: 'GET', timestamp: '2023-01-02T10:00:00Z', user: 'user2', resourceState: {} },
      { id: '3', type: 'UPDATE', timestamp: '2023-01-03T10:00:00Z', user: 'user3', resourceState: {} },
      { id: '4', type: 'LIST', timestamp: '2023-01-04T10:00:00Z', user: 'user4', resourceState: {} },
      { id: '5', type: 'PATCH', timestamp: '2023-01-05T10:00:00Z', user: 'user5', resourceState: {} },
      { id: '6', type: 'WATCH', timestamp: '2023-01-06T10:00:00Z', user: 'user6', resourceState: {} },
    ];

    const filtered = filterEvents(events, true);

    expect(filtered).toHaveLength(3);
    expect(filtered.map(e => e.type)).toEqual(['CREATE', 'UPDATE', 'PATCH']);
    expect(filtered.map(e => e.id)).toEqual(['1', '3', '5']);
  });

  it('should return all events when hideReadOnly is false', () => {
    const events: LifecycleEvent[] = [
      { id: '1', type: 'CREATE', timestamp: '2023-01-01T10:00:00Z', user: 'user1', resourceState: {} },
      { id: '2', type: 'GET', timestamp: '2023-01-02T10:00:00Z', user: 'user2', resourceState: {} },
      { id: '3', type: 'DELETE', timestamp: '2023-01-03T10:00:00Z', user: 'user3', resourceState: {} },
    ];

    const filtered = filterEvents(events, false);

    expect(filtered).toHaveLength(3);
    expect(filtered).toEqual(events);
  });

  it('should preserve chronological order after filtering', () => {
    const events: LifecycleEvent[] = [
      { id: '1', type: 'UPDATE', timestamp: '2023-01-05T10:00:00Z', user: 'user1', resourceState: {} },
      { id: '2', type: 'GET', timestamp: '2023-01-04T10:00:00Z', user: 'user2', resourceState: {} },
      { id: '3', type: 'CREATE', timestamp: '2023-01-03T10:00:00Z', user: 'user3', resourceState: {} },
      { id: '4', type: 'LIST', timestamp: '2023-01-02T10:00:00Z', user: 'user4', resourceState: {} },
      { id: '5', type: 'PATCH', timestamp: '2023-01-01T10:00:00Z', user: 'user5', resourceState: {} },
    ];

    const filtered = filterEvents(events, true);

    expect(filtered.map(e => e.id)).toEqual(['1', '3', '5']);
    expect(filtered[0].timestamp > filtered[1].timestamp).toBe(true);
    expect(filtered[1].timestamp > filtered[2].timestamp).toBe(true);
  });

  it('should return empty array when input is empty', () => {
    const events: LifecycleEvent[] = [];

    const filteredHide = filterEvents(events, true);
    const filteredShow = filterEvents(events, false);

    expect(filteredHide).toEqual([]);
    expect(filteredShow).toEqual([]);
  });

  it('should return empty array when all events are read-only and hideReadOnly is true', () => {
    const events: LifecycleEvent[] = [
      { id: '1', type: 'GET', timestamp: '2023-01-01T10:00:00Z', user: 'user1', resourceState: {} },
      { id: '2', type: 'LIST', timestamp: '2023-01-02T10:00:00Z', user: 'user2', resourceState: {} },
      { id: '3', type: 'WATCH', timestamp: '2023-01-03T10:00:00Z', user: 'user3', resourceState: {} },
    ];

    const filtered = filterEvents(events, true);

    expect(filtered).toEqual([]);
    expect(filtered).toHaveLength(0);
  });

  it('should return all events when none are read-only', () => {
    const events: LifecycleEvent[] = [
      { id: '1', type: 'CREATE', timestamp: '2023-01-01T10:00:00Z', user: 'user1', resourceState: {} },
      { id: '2', type: 'UPDATE', timestamp: '2023-01-02T10:00:00Z', user: 'user2', resourceState: {} },
      { id: '3', type: 'DELETE', timestamp: '2023-01-03T10:00:00Z', user: 'user3', resourceState: {} },
      { id: '4', type: 'PATCH', timestamp: '2023-01-04T10:00:00Z', user: 'user4', resourceState: {} },
    ];

    const filtered = filterEvents(events, true);

    expect(filtered).toEqual(events);
    expect(filtered).toHaveLength(4);
  });

  it('should handle case-insensitive event types', () => {
    const events: LifecycleEvent[] = [
      { id: '1', type: 'CREATE', timestamp: '2023-01-01T10:00:00Z', user: 'user1', resourceState: {} },
      { id: '2', type: 'GET', timestamp: '2023-01-02T10:00:00Z', user: 'user2', resourceState: {} },
      { id: '3', type: 'get' as EventType, timestamp: '2023-01-03T10:00:00Z', user: 'user3', resourceState: {} },
    ];

    const filtered = filterEvents(events, true);

    expect(filtered).toHaveLength(1);
    expect(filtered[0].id).toBe('1');
  });

  it('should not mutate the original events array', () => {
    const events: LifecycleEvent[] = [
      { id: '1', type: 'CREATE', timestamp: '2023-01-01T10:00:00Z', user: 'user1', resourceState: {} },
      { id: '2', type: 'GET', timestamp: '2023-01-02T10:00:00Z', user: 'user2', resourceState: {} },
    ];

    const originalLength = events.length;
    const originalFirstId = events[0].id;

    const filtered = filterEvents(events, true);

    expect(events).toHaveLength(originalLength);
    expect(events[0].id).toBe(originalFirstId);
    expect(events[1].type).toBe('GET');
    expect(filtered).not.toBe(events);
    expect(filtered).toHaveLength(1);
  });

  it('should satisfy acceptance criterion: hide read-only events by default', () => {
    const events: LifecycleEvent[] = [
      { id: '1', type: 'CREATE', timestamp: '2023-01-01T10:00:00Z', user: 'admin', resourceState: {} },
      { id: '2', type: 'GET', timestamp: '2023-01-02T10:00:00Z', user: 'viewer', resourceState: {} },
      { id: '3', type: 'UPDATE', timestamp: '2023-01-03T10:00:00Z', user: 'admin', resourceState: {} },
    ];

    const defaultFiltered = filterEvents(events, true);

    expect(defaultFiltered.map(e => e.type)).toEqual(['CREATE', 'UPDATE']);
    expect(defaultFiltered.some(e => e.type === 'GET')).toBe(false);
  });
});
