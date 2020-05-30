import Koa from 'koa'
import serve from 'koa-static'
import config from './config'
import { resolve } from 'path'
import morgan from 'koa-morgan'
import Router from '@koa/router'
import cors from '@koa/cors'
import bodyParser from 'koa-body'
import DocumentRoute from './routes/document.route'

const app = new Koa()
const router = new Router({
  prefix: '/api/v1'
})

// Setup app middleware
app
  .use(cors())
  .use(bodyParser())
  .use(morgan('tiny'))
  .use(router.routes())
  .use(router.allowedMethods())
  .use(serve(resolve(process.cwd(), 'static'), {
    gzip: config.options.useGzip,
    brotli: config.options.useBrotli,
    maxAge: config.options.staticMaxAge
  }))

// Register routes
DocumentRoute(router)

try {
  app.listen(config.options.port, config.options.host)

  console.log(`Glue started on ${config.options.host}:${config.options.port}`)
} catch (err) {
  throw new Error(err)
}
