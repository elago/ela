package ela

import (
	"regexp"
	"sort"
)

type uriMode struct {
	mode int // 0:direct mode; 1:arg parse mode
	exp  string
	raw  string
	fun  interface{}
}

var uriMapping map[string]uriMode

func init() {
	uriMapping = make(map[string]uriMode)
}

// put uri-func mapping into map
func Router(uri string, f interface{}) {
	if isArgMode(uri) {
		uriMapping[uri] = uriMode{mode: 1, raw: uri, exp: getArgParseExp(uri), fun: f}
	} else {
		uriMapping[uri] = uriMode{mode: 0, raw: uri, exp: uri, fun: f}
	}

}

// get controller from router map
func getController(uri string) interface{} {

	// using direct match mode
	for k := range uriMapping {
		routerElement := uriMapping[k]
		if routerElement.mode == 0 && routerElement.raw == uri {
			return routerElement.fun
		}
	}

	// using param parse mode

	sorted_keys := make([]string, 0)
	for k, _ := range uriMapping {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	for _, k := range sorted_keys {
		routerElement := uriMapping[k]
		expression := routerElement.exp
		if routerElement.mode == 1 {
			matched, _ := regexp.MatchString(expression, uri)
			if matched {
				return routerElement.fun
			}
		}

	}

	return nil
}

func isArgMode(uri string) bool {
	regexpress := `/:([\D]{1}[\d\D][^:\/\n\r]*)`
	matched, _ := regexp.MatchString(regexpress, uri)
	return matched
}

func getArgParseExp(argUriKey string) string {
	regexpress := `/:([\D]{1}[\d\D][^:\/\n\r]*)`
	paramMatchExp := `/[\d\D][^\/]+`
	reg := regexp.MustCompile(regexpress)
	return `^` + reg.ReplaceAllString(argUriKey, paramMatchExp) + `$`
}
