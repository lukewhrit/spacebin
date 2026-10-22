/*
 * Copyright 2020-2026 Luke Whritenour

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
	"encoding/base64"
	"encoding/json"
	"html"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lukewhrit/spacebin/internal/database"
	"github.com/lukewhrit/spacebin/internal/database/databasefakes"
	"github.com/lukewhrit/spacebin/internal/server"
	"github.com/stretchr/testify/require"
)

// envelope is a well-formed encrypted document. It doesn't decrypt to anything, but the server can't tell.
var envelope = func() string {
	encoded, _ := json.Marshal(map[string]any{
		"v":    1,
		"iter": 600000,
		"salt": base64.StdEncoding.EncodeToString(make([]byte, 16)),
		"iv":   base64.StdEncoding.EncodeToString(make([]byte, 12)),
		"ct":   base64.StdEncoding.EncodeToString(make([]byte, 32)),
	})
	return string(encoded)
}()

func encryptedDoc() database.Document {
	doc := ownedDoc()
	doc.Content = envelope
	doc.Encrypted = true
	return doc
}

func encryptedCreateBody(t *testing.T, content string) *strings.Reader {
	t.Helper()
	body, err := json.Marshal(map[string]any{"content": content, "encrypted": true})
	require.NoError(t, err)
	return strings.NewReader(string(body))
}

// --- Creating ---

func TestCreateEncryptedDocument(t *testing.T) {
	db := &databasefakes.FakeDatabase{}
	db.GetDocumentReturns(encryptedDoc(), nil)
	s := server.NewServer(&mockConfig, db)
	s.MountHandlers()

	req := httptest.NewRequest(http.MethodPost, "/api/", encryptedCreateBody(t, envelope))
	req.Header.Set("Content-Type", "application/json")
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), `"encrypted":true`)

	_, _, content, _, encrypted := db.CreateDocumentArgsForCall(0)
	require.Equal(t, envelope, content)
	require.True(t, encrypted)
}

func TestCreatePlainDocumentIsNotEncrypted(t *testing.T) {
	db := &databasefakes.FakeDatabase{}
	db.GetDocumentReturns(ownedDoc(), nil)
	s := server.NewServer(&mockConfig, db)
	s.MountHandlers()

	// Plain documents may look like envelopes, they just aren't flagged as encrypted
	req := httptest.NewRequest(http.MethodPost, "/api/", strings.NewReader(`{"content":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	_, _, _, _, encrypted := db.CreateDocumentArgsForCall(0)
	require.False(t, encrypted)
}

func TestCreateEncryptedDocumentInvalidEnvelope(t *testing.T) {
	db := &databasefakes.FakeDatabase{}
	s := server.NewServer(&mockConfig, db)
	s.MountHandlers()

	req := httptest.NewRequest(http.MethodPost, "/api/", encryptedCreateBody(t, "this is plaintext"))
	req.Header.Set("Content-Type", "application/json")
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "not a valid envelope")
	require.Zero(t, db.CreateDocumentCallCount())
}

// --- Viewing ---

func TestStaticEncryptedDocument(t *testing.T) {
	for _, path := range []string{"/12345678", "/12345678.go", "/12345678?reader=true"} {
		t.Run(path, func(t *testing.T) {
			db := &databasefakes.FakeDatabase{}
			db.GetDocumentReturns(encryptedDoc(), nil)
			s := server.NewServer(&mockConfig, db)
			s.MountHandlers()

			req, _ := http.NewRequest(http.MethodGet, path, nil)
			res := executeRequest(req, s)

			checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
			body := res.Body.String()
			require.Contains(t, body, `id="decrypt-form"`)
			require.Contains(t, body, `src="/static/crypto.js"`)
			require.Contains(t, body, `data-envelope="`+html.EscapeString(envelope)+`"`)
			require.NotContains(t, body, `class="chroma"`)
			require.NotContains(t, body, "/12345678/edit")
			require.NotContains(t, body, "{{")
		})
	}
}

func TestStaticEncryptedDocumentOwnerCanDelete(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(encryptedDoc(), nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := authedRequest(http.MethodGet, "/12345678", nil, userToken)
	res := executeRequest(req, s)

	require.Contains(t, res.Body.String(), `action="/12345678/delete"`)
	require.NotContains(t, res.Body.String(), "/12345678/edit")
}

func TestHomePageOffersEncryption(t *testing.T) {
	s := server.NewServer(&mockConfig, &databasefakes.FakeDatabase{})
	s.MountHandlers()
	s.MountStatic()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	res := executeRequest(req, s)

	body := res.Body.String()
	require.Contains(t, body, `id="encrypt"`)
	require.Contains(t, body, `src="/static/crypto.js"`)

	// The password field must never be submitted with the form
	field := body[strings.Index(body, `<input id="password"`):]
	field = field[:strings.Index(field, ">")]
	require.NotContains(t, field, "name=")
	require.NotContains(t, field, "form=")
}

// --- Editing ---

func TestEditEncryptedDocument(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		t.Run(method, func(t *testing.T) {
			cfg := mockConfig
			cfg.AccountsEnabled = true

			db, userToken := sessionFakeDB(t, "owner")
			db.GetDocumentReturns(encryptedDoc(), nil)
			s := server.NewServer(&cfg, db)
			s.MountHandlers()

			buf, ct := multipartBody(t, "content", "new content")
			req := authedRequest(method, "/12345678/edit", buf, userToken)
			req.Header.Set("Content-Type", ct)
			res := executeRequest(req, s)

			checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
			require.Contains(t, res.Body.String(), "encrypted documents can&#39;t be edited")
			require.Zero(t, db.UpdateDocumentCallCount())
		})
	}
}

// --- Collections and the account page ---

func TestStaticCollectionWithEncryptedDocument(t *testing.T) {
	db := &databasefakes.FakeDatabase{}
	db.GetCollectionReturns(ownedCollection(), nil)
	db.GetDocumentReturnsOnCall(0, encryptedDoc(), nil)
	db.GetDocumentReturnsOnCall(1, database.Document{ID: "abcdefgh", Content: "second document"}, nil)
	s := server.NewServer(&mockConfig, db)
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/c/colle123", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	body := res.Body.String()
	require.Contains(t, body, "This document is encrypted.")
	require.NotContains(t, body, `"ct"`)
	require.Contains(t, body, "second document")
}

func TestStaticSettingsPageEncryptedDocument(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentsByUsernameReturns([]database.Document{encryptedDoc()}, nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := authedRequest(http.MethodGet, "/account", nil, userToken)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), `<td class="preview-content">Encrypted</td>`)
	require.NotContains(t, res.Body.String(), "/12345678/edit")
}
