import 'reflect-metadata'
import Koa from 'koa'
import * as config from './controllers/config.controller'
import morgan from 'koa-morgan'
import cors from '@koa/cors'
import bodyParser from 'koa-bodyparser'
import ratelimit from 'koa-ratelimit'
import helmet from 'koa-helmet'
import { router } from './routes/document.route'

const app = new Koa()

// Setup app middleware
app
  .use(ratelimit({
    driver: 'memory',
    db: new Map(),
    duration: config.rateLimits.duration,
    max: config.rateLimits.requests
  }))
  .use(cors())
  .use(bodyParser())
  .use(morgan('tiny'))
  .use(router.middleware())
  .use(helmet({
    contentSecurityPolicy: config.useCSP || true,
    referrerPolicy: true
  }))
  .use(helmet.contentSecurityPolicy({
    directives: {
      defaultSrc: ['https', 'unsafe-eval', 'unsafe-inline'],
      objectSrc: ['none'],
      styleSrc: ['self', 'unsafe-inline'],
      frameAncestors: ['none']
    }
  }))

// Try to spawn server
try {
  app.listen(config.port, config.host)

  console.log(`Spacebin started on ${config.host}:${config.port}`)
} catch (err) {
  throw new Error(err)
}
