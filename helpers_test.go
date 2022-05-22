package scraper

import (
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatString(t *testing.T) {
	inputAndExpected := make(map[string]string, 0)
	fillMap := func(mapp map[string]string) map[string]string {
		mapp["{{PATH}} should be replaced"] = "path/to/somewhere should be replaced"

		y, m, d := time.Now().Date()
		var mon string
		mo := strconv.Itoa(int(m))
		if len(mo) == 1 {
			for i := 0; i < len(mo); i++ {
				mon += "0"
			}
		}
		mon += strconv.Itoa(int(m))
		mapp["Todays date is {{DATE}}"] = "Todays date is " + strconv.Itoa(d) + "." + mon + "." + strconv.Itoa(y)

		return mapp
	}
	inputAndExpected = fillMap(inputAndExpected)
	var vars []interface{}
	vars = append(vars, time.Now())

	funcs := make(map[string]interface{}, 0)
	funcs["{{PATH}}"] = func(str string) string {
		demoPath := path.Join("path", "to", "somewhere")
		return strings.ReplaceAll(str, "{{PATH}}", demoPath)
	}
	funcs["{{DATE}}"] = func(str string, vars []interface{}) string {
		y, m, d := vars[0].(time.Time).Date()
		var mon string
		mo := strconv.Itoa(int(m))
		if len(mo) == 1 {
			for i := 0; i < len(mo); i++ {
				mon += "0"
			}
		}
		mon += strconv.Itoa(int(m))
		return strings.ReplaceAll(str, "{{DATE}}", strconv.Itoa(d)+"."+mon+"."+strconv.Itoa(y))
	}

	for k, v := range inputAndExpected {
		assert.Equal(t, v, formatString(k, funcs, vars...))
	}
}
