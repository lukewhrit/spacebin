package server_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lukewhrit/spacebin/internal/database"
	"github.com/lukewhrit/spacebin/internal/database/databasefakes"
	"github.com/lukewhrit/spacebin/internal/server"
	"github.com/lukewhrit/spacebin/internal/util"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

func TestSignUpSuccess(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturnsOnCall(0, database.Account{}, sql.ErrNoRows)
	fakeDB.GetAccountByUsernameReturnsOnCall(1, database.Account{
		ID:       1,
		Username: "newuser",
	}, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{
		"username": "newuser",
		"password": "strongpassword",
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
}

func TestSignUpDuplicateUsername(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturns(database.Account{ID: 1, Username: "existing"}, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{
		"username": "existing",
		"password": "strongpassword",
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusConflict, res.Result().StatusCode)
}

func TestSignInInvalidCredentials(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturns(database.Account{}, sql.ErrNoRows)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{
		"username": "missing",
		"password": "password",
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusUnauthorized, res.Result().StatusCode)
}

func TestSignInPasswordMismatch(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturns(database.Account{
		ID:       1,
		Username: "user",
		Password: string(hashedPassword),
	}, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{
		"username": "user",
		"password": "wrong-password",
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusUnauthorized, res.Result().StatusCode)
}

func TestSignInSetsCookieAndSessionUsername(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturns(database.Account{
		ID:       1,
		Username: "user",
		Password: string(hashedPassword),
	}, nil)

	var capturedUsername string
	fakeDB.CreateSessionStub = func(ctx context.Context, public, token, secret, username string) error {
		capturedUsername = username
		return nil
	}

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{
		"username": "user",
		"password": "correct-password",
	})

	req, _ := http.NewRequest(http.MethodPost, "/api/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Equal(t, "user", capturedUsername)

	duration := time.Duration(cfg.SessionTTLHours) * time.Hour
	minExpiry := start.Add(duration - time.Second)
	maxExpiry := start.Add(duration + time.Second)
	expectedMaxAge := int(duration.Seconds())

	foundCookie := false
	for _, c := range res.Result().Cookies() {
		if c.Name != "spacebin_token" || c.Value == "" {
			continue
		}

		foundCookie = true
		require.Equal(t, "/", c.Path)
		require.Equal(t, http.SameSiteLaxMode, c.SameSite)
		require.Equal(t, expectedMaxAge, c.MaxAge)
		require.True(t, c.Expires.After(minExpiry) && c.Expires.Before(maxExpiry))
		require.True(t, c.HttpOnly)
		require.False(t, c.Secure)
		require.Empty(t, c.Domain)
	}

	require.True(t, foundCookie)
}

func TestAuthenticationDisabled(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = false

	fakeDB := &databasefakes.FakeDatabase{}

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{
		"username": "user",
		"password": "password",
	})

	signUpReq, _ := http.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(body))
	signUpReq.Header.Set("Content-Type", "application/json")

	signInReq, _ := http.NewRequest(http.MethodPost, "/api/signin", bytes.NewBuffer(body))
	signInReq.Header.Set("Content-Type", "application/json")

	signUpRes := executeRequest(signUpReq, s)
	signInRes := executeRequest(signInReq, s)

	checkResponseCode(t, http.StatusNotFound, signUpRes.Result().StatusCode)
	checkResponseCode(t, http.StatusNotFound, signInRes.Result().StatusCode)
}

func TestStaticIndexAuthenticated(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    userToken,
		Secret:   serverToken,
		Username: "tester",
	}, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountStatic()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "tester")
}

