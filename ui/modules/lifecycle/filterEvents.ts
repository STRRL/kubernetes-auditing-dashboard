/**
 * Utility functions for filtering lifecycle events
 * Implements client-side filtering for read-only Kubernetes audit events
 */

/**
 * Read-only event types that should be hidden by default
 * Based on Kubernetes audit verbs: get, list, watch
 */
export const READ_ONLY_EVENTS = ['get', 'list', 'watch'] as const;

/**
 * Event type from GraphQL codegen output
 */
type EventType = 'CREATE' | 'UPDATE' | 'DELETE' | 'GET' | 'LIST' | 'WATCH' | 'PATCH';

/**
 * Lifecycle event interface (matches GraphQL schema)
 */
export interface LifecycleEvent {
  id: string;
  type: EventType;
  timestamp: string;
  user: string;
  resourceState: any;
  diff?: any;
}

/**
 * Filters lifecycle events based on read-only status
 *
 * @param events - Array of lifecycle events to filter
 * @param hideReadOnly - Whether to hide read-only events (get, list, watch)
 * @returns Filtered array of events, preserving chronological order
 *
 * @example
 * const events = [
 *   { id: '1', type: 'CREATE', ... },
 *   { id: '2', type: 'GET', ... },
 *   { id: '3', type: 'UPDATE', ... }
 * ];
 * const filtered = filterEvents(events, true); // Returns only CREATE and UPDATE
 */
export function filterEvents(
  events: LifecycleEvent[],
  hideReadOnly: boolean
): LifecycleEvent[] {
  if (!hideReadOnly) {
    return events;
  }

  return events.filter(event => {
    const eventType = event.type.toLowerCase();
    return !READ_ONLY_EVENTS.includes(eventType as any);
  });
}
