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

func (t *Utf8Trie) getOrMakeNext(c rune) *utf8TrieNode {
	// child may exist
	if next, exists := t.children[c]; exists {
		return next
	} else {
		if t.children == nil {
			// children uninitialized so initialize it
			t.children = map[rune]*utf8TrieNode{}
		}
		// child doesn't exist so make it
		t.children[c] = &utf8TrieNode{c, nil, nil}
		return t.children[c]
	}
}

func (n *utf8TrieNode) getOrMakeNext(c rune) *utf8TrieNode {
	// child may exist
	if next, exists := n.children[c]; exists {
		return next
	} else {
		if n.children == nil {
			// children uninitialized so initialize it
			n.children = map[rune]*utf8TrieNode{}
		}
		// child doesn't exist so make it
		n.children[c] = &utf8TrieNode{c, nil, nil}
		return n.children[c]
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
		return []string{}
	}
	// otherwise continue the climb
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

func (n *utf8TrieNode) delChild(key rune) {
	delete(n.children, key)
	if len(n.children) > 0 {
		n.children = nil
	}
}

func (t *Utf8Trie) delChild(key rune) {
	delete(t.children, key)
	if len(t.children) > 0 {
		t.children = nil
	}
}

func (t *Utf8Trie) pruneBranch(trajectory []*utf8TrieNode) {
	// given the trajectory walk backward through it and delete nodes until the branch is pruned
	var child *utf8TrieNode
	for i := len(trajectory) - 1; i >= 0; i-- {
		n := trajectory[i]
		// if there was a leaf then delete it
		if child != nil {
			n.delChild(child.r)
		}
		if n.docIds == nil && n.children == nil {
			// continue up the tree and delete the
			child = n
		} else {
			return
		}
	}

	// at this point you made it allll the way up the trajectory to the root,
	// so it should be a child of t so you need to delete the child of t
	if child != nil {
		// first check if child is empty
		if child.children == nil && child.docIds == nil {
			// delete final child from Trie
			t.delChild(child.r)
		}
	}
}

func (n *utf8TrieNode) delID(id string) {
	// TODO: sort docIds
	// if you sort the ids in nodes then you can binary search them so that removals and checks are faster

	// find the id ur removing
	ir := -1
	for i, v := range n.docIds {
		if v == id {
			ir = i
			break
		}
	}

	// if id was found remove the element
	if ir != -1 {
		if ir < len(n.docIds)-1 && ir > 0 {
			n.docIds = append(n.docIds[:ir], n.docIds[ir+1:]...)
		} else {
			// otherwise ur removing either the front or back element so just chop the one element
			if ir == len(n.docIds)-1 {
				n.docIds = n.docIds[:ir]
			} else {
				n.docIds = n.docIds[1:]
			}
		}

		// if slice is now empty, set it nil
		if len(n.docIds) == 0 {
			n.docIds = nil
		}
	}
}

// DelKey deletes a key if it exists in the tree and clean up it's branch
func (t *Utf8Trie) DelKey(key string) {
	// traverse the tree down to the key, keeping track of previous node that is the end of a key (has docIds)
	/*
	* steps:
	* 1. climb to the node of key
	* 2. remove key from node.docIds
	* 3. cleanup tree
	* 	if node.docIds len() == 0
	* 		node is empty so remove doc ids (set nil)
	* 		if node is leaf (has no children)
	* 			prune branch
	* 				how tf do you prune a branch in this god forsaken data structure?
	 */

	// traj := []string{}

	// climb to node
	// ...
}

// TODO: .Keys method for trie which returns an iterator
