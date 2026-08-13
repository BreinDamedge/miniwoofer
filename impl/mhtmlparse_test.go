package impl

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	x, err := os.ReadFile(filepath.Join("..", "corpus", "vlm.mht"))
	if err != nil {
		t.Fatalf("read corpus: %v", err)
	}

	r := bytes.NewReader(x)
	tokens, err := TokenizeMhtml(r)

	if err != nil {
		t.Fatalf("tokenize: %v", err)
	}

	if len(tokens) == 0 {
		t.Fatalf("no tokens")
	}
}
