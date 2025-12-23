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
	"errors"
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

func TestCreateDocumentSuite(t *testing.T) {
	suite.Run(t, new(CreateDocumentSuite))
}

// TestCreateBadContentDocument tests creating a document with invalid content
func TestCreateBadContentDocument(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	// Test with content too short
	req, _ := http.NewRequest(http.MethodPost, "/api/",
		bytes.NewReader([]byte(`{"content": "x"}`)),
	)
	req.Header.Set("Content-Type", "application/json")
	rr := executeRequest(req, srv)

	require.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)

	x, _ := io.ReadAll(rr.Result().Body)
	var body DocumentResponse
	json.Unmarshal(x, &body)

	require.Contains(t, body.Error, "bad request")
}

// TestCreateEmptyContentDocument tests creating a document with empty content
func TestCreateEmptyContentDocument(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodPost, "/api/",
		bytes.NewReader([]byte(`{"content": ""}`)),
	)
	req.Header.Set("Content-Type", "application/json")
	rr := executeRequest(req, srv)

	require.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)

	x, _ := io.ReadAll(rr.Result().Body)
	var body DocumentResponse
	json.Unmarshal(x, &body)

	require.Contains(t, body.Error, "bad request")
}

// TestStaticCreateBadContentDocument tests static create with bad content
func TestStaticCreateBadContentDocument(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	// Setup multipart/form-data body with content too short
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	mw.WriteField("content", "x")
	mw.Close()

	req, _ := http.NewRequest(http.MethodPost, "/", &b)
	req.Header.Add("Content-Type", mw.FormDataContentType())
	rr := executeRequest(req, srv)

	require.Equal(t, http.StatusBadRequest, rr.Result().StatusCode)
}

// TestCreateDocumentDatabaseError tests creating a document when database fails
func TestCreateDocumentDatabaseError(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.CreateDocumentReturns(errors.New("database error"))

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodPost, "/api/",
		bytes.NewReader([]byte(`{"content": "test"}`)),
	)
	req.Header.Set("Content-Type", "application/json")
	rr := executeRequest(req, srv)

	require.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)

	x, _ := io.ReadAll(rr.Result().Body)
	var body DocumentResponse
	json.Unmarshal(x, &body)

	require.Equal(t, "database error", body.Error)
}

// TestStaticCreateDocumentDatabaseError tests static create when database fails
func TestStaticCreateDocumentDatabaseError(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.CreateDocumentReturns(errors.New("database error"))

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	mw.WriteField("content", "test")
	mw.Close()

	req, _ := http.NewRequest(http.MethodPost, "/", &b)
	req.Header.Add("Content-Type", mw.FormDataContentType())
	rr := executeRequest(req, srv)

	require.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
}

// TestCreateDocumentGetDocumentError tests when GetDocument fails after CreateDocument
func TestCreateDocumentGetDocumentError(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{}, errors.New("get document error"))

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodPost, "/api/",
		bytes.NewReader([]byte(`{"content": "test"}`)),
	)
	req.Header.Set("Content-Type", "application/json")
	rr := executeRequest(req, srv)

	require.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)

	x, _ := io.ReadAll(rr.Result().Body)
	var body DocumentResponse
	json.Unmarshal(x, &body)

	require.Equal(t, "get document error", body.Error)
}

// TestStaticCreateDocumentGetDocumentError tests when GetDocument fails after StaticCreateDocument
func TestStaticCreateDocumentGetDocumentError(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{}, errors.New("get document error"))

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	mw.WriteField("content", "test")
	mw.Close()

	req, _ := http.NewRequest(http.MethodPost, "/", &b)
	req.Header.Add("Content-Type", mw.FormDataContentType())
	rr := executeRequest(req, srv)

	require.Equal(t, http.StatusInternalServerError, rr.Result().StatusCode)
}
