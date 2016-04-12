package ela

import (
	"fmt"
	"regexp"
	"sort"
)

type uriMode struct {
	mode int         // 0:direct mode; 1:arg parse mode
	exp  string      // parsed uri pattern
	raw  string      // raw request uri
	fun  interface{} // controller
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
func getController(uri string) (interface{}, map[string]string) {

	// using direct match mode
	for k := range uriMapping {
		routerElement := uriMapping[k]
		if routerElement.mode == 0 && routerElement.raw == uri {
			return routerElement.fun, nil
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
				argMap, _ := getArgs(uri, routerElement.raw)
				return routerElement.fun, argMap
			}
		}
	}

	return nil, nil
}

func isArgMode(uri string) bool {
	regexpress := `/:([\D]{1}[\d\D][^:\/\n\r]*)`
	matched, _ := regexp.MatchString(regexpress, uri)
	return matched
}

func getArgParseExp(argUriKey string) string {
	regexpress := `/:([\D]{1}[\d\D][^:\/\n\r]*)`
	paramMatchExp := `/([\d\D][^\/]+)`
	reg := regexp.MustCompile(regexpress)
	return `^` + reg.ReplaceAllString(argUriKey, paramMatchExp) + `$`
}

// get arguments from uri
// raw: request uri
// pattern: uri pattern
func getArgs(raw, pattern string) (map[string]string, error) {
	if !isArgMode(pattern) {
		return nil, fmt.Errorf("%s", "pattern is not arg parse mode")
	}

	argsMap := make(map[string]string)
	regRaw := regexp.MustCompile(`/:([\D]{1}[\d\D][^:\/\n\r]*)`)
	regExp := regexp.MustCompile(uriMapping[pattern].exp)
	argNameArray := regRaw.FindAllStringSubmatch(pattern, -1)
	argValueArray := regExp.FindAllStringSubmatch(raw, -1)

	var valueSubArray []string
	var lenOfArgValueSubArray int
	lenOfArgValueArray := len(argValueArray)

	if lenOfArgValueArray > 0 {
		valueSubArray = argValueArray[0]
		lenOfArgValueSubArray = len(valueSubArray)
	}

	for i := 0; i < len(argNameArray); i++ {
		key := argNameArray[i][1]
		var value string

		if i+1 < lenOfArgValueSubArray {
			value = valueSubArray[i+1]
		}

		argsMap[key] = value
	}

	return argsMap, nil
}
