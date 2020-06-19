import joiRouter from 'koa-joi-router'
import { DocumentHandler } from '../controllers/document.controller'
import * as config from '../controllers/config.controller'
import crypto from 'crypto'
import { validators } from '../validators/document.validator'

const router = joiRouter()
const handler = new DocumentHandler(config)

router.prefix(`${config.routePrefix}documents`)

router.post('/', validators.create, async (ctx) => {
  try {
    const { id, content } = await handler.newDocument(ctx.request.body.content, ctx.request.body.extension)

    ctx.body = {
      id,
      contentHash: crypto.createHash('sha256').update(content).digest('hex')
    }

    ctx.status = 201
  } catch (err) {
    ctx.status = 500
    ctx.body = { err }
  }
})

router.post('/verify', validators.verify, async (ctx) => {
  try {
    const doc = await handler.getDocument(ctx.request.body.id)

    if (doc) {
      ctx.status = 200
    } else {
      ctx.status = 404
    }
  } catch (err) {
    ctx.status = 500
    ctx.body = { err }
  }
})

router.get('/:id', validators.read, async (ctx) => {
  try {
    const doc = await handler.getDocument(ctx.params.id)

    if (doc) {
      ctx.status = 200
      ctx.body = doc
    } else {
      ctx.status = 404
    }
  } catch (err) {
    ctx.status = 500
    ctx.body = { err }
  }
})

router.get('/:id/raw', validators.readRaw, async (ctx) => {
  try {
    const doc = await handler.getRawDocument(ctx.params.id)

    if (doc) {
      ctx.status = 200
      ctx.body = doc
    } else {
      ctx.status = 404
    }
  } catch (err) {
    ctx.status = 500
    ctx.body = { err }
  }
})

export { router }
