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

const baseURL = '/api/v1/document'

describe('Document Endpoints', () => {
  it('should fail on strange input', async () => {
    const res = await request(app)
      .post(`${baseURL}/`)
      .send({
        foo: 'this is a test',
        baz: 0,
        bar: true
      })

    expect(res.status).toBe(400)
  })

  it('should succeed on create document request', async () => {
    const res = await request(app)
      .post(`${baseURL}/`)
      .send({
        content: 'this is a test'
      })

    expect(res.body).toHaveProperty('payload.id')
    expect(res.body).toHaveProperty('payload.contentHash', '2e99758548972a8e8822ad47fa1017ff72f06f3ff6a016851f45c398732bc50c')
  })

  it('should properly return document object')

  it('should return raw content of document')

  it('should send status code 200 if document exists', async () => {
    // im unsure how to properly implement this...
    // it would require having a document we're sure exists in the database otherwise it will give a false negative.
  })

  it('should send status code 404 if document does not exist', async () => {
    const res = await request(app)
      .get(`${baseURL}/abcdefgh/verify`)

    expect(res.status).toBe(404)
    expect(res.body).toHaveProperty('payload.exists', false)
  })
})
