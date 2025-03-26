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
	"time"

	"github.com/bep/godartsass/v2"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tduyng/gozzi/internal/config"
	"github.com/tduyng/gozzi/internal/parser"
)

type SiteGenerator struct {
	cfg      *config.GlobalConfig
	minifier *minify.M
	sass     *godartsass.Transpiler
	tmpl     *template.Template
}

func NewSiteGenerator(cfg *config.GlobalConfig) (*SiteGenerator, error) {
	sg := &SiteGenerator{
		cfg:      cfg,
		minifier: minify.New(),
	}

	sg.minifier.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
		KeepEndTags:      true,
	})

	if cfg.CompileSass {
		transpiler, err := godartsass.Start(godartsass.Options{
			DartSassEmbeddedFilename: "", // Use system default
			Timeout:                  30 * time.Second,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Sass compiler: %w", err)
		}
		sg.sass = transpiler
	}

	tmpl, err := template.New("base").Funcs(template.FuncMap{
		"urlize":   URLize,
		"safeHTML": SafeHTML,
	}).ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		return nil, fmt.Errorf("template parsing failed: %w", err)
	}
	sg.tmpl = tmpl

	return sg, nil
}

func (sg *SiteGenerator) Close() error {
	if sg.sass != nil {
		if err := sg.sass.Close(); err != nil {
			return fmt.Errorf("error closing Sass transpiler: %w", err)
		}
	}
	return nil
}

func BuildSite(cfg *config.GlobalConfig) error {
	sg, err := NewSiteGenerator(cfg)
	if err != nil {
		return fmt.Errorf("site generator initialization failed: %w", err)
	}
	defer sg.Close()

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

	if cfg.CompileSass {
		if err := sg.compileSass(); err != nil {
			return fmt.Errorf("sass compilation failed: %w", err)
		}
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
	outputPath := filepath.Join(mergedConfig.OutputDir, page.Slug, "index.html")

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
	data := struct {
		Config *config.MergedConfig
		Page   *parser.Page
	}{
		Config: cfg,
		Page:   page,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	// Minify HTML if enabled
	if cfg.MinifyHTML {
		var minBuf bytes.Buffer
		if err := sg.minifier.Minify("text/html", &minBuf, &buf); err != nil {
			return fmt.Errorf("minification failed: %w", err)
		}
		buf = minBuf
	}

	// Write output file
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("file write failed: %w", err)
	}

	return nil
}

func (sg *SiteGenerator) compileSass() error {
	return filepath.Walk(filepath.Join("static", "scss"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("sass file access error: %w", err)
		}

		if info.IsDir() || filepath.Ext(path) != ".scss" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read SCSS file: %w", err)
		}

		result, err := sg.sass.Execute(godartsass.Args{
			Source:          string(content),
			SourceSyntax:    godartsass.SourceSyntaxSCSS,
			OutputStyle:     godartsass.OutputStyleCompressed,
			IncludePaths:    []string{filepath.Dir(path)},
			EnableSourceMap: false,
		})
		if err != nil {
			return fmt.Errorf("sass compilation error: %w", err)
		}

		outputPath := filepath.Join(
			sg.cfg.OutputDir,
			"css",
			strings.TrimSuffix(info.Name(), ".scss")+".css",
		)

		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("CSS directory creation failed: %w", err)
		}

		if err := os.WriteFile(outputPath, []byte(result.CSS), 0644); err != nil {
			return fmt.Errorf("CSS file write failed: %w", err)
		}

		return nil
	})
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
