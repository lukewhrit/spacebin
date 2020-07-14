/*
 * This file contains various values that rarely need to be changed and are used throughout Spirit.
 */

export const cspConfig = {
  directives: {
    defaultSrc: ["'none'"],
    objectSrc: ["'none'"],
    scriptSrc: ["'self'"],
    styleSrc: ["'self'"],
    frameAncestors: ["'none'"],
    baseUri: ["'none'"],
    formAction: ["'none'"]
  }
}

export const routePrefix = '/api/v1/'
