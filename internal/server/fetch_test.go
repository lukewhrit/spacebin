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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lukewhrit/spacebin/internal/database"
	"github.com/lukewhrit/spacebin/internal/database/databasefakes"
	"github.com/lukewhrit/spacebin/internal/server"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DocumentResponse struct {
	Payload database.Document
	Error   string
}

type FetchDocumentSuite struct {
	suite.Suite

	srv *server.Server
}

func (s *FetchDocumentSuite) SetupTest() {
	mockDB := &databasefakes.FakeDatabase{}

	mockDB.GetDocumentReturns(database.Document{
		ID:        "12345678",
		Content:   "test",
		CreatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
		UpdatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
	}, nil)

	s.srv = server.NewServer(&mockConfig, mockDB)
	s.srv.MountHandlers()
}

func (s *FetchDocumentSuite) TestFetchDocument() {
	req, _ := http.NewRequest(http.MethodGet, "/api/12345678", nil)
	res := executeRequest(req, s.srv)

	require.Equal(s.T(), http.StatusOK, res.Result().StatusCode)

	x, _ := io.ReadAll(res.Result().Body)
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

	require.Equal(s.T(), expectedResponse.Payload, body.Payload)
}

func (s *FetchDocumentSuite) TestFetchRawDocument() {
	req, _ := http.NewRequest(http.MethodGet, "/api/12345678/raw", nil)
	res := executeRequest(req, s.srv)

	require.Equal(s.T(), http.StatusOK, res.Result().StatusCode)
	require.Equal(s.T(), "text/plain", res.Result().Header.Get("Content-Type"))
	require.Equal(s.T(), "test", res.Body.String())
}

// mocked GetDocument always returns a document, so this test needs to be reworked
// func (s *FetchDocumentSuite) TestFetchNotFoundDocument() {
// 	req, _ := http.NewRequest(http.MethodGet, "/api/12345679", nil)
// 	res := executeRequest(req, s.srv)

// 	// require.Equal(s.T(), http.StatusNotFound, res.Result().StatusCode)
// 	require.Equal(s.T(), "application/json", res.Result().Header.Get("Content-Type"))

// 	x, _ := io.ReadAll(res.Result().Body)
// 	var body DocumentResponse
// 	json.Unmarshal(x, &body)

// 	require.Equal(s.T(), "sql: no rows in result set", body.Error)
// }

// TestFetchBadIDDocument tests fetching a document with an invalid ID
func (s *FetchDocumentSuite) TestFetchBadIDDocument() {
	req, _ := http.NewRequest(http.MethodGet, "/api/1234", nil)
	res := executeRequest(req, s.srv)

	require.Equal(s.T(), http.StatusBadRequest, res.Result().StatusCode)
	require.Equal(s.T(), "application/json", res.Result().Header.Get("Content-Type"))

	x, _ := io.ReadAll(res.Result().Body)
	var body DocumentResponse
	json.Unmarshal(x, &body)

	require.Equal(s.T(), "id is of length 4, should be 8", body.Error)
}

func TestFetchDocumentSuite(t *testing.T) {
	suite.Run(t, new(FetchDocumentSuite))
}

// TestStaticDocument tests the static document handler
func TestStaticDocument(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{
		ID:        "12345678",
		Content:   "# Test\n\nThis is a test document.",
		CreatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
		UpdatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
	}, nil)

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	// Test normal document view
	req, _ := http.NewRequest(http.MethodGet, "/12345678", nil)
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "Test")
}

// TestStaticDocumentWithExtension tests static document with file extension
func TestStaticDocumentWithExtension(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{
		ID:        "12345678",
		Content:   "package main\n\nfunc main() {}",
		CreatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
		UpdatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
	}, nil)

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/12345678.go", nil)
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "package main")
}

// TestStaticDocumentReaderMode tests static document in reader mode
func TestStaticDocumentReaderMode(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{
		ID:        "12345678",
		Content:   "# Markdown Title\n\nThis is markdown content.",
		CreatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
		UpdatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
	}, nil)

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/12345678?reader=true", nil)
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "Markdown Title")
}

// TestStaticDocumentNotFound tests static document when not found
func TestStaticDocumentNotFound(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{}, sql.ErrNoRows)

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/12345678", nil)
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusNotFound, res.Result().StatusCode)
}

// TestStaticDocumentBadID tests static document with bad ID
func TestStaticDocumentBadID(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/1234", nil)
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
}

// TestFetchDocumentDatabaseError tests FetchDocument when database returns error
func TestFetchDocumentDatabaseError(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{}, errors.New("database error"))

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/api/12345678", nil)
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)

	x, _ := io.ReadAll(res.Result().Body)
	var body DocumentResponse
	json.Unmarshal(x, &body)

	require.Equal(t, "database error", body.Error)
}

// TestFetchRawDocumentDatabaseError tests FetchRawDocument when database returns error
func TestFetchRawDocumentDatabaseError(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{}, errors.New("database error"))

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/api/12345678/raw", nil)
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "database error")
}

// TestStaticDocumentDatabaseError tests StaticDocument when database returns error
func TestStaticDocumentDatabaseError(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{}, errors.New("database error"))

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/12345678", nil)
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
}

// TestFetchNotFoundDocument tests fetching a non-existent document
func TestFetchNotFoundDocument(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{}, sql.ErrNoRows)

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/api/12345678", nil)
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusNotFound, res.Result().StatusCode)
	require.Equal(t, "application/json", res.Result().Header.Get("Content-Type"))

	x, _ := io.ReadAll(res.Result().Body)
	var body DocumentResponse
	json.Unmarshal(x, &body)

	require.Equal(t, "sql: no rows in result set", body.Error)
}

