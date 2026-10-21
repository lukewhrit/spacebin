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

package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Document struct {
	ID        string    `db:"id" json:"id"`
	Content   string    `db:"content" json:"content"`
	Username  string    `db:"username" json:"username"`
	Encrypted bool      `db:"encrypted" json:"encrypted"` // Content is a ciphertext envelope only browsers can decrypt
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Collection is an ordered group of documents shared under one ID
type Collection struct {
	ID        string    `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Username  string    `db:"username" json:"username"`
	Documents []string  `db:"-" json:"documents"` // Document IDs, in order. Deleted documents are left out
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Account struct {
	ID       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	// Documents []Document `db:"documents" json:"documents"`
}

type Session struct {
	Public   string `db:"public" json:"public"`
	Token    string `db:"token" json:"token"`
	Secret   string `db:"secret" json:"secret"`
	Username string `db:"username" json:"username"`
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Database
type Database interface {
	Migrate(ctx context.Context) error
	Close() error

	GetDocument(ctx context.Context, id string) (Document, error)
	GetDocumentsByUsername(ctx context.Context, username string) ([]Document, error)
	CreateDocument(ctx context.Context, id, content, username string, encrypted bool) error
	UpdateDocument(ctx context.Context, id, content string) error
	DeleteDocument(ctx context.Context, id string) error

	GetCollection(ctx context.Context, id string) (Collection, error)
	// GetCollectionsByUsername returns a user's collections without their Documents
	GetCollectionsByUsername(ctx context.Context, username string) ([]Collection, error)
	CreateCollection(ctx context.Context, id, title, username string, documents []string) error
	DeleteCollection(ctx context.Context, id string) error

	GetAccount(ctx context.Context, id string) (Account, error)
	GetAccountByUsername(ctx context.Context, username string) (Account, error)
	CreateAccount(ctx context.Context, username, password string) error
	// UpdateAccount(ctx context.Context, id, username, password string) error
	DeleteAccount(ctx context.Context, id string) error

	GetSession(ctx context.Context, id string) (Session, error)
	CreateSession(ctx context.Context, public, token, secret, username string) error
	DeleteSession(ctx context.Context, public string) error
}

// scanCollections reads collections from rows selecting id, title, username and created_at.
func scanCollections(rows *sql.Rows) ([]Collection, error) {
	defer rows.Close()

	var cols []Collection
	for rows.Next() {
		var col Collection
		if err := rows.Scan(&col.ID, &col.Title, &col.Username, &col.CreatedAt); err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}
	return cols, rows.Err()
}

// scanDocumentIDs reads a single column of document IDs from rows.
func scanDocumentIDs(rows *sql.Rows) ([]string, error) {
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
