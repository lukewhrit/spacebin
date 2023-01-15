package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/orca-group/spirit/internal/config"
	"github.com/stretchr/testify/require"
)

type Response struct {
	Payload config.Cfg
	Error   string
}

var mockConfig = config.Cfg{
	Host:             "0.0.0.0",
	Port:             9000,
	CompressionLevel: 1,
	Ratelimiter:      "200x5",
	IDLength:         8,
	IDType:           "key",
	MaxSize:          400_000,
	ExpirationAge:    720,
}

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

// checkResponseCode is a simple utility to check the response code
// of the response
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestConfig(t *testing.T) {
	s := NewServer(&mockConfig)
	s.MountHandlers()

	req, _ := http.NewRequest("GET", "/config", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)

	x, _ := io.ReadAll(res.Result().Body)
	var body Response
	json.Unmarshal(x, &body)

	require.Equal(t, mockConfig, body.Payload)
}
