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

package database

import (
	"embed"
	"errors"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	migratedb "github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/postgres/*.sql migrations/mysql/*.sql migrations/sqlite/*.sql
var migrationFS embed.FS

// migrateUp applies all migrations in dir that haven't been applied to the database behind driver yet.
func migrateUp(fsys fs.FS, dir string, driver migratedb.Driver) error {
	source, err := iofs.New(fsys, dir)

	if err != nil {
		return err
	}

	// NewWithInstance only stores its arguments, so it never returns an error
	migrator, _ := migrate.NewWithInstance("iofs", source, "", driver)

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
