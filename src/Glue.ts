import Koa from 'koa'
import serve from 'koa-static'
import config from './config'
import { resolve } from 'path'

const app = new Koa()

console.log(process.cwd())

app.use(serve(resolve(process.cwd(), 'static'), {
  gzip: true,
  brotli: true,
  maxAge: config.options.staticMaxAge
}))

app.listen(config.options.port, config.options.host)
