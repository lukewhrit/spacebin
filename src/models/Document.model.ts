import Database from '../controllers/Database'

const db = Database.getInstance()
const bookshelf = db.getBookshelf()

export class Document extends bookshelf.Model<Document> {
  get tableName () { return 'documents' }

  get hasTimestamps () { return true }

  // Schema values
  public get content (): string { return this.get('content') }
  public set content (value: string) { this.set({ content: value }) }

  public get key (): string { return this.get('key') }
  public set key (value: string) { this.set({ key: value }) }
}
