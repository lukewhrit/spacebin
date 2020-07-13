import fs from 'fs'
import path from 'path'
import { Application } from 'express'
import * as config from './config.controller'

export function sanitize (input: string): string {
  // eslint-disable-next-line no-control-regex
  const ansiSequence = /(?:\x1B[@-_]|[\x80-\x9F])[0-?]*[ -/]*[@-~]/gi

  return input.replace(ansiSequence, '')
}

export function loadRoutes (routesDir: string, app: Application): void {
  fs.readdir(routesDir, async (err, files) => {
    if (err) throw err

    // loop over all files in routesDir
    for (const file of files) {
      // only work with route files
      if (!file.endsWith('route.ts')) { return }

      const filePath = path.join(routesDir, file)
      const loadedRoute = await import(filePath)

      // load route files into express
      app.use(`${config.routePrefix}${loadedRoute.prefix}`, loadedRoute.default)
    }
  })
}
