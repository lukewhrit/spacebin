import randomstring from 'randomstring'
import { ConfigObject } from './config.controller'
import { Document } from '../entities/document.entity'
import { createConnection, Connection, Repository } from 'typeorm'

type Extensions = 'py' | 'js' | 'jsx' | 'go' | 'ts' | 'tsx' | 'kt' | 'java' | 'cpp' | 'sql' |
  'cs' | 'swift' | 'xml' | 'dart' | 'r' | 'rb' | 'c' | 'h' | 'scala' | 'hs' |
  'sh' | 'ps1' | 'php' | 'asm' | 'jl' | 'm' | 'pl' | 'cr' | 'json' | 'yaml' | 'toml' | 'txt'
export class DocumentHandler {
  private config: ConfigObject
  private repository: Repository<Document>

  constructor (config: ConfigObject) {
    this.config = config

    DocumentHandler.connect(this.config)
      .then(con => {
        this.repository = con.getRepository(Document)
      })
  }

  private createID (): string {
    return randomstring.generate(this.config.idLength || 12)
  }

  private static connect (config: ConfigObject): Promise<Connection> {
    const connection = createConnection(config.dbOptions)

    return connection
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

  async newDocument (content: string, extension: Extensions): Promise<Document> {
    const id = await this.chooseID()

    const doc = this.repository.create({
      id,
      content,
      extension
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
