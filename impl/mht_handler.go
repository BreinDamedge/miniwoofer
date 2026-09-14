package impl

import (
	"fmt"
	"io"
	"io/fs"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
)

type MhtHandler struct{}

func (_ MhtHandler) Tokenize(x io.Reader) ([]string, error) {
	return TokenizeMhtml(x)
}

func (_ MhtHandler) ExtractTitle(x io.Reader) (string, error) {

	msg, err := mail.ReadMessage(x)
	if err != nil {
		return "", err
	}

	// oo for title extraction surely
	title := ""

	// reading yup
	content_type := msg.Header.Get("Content-Type")
	_, params, err := mime.ParseMediaType(content_type)
	if err != nil {
		return "", err
	}

	mp_reader := multipart.NewReader(msg.Body, params["boundary"])

	for {
		part, err := mp_reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		if ct, _, err := mime.ParseMediaType(part.Header.Get("Content-Type")); err == nil && ct == "text/html" {
			body_bytes, err := io.ReadAll(part)

			if err != nil {
				return "", err
			}

			re := regexp.MustCompile(`<title>([\s\S]*?)<\/title>`)

			matches := re.FindStringSubmatch(string(body_bytes))

			if len(matches) < 2 {
				continue
			}
			title = matches[1]

			break
		} else if err != nil {
			return "", err
		}
	}

	// setting title if it's not extracted well
	if title == "" {
		title = msg.Header.Get("Subject")
	}

	title = strings.Trim(title, "\r\n")
	return title, nil
}

func (_ MhtHandler) ServeFile(fs fs.FS, w http.ResponseWriter, filepath string, extension string, cfg Config) {
	file, err := fs.Open(filepath)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "%+v", err)
		return
	}

	msg, err := mail.ReadMessage(file)
	if err != nil {
		fmt.Fprintf(w, "Error while parsing file %s: %+v", filepath, err)
		w.WriteHeader(500)
		return
	}

	content_type := msg.Header.Get("Content-Type")
	_, params, err := mime.ParseMediaType(content_type)

	if err != nil {
		fmt.Fprintf(w, "Error while parsing mhtml %s: %+v", filepath, err)
		w.WriteHeader(500)
		return
	}

	mp_reader := multipart.NewReader(msg.Body, params["boundary"])

	// fmt.Fprintf(w, "%s\n", search_bar)

	for {
		part, err := mp_reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(w, "Error while parsing part in file %s: %+v", filepath, err)
			w.WriteHeader(500)
			return
		}

		ct, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))

		if err != nil {
			fmt.Printf("Failed to read media type of %+v: %+v\n", part, err)
			continue
		}

		body_bytes, err := io.ReadAll(part)

		if err != nil {
			fmt.Printf("Failed to read body of %+v: %+v", part.Header, err)
			continue
		}

		switch ct {
		case "text/html":
			fmt.Fprintf(w, "%s\n", string(body_bytes))
		case "text/css":
			fmt.Fprintf(w, "<style>%s</style>\n", string(body_bytes))
		}

	}

	WriteWidget(w)

}
