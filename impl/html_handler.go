package impl

import (
	"io"
	"io/fs"
	"net/http"
	"regexp"
)

type HtmlHandler struct{}

func (_ HtmlHandler) Tokenize(x io.Reader) ([]string, error) {
	return TokenizeHtml(x)
}

func (_ HtmlHandler) ExtractTitle(x io.Reader) (string, error) {

	bytes, err := io.ReadAll(x)

	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`<title>([\s\S]*?)<\/title>`)

	matches := re.FindStringSubmatch(string(bytes))

	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", nil
}

func (_ HtmlHandler) ServeFile(fs fs.FS, w http.ResponseWriter, filepath string, extension string, cfg Config) {
	ServeFile(fs, w, filepath, "html")
	WriteWidget(w, cfg.CorpusDir+filepath)
}
