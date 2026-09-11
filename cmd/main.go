package main

import (
	"fmt"
	"log"
	"miniwoofer/impl"
)

func main() {
	index := impl.NewBm25()
	config := impl.LoadConfig()
	db, err := impl.MetaDbOpen(config)
	if err != nil {
		panic(err)
	}

	force_build_metadb := false // force the build if the index didn't exist (fixes changed bug on first build)
	fmt.Print("Trying to load serialized BM25\n")
	if err := index.Load(config.Bm25Path()); err != nil {
		fmt.Print("Serialized BM25 Does not exist\n")
		if err := impl.ParseCorpus(index, config); err != nil {
			panic(err)
		}
		force_build_metadb = true
	}

	fmt.Print("Checking if corpus changed\n")
	if changed, err := impl.CheckChanged(config.CorpusDir, index.CorpusHash); err != nil || changed || force_build_metadb {
		fmt.Print("Corpus Changed\n")
		if err := impl.ParseCorpus(index, config); err != nil {
			panic(err)
		}
		if err := db.AddCorpus(config); err != nil {
			panic(err)
		}
	}

	web := impl.MiniWooferWeb{}
	if err := web.Run(index, db, config); err != nil {
		log.Fatal(err)
	}

}
