import { graphql } from '@/modules/gql'
import { useQuery } from '@tanstack/react-query'
import request from 'graphql-request'
import Head from 'next/head'
import { use, useEffect, useState } from 'react'
import { useParams } from 'next/navigation';
import { useRouter } from 'next/router'

const moment = require('moment');

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
        if (name === undefined || name === "") {
            return ""
        }
        return `${name} (cluster-scoped)`
    }

    return `${namespace}/${name}`
}

export default function Events() {
    const router = useRouter()

    const [page, setPage] = useState(0)
    const [pageSize, setPageSize] = useState(15)

    useEffect(() => {
        setPage(parseInt(router.query.page as string || '0'))
    }, [router.query.page])


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
                        <h2 className='text-4xl'>Auditing Events</h2>
                    </div>
                    <div className='m-4'>
                        <div className="overflow-x-auto">
                            <table className="table w-full">
                                {/* head*/}
                                <thead>
                                    <tr>
                                        <th>ID</th>
                                        <th>Verb</th>
                                        <th>Resource</th>
                                        <th>Namespaced Name</th>
                                        <th>Component / User-Agent</th>
                                        <th>Time</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {
                                        eventsListQuery.data?.completedRequestResponseAuditEvents.rows?.map((item, index) => {
                                            return (<tr key={item?.id}>
                                                <th>{item?.id}</th>
                                                <td>{(item?.verb as String).toUpperCase()}</td>
                                                <td>{`${resourceName(item?.apigroup, item?.apiversion!, item?.resource!)}`}</td>
                                                <td>{`${namespacedName(item?.namespace, item?.name!)}`}</td>
                                                <td>${item?.useragent}</td>
                                                <td>{moment(item?.stagetimestamp).format('YYYY-MM-DD HH:mm:ss Z')}</td>
                                            </tr>)
                                        })
                                    }
                                </tbody>
                            </table>
                        </div>
                    </div>
                    <div className='flex justify-center'>
                        <div className="btn-group">
                            <button className={`btn ${eventsListQuery.data?.completedRequestResponseAuditEvents.hasPreviousPage ? "" : "btn-disabled"}`} onClick={() => {
                                router.push(`/events/page/${page - 1}`)
                            }}>«</button>
                            <button className="btn">Page {`${page + 1} / ${eventsListQuery.data?.completedRequestResponseAuditEvents.totalPages}`}</button>
                            <button className={`btn ${eventsListQuery.data?.completedRequestResponseAuditEvents.hasNextPage ? "" : "btn-disabled"}`} onClick={() => {
                                router.push(`/events/page/${page + 1}`)
                            }}>»</button>
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
