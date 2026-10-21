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
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lukewhrit/spacebin/internal/database"
	"github.com/lukewhrit/spacebin/internal/database/databasefakes"
	"github.com/lukewhrit/spacebin/internal/server"
	"github.com/stretchr/testify/require"
)

// helpers shared across collection tests

func ownedCollection() database.Collection {
	return database.Collection{
		ID:        "colle123",
		Title:     "My notes",
		Username:  "owner",
		Documents: []string{"12345678", "abcdefgh"},
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func collectionJSON(t *testing.T, title string, documents ...string) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, json.NewEncoder(&buf).Encode(map[string]any{
		"title":     title,
		"documents": documents,
	}))
	return &buf
}

func collectionForm(t *testing.T, title, documents string) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("title", title)
	w.WriteField("documents", documents)
	w.Close()
	return &buf, w.FormDataContentType()
}

func newCollectionServer(db *databasefakes.FakeDatabase) *server.Server {
	s := server.NewServer(&mockConfig, db)
	s.MountHandlers()
	return s
}

// --- CreateCollection ---

func TestCreateCollection(t *testing.T) {
	db := &databasefakes.FakeDatabase{}
	db.GetDocumentReturns(ownedDoc(), nil)
	db.GetCollectionReturns(ownedCollection(), nil)
	s := newCollectionServer(db)

	req := httptest.NewRequest(http.MethodPost, "/api/collections", collectionJSON(t, "  My notes ", "12345678", "abcdefgh", "12345678"))
	req.Header.Set("Content-Type", "application/json")
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), `"id":"colle123"`)

	require.Equal(t, 1, db.CreateCollectionCallCount())
	_, id, title, username, documents := db.CreateCollectionArgsForCall(0)
	require.Len(t, id, mockConfig.IDLength)
	require.Equal(t, "My notes", title)
	require.Equal(t, "", username)
	require.Equal(t, []string{"12345678", "abcdefgh"}, documents)
}

func TestCreateCollectionTracksOwner(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(ownedDoc(), nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := authedRequest(http.MethodPost, "/api/collections", collectionJSON(t, "", "12345678"), userToken)
	req.Header.Set("Content-Type", "application/json")
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	_, _, _, username, _ := db.CreateCollectionArgsForCall(0)
	require.Equal(t, "owner", username)
}

func TestCreateCollectionBadRequests(t *testing.T) {
	tests := map[string]struct {
		body    *bytes.Buffer
		message string
	}{
		"no documents":      {collectionJSON(t, "title"), "at least one document"},
		"title too long":    {collectionJSON(t, strings.Repeat("a", 256), "12345678"), "title must be at most 255"},
		"invalid id":        {collectionJSON(t, "", "short"), "not a valid document id"},
		"invalid JSON body": {bytes.NewBufferString("{invalid"), "bad request"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := &databasefakes.FakeDatabase{}
			db.GetDocumentReturns(ownedDoc(), nil)
			s := newCollectionServer(db)

			req := httptest.NewRequest(http.MethodPost, "/api/collections", test.body)
			req.Header.Set("Content-Type", "application/json")
			res := executeRequest(req, s)

			checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
			require.Contains(t, res.Body.String(), test.message)
			require.Zero(t, db.CreateCollectionCallCount())
		})
	}
}

func TestCreateCollectionTooManyDocuments(t *testing.T) {
	cfg := mockConfig
	cfg.MaxCollectionSize = 1

	db := &databasefakes.FakeDatabase{}
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := httptest.NewRequest(http.MethodPost, "/api/collections", collectionJSON(t, "", "12345678", "abcdefgh"))
	req.Header.Set("Content-Type", "application/json")
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "at most 1 documents")
}

func TestCreateCollectionDocumentNotFound(t *testing.T) {
	db := &databasefakes.FakeDatabase{}
	db.GetDocumentReturns(database.Document{}, sql.ErrNoRows)
	s := newCollectionServer(db)

	req := httptest.NewRequest(http.MethodPost, "/api/collections", collectionJSON(t, "", "12345678"))
	req.Header.Set("Content-Type", "application/json")
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "document 12345678 not found")
}

