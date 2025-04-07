package content

import (
	"html/template"
	"path/filepath"
	"regexp"
	"strings"
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
}

var (
	datePrefixRe  = regexp.MustCompile(`^\d{4}[-_]\d{1,2}[-_]\d{1,2}[-_]`)
	slugCleanerRe = regexp.MustCompile(`[^a-z0-9\-]`)
	multiDashRe   = regexp.MustCompile(`\-+`)
)

func NewContentNode(path string, parent *Node) *Node {
	slug := GenerateSlug(path)
	if path == "." {
		slug = "/"
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

func (n *Node) TemplateChain() []string {
	chain := []string{}
	if templateName, ok := n.Config["template"].(string); ok && templateName != "" {
		chain = append([]string{templateName}, chain...)
	}
	if n.Parent != nil {
		chain = append(n.Parent.TemplateChain(), chain...)
	}
	return chain
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
