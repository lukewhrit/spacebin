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

package util

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"math"
	"net/http"
	"net/url"
	"strings"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rs/zerolog/log"
)

// readAll reads uploaded files. It is a variable so tests can simulate read failures.
var readAll = io.ReadAll

func ValidateBody(maxSize int, body CreateRequest) error {
	return validation.ValidateStruct(&body,
		validation.Field(&body.Content, validation.Required, validation.Length(2, maxSize)),
	)
}

func HandleCreateBody(maxSize int, r *http.Request) (re CreateRequest, e error) {
	// Ignore charset or boundary fields, just get type of content
	switch contentType := strings.Split(r.Header.Get("Content-Type"), ";")[0]; contentType {
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

		// Try to get the "content" field as plain text
		content := r.FormValue("content")

		if content != "" {
			return CreateRequest{
				Content: content,
			}, nil
		}

		// If "content" is not plain text, check for file uploads
		file, _, err := r.FormFile("content") // Access file under the "content" name
		if err != nil {
			return CreateRequest{}, fmt.Errorf("failed to parse content field as file: %w", err)
		}
		defer file.Close()

		// Read the uploaded file's content
		fileContent, err := readAll(file)

		if err != nil {
			return CreateRequest{}, fmt.Errorf("failed to read uploaded file content: %w", err)
		}

		return CreateRequest{
			Content: string(fileContent),
		}, nil
	default:
		return CreateRequest{}, fmt.Errorf("unsupported Content-Type: %s", contentType)
	}
}

// HandleCollectionBody handles the body of a request to create a collection. Multipart forms may list
// documents in one or more "documents" fields, separated by commas or whitespace.
func HandleCollectionBody(maxSize int, r *http.Request) (CollectionRequest, error) {
	// Ignore charset or boundary fields, just get type of content
	switch contentType := strings.Split(r.Header.Get("Content-Type"), ";")[0]; contentType {
	case "application/json":
		var req CollectionRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return CollectionRequest{}, err
		}

		req.Documents = ParseDocumentIDs(strings.Join(req.Documents, " "))

		return req, nil
	case "multipart/form-data":
		err := r.ParseMultipartForm(int64(float64(maxSize) * math.Pow(1024, 2)))

		if err != nil {
			return CollectionRequest{}, err
		}

		return CollectionRequest{
			Title:     r.FormValue("title"),
			Documents: ParseDocumentIDs(strings.Join(r.MultipartForm.Value["documents"], " ")),
		}, nil
	default:
		return CollectionRequest{}, fmt.Errorf("unsupported Content-Type: %s", contentType)
	}
}

// ParseDocumentIDs splits s on commas and whitespace into document IDs. Each entry may also be a link to a
// document, such as "https://spaceb.in/AbCd1234.go" or "/AbCd1234/raw", in which case only the ID is kept.
func ParseDocumentIDs(s string) []string {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || unicode.IsSpace(r)
	})

	ids := make([]string, 0, len(fields))

	for _, field := range fields {
		// Drops the scheme, host, query and fragment of links
		if u, err := url.Parse(field); err == nil {
			field = u.Path
		}

		field = strings.TrimSuffix(strings.TrimRight(field, "/"), "/raw")
		field = field[strings.LastIndex(field, "/")+1:]
		field, _, _ = strings.Cut(field, ".")

		if field != "" {
			ids = append(ids, field)
		}
	}

	return ids
}

// HandleSignupBody handles the body of a Signup request
func HandleSignupBody(maxSize int, r *http.Request) (re SignupRequest, e error) {
	// Ignore charset or boundary fields, just get type of content
	switch strings.Split(r.Header.Get("Content-Type"), ";")[0] {
	case "application/json":
		resp := make(map[string]string)

		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return SignupRequest{}, err
		}

		return SignupRequest{
			Username: resp["username"],
			Password: resp["password"],
		}, nil
	case "multipart/form-data":
		err := r.ParseMultipartForm(int64(float64(maxSize) * math.Pow(1024, 2)))

		if err != nil {
			return SignupRequest{}, err
		}

		return SignupRequest{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}, nil
	}

	return SignupRequest{}, nil
}

// HandleSigninBody handles the body of a Signin request
func HandleSigninBody(maxSize int, r *http.Request) (re SigninRequest, e error) {
	// Ignore charset or boundary fields, just get type of content
	switch strings.Split(r.Header.Get("Content-Type"), ";")[0] {
	case "application/json":
		resp := make(map[string]string)

		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			return SigninRequest{}, err
		}

		return SigninRequest{
			Username: resp["username"],
			Password: resp["password"],
		}, nil
	case "multipart/form-data":
		err := r.ParseMultipartForm(int64(float64(maxSize) * math.Pow(1024, 2)))

		if err != nil {
			return SigninRequest{}, err
		}

		return SigninRequest{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}, nil
	}

	return SigninRequest{}, nil
}

// WriteJSON writes a Request payload (p) to an HTTP response writer (w)
func WriteJSON[R any](w http.ResponseWriter, status int, r R) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"payload": r,
		"error":   "",
	})
}

// WriteError writes an Error object (e) to an HTTP response writer (w)
func WriteError(w http.ResponseWriter, status int, e error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"payload": map[string]interface{}{},
		"error":   e.Error(),
	})

	log.Debug().Err(e).Msg("Request Error")
}

// RenderError renders errors to the client using an HTML template.
func RenderError(r fs.FS, w http.ResponseWriter, status int, err error) error {
	tmpl := template.Must(template.ParseFS(r, "web/error.html"))

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)

	data := struct {
		Status string
		Error  string
	}{
		Status: strings.Join([]string{fmt.Sprintf("%d", status), http.StatusText(status)}, " "),
		Error:  err.Error(),
	}

	return tmpl.Execute(w, data)
}
