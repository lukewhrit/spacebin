/*
 * Copyright 2020-2023 Luke Whritenour, Jack Dorland

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

package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/orca-group/spirit/internal/database"
	"github.com/orca-group/spirit/internal/util"
	"golang.org/x/exp/slices"
)

func getDocument(s *Server, w http.ResponseWriter, id string) database.Document {
	// Retrieve document from the database
	document, err := database.FindDocument(s.Database, id)

	if err != nil {
		// If the document is not found (ErrNoRows), return the error with a 404
		if errors.Is(err, sql.ErrNoRows) {
			util.WriteError(w, http.StatusNotFound, err)
			return document
		}

		// Otherwise, return the error with a 500
		util.WriteError(w, http.StatusInternalServerError, err)
		return document
	}

	return document
}

func (s *Server) StaticDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "document")

	// Validate document ID
	if len(id) != s.Config.IDLength && !slices.Contains(s.Config.Documents, id) {
		err := fmt.Errorf("id is of length %d, should be %d", len(id), s.Config.IDLength)
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Retrieve document from the database
	document := getDocument(s, w, id)

	t, err := template.ParseFS(resources, "web/document.html")

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	data := map[string]interface{}{
		"Lines":   util.CountLines(document.Content),
		"Content": document.Content,
	}

	if err := t.Execute(w, data); err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) FetchDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "document")

	// Validate document ID
	if len(id) != s.Config.IDLength && !slices.Contains(s.Config.Documents, id) {
		err := fmt.Errorf("id is of length %d, should be %d", len(id), s.Config.IDLength)
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	document := getDocument(s, w, id)

	// Try responding with the document and a 200, or write an error if that fails
	if err := util.WriteJSON(w, http.StatusOK, document); err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) FetchRawDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "document")

	// Validate document ID
	if len(id) != s.Config.IDLength && !slices.Contains(s.Config.Documents, id) {
		err := fmt.Errorf("id is of length %d, should be %d", len(id), s.Config.IDLength)
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	document := getDocument(s, w, id)

	// Respond with only the documents content
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(document.Content))
}
