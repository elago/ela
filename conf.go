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
		hasEqualMark, err1 := regexp.Match(`=`, []byte(item))
		hasSectionMark, err2 := regexp.Match(`\[[\d\D][^\[\]]+]$`, []byte(item))

		switch {
		case len(item) <= 0:
			//empty line, skip
		case hasEqualMark && (err1 == nil):
			//normal key value item
			reg := regexp.MustCompile(`([\d\D][^=]+)=([\d\D]+)$`)
			kvArray := reg.FindSubmatch([]byte(item))
			if len(kvArray) > 2 {
				key := strings.TrimSpace(string(kvArray[1]))
				value := strings.TrimSpace(string(kvArray[2]))
				this.conf[currentSection][key] = this.parseValue(value)
			}
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

// serialize config as ini file
func (this *Config) serialize() string {
	sectionContentMap := map[string]string{}
	content := ""
	for section, val := range this.conf {
		if section == "_" {
			sectionMap := val
			sectionContent := ""
			for key, value := range sectionMap {
				sectionContent = sectionContent + fmt.Sprintf("%s = %v\n", key, value)
			}
			sectionContentMap["_"] = sectionContent
		} else {
			// section title
			title := "\n[" + section + "]\n"
			sectionMap := val
			sectionContent := ""
			for key, value := range sectionMap {
				sectionContent = sectionContent + fmt.Sprintf("%s = %v\n", key, value)
			}
			sectionContentMap[section] = title + sectionContent
		}
	}

	for section, sectionVal := range sectionContentMap {
		if section == "_" {
			content = sectionVal + content
		} else {
			content = content + sectionVal
		}
	}

	return content
}

func (this *Config) GetWarnings() []string {
	return this.warning
}

func (this *Config) Get(section, key string) interface{} {
	return this.conf[section][key]
}

func (this *Config) GetBool(section, key string) bool {
	value, ok := this.Get(section, key).(bool)
	if !ok {
		valueString, ok := this.Get(section, key).(string)
		if ok && strings.ToLower(valueString) == "true" {
			return true
		} else {
			return false
		}
	} else {
		return value
	}
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

func (this *Config) set(section, key string, value interface{}) {
	sectionMap, ok := this.conf[section]
	if !ok {
		sectionMap = map[string]interface{}{}
	}
	sectionMap[key] = value
	this.conf[section] = sectionMap
}

func (this *Config) SetInt(section, key string, value int64) {
	this.set(section, key, value)
}

func (this *Config) SetBool(section, key string, value bool) {
	this.set(section, key, value)
}

func (this *Config) SetFloat(section, key string, value float64) {
	this.set(section, key, value)
}

func (this *Config) SetString(section, key string, value string) {
	this.set(section, key, value)
}

func (this *Config) Save(path string) error {
	content := this.serialize()
	this.rawContent = content
	this.path = path
	this.content = content
	this.arraylize()
	return com.WriteFileWithCreatePath(path, this.content)
}
