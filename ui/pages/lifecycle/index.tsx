import React from 'react';
import { useRouter } from 'next/router';
import { useQuery } from '@tanstack/react-query';
import request from 'graphql-request';
import Head from 'next/head';
import Link from 'next/link';
import { graphql } from '@/modules/gql';
import { TimelineView } from '@/modules/lifecycle/TimelineView';
import { EmptyState } from '@/modules/lifecycle/EmptyState';

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
      <div className="drawer drawer-mobile">
        <input id="drawer-indicator" type="checkbox" className="drawer-toggle" />
        <div className="drawer-content flex flex-col p-4">
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
                {data.resourceLifecycle.length === 0 ? (
                  <EmptyState />
                ) : (
                  <TimelineView events={data.resourceLifecycle} />
                )}
              </>
            )}
          </div>
        </div>

        <div className="drawer-side">
          <label htmlFor="drawer-indicator" className="drawer-overlay"></label>
          <ul className="menu p-4 w-80 bg-base-100 text-base-content">
            <li>
              <Link href="/">Home</Link>
            </li>
            <li>
              <Link href="/events">Recent Changes</Link>
            </li>
            <li className="font-bold">
              <Link href="/lifecycle">Resource Lifecycle</Link>
            </li>
          </ul>
        </div>
      </div>
    </>
  );
}
