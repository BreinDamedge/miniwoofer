package impl

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

var textMimeTypes = map[string]struct{}{
	"text/plain": {},
	"text/html":  {},
}

func parsePlainContent(x string) []string {
	tokens := Tokenize(string(x))

	return tokens
}

func parseHtmlContent(x []byte) ([]string, error) {
	dom, err := html.Parse(strings.NewReader(string(x)))
	if err != nil {
		return nil, err
	}

	var text strings.Builder
	var walk func(*html.Node)

	walk = func(node *html.Node) {
		if node == nil {
			return
		}

		if node.Type == html.ElementNode {
			switch node.Data {
			case "head", "script", "style", "noscript":
				return
			}
		} else if node.Type == html.TextNode {
			text.WriteString(node.Data)
			text.WriteByte(' ')
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(dom)
	return parsePlainContent(text.String()), nil

}

func parseTextContent(x []byte, ct string) ([]string, error) {
	switch ct {
	case "text/plain":
		return parsePlainContent(string(x)), nil
	case "text/html":
		return parseHtmlContent(x)
	default:
		return nil, fmt.Errorf("unknown content type: %s", ct)
	}
}

func parsePart(part *multipart.Part) ([]string, error) {
	ct := part.Header.Get("Content-Type")
	if ct == "" {
		ct = "text/plain"
	}

	mediaType, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		return parseMultipart(part, params["boundary"])
	}

	if _, ok := textMimeTypes[mediaType]; !ok {
		return []string{}, nil
	}

	body, err := io.ReadAll(part)
	if err != nil {
		return nil, err
	}

	return parseTextContent(body, mediaType)
}

func parseMultipart(x io.Reader, bound string) ([]string, error) {
	if bound == "" {
		return nil, fmt.Errorf("multipart content missing boundary")
	}

	tokens := []string{}

	reader := multipart.NewReader(x, bound)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			return tokens, nil
		}
		if err != nil {
			return nil, err
		}

		if tok, err := parsePart(part); err != nil {
			return nil, err
		} else {
			tokens = slices.Concat(tokens, tok)
		}
	}
}

func parseMhtml(x io.Reader) ([]string, error) {
	msg, err := mail.ReadMessage(x)
	if err != nil {
		return nil, err
	}

	ct := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		return parseMultipart(msg.Body, params["boundary"])
	} else {
		return nil, fmt.Errorf("expected multipart media, got %s", ct)
	}
}

func TokenizeMhtml(x io.Reader) ([]string, error) {
	return parseMhtml(x)
}

func TokenizeHtml(x io.Reader) ([]string, error) {
	bytes, err := io.ReadAll(x)
	if err != nil {
		return nil, err
	}
	return parseHtmlContent(bytes)
}

func TokenizePlaintext(x io.Reader) ([]string, error) {
	bytes, err := io.ReadAll(x)
	if err != nil {
		return nil, err
	}
	return parsePlainContent(string(bytes)), nil
}
