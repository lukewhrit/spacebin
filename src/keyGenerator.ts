// Draws inspiration from pwgen and http://tools.arantius.com/password

const randOf = (collection: string | string[]): Function => {
  return (): string => {
    return collection[Math.floor(Math.random() * collection.length)]
  }
}

export class PhoneticKeyGenerator {
  // Helper methods to get an random vowel or consonant
  private randVowel = randOf('aeiou')
  private randConsonant = randOf('bcdfghjklmnpqrstvwxyz')

  // Generate a phonetic key of alternating consonant & vowel
  createKey (keyLength: number): string {
    let text = ''

    const start = Math.round(Math.random())

    for (let i = 0; i < keyLength; i++) {
      text += (i % 2 === start) ? this.randConsonant() : this.randVowel()
    }

    return text
  }
}
