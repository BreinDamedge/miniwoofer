package impl

import (
	"io"
	"io/fs"
	"mime"
	"net/http"
)

type FileHandler interface {
	Tokenize(x io.Reader) ([]string, error)
	ExtractTitle(x io.Reader) (string, error)
	ServeFile(fs fs.FS, w http.ResponseWriter, filepath string, extension string, config Config)
}

var FileHandlers = map[string]FileHandler{
	mime.TypeByExtension(".txt"):   TextHandler{},
	mime.TypeByExtension(".html"):  HtmlHandler{},
	"multipart/related":            MhtHandler{},
	mime.TypeByExtension(".md"):    MarkdownHandler{},
	"text/markdown":                MarkdownHandler{},
	"text/markdown; charset=utf-8": MarkdownHandler{},
}
