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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// testDatabases are the databases every test in this file runs against. SQLite always runs in memory,
// Postgres and MySQL only run when their environment variable holds a connection URI, e.g.
//
//	SPACEBIN_TEST_POSTGRES=postgres://user:pass@localhost:5432/spacebin?sslmode=disable
//	SPACEBIN_TEST_MYSQL=mysql://user:pass@tcp(localhost:3306)/spacebin
//
// Both databases must be empty or only used by these tests, as the tests delete all rows.
var testDatabases = []struct {
	name string
	env  string
	open func(string) (Database, error)
}{
	{"sqlite", "", NewSQLite},
	{"postgres", "SPACEBIN_TEST_POSTGRES", NewPostgres},
	{"mysql", "SPACEBIN_TEST_MYSQL", NewMySQL},
}

// eachDatabase runs test against a new connection to every available database.
func eachDatabase(t *testing.T, test func(t *testing.T, db Database)) {
	for _, d := range testDatabases {
		t.Run(d.name, func(t *testing.T) {
			uri := "sqlite://:memory:"

			if d.env != "" {
				if uri = os.Getenv(d.env); uri == "" {
					t.Skipf("%s is not set", d.env)
				}
			}

			db, err := d.open(uri)
			require.NoError(t, err)
			t.Cleanup(func() { db.Close() })

			test(t, db)
		})
	}
}

// eachMigratedDatabase is like eachDatabase, but migrates each database and empties its tables first.
func eachMigratedDatabase(t *testing.T, test func(t *testing.T, db Database)) {
	eachDatabase(t, func(t *testing.T, db Database) {
		ctx := context.Background()
		require.NoError(t, db.Migrate(ctx))

		// Postgres and MySQL keep their data between tests
		for _, table := range []string{"documents", "accounts", "sessions"} {
			_, err := execRaw(db, "DELETE FROM "+table)
			require.NoError(t, err)
		}

		test(t, db)
	})
}

// execRaw runs a query that the Database interface doesn't expose.
func execRaw(db Database, query string) (sql.Result, error) {
	return db.(interface {
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	}).ExecContext(context.Background(), query)
}

// Constructor and migration tests

func TestNewSQLiteFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "spacebin.db")

	// One leading slash is stripped, so an absolute path needs two
	db, err := NewSQLite("sqlite:///" + path)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	require.NoError(t, db.Migrate(context.Background()))
	require.FileExists(t, path)
}

func TestNewSQLiteRelativePath(t *testing.T) {
	// Both forms are relative to the working directory
	for _, uri := range []string{"sqlite://spacebin.db", "sqlite:///spacebin.db"} {
		t.Run(uri, func(t *testing.T) {
			t.Chdir(t.TempDir())

			db, err := NewSQLite(uri)
			require.NoError(t, err)
			t.Cleanup(func() { db.Close() })

			require.NoError(t, db.Migrate(context.Background()))
			require.FileExists(t, "spacebin.db")
		})
	}
}

func TestNewMySQLInvalidDSN(t *testing.T) {
	_, err := NewMySQL("mysql://not-a-dsn")
	require.ErrorContains(t, err, "invalid DSN")
}

func TestMigrateTwice(t *testing.T) {
	eachDatabase(t, func(t *testing.T, db Database) {
		require.NoError(t, db.Migrate(context.Background()))
		require.NoError(t, db.Migrate(context.Background()))
	})
}

// unreachableDatabases can be opened, but fail as soon as they're used.
var unreachableDatabases = []struct {
	name string
	uri  string
	open func(string) (Database, error)
}{
	{"sqlite", "sqlite:////nonexistent/directory/spacebin.db", NewSQLite},
	{"postgres", "postgres://spacebin@127.0.0.1:1/spacebin?sslmode=disable&connect_timeout=1", NewPostgres},
	{"mysql", "mysql://spacebin@tcp(127.0.0.1:1)/spacebin?timeout=1s", NewMySQL},
}

// eachUnreachableDatabase runs test against every unreachable database. Unlike eachDatabase, these don't need a server.
func eachUnreachableDatabase(t *testing.T, test func(t *testing.T, db Database)) {
	for _, d := range unreachableDatabases {
		t.Run(d.name, func(t *testing.T) {
			db, err := d.open(d.uri)
			require.NoError(t, err)
			t.Cleanup(func() { db.Close() })

			test(t, db)
		})
	}
}

func TestMigrateUnreachableDatabase(t *testing.T) {
	eachUnreachableDatabase(t, func(t *testing.T, db Database) {
		require.Error(t, db.Migrate(context.Background()))
	})
}

func newSQLiteMigrationDriver(t *testing.T) *sqlite.Sqlite {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	require.NoError(t, err)

	return driver.(*sqlite.Sqlite)
}

func TestMigrateUpMissingDirectory(t *testing.T) {
	err := migrateUp(fstest.MapFS{}, "migrations/missing", newSQLiteMigrationDriver(t))
	require.Error(t, err)
}

func TestMigrateUpInvalidMigration(t *testing.T) {
	fsys := fstest.MapFS{
		"migrations/broken/1_broken.up.sql": {Data: []byte("NOT VALID SQL;")},
	}

	err := migrateUp(fsys, "migrations/broken", newSQLiteMigrationDriver(t))
	require.ErrorContains(t, err, "NOT VALID SQL")
}

