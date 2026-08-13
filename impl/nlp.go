package impl

import (
	"embed"
	"strings"
	"unicode"
)

// jacob that next line is a pragma I think
//
//go:embed assets
var assets embed.FS

func wordSetFrom(path string) map[string]struct{} {
	corpus, err := assets.ReadFile(path)
	if err != nil {
		panic(err)
	}

	result := make(map[string]struct{})
	for _, line := range strings.Split(string(corpus), "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		}

		result[strings.TrimSpace(line)] = struct{}{}
	}

	return result
}

var stopWordSet = wordSetFrom("assets/stopwords.txt")

func lower(tokens []string) []string {
	for i := range tokens {
		tokens[i] = strings.ToLower(tokens[i])
	}

	return tokens
}

func removeStopwords(tokens []string) []string {
	n := 0

	for _, w := range tokens {
		if _, ok := stopWordSet[w]; !ok {
			tokens[n] = w
			n++
		}
	}

	return tokens[:n]
}

func Tokenize(sentence string) []string {
	tokens := strings.FieldsFunc(sentence, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	tokens = lower(tokens)
	tokens = removeStopwords(tokens)

	return tokens
}
