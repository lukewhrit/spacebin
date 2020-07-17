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

import chalk from 'chalk'
import { Request, Response } from 'express'

export type Keywords = 'error' | 'success' | 'info' | 'verbose' | 'debug'
export interface KeywordObject { background: string; text: string }

export const colorMap: Record<Keywords, KeywordObject> = {
  error: { background: 'indianred', text: 'black' },
  success: { background: 'mediumspringgreen', text: 'black' },
  verbose: { background: 'darkslategray', text: 'white' },
  debug: { background: 'yellow', text: 'black' },
  info: { background: 'cornflowerblue', text: 'cornsilk' }
}

const buildString = (message: string, keyword: Keywords): void => {
  const bgColor = chalk.bold.bgKeyword(colorMap[keyword].background).keyword(colorMap[keyword].text)

  const keywordString = bgColor(` ${keyword.toUpperCase()} `)
  const string = `🔭 ${keywordString} ${message}`

  console.log(string)
}

export const error = (message: string): void => buildString(message, 'error')
export const success = (message: string): void => buildString(message, 'success')
export const verbose = (message: string): void => buildString(message, 'verbose')
export const info = (message: string): void => buildString(message, 'info')
export const debug = (message: string): void => buildString(message, 'debug')

export function express (): (req: Request, res: Response, next: Function) => void {
  return (req: Request, res: Response, next: Function): void => {
    const { method, originalUrl } = req
    const { statusCode } = res

    info(`${method} ${originalUrl} ${statusCode} TODO (Content-Length) - TODO ms`)

    next()
  }
}
