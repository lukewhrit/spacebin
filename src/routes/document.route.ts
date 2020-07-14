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

import express from 'express'
import { ResponseBuilder as Response, SpacebinError } from '../controllers/response.controller'
import { DocumentHandler } from '../controllers/document.controller'
import * as config from '../controllers/config.controller'
import { createHash } from 'crypto'
import { validateCreate } from '../validators/document.validator'

const router = express.Router()
const handler = new DocumentHandler(config)

router.post('/', async (req, res) => {
  const { error } = validateCreate(req.body)
  if (error) {
    return res.status(400).send(new SpacebinError(res, {
      message: error.details[0].message
    }))
  }

  try {
    const { id, content, extension } = await handler.newDocument(
      req.body.content,
      req.body.extension
    )

    res.status(201).send(new Response(res, {
      payload: {
        id,
        contentHash: createHash('sha256').update(content).digest('hex'),
        extension
      }
    }))
  } catch (err) {
    res.send(new SpacebinError(res, {
      message: err
    }))
  }
})

export const prefix = 'document'
export default router
