/*
 * Copyright 2020 Luke Whrit, Jack Dorland; The Spacebin Authors

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

package document

import "github.com/gofiber/fiber"

// Payload is a document object
type Payload struct {
	ContentHash string  `json:"content_hash,omitempty"` // A base64 representation form of the document's content.
	ID          *string `json:"id,omitempty"`           // The document ID.
	Content     *string `json:"content,omitempty"`      // The document content.
	Extension   *string `json:"extension,omitempty"`    // The extension of the document.
	CreatedAt   *int64  `json:"created_at,omitempty"`   // The Unix timestamp of when the document was inserted.
	UpdatedAt   *int64  `json:"updated_at,omitempty"`   // The Unix timestmap of when the document was last modified.
	Exists      *bool   `json:"exists,omitempty"`       // Whether the document does or does not exist.
}

// Response is a Spacebin API response
type Response struct {
	Error   string  `json:"error"` // .Error() should already be called
	Payload Payload `json:"payload"`
	Status  int     `json:"status"`
}

// Register loads all document-related endpoints
func Register(app *fiber.App) {
	api := app.Group("/api/v1/documents")

	registerCreate(api)
	registerRead(api)
}
