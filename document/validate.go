package document

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/spacebin-org/curiosity/config"
)

// CreateRequest represents a valid body object for the create document request
type CreateRequest struct {
	Content   string
	Extension string
}

// Validate performs validation on the body
func (c CreateRequest) Validate() error {
	regex := regexp.MustCompile("^txt$|^js$")

	return validation.ValidateStruct(&c,
		validation.Field(
			&c.Content,
			validation.Required,
			// Enforce length to follow
			validation.Length(2, config.Config.Documents.MaxDocumentLength),
		),
		validation.Field(
			&c.Extension,
			validation.Match(regex),
		),
	)
}
