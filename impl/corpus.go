package impl

/*

Basic Corpus operations

*/
import (
	"crypto/sha1"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
		if err != nil {
			return err
		} else if !strings.HasSuffix(path, ".mht") && !strings.HasSuffix(path, ".mhtml") && !strings.HasSuffix(path, ".html") {
			return nil
		}

		reader, err := os.Open(path)
		if err != nil {
			return err
		}
		var tokens []string

		if strings.HasSuffix(path, ".mht") || strings.HasSuffix(path, ".mhtml") {
			tokens, err = TokenizeMhtml(reader)
			if err != nil {
				return err
			}
		} else {
			tokens, err = TokenizeHtml(reader)
			if err != nil {
				return err
			}
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

	// run optimizer
	// load the tuning data from toml
	// fmt.Println("Loading tuning data... ")
	// tuningData := LoadQueryData()
	// fmt.Println("Done.")

	// do the fitting
	fmt.Println("Fitting b & k1...")
	// b.SetParams(Optimize(b, documents, tuningData))
	fmt.Println("Done.")

	hash, err := HashDir(config.CorpusDir)
	if err != nil {
		return err
	}
	b.CorpusHash = hash

	return b.Save(config.Bm25Path())
}
