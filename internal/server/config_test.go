/*
 * Copyright 2020-2023 Luke Whritenour

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

package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/orca-group/spirit/internal/config"
	"github.com/orca-group/spirit/internal/database"
	"github.com/stretchr/testify/require"
)

type ConfigResponse struct {
	Payload config.Cfg
	Error   string
}

var mockConfig = config.Cfg{
	Host:             "0.0.0.0",
	Port:             9000,
	CompressionLevel: 1,
	Ratelimiter:      "200x5",
	IDLength:         8,
	IDType:           "key",
	MaxSize:          400_000,
	ExpirationAge:    720,
	Headless:         false,
}

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

// checkResponseCode is a simple utility to check the response code
// of the response
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestConfig(t *testing.T) {
	mockDB := database.NewMockDatabase(t)

	s := NewServer(&mockConfig, mockDB)
	s.MountHandlers()

	req, _ := http.NewRequest("GET", "/config", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)

	x, _ := io.ReadAll(res.Result().Body)
	var body ConfigResponse
	json.Unmarshal(x, &body)

	require.Equal(t, mockConfig, body.Payload)
}