func TestSignInRedirectsWithCookie(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturns(database.Account{
		ID:       1,
		Username: "user",
		Password: string(hashedPassword),
	}, nil)

	fakeDB.CreateSessionReturns(nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("username", "user")
	writer.WriteField("password", "correct-password")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/signin", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusSeeOther, res.Result().StatusCode)
	require.Equal(t, "/", res.Result().Header.Get("Location"))

	foundCookie := false
	for _, c := range res.Result().Cookies() {
		if c.Name == "spacebin_token" && c.Value != "" {
			foundCookie = true
		}
	}

	require.True(t, foundCookie)
}

func TestSignInCookieSecureWithHTTPS(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturns(database.Account{
		ID:       1,
		Username: "user",
		Password: string(hashedPassword),
	}, nil)
	fakeDB.CreateSessionReturns(nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{
		"username": "user",
		"password": "correct-password",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.TLS = &tls.ConnectionState{}

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)

	foundCookie := false
	for _, c := range res.Result().Cookies() {
		if c.Name == "spacebin_token" {
			foundCookie = true
			require.True(t, c.Secure)
		}
	}

	require.True(t, foundCookie)
}

func TestSignInCookieConfigurableAttributes(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true
	cfg.SessionTTLHours = 1
	cfg.SessionCookieSecure = true
	cfg.SessionCookieSameSite = "strict"
	cfg.SessionCookieDomain = "example.com"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturns(database.Account{
		ID:       1,
		Username: "user",
		Password: string(hashedPassword),
	}, nil)
	fakeDB.CreateSessionReturns(nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{
		"username": "user",
		"password": "correct-password",
	})

	start := time.Now()
	req := httptest.NewRequest(http.MethodPost, "/api/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)

	duration := time.Duration(cfg.SessionTTLHours) * time.Hour
	minExpiry := start.Add(duration - time.Second)
	maxExpiry := start.Add(duration + time.Second)
	expectedMaxAge := int(duration.Seconds())

	foundCookie := false
	for _, c := range res.Result().Cookies() {
		if c.Name == "spacebin_token" {
			foundCookie = true
			require.True(t, c.Secure)
			require.Equal(t, http.SameSiteStrictMode, c.SameSite)
			require.Equal(t, expectedMaxAge, c.MaxAge)
			require.True(t, c.Expires.After(minExpiry) && c.Expires.Before(maxExpiry))
			require.Equal(t, "example.com", c.Domain)
			require.Equal(t, "/", c.Path)
		}
	}

	require.True(t, foundCookie)
}

func TestSignUpRedirectsToSignIn(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturnsOnCall(0, database.Account{}, sql.ErrNoRows)
	fakeDB.GetAccountByUsernameReturnsOnCall(1, database.Account{
		ID:       1,
		Username: "newuser",
	}, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	writer.WriteField("username", "newuser")
	writer.WriteField("password", "strongpassword")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/signup", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusSeeOther, res.Result().StatusCode)
	require.Equal(t, "/signin", res.Result().Header.Get("Location"))
}

func TestLogoutClearsCookieAndDeletesSessionJSON(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    userToken,
		Secret:   serverToken,
		Username: "tester",
	}, nil)

	var deleted bool
	fakeDB.DeleteSessionStub = func(ctx context.Context, public string) error {
		deleted = true
		require.Equal(t, "publicKey", public)
		return nil
	}

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodPost, "/api/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNoContent, res.Result().StatusCode)
	require.True(t, deleted)

	foundCookie := false
	for _, c := range res.Result().Cookies() {
		if c.Name == "spacebin_token" {
			foundCookie = true
			require.Equal(t, "", c.Value)
			require.True(t, c.Expires.Before(time.Now().Add(time.Second)))
		}
	}

	require.True(t, foundCookie)
}

func TestLogoutRedirectsAndClearsCookie(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    userToken,
		Secret:   serverToken,
		Username: "tester",
	}, nil)

	fakeDB.DeleteSessionReturns(nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusSeeOther, res.Result().StatusCode)
	require.Equal(t, "/", res.Result().Header.Get("Location"))

	foundCookie := false
	for _, c := range res.Result().Cookies() {
		if c.Name == "spacebin_token" {
			foundCookie = true
			require.Equal(t, "", c.Value)
			require.True(t, c.Expires.Before(time.Now().Add(time.Second)))
		}
	}

	require.True(t, foundCookie)
	require.Equal(t, 1, fakeDB.DeleteSessionCallCount())
}

func TestLogoutInvalidTokenHandledGracefully(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	fakeDB := &databasefakes.FakeDatabase{}

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodPost, "/api/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: "invalid-token"})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNoContent, res.Result().StatusCode)
	require.Equal(t, 0, fakeDB.DeleteSessionCallCount())

	foundCookie := false
	for _, c := range res.Result().Cookies() {
		if c.Name == "spacebin_token" {
			foundCookie = true
			require.Equal(t, "", c.Value)
			require.True(t, c.Expires.Before(time.Now().Add(time.Second)))
		}
	}

	require.True(t, foundCookie)
}

