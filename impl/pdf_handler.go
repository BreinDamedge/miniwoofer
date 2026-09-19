package impl

import (
	"io"
	"io/fs"
	"net/http"

	"bytes"
	"github.com/dslipak/pdf"
)

func readPdf(path string) (string, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

type PDFHandler struct{}

func (_ PDFHandler) Tokenize(x io.Reader) ([]string, error) {
	// extract plaintext from PDF w/package
	// ...

	return TokenizePlaintext(x)
}

func (_ PDFHandler) ExtractTitle(x io.Reader) (string, error) {
	return "PDF FILE", nil
}

func (_ PDFHandler) ServeFile(fs fs.FS, w http.ResponseWriter, filepath string, extension string, cfg Config) {
	ServeFile(fs, w, filepath, "pdf")
	WriteWidget(w, filepath)
}
