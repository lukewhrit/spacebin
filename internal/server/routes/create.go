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
	"net/http"

	"github.com/orca-group/spirit/internal/config"
	"github.com/orca-group/spirit/internal/database"
	"github.com/orca-group/spirit/internal/database/models"
	"github.com/orca-group/spirit/internal/util"
)

func CreateDocument(w http.ResponseWriter, r *http.Request) {
	// Parse body from HTML request
	body, err := util.HandleBody(r)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Validate fields of body
	if err := util.ValidateBody(body); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Generate ID and create document with ID and content
	doc := models.Document{
		ID:      util.GenerateID(config.Config.IDType, config.Config.IDLength),
		Content: body.Content,
	}

	// Add Document object to database
	res := database.DBConn.Create(&doc)

	if res.Error != nil {
		util.WriteError(w, http.StatusInternalServerError, res.Error)
		return
	}

	// Respond to request with Document object
	if err := util.WriteJSON(w, http.StatusOK, util.DocumentResponse{
		ID:        doc.ID,
		Content:   doc.Content,
		UpdatedAt: doc.UpdatedAt,
		CreatedAt: doc.CreatedAt,
	}); err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}
