package config

import (
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/spacebin-org/curiosity/structs"
)

var k = koanf.New(".")

var configuration structs.Config

// Load configuration from file
func Load() error {
	// Set some default values
	k.Load(confmap.Provider(map[string]interface{}{
		"Server.Host":                 "0.0.0.0",
		"Server.Port":                 9000,
		"Server.CompressionLevel":     -1,
		"Server.EnableCSP":            true,
		"Server.Ratelimits.Requests":  80,
		"Server.Ratelimits.Duration":  60_000,
		"Documents.IDLength":          8,
		"Documents.MaxDocumentLength": 400_000,
		"Database.Dialect":            "sqlite3",
		"Database.ConnectionURI":      "spacebin.db",
	}, "."), nil)

	f := file.Provider("./config.toml")

	// Load configuration from JSON on top of said default values
	if err := k.Load(f, toml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	k.Unmarshal("", &configuration)

	return nil
}

// GetConfig returns the entire configuration object
func GetConfig() structs.Config {
	return configuration
}
