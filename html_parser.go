package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func addSpace(s_ string) string {
	if s_ != "" {
		return s_ + " "
	}
	return s_
}

// process the details of the product within the <li> element
func processNode(n *html.Node) string {
	text := ""
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

	case "span":
		// check for the span with class "amount"
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "amount") {
				// retrieve the text content of the "amount" span
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						// print product price
						fmt.Println("Price:", c.Data)
					}
				}
			}
		}

	case "img":
		// check for the src attribute in the img tag
		for _, a := range n.Attr {
			if a.Key == "src" {
				// retrieve src value
				ImageURL := a.Val
				// print image URL
				fmt.Println("Image URL:", ImageURL)
			}
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
		text += processNode(c)
	}
	return text
}
