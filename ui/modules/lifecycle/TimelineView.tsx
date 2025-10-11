import { getVerbColor } from '@/lib/verb-colors';
import React from 'react';
import { TimestampDisplay } from '../common/TimestampDisplay';
import { DiffViewer } from './DiffViewer';

interface LifecycleEvent {
  id: string;
  type: string;
  timestamp: string;
  user: string;
  resourceState: string;
  previousState?: string | null;
}

interface TimelineViewProps {
  events: LifecycleEvent[];
}

export const TimelineView: React.FC<TimelineViewProps> = ({ events }) => {
  return (
    <div className="space-y-4">
      {events.map((event, index) => {
        const verbColor = getVerbColor(event.type);
        const isLast = index === events.length - 1;

        return (
          <div key={event.id} className="flex gap-4">
            <div className="flex flex-col items-center">
              <div
                className={`w-10 h-10 rounded-full ${verbColor} flex items-center justify-center`}
              >
              </div>
              {!isLast && <div className="w-0.5 flex-1 bg-gray-300 mt-2"></div>}
            </div>

            <div className="flex-1 pb-8">
              <div className="bg-white border border-gray-200 rounded-lg shadow-sm p-4">
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-2">
                    <span className={`px-3 py-1 rounded-full text-sm font-semibold text-white ${verbColor}`}>
                      {event.type ? event.type.toUpperCase() : 'UNKNOWN'}
                    </span>
                    <span className="text-sm text-gray-600">by {event.user}</span>
                  </div>
                  <TimestampDisplay timestamp={event.timestamp} />
                </div>

                {(event.type === 'update' || event.type === 'patch' || event.type === 'create') && event.previousState && (
                  <DiffViewer
                    currentState={event.resourceState}
                    previousState={event.previousState}
                  />
                )}

                {(event.type === 'create') && (
                  <div>
                    <DiffViewer
                      currentState={event.resourceState}
                      previousState={'{}'}
                    />
                  </div>
                )}
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};
