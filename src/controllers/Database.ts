import Knex from 'knex'
import Bookshelf from 'bookshelf'
// eslint-disable-next-line no-unused-vars
import { ConfigObject } from './Config'
import Config from '../config'

export default class Database {
  private static _instance : Database = new Database(Config.options)
  protected _knex:any = null
  protected _bookshelf: any = null

  options: ConfigObject

  constructor (options: ConfigObject) {
    this.options = options

    if (Database._instance) {
      throw new Error('Error: Instantiation failed: Use Database.getInstance() instead of new.')
    }

    this._knex = Knex(this.options.dbOptions)
    this._bookshelf = Bookshelf(this._knex)
    Database._instance = this
  }

  public static getInstance (): Database {
    return Database._instance
  }

  public getKnex (): any {
    return this._knex
  }

  public getBookshelf (): Bookshelf {
    return this._bookshelf
  }
}
