package server

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/lukewhrit/spacebin/internal/config"
	"github.com/lukewhrit/spacebin/internal/util"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

const sessionCookieName = "spacebin_token"

func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	if !s.Config.AccountsEnabled {
		util.WriteError(w, http.StatusNotFound, errors.New("accounts disabled"))
		return
	}

	// Parse body from HTML request
	body, err := util.HandleSignupBody(s.Config.MaxSize, r)

	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	body.Username = strings.TrimSpace(body.Username)
	if body.Username == "" || body.Password == "" {
		util.WriteError(w, http.StatusBadRequest, errors.New("username and password are required"))
		return
	}

	if len(body.Password) < 8 {
		util.WriteError(w, http.StatusBadRequest, errors.New("password must be at least 8 characters long"))
		return
	}

	// Make sure username does not exist
	_, err = s.Database.GetAccountByUsername(r.Context(), body.Username)

	if err == nil {
		util.WriteError(w, http.StatusConflict, errors.New("username already exists"))
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Create account
	// Encryption handled in Database function
	err = s.Database.CreateAccount(r.Context(), body.Username, body.Password)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Respond on success with account ID and username
	account, err := s.Database.GetAccountByUsername(r.Context(), body.Username)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":       account.ID,
		"username": account.Username,
	})

}

func (s *Server) StaticSignUp(w http.ResponseWriter, r *http.Request) {
	if !s.Config.AccountsEnabled {
		util.WriteError(w, http.StatusNotFound, errors.New("accounts disabled"))
		return
	}

	t, err := template.ParseFS(resources, "web/signup.html")

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	username, err := s.authenticatedUsername(r)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = t.Execute(w, map[string]interface{}{
		"Analytics":     config.Config.Analytics,
		"Authenticated": username != "",
		"Username":      username,
	})

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) SignIn(w http.ResponseWriter, r *http.Request) {
	if !s.Config.AccountsEnabled {
		util.WriteError(w, http.StatusNotFound, errors.New("accounts disabled"))
		return
	}

	// Parse body from HTML request
	body, err := util.HandleSigninBody(s.Config.MaxSize, r)

	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	body.Username = strings.TrimSpace(body.Username)
	if body.Username == "" || body.Password == "" {
		util.WriteError(w, http.StatusBadRequest, errors.New("username and password are required"))
		return
	}

	// Get user from database
	acc, err := s.Database.GetAccountByUsername(r.Context(), body.Username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.WriteError(w, http.StatusUnauthorized, errors.New("invalid username or password"))
			return
		}
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

		// Compare passwords
		if bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(body.Password)) == nil {
			// Generate public, secret keys and salt
			pub, sec, salt, err := util.GenerateStrings([]int{64, 64, 32})

		if err != nil {
			log.Fatal(err)
		}

		// Salt secret key
		buf := []byte(sec + salt)
		secret := make([]byte, 64)
		sha3.ShakeSum256(secret, buf)

		// Create user and server tokens for later comparison
			userToken := util.MakeToken(util.Token{
				Version: "v1",
				Public:  pub,
				Secret:  base64.URLEncoding.EncodeToString([]byte(sec)),
				Salt:    salt,
			})

			serverToken := util.MakeToken(util.Token{
				Version: "v1",
				Public:  pub,
				Secret:  fmt.Sprintf("%x", secret),
				Salt:    salt,
			})

			// Add session to Postgres
			if err := s.Database.CreateSession(r.Context(), pub, userToken, serverToken, acc.Username); err != nil {
				util.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     sessionCookieName,
				Value:    userToken,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})

			util.WriteJSON(w, http.StatusOK, map[string]string{
				"token": userToken,
				"user":  acc.Username,
			})
		} else {
			util.WriteError(w, http.StatusUnauthorized, errors.New("invalid username or password"))
			return
		}
	}

func (s *Server) StaticSignIn(w http.ResponseWriter, r *http.Request) {
	if !s.Config.AccountsEnabled {
		util.WriteError(w, http.StatusNotFound, errors.New("accounts disabled"))
		return
	}

	t, err := template.ParseFS(resources, "web/signin.html")

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	username, err := s.authenticatedUsername(r)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = t.Execute(w, map[string]interface{}{
		"Analytics":     config.Config.Analytics,
		"Authenticated": username != "",
		"Username":      username,
	})

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) authenticatedUsername(r *http.Request) (string, error) {
	if !s.Config.AccountsEnabled {
		return "", nil
	}

	token := getTokenFromRequest(r)
	if token == "" {
		return "", nil
	}

	clientToken, err := util.ParseToken(token)
	if err != nil {
		return "", nil
	}

	session, err := s.Database.GetSession(r.Context(), clientToken.Public)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	serverToken, err := util.ParseToken(session.Secret)
	if err != nil {
		return "", nil
	}

	secretBytes, err := base64.URLEncoding.DecodeString(clientToken.Secret)
	if err != nil {
		return "", nil
	}

	secret := make([]byte, 64)
	sha3.ShakeSum256(secret, append(secretBytes, []byte(clientToken.Salt)...))
	expected := fmt.Sprintf("%x", secret)

	if clientToken.Public != serverToken.Public || clientToken.Salt != serverToken.Salt || expected != serverToken.Secret {
		return "", nil
	}

	if session.Username == "" {
		return "", nil
	}

	return session.Username, nil
}

func getTokenFromRequest(r *http.Request) string {
	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		return cookie.Value
	}

	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return strings.TrimSpace(authHeader[7:])
	}

	return ""
}
