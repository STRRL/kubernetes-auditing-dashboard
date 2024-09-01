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

// 更新 formatUserAgent 函数
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

    // 如果不是已知组件，返回原始字符串
    return { name: userAgent, color: '', isKnown: false };
}

// 添加这个新函数来处理 verb
const formatVerb = (verb: string) => {
    const verbColors: { [key: string]: string } = {
        'get': 'bg-blue-500',
        'list': 'bg-cyan-500',
        'watch': 'bg-yellow-500',
        'create': 'bg-green-500',
        'update': 'bg-indigo-500',
        'patch': 'bg-pink-500',
        'delete': 'bg-red-500',
    };

    const color = verbColors[verb.toLowerCase()] || 'bg-gray-500';
    return { verb: verb.toUpperCase(), color };
}

const formatResource = (apiGroup: string | undefined, apiVersion: string, resource: string, namespace: string | undefined, name: string) => {
    const resourceString = `${apiGroup ? apiGroup + "/" : ""}${apiVersion} ${resource}`;
    const nameString = namespacedName(namespace, name);

    return {
        resource: resourceString,
        name: nameString,
    };
}

const formatTime = (timestamp: string) => {
    const date = moment(timestamp);
    return {
        date: date.format('YYYY-MM-DD'),
        time: date.format('HH:mm:ss'),
        fromNow: date.fromNow()
    };
}

export default function Events() {
    const router = useRouter()

    const [page, setPage] = useState(0)
    const [pageSize, setPageSize] = useState(12)  // 将这里的值从 15 改为 12

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
                        <h2 className='text-4xl'>Recent Changes</h2>
                    </div>
                    <div className='m-4'>
                        <div className="overflow-x-auto">
                            <table className="table w-full">
                                {/* head*/}
                                <thead>
                                    <tr>
                                        <th>Verb</th>
                                        <th>Resource / Name</th>
                                        <th>Component / User-Agent</th>
                                        <th>Time</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {
                                        eventsListQuery.data?.completedRequestResponseAuditEvents.rows?.map((item, index) => {
                                            const userAgent = formatUserAgent(item?.useragent || '');
                                            const verbInfo = formatVerb(item?.verb || '');
                                            const resourceInfo = formatResource(item?.apigroup, item?.apiversion!, item?.resource!, item?.namespace, item?.name!);
                                            const timeInfo = formatTime(item?.stagetimestamp);
                                            return (<tr key={item?.id}>
                                                <td>
                                                    <span className={`px-2 py-1 rounded-full text-xs font-semibold text-white ${verbInfo.color}`}>
                                                        {verbInfo.verb}
                                                    </span>
                                                </td>
                                                <td>
                                                    <div className="flex flex-col">
                                                        <span className="text-sm font-medium text-gray-700">{resourceInfo.resource}</span>
                                                        <span className="text-xs text-gray-500">{resourceInfo.name}</span>
                                                    </div>
                                                </td>
                                                <td>
                                                    {userAgent.isKnown ? (
                                                        <span className={`px-2 py-1 rounded-full text-xs font-semibold text-white ${userAgent.color}`}>
                                                            {userAgent.name}
                                                        </span>
                                                    ) : (
                                                        userAgent.name
                                                    )}
                                                </td>
                                                <td>
                                                    <div className="flex flex-col">
                                                        <span className="text-sm font-medium text-gray-700">{timeInfo.time}</span>
                                                        <span className="text-xs text-gray-500">{timeInfo.date}</span>
                                                        <span className="text-xs text-gray-400">{timeInfo.fromNow}</span>
                                                    </div>
                                                </td>
                                            </tr>)
                                        })
                                    }
                                </tbody>
                            </table>
                        </div>
                    </div>
                    <div className='flex justify-center'>
                        <div className="btn-group">
                            <button className={`btn ${page > 0 ? "" : "btn-disabled"}`} onClick={() => {
                                router.push(`/events/page/0`)
                            }}>««</button>
                            <button className={`btn ${eventsListQuery.data?.completedRequestResponseAuditEvents.hasPreviousPage ? "" : "btn-disabled"}`} onClick={() => {
                                router.push(`/events/page/${page - 1}`)
                            }}>«</button>
                            <button className="btn">Page {`${page + 1} / ${eventsListQuery.data?.completedRequestResponseAuditEvents.totalPages}`}</button>
                            <button className={`btn ${eventsListQuery.data?.completedRequestResponseAuditEvents.hasNextPage ? "" : "btn-disabled"}`} onClick={() => {
                                router.push(`/events/page/${page + 1}`)
                            }}>»</button>
                            <button className={`btn ${page < (eventsListQuery.data?.completedRequestResponseAuditEvents.totalPages || 0) - 1 ? "" : "btn-disabled"}`} onClick={() => {
                                router.push(`/events/page/${(eventsListQuery.data?.completedRequestResponseAuditEvents.totalPages || 1) - 1}`)
                            }}>»»</button>
                        </div>
                    </div>
                </div>

                <div className="drawer-side">
                    <label htmlFor="drawer-indicator" className="drawer-overlay"></label>
                    <ul className="menu p-4 w-80 bg-base-100 text-base-content">
                        <li><a href='/'>Home</a></li>
                        <li><a href='/events'>Recent Changes</a></li>
                        <li><a href='/lifecycle'>Resource Lifecycle(TBD)</a></li>
                    </ul>
                </div>
            </div>
        </>
    )
}
