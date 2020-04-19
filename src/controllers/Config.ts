interface RateLimits {
  requests: number
  every: number
}

interface ConfigObject {
  host: string
  port: number

  keyLength: number
  maxLength: number
  staticMaxAge: number
  recompressStaticAssets: boolean

  rateLimits: RateLimits
}

export class Config {
  options: ConfigObject

  constructor (options: ConfigObject) {
    this.options = options
  }
}
