package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLoad runs the Load() function, ensures it doesn't error and populates the Config struct
func TestLoad(t *testing.T) {
	t.Setenv("SPIRIT_CONNECTION_URI", "host=localhost port=5432 user=spacebin database=spacebin sslmode=disable")

	if Load() != nil {
		t.Fail()
	}

	require.EqualValues(t, Config, Cfg{
		Host:             "0.0.0.0",
		Port:             9000,
		CompressionLevel: 1,
		Ratelimiter:      "200x5",
		IDLength:         8,
		IDType:           "key",
		MaxSize:          400_000,
		ConnectionURI:    "host=localhost port=5432 user=spacebin database=spacebin sslmode=disable",
		ExpirationAge:    720,
	})
}
