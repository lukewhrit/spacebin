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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lukewhrit/spacebin/internal/database/databasefakes"
	"github.com/lukewhrit/spacebin/internal/server"
	"github.com/stretchr/testify/require"
)

func TestMountStatic(t *testing.T) {
	// Create server and mount expected static files
	s := server.NewServer(&mockConfig, &databasefakes.FakeDatabase{})

	s.MountStatic()

	// Check robots.txt
	req, _ := http.NewRequest(http.MethodGet, "/robots.txt", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)

	file, _ := os.ReadFile("web/static/robots.txt")
	require.Equal(t, res.Body.String(), string(file))

	// Check presence of CSS files
	globalCssRequest, _ := http.NewRequest(http.MethodGet, "/static/global.css", nil)
	globalCssResponse := executeRequest(globalCssRequest, s)
	checkResponseCode(t, http.StatusOK, globalCssResponse.Result().StatusCode)

	normalizeCssRequest, _ := http.NewRequest(http.MethodGet, "/static/normalize.css", nil)
	normalizeCssResponse := executeRequest(normalizeCssRequest, s)
	checkResponseCode(t, http.StatusOK, normalizeCssResponse.Result().StatusCode)

	// Check presence of JS files
	appJsRequest, _ := http.NewRequest(http.MethodGet, "/static/app.js", nil)
	appJsResponse := executeRequest(appJsRequest, s)
	checkResponseCode(t, http.StatusOK, appJsResponse.Result().StatusCode)

	// Check presence of image files (logo.svg, favicon.ico)
	faviconRequest, _ := http.NewRequest(http.MethodGet, "/static/favicon.ico", nil)
	faviconResponse := executeRequest(faviconRequest, s)
	checkResponseCode(t, http.StatusOK, faviconResponse.Result().StatusCode)

	logoRequest, _ := http.NewRequest(http.MethodGet, "/static/logo.svg", nil)
	logoResponse := executeRequest(logoRequest, s)
	checkResponseCode(t, http.StatusOK, logoResponse.Result().StatusCode)

	// Check index file renders correctly (no unrendered template syntax)
	indexRequest, _ := http.NewRequest(http.MethodGet, "/", nil)
	indexResponse := executeRequest(indexRequest, s)

	checkResponseCode(t, http.StatusOK, indexResponse.Result().StatusCode)
	require.Contains(t, indexResponse.Body.String(), "Spacebin")
	require.Contains(t, indexResponse.Body.String(), "textarea")
	require.NotContains(t, indexResponse.Body.String(), "{{")
}

func TestRegisterHeaders(t *testing.T) {
	s := server.NewServer(&mockConfig, &databasefakes.FakeDatabase{})

	s.RegisterHeaders()
	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("."))
	})

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	res := executeRequest(req, s)

	// Ensure 200
	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)

	require.Equal(t, "noopen", res.Result().Header.Get("X-Download-Options"))
	require.Equal(t, "off", res.Result().Header.Get("X-DNS-Prefetch-Control"))
	require.Equal(t, "SAMEORIGIN", res.Result().Header.Get("X-Frame-Options"))
	require.Equal(t, "1; mode=block", res.Result().Header.Get("X-XSS-Protection"))
	require.Equal(t, "nosniff", res.Result().Header.Get("X-Content-Type-Options"))
	require.Equal(t, "no-referrer-when-downgrade", res.Result().Header.Get("Referrer-Policy"))
	require.Equal(t, "max-age=31536000; includeSubDomains; preload", res.Result().Header.Get("Strict-Transport-Security"))
	require.Equal(t, mockConfig.ContentSecurityPolicy, res.Result().Header.Get("Content-Security-Policy"))
}

// TestMountMiddleware tests mounting middleware on the server
func TestMountMiddleware(t *testing.T) {
	s := server.NewServer(&mockConfig, &databasefakes.FakeDatabase{})

	s.MountMiddleware()
	s.Router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)

	// Test ping heartbeat endpoint
	pingReq, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	pingRes := executeRequest(pingReq, s)
	checkResponseCode(t, http.StatusOK, pingRes.Result().StatusCode)
	require.Equal(t, ".", pingRes.Body.String())
}

