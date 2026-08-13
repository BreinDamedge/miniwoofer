# Boogaloo -> miniwoofer
This is a BM25 based search engine for mhtml documents created as a class project. Modifying it (fully rewriting) to eventually become miniwoofer.

# Running The Project:
1. ensure golang is installed
1. from the root directory of the project running `go run ./cmd/main.go` will compile and start the project (replace `run` with `build` to compile to an executable instead)

# High Level Architecture:
SQLite database baybeeeee.
data that needs storing in sqlite:
each document has:
- title
- id
- keywords
- content
- tags (optional)

inverted index:
this is handled currently by the bm25 data structure


# Todo:
### BM25 Index Serialization
- [x] serialize the index data structure lol.

### Title Extraction
automatic extraction of document titles when parsing MHTML documents. The titles should be a seprate field stored. Use sqlite or some other serializable data structure to store document titles in a way that can be retrieved given a document ID.  
This also leads nicely into an eventual need for some data structure that can hold arbitrary metadata about a document (that's why i suggest sqlite, tho a custom one is fine too since documents may not all have the same metadata).
- [ ] sqlite storage for document metadata

### Whole Webserver Situation
This entire thing needs a simple frontend with these capabilities:
see `impl/webserver.go`
- [x] trigger a text query
- [x] display document titles in ranked order
  - [ ] cut them off after some limit (configurable but hardcoding to start is fine and dandy)
- [ ] trigger a rescan of the corpus.
  - [ ] update the index & metadata store for any new or removed documents

### Testing Framework (For Evaluating Configurability Of Recommender)
- [ ] design hybrid recommendation system
Eventually having a way to make a dataset of queries and documents that the user wants those queries to retrieve is needed. this kind of already exists in `impl/optim.go`.
Optimization & tuning of the recommender system will become the focus after the rest of the system is working. Then we can focus on configurability of the search algorithm.
