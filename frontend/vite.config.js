const proxy = require('koa-proxy');

/**
 * type {import('vite').UserConfig}
 */
export default {
  configureServer: ({ app }) => app.use(proxy({
    host: 'http://localhost:8080',
    match: /^\/api\//
  })),
  base: './'
};