package content

import (
	"html/template"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type NodeType int

const (
	NodeTypeSection NodeType = iota
	NodeTypePage
)

type Node struct {
	Type     NodeType
	Path     string
	Slug     string
	Config   map[string]any
	Content  template.HTML
	Parent   *Node
	Children []*Node
	PageMeta *PageMeta
	Section  *Section
}

type Section struct {
	Title string
	Pages []*Node
}

type PageMeta struct {
	Date    time.Time
	Updated time.Time
	Tags    []string
	Assets  string
	Draft   bool
	ImgURL  string
}

var (
	datePrefixRe  = regexp.MustCompile(`^\d{4}[-_]\d{1,2}[-_]\d{1,2}[-_]`)
	slugCleanerRe = regexp.MustCompile(`[^a-z0-9\-]`)
	multiDashRe   = regexp.MustCompile(`\-+`)
)

func NewContentNode(path string, parent *Node) *Node {
	slug := GenerateSlug(path)
	if path == "." {
		slug = "home"
	}
	return &Node{
		Path:     path,
		Slug:     slug,
		Parent:   parent,
		Children: make([]*Node, 0),
	}
}

func GenerateSlug(path string) string {
	base := extractBaseName(path)
	if match := datePrefixRe.FindStringSubmatch(base); len(match) > 0 {
		base = strings.TrimPrefix(base, match[0])
	}

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
