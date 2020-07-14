import { VerifyResponse } from '../structures/verifyResponse.struct'
import { DocumentResponse } from '../structures/documentResponse.struct'
import { Response } from 'express'

export interface ResponseBuilderOptions {
  payload?: VerifyResponse | DocumentResponse;
  error?: string;
}

export interface SpacebinErrorOptions {
  status?: number;
  message: string;
}

export class ResponseBuilder {
  constructor (res: Response, options: ResponseBuilderOptions) {
    const { payload, error } = options

    return {
      status: res.statusCode,
      payload: payload || {},
      error: error || {}
    }
  }
}

export class SpacebinError extends ResponseBuilder {
  constructor (res: Response, options: SpacebinErrorOptions) {
    res.status(options.status || 500)

    super(res, {
      error: options.message
    })
  }
}
