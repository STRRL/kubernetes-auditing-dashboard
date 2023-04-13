import { Inter } from 'next/font/google'
import Head from 'next/head'

const inter = Inter({ subsets: ['latin'] })

export default function Home() {
  return (
    <>
      <Head>
        <title>Kubernetes Auditing Dashboard</title>
      </Head>
      <div className="drawer drawer-mobile">
        <input id="drawer-indicator" type="checkbox" className="drawer-toggle" />
        <div className="drawer-content flex flex-col p-4">

          <div className='flex'>
            <div className="flex-1 stats shadow mx-4">
              <div className="stat">
                <div className="stat-title">Total Events</div>
                <div className="stat-value">89,400</div>
                <div className="stat-desc"></div>
              </div>
            </div>
            <div className="flex-1 stats shadow mx-4">
              <div className="stat">
                <div className="stat-title">Total Events</div>
                <div className="stat-value">89,400</div>
                <div className="stat-desc"></div>
              </div>
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
