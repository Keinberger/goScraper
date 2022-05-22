package scraper

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScrape(t *testing.T) {
	testMap := make(map[string]func(t *testing.T), 0)

	testMap["scrapeWebsite_OneElement"] = func(t *testing.T) {
		testWebsite := Website{
			URL: "https://www.wikipedia.org/wiki/Wikipedia",
			LookUpElements: []LookUpElement{
				{
					Element: Element{
						Typ: "h1",
						Tags: []Tag{
							{
								Typ:   "id",
								Value: "firstHeading",
							},
						},
					},
				},
			},
		}

		expected := "Wikipedia"
		content, err := testWebsite.Scrape(nil)
		require.NoError(t, err)
		assert.Equal(t, expected, content)
	}
	testMap["scrapeWebsite_TwoElementsAndSeperator"] = func(t *testing.T) {
		testWebsite := Website{
			Seperator: ", ",
			URL:       "https://www.wikipedia.org/wiki/Wikipedia",
			LookUpElements: []LookUpElement{
				{
					Element: Element{
						Typ: "h1",
						Tags: []Tag{
							{
								Typ:   "id",
								Value: "firstHeading",
							},
						},
					},
				},
				{
					Element: Element{
						Typ: "li",
						Tags: []Tag{
							{
								Typ:   "id",
								Value: "ca-history",
							},
						},
					},
				},
			},
		}

		expected := "Wikipedia, View history"
		content, err := testWebsite.Scrape(nil)
		require.NoError(t, err)
		assert.Equal(t, expected, content)
	}
	testMap["scrapeWebsite_ReplacementFuncs"] = func(t *testing.T) {
		strToBeReplaced := "{{INSERT_WIKIPEDIA}}"
		testWebsite := Website{
			URL: "https://en.wikipedia.org/wiki/" + strToBeReplaced,
			LookUpElements: []LookUpElement{
				{
					Element: Element{
						Typ: "h1",
						Tags: []Tag{
							{
								Typ:   "id",
								Value: "firstHeading",
							},
						},
					},
				},
			},
		}
		funcs := make(map[string]interface{}, 0)
		funcs[strToBeReplaced] = func(str string) string {
			return strings.ReplaceAll(str, strToBeReplaced, "Wikipedia")
		}

		expected := "Wikipedia"
		content, err := testWebsite.Scrape(&funcs)
		require.NoError(t, err)
		assert.Equal(t, expected, content)

	}
	testMap["scrapeWebsite_ReplacementFuncsWithVars"] = func(t *testing.T) {
		numericalDateStr := "{{NUMERICAL_DATE}}" // 2022/05/22
		testWebsite := Website{
			URL: "https://www.nytimes.com/issue/todayspaper/{{NUMERICAL_DATE}}/todays-new-york-times",
			LookUpElements: []LookUpElement{
				{
					Element: Element{
						Typ: "h2",
						Tags: []Tag{
							{
								Typ:   "class",
								Value: "css-q1brm6",
							},
						},
					},
				},
			},
		}

		funcs := make(map[string]interface{}, 0)
		funcs[numericalDateStr] = func(str string, vars []interface{}) string {
			y, m, d := vars[0].(time.Time).Date()
			var mon string
			mo := strconv.Itoa(int(m))
			if len(mo) == 1 {
				for i := 0; i < len(mo); i++ {
					mon += "0"
				}
			}
			mon += strconv.Itoa(int(m))
			return strings.ReplaceAll(str, numericalDateStr, strconv.Itoa(y)+"/"+mon+"/"+strconv.Itoa(d))
		}

		expected := " The Front Page"
		content, err := testWebsite.Scrape(&funcs, time.Now())
		require.NoError(t, err)
		assert.Equal(t, expected, content)
	}

	for testName, testFunc := range testMap {
		t.Run(testName, testFunc)
	}
}

