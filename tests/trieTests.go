// Package tests has some tests for now
package tests

import (
	"fmt"

	"miniwoofer/indexing"
)

func TestTrie() {
	fmt.Println("Creating trie...")
	t := indexing.Utf8Trie{}
	fmt.Println("adding doc at key 'a'...")
	t.AddDoc("a", "a test")
	fmt.Println("Check if t.children is empty...")
	fmt.Println(t.IsEmpty())
	fmt.Println("retrieving doc at 'a'...")
	fmt.Println(t.AtKey("a"))

	fmt.Println("adding doc at key 'a' again...")
	t.AddDoc("a", "a test")
	fmt.Println("retrieving doc at 'a' again...")
	fmt.Println(t.AtKey("a"))

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
