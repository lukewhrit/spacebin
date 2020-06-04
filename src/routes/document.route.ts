import Router from '@koa/router'
import { createConnection } from 'typeorm'
import { Document } from '../entities/document.entity'
import { dbOptions } from '../controllers/database.controller'
import { DocumentHandler } from '../controllers/document.controller'
import config from '../config'

export default async (router: Router): Promise<void> => {
  const connection = await createConnection(dbOptions)
  const documents = connection.getRepository(Document)
  const documentHandler = new DocumentHandler(config.options)

  router.post('create document', '/document', async (ctx) => {
    const document = await documentHandler.newDocument(ctx.request.body.content, documents)

    ctx.body = document
  })

  router.get('get document', '/document', async (ctx) => {
    const doc = await documents.findOne({
      where: { id: ctx.request.body.key }
    })

    ctx.body = doc
  })
}
