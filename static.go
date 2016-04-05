package ela

import (
	"net/http"
	"os"
	"path/filepath"
)

var (
	staticDirectory = "static"
)

func init() {
	staticDirectory = config.GetStringDefault("static", "path", "static")
}

func StaticServ(uri string, writer http.ResponseWriter, request *http.Request) {
	path := uri
	stat, err := os.Stat(filepath.Join(staticDirectory, path))
	if err != nil {
		// 404
		http.Error(writer, "404, File Not Exist", 404)
		return
	}

	if !stat.IsDir() {
		// read file
		servPath(path, writer, request)
	} else {
		path = path + "/index.html"
		servPath(path, writer, request)
	}

}

func servPath(path string, writer http.ResponseWriter, request *http.Request) {
	FileSystem := newStaticFileSystem(staticDirectory)
	f, err := FileSystem.Open(path)

	if err != nil {
		// 404
		http.Error(writer, "404, File Not Exist", 404)
		return
	} else {
		fi, err := f.Stat()
		if err != nil {
			// File exists but fail to open.
			// 404
			http.Error(writer, "404, File Not Exist", 404)
			return
		}

		http.ServeContent(writer, request, path, fi.ModTime(), f)
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

// staticFileSystem implements http.FileSystem interface.
type staticFileSystem struct {
	dir *http.Dir
}

func newStaticFileSystem(directory string) staticFileSystem {
	Root, err := os.Getwd()
	if err != nil {
		panic("error getting work directory: " + err.Error())
	}

	if !filepath.IsAbs(directory) {
		directory = filepath.Join(Root, directory)
	}
	dir := http.Dir(directory)
	return staticFileSystem{&dir}
}

func (fs staticFileSystem) Open(name string) (http.File, error) {
	return fs.dir.Open(name)
}
