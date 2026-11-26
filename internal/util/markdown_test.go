/*
* Copyright 2020-2026 Luke Whritenour

* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at

*     http://www.apache.org/licenses/LICENSE-2.0

* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package util_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/lukewhrit/spacebin/internal/util"
	"github.com/stretchr/testify/require"
)

func TestParseMarkdown(t *testing.T) {
	// Test basic markdown
	input := []byte("# Title\n\nThis is a paragraph.")
	output := util.ParseMarkdown(input)

	require.NotEmpty(t, output)
	require.Contains(t, string(output), "<h1")
	require.Contains(t, string(output), "Title")
	require.Contains(t, string(output), "<p>")
	require.Contains(t, string(output), "This is a paragraph.")
}

func TestParseMarkdownWithLinks(t *testing.T) {
	// Test markdown with links
	input := []byte("[Link](https://example.com)")
	output := util.ParseMarkdown(input)

	require.NotEmpty(t, output)
	require.Contains(t, string(output), "<a")
	require.Contains(t, string(output), "https://example.com")
	require.Contains(t, string(output), "Link")
}

func TestParseMarkdownEmpty(t *testing.T) {
	// Test empty markdown
	input := []byte("")
	output := util.ParseMarkdown(input)

	// Empty input might return nil or empty slice, both are valid
	_ = output
}

// Regression tests for #592: reader mode must not render links or images whose
// destination can run script.
func TestParseMarkdownBlocksUnsafeURLs(t *testing.T) {
	tests := map[string]string{
		"javascript link":        "[click](javascript:alert(1))",
		"mixed case scheme":      "[click](JaVaScRiPt:alert(1))",
		"leading space":          "[click]( javascript:alert(1))",
		"tab inside scheme":      "[click](<java\tscript:alert(1)>)",
		"entity encoded scheme":  "[click](&#106;avascript:alert(1))",
		"entity encoded colon":   "[click](javascript&colon;alert(1))",
		"vbscript link":          "[click](vbscript:msgbox(1))",
		"data html link":         "[click](data:text/html,<script>alert(1)</script>)",
		"reference link":         "[click][r]\n\n[r]: javascript:alert(1)",
		"autolink":               "<javascript:alert(1)>",
		"protocol relative link": "[click](//evil.example)",
		"javascript image":       "![img](javascript:alert(1))",
		"svg data image":         "![img](data:image/svg+xml,<svg onload=alert(1)>)",
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			output := strings.ToLower(string(util.ParseMarkdown([]byte(input))))

			require.NotContains(t, output, "href=")
			require.NotContains(t, output, "src=")
			require.NotContains(t, output, "onclick=")
			require.NotContains(t, output, "onload=")
		})
	}
}

func TestParseMarkdownKeepsSafeURLs(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"https link":     {"[a](https://example.com)", `href="https://example.com"`},
		"http link":      {"[a](http://example.com)", `href="http://example.com"`},
		"mailto link":    {"[a](mailto:me@example.com)", `href="mailto:me@example.com"`},
		"relative link":  {"[a](/about)", `href="/about"`},
		"dot relative":   {"[a](./about)", `href="./about"`},
		"fragment link":  {"[a](#section)", `href="#section"`},
		"https image":    {"![a](https://example.com/a.png)", `src="https://example.com/a.png"`},
		"bare image":     {"![a](a.png)", `src="a.png"`},
		"png data image": {"![a](data:image/png;base64,AAAA)", `src="data:image/png;base64,AAAA"`},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Contains(t, string(util.ParseMarkdown([]byte(tc.input))), tc.want)
		})
	}
}

// eventHandlerAttr matches an on* attribute inside an HTML tag.
var eventHandlerAttr = regexp.MustCompile(`(?i)<[^>]*\son[a-z]+\s*=`)

func TestParseMarkdownDropsEventHandlerAttributes(t *testing.T) {
	inputs := []string{
		"[click](https://example.com){onclick=alert(1)}",
		"{onclick=\"alert(1)\" #id}\n# Heading",
		"{onmouseover=\"alert(1)\"}\nparagraph",
	}

	for _, input := range inputs {
		output := string(util.ParseMarkdown([]byte(input)))

		require.NotRegexp(t, eventHandlerAttr, output, input)
	}
}
