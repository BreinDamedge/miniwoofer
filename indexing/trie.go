// Package indexing implements the index data structures
//
// This file holds the implementation of an index using a trie.
package indexing

import (
	"fmt"
)

/*
* for our inverted index we'll use unicode (utf-8) as our charset to build a trie.
* serializing this data structure won't be as simple as serializing a map/dict
* but it is certainly possible and will provide a fun challenge.
 */

type utf8TrieNode struct {
	r        rune                   // self rune, not sure we need this, but I'll leave it for now
	children map[rune]*utf8TrieNode // since runes are int32 the children "array" in memory needs to be a sparse structure
	docIds   []string
}

type Utf8Trie struct {
	children map[rune]*utf8TrieNode
	// maybe you want sizes and stuff but meh whatever for now
}

func (t *Utf8Trie) IsEmpty() bool {
	return t.children == nil
}

func (t *Utf8Trie) getOrMakeNext(c rune) *utf8TrieNode {
	if t.children == nil {
		fmt.Println("children was nil")
		// child doesn't exist so make it and initialize the children map
		t.children = map[rune]*utf8TrieNode{}
		t.children[c] = &utf8TrieNode{c, nil, nil}
		return t.children[c]
	}

	// children has been initialized so check if it exists
	if next, exists := t.children[c]; exists {
		return next
	}

	// child doesn't exist so make it
	t.children[c] = &utf8TrieNode{c, nil, nil}
	return t.children[c]
}

func (n *utf8TrieNode) getOrMakeNext(c rune) *utf8TrieNode {
	if n.children == nil {
		// child doesn't exist so make it and initialize the children map
		n.children = map[rune]*utf8TrieNode{}
		n.children[c] = &utf8TrieNode{c, nil, nil}
		return n.children[c]
	} else {
		// child may exist
		if next, exists := n.children[c]; exists {
			return next
		} else {
			// child doesn't exist so make it and initialize the children map
			n.children = map[rune]*utf8TrieNode{}
			n.children[c] = &utf8TrieNode{c, nil, nil}
			return n.children[c]
		}
	}
}

func (t *Utf8Trie) AddDoc(key string, docID string) {
	// for character in key: step to node which represents end of key, and add the given doc id
	if key == "" {
		panic("Cannot add docID to node with empty key")
	}

	// navigate to node which represents this keyword
	at := t.getOrMakeNext(rune(key[0]))
	for _, c := range key[1:] {
		// if children contains key c go to that node, otherwise create it
		at = at.getOrMakeNext(c)
	}

	// insert docID at this node (if slice is nil append handles it)
	at.docIds = append(at.docIds, docID)
}

func (t *Utf8Trie) AtKey(key string) []string {
	if key == "" {
		panic("missing key")
	}

	// navigate to node which represents this keyword
	at, exists := t.children[rune(key[0])]
	if !exists {
		fmt.Printf("key '%s' does not exist \n", key)
		return []string{}
	}
	// otherwise continue the climb
	for _, c := range key[1:] {
		// if children contains key c go to that node, otherwise given key isn't in index so return empty slice
		next, exists := at.children[c]
		if exists {
			at = next
		} else {
			fmt.Printf("key '%s' does not exist \n", key)
			return []string{}
		}
	}

	// if we've made it to a node, then we can return it's docIds
	return at.docIds
}

// DelKey deletes a key if it exists in the tree and clean up it's branch
func (t *Utf8Trie) DelKey(key string) {
	// traverse the tree down to the key, keeping track of previous node that is the end of a key (has docIds)
}

// TODO: .Keys method for trie which returns an iterator