// --- StaticSignUp ---

func TestStaticSignUpAccountsEnabled(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/signup", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "sign")
	require.NotContains(t, res.Body.String(), "{{")
}

func TestStaticSignUpAccountsDisabled(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = false

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/signup", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNotFound, res.Result().StatusCode)
}

// --- StaticSignIn ---

func TestStaticSignInAccountsEnabled(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/signin", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "sign")
	require.NotContains(t, res.Body.String(), "{{")
}

func TestStaticSignInAccountsDisabled(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = false

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/signin", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNotFound, res.Result().StatusCode)
}

// --- StaticSettingsPage ---

func TestStaticSettingsPageAccountsEnabled(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    userToken,
		Secret:   serverToken,
		Username: "tester",
	}, nil)
	fakeDB.GetDocumentsByUsernameReturns(nil, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/account", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.NotContains(t, res.Body.String(), "{{")
}

func TestStaticSettingsPageAccountsDisabled(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = false

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/account", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNotFound, res.Result().StatusCode)
}

// --- getTokenFromRequest bearer path ---

func TestStaticIndexAuthenticatedWithBearerToken(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    userToken,
		Secret:   serverToken,
		Username: "tester",
	}, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountStatic()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "tester")
}

// --- wantsJSONResponse via Accept header ---

func TestLogoutWithAcceptJSONHeader(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    userToken,
		Secret:   serverToken,
		Username: "tester",
	}, nil)
	fakeDB.DeleteSessionReturns(nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	// No Content-Type: application/json; only Accept header triggers wantsJSONResponse via Accept path
	req, _ := http.NewRequest(http.MethodPost, "/api/logout", nil)
	req.Header.Set("Accept", "application/json")
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNoContent, res.Result().StatusCode)
}

// --- parseSameSite default/unknown mode ---

func TestSignInCookieSameSiteUnknownMode(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true
	cfg.SessionCookieSameSite = "none" // unknown → falls to default → SameSiteLaxMode

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturns(database.Account{
		ID:       1,
		Username: "user",
		Password: string(hashedPassword),
	}, nil)
	fakeDB.CreateSessionReturns(nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{
		"username": "user",
		"password": "correct-password",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)

	for _, c := range res.Result().Cookies() {
		if c.Name == "spacebin_token" {
			require.Equal(t, http.SameSiteLaxMode, c.SameSite)
		}
	}
}

// --- handleLogout disabled ---

func TestLogoutDisabled(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = false

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodPost, "/api/logout", nil)
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNotFound, res.Result().StatusCode)
}

// --- invalidateSession / authenticatedUsername database error ---

func TestLogoutSessionDatabaseError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, _ := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{}, fmt.Errorf("database error"))

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodPost, "/api/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func TestStaticIndexSessionDatabaseError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, _ := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{}, fmt.Errorf("database error"))

	s := server.NewServer(&cfg, fakeDB)
	s.MountStatic()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

// --- SignUp missing branches ---

func TestSignUpInvalidBody(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodPost, "/api/signup", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
}

func TestSignUpEmptyCredentials(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{"username": "  ", "password": ""})
	req, _ := http.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
}

func TestSignUpShortPassword(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{"username": "user", "password": "short"})
	req, _ := http.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
}

func TestSignUpGetAccountDatabaseError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturns(database.Account{}, fmt.Errorf("database error"))

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{"username": "user", "password": "strongpassword"})
	req, _ := http.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func TestSignUpCreateAccountDatabaseError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturnsOnCall(0, database.Account{}, sql.ErrNoRows)
	fakeDB.CreateAccountReturns(fmt.Errorf("database error"))

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{"username": "newuser", "password": "strongpassword"})
	req, _ := http.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func TestSignUpSecondGetAccountDatabaseError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturnsOnCall(0, database.Account{}, sql.ErrNoRows)
	fakeDB.CreateAccountReturns(nil)
	fakeDB.GetAccountByUsernameReturnsOnCall(1, database.Account{}, fmt.Errorf("database error"))

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{"username": "newuser", "password": "strongpassword"})
	req, _ := http.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

