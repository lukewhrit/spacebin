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
})
