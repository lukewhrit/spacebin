/*
 * Copyright 2020-2023 Luke Whritenour

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

package util

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/rs/zerolog/log"
)

type CreateRequest struct {
	Content string
}

func ValidateBody(maxSize int, body CreateRequest) error {
	return validation.ValidateStruct(&body,
		validation.Field(&body.Content, validation.Required,
			validation.Length(2, maxSize)),
	)
}

func CountLines(v string) template.HTML {
	var x []string

	for i := range strings.Split(v, "\n") {
		x = append(x, fmt.Sprintf("<div>%d</div>", i+1))
	}

	return template.HTML(strings.Join(x, ""))
}

// HandleBody figures out whether a incoming request is in JSON or multipart/form-data and decodes it appropriately
func HandleBody(maxSize int, r *http.Request) (CreateRequest, error) {
	// Ignore charset or boundary fields, just get type of content
	switch strings.Split(r.Header.Get("Content-Type"), ";")[0] {
	case "application/json":
		resp := make(map[string]string)

		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return CreateRequest{}, err
		}

		return CreateRequest{
			Content: resp["content"],
		}, nil
	case "multipart/form-data":
		err := r.ParseMultipartForm(int64(float64(maxSize) * math.Pow(1024, 2)))

		if err != nil {
			return CreateRequest{}, err
		}

		return CreateRequest{
			Content: r.FormValue("content"),
		}, nil
	}

	return CreateRequest{}, nil
}

// WriteJSON writes a Request payload (p) to an HTTP response writer (w)
func WriteJSON[R any](w http.ResponseWriter, status int, r R) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"payload": r,
		"error":   "",
	})

	return nil
}

// WriteError writes an Error object (e) to an HTTP response writer (w)
func WriteError(w http.ResponseWriter, status int, e error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"payload": map[string]interface{}{},
		"error":   e.Error(),
	})

	log.Debug().Err(e).Msg("Request Error")

	return nil
}
