package impl

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"

	"golang.org/x/net/html"
)

var textMimeTypes = map[string]struct{}{
	"text/plain": {},
	"text/html":  {},
}

type multipartParser struct {
	tokens []string
}

func (t *multipartParser) parsePlainContent(x string) error {
	t.tokens = append(t.tokens, Tokenize(string(x))...)
	return nil
}

func (t *multipartParser) parseHtmlContent(x []byte) error {
	dom, err := html.Parse(strings.NewReader(string(x)))
	if err != nil {
		return err
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
	return t.parsePlainContent(text.String())
}

func (t *multipartParser) parseTextContent(x []byte, ct string) error {
	switch ct {
	case "text/plain":
		return t.parsePlainContent(string(x))
	case "text/html":
		return t.parseHtmlContent(x)
	default:
		return fmt.Errorf("unknown content type: %s", ct)
	}
}

func (t *multipartParser) parsePart(part *multipart.Part) error {
	ct := part.Header.Get("Content-Type")
	if ct == "" {
		ct = "text/plain"
	}

	mediaType, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return err
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		return t.parseMultipart(part, params["boundary"])
	}

	if _, ok := textMimeTypes[mediaType]; !ok {
		return nil
	}

	body, err := io.ReadAll(part)
	if err != nil {
		return err
	}

	return t.parseTextContent(body, mediaType)
}

func (t *multipartParser) parseMultipart(x io.Reader, bound string) error {
	if bound == "" {
		return fmt.Errorf("multipart content missing boundary")
	}

	reader := multipart.NewReader(x, bound)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if err := t.parsePart(part); err != nil {
			return err
		}
	}
}

func (t *multipartParser) parseMhtml(x io.Reader) error {
	msg, err := mail.ReadMessage(x)
	if err != nil {
		return err
	}

	ct := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return err
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		return t.parseMultipart(msg.Body, params["boundary"])
	} else {
		return fmt.Errorf("expected multipart media, got %s", ct)
	}
}

func TokenizeMhtml(x io.Reader) ([]string, error) {
	var t multipartParser
	if err := t.parseMhtml(x); err != nil {
		return nil, err
	}

	return t.tokens, nil
}

func TokenizeHtml(x io.Reader) ([]string, error) {
	var t multipartParser
	bytes, err := io.ReadAll(x)
	if err != nil {
		return nil, err
	}
	if err = t.parseHtmlContent(bytes); err != nil {
		return nil, err
	}
	return t.tokens, err
}
