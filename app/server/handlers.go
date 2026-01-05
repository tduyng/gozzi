// Package server provides HTTP handlers for file serving and live reload.
package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type fileHandler struct {
	root     http.Dir
	notFound string
	dev      bool
}

func (h *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean(r.URL.Path)
	if path == "." {
		path = "/"
	}

	f, err := h.root.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			h.serve404(w, r)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = f.Close()
	}()
	stat, _ := f.Stat()

	if stat.IsDir() {
		indexPath := filepath.Join(path, "index.html")
		indexFile, err := h.root.Open(indexPath)
		if err != nil {
			h.serve404(w, r)
			return
		}
		defer func() {
			_ = indexFile.Close()
		}()
		h.serveHTML(indexFile, w)
		return

	}

	if h.dev && strings.HasSuffix(path, ".html") {
		h.serveHTML(f, w)
		return
	}

	http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
}

func (h *fileHandler) serveHTML(f http.File, w http.ResponseWriter) {
	content, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	script := []byte(`<script>
	(new EventSource('/livereload')).onmessage = () => location.reload()
	</script></body>`)
	content = bytes.Replace(content, []byte("</body>"), script, 1)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Expires", "0")
	_, _ = w.Write(content)
}

func (h *fileHandler) serve404(w http.ResponseWriter, r *http.Request) {
	f, err := h.root.Open(h.notFound)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer func() {
		_ = f.Close()
	}()
	h.serveHTML(f, w)
}

func (s *DevServer) handleLiveReload(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	msgChan := make(chan string)
	s.mu.Lock()
	s.clients[msgChan] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, msgChan)
		s.mu.Unlock()
		close(msgChan)
	}()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg := <-msgChan:
			_, _ = fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}
