/*
 * Copyright 2020-2022 Luke Whritenour, Jack Dorland

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
	env "github.com/caarlos0/env/v6"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

// Config is the loaded config object
var Config struct {
	// General
	Host             string         `env:"HOST" envDefault:"0.0.0.0"`
	Port             int            `env:"PORT" envDefault:"9000"`
	CompressionLevel compress.Level `env:"COMPRESS_LEVEL" envDefault:"1"`
	Ratelimiter      string         `env:"RATELIMITER" envDefault:"200x5"` // Requests x Seconds
	ConnectionURI    string         `env:"CONNECTION_URI"`

	// Document
	IDLength      int   `env:"ID_LENGTH" envDefault:"8"`
	MaxSize       int   `env:"MAX_SIZE" envDefault:"400000"`    // in bytes
	ExpirationAge int64 `env:"EXPIRATION_AGE" envDefault:"720"` // in hours
}

// Load configuration from file
func Load() error {
	return env.Parse(&Config, env.Options{
		Prefix:          "SPIRIT_",
		RequiredIfNoDef: true,
	})
}
