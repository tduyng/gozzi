package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/tduyng/gozzi/internal/config"
	"github.com/tduyng/gozzi/internal/parser"
)

type SiteGenerator struct {
	cfg  *config.GlobalConfig
	tmpl *template.Template
}

type TemplateData struct {
	Config *config.MergedConfig
	Page   *parser.Page
	Site   *config.GlobalConfig
}

func NewSiteGenerator(cfg *config.GlobalConfig) (*SiteGenerator, error) {
	sg := &SiteGenerator{
		cfg: cfg,
	}

	tmpl, err := template.New("base").Funcs(template.FuncMap{
		"urlize":   URLize,
		"safeHTML": SafeHTML,
	}).ParseGlob("templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("template parsing failed: %w", err)
	}
	sg.tmpl = tmpl

	return sg, nil
}

func BuildSite(cfg *config.GlobalConfig) error {
	sg, err := NewSiteGenerator(cfg)
	if err != nil {
		return fmt.Errorf("site generator initialization failed: %w", err)
	}

	if err := os.RemoveAll(cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to clean output directory: %w", err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 100)

	err = filepath.Walk("content", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("access error for %q: %w", path, err)
		}

		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			if err := sg.processPage(p); err != nil {
				errChan <- fmt.Errorf("processing %q: %w", p, err)
			}
		}(path)

		return nil
	})
	if err != nil {
		return fmt.Errorf("content directory traversal failed: %w", err)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		log.Printf("Processing error: %v", err)
	}

	if err := sg.copyStatic(); err != nil {
		return fmt.Errorf("static assets copy failed: %w", err)
	}

	return nil
}

func (sg *SiteGenerator) processPage(path string) error {
	page, err := parser.ParseMarkdown(path)
	if err != nil {
		return fmt.Errorf("markdown parsing failed: %w", err)
	}

	sectionPath := filepath.Dir(path)
	sectionCfg, err := config.LoadSectionConfig(sectionPath)
	if err != nil {
		return fmt.Errorf("section config load failed: %w", err)
	}

	mergedConfig := config.MergeConfigs(sg.cfg, sectionCfg, &page.FrontMatter)

	// Get relative path from content directory
	relPath, err := filepath.Rel("content", filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	// Determine output path based on file type
	filename := filepath.Base(path)
	var outputPath string

	switch {
	case filename == "_index.md":
		outputPath = filepath.Join(sg.cfg.OutputDir, relPath, "index.html")
	case strings.HasSuffix(filename, ".md"):
		outputPath = filepath.Join(sg.cfg.OutputDir, relPath, page.Slug, "index.html")
	default:
		return fmt.Errorf("unsupported file type: %s", filename)
	}

	return sg.renderPage(outputPath, page, mergedConfig)
}

func (sg *SiteGenerator) renderPage(outputPath string, page *parser.Page, cfg *config.MergedConfig) error {
	// Create output directory structure
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("directory creation failed: %w", err)
	}

	// Find appropriate template
	tmpl := sg.tmpl.Lookup(cfg.Template)
	if tmpl == nil {
		return fmt.Errorf("template %q not found", cfg.Template)
	}

	// Execute template with context
	var buf bytes.Buffer
	data := TemplateData{
		Config: cfg,
		Page:   page,
		Site:   sg.cfg,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("file write failed: %w", err)
	}

	return nil
}

func URLize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")

	return strings.Trim(s, "-")
}

func SafeHTML(s string) template.HTML {
	return template.HTML(s)
}

func (sg *SiteGenerator) copyStatic() error {
	return filepath.Walk(filepath.Join("static"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("static file access error: %w", err)
		}

		if info.IsDir() || isHiddenFile(path) || shouldSkip(path) {
			return nil
		}

		relPath, err := filepath.Rel("static", path)
		if err != nil {
			return fmt.Errorf("path resolution failed: %w", err)
		}

		destPath := filepath.Join(sg.cfg.OutputDir, relPath)

		if err := copyFile(path, destPath); err != nil {
			return fmt.Errorf("file copy failed: %w", err)
		}

		return nil
	})
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("source open failed: %w", err)
	}
	defer source.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("destination create failed: %w", err)
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return fmt.Errorf("copy operation failed: %w", err)
	}

	return nil
}

func isHiddenFile(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
}

func shouldSkip(path string) bool {
	ext := filepath.Ext(path)
	dir := filepath.Dir(path)
	return strings.Contains(dir, "scss") && ext == ".scss"
}
