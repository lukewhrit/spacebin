import Router from '@koa/router'

export default (router: Router): void => {
  router.post('create document', '/create', async (ctx) => {
    ctx.status = 501
  })

  router.get('read document', '/read/:document', async (ctx) => {
    ctx.status = 501
  })
}
