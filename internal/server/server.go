package server

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/tduyng/gozzi/internal/config"
	"github.com/tduyng/gozzi/internal/generator"
	"github.com/tduyng/gozzi/internal/parser"
)

type DevServer struct {
	configPath     string
	contentDir     string
	site           *config.Site
	gen            *generator.Generator
	parser         *parser.ContentParser
	watcher        *fsnotify.Watcher
	clients        map[chan string]struct{}
	mu             sync.Mutex
	lastConfigHash string
}

func NewDevServer(configPath, contentDir string, site *config.Site, gen *generator.Generator, parser *parser.ContentParser) (*DevServer, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &DevServer{
		configPath: configPath,
		contentDir: contentDir,
		site:       site,
		gen:        gen,
		parser:     parser,
		watcher:    watcher,
		clients:    make(map[chan string]struct{}),
	}, nil
}

func (s *DevServer) Start(port int) {
	if err := s.initialize(); err != nil {
		log.Fatal(err)
	}

	go s.watchChanges()
	mux := http.NewServeMux()
	mux.Handle("/", &fileHandler{
		root:     http.Dir(s.site.OutputDir),
		notFound: "404.html",
		dev:      true,
	})
	mux.HandleFunc("/livereload", s.handleLiveReload)

	log.Printf("Server listening on http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), mux))
}

func (s *DevServer) initialize() error {
	if err := s.parser.Parse(s.contentDir); err != nil {
		return fmt.Errorf("initial content parse failed: %w", err)
	}
	return s.gen.Generate(s.parser.ContentMap["."])
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
	w.Write(content)
}

func (h *fileHandler) serve404(w http.ResponseWriter, r *http.Request) {
	f, err := h.root.Open(h.notFound)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
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
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}
}

func (s *DevServer) notifyClients() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for client := range s.clients {
		select {
		case client <- "reload":
		default:
		}
	}
}

func (s *DevServer) watchChanges() {
	defer s.watcher.Close()
	debounceDuration := 500 * time.Millisecond
	var (
		debounceTimer   *time.Timer
		lastRebuildTime time.Time
	)

	paths := []string{
		filepath.Dir(s.configPath),
		s.contentDir,
		"templates",
		"static",
	}

	for _, path := range paths {
		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil || !info.IsDir() || s.shouldIgnore(p) {
				return nil
			}
			return s.watcher.Add(p)
		})
		if err != nil {
			log.Printf("Error watching %s: %v", path, err)
		}
	}

	debounce := time.NewTimer(0)
	debounce.Stop()

	for {
		select {

		case event, ok := <-s.watcher.Events:
			if !ok || !s.isRelevantChange(event) {
				continue
			}

			// Reset timer on each relevant event
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			debounceTimer = time.AfterFunc(debounceDuration, func() {
				// Only rebuild if enough time passed since last
				if time.Since(lastRebuildTime) > debounceDuration {
					s.triggerRebuild()
					lastRebuildTime = time.Now()
				}
			})
		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (s *DevServer) shouldIgnore(path string) bool {
	absOutput, _ := filepath.Abs(s.site.OutputDir)
	absPath, _ := filepath.Abs(path)

	return strings.HasPrefix(absPath, absOutput) ||
		strings.Contains(path, "/.") || // Unix hidden
		strings.Contains(path, "\\.") // Windows hidden
}

func (s *DevServer) isRelevantChange(event fsnotify.Event) bool {
	// Handle config file using base name match
	base := filepath.Base(event.Name)
	if base == "config.toml" {
		return true
	}

	// Check valid extensions
	ext := filepath.Ext(event.Name)
	relevant := map[string]bool{
		".md":   true,
		".html": true,
		".css":  true,
		".js":   true,
	}
	return relevant[ext]
}

func (s *DevServer) triggerRebuild() {
	start := time.Now()

	if err := s.reloadConfig(); err != nil {
		log.Printf("Config reload error: %v", err)
	}

	if err := s.gen.ReloadTemplates(); err != nil {
		log.Printf("Template reload error: %v", err)
	}

	if err := s.parser.Parse(s.contentDir); err != nil {
		log.Printf("Content parse error: %v", err)
	}

	if err := s.gen.Generate(s.parser.ContentMap["."]); err != nil {
		log.Printf("Build error: %v", err)
	}

	s.notifyClients()
	log.Printf("Change detected, build done in %dms", time.Since(start).Milliseconds())
}

func (s *DevServer) reloadConfig() error {
	content, err := os.ReadFile(s.configPath)
	if err != nil {
		return fmt.Errorf("config read failed: %w", err)
	}

	newHash := fmt.Sprintf("%x", md5.Sum(content))
	if newHash == s.lastConfigHash {
		return nil
	}

	newSite, err := config.LoadSite(s.configPath)
	if err != nil {
		return fmt.Errorf("config reload failed: %w", err)
	}

	newSite.OutputDir = s.site.OutputDir // Maintain output directory

	newParser := parser.NewParser(newSite)
	if err := newParser.Parse(s.contentDir); err != nil {
		return fmt.Errorf("content re-parse failed: %w", err)
	}

	newGen, err := generator.NewGenerator(newSite, newParser)
	if err != nil {
		return fmt.Errorf("generator recreation failed: %w", err)
	}

	s.site = newSite
	s.parser = newParser
	s.gen = newGen
	s.lastConfigHash = newHash

	return nil
}
