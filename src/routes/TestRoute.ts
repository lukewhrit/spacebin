import Router from '@koa/router'

export default (router: Router) => {
  router.get('/test', (ctx, next) => {
    ctx.body = {
      message: 'Hello, world!'
    }
  })
}
