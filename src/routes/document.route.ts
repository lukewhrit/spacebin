import Router from '@koa/router'

export default (router: Router): void => {
  router.post('create document', '/document', async (ctx) => {
    ctx.status = 501
  })

  router.get('get document', '/document', async (ctx) => {
    ctx.status = 501
  })
}
