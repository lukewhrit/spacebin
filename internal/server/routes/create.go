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

package routes

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/orca-group/spirit/internal/config"
	"github.com/orca-group/spirit/internal/database"
	"github.com/orca-group/spirit/internal/database/models"
	"github.com/orca-group/spirit/internal/util"
)

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateId() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, config.Config.Documents.IDLength)

	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(b)
}

func CreateDocument(w http.ResponseWriter, r *http.Request) {
	// Parse body from HTML request
	body, err := util.HandleBody(r)

	if err != nil {
		util.WriteError(err, w, http.StatusBadRequest)
	}

	// Validate fields of body
	if err := util.ValidateBody(body); err != nil {
		util.WriteError(err, w, http.StatusBadRequest)
	}

	// Generate ID and append it to content and extension as a Document object
	id := generateId()
	doc := models.Document{
		ID:        id,
		Content:   body.Content,
		Extension: body.Extension,
	}

	// Add Document object to database
	res := database.DBConn.Create(&doc)

	if res.Error != nil {
		util.WriteError(res.Error, w, http.StatusInternalServerError)
	}

	// Respond to request with Document object
	payload := util.Payload{
		ID:          doc.ID,
		Content:     doc.Content,
		ContentHash: "",
		Extension:   doc.Extension,
		UpdatedAt:   doc.UpdatedAt,
		CreatedAt:   doc.CreatedAt,
	}

	payload.WriteJSON(w, http.StatusOK)
}
