package config

import (
	"github.com/spf13/viper"
)

var configuration *config

type config struct {
	Server struct {
		Port           int
		UseCSP         bool
		CompresssLevel int

		Ratelimits struct {
			Requests int
			Duration int
		}
	}

	Documents struct {
		IDLength          int
		MaxDocumentLength int
	}
}

// RatelimitsStruct contains values for ratelimiting configuration
type RatelimitsStruct struct {
	Requests int
	Duration int
}

// Load configuration from file
func Load() error {
	c := new(config)

	// Set defaults
	viper.SetDefault("server.Port", 77223)
	viper.SetDefault("server.UseCSP", true)
	viper.SetDefault("server.CompressLevel", 1)

	viper.SetDefault("server.ratelimits.requests", 500)
	viper.SetDefault("server.ratelimits.duration", 60000)

	viper.SetDefault("documents.IDLength", 12)
	viper.SetDefault("documents.MaxDocumentLength", 400000)

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

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

// Port returns the port for the server to listen on
func Port() int {
	return configuration.Server.Port
}

// UseCSP returns a boolean indicating whether to use CSP or not
func UseCSP() bool {
	return configuration.Server.UseCSP
}

// CompressLevel returns the level of compression to use
func CompressLevel() int {
	return configuration.Server.CompresssLevel
}

// Ratelimits returns the ratelimits object from the config
func Ratelimits() RatelimitsStruct {
	return configuration.Server.Ratelimits
}
