/*
 * Copyright (C) 2020 The Spacebin Authors: notably Luke Whrit, Jack Dorland

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

import { ValidationError } from '@hapi/joi'
import { NextFunction, Response, Request } from 'express'
import { SpacebinError } from '../controllers/response.controller'
import { validateCreate, validateRead, validateRawRead, validateVerify } from './document.validator'

export function validate (validator: 'create' | 'verify' | 'read' | 'readRaw'): (
    req: Request,
    res: Response,
    next: NextFunction
  ) => void {
  const validate = (res: Response, next: NextFunction, error?: ValidationError): void => {
    if (error) {
      res.status(400).send(new SpacebinError(res, {
        message: error.details[0].message
      }))
    } else {
      next()
    }
  }

  return (req: Request, res: Response, next: NextFunction): void => {
    switch (validator) {
      case 'create':
        validate(res, next, validateCreate(req.body).error)
        break
      case 'read':
        validate(res, next, validateRead(req.body).error)
        break
      case 'readRaw':
        validate(res, next, validateRawRead(req.body).error)
        break
      case 'verify':
        validate(res, next, validateVerify(req.body).error)
        break
    }
  }
}
