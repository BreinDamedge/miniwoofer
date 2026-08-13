package impl

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
)

// bm25 default parameters
const (
	defaultK1 = 1.5
	defaultB  = 0.75
)

type bm25Posting struct {
	Id string
	Tf int
}

// holds metadata for a doc
type bm25Meta struct {
	Len int
}
type Bm25 struct {
	K1, B float64

	TotalDocLen int
	AvgDocLen   float64

	Metas    map[string]bm25Meta
	Postings map[string][]bm25Posting
	Freq     map[string]int

	CorpusHash string // Should probably be put elsewhere
}

func NewBm25() *Bm25 {
	b := &Bm25{}
	b.K1 = defaultK1
	b.B = defaultB
	b.Metas = make(map[string]bm25Meta)
	b.Postings = make(map[string][]bm25Posting)
	b.Freq = make(map[string]int)

	return b
}

// TODO: change tokens to ints, and have a map of int -> str for token id -> string representation
func (b *Bm25) Append(id string, tokens []string) error {
	if _, ok := b.Metas[id]; ok {
		return fmt.Errorf("document %s already ok", id)
	}

	// count term frequencies for this doc
	// TODO: consider storing these, so that index freq can be updated on document removal w/out re-parsing. Alternatively accept the tradeoff in cpu time for memory savings
	freq := make(map[string]int)
	for _, token := range tokens {
		freq[token]++
	}

	// record the new document length, and update index stats
	b.Metas[id] = bm25Meta{Len: len(tokens)}
	b.TotalDocLen += len(tokens)
	b.AvgDocLen = float64(b.TotalDocLen) / float64(len(b.Metas))

	// update the term frequencies of the entire corpus
	for term, tf := range freq {
		b.Freq[term]++
		b.Postings[term] = append(b.Postings[term], bm25Posting{
			Id: id,
			Tf: tf,
		})
	}

	return nil
}

func (b *Bm25) Search(query []string) ([]Search, error) {
	if len(b.Metas) == 0 {
		return nil, nil
	}

	scores := make(map[string]float64)
	for _, term := range query {
		idf := b.calculateIdf(term)

		for _, posting := range b.Postings[term] {
			meta := b.Metas[posting.Id]
			scores[posting.Id] += b.calculateScore(idf, posting.Tf, meta.Len)
		}
	}

	if len(scores) == 0 {
		return nil, nil
	}

	ranked := make([]Search, 0, len(scores))
	for id, score := range scores {
		ranked = append(ranked, Search{Id: id, Score: score})
	}

	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Score == ranked[j].Score {
			return ranked[i].Id < ranked[j].Id
		} else {
			return ranked[i].Score > ranked[j].Score
		}
	})

	return ranked, nil
}

func (b *Bm25) calculateIdf(term string) float64 {
	total := float64(len(b.Metas))
	freq := float64(b.Freq[term])
	return math.Log1p((total - freq + 0.5) / (freq + 0.5))
}

func (b *Bm25) calculateScore(idf float64, tf int, docLen int) float64 {
	if idf == 0 || tf == 0 || b.AvgDocLen == 0 {
		return 0
	}

	numerator := float64(tf) * (b.K1 + 1)
	denominator := float64(tf) + b.K1*(1-b.B+b.B*(float64(docLen)/b.AvgDocLen))

	if denominator != 0 {
		return idf * numerator / denominator
	} else {
		return 0
	}
}

func (b *Bm25) SetB(newB float64) {
	b.B = newB
}

func (b *Bm25) SetK1(newK1 float64) {
	b.K1 = newK1
}

func (b *Bm25) Params() (float64, float64) {
	return b.B, b.K1
}

func (b *Bm25) SetParams(cfg Bm25Parameters) {
	b.B = cfg.B
	b.K1 = cfg.K1
}

// Saves the bm25 to a json file, pretty simple
func (b *Bm25) Save(path string) error {
	jsonData, err := json.Marshal(*b)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return err
	}

	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s", string(jsonData))
	if err != nil {
		return err
	}
	return nil
}

// Inverse of Save
func (b *Bm25) Load(path string) error {
	jsonData, err := os.ReadFile(path)

	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, b)
}
