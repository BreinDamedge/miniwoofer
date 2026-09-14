package impl

import (
	"io"
	"io/fs"
	"net/http"
)

type TextHandler struct{}

func (_ TextHandler) Tokenize(x io.Reader) ([]string, error) {
	return TokenizePlaintext(x)
}

func (_ TextHandler) ExtractTitle(_ io.Reader) (string, error) {
	return "", nil
}

func (_ TextHandler) ServeFile(fs fs.FS, w http.ResponseWriter, filepath string, extension string, cfg Config) {
	if cfg.StylePlaintext {
		// do logic for html styling here
		ServeFile(fs, w, filepath, "html")
		WriteWidget(w)

	} else {
		ServeFile(fs, w, filepath, "txt")
	}
}
