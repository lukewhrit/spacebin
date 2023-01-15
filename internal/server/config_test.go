package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/orca-group/spirit/internal/config"
	"github.com/stretchr/testify/require"
)

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
	s := NewServer()
	s.MountHandlers()

	req, _ := http.NewRequest("GET", "/config", nil)
	res := executeRequest(req, s)

	checkResponseCode(t, http.StatusOK, res.Result().StatusCode)

	var body config.Cfg
	json.Unmarshal(res.Body.Bytes(), &body)

	require.Equal(t, config.Cfg{
		Host:             config.Config.Host,
		Port:             config.Config.Port,
		CompressionLevel: config.Config.CompressionLevel,
		Ratelimiter:      config.Config.Ratelimiter,
		IDLength:         config.Config.IDLength,
		IDType:           config.Config.IDType,
		MaxSize:          config.Config.MaxSize,
		ExpirationAge:    config.Config.ExpirationAge,
	}, body)
}
