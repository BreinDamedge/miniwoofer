package main

import (
	"fmt"
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
