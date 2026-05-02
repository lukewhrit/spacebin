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
	"errors"
	"net/url"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/lukewhrit/spacebin/internal/util"
	_ "modernc.org/sqlite"
)

type SQLite struct {
	*sql.DB
	sync.RWMutex
}

func NewSQLite(uri *url.URL) (Database, error) {
	dbPath := uri.Path

	if uri.Scheme == "sqlite" && uri.Host == ":memory:" {
		dbPath = ":memory:"
	} else {
		dbPath = uri.Path
		if len(dbPath) > 0 && dbPath[0] == '/' {
			dbPath = dbPath[1:]
		}
	}

	db, err := sql.Open("sqlite", dbPath)

	return &SQLite{db, sync.RWMutex{}}, err
}

func (s *SQLite) Migrate(ctx context.Context) error {
	_ = ctx

	s.Lock()
	defer s.Unlock()

	driver, err := sqlite.WithInstance(s.DB, &sqlite.Config{})

	if err != nil {
		return err
	}

	source, err := iofs.New(migrationFS, "migrations/sqlite")

	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "sqlite", driver)

	if err != nil {
		return err
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (s *SQLite) GetDocument(ctx context.Context, id string) (Document, error) {
	s.RLock()
	defer s.RUnlock()

	doc := new(Document)
	row := s.QueryRow("SELECT id, content, username, created_at, updated_at FROM documents WHERE id=?", id)
	err := row.Scan(&doc.ID, &doc.Content, &doc.Username, &doc.CreatedAt, &doc.UpdatedAt)

	return *doc, err
}

func (s *SQLite) GetDocumentsByUsername(ctx context.Context, username string) ([]Document, error) {
	s.RLock()
	defer s.RUnlock()

	rows, err := s.QueryContext(ctx, "SELECT id, content, username, created_at, updated_at FROM documents WHERE username=? ORDER BY created_at DESC", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []Document
	for rows.Next() {
		var doc Document
		if err := rows.Scan(&doc.ID, &doc.Content, &doc.Username, &doc.CreatedAt, &doc.UpdatedAt); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, rows.Err()
}

func (s *SQLite) CreateDocument(ctx context.Context, id, content, username string) error {
	s.Lock()
	defer s.Unlock()

	tx, err := s.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO documents (id, content, username) VALUES (?, ?, ?)",
		id, content, username)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SQLite) UpdateDocument(ctx context.Context, id, content string) error {
	s.Lock()
	defer s.Unlock()

	tx, err := s.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE documents SET content=?, updated_at=CURRENT_TIMESTAMP WHERE id=?", content, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SQLite) DeleteDocument(ctx context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	tx, err := s.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM documents WHERE id=?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SQLite) GetAccount(ctx context.Context, id string) (Account, error) {
	s.RLock()
	defer s.RUnlock()

	acc := new(Account)
	row := s.QueryRow("SELECT * FROM accounts WHERE id=$1", id)
	err := row.Scan(&acc.ID, &acc.Username, &acc.Password)

	return *acc, err
}

func (s *SQLite) GetAccountByUsername(ctx context.Context, username string) (Account, error) {
	account := new(Account)
	row := s.QueryRow("SELECT * FROM accounts WHERE username=$1", username)
	err := row.Scan(&account.ID, &account.Username, &account.Password)

	return *account, err
}

func (s *SQLite) CreateAccount(ctx context.Context, username, password string) error {
	s.Lock()
	defer s.Unlock()

	tx, err := s.Begin()

	if err != nil {
		return err
	}

	// Add account to database
	// Hash and salt the password
	_, err = tx.Exec("INSERT INTO accounts (username, password) VALUES ($1, $2)",
		username, util.HashAndSalt([]byte(password)))

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SQLite) DeleteAccount(ctx context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	tx, err := s.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM accounts WHERE id=$1", id)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SQLite) GetSession(ctx context.Context, id string) (Session, error) {
	s.RLock()
	defer s.RUnlock()

	session := new(Session)
	row := s.QueryRow("SELECT public, token, secret, username FROM sessions WHERE public=?", id)
	err := row.Scan(&session.Public, &session.Token, &session.Secret, &session.Username)

	return *session, err
}

func (s *SQLite) CreateSession(ctx context.Context, public, token, secret, username string) error {
	s.Lock()
	defer s.Unlock()

	tx, err := s.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO sessions (public, token, secret, username) VALUES ($1, $2, $3, $4)",
		public, token, secret, username)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SQLite) DeleteSession(ctx context.Context, public string) error {
	s.Lock()
	defer s.Unlock()

	tx, err := s.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM sessions WHERE public=$1", public)

	if err != nil {
		return err
	}

	return tx.Commit()
}