// --- SignIn missing branches ---

func TestSignInInvalidBody(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodPost, "/api/signin", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
}

func TestSignInEmptyCredentials(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{"username": "", "password": ""})
	req, _ := http.NewRequest(http.MethodPost, "/api/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusBadRequest, res.Result().StatusCode)
}

func TestSignInGetAccountDatabaseError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturns(database.Account{}, fmt.Errorf("database error"))

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{"username": "user", "password": "strongpassword"})
	req, _ := http.NewRequest(http.MethodPost, "/api/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func TestSignInCreateSessionError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturns(database.Account{
		ID:       1,
		Username: "user",
		Password: string(hashedPassword),
	}, nil)
	fakeDB.CreateSessionReturns(fmt.Errorf("database error"))

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{
		"username": "user",
		"password": "correct-password",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

// --- authenticatedUsername: empty session username ---

func TestStaticIndexEmptySessionUsername(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    userToken,
		Secret:   serverToken,
		Username: "", // empty — should be treated as unauthenticated
	}, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountStatic()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	// Page renders but shows no username (unauthenticated view)
	require.NotContains(t, res.Body.String(), "{{")
}

// --- SignIn X-Forwarded-Proto HTTPS ---

func TestSignInCookieSecureWithForwardedProto(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetAccountByUsernameReturns(database.Account{
		ID:       1,
		Username: "user",
		Password: string(hashedPassword),
	}, nil)
	fakeDB.CreateSessionReturns(nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	body, _ := json.Marshal(map[string]string{
		"username": "user",
		"password": "correct-password",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/signin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-Proto", "https")

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)

	for _, c := range res.Result().Cookies() {
		if c.Name == "spacebin_token" {
			require.True(t, c.Secure)
		}
	}
}

// --- Token validation mismatch (covers the if-mismatched branch in authenticatedUsername/invalidateSession) ---

func TestStaticIndexTokenValidationMismatch(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken := util.MakeToken(util.Token{
		Version: "v1",
		Public:  "publicKey",
		Secret:  base64.URLEncoding.EncodeToString([]byte("secret")),
		Salt:    "salt",
	})
	// serverToken has a mismatched Public, so validation fails
	wrongServerToken := util.MakeToken(util.Token{
		Version: "v1",
		Public:  "wrongPublic",
		Secret:  "somesecret",
		Salt:    "salt",
	})

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    userToken,
		Secret:   wrongServerToken,
		Username: "tester",
	}, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountStatic()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.NotContains(t, res.Body.String(), "tester") // mismatch → no username rendered
}

func TestLogoutTokenValidationMismatch(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken := util.MakeToken(util.Token{
		Version: "v1",
		Public:  "publicKey",
		Secret:  base64.URLEncoding.EncodeToString([]byte("secret")),
		Salt:    "salt",
	})
	wrongServerToken := util.MakeToken(util.Token{
		Version: "v1",
		Public:  "wrongPublic",
		Secret:  "somesecret",
		Salt:    "salt",
	})

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    userToken,
		Secret:   wrongServerToken,
		Username: "tester",
	}, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodPost, "/api/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNoContent, res.Result().StatusCode)
	require.Equal(t, 0, fakeDB.DeleteSessionCallCount()) // no delete because validation failed
}

// --- base64 decode error in authenticatedUsername ---

func TestStaticIndexInvalidBase64Secret(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	// Token with invalid base64 in the Secret field
	invalidToken := util.MakeToken(util.Token{
		Version: "v1",
		Public:  "publicKey",
		Secret:  "not-valid-base64!!!!",
		Salt:    "salt",
	})
	// Session exists but the client's Secret can't be base64-decoded
	serverToken := util.MakeToken(util.Token{
		Version: "v1",
		Public:  "publicKey",
		Secret:  "somehash",
		Salt:    "salt",
	})

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    invalidToken,
		Secret:   serverToken,
		Username: "tester",
	}, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountStatic()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: invalidToken})

	res := executeRequest(req, s)

	// base64 error → authenticatedUsername returns "", nil → renders unauthenticated
	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.NotContains(t, res.Body.String(), "tester")
}

