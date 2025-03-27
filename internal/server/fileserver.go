package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type FileServer struct {
	rootDir      string
	notFoundPath string
}

func NewFileServer(rootDir string) *FileServer {
	return &FileServer{
		rootDir:      rootDir,
		notFoundPath: filepath.Join(rootDir, "404.html"),
	}
}

func (fs *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestPath := filepath.Join(fs.rootDir, r.URL.Path)

	if fs.pathExists(requestPath) || fs.isValidDirectory(requestPath) {
		http.FileServer(http.Dir(fs.rootDir)).ServeHTTP(w, r)
		return
	}

	fs.serveNotFound(w, r)
}

func (fs *FileServer) pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (fs *FileServer) isValidDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}
	return fs.pathExists(filepath.Join(path, "index.html"))
}

func (fs *FileServer) serveNotFound(w http.ResponseWriter, r *http.Request) {
	if fs.pathExists(fs.notFoundPath) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		http.ServeFile(w, r, fs.notFoundPath)
		return
	}

	fs.servePlainNotFound(w)
}

func (fs *FileServer) servePlainNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "404 Not Found")
}
