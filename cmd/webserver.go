package main

import (
	"fmt"
	"log"
	"miniwoofer/impl"
)

func main() {
	test := impl.NewBm25()
	config := impl.LoadConfig()

	fmt.Print("Trying to load serialized BM25\n")
	if err := test.Load(config.Bm25File); err != nil {
		fmt.Print("Serialized BM25 Does not exist\n")
		if err := impl.ParseCorpus(test, config); err != nil {
			panic(err)
		}
	}

	fmt.Print("Checking if corpus changed\n")
	if changed, err := impl.CheckChanged(config.CorpusDir, test.CorpusHash); err != nil || changed {
		fmt.Print("Corpus Changed\n")
		if err := impl.ParseCorpus(test, config); err != nil {
			panic(err)
		}
	}

	web := impl.BoogalooWeb{}
	if err := web.Run(test, config); err != nil {
		log.Fatal(err)
	}

}
