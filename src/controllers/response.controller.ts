/*
 * Copyright (C) 2020 The Spacebin Authors

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.

 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
      error: error || ''
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
