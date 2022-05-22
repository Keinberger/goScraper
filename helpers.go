package scraper

import (
	"strings"
)

type ErrType int

const (
	// ErrMissingElement will be returned if the element is issing
	ErrMissingElement = iota
	// ErrNoNodeFound will be returned if no element was found
	ErrNoNodeFound
	// ErrIdxOutOfRange will be returned if the index of an array is out of range
	ErrIdxOutOfRange
)

// Error defines the data structure for a custom error
type Error struct {
	ErrType
	msg string
}

// Error returns the error msg of an error
func (e Error) Error() string {
	return e.msg
}

// newErr creates a new err of type typ, including a msg
func newErr(typ ErrType, msg string) Error {
	return Error{ErrType: typ, msg: msg}
}

// formatString replaces str with the return value of a func specified in funcs,
// if str contains the map key of funcs one may use an array of variables
// that can be passed into the function, as well, hence a function has to have
// a string as the first argument (where str gets passed in) and an array
// of interface{} as an optional second argument, for variables to be used inside of the func
// e.g.: One may use {{path}} inside of the url, which will be replaced by a certain string
// coded into the function of the map at key {{path}}
// e.g.: One may use {{date}} inside of the url, which will be replaced by a time.Now() constant
// specified in the constants array, coded into the function of the map at key {{date}}
func formatString(str string, funcs map[string]interface{}, vars ...interface{}) string {
	for k, v := range funcs {
		if strings.Contains(str, k) {
			if fn, ok := v.(func(string, []interface{}) string); ok {
				str = fn(str, vars)
			} else if fn, ok := v.(func(string) string); ok {
				str = fn(str)
			}
		}
	}
	return str
}
