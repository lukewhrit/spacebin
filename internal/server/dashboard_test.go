package server_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lukewhrit/spacebin/internal/database"
	"github.com/lukewhrit/spacebin/internal/database/databasefakes"
	"github.com/lukewhrit/spacebin/internal/server"
	"github.com/stretchr/testify/require"
)

// helpers shared across dashboard tests

func authedRequest(method, target string, body *bytes.Buffer, token string) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, target, body)
	} else {
		req, _ = http.NewRequest(method, target, nil)
	}
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: token})
	return req
}

func sessionFakeDB(t *testing.T, username string) (*databasefakes.FakeDatabase, string) {
	t.Helper()
	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "pubkey")
	db := &databasefakes.FakeDatabase{}
	db.GetSessionReturns(database.Session{
		Public:   "pubkey",
		Token:    userToken,
		Secret:   serverToken,
		Username: username,
	}, nil)
	return db, userToken
}

func ownedDoc() database.Document {
	return database.Document{
		ID:        "12345678",
		Content:   "hello world content",
		Username:  "owner",
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func multipartBody(t *testing.T, field, value string) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField(field, value)
	w.Close()
	return &buf, w.FormDataContentType()
}

// --- StaticEditPage ---

func TestStaticEditPageAccountsDisabled(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = false

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/12345678/edit", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNotFound, res.Result().StatusCode)
}

func TestStaticEditPageUnauthenticated(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/12345678/edit", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusSeeOther, res.Result().StatusCode)
	require.Equal(t, "/signin", res.Result().Header.Get("Location"))
}

func TestStaticEditPageInvalidID(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := authedRequest(http.MethodGet, "/1234/edit", nil, userToken)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
}

func TestStaticEditPageDocumentNotFound(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(database.Document{}, sql.ErrNoRows)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := authedRequest(http.MethodGet, "/12345678/edit", nil, userToken)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNotFound, res.Result().StatusCode)
}

func TestStaticEditPageDocumentDatabaseError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(database.Document{}, fmt.Errorf("db error"))
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := authedRequest(http.MethodGet, "/12345678/edit", nil, userToken)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func TestStaticEditPageForbidden(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "other")
	doc := ownedDoc() // belongs to "owner", not "other"
	db.GetDocumentReturns(doc, nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := authedRequest(http.MethodGet, "/12345678/edit", nil, userToken)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusForbidden, res.Result().StatusCode)
}

func TestStaticEditPageSuccess(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(ownedDoc(), nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := authedRequest(http.MethodGet, "/12345678/edit", nil, userToken)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "hello world content")
	require.NotContains(t, res.Body.String(), "{{")
}

// --- EditDocument ---

func TestEditDocumentAccountsDisabled(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = false

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	buf, ct := multipartBody(t, "content", "updated content")
	req := httptest.NewRequest(http.MethodPost, "/12345678/edit", buf)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNotFound, res.Result().StatusCode)
}

func TestEditDocumentUnauthenticated(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	buf, ct := multipartBody(t, "content", "updated content")
	req := httptest.NewRequest(http.MethodPost, "/12345678/edit", buf)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusUnauthorized, res.Result().StatusCode)
}

func TestEditDocumentInvalidID(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "content", "updated content")
	req := authedRequest(http.MethodPost, "/1234/edit", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
}

func TestEditDocumentNotFound(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(database.Document{}, sql.ErrNoRows)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "content", "updated content")
	req := authedRequest(http.MethodPost, "/12345678/edit", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNotFound, res.Result().StatusCode)
}

func TestEditDocumentGetDocumentDatabaseError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(database.Document{}, fmt.Errorf("db error"))
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "content", "updated content")
	req := authedRequest(http.MethodPost, "/12345678/edit", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func TestEditDocumentForbidden(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "other")
	db.GetDocumentReturns(ownedDoc(), nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "content", "updated content")
	req := authedRequest(http.MethodPost, "/12345678/edit", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusForbidden, res.Result().StatusCode)
}

func TestEditDocumentBadContent(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(ownedDoc(), nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	// Content too short (min 2 chars)
	buf, ct := multipartBody(t, "content", "x")
	req := authedRequest(http.MethodPost, "/12345678/edit", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
}

func TestEditDocumentUpdateDatabaseError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(ownedDoc(), nil)
	db.UpdateDocumentReturns(fmt.Errorf("db error"))
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "content", "updated content")
	req := authedRequest(http.MethodPost, "/12345678/edit", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func TestEditDocumentSuccessRedirect(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(ownedDoc(), nil)
	db.UpdateDocumentReturns(nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "content", "updated content here")
	req := authedRequest(http.MethodPost, "/12345678/edit", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusSeeOther, res.Result().StatusCode)
	require.Equal(t, "/12345678", res.Result().Header.Get("Location"))
}

func TestEditDocumentSuccessJSON(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(ownedDoc(), nil)
	db.UpdateDocumentReturns(nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "content", "updated content here")
	req := authedRequest(http.MethodPost, "/12345678/edit", buf, userToken)
	req.Header.Set("Content-Type", ct)
	req.Header.Set("Accept", "application/json")
	// Override to JSON for response detection
	req2 := httptest.NewRequest(http.MethodPost, "/12345678/edit", buf)
	req2.Header.Set("Content-Type", "application/json")
	req2.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})
	// Use multipart since AllowContentType requires it, but set Accept: json
	_ = req2
	// JSON path: set Content-Type to application/json triggers wantsJSON
	var jsonBuf bytes.Buffer
	jsonBuf.WriteString(`{"content": "updated content here"}`)
	req3 := httptest.NewRequest(http.MethodPost, "/12345678/edit", &jsonBuf)
	req3.Header.Set("Content-Type", "application/json")
	req3.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})
	res := executeRequest(req3, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "12345678")
}

