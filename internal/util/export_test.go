package util

import (
	"io"
	"testing"
)

// SetReadAll replaces the function used to read uploaded files for the duration of a test.
func SetReadAll(t testing.TB, f func(io.Reader) ([]byte, error)) {
	old := readAll
	readAll = f
	t.Cleanup(func() { readAll = old })
}
