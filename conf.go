package ela

import (
	"fmt"
	"github.com/gogather/com"
	// "github.com/gogather/com/log"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	conf              map[string]map[string]interface{}
	path              string
	rawContent        string
	content           string
	rawArrayContainer []string
	warning           []string
}

func NewConfig(path string) Config {
	conf := Config{path: path}
	conf.parseIniFile()
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
	count := len(this.rawArrayContainer)
	this.conf = map[string]map[string]interface{}{}
	this.conf["_"] = map[string]interface{}{}
	this.warning = nil

	currentSection := "_"
	for i := 0; i < count; i++ {
		item := this.rawArrayContainer[i]
		item = strings.TrimSpace(item)
		hasEqualMark, err1 := regexp.Match("=", []byte(item))
		hasSectionMark, err2 := regexp.Match(`\[[\d\D][^\[\]]+]$`, []byte(item))

		switch {
		case len(item) <= 0:
			//empty line, skip
		case hasEqualMark && (err1 == nil):
			//normal key value item
			kvArray := strings.Split(item, "=")
			key := strings.TrimSpace(kvArray[0])
			value := strings.TrimSpace(kvArray[1])
			this.conf[currentSection][key] = this.parseValue(value)
		case hasSectionMark && (err2 == nil):
			// section mark line
			reg := regexp.MustCompile(`\[([\d\D][^\[\]]+)]$`)
			result := reg.FindSubmatch([]byte(item))
			if len(result) > 1 {
				currentSection = string(result[1])
				this.conf[currentSection] = map[string]interface{}{}
			}
		default:
			this.warning = append(this.warning, fmt.Sprintf("INI file SyntaxError in Line %d", i+1))
		}

	}
}

// parse value
func (this *Config) parseValue(content string) interface{} {
	reg := regexp.MustCompile(`\"([\d\D][^\"]+)"$`)
	result := reg.FindSubmatch([]byte(content))

	if len(result) > 1 {
		return string(result[1])
	}

	boolValue, err := strconv.ParseBool(content)
	if err == nil {
		return boolValue
	}

	intValue, err := strconv.ParseInt(content, 0, 64)
	if err == nil {
		return intValue
	}

	floatValue, err := strconv.ParseFloat(content, 64)
	if err == nil {
		return floatValue
	}

	return content
}

// parse ini file
func (this *Config) parseIniFile() (map[string]map[string]interface{}, error) {
	_, err := this.readConfigFile()
	if err != nil {
		return nil, err
	} else {
		this.filterComment()
		this.arraylize()
		this.parseItems()
		return this.conf, nil
	}
}

func (this *Config) GetWarnings() []string {
	return this.warning
}

func (this *Config) Get(section, key string) interface{} {
	return this.conf[section][key]
}

func (this *Config) GetBool(section, key string) bool {
	return this.Get(section, key).(bool)
}

func (this *Config) GetInt(section, key string) int64 {
	return this.Get(section, key).(int64)
}

func (this *Config) GetFloat(section, key string) float64 {
	return this.Get(section, key).(float64)
}

func (this *Config) GetString(section, key string) string {
	return this.Get(section, key).(string)
}
