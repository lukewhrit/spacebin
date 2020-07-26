package config

import (
	"github.com/spf13/viper"
)

// Ratelimits contains values for ratelimiting configuration
type Ratelimits struct {
	Requests int
	Duration int
}

// Documents hold values related to document IDs
type Documents struct {
	IDLength          int
	MaxDocumentLength int
}

// Database holds the required information for connecting to a database via Gorm
type Database struct {
	Dialect       string
	ConnectionURI string
}

// Config is the configuration object
type Config struct {
	Server struct {
		Host           string
		Port           int
		UseCSP         bool
		CompresssLevel int

		Ratelimits
	}

	Documents

	Database
}

var configuration *Config

// Load configuration from file
func Load() error {
	c := new(Config)

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	// Set default values for server block
	viper.SetDefault("server.Host", "0.0.0.0")
	viper.SetDefault("server.Port", 9000)
	viper.SetDefault("server.UseCSP", true)
	viper.SetDefault("server.CompressLevel", 0)

	// Set default values for server ratelimits block
	viper.SetDefault("server.ratelimits.requests", 80)
	viper.SetDefault("server.ratelimits.duration", 60000)

	// Set default values for documents block
	viper.SetDefault("documents.IDLength", 12)
	viper.SetDefault("documents.MaxDocumentLength", 400000)

	// Set default values for database block
	viper.SetDefault("database.Dialect", "sqlite3")
	viper.SetDefault("database.ConnectionURI", "spacebin.db")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return err
	}

	configuration = c

	return nil
}

// GetPort returns the port for the server to listen on
func GetPort() int {
	return configuration.Server.Port
}

// GetHost returns the host for the server to listen on
func GetHost() string {
	return configuration.Server.Host
}

// GetUseCSP returns a boolean indicating whether to use CSP or not
func GetUseCSP() bool {
	return configuration.Server.UseCSP
}

// GetCompressLevel returns the level of compression to use
func GetCompressLevel() int {
	return configuration.Server.CompresssLevel
}

// GetRatelimits returns the ratelimits object from the config
func GetRatelimits() Ratelimits {
	return configuration.Server.Ratelimits
}

// GetDatabase returns information for connecting to a database
func GetDatabase() Database {
	return configuration.Database
}
