package structs

// Config is the configuration object
type Config struct {
	Server struct {
		Host              string `koanf:"host"`
		Port              int    `koanf:"port"`
		UseCSP            bool   `koanf:"useCSP"`
		CompresssionLevel int    `koanf:"compressionLevel"`
		Prefork           bool   `koanf:"prefork"`

		Ratelimits struct {
			Requests int `koanf:"requests"`
			Duration int `koanf:"duration"`
		} `koanf:"ratelimits"`
	}

	Documents struct {
		IDLength          int `koanf:"idLength"`
		MaxDocumentLength int `koanf:"maxDocunentLength"`
	} `koanf:"documents"`

	Database struct {
		Dialect       string `koanf:"dialect"`
		ConnectionURI string `koanf:"connection_uri"`
	} `koanf:"database"`
}
