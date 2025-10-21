import { getVerbColor } from '@/lib/verb-colors';
import React from 'react';
import { TimestampDisplay } from '../common/TimestampDisplay';
import { DiffViewer } from './DiffViewer';

interface LifecycleEvent {
  id: string;
  type: string;
  timestamp: string;
  user: string;
  userAgent: string;
  resourceState: string;
  previousState?: string | null;
  diff?: {
    added?: any;
    removed?: any;
    modified: Array<{
      path: string;
      oldValue: any;
      newValue: any;
    }>;
  } | null;
}

interface TimelineViewProps {
  events: LifecycleEvent[];
}

type EventType = 'CREATE' | 'UPDATE' | 'DELETE' | 'GET';

const eventTypeConfig: Record<EventType, { color: string; label: string }> = {
  CREATE: { color: 'bg-green-500', label: 'CREATE' },
  UPDATE: { color: 'bg-blue-500', label: 'UPDATE' },
  DELETE: { color: 'bg-red-500', label: 'DELETE' },
  GET: { color: 'bg-gray-500', label: 'GET' },
};

const formatUserAgent = (userAgent: string) => {
  const wellKnownComponents = {
    'kubelet': { name: 'Kubelet', color: 'bg-blue-500' },
    'kube-apiserver': { name: 'API Server', color: 'bg-green-500' },
    'kube-controller-manager': { name: 'Controller Manager', color: 'bg-yellow-500' },
    'kube-scheduler': { name: 'Scheduler', color: 'bg-purple-500' },
    'kube-proxy': { name: 'Kube Proxy', color: 'bg-red-500' },
    'storage-provisioner': { name: 'Minikube Storage Provisioner', color: 'bg-indigo-500' },
    'kubectl': { name: 'Kubectl', color: 'bg-teal-500' }
  };

  for (const [key, value] of Object.entries(wellKnownComponents)) {
    if (userAgent.toLowerCase().includes(key)) {
      return { name: value.name, color: value.color, isKnown: true };
    }
  }

  return { name: userAgent, color: '', isKnown: false };
};

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
                  </div>
                  <TimestampDisplay timestamp={event.timestamp} />
                </div>

                <div className="grid grid-cols-2 gap-4 mt-3 mb-2">
                  <div className="flex flex-col">
                    <span className="text-xs text-gray-500 mb-1">Component</span>
                    {(() => {
                      const component = formatUserAgent(event.userAgent);
                      return component.isKnown ? (
                        <span className={`px-2 py-1 rounded-full text-xs font-semibold text-white ${component.color} inline-block w-fit`}>
                          {component.name}
                        </span>
                      ) : (
                        <span className="text-sm text-gray-700">{component.name}</span>
                      );
                    })()}
                  </div>
                  <div className="flex flex-col">
                    <span className="text-xs text-gray-500 mb-1">RBAC User</span>
                    <span className="text-sm text-gray-700">{event.user}</span>
                  </div>
                </div>

                {(event.type.toUpperCase() === 'UPDATE' || event.type.toUpperCase() === 'PATCH') && event.previousState && (
                  <DiffViewer
                    currentState={event.resourceState}
                    previousState={event.previousState}
                  />
                )}

                {event.type.toUpperCase() === 'CREATE' && (
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
