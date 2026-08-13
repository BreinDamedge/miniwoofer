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
	"strings"
)

type BoogalooWeb struct{}

// this abuses that when a browser sees 2 <body> tags it merges them, hope this works
var search_bar string = `
<body>
<h2>Boogaloo</h2>
<form action="/" method="GET">
<label for="search">Search:</label>
<input type=text id="search" name=query></br>
</form>
</body>
`

func serve_root(b *Bm25, w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Fprintf(w, `
		<head>
		<title>Boogaloo</title>
		</head>
		%s
		<body>
	`, search_bar)

	defer fmt.Fprintf(w, "</body>")

	if req.Form["query"] != nil {
		results, err := b.Search(req.Form["query"])
		if err != nil {
			fmt.Fprintf(w, "Error while searching %+v: %+v", req.Form["query"], err)
			w.WriteHeader(500)
			return
		}
		for _, result := range results {
			fmt.Fprintf(w, `
				<a href=%s>%s</a></br>
			`,
				result.Id,
				result.Id,
			)
		}
	}
}

func serve_corpus(fs fs.FS, w http.ResponseWriter, req *http.Request) {
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

	fmt.Fprintf(w, "%s\n", search_bar)

	for {
		part, err := mp_reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(w, "Error while parsing part in file %s: %+v", file_name, err)
			w.WriteHeader(500)
			return
		}

		ct := part.Header.Get("Content-Type")

		body_bytes, _ := io.ReadAll(part)
		if strings.Contains(ct, "html") {

			fmt.Fprintf(w, "%s\n", string(body_bytes))
		} else if strings.Contains(ct, "css") {
			fmt.Fprintf(w,
				`
		<style>%s</style>\n
		`, string(body_bytes))
		}

	}

}

func (web *BoogalooWeb) Run(b *Bm25, config Config) error {
	fs := os.DirFS(config.CorpusDir)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { serve_root(b, w, r) })
	http.HandleFunc("/corpus/{file}", func(w http.ResponseWriter, r *http.Request) { serve_corpus(fs, w, r) })
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.WebserverPort), nil)
}
