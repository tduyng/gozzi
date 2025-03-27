package parser

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/tduyng/gozzi/internal/config"
	"github.com/yuin/goldmark"
	highlight "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type Page struct {
	FrontMatter config.PageConfig
	Content     template.HTML
	Slug        string
	FilePath    string
	ModTime     time.Time
	IsNotFound  bool
}

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Footnote,
		highlight.NewHighlighting(
			highlight.WithGuessLanguage(true),
		),
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
)

var (
	datePrefixRe  = regexp.MustCompile(`^\d{4}[-_]\d{1,2}[-_]\d{1,2}[-_]`)
	slugCleanerRe = regexp.MustCompile(`[^a-z0-9\-]`)
	multiDashRe   = regexp.MustCompile(`\-+`)
)

func ParseMarkdown(path string) (*Page, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	pageConfig, err := config.LoadPageConfig(content)
	if err != nil {
		return nil, fmt.Errorf("front matter error: %w", err)
	}

	// split content into front matter and markdown
	parts := bytes.SplitN(content, []byte("+++"), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid front matter format in %s", path)
	}

	var buf bytes.Buffer
	if err := md.Convert(parts[2], &buf); err != nil {
		return nil, fmt.Errorf("failed to convert markdown: %w", err)
	}
	slug := generateSlug(path)

	return &Page{
		FrontMatter: *pageConfig,
		Content:     template.HTML(buf.String()),
		Slug:        slug,
		FilePath:    path,
		ModTime:     info.ModTime(),
	}, nil
}

func generateSlug(path string) string {
	base := extractBaseName(path)
	base = datePrefixRe.ReplaceAllString(base, "")

	slug := strings.ToLower(base)
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = slugCleanerRe.ReplaceAllString(slug, "-")
	slug = multiDashRe.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	if slug == "" {
		return "untitled"
	}
	return slug
}

func extractBaseName(path string) string {
	dir, file := filepath.Split(path)
	base := strings.TrimSuffix(file, filepath.Ext(file))

	if base == "index" {
		parentDir := strings.TrimRight(dir, string(filepath.Separator))
		if parentDir != "" {
			return filepath.Base(parentDir)
		}
	}
	return base
}