func TestCreateCollectionDatabaseErrors(t *testing.T) {
	tests := map[string]func(db *databasefakes.FakeDatabase){
		"get document": func(db *databasefakes.FakeDatabase) {
			db.GetDocumentReturns(database.Document{}, errors.New("db down"))
		},
		"create collection": func(db *databasefakes.FakeDatabase) {
			db.GetDocumentReturns(ownedDoc(), nil)
			db.CreateCollectionReturns(errors.New("db down"))
		},
		"get collection": func(db *databasefakes.FakeDatabase) {
			db.GetDocumentReturns(ownedDoc(), nil)
			db.GetCollectionReturns(database.Collection{}, errors.New("db down"))
		},
	}

	for name, setup := range tests {
		t.Run(name, func(t *testing.T) {
			db := &databasefakes.FakeDatabase{}
			setup(db)
			s := newCollectionServer(db)

			req := httptest.NewRequest(http.MethodPost, "/api/collections", collectionJSON(t, "", "12345678"))
			req.Header.Set("Content-Type", "application/json")
			res := executeRequest(req, s)

			checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
		})
	}
}

// --- StaticCreateCollection ---

func TestStaticCreateCollection(t *testing.T) {
	db := &databasefakes.FakeDatabase{}
	db.GetDocumentReturns(ownedDoc(), nil)
	s := newCollectionServer(db)

	buf, ct := collectionForm(t, "My notes", "https://spaceb.in/12345678.go\nabcdefgh")
	req := httptest.NewRequest(http.MethodPost, "/c/new", buf)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusSeeOther, res.Result().StatusCode)

	_, id, _, _, documents := db.CreateCollectionArgsForCall(0)
	require.Equal(t, "/c/"+id, res.Result().Header.Get("Location"))
	require.Equal(t, []string{"12345678", "abcdefgh"}, documents)
}

func TestStaticCreateCollectionBadRequest(t *testing.T) {
	s := newCollectionServer(&databasefakes.FakeDatabase{})

	buf, ct := collectionForm(t, "My notes", "")
	req := httptest.NewRequest(http.MethodPost, "/c/new", buf)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "at least one document")
}

func TestStaticCreateCollectionDatabaseError(t *testing.T) {
	db := &databasefakes.FakeDatabase{}
	db.GetDocumentReturns(ownedDoc(), nil)
	db.CreateCollectionReturns(errors.New("db down"))
	s := newCollectionServer(db)

	buf, ct := collectionForm(t, "", "12345678")
	req := httptest.NewRequest(http.MethodPost, "/c/new", buf)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

// --- FetchCollection ---

func TestFetchCollection(t *testing.T) {
	db := &databasefakes.FakeDatabase{}
	db.GetCollectionReturns(ownedCollection(), nil)
	s := newCollectionServer(db)

	req, _ := http.NewRequest(http.MethodGet, "/api/collections/colle123", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)

	var body struct {
		Payload database.Collection
		Error   string
	}
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &body))
	require.Equal(t, "My notes", body.Payload.Title)
	require.Equal(t, []string{"12345678", "abcdefgh"}, body.Payload.Documents)
}

func TestFetchCollectionErrors(t *testing.T) {
	tests := map[string]struct {
		path   string
		err    error
		status int
	}{
		"invalid id":     {"/api/collections/short", nil, http.StatusBadRequest},
		"not found":      {"/api/collections/colle123", sql.ErrNoRows, http.StatusNotFound},
		"database error": {"/api/collections/colle123", errors.New("db down"), http.StatusInternalServerError},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := &databasefakes.FakeDatabase{}
			db.GetCollectionReturns(database.Collection{}, test.err)
			s := newCollectionServer(db)

			req, _ := http.NewRequest(http.MethodGet, test.path, nil)
			res := executeRequest(req, s)

			checkResponseCode(t, test.status, res.Result().StatusCode)
		})
	}
}

