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

package database

import (
	"log"

	"github.com/coral-dev/spirit/internal/pkg/config"
	"github.com/coral-dev/spirit/internal/pkg/database/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBConn holds the current connection to the database
var DBConn *gorm.DB

// Init opens a connection to the database
func Init() {
	var err error
	var dialect gorm.Dialector

	switch config.Config.Database.Dialect {
	case "sqlite":
		dialect = sqlite.Open(config.Config.Database.ConnectionURI)
	case "postgresql":
		dialect = postgres.Open(config.Config.Database.ConnectionURI)
	case "mysql":
		dialect = mysql.Open(config.Config.Database.ConnectionURI)
	}

	DBConn, err = gorm.Open(dialect, &gorm.Config{})

	DBConn.AutoMigrate(&models.Document{})

	if err != nil {
		log.Fatalf("Failed to connect to database: %e", err)
	}
}
