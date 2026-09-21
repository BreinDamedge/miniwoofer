package impl

import (
	"io"
	"io/fs"
	"net/http"

	"bytes"

	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

func md_to_html(x io.Reader) ([]byte, error) {
	// take in the markdown io.Reader
	var buf bytes.Buffer

	file_bytes, err := io.ReadAll(x) // probably a better way to do this
	if err != nil {
		return buf.Bytes(), err
	}

	// convert the markdown to html
	p := parser.New(parser.WithAttribute(), parser.WithExtensions(extension.GFMParser))
	r := html.New(html.WithXHTML(), html.WithUnsafe(), html.WithExtensions(extension.GFMHTMLRenderer))

	doc := p.Parse(file_bytes)
	if err := r.Render(&buf, file_bytes, doc); err != nil {
		return buf.Bytes(), err
	}

	// success
	return buf.Bytes(), nil
}

type MarkdownHandler struct{}

func (_ MarkdownHandler) Tokenize(x io.Reader) ([]string, error) {
	// parse md as raw text
	return TokenizePlaintext(x)
}

func (_ MarkdownHandler) ExtractTitle(x io.Reader) (string, error) {
	return "MARKDOWN FILE", nil
}

func (_ MarkdownHandler) ServeFile(fs fs.FS, w http.ResponseWriter, filepath string, extension string, cfg Config) {
	WriteCSS(w, filepath)
	ServeFile(fs, w, filepath, "md")
	WriteWidget(w, filepath)
}