// TestFetchRawNotFoundDocument tests fetching a non-existent document in raw format
func TestFetchRawNotFoundDocument(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{}, sql.ErrNoRows)

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/api/12345678/raw", nil)
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusNotFound, res.Result().StatusCode)
	require.Equal(t, "text/plain", res.Result().Header.Get("Content-Type"))
	require.Contains(t, res.Body.String(), "Document with ID 12345678 not found")
}

// TestStaticDocumentGetUsernameError tests StaticDocument when authenticatedUsername
// returns an error (e.g., GetSession fails). getUsername swallows the error and
// returns a placeholder, so the document still renders with 200.
func TestStaticDocumentGetUsernameError(t *testing.T) {
	userToken, _ := buildSessionTokens(t, "secret", "salt", "publicKey")

	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{
		ID:        "12345678",
		Content:   "hello world",
		CreatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
		UpdatedAt: time.Date(1970, 1, 1, 1, 1, 1, 1, time.UTC),
	}, nil)
	mockDB.GetSessionReturns(database.Session{}, fmt.Errorf("database error"))

	cfg := mockConfig
	cfg.AccountsEnabled = true

	srv := server.NewServer(&cfg, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/12345678", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, srv)

	require.Equal(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "hello world")
}

// TestFetchRawBadIDDocument tests fetching a document with bad ID in raw format
func TestFetchRawBadIDDocument(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/api/1234/raw", nil)
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	require.Equal(t, "application/json", res.Result().Header.Get("Content-Type"))

	x, _ := io.ReadAll(res.Result().Body)
	var body DocumentResponse
	json.Unmarshal(x, &body)

	require.Equal(t, "id is of length 4, should be 8", body.Error)
}

// TestFetchDocumentQRSuccess tests QR code generation for a document
func TestFetchDocumentQRSuccess(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{
		ID:      "12345678",
		Content: "hello",
	}, nil)

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req := httptest.NewRequest(http.MethodGet, "/12345678/qr", nil)
	req.Host = "example.com"
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusOK, res.Result().StatusCode)
	require.Equal(t, "image/png", res.Result().Header.Get("Content-Type"))
	require.Greater(t, res.Body.Len(), 0)
}

// TestFetchDocumentQRNotFound tests QR code when document does not exist
func TestFetchDocumentQRNotFound(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{}, sql.ErrNoRows)

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req := httptest.NewRequest(http.MethodGet, "/12345678/qr", nil)
	req.Host = "example.com"
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusNotFound, res.Result().StatusCode)
}

// TestFetchDocumentQRNoHost tests QR code when Host header is empty
func TestFetchDocumentQRNoHost(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{
		ID:      "12345678",
		Content: "hello",
	}, nil)

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	// Host is empty — should return 500
	req := httptest.NewRequest(http.MethodGet, "/12345678/qr", nil)
	req.Host = ""
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
}

// TestFetchDocumentQRHTTPS tests QR code with X-Forwarded-Proto https
func TestFetchDocumentQRHTTPS(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{
		ID:      "12345678",
		Content: "hello",
	}, nil)

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req := httptest.NewRequest(http.MethodGet, "/12345678/qr", nil)
	req.Host = "example.com"
	req.Header.Set("X-Forwarded-Proto", "https")
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusOK, res.Result().StatusCode)
	require.Equal(t, "image/png", res.Result().Header.Get("Content-Type"))
}

// TestFetchDocumentQRForwardedHeader tests QR with Forwarded header
func TestFetchDocumentQRForwardedHeader(t *testing.T) {
	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetDocumentReturns(database.Document{
		ID:      "12345678",
		Content: "hello",
	}, nil)

	srv := server.NewServer(&mockConfig, mockDB)
	srv.MountHandlers()

	req := httptest.NewRequest(http.MethodGet, "/12345678/qr", nil)
	req.Host = "example.com"
	req.Header.Set("Forwarded", "proto=https; host=example.com")
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusOK, res.Result().StatusCode)
}

// TestStaticDocumentOwnerButtons tests that edit/delete buttons appear for the owner
func TestStaticDocumentOwnerButtons(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "pubkey")

	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetSessionReturns(database.Session{
		Public:   "pubkey",
		Token:    userToken,
		Secret:   serverToken,
		Username: "owner",
	}, nil)
	mockDB.GetDocumentReturns(database.Document{
		ID:       "12345678",
		Content:  "owned content",
		Username: "owner",
	}, nil)

	srv := server.NewServer(&cfg, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/12345678", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "/12345678/edit")
	require.Contains(t, res.Body.String(), "/12345678/delete")
}

// TestStaticDocumentNonOwnerNoButtons tests that non-owners don't see edit/delete
func TestStaticDocumentNonOwnerNoButtons(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "pubkey")

	mockDB := &databasefakes.FakeDatabase{}
	mockDB.GetSessionReturns(database.Session{
		Public:   "pubkey",
		Token:    userToken,
		Secret:   serverToken,
		Username: "other",
	}, nil)
	mockDB.GetDocumentReturns(database.Document{
		ID:       "12345678",
		Content:  "someone elses content",
		Username: "owner",
	}, nil)

	srv := server.NewServer(&cfg, mockDB)
	srv.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/12345678", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})
	res := executeRequest(req, srv)

	require.Equal(t, http.StatusOK, res.Result().StatusCode)
	require.NotContains(t, res.Body.String(), "/12345678/edit")
	require.NotContains(t, res.Body.String(), "/12345678/delete")
}
