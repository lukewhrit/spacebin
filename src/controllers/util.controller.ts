export function sanitize (input: string): string {
  // eslint-disable-next-line no-control-regex
  const ansiSequence = /(?:\x1B[@-_]|[\x80-\x9F])[0-?]*[ -/]*[@-~]/gi

  return input.replace(ansiSequence, '')
}
