// eslint-disable-next-line no-unused-vars
import Router from '@koa/router'

export default (router: Router) => {
  router.post('create document', '/create', async (ctx, next) => {
    ctx.status = 501
  })

  router.get('read document', '/read/:document', async (ctx, next) => {
    ctx.status = 501
  })
}
