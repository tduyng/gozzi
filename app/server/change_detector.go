// This file implements the file change detection and classification system.
package server

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ChangeType represents the category of file change
type ChangeType int

const (
	ChangeTypeIgnored ChangeType = iota
	ChangeTypeConfig
	ChangeTypeContent
	ChangeTypeTemplate
	ChangeTypeStatic
)

func (ct ChangeType) String() string {
	switch ct {
	case ChangeTypeConfig:
		return "config"
	case ChangeTypeContent:
		return "content"
	case ChangeTypeTemplate:
		return "template"
	case ChangeTypeStatic:
		return "static"
	case ChangeTypeIgnored:
		return "ignored"
	default:
		return "unknown"
	}
}

// FileChange represents a detected file change with its classification
type FileChange struct {
	Path       string
	Type       ChangeType
	RelPath    string // Relative path from the appropriate base directory
	IsNew      bool
	ContentMD5 string
}

// ChangeDetector detects and classifies file changes
type ChangeDetector struct {
	contentDir string
	outputDir  string
	configPath string
	fileHashes map[string]string
}

// NewChangeDetector creates a new change detector
func NewChangeDetector(contentDir, outputDir, configPath string) *ChangeDetector {
	return &ChangeDetector{
		contentDir: contentDir,
		outputDir:  outputDir,
		configPath: configPath,
		fileHashes: make(map[string]string),
	}
}

// ClassifyChange determines what type of change this file represents
func (cd *ChangeDetector) ClassifyChange(path string) ChangeType {
	// Normalize paths for comparison
	absPath, _ := filepath.Abs(path)
	absOutput, _ := filepath.Abs(cd.outputDir)
	absContent, _ := filepath.Abs(cd.contentDir)
	absConfig, _ := filepath.Abs(cd.configPath)

	// Ignore output directory
	if strings.HasPrefix(absPath, absOutput) {
		return ChangeTypeIgnored
	}

	// Ignore system and temporary files
	base := filepath.Base(path)
	if cd.isSystemFile(base) {
		return ChangeTypeIgnored
	}

	ext := filepath.Ext(path)

	// Check if it's the config file
	if absPath == absConfig {
		return ChangeTypeConfig
	}

	// Check if it's in content directory
	if strings.HasPrefix(absPath, absContent) {
		// .md files are content
		if ext == ".md" {
			return ChangeTypeContent
		}
		// ALL other files in content directory (except system/temp files) are treated as static files
		// This includes images, videos, PDFs, etc. No need to hardcode extensions.
		// Examples: content/books/post-name/img/cover.webp, content/docs/file.pdf
		// The flexible approach: if it's not markdown and not a system file, copy it as static
		return ChangeTypeStatic
	}

	// Check if it's a template (must be under templates/ directory)
	templatesPath, _ := filepath.Abs("templates")
	if strings.HasPrefix(absPath, templatesPath) && ext == ".html" {
		return ChangeTypeTemplate
	}

	// Check if it's static assets (must be under static/ directory)
	staticPath, _ := filepath.Abs("static")
	if strings.HasPrefix(absPath, staticPath) {
		return ChangeTypeStatic
	}

	// Everything else is ignored
	return ChangeTypeIgnored
}

// isSystemFile checks if a filename should be ignored (system/temp files)
func (cd *ChangeDetector) isSystemFile(filename string) bool {
	// Hidden files
	if strings.HasPrefix(filename, ".") {
		return true
	}

	// Backup files
	if strings.HasPrefix(filename, "~") || strings.HasSuffix(filename, "~") {
		return true
	}

	// Editor temporary files
	if strings.HasSuffix(filename, ".swp") ||
		strings.HasSuffix(filename, ".swo") ||
		strings.HasSuffix(filename, ".swn") ||
		strings.HasSuffix(filename, ".tmp") ||
		strings.HasSuffix(filename, ".bak") ||
		strings.HasSuffix(filename, ".lock") {
		return true
	}

	// Vim/Neovim numbered backup files (e.g., "4913", "1234")
	// These have no extension and are purely numeric
	if filepath.Ext(filename) == "" && isNumericFilename(filename) {
		return true
	}

	// Log files
	if filepath.Ext(filename) == ".log" {
		return true
	}

	return false
}

// isNumericFilename checks if a filename consists only of digits
func isNumericFilename(filename string) bool {
	if len(filename) == 0 {
		return false
	}
	for _, char := range filename {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// shouldIgnoreDir checks if a directory should be ignored
func (cd *ChangeDetector) shouldIgnoreDir(path string) bool {
	absOutput, _ := filepath.Abs(cd.outputDir)
	absPath, _ := filepath.Abs(path)

	// Ignore output directory
	if strings.HasPrefix(absPath, absOutput) {
		return true
	}

	// Ignore hidden directories
	base := filepath.Base(path)
	return strings.HasPrefix(base, ".")
}

// HasChanged checks if a file's content has actually changed using MD5 hash
func (cd *ChangeDetector) HasChanged(path string) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		// File was deleted or can't be read
		delete(cd.fileHashes, path)
		return true, nil
	}

	newHash := fmt.Sprintf("%x", md5.Sum(content))
	oldHash, exists := cd.fileHashes[path]

	if !exists || oldHash != newHash {
		cd.fileHashes[path] = newHash
		return true, nil
	}

	return false, nil
}

// DetectChange combines classification and change detection
func (cd *ChangeDetector) DetectChange(path string) (*FileChange, error) {
	changeType := cd.ClassifyChange(path)

	if changeType == ChangeTypeIgnored {
		return &FileChange{
			Path: path,
			Type: ChangeTypeIgnored,
		}, nil
	}

	changed, err := cd.HasChanged(path)
	if err != nil {
		return nil, err
	}

	if !changed {
		return &FileChange{
			Path: path,
			Type: ChangeTypeIgnored, // No actual change
		}, nil
	}

	// Calculate relative path based on type
	var relPath string
	switch changeType {
	case ChangeTypeContent:
		relPath, _ = filepath.Rel(cd.contentDir, path)
	case ChangeTypeTemplate:
		relPath, _ = filepath.Rel("templates", path)
	case ChangeTypeStatic:
		relPath = path
	case ChangeTypeConfig:
		relPath = filepath.Base(path)
	}

	return &FileChange{
		Path:    path,
		Type:    changeType,
		RelPath: relPath,
	}, nil
}

// InitializeHashes pre-populates file hashes for watched directories
// This prevents spurious rebuilds on first file touch after server start
func (cd *ChangeDetector) InitializeHashes() error {
	paths := []string{
		cd.contentDir,
		"templates",
		"static",
		cd.configPath,
	}

	for _, path := range paths {
		if err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return nil // Skip paths with errors
			}

			// Skip directories
			if d.IsDir() {
				// Skip ignored directories
				if cd.shouldIgnoreDir(p) {
					return filepath.SkipDir
				}
				return nil
			}

			// Skip system/temp files
			if cd.isSystemFile(filepath.Base(p)) {
				return nil
			}

			// Only hash files we actually care about
			changeType := cd.ClassifyChange(p)
			if changeType == ChangeTypeIgnored {
				return nil
			}

			// Pre-compute and store hash
			content, err := os.ReadFile(p)
			if err != nil {
				return nil // Skip files we can't read
			}
			cd.fileHashes[p] = fmt.Sprintf("%x", md5.Sum(content))
			return nil
		}); err != nil {
			// Don't fail if a path doesn't exist
			if !os.IsNotExist(err) {
				return err
			}
		}
	}

	return nil
}
