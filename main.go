package main

import (
	"miniwoofer/tests"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	tests.TestTrie()
}
