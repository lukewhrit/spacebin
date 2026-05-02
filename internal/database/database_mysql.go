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
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/lukewhrit/spacebin/internal/util"
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
	_ = ctx

	driver, err := mysql.WithInstance(m.DB, &mysql.Config{})

	if err != nil {
		return err
	}

	source, err := iofs.New(migrationFS, "migrations/mysql")

	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "mysql", driver)

	if err != nil {
		return err
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (m *MySQL) GetDocument(ctx context.Context, id string) (Document, error) {
	doc := new(Document)
	row := m.QueryRow("SELECT id, content, username, created_at, updated_at FROM documents WHERE id=?", id)
	err := row.Scan(&doc.ID, &doc.Content, &doc.Username, &doc.CreatedAt, &doc.UpdatedAt)

	return *doc, err
}

func (m *MySQL) GetDocumentsByUsername(ctx context.Context, username string) ([]Document, error) {
	rows, err := m.QueryContext(ctx, "SELECT id, content, username, created_at, updated_at FROM documents WHERE username=? ORDER BY created_at DESC", username)
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

func (m *MySQL) CreateDocument(ctx context.Context, id, content, username string) error {
	tx, err := m.Begin()

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

func (m *MySQL) UpdateDocument(ctx context.Context, id, content string) error {
	tx, err := m.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE documents SET content=?, updated_at=NOW() WHERE id=?", content, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *MySQL) DeleteDocument(ctx context.Context, id string) error {
	tx, err := m.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM documents WHERE id=?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *MySQL) GetAccount(ctx context.Context, id string) (Account, error) {
	acc := new(Account)
	row := m.QueryRow("SELECT * FROM accounts WHERE id=?", id)
	err := row.Scan(&acc.ID, &acc.Username, &acc.Password)

	return *acc, err
}

func (m *MySQL) GetAccountByUsername(ctx context.Context, username string) (Account, error) {
	account := new(Account)
	row := m.QueryRow("SELECT * FROM accounts WHERE username=?", username)
	err := row.Scan(&account.ID, &account.Username, &account.Password)

	return *account, err
}

func (m *MySQL) CreateAccount(ctx context.Context, username, password string) error {
	tx, err := m.Begin()

	if err != nil {
		return err
	}

	// Add account to database
	// Hash and salt the password
	_, err = tx.Exec("INSERT INTO accounts (username, password) VALUES (?, ?)",
		username, util.HashAndSalt([]byte(password)))

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *MySQL) DeleteAccount(ctx context.Context, id string) error {
	tx, err := m.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM accounts WHERE id=?", id)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *MySQL) GetSession(ctx context.Context, id string) (Session, error) {
	session := new(Session)
	row := m.QueryRow("SELECT public, token, secret, username FROM sessions WHERE public=?", id)
	err := row.Scan(&session.Public, &session.Token, &session.Secret, &session.Username)

	return *session, err
}

func (m *MySQL) CreateSession(ctx context.Context, public, token, secret, username string) error {
	tx, err := m.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO sessions (public, token, secret, username) VALUES (?, ?, ?, ?)",
		public, token, secret, username)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *MySQL) DeleteSession(ctx context.Context, public string) error {
	tx, err := m.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM sessions WHERE public=?", public)

	if err != nil {
		return err
	}

	return tx.Commit()
}
