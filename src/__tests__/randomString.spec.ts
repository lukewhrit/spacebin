import randomstring from 'randomstring'

test('ids generate with proper length', () => {
  expect(randomstring.generate(8).length).toBe(8)
})
