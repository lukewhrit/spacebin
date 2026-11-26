package util

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// isSafeURL allows the URLs gomarkdown's Safelink accepts (http, https, ftp,
// mailto and relative paths) plus same-document fragment links such as
// "#heading", which can never carry a scheme.
func isSafeURL(dest []byte) bool {
	return (len(dest) > 0 && dest[0] == '#') || parser.IsSafeURL(dest)
}

func ParseMarkdown(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions. SkipHTML drops raw HTML and
	// Safelink drops links and images with untrusted schemes (javascript:,
	// data:, vbscript:, ...), since pastes are untrusted input.
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.SkipHTML | html.Safelink
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	renderer.IsSafeURLOverride = isSafeURL

	return markdown.Render(doc, renderer)
}
