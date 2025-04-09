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
	Type      NodeType
	Path      string
	Slug      string
	Permalink string
	URL       string
	Config    map[string]any
	Content   template.HTML
	Parent    *Node
	Children  []*Node
	Lower     *Node
	Higher    *Node
	WordCount int
	ReadTime  int
}

var (
	datePrefixRe  = regexp.MustCompile(`^\d{4}[-_]\d{1,2}[-_]\d{1,2}[-_]`)
	slugCleanerRe = regexp.MustCompile(`[^a-z0-9\-]`)
	multiDashRe   = regexp.MustCompile(`\-+`)
)

func NewContentNode(path string, parent *Node) *Node {
	slug := GenerateSlug(path, parent)
	return &Node{
		Path:     path,
		Slug:     slug,
		Parent:   parent,
		Children: make([]*Node, 0),
	}
}

func (node *Node) ToMap() map[string]any {
	return map[string]any{
		"Type":      node.Type,
		"Path":      node.Path,
		"Slug":      node.Slug,
		"Permalink": node.Permalink,
		"URL":       node.URL,
		"Config":    node.Config,
		"Content":   node.Content,
		"Parent":    node.Parent,
		"Children":  node.Children,
		"Higher":    node.Higher,
		"Lower":     node.Lower,
		"WordCount": node.WordCount,
		"ReadTime":  node.ReadTime,
	}
}

func GenerateSlug(path string, parent *Node) string {
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

func (n *Node) TemplateChain() []string {
	chain := []string{"default.html"}

	if n.Parent != nil {
		chain = append(n.Parent.TemplateChain(), chain...)
	}

	if templateName, ok := n.Config["template"].(string); ok && templateName != "" {
		chain = append([]string{templateName}, chain...)
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
