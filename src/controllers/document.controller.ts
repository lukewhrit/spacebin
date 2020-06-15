import randomstring from 'randomstring'
import { ConfigObject } from './config.controller'
import { Document } from '../entities/document.entity'
import { Repository } from 'typeorm'

export class DocumentHandler {
  private options: ConfigObject
  private repository: Repository<Document>

  constructor (options: ConfigObject, docsRepo: Repository<Document>) {
    this.options = options
    this.repository = docsRepo
  }

  private createID (): string {
    return randomstring.generate(this.options.idLength || 12)
  }

  private chooseID (): Promise<string> {
    let id = this.createID()

    return new Promise((resolve) => {
      const doc = this.getDocument(id)

      doc.then(doc => {
        if (!doc) { // If ID is not found in DB..
          resolve(id)
        } else { // Otherwise re-call function
          id = this.createID()
          this.chooseID()
        }
      })
    })
  }

  async newDocument (content: string): Promise<Document> {
    const id = await this.chooseID()

    const doc = this.repository.create({
      id,
      content
    })

    this.repository.save(doc)

    return { ...doc }
  }

  async getDocument (id: string): Promise<Document | undefined> {
    const doc = await this.repository.findOne({
      where: { id }
    })

    return doc
  }

  async getRawDocument (id: string): Promise<string | undefined> {
    const documentObject = await this.getDocument(id)

    return documentObject?.content
  }
}
