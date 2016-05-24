package ela

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type I18n struct {
	basePath    string
	lang        map[string]string
	data        map[string]*Config
	currentLang string
}

// new an empty i18n module
func NewEmptyI18n() *I18n {
	return &I18n{currentLang: ""}
}

// new an i18n module
func NewI18n(basePath string) *I18n {
	i18n := NewEmptyI18n()
	i18n.Load(basePath)
	return i18n
}

func InitModule(ctx *Context, basePath string) *I18n {
	i18n := NewI18n(basePath)
	i18n.initModule(ctx)
	return i18n
}

// load i18n files for i18n module
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

func (i *I18n) Tr(args ...string) string {
	section := "_"
	key := args[0]

	if len(args) > 1 {
		section = args[0]
		key = args[1]
	}

	lang := i.Lang()
	langData, ok := i.data[lang]
	if ok {
		return langData.GetStringDefault(section, key, key)
	} else {
		return key
	}
}

func (i *I18n) InitI18nModule(ctx *Context) {
	lang := ctx.GetParam("lang")
	if strings.TrimSpace(lang) == "" {
		lang = ctx.GetCookie("lang")
		if strings.TrimSpace(lang) == "" {
			acceptLang := ctx.GetRequest().Header.Get("Accept-Language")
			reg1 := regexp.MustCompile(`;`)
			reg2 := regexp.MustCompile(`-`)
			acceptLang = reg1.ReplaceAllString(acceptLang, ",")
			acceptLang = reg2.ReplaceAllString(acceptLang, "_")
			langArray := strings.Split(acceptLang, ",")
			for idx := 0; idx < len(langArray); idx++ {
				value := langArray[idx]
				_, ok := i.lang[value]
				if ok {
					lang = value
					break
				}
			}
		}
	}

	i.currentLang = strings.TrimSpace(lang)

	if i.currentLang == "" {
		i.currentLang = "en_US"
	}

	lang = i.Lang()

	ctx.Data["Lang"] = lang
	ctx.Data["Tr"] = i.Tr
}

func (i *I18n) SetLang(lang string) {
	i.currentLang = lang
}

func (i *I18n) Lang() string {
	if i.currentLang == "" {
		return "en_US"
	} else {
		return i.currentLang
	}
}
