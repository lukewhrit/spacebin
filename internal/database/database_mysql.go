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

package database

import (
	"context"
	"database/sql"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	*sql.DB
}

func NewMySQL(uri *url.URL) (Database, error) {
	_, uriTrimmed, _ := strings.Cut(uri.String(), uri.Scheme+"://")
	db, err := sql.Open("mysql", uriTrimmed)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySQL{db}, err
}

func (m *MySQL) Migrate(ctx context.Context) error {
	_, err := m.Exec(`
CREATE TABLE IF NOT EXISTS documents (
	id VARCHAR(255) PRIMARY KEY,
	content TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`)

	return err
}

func (m *MySQL) GetDocument(ctx context.Context, id string) (Document, error) {
	doc := new(Document)
	row := m.QueryRow("SELECT * FROM documents WHERE id=?", id)
	err := row.Scan(&doc.ID, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt)

	return *doc, err
}

func (m *MySQL) CreateDocument(ctx context.Context, id, content string) error {
	tx, err := m.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO documents (id, content) VALUES (?, ?)",
		id, content) // created_at and updated_at are auto-generated

	if err != nil {
		return err
	}

	return tx.Commit()
}
