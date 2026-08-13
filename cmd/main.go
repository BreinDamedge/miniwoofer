package main

import (
	"bufio"
	"fmt"
	"miniwoofer/impl"
	"os"
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

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		tokenized := impl.Tokenize(line)

		fmt.Printf("You searched: %s\n", line)
		fmt.Printf("We tokenized: %s\n", tokenized)

		result, err := test.Search(tokenized)
		if err != nil {
			panic(err)
		}

		fmt.Printf("We found %d results:\n", len(result))
		for i, x := range result {
			fmt.Printf(" - #%d doc=%s, score=%f\n", i+1, x.Id, x.Score)
		}
	}
}
