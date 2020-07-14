import { resolve } from 'path'
import { ConfigObject } from './controllers/config.controller'
import { Document } from './entities/document.entity'

export const config: ConfigObject = {
  useCSP: true,
  dbOptions: {
    type: 'sqlite',
    database: resolve(__dirname, '..', 'data', 'db.sqlite'),
    synchronize: process.env.NODE_ENV === 'development',
    logging: false,
    entities: [
      Document
    ]
  }
}
