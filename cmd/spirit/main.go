/*
 * Copyright 2020-2022 Luke Whrit, Jack Dorland

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
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/orca-group/spirit/internal/config"
	"github.com/orca-group/spirit/internal/database"
	"github.com/orca-group/spirit/internal/server"
	"github.com/robfig/cron"
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

	// Start server and initialize database
	if err := database.Init(); err != nil {
		log.Fatal().
			Err(err).
			Msg("Could not connect to database")
	}

	// Start expire document cron job
	c := cron.New()

	c.AddFunc("@every 3hr", expirationJob)
}

func main() {
	defer (func() {
		err := http.ListenAndServe(
			fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port),
			server.Router(),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not start HTTP server")
			return
		}
	})()

	log.Info().
		Str("host", config.Config.Host).
		Int("port", config.Config.Port).
		Msg("Successfully started HTTP server")
}

func expirationJob() {
	model := database.Connection
	row, err := model.Rows()

	if err != nil {
		panic(err)
	}

	for row.Next() {
		document := models.Document{}
		database.DBConn.ScanRows(row, &document)

		if time.Now().Unix()-document.CreatedAt >= config.Config.ExpirationAge {
			database.DBConn.Delete(document)
		}

		continue
	}
}
