import Router from '@koa/router'
import { createConnection } from 'typeorm'
import { Document } from '../entities/document.entity'
import { DocumentHandler } from '../controllers/document.controller'
import config from '../config'

const router = new Router({
  prefix: '/api/v1'
})

// needs to be a function for async/await
const main = async (): Promise<void> => {
  // get document repository
  const connection = await createConnection(config.options.dbOptions)
  const documents = connection.getRepository(Document)

  // create instance of document handler
  const documentHandler = new DocumentHandler(config.options)

  router.post('/document', async (ctx) => {
    try {
      // create new document with contents of request body content in repository documents
      const document = await documentHandler.newDocument(ctx.request.body.content, documents)

      ctx.body = document
    } catch (err) {
      ctx.body = { err }
    }
  })

  router.get('/document/:id', async (ctx) => {
    try {
      // find document where id = id in request body
      const doc = await documents.findOne({
        where: { id: ctx.params.id }
      })

      ctx.body = doc
    } catch (err) {
      ctx.body = { err }
    }
  })
}

main()

export { router }
