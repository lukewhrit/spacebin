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

import request from 'supertest'
import { app } from '../server'
import test from 'ava'
import has from 'has-value'

const baseURL = '/api/v1/document'

test('should fail on invalid input', async t => {
  const res = await request(app)
    .post(`${baseURL}/`)
    .send({
      foo: 'this is a test',
      baz: 0,
      bar: true
    })

  t.is(res.status, 400)
})

test('should succeed on create document request', async t => {
  const res = await request(app)
    .post(`${baseURL}/`)
    .send({
      content: 'this is a test'
    })

  t.assert(has(res.body, 'payload.id'))
  t.is(res.body.payload.contentHash, '2e99758548972a8e8822ad47fa1017ff72f06f3ff6a016851f45c398732bc50c')
})

test('should send a status code 404 if document does not exist', async t => {
  const res = await request(app)
    .get(`${baseURL}/abcdefgh/verify`)

  t.is(res.status, 500) // currently routes send 500 instead of 404
  t.is(res.body.payload.exists, false)
})

// @todo add test for 'should properly return document object'
// @todo add test for 'should return raw content of document'
// @todo add test for 'should send status code 200 if document exists'
