package util

import (
	"testing"
)

func TestHighlight(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		extension   string
		wantHTML    string
		wantCSS     string
		expectError bool
	}{
		{
			name:        "Valid Go Code",
			code:        "package main\nfunc main() {}",
			extension:   "go",
			expectError: false,
		},
		{
			name:        "Valid Python Code",
			code:        "def main():\n    pass",
			extension:   "py",
			expectError: false,
		},
		{
			name:        "Plaintext Code",
			code:        "Just some text.",
			extension:   "",
			expectError: false,
		},
		{
			name:        "Invalid Extension",
			code:        "Invalid extension test",
			extension:   "invalid",
			expectError: false, // Lexer should fallback
		},
		{
			name:        "Empty Code",
			code:        "",
			extension:   "",
			expectError: false,
		},
		{
			name:        "Code with no extension - lexer analyse",
			code:        "console.log('test');",
			extension:   "",
			expectError: false,
		},
		{
			name:        "Extremely long extension that doesn't exist",
			code:        "test content",
			extension:   "thisdoesnotexistatall123456789",
			expectError: false, // Should fallback to default lexer
		},
		{
			name:        "Various programming languages",
			code:        "import java.util.*;",
			extension:   "java",
			expectError: false,
		},
		{
			name:        "C code",
			code:        "#include <stdio.h>\nint main() { return 0; }",
			extension:   "c",
			expectError: false,
		},
		{
			name:        "JavaScript code",
			code:        "function test() { return true; }",
			extension:   "js",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHTML, gotCSS, err := Highlight(tt.code, tt.extension)
			if (err != nil) != tt.expectError {
				t.Errorf("Highlight() error = %v, wantErr %v", err, tt.expectError)
				return
			}

			// Since the output is dependent on the Chroma library, it's tricky to hardcode expected values.
			// We can, however, validate that the result is not empty when it should be valid.
			if !tt.expectError && (gotHTML == "" || gotCSS == "") {
				t.Errorf("Expected non-empty HTML and CSS, gotHTML = %v, gotCSS = %v", gotHTML, gotCSS)
			}
		})
	}
}

// TestHighlightNilLexer tests the fallback when lexer is nil
func TestHighlightNilLexer(t *testing.T) {
	// Test with an extension that doesn't exist to trigger nil lexer
	html, css, err := Highlight("some random text", "nonexistentextension")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if html == "" || css == "" {
		t.Error("Expected non-empty output even with nil lexer (should use fallback)")
	}
}

// TestHighlightWithAnalyse tests code analysis without extension
func TestHighlightWithAnalyse(t *testing.T) {
	// Test various code snippets to ensure Analyse path is covered
	testCases := []string{
		"def foo(): pass",           // Python
		"function test() {}",        // JavaScript
		"<html><body></body></html>", // HTML
		"SELECT * FROM users;",      // SQL
	}
	
	for _, code := range testCases {
		html, css, err := Highlight(code, "")
		if err != nil {
			t.Errorf("Unexpected error for code %q: %v", code, err)
		}
		if html == "" || css == "" {
			t.Errorf("Expected non-empty output for code %q", code)
		}
	}
}
