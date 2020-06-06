import test from 'ava'
import randomstring from 'randomstring'

test('ids generate with proper length', t => {
  t.is(randomstring.generate(8).length, 8)
})
