# miniwoofer
This is a BM25 based search engine for mhtml documents.

# Running The Project:
1. ensure golang is installed
1. from the root directory of the project running `go run ./cmd/main.go` will compile and start the project (replace `run` with `build` to compile to an executable instead)
> [!NOTE]
> If you are getting errors when you try to search, start by deleting your `.metadata` folder and reinitializing the program.

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


