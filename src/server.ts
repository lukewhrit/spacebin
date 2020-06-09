import 'reflect-metadata'
import Koa from 'koa'
import { host, port, enableCSP, rateLimits } from './controllers/config.controller'
import morgan from 'koa-morgan'
import cors from '@koa/cors'
import bodyParser from 'koa-bodyparser'
import ratelimit from 'koa-ratelimit'
import helmet from 'koa-helmet'
import { router } from './routes/document.route'

const app = new Koa()

// setup app middleware
app
  .use(ratelimit({
    driver: 'memory',
    db: new Map(),
    duration: rateLimits.duration,
    max: rateLimits.requests
  }))
  .use(cors())
  .use(bodyParser())
  .use(morgan('tiny'))
  .use(router.routes())
  .use(router.allowedMethods())
  .use(helmet({
    contentSecurityPolicy: enableCSP || false,
    referrerPolicy: true
  }))

// spawn server
try {
  app.listen(port, host)

  console.log(`Spacebin started on ${host}:${port}`)
} catch (err) {
  throw new Error(err)
}
