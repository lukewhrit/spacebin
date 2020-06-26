import joiRouter from 'koa-joi-router'
import { DocumentHandler } from '../controllers/document.controller'
import * as config from '../controllers/config.controller'
import crypto from 'crypto'
import { validators } from '../validators/document.validator'
import { ResponseBuilder as Response } from '../controllers/response.controller'

const router = joiRouter()
const handler = new DocumentHandler(config)

router.prefix(`${config.routePrefix}document`)

router.post('/', validators.create, async (ctx) => {
  try {
    const { id, content, extension } = await handler.newDocument(ctx.request.body.content, ctx.request.body.extension)

    ctx.body = new Response(ctx, {
      payload: {
        id,
        contentHash: crypto.createHash('sha256').update(content).digest('hex'),
        extension
      }
    })

    ctx.status = 201
  } catch (error) {
    ctx.status = 500
    ctx.body = new Response(ctx, { error })
  }
})

router.get('/:id/verify', validators.verify, async (ctx) => {
  try {
    const doc = await handler.getDocument(ctx.params.id)

    ctx.status = doc ? 200 : 404

    ctx.body = new Response(ctx, {
      payload: {
        exists: !!doc
      }
    })
  } catch (error) {
    ctx.status = 500
    ctx.body = new Response(ctx, { error })
  }
})

router.get('/:id', validators.read, async (ctx) => {
  try {
    const doc = await handler.getDocument(ctx.params.id)

    if (doc) {
      ctx.status = 200
      ctx.body = new Response(ctx, {
        payload: {
          ...doc
        }
      })
    } else {
      ctx.status = 404
    }
  } catch (error) {
    ctx.status = 500
    ctx.body = new Response(ctx, { error })
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
  } catch (error) {
    ctx.status = 500
    ctx.body = new Response(ctx, { error })
  }
})

export { router }
