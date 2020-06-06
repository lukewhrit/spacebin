import { Config } from './controllers/config.controller'
import { resolve } from 'path'
import { Document } from './entities/document.entity'

export default new Config({
  host: '0.0.0.0',
  port: 7777,

  keyLength: 12,
  maxDocumentLength: 400_000,
  staticMaxAge: 86_400,

  useBrotli: true,
  useGzip: true,
  enableCSP: false,

  rateLimits: {
    requests: 500,
    duration: 60_000
  },

  dbOptions: {
    type: 'sqlite',
    database: resolve(__dirname, '..', 'data', 'db.sqlite'),
    synchronize: true,
    logging: false,
    entities: [
      Document
    ]
  }
})
