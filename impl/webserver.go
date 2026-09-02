package impl

import (
	"fmt"
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

// this abuses that when a browser sees 2 <body> tags it merges them, hope this works
const search_bar string = `
<div>
<h2>MiniWoofer</h2>
<form action="/" method="GET">
<label for="search">Search:</label>
<input type=text id="search" name=query value=%s>
<button formmethod="GET" formtarget="search">Go!</button>
</form>
</div>
`

func serve_root(b *Bm25, db *MetaDb, w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Fprint(w, `
		<head>
		<title>miniwoofer</title>
		</head>
		<body>
	`)

	defer fmt.Fprintf(w, "</body>")

	joined_terms := strings.Join(req.Form["query"], " ")

	re := regexp.MustCompile(`((?:".*?" *)|(?:(?:\S+? +)))`)

	search_terms := []string{}

	for _, match := range re.FindAllStringSubmatch(joined_terms+" ", -1) {
		if len(match) > 1 {
			search_terms = append(search_terms, strings.Trim(strings.ToLower(match[1]), ` "`))
		}
	}

	fmt.Fprintf(w, search_bar, joined_terms)
	if req.Form["query"] != nil {
		results, err := b.Search(search_terms)
		if err != nil {
			fmt.Fprintf(w, "Error while searching %+v: %+v", req.Form["query"], err)
			w.WriteHeader(500)
			return
		}
		for _, result := range results {
			doc, _ := db.GetDocument(result.Id)
			fmt.Fprintf(w, `
				<a href="%s">%s</a></br>
			`,
				doc.Id,
				doc.Title,
			)
		}
	}
}

func serve_corpus(fs fs.FS, w http.ResponseWriter, req *http.Request) {
	file_name := req.PathValue("file")

	_, ext, _ := strings.Cut(file_name, ".")

	switch ext {
	case "mht", "mhtml":
		serve_mht(fs, w, req)
	default:
		serve_file(fs, w, file_name)
	}
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

}

func (web *MiniWooferWeb) Run(b *Bm25, db *MetaDb, config Config) error {
	fs := os.DirFS(config.CorpusDir)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { serve_root(b, db, w, r) })
	http.HandleFunc("/corpus/{file...}", func(w http.ResponseWriter, r *http.Request) { serve_corpus(fs, w, r) })
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.WebserverPort), nil)
}
