'use client';
import { graphql } from '@/modules/gql'
import { useQuery } from '@tanstack/react-query'
import request from 'graphql-request'
import Head from 'next/head'
import { useState } from 'react'
import { useParams } from 'next/navigation';

const eventsCountDocumentations = graphql(/* GraphQL */ `
  query eventsCount{
    auditEvents{
      totalCount
      pageInfo{
        hasNextPage
        hasPreviousPage
        startCursor
        endCursor
      }
    }
  }
`)

const completedRequestResponseAuditEventsDocumentations = graphql(/* GraphQL */ `
  query completedRequestResponseAuditEvents($page: Int, $pageSize: Int){
    completedRequestResponseAuditEvents(page: $page, pageSize: $pageSize) {
    total
    page
    pageSize
    totalPages
    hasNextPage
    hasPreviousPage
    rows {
      id
      level
      stage
      verb
      useragent
      resource
      namespace
      name
      stagetimestamp
      apigroup
      apiversion
    }
  }
  }
`)
const resourceName = (apiGroup: string | undefined, apiVersion: string, resourceName: string) => {
  return `${apiGroup ? apiGroup + "/" : ""}${apiVersion} ${resourceName}`
}
const namespacedName = (namespace: string | undefined, name: string) => {
  if (namespace === undefined || namespace === "") {
    return `${name} (cluster-scoped)`
  }
  return `${namespace}/${name}`
}
export default function Home() {

  const params = useParams();
  console.log("params", params)

  const [page, setPage] = useState(0)
  const [pageSize, setPageSize] = useState(15)

  const eventsCountQuery = useQuery({
    queryKey: ['eventsCount'],
    queryFn: async () => request('/api/query', eventsCountDocumentations)
  })

  const eventsListQuery = useQuery({
    queryKey: ['eventsList', { page: page, pageSize: pageSize }],
    queryFn: async ({ queryKey }) => request('/api/query', completedRequestResponseAuditEventsDocumentations, {
      page: page,
      pageSize: pageSize,
    })
  })

  return (
    <>
      <Head>
        <title>Kubernetes Auditing Dashboard</title>
      </Head>
      <div className="drawer drawer-mobile">
        <input id="drawer-indicator" type="checkbox" className="drawer-toggle" />
        <div className="drawer-content flex flex-col p-4">
          <div className='m-4'>
            <h2 className='text-4xl'>Discovery</h2>
          </div>
          <div className='flex'>
            <div className="flex stats shadow mx-4 w-1/4">
              <div className="stat">
                <div className="stat-title">Total Events</div>
                <div className="stat-value">{eventsCountQuery.data?.auditEvents.totalCount.toLocaleString()}</div>
                <div className="stat-desc"></div>
              </div>
            </div>
          </div>
          <div className='m-4'>
            <div className="overflow-x-auto">
            </div>
          </div>
        </div>

        <div className="drawer-side">
          <label htmlFor="drawer-indicator" className="drawer-overlay"></label>
          <ul className="menu p-4 w-80 bg-base-100 text-base-content">
            <li><a href='/'>Home</a></li>
            <li><a href='/events'>Auditing Events</a></li>
            <li><a href='/lifecycle'>Resource Lifecycle</a></li>
          </ul>
        </div>
      </div>
    </>
  )
}
