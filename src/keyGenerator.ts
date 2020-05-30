// Draws inspiration from pwgen and http://tools.arantius.com/password

export class PhoneticKeyGenerator {
  private randOf = (collection: string | string[]): Function => {
    return (): string => {
      return collection[Math.floor(Math.random() * collection.length)]
    }
  }

  // Helper methods to get an random vowel or consonant
  private randVowel = this.randOf('aeiou')
  private randConsonant = this.randOf('bcdfghjklmnpqrstvwxyz')

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
