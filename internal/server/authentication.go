package server

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lukewhrit/spacebin/internal/database"
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

	// bcrypt can't hash anything longer
	if len(body.Password) > 72 {
		util.WriteError(w, http.StatusBadRequest, errors.New("password must be at most 72 bytes long"))
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

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		util.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"id":       account.ID,
			"username": account.Username,
		})
		return
	}

	http.Redirect(w, r, "/signin", http.StatusSeeOther)

}

func (s *Server) StaticSignUp(w http.ResponseWriter, r *http.Request) {
	if !s.Config.AccountsEnabled {
		util.RenderError(&resources, w, http.StatusNotFound, errors.New("accounts disabled"))
		return
	}

	username, err := s.authenticatedUsername(r)
	if err != nil {
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	render(w, "signup.html", s.pageData(username, nil))
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
		pub, sec, salt := util.GenerateStrings()

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

		http.SetCookie(w, s.buildSessionCookie(r, userToken))

		if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			util.WriteJSON(w, http.StatusOK, map[string]string{
				"token": userToken,
				"user":  acc.Username,
			})
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	util.WriteError(w, http.StatusUnauthorized, errors.New("invalid username or password"))
}

func (s *Server) StaticSignIn(w http.ResponseWriter, r *http.Request) {
	if !s.Config.AccountsEnabled {
		util.RenderError(&resources, w, http.StatusNotFound, errors.New("accounts disabled"))
		return
	}

	username, err := s.authenticatedUsername(r)
	if err != nil {
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	render(w, "signin.html", s.pageData(username, nil))
}

func (s *Server) StaticSettingsPage(w http.ResponseWriter, r *http.Request) {
	if !s.Config.AccountsEnabled {
		util.RenderError(&resources, w, http.StatusNotFound, errors.New("accounts disabled"))
		return
	}

	username, err := s.authenticatedUsername(r)
	if err != nil {
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	if username == "" {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	docs, err := s.Database.GetDocumentsByUsername(r.Context(), username)
	if err != nil {
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	render(w, "account.html", s.pageData(username, map[string]any{
		"Documents": docs,
	}))
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	s.handleLogout(w, r, false)
}

func (s *Server) StaticLogout(w http.ResponseWriter, r *http.Request) {
	s.handleLogout(w, r, true)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request, forceRedirect bool) {
	if !s.Config.AccountsEnabled {
		util.WriteError(w, http.StatusNotFound, errors.New("accounts disabled"))
		return
	}

	if err := s.invalidateSession(r); err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	clearSessionCookie(w)

	if forceRedirect || !wantsJSONResponse(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) authenticatedUsername(r *http.Request) (string, error) {
	if !s.Config.AccountsEnabled {
		return "", nil
	}

	session, _, err := s.requestSession(r)

	return session.Username, err
}

// requestSession returns the session belonging to the request's token. ok is false
// if the request has no valid token, in which case session is empty.
func (s *Server) requestSession(r *http.Request) (session database.Session, ok bool, err error) {
	// Also fails for requests without a token
	clientToken, err := util.ParseToken(getTokenFromRequest(r))
	if err != nil {
		return database.Session{}, false, nil
	}

	session, err = s.Database.GetSession(r.Context(), clientToken.Public)
	if errors.Is(err, sql.ErrNoRows) {
		return database.Session{}, false, nil
	}

	if err != nil {
		return database.Session{}, false, err
	}

	serverToken, err := util.ParseToken(session.Secret)
	if err != nil {
		return database.Session{}, false, nil
	}

	secretBytes, err := base64.URLEncoding.DecodeString(clientToken.Secret)
	if err != nil {
		return database.Session{}, false, nil
	}

	secret := make([]byte, 64)
	sha3.ShakeSum256(secret, append(secretBytes, []byte(clientToken.Salt)...))
	expected := fmt.Sprintf("%x", secret)

	if clientToken.Public != serverToken.Public || clientToken.Salt != serverToken.Salt || expected != serverToken.Secret {
		return database.Session{}, false, nil
	}

	return session, true, nil
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

func (s *Server) invalidateSession(r *http.Request) error {
	session, ok, err := s.requestSession(r)
	if !ok {
		return err
	}

	return s.Database.DeleteSession(r.Context(), session.Public)
}

func (s *Server) buildSessionCookie(r *http.Request, token string) *http.Cookie {
	duration := time.Duration(s.Config.SessionTTLHours) * time.Hour
	secure := s.Config.SessionCookieSecure

	if r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https") {
		secure = true
	}

	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(duration.Seconds()),
		Expires:  time.Now().Add(duration),
		HttpOnly: true,
		Secure:   secure,
		SameSite: parseSameSite(s.Config.SessionCookieSameSite),
	}

	if s.Config.SessionCookieDomain != "" {
		cookie.Domain = s.Config.SessionCookieDomain
	}

	return cookie
}

func parseSameSite(mode string) http.SameSite {
	switch strings.ToLower(mode) {
	case "strict":
		return http.SameSiteStrictMode
	case "lax", "default", "":
		return http.SameSiteLaxMode
	default:
		return http.SameSiteLaxMode
	}
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	})
}

func wantsJSONResponse(r *http.Request) bool {
	if strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")), "application/json") {
		return true
	}

	return strings.Contains(strings.ToLower(r.Header.Get("Accept")), "application/json")
}
