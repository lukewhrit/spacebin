import joiRouter from 'koa-joi-router'
import { createConnection } from 'typeorm'
import { Document } from '../entities/document.entity'
import { DocumentHandler } from '../controllers/document.controller'
import * as config from '../controllers/config.controller'
import crypto from 'crypto'

const router = joiRouter()
const Joi = joiRouter.Joi

router.prefix(config.routePrefix + 'document')

// This needs to be a function for async/await
const main = async (): Promise<void> => {
  // Setup document handler
  const connection = await createConnection(config.dbOptions)
  const documents = connection.getRepository(Document)
  const handler = new DocumentHandler(config, documents)

  router.post('/', {
    validate: {
      body: {
        content: Joi.string().max(config.maxDocumentLength).required(),
        extension: Joi.string().lowercase().optional().default('txt')
      },
      type: 'json',
      output: {
        201: {
          body: {
            id: Joi.string().length(config.idLength),
            contentHash: Joi.string().hex()
          }
        },
        500: {
          body: {
            err: Joi.string()
          }
        }
      }
    }
  }, async (ctx) => {
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

  router.post('/verify', {
    validate: {
      body: {
        id: Joi.string().max(config.maxDocumentLength).required()
      },
      type: 'json',
      output: { // Couldn't figure out how to properly include 200 and 404, so they're just left out.
        500: {
          body: {
            error: Joi.string()
          }
        }
      }
    }
  }, async (ctx) => {
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

  router.get('/:id', async (ctx) => {
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

  router.get('/:id/raw', async (ctx) => {
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
}

main()

export { router }
