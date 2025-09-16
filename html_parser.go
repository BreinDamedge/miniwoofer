package main

import (
	"fmt"
	"os"
	"unicode"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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

// if this node is an anchor tag with a link, this function will return the link
// otherwise, an empty stream is returned
func isLink(n *html.Node) string {
	if n.Type == html.ElementNode && n.DataAtom == atom.A {
		for _, a := range n.Attr {
			if a.Key == "href" {
				return a.Val
			}
		}
	}
	return ""
}

func addNode2DocData(n *html.Node, data *DocData) {
	// do parsing here
	// get data out of this node
	// ...
	// !!!srcURL will probably only be present in archived files if they're mhtml!!!
	// if outlink
	if link := isLink(n); link != "" {
		// this is a link tag, determine if it is an outlink
		data.outlinks = append(data.outlinks, link)
	}

	// if title
	// if contains (text-content)
}

func HTMLFile2DocData(docPath string) DocData {
	data := DocData{}

	file, err := os.Open(docPath) // For read access.
	check(err)
	doc, err := html.Parse(file)
	check(err)

	// process data in doc (head/root)
	addNode2DocData(doc, &data)

	// get data from child nodes
	// ...
	// use n.Decendents() to get an iterator and you won't need to recurse
	for n := range doc.Descendants() {
		addNode2DocData(n, &data)
	}

	return data
}
