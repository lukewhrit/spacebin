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
	"strings"
	"time"

	driver "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/lukewhrit/spacebin/internal/util"
)

type MySQL struct {
	*sql.DB
	cfg *driver.Config
}

// NewMySQL opens a MySQL database from uri, which is a go-sql-driver DSN prefixed with a scheme,
// such as "mysql://user:password@tcp(localhost:3306)/spacebin".
func NewMySQL(uri string) (Database, error) {
	_, dsn, _ := strings.Cut(uri, "://")
	cfg, err := driver.ParseDSN(dsn)

	if err != nil {
		return nil, err
	}

	// Required to scan timestamps into time.Time
	cfg.ParseTime = true

	// NewConnector only fails for configs that ParseDSN already rejects
	connector, _ := driver.NewConnector(cfg)

	db := sql.OpenDB(connector)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySQL{db, cfg}, nil
}

func (m *MySQL) Migrate(ctx context.Context) error {
	// Migrations contain multiple statements, which MySQL only runs on connections that allow it
	cfg := m.cfg.Clone()
	cfg.MultiStatements = true

	connector, _ := driver.NewConnector(cfg)
	db := sql.OpenDB(connector)
	defer db.Close()

	migrationDriver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		return err
	}

	return migrateUp(migrationFS, "migrations/mysql", migrationDriver)
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
	_, err := m.ExecContext(ctx, "INSERT INTO documents (id, content, username) VALUES (?, ?, ?)", id, content, username)

	return err
}

func (m *MySQL) UpdateDocument(ctx context.Context, id, content string) error {
	_, err := m.ExecContext(ctx, "UPDATE documents SET content=?, updated_at=NOW() WHERE id=?", content, id)

	return err
}

func (m *MySQL) DeleteDocument(ctx context.Context, id string) error {
	_, err := m.ExecContext(ctx, "DELETE FROM documents WHERE id=?", id)

	return err
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
	// Hash and salt the password
	hash, err := util.HashAndSalt([]byte(password))

	if err != nil {
		return err
	}

	// Add account to database
	_, err = m.ExecContext(ctx, "INSERT INTO accounts (username, password) VALUES (?, ?)", username, hash)

	return err
}

func (m *MySQL) DeleteAccount(ctx context.Context, id string) error {
	_, err := m.ExecContext(ctx, "DELETE FROM accounts WHERE id=?", id)

	return err
}

func (m *MySQL) GetSession(ctx context.Context, id string) (Session, error) {
	session := new(Session)
	row := m.QueryRow("SELECT public, token, secret, username FROM sessions WHERE public=?", id)
	err := row.Scan(&session.Public, &session.Token, &session.Secret, &session.Username)

	return *session, err
}

func (m *MySQL) CreateSession(ctx context.Context, public, token, secret, username string) error {
	_, err := m.ExecContext(ctx, "INSERT INTO sessions (public, token, secret, username) VALUES (?, ?, ?, ?)", public, token, secret, username)

	return err
}

func (m *MySQL) DeleteSession(ctx context.Context, public string) error {
	_, err := m.ExecContext(ctx, "DELETE FROM sessions WHERE public=?", public)

	return err
}
