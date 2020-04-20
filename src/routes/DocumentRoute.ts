import Router from '@koa/router'

export default (router: Router) => {
  router.post('/create', (ctx, next) => {
    ctx.body = ctx.request.body
  })

  router.get('read/:document', async (ctx, next) => {
    ctx.body = {
      message: 'Hello, world!'
    }
  })
}
