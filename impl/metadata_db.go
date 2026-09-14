package impl

import (
	"database/sql"
	"fmt"
	"io/fs"
	"mime"
	"os"
	"path/filepath"

	_ "github.com/ncruces/go-sqlite3/driver"
)

type MetaDb struct {
	db *sql.DB
}

type DocumentMeta struct {
	Id    string
	Title string
}

func ForwardBackslashes(windows_path_ string) string {
	// replace all the backslashes in windows_path_ with forward slashes
	// going to treat it as a utf-8 string for now
	// convert to a rune slice then replace chars
	runes := []rune(windows_path_)
	num_chars := len(runes)
	backslash := rune('\\')
	forwardslash := rune('/')
	for i := range num_chars {
		if runes[i] == backslash {
			runes[i] = forwardslash
		}
	}

	return string(runes)
}

func (dm *DocumentMeta) normId() {
	// "normalize" the windows doc path to unix style
	dm.Id = ForwardBackslashes(dm.Id)
}

func MetaDbOpen(config Config) (*MetaDb, error) {
	should_generate := false
	if err := os.MkdirAll(filepath.Dir(config.DatabasePath()), 0777); err != nil && err != os.ErrExist {
		fmt.Println(err)
		return nil, err
	}
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
	doc.normId()
	_, err := md.db.Exec("insert into documents values (?, ?) on conflict do nothing;", doc.Id, doc.Title)
	return err
}

func (md *MetaDb) UpsertDocument(doc DocumentMeta) error {
	doc.normId()
	_, err := md.db.Exec(`insert into documents values (?, ?) on conflict update title = ?;`, doc.Id, doc.Title, doc.Title)
	return err
}

func (md *MetaDb) AddCorpus(config Config) error {
	defer md.db.Exec("COMMIT;") // commit changes once done parsing TODO: is this proper defer usage?
	return filepath.WalkDir(config.CorpusDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		mime_type := mime.TypeByExtension(ext)
		title := ""

		if handler, ok := FileHandlers[mime_type]; ok {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			title, err = handler.ExtractTitle(f)
			if err != nil {
				return err
			}
		}

		if title == "" {
			title = filepath.Base(path)
		}

		md.AddDocument(DocumentMeta{Id: path, Title: title})
		return nil
	})
}
