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

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lukewhrit/spacebin/internal/config"
	"github.com/lukewhrit/spacebin/internal/database"
	"github.com/lukewhrit/spacebin/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// Setup zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Load config
	if err := config.Load(); err != nil {
		log.Fatal().
			Err(err).
			Msg("Could not load config")
	}
}

func main() {
	var db database.Database

	u, err := url.Parse(config.Config.ConnectionURI)
	if err != nil {
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("not a walid ConnectionURI")
		}
	}

	switch u.Scheme {
	case "file":
		sq, err := database.NewSqlite(u.Host)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Could not connect to database")
		}
		db = sq
	case "postgresql":
		pg, err := database.NewPostgres()
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Could not connect to database")
		}
		db = pg
	}

	if err := db.Migrate(context.Background()); err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed migrations; Could not create DOCUMENTS tables.")
	}

	m := server.NewServer(&config.Config, db)

	m.MountMiddleware()
	m.RegisterHeaders()

	if !config.Config.Headless {
		m.MountStatic()
	}

	m.MountHandlers()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port),
		Handler: m.Router,
	}

	srvCtx, srvStopCtx := context.WithCancel(context.Background())

	// Watch for OS signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(srvCtx, 30*time.Second)
		defer shutdownCtxCancel() // release srvCtx if we take too long to shut down

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Warn().Msg("Graceful shutdown timed out... forcing regular exit.")
			}
		}()

		// Gracefully shut down services
		log.Info().Msg("Killing services")

		// Web server
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed shutting HTTP listener down")
		}

		// Database
		err := db.Close()

		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed closing database connection")
		}

		srvStopCtx()
	}()

	log.Info().
		Str("host", config.Config.Host).
		Int("port", config.Config.Port).
		Msg("Starting HTTP listener")

	// Start the server
	err = srv.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.Fatal().
			Err(err).
			Msg("Failed to start HTTP listener")
	}

	<-srvCtx.Done()
	log.Info().Msg("Successfully and cleanly shut down all Spirit services")
}
