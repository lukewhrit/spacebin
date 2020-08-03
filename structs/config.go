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
}
