package main

import "fmt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	data := HTMLFile2DocData("golanghtml.html")
	// fmt.Println(data.outlinks)
	fmt.Println(data.title)
}
