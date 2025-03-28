package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
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
	// Initial build
	if err := generator.BuildSite(s.cfg); err != nil {
		log.Fatalf("Initial build failed: %v", err)
	}

	// Setup file watcher
	go s.watchChanges()

	// Start HTTP server
	fs := http.FileServer(http.Dir(s.cfg.OutputDir))
	log.Printf("Server listening on http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), fs))
}

func (s *DevServer) watchChanges() {
	defer s.watcher.Close()

	// Watch important directories and config file
	watchPaths := []string{
		"content",
		"templates",
		"static",
		"config.toml",
	}

	for _, path := range watchPaths {
		if isFile(path) {
			// Add parent directory for files
			dir := filepath.Dir(path)
			if err := s.watcher.Add(dir); err != nil {
				log.Printf("Failed to watch %s: %v", dir, err)
			}
		} else {
			// Add directory recursively
			filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
				if err == nil && info.IsDir() {
					if err := s.watcher.Add(p); err != nil {
						log.Printf("Failed to watch %s: %v", p, err)
					}
				}
				return nil
			})
		}
	}

	// Debounce rebuilds
	lastRebuild := time.Now()
	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}

			// Only trigger on writes and creates
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				// Check if it's config.toml or in watched directories
				base := filepath.Base(event.Name)
				if base == "config.toml" || contains(watchPaths, filepath.Dir(event.Name)) {
					// Debounce logic
					if time.Since(lastRebuild) > 1*time.Second {
						s.triggerRebuild()
						lastRebuild = time.Now()
					}
				}
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
	s.mu.Lock()
	defer s.mu.Unlock()

	// Cancel any pending rebuild
	if s.debounce != nil {
		s.debounce.Stop()
	}

	// Schedule new rebuild with debounce
	s.debounce = time.AfterFunc(500*time.Millisecond, func() {
		log.Println("Detected changes, rebuilding site...")
		if err := generator.BuildSite(s.cfg); err != nil {
			log.Printf("Rebuild error: %v", err)
		}
	})
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
