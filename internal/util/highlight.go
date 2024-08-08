package util

import (
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// Highlight uses Chroma to highlight code in documents.
func Highlight(code string, extension string) (string, string, error) {
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

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(html.WithLineNumbers(true), html.WithLinkableLineNumbers(true, "L"), html.WithClasses(true))
	iterator, err := lexer.Tokenise(nil, code)

	if err != nil {
		return "", "", err
	}

	w := new(strings.Builder)
	err = formatter.Format(w, style, iterator)

	if err != nil {
		return "", "", err
	}

	css := new(strings.Builder)
	err = formatter.WriteCSS(css, style)

	if err != nil {
		return "", "", err
	}

	return w.String(), css.String(), nil
}
