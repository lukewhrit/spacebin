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
	"net/http"
	"os"
	"testing"

	"github.com/orca-group/spirit/internal/database"
	"github.com/orca-group/spirit/internal/server"
	"github.com/stretchr/testify/require"
)

func TestMountStatic(t *testing.T) {
	// Create server and mount expected static files
	s := server.NewServer(&mockConfig, &database.MockDatabase{})

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

	monokaiCssRequest, _ := http.NewRequest(http.MethodGet, "/static/monokai.min.css", nil)
	monokaiCssResponse := executeRequest(monokaiCssRequest, s)
	checkResponseCode(t, http.StatusOK, monokaiCssResponse.Result().StatusCode)

	normalizeCssRequest, _ := http.NewRequest(http.MethodGet, "/static/normalize.css", nil)
	normalizeCssResponse := executeRequest(normalizeCssRequest, s)
	checkResponseCode(t, http.StatusOK, normalizeCssResponse.Result().StatusCode)

	// Check presence of JS files
	appJsRequest, _ := http.NewRequest(http.MethodGet, "/static/app.js", nil)
	appJsResponse := executeRequest(appJsRequest, s)
	checkResponseCode(t, http.StatusOK, appJsResponse.Result().StatusCode)

	highlightJsRequest, _ := http.NewRequest(http.MethodGet, "/static/highlight.min.js", nil)
	highlightJsResponse := executeRequest(highlightJsRequest, s)
	checkResponseCode(t, http.StatusOK, highlightJsResponse.Result().StatusCode)

	// Check presence of image files (logo.svg, favicon.ico)
	faviconRequest, _ := http.NewRequest(http.MethodGet, "/static/favicon.ico", nil)
	faviconResponse := executeRequest(faviconRequest, s)
	checkResponseCode(t, http.StatusOK, faviconResponse.Result().StatusCode)

	logoRequest, _ := http.NewRequest(http.MethodGet, "/static/logo.svg", nil)
	logoResponse := executeRequest(logoRequest, s)
	checkResponseCode(t, http.StatusOK, logoResponse.Result().StatusCode)

	// Check index file exists and returns the correct content
	indexRequest, _ := http.NewRequest(http.MethodGet, "/", nil)
	indexResponse := executeRequest(indexRequest, s)

	checkResponseCode(t, http.StatusOK, indexResponse.Result().StatusCode)

	indexFile, _ := os.ReadFile("./web/index.html")

	require.Equal(t, string(indexFile), indexResponse.Body.String())
}

func TestRegisterHeaders(t *testing.T) {
	s := server.NewServer(&mockConfig, &database.MockDatabase{})

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
	require.Equal(t, "default-src 'self'; frame-ancestors 'none'; base-uri 'none'; form-action 'self'; script-src 'self' 'unsafe-inline';", res.Result().Header.Get("Content-Security-Policy"))
}
