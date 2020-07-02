import { Context as KoaContext } from 'koa'
import { VerifyResponse } from '../structures/verifyResponse.struct'
import { DocumentResponse } from '../structures/documentResponse.struct'

export interface ResponseBuilderOptions {
  payload?: VerifyResponse | DocumentResponse;
  error?: string;
}

export interface SpacebinErrorOptions {
  status?: number;
  message: string;
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

export class SpacebinError extends ResponseBuilder {
  constructor (ctx: KoaContext, options: SpacebinErrorOptions) {
    ctx.status = options.status || 500

    super(ctx, {
      error: options.message
    })
  }
}
