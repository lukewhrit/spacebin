import test from 'ava'
import { PhoneticKeyGenerator } from '../controllers/KeyGenerator'

test('createKey() returns key of proper length ', t => {
  const g = new PhoneticKeyGenerator()

  t.is(g.createKey(6).length, 6)
})

const vowels = 'aeiou'
const consonants = 'bcdfghjklmnpqrstvwxyz'

test('createKey() alternates consonants and vowels', t => {
  const g = new PhoneticKeyGenerator()
  const k = g.createKey(3)

  t.true(consonants.includes(k[0]))
  t.true(consonants.includes(k[2]))
  t.true(vowels.includes(k[1]))
})
