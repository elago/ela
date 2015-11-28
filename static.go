package ela

import (
	"github.com/gogather/com"
	"net/http"
	"os"
)

const (
	staticDirectory = "static"
)

func StaticServ(uri string, writer http.ResponseWriter) {
	// read file
	data, err := com.ReadFileByte(staticDirectory + uri)
	if err == nil {
		writer.Write([]byte(data))
	} else {
		// 404
		http.Error(writer, "404, File Not Exist", 404)
	}
}

func StaticExist(uri string) bool {
	_, err := os.Stat(staticDirectory + uri)
	if err != nil {
		return false
	} else {
		return true
	}
}
