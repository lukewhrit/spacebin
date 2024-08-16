/*
 * Copyright 2020-2024 Luke Whritenour

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

package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

	"github.com/lukewhrit/spacebin/internal/database"
	"github.com/lukewhrit/spacebin/internal/database/databasefakes"
	"github.com/lukewhrit/spacebin/internal/server"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CreateDocumentSuite struct {
	suite.Suite

	srv *server.Server
}

func (s *CreateDocumentSuite) SetupTest() {
	mockDB := &databasefakes.FakeDatabase{}

	mockDB.GetDocumentReturns(database.Document{
		ID:        "12345678",
		Content:   "test",
		CreatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
		UpdatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
	}, nil)

	// generate a struct of tests and a function to run them

	s.srv = server.NewServer(&mockConfig, mockDB)
	s.srv.MountHandlers()
}

func (s *CreateDocumentSuite) TestCreateDocument() {
	req, _ := http.NewRequest(http.MethodPost, "/api/",
		bytes.NewReader([]byte(`{"content": "test"}`)),
	)
	req.Header.Set("Content-Type", "application/json")
	rr := executeRequest(req, s.srv)

	x, _ := io.ReadAll(rr.Result().Body)
	var body DocumentResponse
	json.Unmarshal(x, &body)

	expectedResponse := DocumentResponse{
		Payload: database.Document{
			ID:        "12345678",
			Content:   "test",
			CreatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
			UpdatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
		},
	}

	require.Equal(s.T(), rr.Result().StatusCode, http.StatusOK)
	require.Equal(s.T(), expectedResponse.Payload, body.Payload)
}

func (s *CreateDocumentSuite) TestCreateMultipartDocument() {
	// Setup multipart/form-data body
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	mw.WriteField("content", "test")
	mw.Close()

	// Send request
	req, _ := http.NewRequest(http.MethodPost, "/api/", &b)
	req.Header.Add("Content-Type", mw.FormDataContentType())
	rr := executeRequest(req, s.srv)

	// Assertions
	x, _ := io.ReadAll(rr.Result().Body)
	var body DocumentResponse
	json.Unmarshal(x, &body)

	expectedResponse := DocumentResponse{
		Payload: database.Document{
			ID:        "12345678",
			Content:   "test",
			CreatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
			UpdatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
		},
	}

	require.Equal(s.T(), http.StatusOK, rr.Result().StatusCode)
	require.Equal(s.T(), expectedResponse.Payload, body.Payload)
}

func (s *CreateDocumentSuite) TestStaticCreateDocument() {
	// Setup multipart/form-data body
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	mw.WriteField("content", "test")
	mw.Close()

	// Send request
	req, _ := http.NewRequest(http.MethodPost, "/", &b)
	req.Header.Add("Content-Type", mw.FormDataContentType())
	rr := executeRequest(req, s.srv)

	// Assertions
	require.Equal(s.T(), http.StatusMovedPermanently, rr.Result().StatusCode)
	require.Equal(s.T(), "/12345678", rr.Result().Header.Get("Location"))
	// add a test for content-type and body?
}

// same as TestFetchNotFoundDocument; mocked GetDocument always returns a document, so this test needs to be reworked
// func (s *CreateDocumentSuite) TestCreateBadDocument() {
// 	req, _ := http.NewRequest(http.MethodPost, "/api/",
// 		bytes.NewReader([]byte(`{"content": "1"}`)),
// 	)
// 	req.Header.Set("Content-Type", "application/json")
// 	rr := executeRequest(req, s.srv)

// 	x, _ := io.ReadAll(rr.Result().Body)
// 	var body DocumentResponse
// 	json.Unmarshal(x, &body)

// 	require.Equal(s.T(), http.StatusBadRequest, rr.Result().StatusCode)
// 	require.Equal(s.T(), "Content: the length must be between 2 and 400000.", body.Error)
// }

func TestCreateDocumentSuite(t *testing.T) {
	suite.Run(t, new(CreateDocumentSuite))
}
