package scraper

import (
	"fmt"
	"strings"
)

// logError() logs an error into the console
func logError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err) // add loggin error to log file here
	}
}

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
			if _, ok := v.(func(string, []interface{}) string); ok {
				str = v.(func(string, []interface{}) string)(str, constants)
			} else if _, ok := v.(func(string) string); ok {
				str = v.(func(string) string)(str)
			}
		}
	}
	return str
}

// checkKey checks if the key exists inside of the string array strArr
// one may use -1 as the key to return the last element of the array
func checkKey(strArr []string, key int) string {
	if key >= len(strArr) {
		return ""
	} else if key == -1 {
		return strArr[len(strArr)-1]
	}
	return strArr[key]
}
