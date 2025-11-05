// Package server provides file watching and automatic rebuild functionality for development server.
// Monitors content, templates, static files, and config for changes and triggers rebuilds.
package server

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/tduyng/gozzi/app/builder"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/parser"
	"github.com/tduyng/gozzi/shared"
)

func (s *DevServer) watchChanges() {
	defer func() {
		_ = s.watcher.Close()
	}()
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
		err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil || !d.IsDir() || s.shouldIgnore(p) {
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
		return shared.WrapWithContext(shared.ErrFileSystem, err, shared.ErrorContext{
			Operation: "read_config_file",
			Component: "dev_server",
			Path:      s.configPath,
		})
	}

	newHash := fmt.Sprintf("%x", md5.Sum(content))
	if newHash == s.lastConfigHash {
		return nil
	}

	newSite, err := config.LoadSite(s.configPath)
	if err != nil {
		return shared.WrapWithContext(shared.ErrConfig, err, shared.ErrorContext{
			Operation: "reload_config",
			Component: "dev_server",
			Path:      s.configPath,
		})
	}

	newSite.OutputDir = s.site.OutputDir // Maintain output directory

	newParser := parser.NewParser(newSite)
	if err := newParser.Parse(s.contentDir); err != nil {
		return shared.WrapWithContext(shared.ErrContent, err, shared.ErrorContext{
			Operation: "reparse_content_after_config_reload",
			Component: "dev_server",
			Path:      s.contentDir,
		})
	}

	newGen, err := builder.NewBuilder(newSite, newParser)
	if err != nil {
		return shared.WrapWithContext(shared.ErrTemplate, err, shared.ErrorContext{
			Operation: "recreate_builder_after_config_reload",
			Component: "dev_server",
		})
	}

	s.site = newSite
	s.parser = newParser
	s.gen = newGen
	s.lastConfigHash = newHash

	return nil
}
