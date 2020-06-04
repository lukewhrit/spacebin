import Router from '@koa/router'
import { createConnection } from 'typeorm'
import { Document } from '../entities/document.entity'
import { dbOptions } from '../controllers/database.controller'
import { DocumentHandler } from '../controllers/document.controller'
import config from '../config'

const router = new Router({
  prefix: '/api/v1'
})

const main = async (): Promise<void> => {
  // get document repository
  const connection = await createConnection(dbOptions)
  const documents = connection.getRepository(Document)

  // create instance of document handler
  const documentHandler = new DocumentHandler(config.options)

  router.post('/document', async (ctx) => {
    const document = await documentHandler.newDocument(ctx.request.body.content, documents)

    ctx.body = document
  })

  router.get('/document', async (ctx) => {
    const doc = await documents.findOne({
      where: { id: ctx.request.body.key }
    })

    ctx.body = doc
  })
}

main()

export { router }
