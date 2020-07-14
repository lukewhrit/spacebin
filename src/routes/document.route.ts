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
