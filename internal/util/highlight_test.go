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
