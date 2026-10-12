package util

import (
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// Highlight uses Chroma to highlight code in documents.
func Highlight(code string, extension string) (string, string) {
	// This function is used to highlight code in documents.
	// It uses the Chroma library to parse and highlight code.
	// The Chroma lexer is determined by the document's language.
	// If the document does not have a language, the lexer is set to plaintext.
	// The highlighted code is then returned as a string containing HTML.

	var lexer chroma.Lexer

	if extension != "" {
		lexer = lexers.Get(extension)
	} else {
		lexer = lexers.Analyse(code)
	}

	if lexer == nil {
		lexer = lexers.Fallback
	}

	return highlightWith(lexer, code)
}

// highlightWith highlights code using lexer, falling back to plaintext if the lexer fails.
func highlightWith(lexer chroma.Lexer, code string) (string, string) {
	// Get falls back to a default style itself, so it never returns nil
	style := styles.Get("monokai")

	formatter := html.New(html.WithLineNumbers(true), html.WithLinkableLineNumbers(true, "L"), html.WithClasses(true))
	iterator, err := lexer.Tokenise(nil, code)

	if err != nil {
		// The plaintext lexer matches any input, so it can't fail
		iterator, _ = lexers.Fallback.Tokenise(nil, code)
	}

	// Writing to a strings.Builder never fails, so neither can formatting
	w := new(strings.Builder)
	formatter.Format(w, style, iterator)

	css := new(strings.Builder)
	formatter.WriteCSS(css, style)

	return w.String(), css.String()
}
