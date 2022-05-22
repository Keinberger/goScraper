<!-- implement 
build_status, (https://travis-ci.org) 
godoc, (https://godoc.org), (https://pkg.go.dev)
go_report_card here (goreportcard.com)
-->


# goScraper

goScraper is a small web-scraping library for Go.

## Installation

Package can be installed manually using <br />
```go
go get github.com/keinberger/goScraper
```

But may also be normally imported when using go modules<br />
```go
import "github.com/keinberger/goScraper"
```

## Usage

The package provides several exported functions to provide high functionality.<br />
However, the main scrape functions 
```go
func (w Website) Scrape(funcs map[string]interface{}, vars ...interface{}) (string, error)
```
```go
func (el lookUpElement) ScrapeTreeForElement(node *html.Node) (string, error)
```
```go
func (e *Element) GetElementNodes(doc *html.Node) ([]*html.Node, error)
```
should be the preffered way to use the scraper library.

As these functions use the other exported functions, as well, it provides all the features of the library packed together
and guided by only having to provide a minimal amount of input. For the main `Scrape()` function, the user input is scoped to only having to provide a custom Website variable.

### Example using `Scrape()`

This example provides a tutorial on how to scrape a website for specific html elements. The html elements will be returned chained-together, seperated by a custom seperator.

The example will use a custom website variable, where the `Scrape()` function will be called upon. The arguments of the `Scrape()` function are optional and will not be needed in this example.
```go
package main

import (
	"fmt"
	"github.com/keinberger/goScraper"
)

func main() {
	website := scraper.Website{
		URL: "https://wikipedia.org/wiki/wikipedia",
		LookUpElements: []scraper.LookUpElement{
			{
				Element: scraper.Element{
					Typ: "h1",
					Tags: []scraper.Tag{
						{
							Typ:   "id",
							Value: "firstHeading",
						},
					},
				},
			},
			{
				Element: scraper.Element{
					Typ: "td",
					Tags: []scraper.Tag{
						{
							Typ:   "class",
							Value: "infobox-data",
						},
					},
				},
				Index: 0,
			},
		},
		Seperator: ", ",
	}

	scraped, err := website.Scrape(nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(scraped)
}
```

### Example using `ScrapeTreeForElement()`
This example will use ScrapeTreeForElement, which will return the content of an html element (*html.Node) inside of a bigger node tree. This function is especially useful, if one only wants one html element from a website, but still wants to retain control over formatting settings.
```go
package main

import (
	"fmt"
	"github.com/keinberger/scraper"
)

func main() {
	htmlNode, err := scraper.GetHTMLNode("https://wikipedia.org/wiki/wikipedia")
	if err != nil {
		panic(err)
	}

	lookUpElement := scraper.LookUpElement{
		Element: scraper.Element{
			Typ: "li",
			Tags: []scraper.Tag{
				{
					Typ:   "id",
					Value: "ca-viewsource",
				},
			},
		},
	}
	content, err := lookUpElement.ScrapeTreeForElement(htmlNode)
	if err != nil {
		panic(err)
	}
	fmt.Println(content)
}
```

### Other exported functions
GetElementNodes returns all html elements `[]*html.Node` found in an html code `htmlNode *html.Node` with the same properties as `e *Element`
```go
func (e *Element) GetElementNodes(htmlNode *html.Node) ([]*html.Node, error)
```
GetTextOfNodes returns the content of an html element `node *html.Node`
```go
func GetTextOfNode(node *html.Node, notRecursive bool) (text string) 
```
RenderNode returns the string representation of a `node *html.Node`
```go
func RenderNode(node *html.Node) string
```
GetHTMLNode returns the node tree `*html.Node` of the html string data
```go
func GetHTMLNode(data string) (*html.Node, error)
```
GetHTML returns the HTML data of URL
```go
func GetHTML(URL string) (string, error)
```

## Contributions

I created this project as a side-project from my normal work. Any contributions are very welcome. Just open up new issues or create a pull request if you want to contribute.
