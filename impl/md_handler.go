package impl

import (
	"io"
	"io/fs"
	"net/http"
	"regexp"

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
	// convert the markdown to html
	html_data, err := md_to_html(x)
	if err != nil {
		return []string{}, err
	}
	// pass to the html tokenizer
	return TokenizeHtml(bytes.NewReader(html_data))
}

func (_ MarkdownHandler) ExtractTitle(x io.Reader) (string, error) {
	// convert the markdown to html
	html_data, err := md_to_html(x)
	if err != nil {
		return "EMPTY_TITLE", err
	}

	// use the same title extraction method as html
	re := regexp.MustCompile(`<title>([\s\S]*?)<\/title>`)

	matches := re.FindStringSubmatch(string(html_data))

	if len(matches) > 1 {
		return matches[1], nil
	}
	return "EMPTY_TITLE", nil
}

func (_ MarkdownHandler) ServeFile(fs fs.FS, w http.ResponseWriter, filepath string, extension string, cfg Config) {
	ServeFile(fs, w, filepath, "md")
}
