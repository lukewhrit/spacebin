/*
 * Copyright 2020 Luke Whrit, Jack Dorland; The Spacebin Authors

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 *     http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package document

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber"
	"github.com/spacebin-org/spirit/server"
)

type Test struct {
	description string
	route       string
	method      string
	inputBody   map[string]interface{}

	expectedError bool
	expectedCode  int
	expectedBody  string
}

func TestCreate(t *testing.T) {
	// Currently requires a document with the ID ""
	// in the database. It'll be like this until we find a
	// better solution.

	// We should probably get the ID we retrieved
	// after the "create a document" test succeeds.

	// Currently, we can't really implement the read tests
	// in a good way without what I previously mentioned.

	ExecuteTest(t, server.Start(), []*Test{
		// We can't properly implement this test without a way
		// of getting the id it returns.
		{
			description: "upload a document",
			route:       "/api/v1/documents/",
			method:      "POST",
			inputBody: map[string]interface{}{
				"content":   "this is a test",
				"extension": "txt",
			},

			expectedBody:  `{}`,
			expectedCode:  200,
			expectedError: false,
		},
		{
			description: "get a document's plain text content",
			route:       "/api/v1/documents/:id",
			method:      "GET",

			expectedBody:  "this is a test",
			expectedCode:  200,
			expectedError: false,
		},
		{
			description: "get a document",
			route:       "/api/v1/documents/:id",
			method:      "GET",

			expectedBody: fmt.Sprintf(`{
				"status": 200,
				"payload": {
					"id":        %s,
					"content":   "this is a test",
					"extension": "txt",
					"created_at": %d,
					"updated_at": %d,
				},
				"error": "",
			}`, "", time.Now().Unix(), time.Now().Unix()),
			expectedCode:  200,
			expectedError: false,
		},
	})
}

// ExecuteTest runs tests defined in a tests object
func ExecuteTest(t *testing.T, app *fiber.App, tests []*Test) {
	t.Helper()

	for _, test := range tests {
		req, _ := http.NewRequest(
			test.method,
			test.route,
			nil)

		res, err := app.Test(req, -1)

		if !test.expectedError && (err != nil) {
			t.Errorf("did not expect error. got=%s [%s]", err, test.description)
			continue
		}

		// If this test was expected to fail and did fail go to the next iteration in the loop
		if test.expectedError {
			continue
		}

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			t.Errorf("body read error. %s [%s]", err, test.description)
			continue
		}

		actual := string(body)
		if test.expectedBody != actual {
			t.Errorf("body has wrong value. got=%s, want=%s [%s]",
				test.expectedBody, actual, test.description)
			continue
		}

		if test.expectedCode != res.StatusCode {
			t.Errorf("statusCode has wrong value. got=%d, want=%d [%s]",
				test.expectedCode, res.StatusCode, test.description)
			continue
		}
	}
}
