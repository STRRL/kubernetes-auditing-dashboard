import React, { useState } from 'react';
import { useRouter } from 'next/router';
import { useQuery } from '@tanstack/react-query';
import request from 'graphql-request';
import Head from 'next/head';
import Link from 'next/link';
import { graphql } from '@/modules/gql';
import { TimelineView } from '@/modules/lifecycle/TimelineView';
import { EmptyState } from '@/modules/lifecycle/EmptyState';
import { Sidebar } from '@/components/Sidebar';
import { Switch } from '@/components/ui/switch';

const getResourceLifecycleQuery = graphql(/* GraphQL */ `
  query GetResourceLifecycle(
    $apiGroup: String!
    $version: String!
    $kind: String!
    $namespace: String
    $name: String!
  ) {
    resourceLifecycle(
      apiGroup: $apiGroup
      version: $version
      kind: $kind
      namespace: $namespace
      name: $name
    ) {
      id
      type
      timestamp
      user
      resourceState
      previousState
      diff {
        added
        removed
        modified {
          path
          oldValue
          newValue
        }
      }
    }
  }
`);

export default function LifecyclePage() {
  const router = useRouter();
  const { group, version, kind, namespace, name } = router.query;
  const [showReadOnlyEvents, setShowReadOnlyEvents] = useState(false);

  const apiGroup = (group as string) || '';
  const apiVersion = version as string;
  const resourceKind = kind as string;
  const resourceNamespace = namespace as string | undefined;
  const resourceName = name as string;

  const isValid = apiVersion && resourceKind && resourceName;

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ['resourceLifecycle', { apiGroup, apiVersion, resourceKind, resourceNamespace, resourceName }],
    queryFn: async () => {
      if (!isValid) throw new Error('Invalid URL parameters');

      return request('/api/query', getResourceLifecycleQuery, {
        apiGroup,
        version: apiVersion,
        kind: resourceKind,
        namespace: resourceNamespace || null,
        name: resourceName,
      });
    },
    enabled: !!isValid,
  });

  const filteredEvents = React.useMemo(() => {
    if (!data?.resourceLifecycle) return [];
    if (showReadOnlyEvents) return data.resourceLifecycle;

    const readOnlyEventTypes = ['GET', 'LIST', 'WATCH'];
    return data.resourceLifecycle.filter(
      (event) => !readOnlyEventTypes.includes(event.type)
    );
  }, [data?.resourceLifecycle, showReadOnlyEvents]);

  if (!isValid) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-red-600 mb-2">Invalid URL</h2>
          <p className="text-gray-600">
            Please navigate to this page from the Recent Changes list.
          </p>
        </div>
      </div>
    );
  }

  const displayNamespace = resourceNamespace ? `${resourceNamespace}/` : '';
  const scopeLabel = resourceNamespace ? '' : ' (cluster-scoped)';
  const displayName = `${displayNamespace}${resourceName}${scopeLabel}`;

  const resourceTitle = `${apiGroup ? apiGroup + '/' : ''}${apiVersion} ${resourceKind}`;

  return (
    <>
      <Head>
        <title>Lifecycle: {displayName} | Kubernetes Auditing Dashboard</title>
      </Head>
      <Sidebar>
        <div className="p-4">
          <div className="m-4">
            <h2 className="text-4xl font-bold text-gray-800">Resource Lifecycle</h2>
            <div className="mt-2 text-lg text-gray-600">
              <span className="font-mono text-sm bg-gray-100 px-2 py-1 rounded">
                {resourceTitle}
              </span>
              <span className="mx-2">/</span>
              <span className="font-semibold">{displayName}</span>
            </div>
          </div>

          <div className="m-4 flex items-center gap-3 p-3 bg-gray-50 rounded-lg border border-gray-200">
            <Switch
              id="show-readonly"
              checked={showReadOnlyEvents}
              onCheckedChange={setShowReadOnlyEvents}
            />
            <label
              htmlFor="show-readonly"
              className="text-sm font-medium text-gray-700 cursor-pointer"
            >
              Show read-only events (GET, LIST, WATCH)
            </label>
          </div>

          <div className="m-4">
            {isLoading && (
              <div className="flex items-center justify-center min-h-[400px]">
                <div className="text-center">
                  <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
                  <p className="mt-4 text-gray-600">Loading lifecycle events...</p>
                </div>
              </div>
            )}

            {isError && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                <h3 className="text-lg font-semibold text-red-800 mb-2">Error Loading Events</h3>
                <p className="text-red-700">
                  {error instanceof Error ? error.message : 'An unknown error occurred'}
                </p>
              </div>
            )}

            {!isLoading && !isError && data?.resourceLifecycle && (
              <>
                {filteredEvents.length === 0 ? (
                  <EmptyState />
                ) : (
                  <TimelineView events={filteredEvents} />
                )}
              </>
            )}
          </div>
        </div>
      </Sidebar>
    </>
  );
}
