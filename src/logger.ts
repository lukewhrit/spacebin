import chalk from 'chalk'
import { Request, Response } from 'express'

export type Keywords = 'error' | 'success' | 'info' | 'verbose' | 'debug'
export interface KeywordObject { background: string; text: string }

export const colorMap: Record<Keywords, KeywordObject> = {
  error: {
    background: 'indianred',
    text: 'black'
  },
  success: {
    background: 'mediumspringgreen',
    text: 'black'
  },
  verbose: {
    background: 'darkslategray',
    text: 'white'
  },
  debug: {
    background: 'yellow',
    text: 'black'
  },
  info: {
    background: 'cornflowerblue',
    text: 'cornsilk'
  }
}

const buildString = (message: string, keyword: Keywords): void => {
  const bgText = colorMap[keyword].text
  const bgColor = chalk.bold.bgKeyword(colorMap[keyword].background).keyword(bgText)

  const keywordString = bgColor(` ${keyword.toUpperCase()} `)

  const string = `🔭 ${keywordString} ${message}`

  console.log(string)
}

export const error = (message: string): void => buildString(message, 'error')
export const success = (message: string): void => buildString(message, 'success')
export const verbose = (message: string): void => buildString(message, 'verbose')
export const info = (message: string): void => buildString(message, 'info')
export const debug = (message: string): void => buildString(message, 'debug')

export function express (req: Request, res: Response, next: Function): void {
  const { method, originalUrl, get } = req
  const { statusCode } = res

  info(`${method} ${originalUrl} ${statusCode} ${get('Content-Length')} - ${res.get('X-Response-Time')} ms`)

  next()
}
