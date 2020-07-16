/*
 * Copyright (C) 2020 The Spacebin Authors: notably Luke Whrit, Jack Dorland

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.

 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import { ConnectionOptions } from 'typeorm'
import { config } from '../config'

interface RateLimits {
  requests: number;
  duration: number;
}

interface SSLOptions {
  cert: string;
  key: string;
}

export interface ConfigObject {
  host?: string;
  port?: number;

  idLength?: number;
  maxDocumentLength?: number;

  useBrotli?: boolean;
  useGzip?: boolean;
  useCSP?: boolean;
  useSSL?: boolean;

  rateLimits?: RateLimits;
  sslOptions?: SSLOptions;

  dbOptions: ConnectionOptions;
}

export const { // https://wesbos.com/destructuring-default-values
  host = '0.0.0.0',
  port = 7777,

  idLength = 12,
  maxDocumentLength = 400_000,

  useBrotli = true,
  useGzip = true,
  useCSP = false,
  useSSL = false,

  rateLimits = {
    requests: 500,
    duration: 60_000
  },

  sslOptions,

  dbOptions
} = config