// Document tests

func TestCreateAndGetDocument(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		ctx := context.Background()

		require.NoError(t, db.CreateDocument(ctx, "abc123", "hello world", "alice"))

		doc, err := db.GetDocument(ctx, "abc123")
		require.NoError(t, err)
		require.Equal(t, "abc123", doc.ID)
		require.Equal(t, "hello world", doc.Content)
		require.Equal(t, "alice", doc.Username)
		require.NotZero(t, doc.CreatedAt)
	})
}

func TestCreateDuplicateDocument(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		ctx := context.Background()

		require.NoError(t, db.CreateDocument(ctx, "abc123", "first", ""))
		require.Error(t, db.CreateDocument(ctx, "abc123", "second", ""))
	})
}

func TestGetDocumentNotFound(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		_, err := db.GetDocument(context.Background(), "nope")
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestUpdateDocument(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		ctx := context.Background()

		require.NoError(t, db.CreateDocument(ctx, "doc1", "original", "alice"))
		require.NoError(t, db.UpdateDocument(ctx, "doc1", "updated content"))

		doc, err := db.GetDocument(ctx, "doc1")
		require.NoError(t, err)
		require.Equal(t, "updated content", doc.Content)
	})
}

func TestDeleteDocument(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		ctx := context.Background()

		require.NoError(t, db.CreateDocument(ctx, "doc2", "to delete", "bob"))
		require.NoError(t, db.DeleteDocument(ctx, "doc2"))

		_, err := db.GetDocument(ctx, "doc2")
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestGetDocumentsByUsername(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		ctx := context.Background()

		require.NoError(t, db.CreateDocument(ctx, "d1", "content1", "alice"))
		require.NoError(t, db.CreateDocument(ctx, "d2", "content2", "alice"))
		require.NoError(t, db.CreateDocument(ctx, "d3", "content3", "bob"))

		aliceDocs, err := db.GetDocumentsByUsername(ctx, "alice")
		require.NoError(t, err)
		require.Len(t, aliceDocs, 2)

		bobDocs, err := db.GetDocumentsByUsername(ctx, "bob")
		require.NoError(t, err)
		require.Len(t, bobDocs, 1)
		require.Equal(t, "d3", bobDocs[0].ID)

		noDocs, err := db.GetDocumentsByUsername(ctx, "nobody")
		require.NoError(t, err)
		require.Empty(t, noDocs)
	})
}

func TestGetDocumentsByUsernameScanError(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		// created_at is nullable, but can't be scanned into a time.Time when NULL
		_, err := execRaw(db, "INSERT INTO documents (id, content, username, created_at) VALUES ('broken', 'content', 'alice', NULL)")
		require.NoError(t, err)

		_, err = db.GetDocumentsByUsername(context.Background(), "alice")
		require.Error(t, err)
	})
}

func TestGetDocumentsByUsernameUnreachableDatabase(t *testing.T) {
	eachUnreachableDatabase(t, func(t *testing.T, db Database) {
		_, err := db.GetDocumentsByUsername(context.Background(), "alice")
		require.Error(t, err)
	})
}

// Account tests

func TestCreateAndGetAccount(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		ctx := context.Background()

		require.NoError(t, db.CreateAccount(ctx, "alice", "password123"))

		acc, err := db.GetAccountByUsername(ctx, "alice")
		require.NoError(t, err)
		require.Equal(t, "alice", acc.Username)
		require.NoError(t, bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte("password123")))

		acc2, err := db.GetAccount(ctx, fmt.Sprint(acc.ID))
		require.NoError(t, err)
		require.Equal(t, acc, acc2)
	})
}

func TestCreateAccountPasswordTooLong(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		err := db.CreateAccount(context.Background(), "alice", strings.Repeat("a", 73))
		require.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
	})
}

func TestGetAccountNotFound(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		_, err := db.GetAccountByUsername(context.Background(), "ghost")
		require.ErrorIs(t, err, sql.ErrNoRows)

		_, err = db.GetAccount(context.Background(), "12345")
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestDeleteAccount(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		ctx := context.Background()

		require.NoError(t, db.CreateAccount(ctx, "carol", "pw"))

		acc, err := db.GetAccountByUsername(ctx, "carol")
		require.NoError(t, err)

		require.NoError(t, db.DeleteAccount(ctx, fmt.Sprint(acc.ID)))

		_, err = db.GetAccountByUsername(ctx, "carol")
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

// Session tests

func TestCreateAndGetSession(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		ctx := context.Background()

		require.NoError(t, db.CreateSession(ctx, "pub1", "tok1", "sec1", "alice"))

		sess, err := db.GetSession(ctx, "pub1")
		require.NoError(t, err)
		require.Equal(t, Session{Public: "pub1", Token: "tok1", Secret: "sec1", Username: "alice"}, sess)
	})
}

func TestGetSessionNotFound(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		_, err := db.GetSession(context.Background(), "nosuchsession")
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestDeleteSession(t *testing.T) {
	eachMigratedDatabase(t, func(t *testing.T, db Database) {
		ctx := context.Background()

		require.NoError(t, db.CreateSession(ctx, "pub2", "tok2", "sec2", "bob"))
		require.NoError(t, db.DeleteSession(ctx, "pub2"))

		_, err := db.GetSession(ctx, "pub2")
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}
