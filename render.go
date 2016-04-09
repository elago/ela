package ela

import (
	"github.com/gogather/com"
	// "github.com/gogather/com/log"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var templates map[string]string
var templatesName []string
var templatefolder string
var templatesFilterOut []string

func init() {
	templates = make(map[string]string)
	templatefolder = `view`
	listFile(templatefolder)

	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)

	println(path)
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
