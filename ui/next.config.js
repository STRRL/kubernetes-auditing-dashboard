const API_URL = process.env.API_URL || 'http://localhost:23333/api'

/** @type {import('next').NextConfig} */
const nextConfig = {
  rewrites: async () => [
    {
      source: '/api/:path*',
      destination: `${API_URL}/:path*`,
    },
  ],
  reactStrictMode: true,
}

module.exports = nextConfig
