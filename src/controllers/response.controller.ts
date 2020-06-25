import { Context as KoaContext } from 'koa'
import { VerifyResponse } from '../structures/verifyResponse.struct'
import { DocumentResponse } from '../structures/documentResponse.struct'

export interface ResponseBuilderOptions {
  payload?: VerifyResponse | DocumentResponse;
  error?: Error;
}

export class ResponseBuilder {
  constructor (ctx: KoaContext, options: ResponseBuilderOptions) {
    const { payload, error } = options

    return {
      status: ctx.status,
      payload: payload || {},
      error: error || {}
    }
  }
}
