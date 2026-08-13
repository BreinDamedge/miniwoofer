package impl

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

/*
* pseudocode from wiki:
* -----------------------------------------------------
* let s = s_0
* for k = 0 through k_max (exclusive)
* 	T <- temperature(1-(k+1)/k_max)
* 	Pick a random neighbour, s_new <- neighbour(s)
* 	if P(E(s), E(s_new), T) >= random(0,1):
* 		s <- s_new
* output the final state s
* -----------------------------------------------------
*
* alrighty so questions:
* - how are you picking neighbours?
* 	- since each candidate solution is a pair b, k, each candidate will be a pair b, k. set a min and max bound for each and a step size.
* - what is P(E(s), E(s_new), T)?
* - how are you controlling tempurature (what function is it)?
*
*
*
*
* what does your
*
*
*
* fit queries is a map loaded in from toml
* where each key is a doc id and at that position is ig another map and each of those maps has one entry which is "queries" and the type of that field is a []string
*
 */

type Bm25Parameters struct {
	B  float64
	K1 float64
}

func LoadCfg() Bm25Parameters {
	var conf Bm25Parameters

	// DecodeFile parses the file directly into the struct pointer
	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		panic(err)
	}

	return conf
}

type helpMe2 struct {
	Queries []string `toml:"queries"`
}
type helpMe struct {
	Docs map[string]helpMe2 `toml:"docs"`
}

func LoadQueryData() map[string][]string {
	var neoVim12 helpMe

	// DecodeFile parses the file directly into the struct pointer
	if _, err := toml.DecodeFile("./queries.toml", &neoVim12); err != nil {
		panic(err)
	}

	out := make(map[string][]string)
	for k, v := range neoVim12.Docs {
		out[k] = v.Queries
	}

	return out
}

// gridsearch optimizer
func Optimize(b *Bm25, docs []Doc, queries map[string][]string) Bm25Parameters {
	// how do we optimize now...
	// for each thing in tuning Data we
	//	run the search query
	//	linear search for a doc w/the target id until we find it and then add it's normalized score to the current cost (if not found, cost = 1)(preferably faster than linear and also we cache tfidf results bc we're going to make the same search multiple times)
	//	find combo of b&k1 w/ minimum cost

	// first tokenize queries to avoid duplicate work
	tok_queries := make(map[string][][]string)
	for id, qs := range queries {
		// id is doc id, qs are the queries
		tok_queries[id] = [][]string{}
		for _, q := range qs {
			tok_queries[id] = append(tok_queries[id], Tokenize(q))
		}
	}

	// optimier params
	B, K1 := b.Params()
	best_cfg := Bm25Parameters{B, K1}
	best_cost := -1.
	done := false

	// k1 range 1.2 - 2.0
	// b range 0.6 - 0.8
	step_size := 0.05
	for cur_b := 0.6; cur_b <= 0.801; cur_b += step_size {
		for cur_k1 := 1.2; cur_k1 <= 2.001; cur_k1 += step_size {
			// remember to set the params
			b.SetB(cur_b)
			b.SetK1(cur_k1)

			// for each doc in tuning data, and each query there, run a search and calculate the total cost of the optimizer under current params
			var cost float64 = 0.0
			for id, qs := range tok_queries {
				// id is doc id, qs are the queries
				for _, q := range qs {
					// run a search
					result, err := b.Search(q)
					if err != nil {
						panic(err)
					}
					// get normalized position in the optimizer
					num_retrieved := len(result)
					if num_retrieved == 0 {
						cost += 1 // doc was not successfully retrieved
					} else {
						var retrieved bool = false
						for i, x := range result {
							if x.Id == id {
								cost += float64(i) / float64(num_retrieved)
								retrieved = true
								break
							}
						}
						// doc was not successfully retrieved
						if retrieved == false {
							cost += 1
						}
					}
				}
			}

			// pull current params
			B, K1 := b.Params()

			// potentially update best_cfg
			if best_cost == -1 || cost < best_cost {
				best_cfg = Bm25Parameters{B, K1}
				best_cost = cost
			}

			// print for fun
			// fmt.Printf("Cost test: current params (%f, %f) yeild cost %f\n", B, K1, cost)

			// early exit
			if cost == 0.0 {
				fmt.Println("you're so perfect gang")
				done = true
				break
			}
		}
		if done {
			break
		}
	}
	fmt.Printf("Optimizer: End Params (%f, %f) have cost %f\n", best_cfg.B, best_cfg.K1, best_cost)
	return best_cfg
}

func Evaluate(b *Bm25, docs []Doc, queries map[string][]string) float64 {
	// first tokenize queries to avoid duplicate work
	tok_queries := make(map[string][][]string)
	for id, qs := range queries {
		// id is doc id, qs are the queries
		tok_queries[id] = [][]string{}
		for _, q := range qs {
			tok_queries[id] = append(tok_queries[id], Tokenize(q))
		}
	}

	// for each doc in tuning data, and each query there, run a search and calculate the total cost of the optimizer under current params
	var cost float64 = 0.0
	for id, qs := range tok_queries {
		// id is doc id, qs are the queries
		for _, q := range qs {
			// run a search
			result, err := b.Search(q)
			if err != nil {
				panic(err)
			}
			// get normalized position in the optimizer
			num_retrieved := len(result)
			if num_retrieved == 0 {
				cost += 1 // doc was not successfully retrieved
			} else {
				var retrieved bool = false
				for i, x := range result {
					if x.Id == id {
						cost += float64(i) / float64(num_retrieved)
						retrieved = true
						break
					}
				}
				// doc was not successfully retrieved
				if retrieved == false {
					cost += 1
				}
			}
		}
	}
	return cost
}
