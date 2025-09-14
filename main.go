package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

func main() {
	//.. HTTP request
	file, err := os.Open("testfile.html") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	doc, err := html.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	// find all <li> elements
	var processAllProduct func(*html.Node) string
	processAllProduct = func(n *html.Node) string {
		// fmt.Println(n)
		// fmt.Println(n.Type)
		// fmt.Println(n.Data)
		// fmt.Println()
		text := ""
		if n.Type == html.ElementNode && n.Data == "body" {
			// process the Product details within each <li> element
			newText := processNode(n)
			if newText != "" {
				text += newText + " "
			}
		}
		// traverse the child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			newText := processAllProduct(c)
			if newText != "" {
				text += newText + " "
			}
		}
		return text
	}
	// make a recursive call to your function
	fmt.Println(processAllProduct(doc))
}
