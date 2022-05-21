package scraper

import (
	"strings"
)

// formatString replaces str with a func specified in funcs, if str contains the map key of funcs
// one may use an array of constant values that can be passed into the function, as well, hence
// a function has to have a string as the first argument (where str gets passed in) and an array
// of interface{} as an optional second argument.
// e.g.: One may use {{path}} inside of the url, which will be replaced by a certain string
// coded into the function of the map at key {{path}}
// e.g.: One may use {{date}} inside of the url, which will be replaced by a time.Now() constant
// specified in the constants array, coded into the function of the map at key {{date}}
func formatString(str string, funcs map[string]interface{}, constants []interface{}) string {
	for k, v := range funcs {
		if strings.Contains(str, k) {
			if fn, ok := v.(func(string, []interface{}) string); ok {
				str = fn(str, constants)
			} else if fn, ok := v.(func(string) string); ok {
				str = fn(str)
			}
		}
	}
	return str
}
