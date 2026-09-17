package impl

/*

Basic Corpus operations

*/
import (
	"crypto/sha1"
	"fmt"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
)

// Sha1 hash of File name and modification date of all files in "path"
func HashDir(path string) (string, error) {

	hasher := sha1.New()
	hash_bytes := []byte{}

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		mod_bytes, err := info.ModTime().Local().MarshalBinary()
		if err != nil {
			return err
		}

		name_bytes := []byte(info.Name())

		_, err = hasher.Write(append(mod_bytes, name_bytes...))

		return err
	})

	if err != nil {
		return "", err
	}

	hash_bytes = hasher.Sum(hash_bytes)
	return fmt.Sprintf("%x", hash_bytes), nil
}

// Compare a previous directory hash to another to see if a new file was added or changed
func CheckChanged(path string, old_hash string) (bool, error) {
	current_hash, err := HashDir(path)
	if err != nil {
		return true, err
	}

	return current_hash != old_hash, nil

}

// This is just yoinked from cmd/main.go for reuse purposes, should be made much more configurable lol
func ParseCorpus(b *Bm25, config Config) error {
	*b = *NewBm25()
	b.SetParams(Bm25Parameters{B: config.B, K1: config.K1})

	fmt.Println("Parsing corpus...")
	documents := []Doc{}
	if err := filepath.WalkDir(config.CorpusDir, func(path string, d fs.DirEntry, err error) error {

		// do nothing if we're looking at the dir instead of a file
		if path == "corpus/" {
			return nil
		}

		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		mime_type := mime.TypeByExtension(ext)

		handler, ok := FileHandlers[mime_type]
		if !ok {
			fmt.Printf("not ok while parsing '%s'\n", path)
			fmt.Printf("mime type was: '%s'\n", mime_type)
			fmt.Printf("ext was: '%s'\n", ext)
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}

		tokens, err := handler.Tokenize(f)

		if err != nil {
			return err
		}

		documents = append(documents, Doc{
			Id:  path,
			Tok: tokens,
		})

		return nil
	}); err != nil {
		return err
	}
	fmt.Println("Done")

	fmt.Println("Populating Index...")
	// add the docs to the bm25 interface (it stores)
	for _, doc := range documents {
		if err := b.Append(doc.Id, doc.Tok); err != nil {
			panic(err)
		}
	}
	fmt.Println("Done")

	hash, err := HashDir(config.CorpusDir)
	if err != nil {
		return err
	}
	b.CorpusHash = hash

	return b.Save(config.Bm25Path())
}
