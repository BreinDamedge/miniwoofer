// Package indexing implements the index data structures
//
// This file holds the implementation of an index using a trie.
package indexing

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

func MakeUtf8Trie() Utf8Trie {
	t := Utf8Trie{}
	t.children = map[rune]*utf8TrieNode{}
	return t
}

func makeUtf8TrieNode() utf8TrieNode {
	n := utf8TrieNode{}
	n.children = map[rune]*utf8TrieNode{}
	n.docIds = []string{}
	return n
}

func (t Utf8Trie) getOrMakeNext(c rune) *utf8TrieNode {
	if next, exists := t.children[c]; exists {
		return next
	} else {
		newChild := makeUtf8TrieNode()
		t.children[c] = &newChild
		next = t.children[c]
		next.r = c
		return next
	}
}

func (n utf8TrieNode) getOrMakeNext(c rune) *utf8TrieNode {
	if next, exists := n.children[c]; exists {
		return next
	} else {
		newChild := makeUtf8TrieNode()
		n.children[c] = &newChild
		next = n.children[c]
		next.r = c
		return next
	}
}

func (t Utf8Trie) AddDoc(key string, docID string) {
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

	// insert docID at this node
	at.docIds = append(at.docIds, docID)
}

func (t Utf8Trie) AtKey(key string) []string {
	if key == "" {
		panic("missing key")
	}

	// navigate to node which represents this keyword
	at := t.getOrMakeNext(rune(key[0]))
	for _, c := range key[1:] {
		// if children contains key c go to that node, otherwise given key isn't in index so return empty slice
		next, exists := at.children[c]
		if exists {
			at = next
		} else {
			return []string{}
		}
	}

	// if we've made it to a node, then we can return it's docIds
	return at.docIds
}

// TODO: .Keys method for trie which returns an iterator
// TODO: index interface definition
