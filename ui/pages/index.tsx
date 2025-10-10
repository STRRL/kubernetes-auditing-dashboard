'use client';
import { graphql } from '@/modules/gql'
import { useQuery } from '@tanstack/react-query'
import request from 'graphql-request'
import Head from 'next/head'
import { useState } from 'react'
import { useParams } from 'next/navigation'
import { Sidebar } from '@/components/Sidebar'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'

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

const nonGetEventsCountGQLDoc = graphql(/* GraphQL */ `
  query eventsCountNonGet{
    auditEvents(where: {
        verb: "watch"
      })
    {
      totalCount
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
    queryFn: async () => request('/api/query', eventsCountDocumentations),
    // refresh every 15 seconds
    refetchInterval: 15000,
  })

  const nonGetEventsCountQuery = useQuery({
    queryKey: ['eventsCountNonGet'],
    queryFn: async () => request('/api/query', nonGetEventsCountGQLDoc),
    // refresh every 15 seconds
    refetchInterval: 15000,
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
      <Sidebar>
        <div className="p-4">
          <div className='m-4'>
            <h2 className='text-4xl font-bold'>Discovery</h2>
          </div>
          <div className='flex gap-4 mx-4'>
            <Card className="w-1/4">
              <CardHeader>
                <CardDescription>Total Events</CardDescription>
                <CardTitle className="text-3xl">{eventsCountQuery.data?.auditEvents.totalCount.toLocaleString()}</CardTitle>
              </CardHeader>
            </Card>
            <Card className="w-1/4">
              <CardHeader>
                <CardDescription>Non-Get Events</CardDescription>
                <CardTitle className="text-3xl">{nonGetEventsCountQuery.data?.auditEvents.totalCount.toLocaleString()}</CardTitle>
              </CardHeader>
            </Card>
          </div>
          <div className='m-4'>
            <div className="overflow-x-auto">
            </div>
          </div>
        </div>
      </Sidebar>
    </>
  )
}
