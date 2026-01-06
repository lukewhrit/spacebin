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
