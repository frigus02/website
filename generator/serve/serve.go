package serve

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Serve starts a webserver at the specified directory.
func Serve(dir string) {
	handler := &fileHandler{dir}

	log.Printf("Server listening on http://localhost:8080\n")
	http.ListenAndServe(":8080", handler)
}

type fileHandler struct {
	dir string
}

func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path

	if filepath.Separator != '/' && strings.ContainsRune(upath, filepath.Separator) {
		http.Error(w, "invalid character in file path", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}

	testpath := upath + "/index.html"
	if f.IsFile(testpath) {
		upath = testpath
	}

	if !strings.HasSuffix(upath, ".html") && !strings.HasSuffix(upath, "/") {
		testpath := upath + ".html"
		if f.IsFile(testpath) {
			upath = testpath
		}
	}

	if !f.IsFile(upath) {
		upath = "404.html"
		if !f.IsFile(upath) {
			http.Error(w, "Not found, but 404.html does not exist", http.StatusInternalServerError)
			return
		}
	}

	file, err := f.Open(upath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileinfo, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, fileinfo.Name(), fileinfo.ModTime(), file)
}

func (f *fileHandler) LocalPath(name string) string {
	return filepath.Join(f.dir, filepath.FromSlash(path.Clean("/"+name)))
}

func (f *fileHandler) Open(name string) (http.File, error) {
	fullName := f.LocalPath(name)
	return os.Open(fullName)
}

func (f *fileHandler) IsFile(name string) bool {
	fullName := f.LocalPath(name)
	fileinfo, err := os.Stat(fullName)
	if err != nil {
		return false
	}

	return !fileinfo.IsDir()
}
