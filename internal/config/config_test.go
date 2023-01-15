/*
 * Copyright 2020-2023 Luke Whritenour, Jack Dorland

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
