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
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func newTestSQLite(t *testing.T) *SQLite {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	s := &SQLite{db, sync.RWMutex{}}
	require.NoError(t, s.Migrate(context.Background()))
	return s
}

// Document tests

func TestSQLiteCreateAndGetDocument(t *testing.T) {
	s := newTestSQLite(t)
	ctx := context.Background()

	err := s.CreateDocument(ctx, "abc123", "hello world", "alice")
	require.NoError(t, err)

	doc, err := s.GetDocument(ctx, "abc123")
	require.NoError(t, err)
	require.Equal(t, "abc123", doc.ID)
	require.Equal(t, "hello world", doc.Content)
	require.Equal(t, "alice", doc.Username)
	require.NotZero(t, doc.CreatedAt)
}

func TestSQLiteGetDocumentNotFound(t *testing.T) {
	s := newTestSQLite(t)
	_, err := s.GetDocument(context.Background(), "nope")
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSQLiteUpdateDocument(t *testing.T) {
	s := newTestSQLite(t)
	ctx := context.Background()

	require.NoError(t, s.CreateDocument(ctx, "doc1", "original", "alice"))

	require.NoError(t, s.UpdateDocument(ctx, "doc1", "updated content"))

	doc, err := s.GetDocument(ctx, "doc1")
	require.NoError(t, err)
	require.Equal(t, "updated content", doc.Content)
}

func TestSQLiteDeleteDocument(t *testing.T) {
	s := newTestSQLite(t)
	ctx := context.Background()

	require.NoError(t, s.CreateDocument(ctx, "doc2", "to delete", "bob"))
	require.NoError(t, s.DeleteDocument(ctx, "doc2"))

	_, err := s.GetDocument(ctx, "doc2")
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSQLiteGetDocumentsByUsername(t *testing.T) {
	s := newTestSQLite(t)
	ctx := context.Background()

	require.NoError(t, s.CreateDocument(ctx, "d1", "content1", "alice"))
	require.NoError(t, s.CreateDocument(ctx, "d2", "content2", "alice"))
	require.NoError(t, s.CreateDocument(ctx, "d3", "content3", "bob"))

	aliceDocs, err := s.GetDocumentsByUsername(ctx, "alice")
	require.NoError(t, err)
	require.Len(t, aliceDocs, 2)

	bobDocs, err := s.GetDocumentsByUsername(ctx, "bob")
	require.NoError(t, err)
	require.Len(t, bobDocs, 1)
	require.Equal(t, "d3", bobDocs[0].ID)

	noDocs, err := s.GetDocumentsByUsername(ctx, "nobody")
	require.NoError(t, err)
	require.Empty(t, noDocs)
}

// Account tests

func TestSQLiteCreateAndGetAccount(t *testing.T) {
	s := newTestSQLite(t)
	ctx := context.Background()

	err := s.CreateAccount(ctx, "alice", "password123")
	require.NoError(t, err)

	// GetAccount requires the integer ID — fetch it via username first
	acc, err := s.GetAccountByUsername(ctx, "alice")
	require.NoError(t, err)
	require.Equal(t, "alice", acc.Username)
	require.NotEmpty(t, acc.Password)

	// Verify GetAccount by ID also works
	acc2, err := s.GetAccount(ctx, fmt.Sprint(acc.ID))
	require.NoError(t, err)
	require.Equal(t, acc.ID, acc2.ID)
	require.Equal(t, "alice", acc2.Username)
}

func TestSQLiteGetAccountByUsername(t *testing.T) {
	s := newTestSQLite(t)
	ctx := context.Background()

	require.NoError(t, s.CreateAccount(ctx, "bob", "secret"))

	acc, err := s.GetAccountByUsername(ctx, "bob")
	require.NoError(t, err)
	require.Equal(t, "bob", acc.Username)
}

func TestSQLiteGetAccountNotFound(t *testing.T) {
	s := newTestSQLite(t)
	_, err := s.GetAccountByUsername(context.Background(), "ghost")
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSQLiteDeleteAccount(t *testing.T) {
	s := newTestSQLite(t)
	ctx := context.Background()

	require.NoError(t, s.CreateAccount(ctx, "carol", "pw"))

	acc, err := s.GetAccountByUsername(ctx, "carol")
	require.NoError(t, err)

	require.NoError(t, s.DeleteAccount(ctx, fmt.Sprint(acc.ID)))

	_, err = s.GetAccountByUsername(ctx, "carol")
	require.ErrorIs(t, err, sql.ErrNoRows)
}

// Session tests

func TestSQLiteCreateAndGetSession(t *testing.T) {
	s := newTestSQLite(t)
	ctx := context.Background()

	err := s.CreateSession(ctx, "pub1", "tok1", "sec1", "alice")
	require.NoError(t, err)

	sess, err := s.GetSession(ctx, "pub1")
	require.NoError(t, err)
	require.Equal(t, "pub1", sess.Public)
	require.Equal(t, "tok1", sess.Token)
	require.Equal(t, "sec1", sess.Secret)
	require.Equal(t, "alice", sess.Username)
}

func TestSQLiteGetSessionNotFound(t *testing.T) {
	s := newTestSQLite(t)
	_, err := s.GetSession(context.Background(), "nosuchsession")
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSQLiteDeleteSession(t *testing.T) {
	s := newTestSQLite(t)
	ctx := context.Background()

	require.NoError(t, s.CreateSession(ctx, "pub2", "tok2", "sec2", "bob"))
	require.NoError(t, s.DeleteSession(ctx, "pub2"))

	_, err := s.GetSession(ctx, "pub2")
	require.ErrorIs(t, err, sql.ErrNoRows)
}
