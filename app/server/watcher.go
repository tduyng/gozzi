// Package server provides file watching and automatic rebuild for development.
package server

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/tduyng/gozzi/app/builder"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/parser"
	"github.com/tduyng/gozzi/app/utils"
)

func (s *DevServer) watchChanges() {
	defer func() {
		_ = s.watcher.Close()
	}()
	debounceDuration := 500 * time.Millisecond
	var (
		debounceTimer     *time.Timer
		lastRebuildTime   time.Time
		changedFiles      []string
		changedFilesMutex sync.Mutex
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

			// Check if file content actually changed (skip no-op writes like vim :w)
			if !s.hasFileChanged(event.Name) {
				continue
			}

			// Track which files changed
			changedFilesMutex.Lock()
			changedFiles = append(changedFiles, event.Name)
			changedFilesMutex.Unlock()

			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			debounceTimer = time.AfterFunc(debounceDuration, func() {
				if time.Since(lastRebuildTime) > debounceDuration {
					changedFilesMutex.Lock()
					files := make([]string, len(changedFiles))
					copy(files, changedFiles)
					changedFiles = changedFiles[:0]
					changedFilesMutex.Unlock()

					s.triggerRebuild(files)
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
		strings.Contains(path, "/.") ||
		strings.Contains(path, "\\.")
}

func (s *DevServer) isRelevantChange(event fsnotify.Event) bool {
	base := filepath.Base(event.Name)
	if base == "config.toml" {
		return true
	}

	ext := filepath.Ext(event.Name)
	relevant := map[string]bool{
		".md":   true,
		".html": true,
		".css":  true,
		".js":   true,
	}
	return relevant[ext]
}

func (s *DevServer) triggerRebuild(changedFiles []string) {
	start := time.Now()

	// If no files changed, skip rebuild
	if len(changedFiles) == 0 {
		return
	}

	// Categorize changed files
	hasConfigChange := false
	contentFiles := []string{}
	templateFiles := []string{}

	for _, file := range changedFiles {
		switch {
		case strings.Contains(file, "templates") && filepath.Ext(file) == ".html":
			// Track specific template files for selective cache invalidation
			if relPath, err := filepath.Rel("templates", file); err == nil {
				templateFiles = append(templateFiles, filepath.ToSlash(relPath))
			}
		case filepath.Base(file) == "config.toml":
			hasConfigChange = true
		case strings.Contains(file, "static"):
			// Static files will be copied during generate
		case strings.Contains(file, s.contentDir) && (filepath.Ext(file) == ".md" || filepath.Base(file) == "_index.md"):
			contentFiles = append(contentFiles, file)
		}
	}

	// Only reload config if it actually changed
	if hasConfigChange {
		if err := s.reloadConfig(); err != nil {
			log.Printf("Config reload error: %v", err)
		}
	}

	// Handle template changes with selective cache invalidation
	if len(templateFiles) > 0 {
		if err := s.gen.ReloadTemplates(); err != nil {
			log.Printf("Template reload error: %v", err)
		} else {
			// Invalidate only the templates that changed
			count := s.gen.InvalidateTemplateCache(templateFiles)
			if count > 0 {
				log.Printf("Invalidated %d cached render(s) for templates: %v", count, templateFiles)
			}
		}
	}

	// Reset stats before parsing
	s.parser.ResetStats()

	// Incremental content parsing: only parse changed markdown files
	if len(contentFiles) > 0 {
		if err := s.parser.ParseFiles(s.contentDir, contentFiles); err != nil {
			log.Printf("Content parse error: %v", err)
		}
	} else if hasConfigChange {
		// For config changes, do a full parse since they may affect all pages
		if err := s.parser.Parse(s.contentDir); err != nil {
			log.Printf("Content parse error: %v", err)
		}
	}
	// Note: Template and static changes don't require re-parsing markdown content.
	// The existing parsed content will be re-rendered with the new templates/assets.

	if err := s.gen.Generate(s.parser.ContentMap["."]); err != nil {
		log.Printf("Build error: %v", err)
	}

	s.notifyClients()

	// Log build completion with optional debug stats
	if s.site.Debug {
		// Log incremental parsing and render cache stats (debug mode only)
		stats := s.parser.GetStats()
		hashCacheStats := s.parser.GetHashCache().Stats()
		renderCacheStats := s.gen.GetRenderCacheStats()

		total := stats.TotalFiles.Load()
		skipped := stats.FilesSkipped.Load()
		parsed := stats.FilesParsed.Load()

		if total > 0 {
			skipRate := float64(skipped) / float64(total) * 100
			log.Printf("Change detected, build done in %dms (parsed: %d, skipped: %d/%.0f%%, %s, %s)",
				time.Since(start).Milliseconds(),
				parsed,
				skipped,
				skipRate,
				hashCacheStats.String(),
				renderCacheStats.String(),
			)
		} else {
			log.Printf("Change detected, build done in %dms (%s)",
				time.Since(start).Milliseconds(),
				renderCacheStats.String(),
			)
		}
	} else {
		// Simple output for normal mode
		log.Printf("Change detected, build done in %dms",
			time.Since(start).Milliseconds(),
		)
	}
}

// hasFileChanged checks if a file's content has actually changed by comparing hashes.
// Returns true if the file is new or its content changed, false for no-op writes.
func (s *DevServer) hasFileChanged(filePath string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		// If we can't read the file, assume it changed (e.g., deleted)
		return true
	}

	newHash := fmt.Sprintf("%x", md5.Sum(content))

	s.fileHashesMu.Lock()
	defer s.fileHashesMu.Unlock()

	oldHash, exists := s.fileHashes[filePath]
	if !exists || oldHash != newHash {
		// File is new or changed, update hash
		s.fileHashes[filePath] = newHash
		return true
	}

	// File content unchanged (no-op write from vim :w, etc.)
	return false
}

func (s *DevServer) reloadConfig() error {
	content, err := os.ReadFile(s.configPath)
	if err != nil {
		return utils.WrapWithContext(utils.ErrFileSystem, err, utils.ErrorContext{
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
		return utils.WrapWithContext(utils.ErrConfig, err, utils.ErrorContext{
			Operation: "reload_config",
			Component: "dev_server",
			Path:      s.configPath,
		})
	}

	newSite.OutputDir = s.site.OutputDir

	newParser := parser.NewParser(newSite)
	if err := newParser.Parse(s.contentDir); err != nil {
		return utils.WrapWithContext(utils.ErrContent, err, utils.ErrorContext{
			Operation: "reparse_content_after_config_reload",
			Component: "dev_server",
			Path:      s.contentDir,
		})
	}

	newGen, err := builder.NewBuilder(newSite, newParser)
	if err != nil {
		return utils.WrapWithContext(utils.ErrTemplate, err, utils.ErrorContext{
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
