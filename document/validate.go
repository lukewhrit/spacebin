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
	/*
	 * This regex matches the file extension for various languages.

	 * Languages including:
	 * 	Python, JavaScript(React), TypeScript(React), Go, Kotlin,
	 *	C++, SQL, C-Sharp, C (including `.h`), Scala, Haskell,
	 * 	Shell, PowerShell, PHP, Assembly, Julia, Objective-C (.m),
	 *	Perl, Crystal, JSON, YAML, TOML, and Plain Text

	 * For any unsupported formats Plain Text should be used.
	 */
	regex := regexp.MustCompile("^py$|^[tj](sx|s)$|^go$|^kt$|^java$|^c$|^c(pp|[rs])$|^sql$|^swift$|^dart$|^r$|^r[sb]$|^h$|^scala$|^hs$|^sh$|^p(s1|hp)$|^asm$|^jl$|^m$|^txt$|^pl$|^(x|ya|to)ml$")

	return validation.ValidateStruct(&c,
		validation.Field(
			&c.Content,
			validation.Required,
			// Enforce length to follow what's set in the config
			validation.Length(2, config.Config.Documents.MaxDocumentLength),
		),
		// The purpose of this field is to support client's that perform
		// syntax highlighting and need to know what highlighter to use.
		validation.Field(
			&c.Extension,
			validation.Match(regex),
			validation.Required,
		),
	)
}
