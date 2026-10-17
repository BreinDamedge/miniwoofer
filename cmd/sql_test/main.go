package main

// main file for μWoofer experiment

import (
	// general os stuff
	"fmt"

	// for search
	"mime"
	"miniwoofer/impl"
	"regexp"
	"strings"
)

type search_result struct {
	score float64
	meta  *impl.DocumentMeta
}

func main() {
	index := impl.NewBm25()
	config := impl.LoadConfig()

	// add mime extentions that may not exist by default
	mime.AddExtensionType(".mht", "multipart/related")
	mime.AddExtensionType(".mhtml", "multipart/related")

	db, err := impl.MetaDbOpen(config)
	if err != nil {
		panic(err)
	}

	fmt.Print("Attempt to load serialized BM25... ")
	if err := index.Load(config.Bm25Path()); err != nil {
		fmt.Println("Serialized BM25 Does not exist.")
		if err := impl.ParseCorpus(index, config); err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Success.")
	}

	fmt.Print("Checking if corpus changed... ")
	if changed, err := impl.CheckChanged(config.CorpusDir, index.CorpusHash); err != nil || changed {
		fmt.Print("Corpus Changed\n")
		if err := impl.ParseCorpus(index, config); err != nil {
			panic(err)
		}
		if err := db.AddCorpus(config); err != nil {
			panic(err)
		}
	} else {
		fmt.Print("No Change.\n")
	}

	// DO TEST QUERY
	var query string = "test"

	// search and populate search results
	re := regexp.MustCompile(`((?:".*?" *)|(?:(?:\S+? +)))`) // stemming ig?
	search_terms := []string{}
	for _, match := range re.FindAllStringSubmatch(query+" ", -1) {
		if len(match) > 1 {
			search_terms = append(search_terms, strings.Trim(strings.ToLower(match[1]), ` "`))
		}
	}

	// retrieve and rank docs w/inverted index and bm25
	results, err := index.Search(search_terms)
	if err != nil {
		panic(err)
	}

	// render the results
	if len(results) == 0 {
		fmt.Print("No Results for \"" + query + "\"")
	} else {
		// there are search results so you should make that list of titles and render them.
		// titles can be retrieved from the metadata db given the ids of the results

		// perform search
		// populate the model's search_results field
		search_results := []search_result{} // reset it first
		for _, s := range results {
			var sr search_result
			sr.score = s.Score

			// get title from metadata HERE IS WHERE PANIC HAPPENS
			doc_md, err := db.GetDocument(s.Id)
			if err != nil {
				panic(err)
			}
			sr.meta = doc_md
			// now put it in the model for state
			search_results = append(search_results, sr)
		}

	}

}
