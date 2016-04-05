package ela

var uriMapping map[string]interface{}

func init() {
	uriMapping = make(map[string]interface{})
}

// put uri-func mapping into map
func Router(uri string, f func(Context)) {
	uriMapping[uri] = f
}
