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

// GetHTMLdata returns the HTML-body of a website
func GetHTMLdata(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		logError(err, "Could not fetch data from: "+url)
		return ""
	}
	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	logError(err, "Could not get HTML data from: "+url)

	return string(html)
}

// GetHTMLElementContent returns the content of an HTML element including the HTML tags of that element
func GetHTMLElementContent(doc *html.Node, el Element) (*html.Node, string, error) {
	var content *html.Node
	var container string
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == el.Typ {
			if el.Tag != (Tag{}) {
				for _, e := range node.Attr {
					if e.Key == el.Tag.Name && e.Val == el.Tag.Value {
						container = `<` + el.Typ + ` ` + el.Tag.Name + `="` + el.Tag.Value + `"`
						content = node

						for _, e := range node.Attr {
							if !strings.Contains(container, e.Key+`="`+e.Val+`"`) {
								container += " " + e.Key + `="` + e.Val + `"`
							}
						}
						container += `>`
						break
					}
				}
			} else {
				container = `<` + node.Data + `>`
				content = node
			}
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if content != nil {
		return content, container, nil
	}
	return nil, container, errors.New("Missing " + el.Typ + " with: " + el.Tag.Name + "='" + el.Tag.Value + "' in the node tree")
}

// renderNode converts a rendered html into string
func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

// GetContentInsideElement returns the inside of an html element
func GetContentInsideElement(htm string, el Element) string {
	doc, _ := html.Parse(strings.NewReader(htm))
	bn, beg, err := GetHTMLElementContent(doc, el)
	if err != nil {
		logError(err, "Error getting html element: ")
		return ""
	}
	return strings.Trim(renderNode(bn), beg)
}

// GetNestedHTMLElement returns the content of a html element inside of a bigger one
func GetNestedHTMLElement(htm string, elements []Element) string {
	var output string
	for _, v := range elements {
		if output != "" {
			output = GetContentInsideElement(output, v)
		} else {
			output = GetContentInsideElement(htm, v)
		}
	}
	return output
}
