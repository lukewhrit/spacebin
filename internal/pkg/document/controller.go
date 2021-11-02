/*
 * Copyright 2020-2021 Luke Whrit, Jack Dorland

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

package document

import (
	"math/rand"
	"time"

	"github.com/coral-dev/spirit/internal/pkg/config"
	"github.com/coral-dev/spirit/internal/pkg/database"
	"github.com/coral-dev/spirit/internal/pkg/database/models"
	"github.com/robfig/cron/v3"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// CreateID generates a random string of length `length` using the unix timestamp
func CreateID(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, length)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// GetDocument retrieves a document record from the database via `id`
func GetDocument(id string) (*models.Document, error) {
	document := models.Document{}
	err := database.DBConn.Where("id = ?", id).First(&document)

	return &document, err.Error
}

// NewDocument creates a new document record in the database
func NewDocument(content string, extension string) (string, error) {
	id := CreateID(config.Config.Documents.IDLength)

	doc := models.Document{
		ID:        id,
		Content:   content,
		Extension: extension,
	}

	// Create new record in database
	res := database.DBConn.Create(&doc)

	return doc.ID, res.Error
}

// ExpireDocument registers a cron job to delete documents after they get too old
func ExpireDocument() *cron.Cron {
	c := cron.New()

	c.AddFunc("@every 3hr", func() {
		model := database.DBConn.Model(&models.Document{})
		row, err := model.Rows()

		if err != nil {
			panic(err)
		}

		for row.Next() {
			document := models.Document{}
			database.DBConn.ScanRows(row, &document)

			if time.Now().Unix()-document.CreatedAt >= config.Config.Documents.MaxAge {
				database.DBConn.Delete(document)
			}

			continue
		}
	})

	return c
}
