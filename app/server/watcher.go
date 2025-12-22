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
			if !ok {
				continue
			}

			// Handle directory creation - add new directories to watch list
			if event.Op&fsnotify.Create == fsnotify.Create {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() && !s.shouldIgnore(event.Name) {
					if err := s.watcher.Add(event.Name); err != nil {
						log.Printf("Error watching new directory %s: %v", event.Name, err)
					} else {
						log.Printf("Now watching new directory: %s", event.Name)
					}
				}
			}

			if !s.isRelevantChange(event) {
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
		// Content and templates
		".md":   true,
		".html": true,
		// Stylesheets and scripts
		".css": true,
		".js":  true,
		// Images
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
		".webp": true,
		".svg":  true,
		// Other assets
		".ico":  true,
		".xml":  true,
		".json": true,
		".txt":  true,
	}
	return relevant[ext]
}

// findParentMarkdownFile finds the markdown file that owns an asset
// For content/blog/post-name/img/image.png -> content/blog/post-name/index.md
func (s *DevServer) findParentMarkdownFile(assetPath string) string {
	// Get directory of the asset (e.g., content/blog/post/img)
	dir := filepath.Dir(assetPath)

	// Walk up the directory tree looking for index.md
	for dir != s.contentDir && dir != "." && dir != "/" {
		// Go up one level (from img -> post-name)
		parentDir := filepath.Dir(dir)

		// Check for index.md in this directory
		indexPath := filepath.Join(parentDir, "index.md")
		if _, err := os.Stat(indexPath); err == nil {
			return indexPath
		}

		// Move up for next iteration
		dir = parentDir
	}

	// No markdown file found - return empty string to skip
	return ""
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
	assetFiles := []string{}

	for _, file := range changedFiles {
		switch {
		case strings.Contains(file, "templates") && filepath.Ext(file) == ".html":
			if relPath, err := filepath.Rel("templates", file); err == nil {
				templateFiles = append(templateFiles, filepath.ToSlash(relPath))
			}
		case filepath.Base(file) == "config.toml":
			hasConfigChange = true
		case strings.HasPrefix(file, "static/") || strings.Contains(file, "/static/"):
			assetFiles = append(assetFiles, file)
		case strings.Contains(file, s.contentDir) && (filepath.Ext(file) == ".md" || filepath.Base(file) == "_index.md"):
			contentFiles = append(contentFiles, file)
		case strings.Contains(file, s.contentDir):
			parentMd := s.findParentMarkdownFile(file)
			if parentMd != "" {
				contentFiles = append(contentFiles, parentMd)
			} else {
				assetFiles = append(assetFiles, file)
			}
		}
	}

	log.Printf("Rebuild triggered - config: %v, content: %d, templates: %d, assets: %d",
		hasConfigChange, len(contentFiles), len(templateFiles), len(assetFiles))

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

	// Reset cache stats to get per-build statistics
	s.gen.ResetCacheStats()

	// Incremental content parsing: only parse changed markdown files
	parseStart := time.Now()
	if len(contentFiles) > 0 {
		if err := s.parser.ParseFiles(s.contentDir, contentFiles); err != nil {
			log.Printf("Content parse error: %v", err)
		}
		log.Printf("ParseFiles took %dms", time.Since(parseStart).Milliseconds())
	} else if hasConfigChange {
		// For config changes, do a full parse since they may affect all pages
		// CRITICAL: Clear hash cache so all files are reparsed even if content hasn't changed
		// Config changes can affect all pages (e.g., site.extra.img used in templates)
		s.parser.GetHashCache().Clear()
		if err := s.parser.Parse(s.contentDir); err != nil {
			log.Printf("Content parse error: %v", err)
		}
		log.Printf("Full parse took %dms", time.Since(parseStart).Milliseconds())
	}
	// Note: Template and static changes don't require re-parsing markdown content.
	// The existing parsed content will be re-rendered with the new templates/assets.

	genStart := time.Now()
	var err error

	// Use incremental builds when only content files changed (no config/template changes)
	if len(contentFiles) > 0 && !hasConfigChange && len(templateFiles) == 0 {
		// Incremental build: only regenerate changed content and affected dependencies
		err = s.gen.GenerateWithOptions(s.parser.ContentMap["."], builder.GenerateOptions{
			Incremental:  true,
			ChangedFiles: contentFiles,
			ContentDir:   s.contentDir,
		})
		if err != nil {
			log.Printf("Incremental build error: %v", err)
		}
	} else {
		// Full build: config or template changes affect everything
		err = s.gen.Generate(s.parser.ContentMap["."])
		if err != nil {
			log.Printf("Full build error: %v", err)
		}
	}
	genDuration := time.Since(genStart).Milliseconds()

	// Log cache stats to understand how well incremental builds are working
	cacheStats := s.gen.GetCacheStats()
	buildType := "full"
	if len(contentFiles) > 0 && !hasConfigChange && len(templateFiles) == 0 {
		buildType = "incremental"
	}
	log.Printf("Generate (%s) took %dms (cache: %d hits, %d misses, %.1f%% hit rate)",
		buildType,
		genDuration,
		cacheStats.Hits,
		cacheStats.Misses,
		cacheStats.HitRate)

	s.notifyClients()

	// Log build completion
	log.Printf("Change detected, build done in %dms", time.Since(start).Milliseconds())
}

// hasFileChanged checks if a file's content has actually changed by comparing hashes.
// Returns true if the file is new or its content changed, false for no-op writes.
func (s *DevServer) hasFileChanged(filePath string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		// If we can't read the file, assume it changed (e.g., deleted)
		// Clean up the hash for deleted files
		s.fileHashesMu.Lock()
		delete(s.fileHashes, filePath)
		s.fileHashesMu.Unlock()
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
