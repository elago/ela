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

func listFile(dir string) {
	templatesName = nil
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		subfile := dir + "/" + file.Name()
		if file.IsDir() {
			listFile(subfile)
		} else {
			content, err := com.ReadFileString(subfile)
			if err != nil {
				templates[subfile] = ""
			} else {
				content = "{{define \"" + subfile + "\"}}\n" + content + "\n{{end}}"
				content = lefTplDir(content, templatefolder)
				subfile = lefTplDir(subfile, templatefolder)
				templates[subfile] = content
				templatesName = append(templatesName, subfile)
			}
		}
	}
}

func lefTplDir(dir string, tplDir string) string {
	return strings.Replace(dir, tplDir+"/", "", 1)
}
