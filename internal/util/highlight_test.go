package util

import (
	"errors"
	"strings"
	"testing"

	"github.com/alecthomas/chroma/v2"
)

func TestHighlight(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		extension string
		wantHTML  string
		wantCSS   string
	}{
		{
			name:      "Valid Go Code",
			code:      "package main\nfunc main() {}",
			extension: "go",
		},
		{
			name:      "Valid Python Code",
			code:      "def main():\n    pass",
			extension: "py",
		},
		{
			name:      "Plaintext Code",
			code:      "Just some text.",
			extension: "",
		},
		{
			name:      "Invalid Extension",
			code:      "Invalid extension test",
			extension: "invalid",
		},
		{
			name:      "Empty Code",
			code:      "",
			extension: "",
		},
		{
			name:      "Code with no extension - lexer analyse",
			code:      "console.log('test');",
			extension: "",
		},
		{
			name:      "Extremely long extension that doesn't exist",
			code:      "test content",
			extension: "thisdoesnotexistatall123456789",
		},
		{
			name:      "Various programming languages",
			code:      "import java.util.*;",
			extension: "java",
		},
		{
			name:      "C code",
			code:      "#include <stdio.h>\nint main() { return 0; }",
			extension: "c",
		},
		{
			name:      "JavaScript code",
			code:      "function test() { return true; }",
			extension: "js",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHTML, gotCSS := Highlight(tt.code, tt.extension)

			// Since the output is dependent on the Chroma library, it's tricky to hardcode expected values.
			// We can, however, validate that the result is not empty.
			if gotHTML == "" || gotCSS == "" {
				t.Errorf("Expected non-empty HTML and CSS, gotHTML = %v, gotCSS = %v", gotHTML, gotCSS)
			}
		})
	}
}

// TestHighlightNilLexer tests the fallback when lexer is nil
func TestHighlightNilLexer(t *testing.T) {
	// Test with an extension that doesn't exist to trigger nil lexer
	html, css := Highlight("some random text", "nonexistentextension")
	if html == "" || css == "" {
		t.Error("Expected non-empty output even with nil lexer (should use fallback)")
	}
}

// TestHighlightWithAnalyse tests code analysis without extension
func TestHighlightWithAnalyse(t *testing.T) {
	// Test various code snippets to ensure Analyse path is covered
	testCases := []string{
		"def foo(): pass",            // Python
		"function test() {}",         // JavaScript
		"<html><body></body></html>", // HTML
		"SELECT * FROM users;",       // SQL
	}

	for _, code := range testCases {
		html, css := Highlight(code, "")
		if html == "" || css == "" {
			t.Errorf("Expected non-empty output for code %q", code)
		}
	}
}

type failingLexer struct{ chroma.Lexer }

func (failingLexer) Tokenise(*chroma.TokeniseOptions, string) (chroma.Iterator, error) {
	return nil, errors.New("tokenise error")
}

// TestHighlightLexerError tests the plaintext fallback when a lexer fails
func TestHighlightLexerError(t *testing.T) {
	html, css := highlightWith(failingLexer{}, "<some code>")

	if !strings.Contains(html, "&lt;some code&gt;") {
		t.Errorf("Expected escaped plaintext output, got %q", html)
	}
	if css == "" {
		t.Error("Expected non-empty CSS")
	}
}
