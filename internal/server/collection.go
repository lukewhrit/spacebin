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

package server

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lukewhrit/spacebin/internal/util"
	"golang.org/x/exp/slices"
)

// maxTitleLength is the longest title a collection can have, limited by MySQL's varchar(255)
const maxTitleLength = 255

// validID reports whether id has the shape of a document or collection ID.
func (s *Server) validID(id string) bool {
	return len(id) == s.Config.IDLength || slices.Contains(s.Config.Documents, id)
}

// createCollection handles the shared logic between the CreateCollection and StaticCreateCollection handlers.
func createCollection(s *Server, r *http.Request) (string, error) {
	body, err := util.HandleCollectionBody(s.Config.MaxSize, r)

	if err != nil {
		return "", fmt.Errorf("bad request: %v", err)
	}

	title := strings.TrimSpace(body.Title)

	if len(title) > maxTitleLength {
		return "", fmt.Errorf("bad request: title must be at most %d characters", maxTitleLength)
	}

	// Drop repeated documents, keeping the first position of each
	var documents []string

	for _, id := range body.Documents {
		if !slices.Contains(documents, id) {
			documents = append(documents, id)
		}
	}

	if len(documents) == 0 {
		return "", errors.New("bad request: a collection needs at least one document")
	}

	if len(documents) > s.Config.MaxCollectionSize {
		return "", fmt.Errorf("bad request: a collection can hold at most %d documents", s.Config.MaxCollectionSize)
	}

	// Every document must exist
	for _, id := range documents {
		if !s.validID(id) {
			return "", fmt.Errorf("bad request: %q is not a valid document id", id)
		}

		if _, err := s.Database.GetDocument(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", fmt.Errorf("bad request: document %s not found", id)
			}

			return "", err
		}
	}

	id := util.GenerateID(s.Config.IDType, s.Config.IDLength)

	// Track owner if authenticated
	username, _ := s.authenticatedUsername(r)

	if err := s.Database.CreateCollection(r.Context(), id, title, username, documents); err != nil {
		return "", err
	}

	return id, nil
}

func (s *Server) CreateCollection(w http.ResponseWriter, r *http.Request) {
	id, err := createCollection(s, r)

	if err != nil {
		if strings.Contains(err.Error(), "bad request:") {
			util.WriteError(w, http.StatusBadRequest, err)
			return
		}

		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	collection, err := s.Database.GetCollection(r.Context(), id)

	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, collection)
}

func (s *Server) StaticCreateCollection(w http.ResponseWriter, r *http.Request) {
	id, err := createCollection(s, r)

	if err != nil {
		if strings.Contains(err.Error(), "bad request:") {
			util.RenderError(&resources, w, http.StatusBadRequest, err)
			return
		}

		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/c/%s", id), http.StatusSeeOther)
}

func (s *Server) FetchCollection(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "collection")

	if !s.validID(id) {
		err := fmt.Errorf("id is of length %d, should be %d", len(id), s.Config.IDLength)
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	collection, err := s.Database.GetCollection(r.Context(), id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.WriteError(w, http.StatusNotFound, err)
			return
		}

		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, collection)
}

// collectionDocument is a highlighted document shown on a collection's page
type collectionDocument struct {
	ID          string
	Highlighted template.HTML
}

func (s *Server) StaticCollection(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "collection")

	if !s.validID(id) {
		err := fmt.Errorf("id is of length %d, should be %d", len(id), s.Config.IDLength)
		util.RenderError(&resources, w, http.StatusBadRequest, err)
		return
	}

	username := s.getUsername(r)

	collection, err := s.Database.GetCollection(r.Context(), id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.RenderError(&resources, w, http.StatusNotFound, err)
			return
		}

		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	var (
		documents  []collectionDocument
		stylesheet string
	)

	for _, documentID := range collection.Documents {
		document, err := s.Database.GetDocument(r.Context(), documentID)

		if err != nil {
			// The document was deleted after the collection was read
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}

			util.RenderError(&resources, w, http.StatusInternalServerError, err)
			return
		}

		// Every document uses the same style, so the stylesheet is the same each time
		var highlighted string
		highlighted, stylesheet = util.Highlight(document.Content, "")

		documents = append(documents, collectionDocument{
			ID:          document.ID,
			Highlighted: template.HTML(highlighted),
		})
	}

	render(w, "collection.html", s.pageData(username, map[string]any{
		"ID":         collection.ID,
		"Title":      collection.Title,
		"Documents":  documents,
		"Stylesheet": template.CSS(stylesheet),
		"IsOwner":    username != "" && collection.Username == username,
	}))
}

func (s *Server) StaticNewCollection(w http.ResponseWriter, r *http.Request) {
	username := s.getUsername(r)

	// Lets document pages link here with themselves already filled in
	documents := strings.Join(util.ParseDocumentIDs(r.URL.Query().Get("documents")), "\n")

	render(w, "collection_new.html", s.pageData(username, map[string]any{
		"Documents":         documents,
		"MaxCollectionSize": s.Config.MaxCollectionSize,
	}))
}

func (s *Server) RemoveCollection(w http.ResponseWriter, r *http.Request) {
	if !s.Config.AccountsEnabled {
		util.WriteError(w, http.StatusNotFound, errors.New("accounts disabled"))
		return
	}

	username, err := s.authenticatedUsername(r)
	if err != nil || username == "" {
		util.WriteError(w, http.StatusUnauthorized, errors.New("authentication required"))
		return
	}

	id := chi.URLParam(r, "collection")
	if !s.validID(id) {
		util.RenderError(&resources, w, http.StatusBadRequest, fmt.Errorf("invalid collection id"))
		return
	}

	collection, err := s.Database.GetCollection(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			util.RenderError(&resources, w, http.StatusNotFound, err)
			return
		}
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	if collection.Username != username {
		util.RenderError(&resources, w, http.StatusForbidden, errors.New("you do not own this collection"))
		return
	}

	if err := s.Database.DeleteCollection(r.Context(), id); err != nil {
		util.RenderError(&resources, w, http.StatusInternalServerError, err)
		return
	}

	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.Redirect(w, r, "/account", http.StatusSeeOther)
}
