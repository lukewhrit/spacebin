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
	env "github.com/caarlos0/env/v6"
)

type Cfg struct {
	// General
	Host             string `env:"HOST" envDefault:"0.0.0.0" json:"host"`
	Port             int    `env:"PORT" envDefault:"9000" json:"port"`
	CompressionLevel int    `env:"COMPRESS_LEVEL" envDefault:"1" json:"compression_level"`
	Ratelimiter      string `env:"RATELIMITER" envDefault:"200x5" json:"ratelimiter"` // Requests x Seconds
	ConnectionURI    string `env:"CONNECTION_URI" json:"-"`

	// Document
	IDLength      int      `env:"ID_LENGTH" envDefault:"8" json:"id_length"`
	IDType        string   `env:"ID_TYPE" envDefault:"key" json:"id_type"`
	MaxSize       int      `env:"MAX_SIZE" envDefault:"400000" json:"max_size"`          // in bytes
	ExpirationAge int64    `env:"EXPIRATION_AGE" envDefault:"720" json:"expiration_age"` // in hours
	Documents     []string `env:"DOCUMENTS" envDefault:"" json:"documents"`
}

// Config is the loaded config object
var Config Cfg

// Load configuration from file
func Load() error {
	return env.Parse(&Config, env.Options{
		Prefix:          "SPIRIT_",
		RequiredIfNoDef: true,
	})
}
