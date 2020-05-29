import test from 'ava'
import { PhoneticKeyGenerator } from '../keyGenerator'

test('createKey() returns key of proper length ', t => {
  const g = new PhoneticKeyGenerator()

  t.is(g.createKey(6).length, 6)
})
