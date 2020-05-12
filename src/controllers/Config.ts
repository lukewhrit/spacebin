interface RateLimits {
  requests: number;
  every: number;
}

export interface ConfigObject {
  host: string;
  port: number;

  keyLength: number;
  maxLength: number;
  staticMaxAge: number;
  recompressStaticAssets: boolean;

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
