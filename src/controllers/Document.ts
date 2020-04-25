import { PhoneticKeyGenerator } from './KeyGenerator'
// eslint-disable-next-line no-unused-vars
import { ConfigObject } from './Config'
import { Document } from '../models/Document.model'

export class DocumentHandler {
  options: ConfigObject

  constructor (options: ConfigObject) {
    this.options = options
  }

  private async createKey () {
    const keyGenerator = new PhoneticKeyGenerator()

    return keyGenerator.createKey(this.options.keyLength)
  }

  async chooseKey () {
    return new Promise((resolve, reject) => {
      this.createKey()
        .then((key) => {
          // @todo reject if key already exists

          resolve(key)
        })
        .catch(reject)
    })
  }

  async newDocument (content: string) {
    const doc = new Document()

    doc.set({ content, key: this.chooseKey() })
      .save()
  }
}
