package ela

import (
	"net/http"
	"os"
	"path/filepath"
)

var (
	staticDirectory = "static"
	specialStatic   []string
)

func init() {
	staticDirectory = config.GetStringDefault("static", "path", "static")
}

func staticServ(uri string, ctx Context) {
	// writer := ctx.w
	// request := ctx.r
	path := uri
	stat, err := os.Stat(filepath.Join(staticDirectory, path))
	if err != nil {
		// 404
		servError(ctx, "<h2>404, File N-ot Exist</h2>", 404, false)
		// http.Error(writer, "404, File Not Exist", 404)
		return
	}

	if !stat.IsDir() {
		// read file
		servPath(path, ctx)
	} else {
		path = path + "/index.html"
		servPath(path, ctx)
	}

}

func servPath(path string, ctx Context) {
	writer := ctx.w
	request := ctx.r
	FileSystem := newStaticFileSystem(staticDirectory)
	f, err := FileSystem.Open(path)

	if err != nil {
		// 404
		servError(ctx, "<h2>404, File N*ot Exist</h2>", 404, false)
		// http.Error(writer, "404, File Not Exist", 404)
		return
	} else {
		fi, err := f.Stat()
		if err != nil {
			// File exists but fail to open.
			// 404
			servError(ctx, "<h2>404, File N/ot Exist</h2>", 404, false)
			// http.Error(writer, "404, File Not Exist", 404)
			return
		}

		http.ServeContent(writer, request, path, fi.ModTime(), f)
	}
}

func staticExist(uri string) bool {
	path := filepath.Join(staticDirectory, uri)
	_, err := os.Stat(path)
	if err != nil {
		return false
	} else {
		return true
	}
}

// add special static files into list
func addSpecialStatic(path string) {
	specialStatic = append(specialStatic, path)
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
