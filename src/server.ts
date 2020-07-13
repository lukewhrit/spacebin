import 'reflect-metadata' // For TypeORM

import express from 'express'
import * as config from './controllers/config.controller'
import https from 'https'
import * as log from './logger'

// Middlewares
import responseTime from 'response-time'
import cors from 'cors'
import bodyParser from 'body-parser'
import rateLimit from 'express-rate-limit'
import helmet from 'helmet'

const app = express()

// Initialize middleware
app
  .use(cors())
  .use(responseTime())
  .use(bodyParser.json({ limit: config.maxDocumentLength * 1000 }))
  .use(log.express)
  .use(rateLimit({
    windowMs: config.rateLimits.duration,
    max: config.rateLimits.requests
  }))
  .use(helmet({
    contentSecurityPolicy: config.useCSP ? {
      directives: {
        defaultSrc: ["'none'"],
        objectSrc: ["'none'"],
        scriptSrc: ["'self'"],
        styleSrc: ["'self'"],
        frameAncestors: ["'none'"],
        baseUri: ["'none'"],
        formAction: ["'none'"]
      }
    } : false,
    referrerPolicy: true
  }))

// correctly register IPs when behind proxies
app.set('trust proxy', 1)

app.get('/', (req, res) => {
  res.send({ message: 'body' })
})

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
