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
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/orca-group/spirit/internal/config"
	"github.com/orca-group/spirit/internal/database"
	"github.com/orca-group/spirit/internal/util"
)

func FetchDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "document")

	if len(id) != config.Config.IDLength {
		err := fmt.Errorf("id is of length %d, should be %d", len(id), config.Config.IDLength)
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	document, err := database.FindDocument(id)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			util.WriteError(w, http.StatusNotFound, err)
		} else {
			util.WriteError(w, http.StatusInternalServerError, err)
		}

		return
	}

	if err := util.WriteJSON(w, http.StatusOK, document); err != nil {
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

	document, err := database.FindDocument(id)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(document.Content))
}
