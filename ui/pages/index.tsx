import { graphql } from '@/modules/gql'
import { useQuery } from '@tanstack/react-query'
import request from 'graphql-request'
import Head from 'next/head'

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

const eventsListDocumentations = graphql(/* GraphQL */ `
query eventsList($first: Int, $after: Cursor){
  auditEvents(first: $first, after: $after){
    totalCount
    pageInfo{
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
    edges{
      node{
        id
        level
        stage
        verb
        useragent
        resource
        namespace
        name
        stagetimestamp
      }
    }
  }}
`)

export default function Home() {
  const eventsCountQuery = useQuery({
    queryKey: ['eventsCount'],
    queryFn: async () => request('/api/query', eventsCountDocumentations)
  })

  const eventsListQuery = useQuery({
    queryKey: ['eventsList', { first: 10, after: null }],
    queryFn: async ({ queryKey }) => request('/api/query', eventsListDocumentations, {
      first: 20,
      after: null
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

          <div className='flex'>
            <div className="flex stats shadow mx-4 w-1/4">
              <div className="stat">
                <div className="stat-title">Total Events</div>
                <div className="stat-value">{eventsCountQuery.data?.auditEvents.totalCount.toLocaleString()}</div>
                <div className="stat-desc"></div>
              </div>
            </div>
          </div>

          <div>
            {eventsListQuery.data?.auditEvents.edges?.map((edge) => {
              return (
                <div key={`event-${edge?.node?.id}`} className='shadow m-4'>
                  <div className="card card-side bg-base-100 shadow-xl">
                    <div className="card-body">
                      <h2 className="card-title">{edge?.node?.verb.toUpperCase()}: {edge?.node?.resource} {edge?.node?.namespace}/{edge?.node?.name}</h2>
                      <p>UserAgent: {edge?.node?.useragent}</p>
                      <div className='text-right'>
                        <p className=''>StagedTimestamp: {edge?.node?.stagetimestamp}</p>
                      </div>
                      <div className="card-actions justify-end">
                      </div>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>

        </div>

        <div className="drawer-side">
          <label htmlFor="drawer-indicator" className="drawer-overlay"></label>
          <ul className="menu p-4 w-80 bg-base-100 text-base-content">
            <li><a>Auditing Event List</a></li>
          </ul>
        </div>
      </div>
    </>
  )
}
