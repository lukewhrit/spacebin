import { ConnectionOptions } from 'typeorm'
import path from 'path'
import { Document } from '../entities/docuent.entity'

// @todo Handle this in config
export const dbOptions: ConnectionOptions = {
  type: 'sqlite',
  database: path.resolve(__dirname, '..', '..', 'data', 'db.sqlite'),
  entities: [
    Document
  ],
  synchronize: true,
  logging: false
}
