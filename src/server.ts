import Koa from 'koa'
import config from './config'
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
    duration: config.options.rateLimits.duration,
    max: config.options.rateLimits.requests
  }))
  .use(cors())
  .use(bodyParser())
  .use(morgan('tiny'))
  .use(router.routes())
  .use(router.allowedMethods())
  .use(helmet({
    contentSecurityPolicy: config.options.enableCSP || false,
    referrerPolicy: true
  }))

// spawn server
try {
  app.listen(config.options.port, config.options.host)

  console.log(`Spacebin started on ${config.options.host}:${config.options.port}`)
} catch (err) {
  throw new Error(err)
}
