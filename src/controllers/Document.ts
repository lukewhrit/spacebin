import { PhoneticKeyGenerator } from './KeyGenerator'
import { ConfigObject } from './Config'

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
          console.log(key)
        })
        .catch(reject)
    })
  }

  async newDocument () {

  }
}
