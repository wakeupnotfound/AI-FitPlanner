/**
 * Vite Plugin for Content Security Policy
 * Adds CSP headers to responses in development mode
 */

export function cspPlugin() {
  return {
    name: 'vite-plugin-csp',
    configureServer(server) {
      server.middlewares.use((req, res, next) => {
        // Set CSP header for development
        const csp = [
          "default-src 'self'",
          "script-src 'self' 'unsafe-inline' 'unsafe-eval'", // unsafe-eval needed for Vue dev tools
          "style-src 'self' 'unsafe-inline'",
          "img-src 'self' data: https: blob:",
          "font-src 'self' data:",
          "connect-src 'self' http://localhost:9999 ws://localhost:* https://api.openai.com https://aip.baidubce.com https://dashscope.aliyuncs.com",
          "frame-src 'none'",
          "object-src 'none'",
          "base-uri 'self'",
          "form-action 'self'"
        ].join('; ')

        res.setHeader('Content-Security-Policy', csp)
        
        // Additional security headers
        res.setHeader('X-Content-Type-Options', 'nosniff')
        res.setHeader('X-Frame-Options', 'DENY')
        res.setHeader('X-XSS-Protection', '1; mode=block')
        res.setHeader('Referrer-Policy', 'strict-origin-when-cross-origin')
        res.setHeader('Permissions-Policy', 'geolocation=(), microphone=(), camera=()')
        
        next()
      })
    }
  }
}