// --- RemoveDocument ---

func TestRemoveDocumentAccountsDisabled(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = false

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := httptest.NewRequest(http.MethodPost, "/12345678/delete", buf)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNotFound, res.Result().StatusCode)
}

func TestRemoveDocumentUnauthenticated(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := httptest.NewRequest(http.MethodPost, "/12345678/delete", buf)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusUnauthorized, res.Result().StatusCode)
}

func TestRemoveDocumentInvalidID(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := authedRequest(http.MethodPost, "/1234/delete", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
}

func TestRemoveDocumentNotFound(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(database.Document{}, sql.ErrNoRows)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := authedRequest(http.MethodPost, "/12345678/delete", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNotFound, res.Result().StatusCode)
}

func TestRemoveDocumentGetDocumentDatabaseError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(database.Document{}, fmt.Errorf("db error"))
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := authedRequest(http.MethodPost, "/12345678/delete", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func TestRemoveDocumentForbidden(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "other")
	db.GetDocumentReturns(ownedDoc(), nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := authedRequest(http.MethodPost, "/12345678/delete", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusForbidden, res.Result().StatusCode)
}

func TestRemoveDocumentDeleteDatabaseError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(ownedDoc(), nil)
	db.DeleteDocumentReturns(fmt.Errorf("db error"))
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := authedRequest(http.MethodPost, "/12345678/delete", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func TestRemoveDocumentSuccess(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(ownedDoc(), nil)
	db.DeleteDocumentReturns(nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := authedRequest(http.MethodPost, "/12345678/delete", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusSeeOther, res.Result().StatusCode)
	require.Equal(t, "/account", res.Result().Header.Get("Location"))
	require.Equal(t, 1, db.DeleteDocumentCallCount())
}

func TestRemoveDocumentSuccessJSON(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetDocumentReturns(ownedDoc(), nil)
	db.DeleteDocumentReturns(nil)
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	req := httptest.NewRequest(http.MethodPost, "/12345678/delete", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNoContent, res.Result().StatusCode)
}

// --- RemoveAccount ---

func TestRemoveAccountAccountsDisabled(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = false

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := httptest.NewRequest(http.MethodPost, "/account/delete", buf)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNotFound, res.Result().StatusCode)
}

func TestRemoveAccountUnauthenticated(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := httptest.NewRequest(http.MethodPost, "/account/delete", buf)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusSeeOther, res.Result().StatusCode)
	require.Equal(t, "/signin", res.Result().Header.Get("Location"))
}

func TestRemoveAccountGetAccountDatabaseError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	db, userToken := sessionFakeDB(t, "owner")
	db.GetAccountByUsernameReturns(database.Account{}, fmt.Errorf("db error"))
	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := authedRequest(http.MethodPost, "/account/delete", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func TestRemoveAccountDeleteDatabaseError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "pubkey")
	db := &databasefakes.FakeDatabase{}
	db.GetSessionReturns(database.Session{
		Public:   "pubkey",
		Token:    userToken,
		Secret:   serverToken,
		Username: "owner",
	}, nil)
	db.GetAccountByUsernameReturns(database.Account{ID: 1, Username: "owner"}, nil)
	db.DeleteAccountReturns(fmt.Errorf("db error"))

	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := authedRequest(http.MethodPost, "/account/delete", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func TestRemoveAccountSuccess(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "pubkey")
	db := &databasefakes.FakeDatabase{}
	db.GetSessionReturns(database.Session{
		Public:   "pubkey",
		Token:    userToken,
		Secret:   serverToken,
		Username: "owner",
	}, nil)
	db.GetAccountByUsernameReturns(database.Account{ID: 1, Username: "owner"}, nil)
	db.DeleteAccountReturns(nil)
	db.DeleteSessionReturns(nil)

	s := server.NewServer(&cfg, db)
	s.MountHandlers()

	buf, ct := multipartBody(t, "", "")
	req := authedRequest(http.MethodPost, "/account/delete", buf, userToken)
	req.Header.Set("Content-Type", ct)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusSeeOther, res.Result().StatusCode)
	require.Equal(t, "/", res.Result().Header.Get("Location"))
	require.Equal(t, 1, db.DeleteAccountCallCount())

	// Session cookie should be cleared
	foundCookie := false
	for _, c := range res.Result().Cookies() {
		if c.Name == "spacebin_token" {
			foundCookie = true
			require.Equal(t, "", c.Value)
		}
	}
	require.True(t, foundCookie)
}
