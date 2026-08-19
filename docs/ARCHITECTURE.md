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

# Project File Structure
```
.metadata/
cmd/
corpus/
docs/
impl/
config.toml
queries.toml
```

### .metadata/
.metadata is to contain the serialized index files, and sqlite data. (so far `bm25.json` & `database.db`)

### corpus/
All indexed documents go here.

### cmd/
this folder contains the .go files which are expected to be run when using this software.

# Optimizers
not yet...
