import Router from '@koa/router'
import { createConnection } from 'typeorm'
import { Document } from '../entities/document.entity'
import { DocumentHandler } from '../controllers/document.controller'
import * as config from '../controllers/config.controller'
import constants from '../const'

const router = new Router({
  prefix: constants.prefix
})

// needs to be a function for async/await
const main = async (): Promise<void> => {
  // setup document handler
  const connection = await createConnection(config.dbOptions)
  const documents = connection.getRepository(Document)
  const handler = new DocumentHandler(config, documents)

  router.post('/document', async (ctx) => {
    try {
      // create new document with contents of `ctx.request.body.content` in repository documents
      const doc = await handler.newDocument(ctx.request.body.content)

      ctx.body = doc
    } catch (err) {
      ctx.body = { err }
    }
  })

  router.get('/document/:id', async (ctx) => {
    try {
      const doc = await handler.getDocument(ctx.params.id)

      ctx.body = doc
    } catch (err) {
      ctx.body = { err }
    }
  })
}

main()

export { router }
