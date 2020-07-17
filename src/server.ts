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

import 'reflect-metadata' // For TypeORM

import express from 'express'
import https from 'https'
import path from 'path'
import { loadRoutes } from './controllers/util.controller'
import * as log from './logger'
import * as config from './controllers/config.controller'
import { cspConfig } from './consts'

import cors from 'cors'
import bodyParser from 'body-parser'
import rateLimit from 'express-rate-limit'
import helmet from 'helmet'

const app = express()

// Initialize middleware
app
  .use(cors())
  .use(bodyParser.json({ limit: config.maxDocumentLength * 1000 }))
  .use(bodyParser.urlencoded({ extended: false }))
  .use(log.express())
  .use(rateLimit({
    windowMs: config.rateLimits.duration,
    max: config.rateLimits.requests
  }))
  .use(helmet({
    contentSecurityPolicy: config.useCSP ? cspConfig : false,
    referrerPolicy: true
  }))

// correctly register IPs when behind proxies
app.set('trust proxy', 1)

loadRoutes(path.join(__dirname, 'routes'), app)

// Use an HTTPs server if SSL is enabled, otherwise use `app`
const server = config.useSSL
  ? https.createServer({
    key: config.sslOptions?.key,
    cert: config.sslOptions?.cert
  }, app)
  : app

if (!module.parent) {
  // Spawn server
  try {
    server.listen(config.port, config.host)

    log.success(`Spacebin started on ${config.host}:${config.port}`)
  } catch (err) {
    throw new Error(err)
  }
}

export { app }
