package scraper

import (
	"errors"
	"reflect"
	"strings"

	"golang.org/x/net/html"
)

// ReplaceObj defines the data structure for an object, that has to be replaced
type ReplaceObj struct {
	ToBeReplaced string `json:"toBeReplaced"`
	Replacement  string `json:"replacement"`
}

// FormatSettings defines the data structure for optional formatting settings of a LookUpElement
type FormatSettings struct {
	Replacements []ReplaceObj `json:"replacements"`
	Trim         []string     `json:"trim"`
	AddBefore    string       `json:"addBefore"`
	AddAfter     string       `json:"addAfter"`
}

// Settings defines the data structure for optional settings of a LookUpElement
type Settings struct {
	FormatSettings           FormatSettings `json:"formatting"`
	DisallowRecursiveContent bool           `json:"disallowRecursiveContent"`
}

// Tag defines the data structure for an HTML Tag
type Tag struct {
	Typ   string `json:"typ"`
	Value string `json:"value"`
}

// Element defines the data structure for an HTML element
type Element struct {
	Typ   string `json:"typ"`
	Tags  []Tag  `json:"tags"`
	Index int    `json:"index"`
}

// LookUpElement defines the data structure for an element to be looked up by the scraper
type LookUpElement struct {
	Element          Element  `json:"element"` // make one element
	Settings         Settings `json:"settings"`
	ContentIsWebsite *Website `json:"followURL"`
}

// Website defines the website data type for the scraper
type Website struct {
	Seperator      string          `json:"seperator"`
	URL            string          `json:"URL"`
	LookUpElements []LookUpElement `json:"lookUpElements"`
}

// Scrape scrapes the website w, returning the found elements in a string each seperated by Seperator
func (w *Website) Scrape(funcs map[string]interface{}, cons ...interface{}) (string, error) {
	copyW := *w // copyWebsite (copy of website value)
	if len(funcs) > 0 {
		vls := reflect.ValueOf(&copyW).Elem()
		for i := 0; i < vls.NumField(); i++ {
			if vls.Field(i).Kind() == reflect.String {
				vls.Field(i).Set(reflect.ValueOf(formatString(vls.Field(i).String(), funcs, cons)))
			}
		}
	}

	node, err := GetHTMLNode(copyW.URL)
	if err != nil {
		return "", err
	}

	var elements []string
	for _, el := range copyW.LookUpElements {
		if content, err := el.ScrapeTreeForElement(node); err != nil {
			return "", err
		} else {
			elements = append(elements, content)
		}
	}

	var elementString string
	for k, v := range elements {
		elementString += v
		if k != len(elements)-1 {
			elementString += copyW.Seperator
		}
	}

	return elementString, nil
}

// ScrapeTreeForElement scraped the node tree for a lookUpElement.Element and formats the content of it accordingly
func (e *LookUpElement) ScrapeTreeForElement(nodeTree *html.Node) (content string, err error) {
	nodes, err := e.Element.GetElementNodes(nodeTree)
	if err != nil {
		return "", err
	}

	// no node found or index out of range
	if len(nodes)-1 < e.Element.Index {
		return "", errors.New("element index out of range or no node found")
	}
	content = GetTextOfNode(nodes[e.Element.Index], e.Settings.DisallowRecursiveContent)

	for _, r := range e.Settings.FormatSettings.Replacements {
		content = strings.ReplaceAll(content, r.ToBeReplaced, r.Replacement)
	}

	for _, v := range e.Settings.FormatSettings.Trim {
		content = strings.Trim(content, v)
	}

	// add changes here

	if len(e.Settings.FormatSettings.AddAfter) > 0 {
		content += e.Settings.FormatSettings.AddAfter
	}

	if len(e.Settings.FormatSettings.AddBefore) > 0 {
		content = e.Settings.FormatSettings.AddBefore + content
	}

	if e.ContentIsWebsite != nil {
		e.ContentIsWebsite.URL = content
		return e.ContentIsWebsite.Scrape(make(map[string]interface{}, 0))
	}

	return content, nil
}
