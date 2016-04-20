package ela

import (
	"github.com/gogather/com"
	// "github.com/gogather/com/log"
	"fmt"
	"html/template"
	"io/ioutil"
	"sort"
	"strings"
)

var (
	templates          map[string]string
	templatesName      []string
	templatefolder     string
	templatesFilterOut []string
	funcMap            template.FuncMap
)

func init() {
	templates = make(map[string]string)
	templatefolder = `view`
	listFile(templatefolder)

	// initialize template functions
	funcMap = template.FuncMap{}
	AddTemplateFunc("tplFunc", tplFunc)
}

func reloadTemplate() {
	listFile(templatefolder)
}

func addFilterOut(filename string) {
	templatesFilterOut = append(templatesFilterOut, filename)
}

func notInFilterOut(filename string) bool {
	for i := 0; i < len(templatesFilterOut); i++ {
		if templatesFilterOut[i] == filename {
			return false
		}
	}
	return true
}

func listFile(dir string) {
	addFilterOut(".DS_Store")

	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		subfile := dir + "/" + file.Name()
		if file.IsDir() {
			listFile(subfile)
		} else if notInFilterOut(file.Name()) {
			content, err := com.ReadFileString(subfile)
			if err != nil {
				templates[subfile] = ""
			} else {
				content = "{{define \"" + subfile + "\"}}" + content + "{{end}}"
				content = lefTplDir(content, templatefolder)
				subfile = lefTplDir(subfile, templatefolder)
				templates[subfile] = content
			}
		}
	}

	// get template name list
	getTemplateNames()
}

func getTemplateNames() {
	templatesName = make([]string, 0, len(templates))
	for k := range templates {
		templatesName = append(templatesName, k)
	}
}

func lefTplDir(dir string, tplDir string) string {
	return strings.Replace(dir, tplDir+"/", "", 1)
}

func parseFiles(filenames ...string) (*template.Template, error) {
	var t *template.Template = nil
	var err error = nil

	if len(filenames) == 0 {
		return nil, fmt.Errorf("html/template: no files named in call to ParseFiles")
	}

	sort.Strings(filenames)

	for _, filename := range filenames {
		if t == nil {
			t = template.New(filename)
		}
		if filename != t.Name() {
			t = t.New(filename)
		}
		_, err = t.Funcs(funcMap).Parse(templates[filename])

		// anyone template syntax error throw panic
		if err != nil {
			panic(err)
		}
	}
	return t, err
}

// add template function definition
func AddTemplateFunc(functionName string, function interface{}) {
	funcMap[functionName] = function
}

// test template function define
func tplFunc(test string) string {
	return "template function, " + test
}
