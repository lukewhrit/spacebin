package util_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lukewhrit/spacebin/internal/util"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	old := log.Logger
	log.Logger = log.Output(&buf)
	t.Cleanup(func() { log.Logger = old })

	handler := util.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("hello"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/some/page", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusTeapot, rr.Code)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	require.Equal(t, "HTTP Request", entry["message"])
	require.Equal(t, "GET", entry["method"])
	require.Equal(t, "/some/page", entry["page"])
	require.EqualValues(t, http.StatusTeapot, entry["status"])
	require.EqualValues(t, 5, entry["size"])
}
