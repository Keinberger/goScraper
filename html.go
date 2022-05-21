package scraper

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// GetHTMLBody returns the HTML-body of the html data at URL
func GetHTMLBody(URL string) (string, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	return string(html), nil
}

// GetHTMLNode returns a node tree representing html data at URL
func GetHTMLNode(URL string) (*html.Node, error) {
	htmlData, err := GetHTMLBody(URL)
	if err != nil {
		return nil, err
	}
	return html.Parse(strings.NewReader(htmlData))
}

// GetElementNodes returns an array of html.Node iniside of doc having the same attributes as element E
func (e *Element) GetElementNodes(doc *html.Node) ([]*html.Node, error) {
	var crawler func(*html.Node) []*html.Node
	crawler = func(node *html.Node) (elements []*html.Node) {
		if node.Type == html.ElementNode && node.Data == e.Typ {
			var foundTags int
			for _, tag := range e.Tags {
			C:
				for _, attr := range node.Attr {
					if attr.Key == tag.Typ && attr.Val == tag.Value {
						foundTags += 1
						break C
					}
				}
			}
			if len(e.Tags) == foundTags { // has found all tags, return node
				return append(elements, node)
			}
		}
		if child := node.FirstChild; child != nil {
			elements = append(elements, crawler(child)...)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			elements = append(elements, crawler(child)...)
		}
		return elements
	}
	if el := crawler(doc); len(el) > 0 {
		return el, nil
	}
	return nil, errors.New("missing " + e.Typ + " in the node tree")
}

// RenderNode converts a rendered html into string
func RenderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

// GetTextOfNode returns the inside of an html element
func GetTextOfNode(node *html.Node, notRecursive bool) string {
	var finished string
	if node.Type == html.TextNode {
		finished += node.Data
	}
	if !notRecursive {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			finished += GetTextOfNode(c, notRecursive)
		}
	}
	return finished
}
