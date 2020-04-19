import Koa from 'koa'
import serve from 'koa-static'
import config from './config'
import { resolve } from 'path'
import morgan from 'koa-morgan'
import Router from '@koa/router'
import cors from '@koa/cors'
import bodyParser from 'koa-body'
import HelloRoute from './routes/TestRoute'

const app = new Koa()
const router = new Router({
  prefix: '/api/v1'
})

// Serve static files
app.use(serve(resolve(process.cwd(), 'static'), {
  gzip: true,
  brotli: true,
  maxAge: config.options.staticMaxAge
}))

// Setup app middleware
app
  .use(cors())
  .use(bodyParser())
  .use(morgan('tiny'))
  .use(router.routes())
  .use(router.allowedMethods())

// Register routes
HelloRoute(router)

app.listen(config.options.port, config.options.host)
