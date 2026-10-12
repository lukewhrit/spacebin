package server

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

type failingWriter struct{ *httptest.ResponseRecorder }

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("write error") }

func TestRenderWriteError(t *testing.T) {
	var buf bytes.Buffer
	old := log.Logger
	log.Logger = log.Output(&buf)
	t.Cleanup(func() { log.Logger = old })

	render(failingWriter{httptest.NewRecorder()}, "index.html", map[string]any{})

	require.Contains(t, buf.String(), "write error")
	require.Contains(t, buf.String(), "index.html")
}
