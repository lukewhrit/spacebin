/*
 * Copyright 2020-2022 Luke Whritenour, Jack Dorland

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

package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/orca-group/spirit/internal/server/routes"
	"github.com/orca-group/spirit/internal/util"
)

// Start initializes the server
func Router() *chi.Mux {
	// Create Mux
	r := chi.NewRouter()

	// Register middleware
	r.Use(util.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.AllowContentType("application/json", "multipart/form-data"))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Recoverer)

	// Headers
	r.Use(middleware.SetHeader("X-Download-Options", "noopen"))
	r.Use(middleware.SetHeader("X-DNS-Prefetch-Control", "off"))
	r.Use(middleware.SetHeader("X-Frame-Options", "SAMEORIGIN"))
	r.Use(middleware.SetHeader("X-XSS-Protection", "1; mode=block"))
	r.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	r.Use(middleware.SetHeader("Referrer-Policy", "no-referrer-when-downgrade"))
	r.Use(middleware.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload"))
	r.Use(middleware.SetHeader("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'none';"))

	// Register routes
	r.Get("/config", routes.Config)

	r.Post("/", routes.CreateDocument)
	r.Get("/{document}", routes.CreateDocument)
	r.Get("/{document}/raw", routes.FetchRawDocument)

	// Old routes
	r.Post("/v1/documents/", routes.CreateDocument)
	r.Get("/v1/documents/{document}", routes.FetchDocument)
	r.Get("/v1/documents/{document}/raw", routes.FetchRawDocument)

	return r
}
