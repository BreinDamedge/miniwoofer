package impl

import (
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "github.com/ncruces/go-sqlite3/driver"
)

type MetaDb struct {
	db *sql.DB
}

type DocumentMeta struct {
	Id    string
	Title string
}

func MetaDbOpen(config Config) (*MetaDb, error) {

	should_generate := false

	if _, err := os.Stat(config.DatabasePath()); err != nil {
		should_generate = true
	}

	db, err := sql.Open("sqlite3", "file:"+config.DatabasePath())
	if err != nil {
		return nil, err
	}
	res := &MetaDb{db: db}

	if should_generate {
		res.initDb()

		res.AddCorpus(config)
	}
	return res, nil
}

func (md *MetaDb) initDb() error {
	table_rows, err := md.db.Query("select name from sqlite_master where type='table' and name='documents'")
	if err != nil {
		return err
	}
	if table_rows.Next() {
		return nil
	}
	table_rows.Close()

	_, err = md.db.Exec(`
		create table documents (
  		id varchar(500) primary key,
  		title varchar(1000) 
		);	
	`)

	return err
}

func (md *MetaDb) GetDocument(id string) (*DocumentMeta, error) {
	rows, err := md.db.Query("select * from documents where id = ?;", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	doc := &DocumentMeta{}
	rows.Next()
	if err := rows.Scan(&doc.Id, &doc.Title); err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}

	return doc, nil
}

func (md *MetaDb) AddDocument(doc DocumentMeta) error {

	_, err := md.db.Exec("insert into documents values (?, ?) on conflict do nothing;", doc.Id, doc.Title)
	return err
}

func (md *MetaDb) UpsertDocument(doc DocumentMeta) error {
	_, err := md.db.Exec(`insert into documents values (?, ?) on conflict update title = ?;`, doc.Id, doc.Title, doc.Title)
	return err
}

func (md *MetaDb) AddCorpus(config Config) error {
	return filepath.WalkDir(config.CorpusDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".mht") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			msg, err := mail.ReadMessage(file)
			if err != nil {
				return err
			}

			title := ""

			content_type := msg.Header.Get("Content-Type")
			_, params, err := mime.ParseMediaType(content_type)
			if err != nil {
				return err
			}

			mp_reader := multipart.NewReader(msg.Body, params["boundary"])

			for {
				part, err := mp_reader.NextPart()
				if err == io.EOF {
					break
				} else if err != nil {
					return err
				}

				if ct, _, err := mime.ParseMediaType(part.Header.Get("Content-Type")); err == nil && ct == "text/html" {
					body_bytes, err := io.ReadAll(part)

					if err != nil {
						return err
					}

					re, err := regexp.Compile(`<title>([\s\S]*?)<\/title>`)

					if err != nil {
						return err
					}

					matches := re.FindStringSubmatch(string(body_bytes))

					if len(matches) < 2 {
						continue
					}
					title = matches[1]

					break
				} else if err != nil {
					return err
				}
			}

			if title == "" {
				title = msg.Header.Get("Subject")
			}
			if title == "" {
				title = path
			}

			title = strings.Trim(title, "\r\n")

			return md.AddDocument(DocumentMeta{Id: path, Title: title})
		}
		return nil
	})
}
