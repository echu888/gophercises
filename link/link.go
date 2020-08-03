// implementation of Gophercise Exercise #4
//   https://github.com/gophercises/link/
//
// primary features:
//	use x/net/html for parsing html
//	use go style unit testing
package main

import (
	"flag"
	"fmt"
	"os"
	//"io"
	//"bytes"
	"golang.org/x/net/html"
	"log"
	"strings"
)

//const (
// 0   ErrorNode NodeType = iota
// 1   TextNode
// 2   DocumentNode
// 3   ElementNode
// 4   CommentNode
// 5   DoctypeNode
//    // RawNode nodes are not returned by the parser, but can be part of the
//    // Node tree passed to func Render to insert raw HTML (without escaping).
//    // If so, this package makes no guarantee that the rendered HTML is secure
//    // (from e.g. Cross Site Scripting attacks) or well-formed.
//    RawNode
//)

type Link struct {
	Href string
	Text string
}

type Links []Link

// convert a node into a string
//func renderNode(node *html.Node) string {
//    var buf bytes.Buffer
//    w := io.Writer(&buf)
//    html.Render(w, node)
//    return buf.String()
//}

// grab all content of a node, filtering out <tags> and <!-- comments -->
func getContent(node *html.Node) string {
	var content string

	// ignores <all_elements/> and <!-- all comments -->
	if node.Type != html.ElementNode &&
		node.Type != html.CommentNode {
		content = node.Data
	}

	//log.Println("getContent:", content)
	var postcontent string

	// recursively process all content available in this node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		//log.Printf("-- CHILD OF %s\n", node.Data)
		postcontent = postcontent + getContent(child)
	}
	return content + postcontent
}

func findLinks(root *html.Node) Links {
	links := []Link{}
	var findLinkFunc func(*html.Node)
	findLinkFunc = func(node *html.Node) {
		//if node.Type == html.ElementNode {
		//log.Printf("node: [%d] [%s]\n", node.Type, node.Data)
		//}
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, a := range node.Attr {
				//log.Println(" -- " , a)
				if a.Key == "href" {
					var link = Link{a.Val, strings.TrimSpace(getContent(node))}
					//log.Println(node)
					links = append(links, link)
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			//log.Printf("-- CHILD OF %s\n", node.Data)
			findLinkFunc(child)
		}
	}
	findLinkFunc(root)
	return links
}

func checkErrors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func loadFile(filename string) *html.Node {
	fmt.Println("Loading file:", filename)
	file, err := os.Open(filename)
	checkErrors(err)

	doc, err := html.Parse(file)
	checkErrors(err)

	return doc
}

func displayLinks(links Links) {
	fmt.Println("Links found:")
	for _, link := range links {
		fmt.Printf("Link: [%s] %s\n", link.Href, link.Text)
	}
}

func main() {
	filenamePtr := flag.String("f", "example.html", "html file to parse for links")
	flag.Parse()

	doc := loadFile(*filenamePtr)
	links := findLinks(doc)
	displayLinks(links)
}
