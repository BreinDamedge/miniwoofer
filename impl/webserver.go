package impl

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"strings"
)

type MiniWooferWeb struct{}

//go:embed html/*
var EmbededResources embed.FS

func serve_root(b *Bm25, db *MetaDb, w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	tmpl, err := template.ParseFS(EmbededResources, "html/default_page.html")

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	joined_terms := strings.Join(req.Form["query"], " ")

	re := regexp.MustCompile(`((?:".*?" *)|(?:(?:\S+? +)))`)

	search_terms := []string{}
	for _, match := range re.FindAllStringSubmatch(joined_terms+" ", -1) {
		if len(match) > 1 {
			search_terms = append(search_terms, strings.Trim(strings.ToLower(match[1]), ` "`))
		}
	}

	results, err := b.Search(search_terms)
	search_results := []DocumentMeta{}

	for _, result := range results {
		doc, err := db.GetDocument(result.Id)

		if err != nil {
			w.WriteHeader(404)
			fmt.Println(err)
			return
		}

		search_results = append(search_results, *doc)
	}

	if err := tmpl.Execute(w, struct {
		Joined_terms string
		Results      []DocumentMeta
	}{joined_terms, search_results}); err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}
}

func serve_corpus(fs fs.FS, w http.ResponseWriter, req *http.Request) {
	file_name := req.PathValue("file")

	_, ext, _ := strings.Cut(file_name, ".")

	switch ext {
	case "mht", "mhtml":
		serve_mht(fs, w, req)
	case "html":
		serve_html(fs, w, file_name)
	default:
		serve_file(fs, w, file_name)
	}
}

func serve_html(fs fs.FS, w http.ResponseWriter, filename string) {
	serve_file(fs, w, filename)
	widget, err := EmbededResources.Open("html/widget.html")
	if err != nil {
		fmt.Println("Couldnt find widget!")
		return
	}

	io.Copy(w, widget)
}

func serve_file(fs fs.FS, w http.ResponseWriter, file_name string) {
	_, ext, _ := strings.Cut(file_name, ".")
	mime_type := mime.TypeByExtension("." + ext)
	file, err := fs.Open(file_name)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", mime_type)
	io.Copy(w, file)

}

func serve_mht(fs fs.FS, w http.ResponseWriter, req *http.Request) {
	file_name := req.PathValue("file")
	file, err := fs.Open(file_name)
	if err != nil {

		w.WriteHeader(404)
		fmt.Fprintf(w, "%+v", err)
		return
	}

	msg, err := mail.ReadMessage(file)
	if err != nil {
		fmt.Fprintf(w, "Error while parsing file %s: %+v", file_name, err)
		w.WriteHeader(500)
		return
	}

	content_type := msg.Header.Get("Content-Type")
	_, params, err := mime.ParseMediaType(content_type)

	if err != nil {
		fmt.Fprintf(w, "Error while parsing mhtml %s: %+v", file_name, err)
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
			fmt.Fprintf(w, "Error while parsing part in file %s: %+v", file_name, err)
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
	widget, err := EmbededResources.Open("html/widget.html")
	if err != nil {
		fmt.Println("Couldnt find widget")
		return
	}

	io.Copy(w, widget)
}

func (web *MiniWooferWeb) Run(b *Bm25, db *MetaDb, config Config) error {
	fs := os.DirFS(config.CorpusDir)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { serve_root(b, db, w, r) })
	http.HandleFunc("/corpus/{file...}", func(w http.ResponseWriter, r *http.Request) { serve_corpus(fs, w, r) })
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.WebserverPort), nil)
}
