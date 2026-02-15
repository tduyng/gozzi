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
	contentDir    string
	outputDir     string
	configPath    string
	templatesPath string
	staticPath    string
	fileHashes    map[string]string
}

// NewChangeDetector creates a new change detector
func NewChangeDetector(contentDir, outputDir, configPath string) *ChangeDetector {
	// Compute absolute paths to ensure consistent comparison regardless of CWD
	absConfig, _ := filepath.Abs(configPath)
	projectRoot := filepath.Dir(absConfig)

	// Resolve contentDir and outputDir to absolute paths
	absContentDir, _ := filepath.Abs(contentDir)
	absOutputDir, _ := filepath.Abs(outputDir)

	templatesPath := filepath.Join(projectRoot, "templates")
	staticPath := filepath.Join(projectRoot, "static")

	return &ChangeDetector{
		contentDir:    absContentDir,
		outputDir:     absOutputDir,
		configPath:    absConfig,
		templatesPath: templatesPath,
		staticPath:    staticPath,
		fileHashes:    make(map[string]string),
	}
}

// ClassifyChange determines what type of change this file represents
func (cd *ChangeDetector) ClassifyChange(path string) ChangeType {
	absPath, _ := filepath.Abs(path)
	absOutput, _ := filepath.Abs(cd.outputDir)
	absContent, _ := filepath.Abs(cd.contentDir)
	absConfig, _ := filepath.Abs(cd.configPath)

	if strings.HasPrefix(absPath, absOutput) {
		return ChangeTypeIgnored
	}

	if strings.Contains(absPath, "/.git/") || strings.HasSuffix(absPath, "/.git") {
		return ChangeTypeIgnored
	}

	base := filepath.Base(path)
	if cd.isSystemFile(base) {
		return ChangeTypeIgnored
	}

	ext := filepath.Ext(path)

	if absPath == absConfig {
		return ChangeTypeConfig
	}

	if strings.HasPrefix(absPath, absContent) {
		if ext == ".md" {
			return ChangeTypeContent
		}
		return ChangeTypeStatic
	}

	templatesPath, _ := filepath.Abs(cd.templatesPath)
	staticPath, _ := filepath.Abs(cd.staticPath)
	if strings.HasPrefix(absPath, templatesPath) && ext == ".html" {
		return ChangeTypeTemplate
	}

	if strings.HasPrefix(absPath, staticPath) {
		return ChangeTypeStatic
	}

	return ChangeTypeIgnored
}

func (cd *ChangeDetector) isSystemFile(filename string) bool {
	if strings.HasPrefix(filename, ".") {
		return true
	}

	if strings.HasPrefix(filename, "~") || strings.HasSuffix(filename, "~") {
		return true
	}

	if strings.HasSuffix(filename, ".swp") ||
		strings.HasSuffix(filename, ".swo") ||
		strings.HasSuffix(filename, ".swn") ||
		strings.HasSuffix(filename, ".tmp") ||
		strings.HasSuffix(filename, ".bak") ||
		strings.HasSuffix(filename, ".lock") {
		return true
	}

	if filepath.Ext(filename) == "" && isNumericFilename(filename) {
		return true
	}

	if filepath.Ext(filename) == ".log" {
		return true
	}

	return false
}

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

func (cd *ChangeDetector) shouldIgnoreDir(path string) bool {
	absOutput, _ := filepath.Abs(cd.outputDir)
	absPath, _ := filepath.Abs(path)

	if strings.HasPrefix(absPath, absOutput) {
		return true
	}

	base := filepath.Base(path)

	return strings.HasPrefix(base, ".")
}

func (cd *ChangeDetector) HasChanged(path string) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
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

	var relPath string
	switch changeType {
	case ChangeTypeContent:
		relPath, _ = filepath.Rel(cd.contentDir, path)
	case ChangeTypeTemplate:
		relPath, _ = filepath.Rel(cd.templatesPath, path)
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

func (cd *ChangeDetector) InitializeHashes() error {
	paths := []string{
		cd.contentDir,
		cd.templatesPath,
		cd.staticPath,
		cd.configPath,
	}

	for _, path := range paths {
		if err := filepath.WalkDir(path, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return nil // Skip paths with errors
			}

			if d.IsDir() {
				if cd.shouldIgnoreDir(p) {
					return filepath.SkipDir
				}
				return nil
			}

			if cd.isSystemFile(filepath.Base(p)) {
				return nil
			}

			changeType := cd.ClassifyChange(p)
			if changeType == ChangeTypeIgnored {
				return nil
			}

			content, err := os.ReadFile(p)
			if err != nil {
				return nil
			}
			cd.fileHashes[p] = fmt.Sprintf("%x", md5.Sum(content))
			return nil
		}); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		}
	}

	return nil
}
