/*
 * Copyright 2020-2024 Luke Whritenour

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
	"fmt"
	"net/http"

	"github.com/lukewhrit/spacebin/internal/util"
)

func createDocument(s *Server, w http.ResponseWriter, r *http.Request) string {
	// Parse body from HTML request
	body, err := util.HandleBody(s.Config.MaxSize, r)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return ""
	}

	// Validate fields of body
	if err := util.ValidateBody(s.Config.MaxSize, body); err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return ""
	}

	// Add Document object to database
	id := util.GenerateID(s.Config.IDType, s.Config.IDLength)

	if err := s.Database.CreateDocument(
		r.Context(),
		id,
		body.Content,
	); err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return ""
	}

	return id
}

func (s *Server) CreateDocument(w http.ResponseWriter, r *http.Request) {
	// Create document, then pull it from the database
	id := createDocument(s, w, r)
	document, err := s.Database.GetDocument(r.Context(), id)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Respond to request with Document object
	if err := util.WriteJSON(w, http.StatusOK, document); err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) StaticCreateDocument(w http.ResponseWriter, r *http.Request) {
	// Create document, then pull it from the database
	id := createDocument(s, w, r)
	document, err := s.Database.GetDocument(r.Context(), id)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s", document.ID), http.StatusMovedPermanently)
}
