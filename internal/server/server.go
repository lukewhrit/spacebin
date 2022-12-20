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
)

// Start initializes the server
func Router() *chi.Mux {
	// Create Mux
	r := chi.NewRouter()

	// Register middleware
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Recoverer)

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
