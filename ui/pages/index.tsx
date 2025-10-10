import { graphql } from '@/modules/gql'
import { useQuery } from '@tanstack/react-query'
import request from 'graphql-request'
import Head from 'next/head'
import { Sidebar } from '@/components/Sidebar'
import { Card, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'

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

export default function Home() {

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
        </div>
      </Sidebar>
    </>
  )
}
