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
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/orca-group/spirit/internal/config"
	"github.com/orca-group/spirit/internal/server/routes"
	"github.com/orca-group/spirit/internal/util"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Router *chi.Mux
}

func NewServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	return s
}

// These functions should be executed in the order they are defined, that is:
//  1. Mount middleware - MountMiddleware()
//  2. Add security headers - RegisterHeaders()
//  3. Mount actual routes - MountHandlers()

func (s *Server) MountMiddleware() {
	// Register middleware
	s.Router.Use(util.Logger)
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.AllowContentType("application/json", "multipart/form-data"))

	// Ratelimiter
	reqs, per, err := util.ParseRatelimiterString(config.Config.Ratelimiter)

	if err != nil {
		log.Error().
			Err(err).
			Msg("Parse Ratelimiter Error")
	}

	s.Router.Use(httprate.LimitAll(reqs, per))
	s.Router.Use(middleware.Heartbeat("/ping"))
	s.Router.Use(middleware.Recoverer)

	// CORS
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}

func (s *Server) RegisterHeaders() {
	s.Router.Use(middleware.SetHeader("X-Download-Options", "noopen"))
	s.Router.Use(middleware.SetHeader("X-DNS-Prefetch-Control", "off"))
	s.Router.Use(middleware.SetHeader("X-Frame-Options", "SAMEORIGIN"))
	s.Router.Use(middleware.SetHeader("X-XSS-Protection", "1; mode=block"))
	s.Router.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	s.Router.Use(middleware.SetHeader("Referrer-Policy", "no-referrer-when-downgrade"))
	s.Router.Use(middleware.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload"))
	s.Router.Use(middleware.SetHeader("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'none';"))
}

func (s *Server) MountHandlers() {
	// Register routes
	s.Router.Get("/config", routes.Config)

	s.Router.Post("/", routes.CreateDocument)
	s.Router.Get("/{document}", routes.FetchDocument)
	s.Router.Get("/{document}/raw", routes.FetchRawDocument)

	// Legacy routes
	s.Router.Post("/v1/documents/", routes.CreateDocument)
	s.Router.Get("/v1/documents/{document}", routes.FetchDocument)
	s.Router.Get("/v1/documents/{document}/raw", routes.FetchRawDocument)
}
