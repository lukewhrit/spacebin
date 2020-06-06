import { ConnectionOptions } from 'typeorm'

interface RateLimits {
  requests: number;
  duration: number;
}

export interface ConfigObject {
  host: string;
  port: number;

  idLength: number;
  maxDocumentLength: number;
  staticMaxAge: number;

  useBrotli: boolean;
  useGzip: boolean;
  enableCSP?: boolean;

  rateLimits: RateLimits;

  dbOptions: ConnectionOptions;
}

export class Config {
  options: ConfigObject

  constructor (options: ConfigObject) {
    this.options = options
  }
}