func TestScrapeTreeForElement(t *testing.T) {
	testMap := make(map[string]func(t *testing.T), 0)
	nodeTree, err := GetHTMLNode(testHTML)
	require.NoError(t, err)

	testLookUpElement := LookUpElement{
		Element: Element{
			Typ: "",
			Tags: []Tag{
				{
					Typ:   "",
					Value: "",
				},
			},
		},
		Settings: Settings{},
	}
	_ = testLookUpElement

	testMap["scrapeTreeForContent"] = func(t *testing.T) {
		testLookUpElement := LookUpElement{
			Element: Element{
				Typ: "p",
				Tags: []Tag{
					{
						Typ:   "id",
						Value: "singleElement_OneTag",
					},
				},
			},
		}
		expected := "This is the single element with one tag"
		actual, err := testLookUpElement.ScrapeTreeForElement(nodeTree)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
	testMap["elementIndexOutOfRange"] = func(t *testing.T) {
		testLookUpElement := LookUpElement{
			Element: Element{
				Typ: "div",
				Tags: []Tag{
					{
						Typ:   "class",
						Value: "hasDuplicate",
					},
				},
			},
			Index: 2,
		}
		_, err := testLookUpElement.ScrapeTreeForElement(nodeTree)
		require.Error(t, err)
		assert.Equal(t, "element index out of range", err.Error())
	}
	testMap["settingsReplacements"] = func(t *testing.T) {
		testLookUpElement := LookUpElement{
			Element: Element{
				Typ: "p",
				Tags: []Tag{
					{
						Typ:   "id",
						Value: "singleElement_OneTag",
					},
				},
			},
			Settings: Settings{
				FormatSettings: FormatSettings{
					Replacements: []ReplaceObj{
						{
							ToBeReplaced: " ",
							Replacement:  "_",
						},
					},
				},
			},
		}
		expected := "This_is_the_single_element_with_one_tag"
		actual, err := testLookUpElement.ScrapeTreeForElement(nodeTree)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
	testMap["settingsTrim"] = func(t *testing.T) {
		testLookUpElement := LookUpElement{
			Element: Element{
				Typ: "p",
				Tags: []Tag{
					{
						Typ:   "id",
						Value: "elementToBeTrimmed",
					},
				},
			},
			Settings: Settings{
				FormatSettings: FormatSettings{
					Trim: []string{
						" ",
					},
				},
			},
		}
		expected := "This is the elemnt which needs some trimming"
		actual, err := testLookUpElement.ScrapeTreeForElement(nodeTree)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
	testMap["settingsAddAfter"] = func(t *testing.T) {
		testLookUpElement := LookUpElement{
			Element: Element{
				Typ: "p",
				Tags: []Tag{
					{
						Typ:   "id",
						Value: "singleElement_OneTag",
					},
				},
			},
			Settings: Settings{
				FormatSettings: FormatSettings{
					AddAfter: ", literally only one tag",
				},
			},
		}
		expected := "This is the single element with one tag, literally only one tag"
		actual, err := testLookUpElement.ScrapeTreeForElement(nodeTree)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
	testMap["settingsAddBefore"] = func(t *testing.T) {
		testLookUpElement := LookUpElement{
			Element: Element{
				Typ: "p",
				Tags: []Tag{
					{
						Typ:   "id",
						Value: "singleElement_OneTag",
					},
				},
			},
			Settings: Settings{
				FormatSettings: FormatSettings{
					AddBefore: "The element with only one tag: ",
				},
			},
		}
		expected := "The element with only one tag: This is the single element with one tag"
		actual, err := testLookUpElement.ScrapeTreeForElement(nodeTree)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	}
	testMap["ContentIsFollowURL"] = func(t *testing.T) {
		testLookUpElement := LookUpElement{
			Element: Element{
				Typ: "a",
				Tags: []Tag{
					{
						Typ:   "id",
						Value: "websiteLink",
					},
				},
			},
			ContentIsFollowURL: &Website{
				LookUpElements: []LookUpElement{
					{
						Element: Element{
							Typ: "h1",
							Tags: []Tag{
								{
									Typ:   "id",
									Value: "firstHeading",
								},
							},
						},
					},
				},
			},
		}
		expected := "Wikipedia"
		actual, err := testLookUpElement.ScrapeTreeForElement(nodeTree)
		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	}

	for testName, testFunc := range testMap {
		t.Run(testName, testFunc)
		fmt.Println(testName)
	}
}
