// Draws inspiration from pwgen and http://tools.arantius.com/password

export class PhoneticKeyGenerator {
  /**
   * Returns random value/character from array or string.
   *
   * @param collection - The string or array of which to choose a random value from.
   *
   * @returns A random value from a string or array.
   */
  private randOf = (collection: string | string[]): Function => {
    return (): string => {
      return collection[Math.floor(Math.random() * collection.length)]
    }
  }

  private randVowel = this.randOf('aeiou')
  private randConsonant = this.randOf('bcdfghjklmnpqrstvwxyz')

  /**
   * Generates a pseudo-random pronounceable key
   *
   * @param keyLength - The length of the key to generate
   *
   * @returns - Random key
   */
  createKey (keyLength: number): string {
    let key = ''

    const start = Math.round(Math.random())

    for (let i = 0; i < keyLength; i++) {
      key += (i % 2 === start) ? this.randConsonant() : this.randVowel()
    }

    return key
  }
}
