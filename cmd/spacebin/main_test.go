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

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/lukewhrit/spacebin/internal/database"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

// freePort returns a port that nothing is listening on.
func freePort(t *testing.T) int {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}

// setTestEnv configures Spirit to listen on port with an in-memory database.
func setTestEnv(t *testing.T, port int) {
	t.Setenv("SPIRIT_HOST", "127.0.0.1")
	t.Setenv("SPIRIT_PORT", strconv.Itoa(port))
	t.Setenv("SPIRIT_CONNECTION_URI", "sqlite://:memory:")
}

// get requests path from the server listening on port.
func get(port int, path string) (int, error) {
	res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d%s", port, path))
	if err != nil {
		return 0, err
	}

	res.Body.Close()

	return res.StatusCode, nil
}

// waitForServer blocks until the server listening on port responds.
func waitForServer(t *testing.T, port int) {
	t.Helper()

	require.Eventually(t, func() bool {
		status, err := get(port, "/ping")
		return err == nil && status == http.StatusOK
	}, 10*time.Second, 20*time.Millisecond)
}

func TestRun(t *testing.T) {
	for _, headless := range []bool{false, true} {
		t.Run(fmt.Sprintf("headless=%t", headless), func(t *testing.T) {
			port := freePort(t)
			setTestEnv(t, port)
			t.Setenv("SPIRIT_HEADLESS", strconv.FormatBool(headless))

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			done := make(chan error, 1)
			go func() { done <- run(ctx) }()

			waitForServer(t, port)

			// The homepage is only served when not headless
			status, err := get(port, "/")
			require.NoError(t, err)
			require.Equal(t, !headless, status == http.StatusOK)

			cancel()
			require.NoError(t, <-done)
		})
	}
}

func TestRunErrors(t *testing.T) {
	// Keep a port busy so the HTTP listener can't use it
	busy, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer busy.Close()

	tests := []struct {
		name string
		env  map[string]string
		err  string
	}{
		{"invalid config", map[string]string{"SPIRIT_PORT": "not-a-number"}, "could not load config"},
		{"unsupported database", map[string]string{"SPIRIT_CONNECTION_URI": "redis://localhost"}, `unsupported database type "redis"`},
		{"invalid database URI", map[string]string{"SPIRIT_CONNECTION_URI": "mysql://not-a-dsn"}, "could not connect to database"},
		{"unreachable database", map[string]string{"SPIRIT_CONNECTION_URI": "sqlite:////nonexistent/directory/spacebin.db"}, "failed migrations"},
		{"port in use", map[string]string{"SPIRIT_PORT": strconv.Itoa(busy.Addr().(*net.TCPAddr).Port)}, "failed to start HTTP listener"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTestEnv(t, freePort(t))

			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			require.ErrorContains(t, run(ctx), tt.err)
		})
	}
}

func TestOpenDatabase(t *testing.T) {
	tests := []struct {
		uri  string
		want database.Database
	}{
		{"sqlite://:memory:", &database.SQLite{}},
		{"file://spacebin.db", &database.SQLite{}},
		{"postgres://127.0.0.1:1/spacebin", &database.Postgres{}},
		{"postgresql://127.0.0.1:1/spacebin", &database.Postgres{}},
		{"mysql://spacebin@tcp(127.0.0.1:1)/spacebin", &database.MySQL{}},
		{"mariadb://spacebin@tcp(127.0.0.1:1)/spacebin", &database.MySQL{}},
	}

	for _, tt := range tests {
		t.Run(tt.uri, func(t *testing.T) {
			db, err := openDatabase(tt.uri)
			require.NoError(t, err)
			t.Cleanup(func() { db.Close() })

			require.IsType(t, tt.want, db)
		})
	}

	_, err := openDatabase("localhost:5432")
	require.ErrorContains(t, err, "unsupported database type")
}

// restoreLogger undoes main's changes to the global logger after the test.
func restoreLogger(t *testing.T) {
	old := log.Logger
	t.Cleanup(func() { log.Logger = old })
}

func TestMainShutsDownOnSignal(t *testing.T) {
	restoreLogger(t)

	port := freePort(t)
	setTestEnv(t, port)

	done := make(chan struct{})
	go func() {
		main()
		close(done)
	}()

	waitForServer(t, port)

	// main handles the signal, so it doesn't stop the test process
	process, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	require.NoError(t, process.Signal(os.Interrupt))

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("main did not return after receiving a signal")
	}
}

func TestMainExitsOnError(t *testing.T) {
	restoreLogger(t)
	t.Setenv("SPIRIT_PORT", "not-a-number")

	code := 0
	old := exit
	exit = func(c int) { code = c }
	t.Cleanup(func() { exit = old })

	main()

	require.Equal(t, 1, code)
}
