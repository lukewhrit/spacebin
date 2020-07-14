import * as Joi from '@hapi/joi'
import * as config from '../controllers/config.controller'

export function validateCreate (data: Record<string, unknown>): Joi.ValidationResult {
  const schema = Joi.object({
    content: Joi.string().max(config.maxDocumentLength).required().insensitive(),
    extension: Joi.string().lowercase().optional().default('txt').insensitive()
  })

  return schema.validate(data)
}

export function validateVerify (data: Record<string, unknown>): Joi.ValidationResult {
  const schema = Joi.object({
    id: Joi.string().max(config.maxDocumentLength).required().insensitive()
  })

  return schema.validate(data)
}

export function validateRead (data: Record<string, unknown>): Joi.ValidationResult {
  const schema = Joi.object({
    id: Joi.string().length(config.idLength).insensitive().required()
  })

  return schema.validate(data)
}

export function validateRawRead (data: Record<string, unknown>): Joi.ValidationResult {
  const schema = Joi.object({
    id: Joi.string().length(config.idLength).required()
  })

  return schema.validate(data)
}
