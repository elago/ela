package ela

import (
	"fmt"
	"github.com/gogather/com"
	// "github.com/gogather/com/log"
	"errors"
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

func (cfg *Config) ReloadConfig(path string) {
	cfg.path = path
	cfg.parseIniFile()
}

func (cfg *Config) readConfigFile() (string, error) {
	rawContent, err := com.ReadFileString(cfg.path)
	cfg.rawContent = rawContent
	return cfg.rawContent, err
}

// filter the code comment
func (cfg *Config) filterComment() string {
	reg := regexp.MustCompile(`[#;][\d\D][^\n]*\n*`)
	rep := []byte("\n")
	cfg.content = string(reg.ReplaceAll([]byte(cfg.rawContent), rep))
	return cfg.content
}

// split lines into array
func (cfg *Config) arraylize() []string {
	cfg.rawArrayContainer = strings.Split(cfg.content, "\n")
	return cfg.rawArrayContainer
}

// parse array items as config items
func (cfg *Config) parseItems() {
	count := len(cfg.rawArrayContainer)
	cfg.conf = map[string]map[string]interface{}{}
	cfg.conf["_"] = map[string]interface{}{}
	cfg.warning = nil

	currentSection := "_"
	for i := 0; i < count; i++ {
		item := cfg.rawArrayContainer[i]
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
				cfg.conf[currentSection][key] = cfg.parseValue(value)
			}
		case hasSectionMark && (err2 == nil):
			// section mark line
			reg := regexp.MustCompile(`\[([\d\D][^\[\]]+)]$`)
			result := reg.FindSubmatch([]byte(item))
			if len(result) > 1 {
				currentSection = string(result[1])
				cfg.conf[currentSection] = map[string]interface{}{}
			}
		default:
			cfg.warning = append(cfg.warning, fmt.Sprintf("INI file SyntaxError in Line %d", i+1))
		}

	}
}

// parse value
func (cfg *Config) parseValue(content string) interface{} {
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
func (cfg *Config) parseIniFile() (map[string]map[string]interface{}, error) {
	_, err := cfg.readConfigFile()
	if err != nil {
		return nil, err
	} else {
		cfg.filterComment()
		cfg.arraylize()
		cfg.parseItems()
		return cfg.conf, nil
	}
}

// serialize config as ini file
func (cfg *Config) serialize() string {
	sectionContentMap := map[string]string{}
	content := ""
	for section, val := range cfg.conf {
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

func (cfg *Config) GetWarnings() []string {
	return cfg.warning
}

func (cfg *Config) Get(section, key string) (interface{}, error) {
	sectionMap, ok := cfg.conf[section]
	if !ok {
		return nil, fmt.Errorf("section %s not exist", section)
	}

	value, ok := sectionMap[key]
	if !ok {
		return nil, fmt.Errorf("key %s in section %s not exist", key, section)
	}

	return value, nil
}

func (cfg *Config) GetBool(section, key string) (bool, error) {
	value, err := cfg.Get(section, key)
	if err != nil {
		return false, err
	}

	valueBool, ok := value.(bool)
	if !ok {
		valueString, ok := value.(string)
		if ok {
			if strings.ToLower(valueString) == "true" {
				return true, nil
			} else {
				return false, nil
			}

		} else {
			return false, errors.New("value not bool type")
		}
	} else {
		return valueBool, nil
	}
}

func (cfg *Config) GetInt(section, key string) (int64, error) {
	value, err := cfg.Get(section, key)
	if err != nil {
		return 0, err
	}

	valueInt, ok := value.(int64)
	if ok {
		return valueInt, nil
	} else {
		return 0, errors.New("value not int type")
	}
}

func (cfg *Config) GetFloat(section, key string) (float64, error) {
	value, err := cfg.Get(section, key)
	if err != nil {
		return 0, err
	}

	valueFloat, ok := value.(float64)
	if ok {
		return valueFloat, nil
	} else {
		return 0, errors.New("value not float type")
	}
}

func (cfg *Config) GetString(section, key string) (string, error) {
	value, err := cfg.Get(section, key)
	if err != nil {
		return "", err
	}

	valueString, ok := value.(string)
	if ok {
		return valueString, nil
	} else {
		return "", errors.New("value not string type")
	}
}

func (cfg *Config) GetBoolDefault(section, key string, defaultValue bool) bool {
	value, err := cfg.GetBool(section, key)
	if err != nil {
		return defaultValue
	} else {
		return value
	}
}

func (cfg *Config) GetIntDefault(section, key string, defaultValue int64) int64 {
	value, err := cfg.GetInt(section, key)
	if err != nil {
		return defaultValue
	} else {
		return value
	}
}

func (cfg *Config) GetFloatDefault(section, key string, defaultValue float64) float64 {
	value, err := cfg.GetFloat(section, key)
	if err != nil {
		return defaultValue
	} else {
		return value
	}
}

func (cfg *Config) GetStringDefault(section, key string, defaultValue string) string {
	value, err := cfg.GetString(section, key)
	if err != nil {
		return defaultValue
	} else {
		return value
	}
}

func (cfg *Config) set(section, key string, value interface{}) {
	sectionMap, ok := cfg.conf[section]
	if !ok {
		sectionMap = map[string]interface{}{}
	}
	sectionMap[key] = value
	cfg.conf[section] = sectionMap
}

func (cfg *Config) SetInt(section, key string, value int64) {
	cfg.set(section, key, value)
}

func (cfg *Config) SetBool(section, key string, value bool) {
	cfg.set(section, key, value)
}

func (cfg *Config) SetFloat(section, key string, value float64) {
	cfg.set(section, key, value)
}

func (cfg *Config) SetString(section, key string, value string) {
	cfg.set(section, key, value)
}

func (cfg *Config) Save(path string) error {
	content := cfg.serialize()
	cfg.rawContent = content
	cfg.path = path
	cfg.content = content
	cfg.arraylize()
	return com.WriteFileWithCreatePath(path, cfg.content)
}
