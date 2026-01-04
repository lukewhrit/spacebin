package server_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

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

	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
	require.Equal(t, "user", capturedUsername)

	foundCookie := false
	for _, c := range res.Result().Cookies() {
		if c.Name == "spacebin_token" && c.Value != "" {
			foundCookie = true
		}
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
