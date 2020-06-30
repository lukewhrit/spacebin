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
            status: Joi.number().max(599).min(100),
            payload: {
              id: Joi.string().required(),
              contentHash: Joi.string().hex().required(),
              extension: Joi.string().required()
            },
            error: Joi.object()
          }
        },
        500: {
          body: {
            status: Joi.number().max(599).min(100),
            payload: Joi.object().empty(),
            error: Joi.object()
          }
        }
      }
    }
  },
  verify: {
    validate: {
      params: {
        id: Joi.string().max(config.maxDocumentLength).required().insensitive()
      },
      output: {
        500: {
          body: {
            status: Joi.number().max(599).min(100),
            payload: Joi.object().empty(),
            error: Joi.object()
          }
        },
        200: {
          body: {
            status: Joi.number().max(599).min(100),
            payload: Joi.object().keys({
              exists: Joi.boolean()
            }),
            error: Joi.object()
          }
        },
        404: {
          body: {
            status: Joi.number().max(599).min(100),
            payload: Joi.object().keys({
              exists: Joi.boolean()
            }),
            error: Joi.object()
          }
        }
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
            status: Joi.number().max(599).min(100),
            payload: Joi.object().keys({
              id: Joi.string().required(),
              content: Joi.string().required(),
              dateCreated: Joi.date().required(),
              extension: Joi.string().required()
            }),
            error: Joi.object().empty()
          }
        },
        500: {
          body: {
            status: Joi.number().max(599).min(100),
            payload: Joi.object().empty(),
            error: Joi.object()
          }
        }
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
            status: Joi.number().max(599).min(100),
            payload: Joi.object().empty(),
            error: Joi.object()
          }
        }
      }
    }
  }
}
