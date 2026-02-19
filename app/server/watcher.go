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
		changedFiles      []*FileChange
		changedFilesMutex sync.Mutex
	)

	paths := []string{
		filepath.Dir(s.configPath),
		s.contentDir,
		s.detector.templatesPath,
		s.detector.staticPath,
		s.detector.dataPath,
	}

	for _, path := range paths {
		err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() && !s.detector.shouldIgnoreDir(p) {
				return s.watcher.Add(p)
			}
			return nil
		})
		if err != nil {
			log.Printf("Error watching %s: %v", path, err)
		}
	}

	absConfigPath, absErr := filepath.Abs(s.configPath)
	if absErr != nil {
		log.Printf("Warning: Could not get absolute path for config %s: %v", s.configPath, absErr)
	} else {
		if watchErr := s.watcher.Add(absConfigPath); watchErr != nil {
			log.Printf("Warning: Could not watch config file %s: %v", absConfigPath, watchErr)
		}
	}

	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				continue
			}

			absEventPath, _ := filepath.Abs(event.Name)

			// Handle directories separately
			if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
				// Directory events: add new directories to watcher, removed dirs trigger clean rebuild
				if event.Op&fsnotify.Create == fsnotify.Create {
					if !s.detector.shouldIgnoreDir(event.Name) {
						if err := s.watcher.Add(event.Name); err != nil {
							log.Printf("Error watching new directory %s: %v", event.Name, err)
						}
					}
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Printf("Directory removed: %s - triggering clean rebuild", event.Name)
					s.triggerCleanRebuild()
				}
				continue
			}

			if strings.Contains(absEventPath, "/.git/") || strings.HasSuffix(absEventPath, "/.git") {
				continue
			}
			if s.detector.shouldIgnoreDir(event.Name) {
				continue
			}

			// Handle config file changes (re-watch after atomic replacement)
			absConfigPath, _ := filepath.Abs(s.configPath)
			if absEventPath == absConfigPath && event.Op&fsnotify.Create == fsnotify.Create {
				if err := s.watcher.Add(absConfigPath); err != nil {
					log.Printf("Error re-watching config file: %v", err)
				}
			}

			// Handle file removals - always trigger clean rebuild
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				log.Printf("File removed: %s - triggering clean rebuild", event.Name)
				s.triggerCleanRebuild()
				continue
			}

			change, err := s.detector.DetectChange(event.Name)
			if err != nil {
				log.Printf("Error detecting change for %s: %v", event.Name, err)
				continue
			}

			if change.Type == ChangeTypeIgnored {
				continue
			}

			changedFilesMutex.Lock()
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

func (s *DevServer) triggerRebuild(changes []*FileChange) {
	start := time.Now()

	if len(changes) == 0 {
		return
	}

	var (
		hasConfigChange bool
		hasDataChange   bool
		contentFiles    []string
		templateFiles   []string
		staticFiles     []string
	)

	for _, change := range changes {
		switch change.Type {
		case ChangeTypeConfig:
			hasConfigChange = true
			log.Printf("Processing config change: %s", change.RelPath)
		case ChangeTypeData:
			hasDataChange = true
			log.Printf("Processing data change: %s", change.RelPath)
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

	if hasDataChange {
		log.Println("Data changed - performing full rebuild")
		if err := s.reloadData(); err != nil {
			log.Printf("Data reload error: %v", err)
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

	if len(templateFiles) > 0 {
		log.Printf("Templates changed: %v", templateFiles)
		if err := s.gen.ReloadTemplates(); err != nil {
			log.Printf("Template reload error: %v", err)
			return
		}
	}

	if len(staticFiles) > 0 {
		log.Printf("Static files changed: %d files", len(staticFiles))
		for _, staticFile := range staticFiles {
			// Skip if file no longer exists (deleted/renamed)
			if _, err := os.Stat(staticFile); os.IsNotExist(err) {
				continue
			}
			if err := s.gen.CopyStaticFile(staticFile); err != nil {
				log.Printf("Error copying static file %s: %v", staticFile, err)
			}
		}
	}

	if len(contentFiles) > 0 {
		log.Printf("Content changed: %d files", len(contentFiles))
		for _, f := range contentFiles {
			log.Printf("  - Content file: %s", f)
		}

		// Re-parse content and do regular rebuild (not clean)
		parseStart := time.Now()
		if err := s.parser.Parse(s.contentDir); err != nil {
			log.Printf("Parse error: %v", err)
			return
		}
		log.Printf("Parse took %dms", time.Since(parseStart).Milliseconds())

		genStart := time.Now()
		if err := s.gen.Generate(s.parser.ContentMap["."]); err != nil {
			log.Printf("Generate error: %v", err)
			return
		}
		log.Printf("Generate took %dms", time.Since(genStart).Milliseconds())

		s.notifyClients()
		log.Printf("Rebuild completed in %dms", time.Since(start).Milliseconds())
		return
	}
	genStart := time.Now()
	if err := s.gen.Generate(s.parser.ContentMap["."]); err != nil {
		log.Printf("Generate error: %v", err)
		return
	}
	log.Printf("Generate took %dms", time.Since(genStart).Milliseconds())

	s.notifyClients()
	log.Printf("Rebuild completed in %dms", time.Since(start).Milliseconds())
}

func (s *DevServer) triggerCleanRebuild() {
	start := time.Now()

	log.Println("Performing clean rebuild (removing stale files)")
	parseStart := time.Now()
	if err := s.parser.Parse(s.contentDir); err != nil {
		log.Printf("Parse error: %v", err)
		return
	}
	log.Printf("Parse took %dms", time.Since(parseStart).Milliseconds())

	genStart := time.Now()
	if err := s.gen.GenerateClean(s.parser.ContentMap["."]); err != nil {
		log.Printf("Generate error: %v", err)
		return
	}
	log.Printf("Generate took %dms", time.Since(genStart).Milliseconds())

	s.notifyClients()
	log.Printf("Clean rebuild completed in %dms", time.Since(start).Milliseconds())
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

func (s *DevServer) reloadData() error {
	projectDir := filepath.Dir(s.configPath)

	if err := s.site.LoadDataFiles(projectDir); err != nil {
		return utils.WrapWithContext(utils.ErrContent, err, utils.ErrorContext{
			Operation: "reload_data_files",
			Component: "dev_server",
		})
	}

	if s.site.I18n != nil {
		if err := s.site.I18n.LoadTranslations(); err != nil {
			return utils.WrapWithContext(utils.ErrContent, err, utils.ErrorContext{
				Operation: "reload_i18n_translations",
				Component: "dev_server",
			})
		}
	}

	newGen, err := builder.NewBuilder(s.site, s.parser)
	if err != nil {
		return utils.WrapWithContext(utils.ErrTemplate, err, utils.ErrorContext{
			Operation: "recreate_builder_after_data_reload",
			Component: "dev_server",
		})
	}

	s.gen = newGen

	return nil
}
