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

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lukewhrit/spacebin/internal/config"
	"github.com/lukewhrit/spacebin/internal/database"
	"github.com/lukewhrit/spacebin/internal/server"
	"github.com/lukewhrit/spacebin/internal/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// exit is os.Exit, replaced in tests
var exit = os.Exit

func main() {
	// Health check for Docker
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		exit(healthcheck())
		return
	}

	// Setup zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Shut down gracefully on OS signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	if err := run(ctx); err != nil {
		log.Error().
			Err(err).
			Msg("Spacebin stopped")

		exit(1)
		return
	}

	log.Info().Msg("Successfully and cleanly shut down all Spacebin services")
}

// run starts Spacebin, and blocks until ctx is cancelled and the server has shut down.
func run(ctx context.Context) error {
	if util.EnvContains("SPIRIT_") {
		log.Warn().Msg("SPIRIT_ environment variables are deprecated. Please use SPACEBIN_ instead. (see v2 upgrade notes)")
	}

	// Load config
	if err := config.Load(); err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	db, err := openDatabase(config.Config.ConnectionURI)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	defer db.Close()

	// Perform migrations
	if err := db.Migrate(ctx); err != nil {
		return fmt.Errorf("failed migrations: %w", err)
	}

	// Create a new server and register middleware, security headers, static files, and handlers
	m := server.NewServer(&config.Config, db)

	m.MountMiddleware()
	m.RegisterHeaders()

	if !config.Config.Headless {
		m.MountStatic()
	}

	m.MountHandlers()

	// Delete expired documents in the background, stopping before the database is closed however run returns
	expireCtx, stopExpiring := context.WithCancel(ctx)
	expired := make(chan struct{})

	go func() {
		m.ExpireDocuments(expireCtx, time.Hour)
		close(expired)
	}()

	defer func() {
		stopExpiring()
		<-expired
	}()

	// Create the server on the specified host and port
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port),
		Handler: m.Router,
	}

	log.Info().
		Str("host", config.Config.Host).
		Int("port", config.Config.Port).
		Msg("Starting HTTP listener")

	// Start the server
	listenErr := make(chan error, 1)
	go func() { listenErr <- srv.ListenAndServe() }()

	// ListenAndServe only returns before shutdown if it failed
	select {
	case err := <-listenErr:
		return fmt.Errorf("failed to start HTTP listener: %w", err)
	case <-ctx.Done():
	}

	// Graceful shutdown
	log.Info().Msg("Killing services")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}

// openDatabase connects to the database that uri points to, based on its scheme.
func openDatabase(uri string) (database.Database, error) {
	scheme, _, _ := strings.Cut(uri, "://")

	switch scheme {
	case "file", "sqlite":
		return database.NewSQLite(uri)
	case "postgresql", "postgres":
		return database.NewPostgres(uri)
	case "mysql", "mariadb":
		return database.NewMySQL(uri)
	}

	return nil, fmt.Errorf("unsupported database type %q", scheme)
}

// healthcheck pings the running server and returns an exit code for Docker's HEALTHCHECK.
func healthcheck() int {
	if err := config.Load(); err != nil {
		return 1
	}

	host := config.Config.Host

	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}

	addr := net.JoinHostPort(host, strconv.Itoa(config.Config.Port))
	client := http.Client{Timeout: 2 * time.Second}

	res, err := client.Get("http://" + addr + "/ping")
	if err != nil {
		return 1
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 1
	}

	return 0
}
