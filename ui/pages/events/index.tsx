'use client';
import { useRouter } from 'next/router';
import { useEffect } from 'react';

export default function Events() {
    const router = useRouter();
    useEffect(() => {
        router.push('/events/page/0');
    })
    return (<></>)
}
