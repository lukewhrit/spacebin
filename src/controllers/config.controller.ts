interface RateLimits {
  requests: number;
  duration: number;
}

export interface ConfigObject {
  host: string;
  port: number;

  keyLength: number;
  maxDocumentLength: number;
  staticMaxAge: number;

  useBrotli: boolean;
  useGzip: boolean;

  rateLimits: RateLimits;

  dbHost: string;
  dbOptions: object;
}

export class Config {
  options: ConfigObject

  constructor (options: ConfigObject) {
    this.options = options
  }
}
