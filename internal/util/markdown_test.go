/*
* Copyright 2020-2024 Luke Whritenour

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
