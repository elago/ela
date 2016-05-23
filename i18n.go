package ela

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type I18n struct {
	basePath string
	lang     map[string]string
	data     map[string]*Config
}

func NewEmptyI18n() *I18n {
	return &I18n{}
}

func NewI18n(basePath string) *I18n {
	i18n := NewEmptyI18n()
	i18n.Load(basePath)
	return i18n
}

func (i *I18n) Load(basePath string) error {
	stat, err := os.Stat(filepath.Join(basePath))
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("base path should be a directory\n")
	}

	i.basePath = basePath

	err = i.scanLang()
	if err != nil {
		return fmt.Errorf("scan path error: %s\n", err)
	}
	i.parseIni()
	return err
}

func (i *I18n) scanLang() error {
	i.lang = make(map[string]string)
	err := filepath.Walk(i.basePath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		filename := f.Name()
		if !strings.HasSuffix(filename, ".ini") {
			return nil
		}
		lang := strings.Replace(filename, ".ini", "", 1)
		i.lang[lang] = filepath.Join(i.basePath, filename)
		return nil
	})
	return err
}

func (i *I18n) parseIni() {
	i.data = make(map[string]*Config)
	for lang, path := range i.lang {
		i.data[lang] = NewConfig(path)
	}
}

func (i *I18n) Tr(lang, section, key string) string {
	langData, ok := i.data[lang]
	if ok {
		return langData.GetStringDefault(section, key, key)
	} else {
		return key
	}
}
