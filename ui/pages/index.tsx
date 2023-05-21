import { graphql } from '@/modules/gql'
import { useQuery } from '@tanstack/react-query'
import request from 'graphql-request'
import Head from 'next/head'
import { useState } from 'react'

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
    }
  }
`)

export default function Home() {

  const [pageSize, setPageSize] = useState(10)
  const [endCursor, setEndCursor] = useState<String | null>(null)
  const eventsCountQuery = useQuery({
    queryKey: ['eventsCount'],
    queryFn: async () => request('/api/query', eventsCountDocumentations)
  })

  const eventsListQuery = useQuery({
    queryKey: ['eventsList', { first: pageSize, after: endCursor }],
    queryFn: async ({ queryKey }) => request('/api/query', eventsListDocumentations, {
      first: pageSize,
      after: endCursor
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
            <div className='overflow-x-auto'>
              <div className='grid grid-cols-3'>
                <div>th <div>h1</div></div>
                <div>th</div>
                <div>th</div>
                <div>th</div>
                <div>th</div>
                <div>th</div>
                <div>th</div>
                <div>th</div>
                <div>th</div>
                <div>th</div>
              </div>
            </div>
          </div>
          <div className='m-4'>
            <div className="overflow-x-auto">
              <table className="table w-full">
                {/* head*/}
                <thead>
                  <tr>
                    <th></th>
                    <th>Time</th>
                    <th>Verb</th>
                    <th>Resource</th>
                  </tr>
                </thead>
                <tbody>
                  {/* row 1 */}
                  <tr>
                    <th>1</th>
                    <td>Cy Ganderton</td>
                    <td>Quality Control Specialist</td>
                    <td>Blue</td>
                  </tr>
                  {/* row 2 */}
                  <tr>
                    <th>2</th>
                    <td>Cy Ganderton</td>
                    <td>Quality Control Specialist</td>
                    <td>Blue</td>
                    <td colSpan={1}>
                      <div tabIndex={0} className="collapse border border-base-300 bg-base-100 rounded-box">
                        <div className="collapse-title text-xl font-medium">
                          Focus me to see content
                        </div>
                        <div className="collapse-content">
                          <p>tabIndex={0} attribute is necessary to make the div focusable</p>
                        </div>
                      </div>
                    </td>
                  </tr>
                  {/* row 3 */}
                  <tr>
                    <th>3</th>
                    <td>Brice Swyre</td>
                    <td>Tax Accountant</td>
                    <td>Red</td>
                  </tr>
                </tbody>
              </table>
            </div>
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
