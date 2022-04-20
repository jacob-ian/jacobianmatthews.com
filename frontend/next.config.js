/** @type {import('next').NextConfig} */
module.exports = {
  reactStrictMode: true,
  swcMinify: false,
  async rewrites() {
    const backendUrl = process.env.BACKEND_URL;

    if (!backendUrl) {
      throw new Error("Missing Backend URL for Reverse Proxying");
    }

    return [
      {
        source: "/api",
        destination: backendUrl,
      },
      {
        source: "/api/:slug*",
        destination: `${backendUrl}/:slug*`,
      },
    ];
  },
};
