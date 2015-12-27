package ela

import (
	"github.com/gogather/com"
	"github.com/gogather/com/log"
	"io/ioutil"
)

var templates map[string]string
var templatesName []string

func init() {
	templates = make(map[string]string)
	templatefolder := `view`
	listFile(templatefolder)

	log.Blueln(templatesName)
}

func listFile(dir string) {
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
				templates[subfile] = content
				templatesName = append(templatesName, subfile)
			}
		}
	}

}

func Render() {

}
