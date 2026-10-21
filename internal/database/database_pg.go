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

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"github.com/lukewhrit/spacebin/internal/util"
)

type Postgres struct {
	*sql.DB
}

func NewPostgres(uri string) (Database, error) {
	db, err := sql.Open("postgres", uri)

	return &Postgres{db}, err
}

func (p *Postgres) Migrate(ctx context.Context) error {
	driver, err := postgres.WithInstance(p.DB, &postgres.Config{})

	if err != nil {
		return err
	}

	return migrateUp(migrationFS, "migrations/postgres", driver)
}

func (p *Postgres) GetDocument(ctx context.Context, id string) (Document, error) {
	doc := new(Document)
	row := p.QueryRow("SELECT id, content, username, encrypted, created_at, updated_at FROM documents WHERE id=$1", id)
	err := row.Scan(&doc.ID, &doc.Content, &doc.Username, &doc.Encrypted, &doc.CreatedAt, &doc.UpdatedAt)

	return *doc, err
}

func (p *Postgres) GetDocumentsByUsername(ctx context.Context, username string) ([]Document, error) {
	rows, err := p.QueryContext(ctx, "SELECT id, content, username, encrypted, created_at, updated_at FROM documents WHERE username=$1 ORDER BY created_at DESC", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []Document
	for rows.Next() {
		var doc Document
		if err := rows.Scan(&doc.ID, &doc.Content, &doc.Username, &doc.Encrypted, &doc.CreatedAt, &doc.UpdatedAt); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, rows.Err()
}

func (p *Postgres) CreateDocument(ctx context.Context, id, content, username string, encrypted bool) error {
	_, err := p.ExecContext(ctx, "INSERT INTO documents (id, content, username, encrypted) VALUES ($1, $2, $3, $4)", id, content, username, encrypted)

	return err
}

func (p *Postgres) UpdateDocument(ctx context.Context, id, content string) error {
	_, err := p.ExecContext(ctx, "UPDATE documents SET content=$1, updated_at=now() WHERE id=$2", content, id)

	return err
}

func (p *Postgres) DeleteDocument(ctx context.Context, id string) error {
	_, err := p.ExecContext(ctx, "DELETE FROM documents WHERE id=$1", id)

	return err
}

func (p *Postgres) GetCollection(ctx context.Context, id string) (Collection, error) {
	col := new(Collection)
	row := p.QueryRowContext(ctx, "SELECT id, title, username, created_at FROM collections WHERE id=$1", id)
	if err := row.Scan(&col.ID, &col.Title, &col.Username, &col.CreatedAt); err != nil {
		return *col, err
	}

	// Joining on documents leaves out documents that have been deleted since
	rows, err := p.QueryContext(ctx, "SELECT cd.document_id FROM collection_documents cd JOIN documents d ON d.id = cd.document_id WHERE cd.collection_id=$1 ORDER BY cd.position", id)
	if err != nil {
		return *col, err
	}

	col.Documents, err = scanDocumentIDs(rows)

	return *col, err
}

func (p *Postgres) GetCollectionsByUsername(ctx context.Context, username string) ([]Collection, error) {
	rows, err := p.QueryContext(ctx, "SELECT id, title, username, created_at FROM collections WHERE username=$1 ORDER BY created_at DESC", username)
	if err != nil {
		return nil, err
	}

	return scanCollections(rows)
}

func (p *Postgres) CreateCollection(ctx context.Context, id, title, username string, documents []string) error {
	tx, err := p.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "INSERT INTO collections (id, title, username) VALUES ($1, $2, $3)", id, title, username); err != nil {
		return err
	}

	for i, document := range documents {
		if _, err := tx.ExecContext(ctx, "INSERT INTO collection_documents (collection_id, document_id, position) VALUES ($1, $2, $3)", id, document, i); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (p *Postgres) DeleteCollection(ctx context.Context, id string) error {
	tx, err := p.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM collection_documents WHERE collection_id=$1", id); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM collections WHERE id=$1", id); err != nil {
		return err
	}

	return tx.Commit()
}

func (p *Postgres) GetAccount(ctx context.Context, id string) (Account, error) {
	account := new(Account)
	row := p.QueryRow("SELECT * FROM accounts WHERE id=$1", id)
	err := row.Scan(&account.ID, &account.Username, &account.Password)

	return *account, err
}

func (p *Postgres) GetAccountByUsername(ctx context.Context, username string) (Account, error) {
	account := new(Account)
	row := p.QueryRow("SELECT * FROM accounts WHERE username=$1", username)
	err := row.Scan(&account.ID, &account.Username, &account.Password)

	return *account, err
}

func (p *Postgres) CreateAccount(ctx context.Context, username, password string) error {
	// Hash and salt the password
	hash, err := util.HashAndSalt([]byte(password))

	if err != nil {
		return err
	}

	// Add account to database
	_, err = p.ExecContext(ctx, "INSERT INTO accounts (username, password) VALUES ($1, $2)", username, hash)

	return err
}

func (p *Postgres) DeleteAccount(ctx context.Context, id string) error {
	_, err := p.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1", id)

	return err
}

func (p *Postgres) GetSession(ctx context.Context, id string) (Session, error) {
	session := new(Session)
	row := p.QueryRow("SELECT public, token, secret, username FROM sessions WHERE public=$1", id)
	err := row.Scan(&session.Public, &session.Token, &session.Secret, &session.Username)

	return *session, err
}

func (p *Postgres) CreateSession(ctx context.Context, public, token, secret, username string) error {
	_, err := p.ExecContext(ctx, "INSERT INTO sessions (public, token, secret, username) VALUES ($1, $2, $3, $4)", public, token, secret, username)

	return err
}

func (p *Postgres) DeleteSession(ctx context.Context, public string) error {
	_, err := p.ExecContext(ctx, "DELETE FROM sessions WHERE public=$1", public)

	return err
}
