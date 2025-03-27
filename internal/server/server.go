package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/tduyng/gozzi/internal/config"
	"github.com/tduyng/gozzi/internal/generator"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type LiveReloadServer struct {
	clients    map[*websocket.Conn]struct{}
	mu         sync.Mutex
	watcher    *fsnotify.Watcher
	cfg        *config.GlobalConfig
	debouncer  *time.Timer
	server     *http.Server
	wsUpgrader websocket.Upgrader
}

func NewLiveReloadServer(cfg *config.GlobalConfig) (*LiveReloadServer, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	lrs := &LiveReloadServer{
		clients:    make(map[*websocket.Conn]struct{}),
		watcher:    watcher,
		cfg:        cfg,
		wsUpgrader: upgrader,
	}

	// Watch relevant directories
	watchDirs := []string{"content", "templates", "static", "config.toml"}
	for _, dir := range watchDirs {
		if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil || !info.IsDir() || lrs.isIgnoredDir(path) {
				return nil
			}
			return watcher.Add(path)
		}); err != nil {
			return nil, fmt.Errorf("failed to add watch path: %w", err)
		}
	}

	return lrs, nil
}

func (lrs *LiveReloadServer) Start(port int) error {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(lrs.cfg.OutputDir)))
	mux.HandleFunc("/livereload", lrs.handleWebSocket)

	lrs.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		log.Printf("Starting server on http://localhost:%d", port)
		if err := lrs.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	go lrs.watchChanges()

	return nil
}

func (lrs *LiveReloadServer) watchChanges() {
	defer lrs.watcher.Close()

	for {
		select {
		case event, ok := <-lrs.watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				lrs.handleChange(event)
			}
		case err, ok := <-lrs.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (lrs *LiveReloadServer) handleChange(_ fsnotify.Event) {
	lrs.mu.Lock()
	defer lrs.mu.Unlock()

	if lrs.debouncer != nil {
		lrs.debouncer.Stop()
	}

	lrs.debouncer = time.AfterFunc(500*time.Millisecond, func() {
		log.Println("Change detected, rebuilding site...")
		if err := generator.BuildSite(lrs.cfg); err != nil {
			log.Printf("Rebuild error: %v", err)
		}
		lrs.notifyClients()
	})
}

func (lrs *LiveReloadServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := lrs.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	lrs.mu.Lock()
	lrs.clients[conn] = struct{}{}
	lrs.mu.Unlock()

	// Keep connection alive
	for {
		if _, _, err := conn.NextReader(); err != nil {
			lrs.mu.Lock()
			delete(lrs.clients, conn)
			lrs.mu.Unlock()
			break
		}
	}
}

func (lrs *LiveReloadServer) notifyClients() {
	lrs.mu.Lock()
	defer lrs.mu.Unlock()

	for client := range lrs.clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
			log.Printf("Error sending reload message: %v", err)
			client.Close()
			delete(lrs.clients, client)
		}
	}
}

func (lrs *LiveReloadServer) Shutdown(ctx context.Context) error {
	lrs.mu.Lock()
	defer lrs.mu.Unlock()

	for client := range lrs.clients {
		client.Close()
		delete(lrs.clients, client)
	}

	if lrs.debouncer != nil {
		lrs.debouncer.Stop()
	}

	if lrs.server != nil {
		return lrs.server.Shutdown(ctx)
	}
	return nil
}

func (lrs *LiveReloadServer) isIgnoredDir(path string) bool {
	ignored := []string{".git", "node_modules", "vendor", lrs.cfg.OutputDir}
	for _, dir := range ignored {
		if strings.Contains(path, dir) {
			return true
		}
	}
	return false
}
