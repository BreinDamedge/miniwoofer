// Package tests has some tests for now
package tests

import (
	"fmt"

	"miniwoofer/indexing"
)

func TestTrie() {
	fmt.Println("Creating trie...")
	t := indexing.MakeUtf8Trie()
	fmt.Println("Adding document at key 'test'...")
	t.AddDoc("test", "test_id")
	fmt.Println("Retrieving document at key 'test'...")
	result := t.AtKey("test")
	fmt.Println(result)
	fmt.Println("Adding another id to key 'test'...")
	t.AddDoc("test", "id2")
	fmt.Println("Retrieving documents at key 'test'...")
	result = t.AtKey("test")
	fmt.Println(result)
	fmt.Println("Done.")
}
