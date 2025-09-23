package indexing

/*
* mapdex datastructure is an index data structure that contains both an inverted index and a forward index.
* it should also avoid duplicating keys by storing them all once and referencing them with ids (ints)
*
* whiteboard here plz... design it
*
* so this structure needs a forward index, and inverted index, and a set of known docIds and kwds
* is this also all lame and expensive? why!?!?!? why what has happened?
 */

type mapdex struct {
	inv [][]int // kwd -> docId
	fwd [][]int // docId -> kwd
	// for memory efficiency you should store idx instead of kwd replicas
	docIds []string
	kwds   []string // might even be a more efficient way to store all this in some kind of tree or w/some encoding
}
