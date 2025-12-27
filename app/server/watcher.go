// This file provides file watching and automatic rebuild for development.
// It uses a clean change detection system for flexible and maintainable file watching.
package server

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/tduyng/gozzi/app/builder"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/parser"
	"github.com/tduyng/gozzi/app/utils"
)

// watchChanges monitors file system changes and triggers rebuilds
func (s *DevServer) watchChanges() {
	defer func() {
		_ = s.watcher.Close()
	}()

	debounceDuration := 500 * time.Millisecond
	var (
		debounceTimer     *time.Timer
		lastRebuildTime   time.Time
		changedFiles      []*FileChange
		changedFilesMutex sync.Mutex
	)

	// Watch all necessary directories
	paths := []string{
		filepath.Dir(s.configPath),
		s.contentDir,
		"templates",
		"static",
	}

	for _, path := range paths {
		err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			// Only watch directories, fsnotify will detect file changes in them
			if d.IsDir() && !s.detector.shouldIgnoreDir(p) {
				return s.watcher.Add(p)
			}
			return nil
		})
		if err != nil {
			log.Printf("Error watching %s: %v", path, err)
		}
	}

	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				continue
			}

			// Handle new directory creation
			if event.Op&fsnotify.Create == fsnotify.Create {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					if !s.detector.shouldIgnoreDir(event.Name) {
						if err := s.watcher.Add(event.Name); err != nil {
							log.Printf("Error watching new directory %s: %v", event.Name, err)
						}
					}
				}
			}

			// Detect and classify the change
			change, err := s.detector.DetectChange(event.Name)
			if err != nil {
				log.Printf("Error detecting change for %s: %v", event.Name, err)
				continue
			}

			// Ignore if not a relevant change (includes unchanged files)
			if change.Type == ChangeTypeIgnored {
				continue
			}

			// Only add to changedFiles if it's actually different
			changedFilesMutex.Lock()
			// Deduplicate: check if this file is already in the pending changes
			alreadyPending := false
			for _, existing := range changedFiles {
				if existing.Path == change.Path {
					alreadyPending = true
					break
				}
			}
			if !alreadyPending {
				changedFiles = append(changedFiles, change)
			}
			changedFilesMutex.Unlock()

			// Reset debounce timer
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			debounceTimer = time.AfterFunc(debounceDuration, func() {
				if time.Since(lastRebuildTime) > debounceDuration {
					changedFilesMutex.Lock()
					changes := make([]*FileChange, len(changedFiles))
					copy(changes, changedFiles)
					changedFiles = changedFiles[:0]
					changedFilesMutex.Unlock()

					s.triggerRebuild(changes)
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

// triggerRebuild processes detected changes and rebuilds the site
func (s *DevServer) triggerRebuild(changes []*FileChange) {
	start := time.Now()

	if len(changes) == 0 {
		return
	}

	// Group changes by type
	var (
		hasConfigChange bool
		contentFiles    []string
		templateFiles   []string
		staticFiles     []string
	)

	for _, change := range changes {
		switch change.Type {
		case ChangeTypeConfig:
			hasConfigChange = true
			log.Printf("Processing config change: %s", change.RelPath)
		case ChangeTypeContent:
			contentFiles = append(contentFiles, change.Path)
			log.Printf("Processing content change: %s", change.RelPath)
		case ChangeTypeTemplate:
			templateFiles = append(templateFiles, change.RelPath)
			log.Printf("Processing template change: %s", change.RelPath)
		case ChangeTypeStatic:
			staticFiles = append(staticFiles, change.Path)
			log.Printf("Processing static file change: %s", change.RelPath)
		}
	}

	// Handle static file changes
	if len(staticFiles) > 0 && !hasConfigChange && len(contentFiles) == 0 && len(templateFiles) == 0 {
		// Only static files changed - copy them and notify clients
		for _, staticFile := range staticFiles {
			if err := s.gen.CopyStaticFile(staticFile); err != nil {
				log.Printf("Error copying static file %s: %v", staticFile, err)
			}
		}
		s.notifyClients()
		log.Printf("Static files updated in %dms", time.Since(start).Milliseconds())
		return
	}

	// If we have other changes along with static files, they'll be handled by the full/incremental build below

	// Handle config changes (requires full rebuild)
	if hasConfigChange {
		log.Println("Config changed - performing full rebuild")
		if err := s.reloadConfig(); err != nil {
			log.Printf("Config reload error: %v", err)
			return
		}
		if err := s.gen.Generate(s.parser.ContentMap["."]); err != nil {
			log.Printf("Generate error: %v", err)
			return
		}
		s.notifyClients()
		log.Printf("Full rebuild completed in %dms", time.Since(start).Milliseconds())
		return
	}

	// Handle template changes
	if len(templateFiles) > 0 {
		log.Printf("Templates changed: %v", templateFiles)
		if err := s.gen.ReloadTemplates(); err != nil {
			log.Printf("Template reload error: %v", err)
		} else {
			count := s.gen.InvalidateTemplateCache(templateFiles)
			log.Printf("Invalidated %d cached renders", count)
		}
	}

	// Handle content changes
	if len(contentFiles) > 0 {
		log.Printf("Content changed: %d files", len(contentFiles))

		// Snapshot taxonomy values BEFORE parsing
		oldTaxonomies := s.gen.SnapshotTaxonomyValues(contentFiles, s.contentDir)

		// Parse changed files
		parseStart := time.Now()
		if err := s.parser.ParseFiles(s.contentDir, contentFiles); err != nil {
			log.Printf("Parse error: %v", err)
			return
		}
		log.Printf("Parse took %dms", time.Since(parseStart).Milliseconds())

		// Generate - use incremental build when possible
		genStart := time.Now()
		var err error

		if len(templateFiles) == 0 {
			// Pure content change - use incremental build
			err = s.gen.GenerateWithOptions(s.parser.ContentMap["."], builder.GenerateOptions{
				Incremental:       true,
				ChangedFiles:      contentFiles,
				ContentDir:        s.contentDir,
				OldTaxonomyValues: oldTaxonomies,
			})
		} else {
			// Template changed - do full rebuild
			err = s.gen.Generate(s.parser.ContentMap["."])
		}

		if err != nil {
			log.Printf("Generate error: %v", err)
			return
		}

		log.Printf("Generate took %dms", time.Since(genStart).Milliseconds())
	}

	// Notify live reload clients
	s.notifyClients()
	log.Printf("Rebuild completed in %dms", time.Since(start).Milliseconds())
}

// reloadConfig reloads the site configuration and performs a full rebuild
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
