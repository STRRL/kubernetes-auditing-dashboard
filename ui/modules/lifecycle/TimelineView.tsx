import React from 'react';
import { TimestampDisplay } from '../common/TimestampDisplay';
import { DiffViewer } from './DiffViewer';

type EventType = 'CREATE' | 'UPDATE' | 'DELETE' | 'GET';

interface LifecycleEvent {
  id: string;
  type: EventType;
  timestamp: string;
  user: string;
  resourceState: string;
  diff?: {
    added?: string | null;
    removed?: string | null;
    modified: Array<{
      path: string;
      oldValue: string;
      newValue: string;
    }>;
  } | null;
}

interface TimelineViewProps {
  events: LifecycleEvent[];
  allEvents?: LifecycleEvent[]; // Full unfiltered events for previousState lookup
}

const eventTypeConfig: Record<EventType, { color: string; label: string }> = {
  CREATE: { color: 'bg-green-500', label: 'CREATE' },
  UPDATE: { color: 'bg-blue-500', label: 'UPDATE' },
  DELETE: { color: 'bg-red-500', label: 'DELETE' },
  GET: { color: 'bg-gray-500', label: 'GET' },
};

export const TimelineView: React.FC<TimelineViewProps> = ({ events, allEvents }) => {
  // Use allEvents for previousState lookup if provided, otherwise fall back to events
  const eventsForLookup = allEvents || events;

  return (
    <div className="space-y-4">
      {events.map((event, index) => {
        const config = eventTypeConfig[event.type];
        const isLast = index === events.length - 1;

        // Find the previous state from the full unfiltered events array
        const eventIndexInAll = eventsForLookup.findIndex(e => e.id === event.id);
        const previousState = eventIndexInAll < eventsForLookup.length - 1
          ? eventsForLookup[eventIndexInAll + 1].resourceState
          : undefined;

        return (
          <div key={event.id} className="flex gap-4">
            <div className="flex flex-col items-center">
              <div
                className={`w-10 h-10 rounded-full ${config.color} flex items-center justify-center`}
              >
              </div>
              {!isLast && <div className="w-0.5 flex-1 bg-gray-300 mt-2"></div>}
            </div>

            <div className="flex-1 pb-8">
              <div className="bg-white border border-gray-200 rounded-lg shadow-sm p-4">
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-2">
                    <span className={`px-3 py-1 rounded-full text-sm font-semibold text-white ${config.color}`}>
                      {config.label}
                    </span>
                    <span className="text-sm text-gray-600">by {event.user}</span>
                  </div>
                  <TimestampDisplay timestamp={event.timestamp} />
                </div>

                {event.type === 'UPDATE' && event.diff && (
                  <DiffViewer
                    diff={event.diff}
                    currentState={event.resourceState}
                    previousState={previousState}
                  />
                )}

                {event.type === 'CREATE' && event.resourceState && (
                  <details className="mt-2">
                    <summary className="cursor-pointer text-sm text-blue-600 hover:text-blue-800">
                      View resource state
                    </summary>
                    <pre className="mt-2 p-2 bg-gray-50 rounded text-xs overflow-x-auto">
                      {JSON.stringify(JSON.parse(event.resourceState), null, 2)}
                    </pre>
                  </details>
                )}
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};
