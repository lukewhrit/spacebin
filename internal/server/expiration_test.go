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

package server_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lukewhrit/spacebin/internal/database/databasefakes"
	"github.com/lukewhrit/spacebin/internal/server"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

// captureLogs collects everything logged until the test ends.
func captureLogs(t *testing.T) *bytes.Buffer {
	t.Helper()

	var buf bytes.Buffer
	old := log.Logger
	log.Logger = log.Output(&buf)
	t.Cleanup(func() { log.Logger = old })

	return &buf
}

func expiringServer(db *databasefakes.FakeDatabase, age int64) *server.Server {
	cfg := mockConfig
	cfg.ExpirationAge = age
	cfg.Documents = []string{"about"}
	return server.NewServer(&cfg, db)
}

func TestExpireDocumentsDisabled(t *testing.T) {
	db := &databasefakes.FakeDatabase{}

	// Returns without waiting for ctx
	expiringServer(db, 0).ExpireDocuments(context.Background(), time.Millisecond)

	require.Zero(t, db.DeleteExpiredDocumentsCallCount())
}

func TestExpireDocumentsRunsUntilCancelled(t *testing.T) {
	logs := captureLogs(t)

	db := &databasefakes.FakeDatabase{}
	db.DeleteExpiredDocumentsReturns(3, nil)
	s := expiringServer(db, 24)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		s.ExpireDocuments(ctx, time.Millisecond)
		close(done)
	}()

	// Runs straight away, then again every interval
	require.Eventually(t, func() bool { return db.DeleteExpiredDocumentsCallCount() >= 2 }, 5*time.Second, time.Millisecond)

	cancel()
	<-done

	_, age, keep := db.DeleteExpiredDocumentsArgsForCall(0)
	require.EqualValues(t, 24, age)
	require.Equal(t, []string{"about"}, keep)
	require.Contains(t, logs.String(), `"deleted":3`)
}

func TestExpireDocumentsLogsErrors(t *testing.T) {
	logs := captureLogs(t)

	db := &databasefakes.FakeDatabase{}
	db.DeleteExpiredDocumentsReturns(0, errors.New("db down"))
	s := expiringServer(db, 24)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		s.ExpireDocuments(ctx, time.Hour)
		close(done)
	}()

	require.Eventually(t, func() bool { return db.DeleteExpiredDocumentsCallCount() == 1 }, 5*time.Second, time.Millisecond)

	cancel()
	<-done

	require.Contains(t, logs.String(), "Error deleting expired documents")
}

func TestExpireDocumentsQuietWhenCancelled(t *testing.T) {
	logs := captureLogs(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	db := &databasefakes.FakeDatabase{}
	db.DeleteExpiredDocumentsReturns(0, context.Canceled)

	expiringServer(db, 24).ExpireDocuments(ctx, time.Hour)

	require.Equal(t, 1, db.DeleteExpiredDocumentsCallCount())
	require.NotContains(t, logs.String(), "Error deleting expired documents")
}
