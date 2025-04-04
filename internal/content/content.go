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
	Section  *Section
}

type Section struct {
	Children []*Node
	Config   map[string]any
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
	chain := []string{"default.html"}
	if templateName, ok := n.Config["template"].(string); ok && templateName != "" {
		chain = append([]string{templateName}, chain...)
	}
	if n.Parent != nil {
		chain = append(n.Parent.TemplateChain(), chain...)
	}
	return chain
}

func (n *Node) Breadcrumbs() []*Node {
	var crumbs []*Node
	for current := n; current != nil; current = current.Parent {
		if current.Slug != "/" {
			crumbs = append([]*Node{current}, crumbs...)
		}
	}
	return crumbs
}

func (n *Node) AllSections() []*Node {
	var sections []*Node
	var traverse func(*Node)

	traverse = func(node *Node) {
		if node.Type == NodeTypeSection {
			sections = append(sections, node)
		}
		for _, child := range node.Children {
			traverse(child)
		}
	}

	traverse(n)
	return sections
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
