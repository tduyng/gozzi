// Package content provides content tree data structures and navigation for gozzi.
package content

import (
	"html/template"
	"path/filepath"
	"regexp"
	"strings"
)

// NodeType represents the type of content node (section or page).
type NodeType int

const (
	// NodeTypeSection represents a content section (directory).
	NodeTypeSection NodeType = iota
	// NodeTypePage represents a content page (markdown file).
	NodeTypePage
)

// Node represents a content node in the hierarchical content tree.
type Node struct {
	Type      NodeType
	Path      string
	Slug      string
	Permalink string
	URL       string
	Config    map[string]any
	Content   template.HTML
	Summary   template.HTML
	Parent    *Node
	Children  []*Node
	Lower     *Node
	Higher    *Node
	WordCount int
	ReadTime  int
	Toc       []map[string]any
	Aliases   []string
}

var (
	datePrefixRe  = regexp.MustCompile(`^\d{4}[-_]\d{1,2}[-_]\d{1,2}[-_]`)
	slugCleanerRe = regexp.MustCompile(`[^a-z0-9\-]`)
	multiDashRe   = regexp.MustCompile(`\-+`)
)

// NewContentNode creates a new content node with the given path and parent.
func NewContentNode(path string, parent *Node) *Node {
	slug := GenerateSlug(path, parent)
	return &Node{
		Path:     path,
		Slug:     slug,
		Parent:   parent,
		Children: make([]*Node, 0),
	}
}

// ToMap converts the node to a map for template rendering.
func (node *Node) ToMap() map[string]any {
	return map[string]any{
		"Type":      node.Type,
		"Path":      node.Path,
		"Slug":      node.Slug,
		"Permalink": node.Permalink,
		"URL":       node.URL,
		"Config":    node.Config,
		"Content":   node.Content,
		"Summary":   node.Summary,
		"Parent":    node.Parent,
		"Children":  node.Children,
		"Higher":    node.Higher,
		"Lower":     node.Lower,
		"WordCount": node.WordCount,
		"ReadTime":  node.ReadTime,
		"Toc":       node.Toc,
	}
}

// ToMapMinimal returns a minimal map representation for cache efficiency.
func (node *Node) ToMapMinimal() map[string]any {
	minimalChildren := make([]*Node, len(node.Children))
	for i, child := range node.Children {
		minimalChildren[i] = &Node{
			Type:      child.Type,
			Path:      child.Path,
			Slug:      child.Slug,
			Permalink: child.Permalink,
			URL:       child.URL,
			Config:    child.Config,
			Summary:   child.Summary,
			WordCount: child.WordCount,
			ReadTime:  child.ReadTime,
		}
	}

	return map[string]any{
		"Type":      node.Type,
		"Path":      node.Path,
		"Slug":      node.Slug,
		"Permalink": node.Permalink,
		"URL":       node.URL,
		"Config":    node.Config,
		"Children":  minimalChildren,
		"Higher":    node.Higher,
		"Lower":     node.Lower,
		"WordCount": node.WordCount,
		"ReadTime":  node.ReadTime,
	}
}

// GenerateSlug generates a URL-friendly slug from a file path.
func GenerateSlug(path string, parent *Node) string {
	base := extractBaseName(path)
	base = datePrefixRe.ReplaceAllString(base, "")
	slug := strings.ToLower(base)
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = slugCleanerRe.ReplaceAllString(slug, "-")
	slug = multiDashRe.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	if parent != nil && parent.Slug != "" {
		return filepath.ToSlash(filepath.Join(parent.Slug, slug))
	}

	return slug
}

// TemplateChain returns the hierarchical chain of template names.
func (node *Node) TemplateChain() []string {
	chain := []string{"default.html"}

	if node.Parent != nil {
		chain = append(node.Parent.TemplateChain(), chain...)
	}

	if templateName, ok := node.Config["template"].(string); ok && templateName != "" {
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
