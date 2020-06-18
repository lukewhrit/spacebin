import { resolve } from 'path'
import { ConfigObject } from './controllers/config.controller'
import { Document } from './entities/document.entity'

export const config: ConfigObject = {
  routePrefix: '/api/v1/',
  dbOptions: {
    type: 'sqlite',
    database: resolve(__dirname, '..', 'data', 'db.sqlite'),
    synchronize: true,
    logging: false,
    entities: [
      Document
    ]
  }
}
