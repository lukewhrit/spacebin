import Router from '@koa/router'
import { createConnection } from 'typeorm'
import { Document } from '../entities/document.entity'
import { dbOptions } from '../controllers/database.controller'

export default async (router: Router): Promise<void> => {
  const connection = await createConnection(dbOptions)
  const documents = connection.getRepository(Document)

  router.post('create document', '/document', async (ctx) => {
    ctx.status = 501
  })

  router.get('get document', '/document', async (ctx) => {
    /*
      Body: {
        "key": "xqm5wXXE"
      }
     */

    const doc = await documents.findOne({
      where: { id: ctx.request.body.key }
    })

    ctx.body = doc
  })
}
