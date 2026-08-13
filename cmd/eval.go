package main

import (
	"fmt"
	"io/fs"
	"miniwoofer/impl"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	test := impl.NewBm25()
	cfg := impl.LoadCfg()
	test.SetParams(cfg)

	documents := []impl.Doc{}

	if err := filepath.WalkDir("corpus", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if !strings.HasSuffix(path, ".mht") && !strings.HasSuffix(path, ".mhtml") {
			return nil
		}

		reader, err := os.Open(path)
		if err != nil {
			return err
		}

		tokens, err := impl.TokenizeMhtml(reader)
		if err != nil {
			return err
		}

		documents = append(documents, impl.Doc{
			Id:  path,
			Tok: tokens,
		})

		return nil
	}); err != nil {
		panic(err)
	}

	// add the docs to the bm25 interface (it stores)
	for _, doc := range documents {
		if err := test.Append(doc.Id, doc.Tok); err != nil {
			panic(err)
		}
	}

	// run optimizer
	// load the tuning data from toml
	tuningData := impl.LoadQueryData()

	// evaluate these params on the dataset
	cost := impl.Evaluate(test, documents, tuningData)
	fmt.Printf("params: %f, %f\ncost: %f\n", cfg.B, cfg.K1, cost)
}
