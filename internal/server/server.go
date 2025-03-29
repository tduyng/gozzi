package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/tduyng/gozzi/internal/config"
	"github.com/tduyng/gozzi/internal/generator"
)

type DevServer struct {
	cfg          *config.GlobalConfig
	watcher      *fsnotify.Watcher
	clients      map[chan []byte]struct{}
	clientMutex  sync.Mutex
	debounce     *time.Timer
	rebuildMutex sync.Mutex
}

func NewDevServer(cfg *config.GlobalConfig) (*DevServer, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &DevServer{
		cfg:     cfg,
		watcher: watcher,
		clients: make(map[chan []byte]struct{}),
	}, nil
}

func (s *DevServer) Start(port int) {
	if err := generator.BuildSite(s.cfg); err != nil {
		log.Fatalf("Initial build failed: %v", err)
	}

	go s.watchChanges()
	mux := http.NewServeMux()
	mux.Handle("/", &fileHandler{
		root:     http.Dir(s.cfg.OutputDir),
		notFound: "404.html",
		dev:      true,
	})
	mux.HandleFunc("/livereload", s.handleLiveReload)

	log.Printf("Server listening on http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), mux))
}

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
	defer f.Close()

	stat, _ := f.Stat()

	if stat.IsDir() {
		indexPath := filepath.Join(path, "index.html")
		indexFile, err := h.root.Open(indexPath)
		if err != nil {
			h.serve404(w, r)
			return
		}
		defer indexFile.Close()
		h.serveHTML(indexFile, w, r)
		return

	}

	if h.dev && strings.HasSuffix(path, ".html") {
		h.serveHTML(f, w, r)
		return
	}

	http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
}

func (h *fileHandler) serveHTML(f http.File, w http.ResponseWriter, _ *http.Request) {
	// Read the file content
	content, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Inject livereload script
	script := []byte(`<script>
	(new EventSource('/livereload')).onmessage = () => location.reload()
	</script></body>`)
	content = bytes.Replace(content, []byte("</body>"), script, 1)

	// Set headers and serve
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

func (h *fileHandler) serve404(w http.ResponseWriter, r *http.Request) {
	f, err := h.root.Open(h.notFound)
	if err == nil {
		defer f.Close()
		h.serveHTML(f, w, r)
		return
	}
	http.NotFound(w, r)
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

	messageChan := make(chan []byte)
	s.clientMutex.Lock()
	s.clients[messageChan] = struct{}{}
	s.clientMutex.Unlock()

	defer func() {
		s.clientMutex.Lock()
		delete(s.clients, messageChan)
		s.clientMutex.Unlock()
		close(messageChan)
	}()

	for {
		select {
		case <-r.Context().Done():
			return
		case msg := <-messageChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}

func (s *DevServer) notifyClients() {
	s.rebuildMutex.Lock()
	defer s.rebuildMutex.Unlock()

	for client := range s.clients {
		select {
		case client <- []byte("reload"):
		default:
		}
	}
}

func (s *DevServer) watchChanges() {
	defer s.watcher.Close()

	watchPaths := []string{
		"content",
		"templates",
		"static",
		filepath.Dir("config.toml"),
	}

	for _, path := range watchPaths {
		absPath, _ := filepath.Abs(path)
		if s.isIgnoredPath(absPath) {
			continue
		}

		filepath.Walk(absPath, func(p string, info os.FileInfo, err error) error {
			if err != nil || !info.IsDir() || s.isIgnoredPath(p) {
				return nil
			}
			if err := s.watcher.Add(p); err != nil {
				log.Printf("Watching %s: %v", p, err)
			}
			return nil
		})
	}

	debounceTime := 500 * time.Millisecond
	var timer *time.Timer

	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}

			// Skip ignored paths and non-content files
			if s.shouldIgnoreEvent(event) {
				continue
			}

			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(debounceTime, func() {
					s.triggerRebuild()
					s.notifyClients()
				})
			}

		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (s *DevServer) triggerRebuild() {
	s.clientMutex.Lock()
	defer s.clientMutex.Unlock()

	log.Println("Detected changes, rebuilding site...")
	if err := generator.BuildSite(s.cfg); err != nil {
		log.Printf("Rebuild error: %v", err)
	}
}

func (s *DevServer) shouldIgnoreEvent(event fsnotify.Event) bool {
	// Ignore output directory changes
	if s.isIgnoredPath(event.Name) {
		return true
	}

	// Ignore non-content files
	ext := filepath.Ext(event.Name)
	return !s.isRelevantExtension(ext)
}

func (s *DevServer) isRelevantExtension(ext string) bool {
	relevant := []string{".md", ".html", ".toml", ".css", ".js", ""} // "" for directories
	return slices.Contains(relevant, ext)
}

func (s *DevServer) isIgnoredPath(path string) bool {
	absOutput, _ := filepath.Abs(s.cfg.OutputDir)
	absPath, _ := filepath.Abs(path)

	// Check if path is in output directory
	if strings.HasPrefix(absPath, absOutput) {
		return true
	}

	// Check other ignored patterns
	ignored := []string{".git", "node_modules", "vendor"}
	for _, pattern := range ignored {
		if strings.Contains(absPath, pattern) {
			return true
		}
	}
	return false
}
