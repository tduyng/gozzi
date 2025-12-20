// Package server provides a development server with live reload functionality.
package server

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/tduyng/gozzi/app/builder"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/parser"
	"github.com/tduyng/gozzi/app/utils"
)

// DevServer provides a development server with file watching and live reload.
type DevServer struct {
	configPath     string
	contentDir     string
	site           *config.Site
	gen            *builder.Builder
	parser         *parser.ContentParser
	watcher        *fsnotify.Watcher
	clients        map[chan string]struct{}
	mu             sync.Mutex
	lastConfigHash string
	fileHashes     map[string]string // Track file content hashes to detect actual changes
	fileHashesMu   sync.RWMutex      // Protect fileHashes map from concurrent access
}

// NewDevServer creates a new development server with file watching enabled.
func NewDevServer(
	configPath, contentDir string,
	site *config.Site,
	gen *builder.Builder,
	parser *parser.ContentParser,
) (*DevServer, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, utils.WrapWithContext(utils.ErrServer, err, utils.ErrorContext{
			Operation: "create_file_watcher",
			Component: "dev_server",
		})
	}

	return &DevServer{
		configPath: configPath,
		contentDir: contentDir,
		site:       site,
		gen:        gen,
		parser:     parser,
		watcher:    watcher,
		clients:    make(map[chan string]struct{}),
		fileHashes: make(map[string]string),
	}, nil
}

// Start starts the development server on the specified port.
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
	// Generate initial site - content has already been parsed by initApp in main.go
	if err := s.gen.Generate(s.parser.ContentMap["."]); err != nil {
		return utils.WrapWithContext(utils.ErrContent, err, utils.ErrorContext{
			Operation: "initial_site_generation",
			Component: "dev_server",
		})
	}

	return nil
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
