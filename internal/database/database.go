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

package database

import (
	"github.com/orca-group/spirit/internal/config"
	"github.com/orca-group/spirit/internal/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConn holds the current connection to the database
var DBConn *gorm.DB

// Init opens a connection to the database
func Init() error {
	DBConn, err := gorm.Open(postgres.Open(config.Config.ConnectionURI),
		&gorm.Config{})

	DBConn.AutoMigrate(&models.Document{})

	return err
}
