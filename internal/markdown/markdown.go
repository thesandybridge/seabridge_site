package markdown

import (
	"bytes"
	"os"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func ConvertToHTML(filePath string) (string, error) {
	markdownContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("gruvbox"),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := markdown.Convert(markdownContent, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
