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
	"path/filepath"
	"regexp"
	"strings"
)

type MiniWooferWeb struct{}

//go:embed html/*
var EmbededResources embed.FS

func write_widget(w http.ResponseWriter) {
	widget, err := EmbededResources.Open("html/widget.html")
	if err != nil {
		fmt.Println("Couldnt find widget!")
		return
	}

	io.Copy(w, widget)

}

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
		// debug print (tho changing this to a "table didn't initialize state" may be prefered)
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

func rescan(_ fs.FS, index *Bm25, mdb *MetaDb, cfg Config, w http.ResponseWriter, _ *http.Request) {
	fmt.Println("rescan triggered!")

	if err := ParseCorpus(index, cfg); err != nil {
		w.WriteHeader(500)
		panic(err)
	}
	if err := mdb.AddCorpus(cfg); err != nil {
		w.WriteHeader(500)
		panic(err)
	}

	w.Header().Set("Location", "/")
	w.WriteHeader(301)
}

func serve_corpus(fs fs.FS, w http.ResponseWriter, req *http.Request, cfg Config) {
	file_name := req.PathValue("file")

	ext := strings.TrimLeft(filepath.Ext(file_name), ".")

	switch ext {
	case "mht", "mhtml":
		serve_mht(fs, w, req)
	case "html":
		serve_html(fs, w, file_name)
	case "txt":
		serve_txt(fs, w, file_name, cfg)
	default:
		serve_file(fs, w, file_name, ext)
	}
}

func serve_html(fs fs.FS, w http.ResponseWriter, filename string) {
	serve_file(fs, w, filename, "html")
	write_widget(w)
}

func serve_txt(fs fs.FS, w http.ResponseWriter, filename string, cfg Config) {

	if cfg.StylePlaintext {
		// do logic for html styling here
		serve_file(fs, w, filename, "html")
		write_widget(w)
	} else {
		serve_file(fs, w, filename, "txt")
	}
}

func serve_file(fs fs.FS, w http.ResponseWriter, file_name string, extension string) {

	mime_type := mime.TypeByExtension("." + extension)
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

	write_widget(w)
}

func (web *MiniWooferWeb) Run(b *Bm25, db *MetaDb, config Config) error {
	fs := os.DirFS(config.CorpusDir)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { serve_root(b, db, w, r) })
	http.HandleFunc("/corpus/{file...}", func(w http.ResponseWriter, r *http.Request) { serve_corpus(fs, w, r, config) })
	http.HandleFunc("/triggers/rescan", func(w http.ResponseWriter, r *http.Request) { rescan(fs, b, db, config, w, r) })

	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.WebserverPort), nil)
}