// --- Invalid server-side token format in session (ParseToken fails on session.Secret) ---

func TestStaticIndexInvalidServerTokenFormat(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken := util.MakeToken(util.Token{
		Version: "v1",
		Public:  "publicKey",
		Secret:  base64.URLEncoding.EncodeToString([]byte("secret")),
		Salt:    "salt",
	})

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    userToken,
		Secret:   "toofewparts", // ParseToken requires at least 3 dot-separated parts
		Username: "tester",
	}, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountStatic()

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.NotContains(t, res.Body.String(), "tester")
}

func TestLogoutInvalidServerTokenFormat(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken := util.MakeToken(util.Token{
		Version: "v1",
		Public:  "publicKey",
		Secret:  base64.URLEncoding.EncodeToString([]byte("secret")),
		Salt:    "salt",
	})

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    userToken,
		Secret:   "toofewparts",
		Username: "tester",
	}, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodPost, "/api/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusNoContent, res.Result().StatusCode)
	require.Equal(t, 0, fakeDB.DeleteSessionCallCount())
}

// --- StaticSignUp/SignIn/SettingsPage session error (covers authenticatedUsername error path) ---

func TestStaticSignUpSessionError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, _ := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{}, fmt.Errorf("database error"))

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/signup", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func TestStaticSignInSessionError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, _ := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{}, fmt.Errorf("database error"))

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/signin", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func TestStaticSettingsPageSessionError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, _ := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{}, fmt.Errorf("database error"))

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/account", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

// TestStaticSettingsPageUnauthenticated tests that unauthenticated users are redirected
func TestStaticSettingsPageUnauthenticated(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/account", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusSeeOther, res.Result().StatusCode)
	require.Equal(t, "/signin", res.Result().Header.Get("Location"))
}

// TestStaticSettingsPageWithDocuments tests that the account page lists documents
func TestStaticSettingsPageWithDocuments(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    userToken,
		Secret:   serverToken,
		Username: "tester",
	}, nil)
	fakeDB.GetDocumentsByUsernameReturns([]database.Document{
		{ID: "abcdefgh", Content: "first document", Username: "tester"},
		{ID: "ijklmnop", Content: "second document", Username: "tester"},
	}, nil)

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/account", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Contains(t, res.Body.String(), "abcdefgh")
	require.Contains(t, res.Body.String(), "ijklmnop")
	require.NotContains(t, res.Body.String(), "{{")
}

// TestStaticSettingsPageGetDocumentsError tests error from GetDocumentsByUsername
func TestStaticSettingsPageGetDocumentsError(t *testing.T) {
	cfg := mockConfig
	cfg.AccountsEnabled = true

	userToken, serverToken := buildSessionTokens(t, "secret", "salt", "publicKey")

	fakeDB := &databasefakes.FakeDatabase{}
	fakeDB.GetSessionReturns(database.Session{
		Public:   "publicKey",
		Token:    userToken,
		Secret:   serverToken,
		Username: "tester",
	}, nil)
	fakeDB.GetDocumentsByUsernameReturns(nil, fmt.Errorf("database error"))

	s := server.NewServer(&cfg, fakeDB)
	s.MountHandlers()

	req, _ := http.NewRequest(http.MethodGet, "/account", nil)
	req.AddCookie(&http.Cookie{Name: "spacebin_token", Value: userToken})
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusInternalServerError, res.Result().StatusCode)
}

func buildSessionTokens(t *testing.T, secret string, salt string, public string) (string, string) {
	t.Helper()

	userToken := util.MakeToken(util.Token{
		Version: "v1",
		Public:  public,
		Secret:  base64.URLEncoding.EncodeToString([]byte(secret)),
		Salt:    salt,
	})

	hashed := make([]byte, 64)
	sha3.ShakeSum256(hashed, []byte(secret+salt))
	serverToken := util.MakeToken(util.Token{
		Version: "v1",
		Public:  public,
		Secret:  fmt.Sprintf("%x", hashed),
		Salt:    salt,
	})

	return userToken, serverToken
}
