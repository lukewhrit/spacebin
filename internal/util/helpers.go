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

package util

import (
	"encoding/json"
	"math"
	"net/http"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/orca-group/spirit/internal/config"
)

type CreateRequest struct {
	Content   string
	Extension string
}

func ValidateBody(body CreateRequest) error {
	regex := regexp.MustCompile("^python$|^javascript$|^jsx$|^typescript$|^tsx$|^go$|^kotlin$|^cpp$|^sql$|^csharp$|^c$|^scala$|^haskell$|^shell-session$|^bash$|^powershell$|^php$|^asm6502$|^julia$|^objc$|^perl$|^crystal$|^json$|^yaml$|^toml$|^none$|^rust$|^ruby$|^markup$|^markdown$|^css$|")

	return validation.ValidateStruct(&body,
		validation.Field(
			&body.Content,
			validation.Required,
			// Enforce length to follow what's set in the config
			validation.Length(2, config.Config.Documents.MaxDocumentLength),
		),
		// The purpose of this field is to support client's that perform
		// syntax highlighting and need to know what highlighter to use.
		validation.Field(
			&body.Extension,
			validation.Match(regex),
			validation.Required,
		),
	)
}

// HandleBody figures out whether a incoming request is in JSON or multipart/form-data and decodes it appropriately
func HandleBody(r *http.Request) (CreateRequest, error) {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		resp := make(map[string]string)

		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return CreateRequest{}, err
		}

		return CreateRequest{
			Content:   resp["content"],
			Extension: resp["extension"],
		}, nil
	case "multipart/form-data":
		err := r.ParseMultipartForm(int64(float64(config.Config.Documents.MaxDocumentLength) * math.Pow(1024, 2)))

		if err != nil {
			return CreateRequest{}, err
		}

		return CreateRequest{
			Content:   r.FormValue("content"),
			Extension: r.FormValue("extension"),
		}, err
	}

	return CreateRequest{}, nil
}

// WriteJSON writes a Request payload (p) to an HTTP response writer (w)
func (p Payload) WriteJSON(w http.ResponseWriter, status int) error {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Payload: p,
		Error:   "",
		Status:  status,
	})

	return nil
}

// WriteError writes an Error object (e) to an HTTP response writer (w)
func WriteError(e error, w http.ResponseWriter, status int) error {
	bytes, err := json.Marshal(Response{
		Error:   e.Error(),
		Payload: Payload{},
		Status:  status,
	})

	if err != nil {
		return err
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)

	return nil
}
