package impl

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"bytes"

	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

type MiniWooferWeb struct{}

//go:embed html/*
var EmbededResources embed.FS

func WriteWidget(w http.ResponseWriter) {
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
	mime_type := mime.TypeByExtension("." + ext)

	if handler, ok := FileHandlers[mime_type]; ok {
		handler.ServeFile(fs, w, file_name, ext, cfg)
	} else {
		ServeFile(fs, w, file_name, ext)
	}

}

func serve_html(fs fs.FS, w http.ResponseWriter, filename string) {
	ServeFile(fs, w, filename, "html")
	WriteWidget(w)
}

func ServeFile(fs fs.FS, w http.ResponseWriter, file_name string, extension string) {

	mime_type := mime.TypeByExtension("." + extension)
	file, err := fs.Open(file_name)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	// refactor this eventually
	if extension == "md" {
		mime_type = "text/html"
		buf, err := markdown_to_html(file)
		if err != nil {
			panic(err)
		}
		io.Copy(w, &buf)
	} else {
		io.Copy(w, file)
	}

	w.Header().Set("Content-Type", mime_type)
}

func markdown_to_html(file fs.File) (bytes.Buffer, error) {
	var buf bytes.Buffer
	file_bytes, err := io.ReadAll(file) // probably a better way to do this
	if err != nil {
		return buf, err
	}

	p := parser.New(parser.WithAttribute(), parser.WithExtensions(extension.GFMParser))
	r := html.New(html.WithXHTML(), html.WithUnsafe(), html.WithExtensions(extension.GFMHTMLRenderer))

	doc := p.Parse(file_bytes)
	if err := r.Render(&buf, file_bytes, doc); err != nil {
		return buf, err
	}

	return buf, nil
}

func serve_markdown(fs fs.FS, w http.ResponseWriter, file_path string) {
	// parse, render to html, and then respond w/html version of markdown. uses goldmark

	file, err := fs.Open(file_path)
	if err != nil {
		w.WriteHeader(500)
		panic(err)
	}

	buf, err := markdown_to_html(file)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/html")
	io.Copy(w, &buf)
}

func (web *MiniWooferWeb) Run(b *Bm25, db *MetaDb, config Config) error {
	fs := os.DirFS(config.CorpusDir)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { serve_root(b, db, w, r) })
	http.HandleFunc("/corpus/{file...}", func(w http.ResponseWriter, r *http.Request) { serve_corpus(fs, w, r, config) })
	http.HandleFunc("/triggers/rescan", func(w http.ResponseWriter, r *http.Request) { rescan(fs, b, db, config, w, r) })

	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.WebserverPort), nil)
}
