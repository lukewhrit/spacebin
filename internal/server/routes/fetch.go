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
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/orca-group/spirit/internal/config"
	"github.com/orca-group/spirit/internal/database"
	"github.com/orca-group/spirit/internal/database/models"
	"github.com/orca-group/spirit/internal/util"
)

func FetchDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "document")

	if len(id) != config.Config.IDLength {
		err := fmt.Errorf("id is of length %d, should be %d", len(id), config.Config.IDLength)
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	doc := models.Document{}

	if err := database.DBConn.Where("id = ?", id).First(&doc).Error; err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}

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

func FetchRawDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "document")

	if len(id) != config.Config.IDLength {
		err := fmt.Errorf("id is of length %d, should be %d", len(id), config.Config.IDLength)
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	document := models.Document{}

	if err := database.DBConn.Where("id = ?", id).First(&document).Error; err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(document.Content))
}
