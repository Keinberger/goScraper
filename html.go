package scraper

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// GetHTML returns the HTML data of URL
func GetHTML(URL string) (string, error) {
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

// GetHTMLNode returns the node tree of the html string data
func GetHTMLNode(data string) (*html.Node, error) {
	return html.Parse(strings.NewReader(data))
}

// GetElementNodes returns an array of html.Node iniside of htmlNode having the same properties as element e
func (e *Element) GetElementNodes(htmlNode *html.Node) ([]*html.Node, error) {
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
			if len(e.Tags) == foundTags { // has found all tags of node
				elements = append(elements, node)
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling { // looking for nodes inside of all other nodes
			elements = append(elements, crawler(child)...)
		}
		return elements
	}
	if el := crawler(htmlNode); len(el) > 0 {
		return el, nil
	}
	return nil, newErr(ErrMissingElement, "missing "+e.Typ+" in the node tree")
}

// RenderNode returns the string representation of an html.Node
func RenderNode(node *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, node)
	return buf.String()
}

// GetTextOfNode returns the content of an html element
func GetTextOfNode(node *html.Node, notRecursive bool) (text string) {
	if node.Type == html.TextNode {
		text += node.Data
	}
	if !notRecursive {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			text += GetTextOfNode(c, notRecursive)
		}
	}
	return
}
