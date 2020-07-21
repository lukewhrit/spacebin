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

import fs from 'fs'
import path from 'path'
import { Application } from 'express'
import { routePrefix } from '../consts'

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
      // don't load sourcemaps
      if (file.endsWith('.map')) { return }

      const filePath = path.join(routesDir, file)
      const loadedRoute = await import(filePath)

      // load route files into express
      app.use(`${routePrefix}${loadedRoute.prefix}`, loadedRoute.default)
    }
  })
}
