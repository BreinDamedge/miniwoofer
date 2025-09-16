package main

import (
	"fmt"
	"os"
	"unicode"

	"golang.org/x/net/html"
)

func isWhiteSpace(s_ string) bool {
	if len(s_) == 0 {
		return true
	}

	for i := 0; i < len(s_); i++ {
		if !unicode.IsSpace(rune(s_[i])) {
			return false
		}
	}
	return true
}

func addSpace(s_ string) string {
	if !isWhiteSpace(s_) {
		if len(s_) > 0 {
			if string(s_[len(s_)-1]) != " " {
				fmt.Printf("Adding ' ' to '%s'\n", s_)
				return s_ + " "
			}
			return s_
		}
	}
	return ""
}

// process the details of the product within the <li> element
func processNode(n *html.Node) string {
	text := ""
	// if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
	// 	// if yes, retrieve FirstChild's data
	// 	text += addSpace(n.FirstChild.Data)
	// }
	switch n.Data {
	case "body":
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			// if yes, retrieve FirstChild's data (name)
			name := n.FirstChild.Data
			// print name
			fmt.Println("Name:", name)
		}

	case "h1":
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			// if yes, retrieve FirstChild's data
			text += addSpace(n.FirstChild.Data)
		}

	case "h2":
		// check if FirstChild node of the h2 element is a text
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			// if yes, retrieve FirstChild's data
			text += addSpace(n.FirstChild.Data)
		}

	case "p":
		if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			// if yes, retrieve FirstChild's data (name)
			name := n.FirstChild.Data
			// print name
			// fmt.Println("Name:", name)
			text += name
		}
	}

	// Traverse child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += addSpace(processNode(c))
	}
	return text
}

type DocData struct {
	srcURL             string
	outlinks           []string
	title, textContent string
}

func addNode2DocData(n *html.Node, data *DocData) {
	// do parsing here
	// get data out of this node
	// ...
	// !!!srcURL will probably only be present in archived files if they're mhtml!!!
	// if outlink
	// if title
	// if contains (text-content)

	// get data from child nodes
	// ...
}

func HTMLFile2DocData(docPath string) DocData {
	data := DocData{}

	file, err := os.Open("golanghtml.html") // For read access.
	check(err)
	doc, err := html.Parse(file)
	check(err)

	// recursively process data in doc (head/root) and children until DocData is complete
	addNode2DocData(doc, &data)

	return data
}
