import { ConnectionOptions } from 'typeorm'
import { config } from '../config'

interface RateLimits {
  requests: number;
  duration: number;
}

export interface ConfigObject {
  host?: string;
  port?: number;

  idLength?: number;
  maxDocumentLength?: number;

  useBrotli?: boolean;
  useGzip?: boolean;
  enableCSP?: boolean;

  rateLimits?: RateLimits;

  dbOptions: ConnectionOptions;
}

export const { // from https://wesbos.com/destructuring-default-values
  host = '0.0.0.0',
  port = 7777,

  idLength = 12,
  maxDocumentLength = 400_000,

  useBrotli = true,
  useGzip = true,
  enableCSP = false,

  rateLimits = {
    requests: 500,
    duration: 60_000
  },

  dbOptions // No defaults for dbOptions but we still export it.
} = config
