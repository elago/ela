package ela

import (
	"fmt"
	"github.com/gogather/com"
	//"github.com/gogather/com/log"
	"html/template"
	"io/ioutil"
	"sort"
	"strings"
)

type Templates struct {
	templates          map[string]string
	templatesName      []string
	templatefolder     string
	templatesFilterOut []string
	funcMap            template.FuncMap
}

func NewTemplates(templatefolder string) {
	t := Templates{}
	t.templates = make(map[string]string)
	t.templatefolder = templatefolder
	// templatefolder := config.GetStringDefault("_", "template", "web/view")
	t.listFile(templatefolder)

	// initialize template functions
	t.funcMap = template.FuncMap{}
	t.AddTemplateFunc("tplFunc", t.tplFunc)
}

func (t *Templates) reloadTemplate() {
	t.listFile(t.templatefolder)
}

func (t *Templates) addFilterOut(filename string) {
	t.templatesFilterOut = append(t.templatesFilterOut, filename)
}

func (t *Templates) notInFilterOut(filename string) bool {
	for i := 0; i < len(t.templatesFilterOut); i++ {
		if t.templatesFilterOut[i] == filename {
			return false
		}
	}
	return true
}

func (t *Templates) listFile(dir string) {
	t.addFilterOut(".DS_Store")

	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		subfile := dir + "/" + file.Name()
		if file.IsDir() {
			t.listFile(subfile)
		} else if t.notInFilterOut(file.Name()) {
			content, err := com.ReadFileString(subfile)
			if err != nil {
				t.templates[subfile] = ""
			} else {
				content = "{{define \"" + subfile + "\"}}" + content + "{{end}}"
				content = t.lefTplDir(content, t.templatefolder)
				subfile = t.lefTplDir(subfile, t.templatefolder)
				t.templates[subfile] = content
			}
		}
	}

	// get template name list
	t.getTemplateNames()
}

func (t *Templates) getTemplateNames() {
	t.templatesName = make([]string, 0, len(t.templates))
	for k := range t.templates {
		t.templatesName = append(t.templatesName, k)
	}
}

func (t *Templates) lefTplDir(dir string, tplDir string) string {
	return strings.Replace(dir, tplDir+"/", "", 1)
}

func (t *Templates) parseFiles(filenames ...string) (*template.Template, error) {
	var temp *template.Template = nil
	var err error = nil

	if len(filenames) == 0 {
		return nil, fmt.Errorf("html/template: no files named in call to ParseFiles")
	}

	sort.Strings(filenames)

	for _, filename := range filenames {
		if temp == nil {
			temp = template.New(filename)
		}
		if filename != temp.Name() {
			temp = temp.New(filename)
		}
		_, err = temp.Funcs(t.funcMap).Parse(t.templates[filename])

		// anyone template syntax error throw panic
		if err != nil {
			panic(err)
		}
	}
	return temp, err
}

func (t *Templates) SetTemplateDir(dir string) {
	exist := com.FileExist(dir)
	if exist {
		t.templatefolder = dir
	}
}

// add template function definition
func (t *Templates) AddTemplateFunc(functionName string, function interface{}) {
	t.funcMap[functionName] = function
}

// test template function define
func (t *Templates) tplFunc(test string) string {
	return "template function, " + test
}
