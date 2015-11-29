package ela

import (
	"github.com/gogather/com"
	"net/http"
	"os"
	"path/filepath"
)

const (
	staticDirectory = "static"
)

func StaticServ(uri string, writer http.ResponseWriter) {
	path := filepath.Join(staticDirectory, uri)
	stat, err := os.Stat(path)
	if err != nil {
		// 404
		http.Error(writer, "404, File Not Exist", 404)
		return
	}

	if !stat.IsDir() {
		// read file
		data, err := com.ReadFileByte(path)
		if err == nil {
			writer.Write([]byte(data))
		} else {
			// 404
			http.Error(writer, "404, File Not Exist", 404)
			return
		}
	} else {
		path := filepath.Join(path, "index.html")
		_, err := os.Stat(path)
		if err != nil {
			// 404
			http.Error(writer, "404, File Not Exist", 404)
			return
		}

		data, err := com.ReadFileByte(path)
		if err == nil {
			writer.Write([]byte(data))
		} else {
			// 404
			http.Error(writer, "404, File Not Exist", 404)
		}
	}

}

func StaticExist(uri string) bool {
	path := filepath.Join(staticDirectory, uri)
	_, err := os.Stat(path)
	if err != nil {
		return false
	} else {
		return true
	}
}
