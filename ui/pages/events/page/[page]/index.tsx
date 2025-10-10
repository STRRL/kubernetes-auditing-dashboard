import { graphql } from '@/modules/gql'
import { useQuery } from '@tanstack/react-query'
import request from 'graphql-request'
import Head from 'next/head'
import { useEffect, useState } from 'react'
import { useRouter } from 'next/router'
import Link from 'next/link'
import { Sidebar } from '@/components/Sidebar'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Button } from '@/components/ui/button'

const moment = require('moment');

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
}

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

const buildLifecycleUrl = (apiGroup: string | undefined, apiVersion: string, resource: string, namespace: string | undefined, name: string) => {
    const kind = resource.charAt(0).toUpperCase() + resource.slice(1, -1);

    const params = new URLSearchParams();
    if (apiGroup && apiGroup !== '') {
        params.set('group', apiGroup);
    }
    params.set('version', apiVersion);
    params.set('kind', kind);
    if (namespace && namespace !== '') {
        params.set('namespace', namespace);
    }
    params.set('name', name);

    return `/lifecycle?${params.toString()}`;
}

export default function Events() {
    const router = useRouter()

    const [page, setPage] = useState(0)
    const [pageSize, setPageSize] = useState(12)

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
            <Sidebar>
                <div className="p-4">
                    <div className='m-4'>
                        <h2 className='text-4xl font-bold'>Recent Changes</h2>
                    </div>
                    <div className='m-4'>
                        <div className="overflow-x-auto">
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead>Verb</TableHead>
                                        <TableHead>Resource / Name</TableHead>
                                        <TableHead>Component / User-Agent</TableHead>
                                        <TableHead>Time</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {
                                        eventsListQuery.data?.completedRequestResponseAuditEvents.rows?.map((item, index) => {
                                            const userAgent = formatUserAgent(item?.useragent || '');
                                            const verbInfo = formatVerb(item?.verb || '');
                                            const resourceInfo = formatResource(item?.apigroup, item?.apiversion!, item?.resource!, item?.namespace, item?.name!);
                                            const timeInfo = formatTime(item?.stagetimestamp);
                                            const lifecycleUrl = buildLifecycleUrl(item?.apigroup, item?.apiversion!, item?.resource!, item?.namespace, item?.name!);
                                            return (<TableRow key={item?.id}>
                                                <TableCell>
                                                    <span className={`px-2 py-1 rounded-full text-xs font-semibold text-white ${verbInfo.color}`}>
                                                        {verbInfo.verb}
                                                    </span>
                                                </TableCell>
                                                <TableCell>
                                                    <Link href={lifecycleUrl} className="hover:opacity-80 transition-opacity">
                                                        <div className="flex flex-col cursor-pointer">
                                                            <span className="text-sm font-medium text-gray-700">{resourceInfo.resource}</span>
                                                            <span className="text-xs text-blue-600 hover:text-blue-800 underline">{resourceInfo.name}</span>
                                                        </div>
                                                    </Link>
                                                </TableCell>
                                                <TableCell>
                                                    {userAgent.isKnown ? (
                                                        <span className={`px-2 py-1 rounded-full text-xs font-semibold text-white ${userAgent.color}`}>
                                                            {userAgent.name}
                                                        </span>
                                                    ) : (
                                                        userAgent.name
                                                    )}
                                                </TableCell>
                                                <TableCell>
                                                    <div className="flex flex-col">
                                                        <span className="text-sm font-medium text-gray-700">{timeInfo.time}</span>
                                                        <span className="text-xs text-gray-500">{timeInfo.date}</span>
                                                        <span className="text-xs text-gray-400">{timeInfo.fromNow}</span>
                                                    </div>
                                                </TableCell>
                                            </TableRow>)
                                        })
                                    }
                                </TableBody>
                            </Table>
                        </div>
                    </div>
                    <div className='flex justify-center gap-2'>
                        <Button
                            variant="outline"
                            disabled={page === 0}
                            onClick={() => router.push(`/events/page/0`)}
                        >
                            ««
                        </Button>
                        <Button
                            variant="outline"
                            disabled={!eventsListQuery.data?.completedRequestResponseAuditEvents.hasPreviousPage}
                            onClick={() => router.push(`/events/page/${page - 1}`)}
                        >
                            «
                        </Button>
                        <Button variant="outline" disabled>
                            Page {`${page + 1} / ${eventsListQuery.data?.completedRequestResponseAuditEvents.totalPages}`}
                        </Button>
                        <Button
                            variant="outline"
                            disabled={!eventsListQuery.data?.completedRequestResponseAuditEvents.hasNextPage}
                            onClick={() => router.push(`/events/page/${page + 1}`)}
                        >
                            »
                        </Button>
                        <Button
                            variant="outline"
                            disabled={page >= (eventsListQuery.data?.completedRequestResponseAuditEvents.totalPages || 0) - 1}
                            onClick={() => router.push(`/events/page/${(eventsListQuery.data?.completedRequestResponseAuditEvents.totalPages || 1) - 1}`)}
                        >
                            »»
                        </Button>
                    </div>
                </div>
            </Sidebar>
        </>
    )
}
