import Koa from 'koa'
import config from './config'
import morgan from 'koa-morgan'
import Router from '@koa/router'
import cors from '@koa/cors'
import bodyParser from 'koa-body'
import ratelimit from 'koa-ratelimit'
import DocumentRoute from './routes/document.route'
import helmet from 'koa-helmet'

const app = new Koa()
const ratelimitDB = new Map()
const router = new Router({
  prefix: '/api/v1'
})

// Setup app middleware
app
  .use(ratelimit({
    driver: 'memory',
    db: ratelimitDB,
    duration: config.options.rateLimits.duration,
    max: config.options.rateLimits.requests
  }))
  .use(cors())
  .use(bodyParser())
  .use(morgan('tiny'))
  .use(router.routes())
  .use(router.allowedMethods())
  .use(helmet())

// Register routes
DocumentRoute(router)

try {
  app.listen(config.options.port, config.options.host)

  console.log(`Spacebin started on ${config.options.host}:${config.options.port}`)
} catch (err) {
  throw new Error(err)
}
