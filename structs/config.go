package structs

// Config is the configuration object
type Config struct {
	Server struct {
		Host           string
		Port           int
		UseCSP         bool
		CompresssLevel int
		Prefork        bool

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
