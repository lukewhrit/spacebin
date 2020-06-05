import { Config } from './controllers/config.controller'
import { resolve } from 'path'

export default new Config({
  host: '0.0.0.0',
  port: 7777,

  keyLength: 12,
  maxDocumentLength: 400_000,
  staticMaxAge: 86_400,

  useBrotli: true,
  useGzip: true,
  enableCSP: true,

  rateLimits: {
    requests: 500,
    duration: 60_000
  },

  dbHost: '',
  dbOptions: {
    client: 'sqlite3',
    connection: {
      filename: resolve(process.cwd(), 'src', '/config.ts')
    }
  }
})
