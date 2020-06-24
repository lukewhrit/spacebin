import { Joi, Config } from 'koa-joi-router'
import * as config from '../controllers/config.controller'

interface Validators {
  create: Config;
  verify: Config;
  read: Config;
  readRaw: Config;
}

export const validators: Validators = {
  create: {
    validate: {
      body: {
        content: Joi.string().max(config.maxDocumentLength).required().insensitive(),
        extension: Joi.string().lowercase().optional().default('txt').insensitive()
      },
      type: 'json',
      output: {
        201: {
          body: {
            id: Joi.string().length(config.idLength),
            contentHash: Joi.string().hex()
          }
        },
        500: { body: { err: Joi.string() } }
      }
    }
  },
  verify: {
    validate: {
      params: {
        id: Joi.string().max(config.maxDocumentLength).required().insensitive()
      },
      type: 'json',
      output: {
        500: { body: { error: Joi.string() } },
        204: { body: { exists: Joi.boolean() } },
        404: { body: { exists: Joi.boolean() } }
      }
    }
  },
  read: {
    validate: {
      params: {
        id: Joi.string().length(config.idLength).insensitive().required()
      },
      type: 'json',
      output: {
        200: {
          body: {
            id: Joi.string().required(),
            content: Joi.string().required(),
            dateCreated: Joi.date().required(),
            extension: Joi.string().required()
          }
        },
        500: { body: { err: Joi.string() } }
      }
    }
  },
  readRaw: {
    validate: {
      params: {
        id: Joi.string().length(config.idLength).required()
      },
      output: {
        200: { body: Joi.string() },
        500: {
          body: {
            err: Joi.string()
          }
        }
      }
    }
  }
}
