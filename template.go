package ela

import (
	"fmt"
	"html/template"
)

func init() {
	// initialize template functions
	AddTemplateFunc("tplFunc", tplFunc)
}

var funcMap = template.FuncMap{}

func parseFiles(filenames ...string) (*template.Template, error) {
	var t *template.Template = nil
	var err error = nil

	if len(filenames) == 0 {
		return nil, fmt.Errorf("html/template: no files named in call to ParseFiles")
	}

	for _, filename := range filenames {
		if t == nil {
			t = template.New(filename)
		}
		if filename != t.Name() {
			t = t.New(filename)
		}
		_, err = t.Funcs(funcMap).Parse(templates[filename])
	}
	return t, err
}

func AddTemplateFunc(functionName string, function interface{}) {
	funcMap[functionName] = function
}

// test template function define
func tplFunc(test string) string {
	return "template function, " + test
}
