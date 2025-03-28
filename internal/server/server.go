package server

import (
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
	cfg      *config.GlobalConfig
	watcher  *fsnotify.Watcher
	debounce *time.Timer
	mu       sync.Mutex
}

func NewDevServer(cfg *config.GlobalConfig) (*DevServer, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &DevServer{
		cfg:     cfg,
		watcher: watcher,
	}, nil
}

func (s *DevServer) Start(port int) {
	if err := generator.BuildSite(s.cfg); err != nil {
		log.Fatalf("Initial build failed: %v", err)
	}

	go s.watchChanges()

	handler := &fileHandler{
		root:     http.Dir(s.cfg.OutputDir),
		notFound: "404.html",
	}

	log.Printf("Server listening on http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), handler))
}

type fileHandler struct {
	root     http.Dir
	notFound string
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

		indexStat, _ := indexFile.Stat()
		http.ServeContent(w, r, indexStat.Name(), indexStat.ModTime(), indexFile)
		return
	}

	http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
}

func (h *fileHandler) serve404(w http.ResponseWriter, r *http.Request) {
	f, err := h.root.Open(h.notFound)
	if err == nil {
		defer f.Close()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		io.Copy(w, f)
		return
	}
	http.NotFound(w, r)
}

func (s *DevServer) triggerRebuild() {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Println("Detected changes, rebuilding site...")
	if err := generator.BuildSite(s.cfg); err != nil {
		log.Printf("Rebuild error: %v", err)
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
				timer = time.AfterFunc(debounceTime, s.triggerRebuild)
			}

		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
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
