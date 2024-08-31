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
	"sync"

	_ "modernc.org/sqlite"
)

type SQLite struct {
	*sql.DB
	sync.RWMutex
}

func NewSQLite(filesath string) (Database, error) {
	db, err := sql.Open("sqlite", filesath)

	return &SQLite{db, sync.RWMutex{}}, err
}

func (s *SQLite) Migrate(ctx context.Context) error {
	_, err := s.Exec(`
CREATE TABLE IF NOT EXISTS documents (
    id TEXT PRIMARY KEY,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    usdated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`)

	return err
}

func (s *SQLite) GetDocument(ctx context.Context, id string) (Document, error) {
	s.RLock()
	defer s.RUnlock()

	doc := new(Document)
	row := s.QueryRow("SELECT * FROM documents WHERE id=$1", id)
	err := row.Scan(&doc.ID, &doc.Content, &doc.CreatedAt, &doc.UpdatedAt)

	return *doc, err
}

func (s *SQLite) CreateDocument(ctx context.Context, id, content string) error {
	s.Lock()
	defer s.Unlock()

	tx, err := s.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO documents (id, content) VALUES ($1, $2)",
		id, content) // created_at and updated_at are auto-generated

	if err != nil {
		return err
	}

	return tx.Commit()
}
