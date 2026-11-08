# miniwoofer
This is a BM25 based search engine for mhtml documents. It is accessible through an api hosted on a webserver on device

# uWoofer
microwoofer uses the same underlying code as miniwoofer, but replaces the web interface with a very basic tui, and opens the files directly on the host system instead of parsing them to html for a web browser. It is for an even more minimal footprint, and also is a secondary interface for testing ranking algorithms quickly. Further development may or may not happen, this was more of an experiment.

# Search Configuration
...

# Running The Project:
1. ensure golang is installed
1. from the root directory of the project running `go run ./cmd/miniwoofer` will compile and start the project (replace `run` with `build` to compile to an executable instead)
> [!NOTE]
> If you are getting errors when you try to search, start by deleting your `.metadata` folder and reinitializing the program.

# Test Corpus
Test corpus has been thrown together. Any files needed for parsing tests should be added to this.  
You can find it [here](https://drive.google.com/file/d/15CUxD-218IN0oBTnLGVwxu5gsAacGfqU/view?usp=sharing)  

# High Level Architecture:
... write this later  

# Attribution
Gotta attribute to w3schools if we use [this](https://www.w3schools.com/css/tryit.asp?filename=trycss_template3) CSS template for markdown document rendering. I assume we'll make our own but I'm not feeling artistic atm so template for testing.