// TestMountMiddlewareWithBasicAuth tests middleware with basic auth
func TestMountMiddlewareWithBasicAuth(t *testing.T) {
	authConfig := mockConfig
	authConfig.Username = "testuser"
	authConfig.Password = "testpass"

	s := server.NewServer(&authConfig, &databasefakes.FakeDatabase{})
	s.MountMiddleware()
	s.Router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("authenticated"))
	})

	// Request without auth should fail
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	res := executeRequest(req, s)
	checkResponseCode(t, http.StatusUnauthorized, res.Result().StatusCode)

	// Request with correct auth should succeed
	authReq, _ := http.NewRequest(http.MethodGet, "/test", nil)
	authReq.SetBasicAuth("testuser", "testpass")
	authRes := executeRequest(authReq, s)
	checkResponseCode(t, http.StatusOK, authRes.Result().StatusCode)
	require.Equal(t, "authenticated", authRes.Body.String())
}

// TestMountMiddlewareWithInvalidRatelimiter tests middleware with invalid ratelimiter
func TestMountMiddlewareWithInvalidRatelimiter(t *testing.T) {
	invalidConfig := mockConfig
	invalidConfig.Ratelimiter = "invalid-format"

	s := server.NewServer(&invalidConfig, &databasefakes.FakeDatabase{})
	s.MountMiddleware() // Should log error but not panic
	s.Router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	res := executeRequest(req, s)
	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)
}

// pingFrom requests /ping as remoteAddr, optionally with an X-Forwarded-For header, and returns the status.
func pingFrom(s *server.Server, remoteAddr, forwardedFor string) int {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.RemoteAddr = remoteAddr

	if forwardedFor != "" {
		req.Header.Set("X-Forwarded-For", forwardedFor)
	}

	return executeRequest(req, s).Result().StatusCode
}

func newRatelimitedServer(trustedProxies ...string) *server.Server {
	cfg := mockConfig
	cfg.Ratelimiter = "1x60"
	cfg.TrustedProxies = trustedProxies

	s := server.NewServer(&cfg, &databasefakes.FakeDatabase{})
	s.MountMiddleware()

	// chi only runs middleware once a route is registered
	s.MountHandlers()
	return s
}

func TestRatelimiterIsPerClientIP(t *testing.T) {
	s := newRatelimitedServer()

	require.Equal(t, http.StatusOK, pingFrom(s, "192.0.2.1:1234", ""))
	require.Equal(t, http.StatusTooManyRequests, pingFrom(s, "192.0.2.1:5678", ""))

	// Another client has its own budget
	require.Equal(t, http.StatusOK, pingFrom(s, "192.0.2.2:1234", ""))
}

func TestRatelimiterIgnoresForwardedForWithoutTrustedProxies(t *testing.T) {
	s := newRatelimitedServer()

	require.Equal(t, http.StatusOK, pingFrom(s, "192.0.2.1:1234", "198.51.100.1"))

	// Changing X-Forwarded-For doesn't get a client a new budget
	require.Equal(t, http.StatusTooManyRequests, pingFrom(s, "192.0.2.1:1234", "198.51.100.2"))
}

func TestRatelimiterTrustsForwardedForFromTrustedProxies(t *testing.T) {
	s := newRatelimitedServer("10.0.0.0/8")

	// Clients behind the proxy each get their own budget
	require.Equal(t, http.StatusOK, pingFrom(s, "10.0.0.1:1234", "198.51.100.1"))
	require.Equal(t, http.StatusOK, pingFrom(s, "10.0.0.1:1234", "198.51.100.2"))
	require.Equal(t, http.StatusTooManyRequests, pingFrom(s, "10.0.0.1:1234", "198.51.100.1"))

	// Trusted hops can't be used to pick another client's budget
	require.Equal(t, http.StatusTooManyRequests, pingFrom(s, "10.0.0.1:1234", "198.51.100.1, 10.0.0.2"))
}
