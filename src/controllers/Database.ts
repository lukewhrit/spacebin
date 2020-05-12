import { connect, connection, Connection } from 'mongoose'
import { GlueDocument, GlueDocumentModel } from '../models/Document.model'
import { ConfigObject } from './Config'
import config from '../config'

declare interface ModelInterfaces {
  GlueDocument: GlueDocumentModel;
}

export class Database {
  private static instance: Database

  private _db: Connection
  private _models: ModelInterfaces

  private constructor (config: ConfigObject) {
    // @todo Get directly from config
    connect(config.dbHost, config.dbOptions)

    this._db = connection
    this._db.on('error', this.error)

    this._models = {
      GlueDocument: new GlueDocument().model
    }
  }

  public static get Models (): ModelInterfaces {
    if (!Database.instance) {
      Database.instance = new Database(config.options)
    }

    return Database.instance._models
  }

  private error (error: Error): void {
    throw new Error(`Mongoose has errored: ${error}`)
  }
}
