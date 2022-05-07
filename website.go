package scraper

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// SplitAfter defines the data structure in order to be able to split the finished element
type SplitAfter struct {
	Phrase    string `json:"phrase"`
	Keys      []int  `json:"keys"`
	Seperator string `json:"seperator"`
}

// Split defines the data structure for splitting the html code at a certain phrase
type Split struct {
	Phrase string `json:"phrase"`
	Key    int    `json:"key"`
}

// ReplaceObj defines the data structure for an object, that has to be replaced
type ReplaceObj struct {
	ToReplace string `json:"replace"`
	With      string `json:"with"`
}

// Tag defines the data structure for an HTML Tag
type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Element defines the data structure for an HTML element
type Element struct {
	Typ string `json:"typ"`
	Tag `json:"tag"`
}

// LookUpElement defines the data structure for an element to be looked up by the scraper
type LookUpElement struct {
	SplitAt      []Split      `json:"splitAt"`
	SplitAfter   []SplitAfter `json:"splitAfter"`
	Replacements []ReplaceObj `json:"replacements"`
	Elements     []Element    `json:"elements"`
	NotFound     string       `json:"notFound"`
	HasToContain string       `json:"lastElementHasToContain"`
	Trim         []string     `json:"trim"`
	LastIsURL    bool         `json:"lastIsURL"`
	FollowURL    *Website     `json:"followURL"`
	AddBefore    string       `json:"addBefore"`
	AddAfter     string       `json:"addAfter"`
}

// Website defines the website data type for the scraper
type Website struct {
	Name           string          `json:"name"`
	Seperator      string          `json:"seperator"`
	URL            string          `json:"URL"`
	LookUpElements []LookUpElement `json:"lookUpElements"`
	Cache          string
}

// Scrape scrapes the website w, returning true and the string of the element, if found
func (w *Website) Scrape(funcs map[string]interface{}, cons ...interface{}) (string, error) {
	pW := *w // parsedWebsite (copy of website value)
	if len(funcs) > 0 {
		vls := reflect.ValueOf(&pW).Elem()
		for i := 0; i < vls.NumField(); i++ {
			if vls.Field(i).Kind() == reflect.String {
				vls.Field(i).Set(reflect.ValueOf(formatString(vls.Field(i).String(), funcs, cons)))
			}
		}
	}

	var body string
	if body = GetHTMLdata(pW.URL); len(body) == 0 {
		return "", fmt.Errorf("error while finding body of %v", pW.URL)
	}

	var finishedElements []string
	for _, notEl := range pW.LookUpElements {
		if contains, err := ScrapeElement(body, notEl); err == nil {
			finishedElements = append(finishedElements, contains)
		} else {
			return "", err
		}
	}

	var entireString string
	for _, v := range finishedElements {
		entireString += v + pW.Seperator
	}

	return strings.Trim(entireString, pW.Seperator), nil
}

// ScrapeElement scrapes the html body for a LookUpElement lookEl and returns true and the string of the element when found
func ScrapeElement(body string, lookEl LookUpElement) (string, error) {
	if lookEl.NotFound != "" {
		if strings.Contains(body, lookEl.NotFound) {
			return "", fmt.Errorf("website contains NotFound (%v)", lookEl.NotFound)
		}
	}

	var err error
	for _, v := range lookEl.SplitAt {
		split := strings.Split(body, v.Phrase)
		body, err = checkKey(split, v.Key)
		if err != nil {
			return "", err
		}
	}

	final := body
	if len(lookEl.Elements) > 0 {
		if finalEl, err := GetNestedHTMLElement(body, lookEl.Elements); err != nil && finalEl != "" {
			final = finalEl
		} else {
			fmt.Println("scrapeElement(): Element could not be found by scraper")
		}
	}

	if lookEl.HasToContain != "" {
		if !strings.Contains(final, lookEl.HasToContain) {
			return "", fmt.Errorf("website does not contain HasToContain (%v)", lookEl.HasToContain)
		}
	}

	for _, v := range lookEl.Trim {
		final = strings.Trim(final, v)
	}
	for _, v := range lookEl.SplitAfter {
		split := strings.Split(final, v.Phrase)
		final = ""
		for _, key := range v.Keys {
			finall, err := checkKey(split, key)
			if err != nil {
				return "", err
			}
			final += finall + v.Seperator
		}
		final = strings.Trim(final, v.Seperator)
	}

	for _, replacement := range lookEl.Replacements {
		final = strings.ReplaceAll(final, replacement.ToReplace, replacement.With)
	}

	// add changes here

	if len(lookEl.AddAfter) > 0 {
		final += lookEl.AddAfter
	}

	if len(lookEl.AddBefore) > 0 {
		final = lookEl.AddBefore + final
	}

	final = strings.Trim(final, " ")

	if lookEl.LastIsURL {
		lookEl.FollowURL.URL = final
		return lookEl.FollowURL.Scrape(make(map[string]interface{}, 0))
	}

	if len(final) > 0 && final != lookEl.AddBefore && final != lookEl.AddAfter {
		return final, nil
	}

	return "", errors.New("could not find element in body")
}
