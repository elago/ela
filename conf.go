package ela

import (
	// "fmt"
	"github.com/gogather/com"
	"regexp"
	"strings"
)

type Config struct {
	Conf              map[string]map[string]interface{}
	path              string
	rawContent        string
	content           string
	rawArrayContainer []string
}

func NewConfig(path string) Config {
	conf := Config{path: path}
	return conf
}

func (this *Config) readConfigFile() (string, error) {
	rawContent, err := com.ReadFileString(this.path)
	this.rawContent = rawContent
	return this.rawContent, err
}

// filter the code comment
func (this *Config) filterComment() string {
	reg := regexp.MustCompile(`#[\d\D][^\n#]*\n`)
	rep := []byte("\n")
	this.content = string(reg.ReplaceAll([]byte(this.rawContent), rep))
	return this.content
}

// split lines into array
func (this *Config) arraylize() []string {
	this.rawArrayContainer = strings.Split(this.content, "\n")
	return this.rawArrayContainer
}

// parse array items as config items
func (this *Config) parseItems() {

}

func (this *Config) parseIniFile() map[string]map[string]interface{} {
	return nil
}