// --- StaticCollection ---

func TestStaticCollection(t *testing.T) {
	db := &databasefakes.FakeDatabase{}
	db.GetCollectionReturns(ownedCollection(), nil)
	db.GetDocumentReturnsOnCall(0, database.Document{ID: "12345678", Content: "first document"}, nil)
	db.GetDocumentReturnsOnCall(1, database.Document{ID: "abcdefgh", Content: "second document"}, nil)
	s := newCollectionServer(db)

	req, _ := http.NewRequest(http.MethodGet, "/c/colle123", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	body := res.Body.String()
	require.Contains(t, body, "My notes")
	require.Contains(t, body, `href="/12345678"`)
	require.Contains(t, body, "second document")
	require.Less(t, strings.Index(body, "first document"), strings.Index(body, "second document"))

	// Each document gets a tab, and only the first starts selected
	require.Contains(t, body, `id="tab-12345678" class="tab-input" checked>`)
	require.Contains(t, body, `id="tab-abcdefgh" class="tab-input" >`)
	require.Contains(t, body, `<label for="tab-abcdefgh" class="tab-label">abcdefgh</label>`)
	require.NotContains(t, body, "/c/colle123/delete")
	require.NotContains(t, body, "{{")
}

func TestStaticCollectionShowsDeleteToOwner(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetCollectionReturns(ownedCollection(), nil)
	db.GetDocumentReturns(ownedDoc(), nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := authedRequest(http.MethodGet, "/c/colle123", nil, userToken)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "/c/colle123/delete")
}

func TestStaticCollectionSkipsDeletedDocument(t *testing.T) {
	db := &databasefakes.FakeDatabase{}
	db.GetCollectionReturns(ownedCollection(), nil)
	db.GetDocumentReturnsOnCall(0, database.Document{}, sql.ErrNoRows)
	db.GetDocumentReturnsOnCall(1, database.Document{ID: "abcdefgh", Content: "second document"}, nil)
	s := newCollectionServer(db)

	req, _ := http.NewRequest(http.MethodGet, "/c/colle123", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.NotContains(t, res.Body.String(), `href="/12345678"`)
	require.Contains(t, res.Body.String(), "second document")
}

func TestStaticCollectionEmpty(t *testing.T) {
	db := &databasefakes.FakeDatabase{}
	db.GetCollectionReturns(database.Collection{ID: "colle123"}, nil)
	s := newCollectionServer(db)

	req, _ := http.NewRequest(http.MethodGet, "/c/colle123", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "has been deleted")
}

func TestStaticCollectionErrors(t *testing.T) {
	tests := map[string]struct {
		path          string
		collectionErr error
		documentErr   error
		status        int
	}{
		"invalid id":              {"/c/short", nil, nil, http.StatusBadRequest},
		"not found":               {"/c/colle123", sql.ErrNoRows, nil, http.StatusNotFound},
		"collection database err": {"/c/colle123", errors.New("db down"), nil, http.StatusInternalServerError},
		"document database err":   {"/c/colle123", nil, errors.New("db down"), http.StatusInternalServerError},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := &databasefakes.FakeDatabase{}
			db.GetCollectionReturns(ownedCollection(), test.collectionErr)
			db.GetDocumentReturns(database.Document{}, test.documentErr)
			s := newCollectionServer(db)

			req, _ := http.NewRequest(http.MethodGet, test.path, nil)
			res := executeRequest(req, s)

			checkResponseCode(t, test.status, res.Result().StatusCode)
		})
	}
}

// --- StaticNewCollection ---

func TestStaticNewCollection(t *testing.T) {
	s := newCollectionServer(&databasefakes.FakeDatabase{})

	req, _ := http.NewRequest(http.MethodGet, "/c/new?documents=12345678,abcdefgh", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "12345678\nabcdefgh</textarea>")
	require.Contains(t, res.Body.String(), "Up to 25 document")
	require.NotContains(t, res.Body.String(), "{{")
}

