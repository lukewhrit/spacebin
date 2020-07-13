import 'reflect-metadata' // For TypeORM

import express from 'express'
import * as config from './controllers/config.controller'
import https from 'https'
import * as log from './logger'
import responseTime from 'response-time'

const app = express()

app
  .use(responseTime())

// Use an HTTPs server if SSL is enabled, otherwise use `app`
const server = config.useSSL
  ? https.createServer({
    key: config.sslOptions?.key,
    cert: config.sslOptions?.cert
  }, app)
  : app

// Spawn server
try {
  server.listen(config.port, config.host)

  log.success(`Spacebin started on ${config.host}:${config.port}`)
} catch (err) {
  throw new Error(err)
}
