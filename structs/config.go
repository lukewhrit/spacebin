package structs

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