func TestHomePageLinksToNewCollection(t *testing.T) {
	s := newCollectionServer(&databasefakes.FakeDatabase{})
	s.MountStatic()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), `href="/c/new"`)
}

func TestStaticDocumentHasNoCollectionButton(t *testing.T) {
	db := &databasefakes.FakeDatabase{}
	db.GetDocumentReturns(ownedDoc(), nil)
	s := newCollectionServer(db)

	req, _ := http.NewRequest(http.MethodGet, "/12345678", nil)
	res := executeRequest(req, s)

	require.NotContains(t, res.Body.String(), "/c/new")
}

// --- RemoveCollection ---

func TestRemoveCollection(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetCollectionReturns(ownedCollection(), nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := authedRequest(http.MethodPost, "/c/colle123/delete", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusSeeOther, res.Result().StatusCode)
	require.Equal(t, "/account", res.Result().Header.Get("Location"))
	require.Equal(t, 1, db.DeleteCollectionCallCount())
}

func TestRemoveCollectionJSON(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetCollectionReturns(ownedCollection(), nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := authedRequest(http.MethodPost, "/c/colle123/delete", nil, userToken)
	req.Header.Set("Content-Type", "application/json")
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNoContent, res.Result().StatusCode)
}

func TestRemoveCollectionErrors(t *testing.T) {
	tests := map[string]struct {
		accounts      bool
		authenticated bool
		path          string
		collection    database.Collection
		getErr        error
		deleteErr     error
		status        int
	}{
		"accounts disabled": {false, true, "/c/colle123/delete", ownedCollection(), nil, nil, http.StatusNotFound},
		"unauthenticated":   {true, false, "/c/colle123/delete", ownedCollection(), nil, nil, http.StatusUnauthorized},
		"invalid id":        {true, true, "/c/short/delete", ownedCollection(), nil, nil, http.StatusBadRequest},
		"not found":         {true, true, "/c/colle123/delete", database.Collection{}, sql.ErrNoRows, nil, http.StatusNotFound},
		"get database err":  {true, true, "/c/colle123/delete", database.Collection{}, errors.New("db down"), nil, http.StatusInternalServerError},
		"not the owner":     {true, true, "/c/colle123/delete", database.Collection{ID: "colle123", Username: "someone"}, nil, nil, http.StatusForbidden},
		"delete db err":     {true, true, "/c/colle123/delete", ownedCollection(), nil, errors.New("db down"), http.StatusInternalServerError},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cfg := mockConfig
			cfg.AccountsEnabled = test.accounts

			db, userToken := sessionFakeDB(t, "owner")
			db.GetCollectionReturns(test.collection, test.getErr)
			db.DeleteCollectionReturns(test.deleteErr)
			s := server.NewServer(&cfg, db)
			s.MountHandlers()

			if !test.authenticated {
				userToken = ""
			}

			req := authedRequest(http.MethodPost, test.path, nil, userToken)
			res := executeRequest(req, s)

			checkResponseCode(t, test.status, res.Result().StatusCode)

			// Only a delete that fails in the database gets as far as deleting
			if test.deleteErr == nil {
				require.Zero(t, db.DeleteCollectionCallCount())
			}
		})
	}
}

// --- Account page ---

func TestStaticSettingsPageListsCollections(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetCollectionsByUsernameReturns([]database.Collection{ownedCollection()}, nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := authedRequest(http.MethodGet, "/account", nil, userToken)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), `href="/c/colle123"`)
	require.Contains(t, res.Body.String(), "My notes")
	require.NotContains(t, res.Body.String(), "{{")
}

func TestStaticSettingsPageGetCollectionsError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetCollectionsByUsernameReturns(nil, errors.New("db down"))
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := authedRequest(http.MethodGet, "/account", nil, userToken)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}
