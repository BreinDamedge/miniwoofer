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
var HtmlFiles embed.FS

func WriteWidget(w http.ResponseWriter, id string) {
	widget, err := template.ParseFS(HtmlFiles, "html/widget.html")

	if err != nil {
		fmt.Println("Couldnt find widget!")
		return
	}

	widget.Execute(w, id)

}

func serve_root(b *Bm25, db *MetaDb, config Config, w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	tmpl, err := template.ParseFS(HtmlFiles, "html/default_page.html")

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
		id, _ := filepath.Rel(config.CorpusDir, result.Id)
		doc, err := db.GetDocument(id)

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
	WriteWidget(w, filename)
}

func ServeFile(fs fs.FS, w http.ResponseWriter, file_name string, extension string) {

	mime_type := mime.TypeByExtension("." + extension)
	file, err := fs.Open(file_name)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", mime_type)
	// refactor this eventually
	if extension == "md" {
		w.Header().Set("Content-Type", "text/html")

		buf, err := markdown_to_html(file)
		if err != nil {
			panic(err)
		}
		io.Copy(w, &buf)
	} else {
		io.Copy(w, file)
	}

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

func management_endpoint(b *MetaDb, w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	doc_id := req.Form.Get("doc")
	if doc_id == "" {
		w.WriteHeader(400)
		return
	}
	tmpl, err := template.ParseFS(HtmlFiles, "html/management.html")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	doc, err := b.GetDocument(doc_id)
	if err != nil {
		w.WriteHeader(404)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, *doc); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}

func delete_document(index *Bm25, db *MetaDb, config Config, w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	doc_id := req.Form.Get("Id")

	if err := db.DeleteDocument(doc_id); err != nil {
		fmt.Println(err)
	}
	os.Remove(filepath.Join(config.CorpusDir, doc_id))

	if err := ParseCorpus(index, config); err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(200)
}

func update_blurb(db *MetaDb, w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	new_blurb := req.Form.Get("blurb")
	doc_id := req.Form.Get("Id")

	if err := db.UpdateBlurb(doc_id, new_blurb); err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)
}

func update_title(db *MetaDb, w http.ResponseWriter, req *http.Request) {
	// lol I was just about to google this but then you wrote it already thanks
	req.ParseForm()
	new_title := req.Form.Get("title")
	doc_id := req.Form.Get("Id")

	if err := db.UpdateTitle(doc_id, new_title); err != nil {
		fmt.Println(err)
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)
}

func (web *MiniWooferWeb) Run(b *Bm25, db *MetaDb, config Config) error {
	fs := os.DirFS(config.CorpusDir)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) { serve_root(b, db, config, w, r) })
	http.HandleFunc("GET /corpus/{file...}", func(w http.ResponseWriter, r *http.Request) { serve_corpus(fs, w, r, config) })
	http.HandleFunc("GET /triggers/rescan", func(w http.ResponseWriter, r *http.Request) { rescan(fs, b, db, config, w, r) }) // this realistically should be a POST, need to fix
	http.HandleFunc("GET /management.html", func(w http.ResponseWriter, r *http.Request) { management_endpoint(db, w, r) })
	http.HandleFunc("POST /triggers/delete", func(w http.ResponseWriter, r *http.Request) { delete_document(b, db, config, w, r) })
	http.HandleFunc("POST /triggers/update_blurb", func(w http.ResponseWriter, r *http.Request) { update_blurb(db, w, r) })
	http.HandleFunc("POST /triggers/update_title", func(w http.ResponseWriter, r *http.Request) { update_title(db, w, r) })

	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.WebserverPort), nil)
}
