// This file implements a clean, rule-based rebuild system for incremental builds.
package builder

import (
	"path/filepath"
	"strings"

	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
)

// RebuildScope defines what needs to be rebuilt
type RebuildScope struct {
	Nodes      map[*content.Node]bool
	Taxonomies map[string]map[string]bool // taxonomy name -> term slugs
	Global     map[string]bool            // global pages: "sitemap", "feed", "robots", "404", "home"
}

// NewRebuildScope creates an empty rebuild scope
func NewRebuildScope() *RebuildScope {
	return &RebuildScope{
		Nodes:      make(map[*content.Node]bool),
		Taxonomies: make(map[string]map[string]bool),
		Global:     make(map[string]bool),
	}
}

// MarkNode marks a node for rebuild
func (rs *RebuildScope) MarkNode(node *content.Node) {
	if node != nil {
		rs.Nodes[node] = true
	}
}

// MarkTaxonomy marks a taxonomy term for rebuild
func (rs *RebuildScope) MarkTaxonomy(taxonomyName, termSlug string) {
	if rs.Taxonomies[taxonomyName] == nil {
		rs.Taxonomies[taxonomyName] = make(map[string]bool)
	}
	rs.Taxonomies[taxonomyName][termSlug] = true
}

// MarkGlobal marks a global page for rebuild (sitemap, feed, etc.)
func (rs *RebuildScope) MarkGlobal(name string) {
	rs.Global[name] = true
}

// IsEmpty returns true if nothing needs rebuilding
func (rs *RebuildScope) IsEmpty() bool {
	return len(rs.Nodes) == 0 && len(rs.Taxonomies) == 0 && len(rs.Global) == 0
}

// GetNodes returns all nodes that need rebuilding
func (rs *RebuildScope) GetNodes() []*content.Node {
	nodes := make([]*content.Node, 0, len(rs.Nodes))
	for node := range rs.Nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// NeedsGlobal checks if a global page needs rebuilding
func (rs *RebuildScope) NeedsGlobal(name string) bool {
	return rs.Global[name]
}

// GetAffectedTaxonomies returns all affected taxonomies with their term slugs
func (rs *RebuildScope) GetAffectedTaxonomies() map[string][]string {
	result := make(map[string][]string)
	for taxonomyName, terms := range rs.Taxonomies {
		slugs := make([]string, 0, len(terms))
		for slug := range terms {
			slugs = append(slugs, slug)
		}
		result[taxonomyName] = slugs
	}
	return result
}

// RebuildAnalyzer analyzes file changes and determines rebuild scope
type RebuildAnalyzer struct {
	contentMap map[string]*content.Node
	parser     *parser.ContentParser
	contentDir string
}

// NewRebuildAnalyzer creates a new analyzer
func NewRebuildAnalyzer(contentMap map[string]*content.Node, p *parser.ContentParser, contentDir string) *RebuildAnalyzer {
	return &RebuildAnalyzer{
		contentMap: contentMap,
		parser:     p,
		contentDir: contentDir,
	}
}

// AnalyzeChanges determines what needs to be rebuilt based on changed files
func (ra *RebuildAnalyzer) AnalyzeChanges(changedFiles []string, oldTaxonomies map[string]map[string]any) *RebuildScope {
	scope := NewRebuildScope()

	for _, file := range changedFiles {
		ra.analyzeFile(file, oldTaxonomies, scope)
	}

	return scope
}

// analyzeFile analyzes a single file change
func (ra *RebuildAnalyzer) analyzeFile(file string, oldTaxonomies map[string]map[string]any, scope *RebuildScope) {
	relPath, err := filepath.Rel(ra.contentDir, file)
	if err != nil {
		return
	}

	node := ra.findNode(relPath)
	if node == nil {
		return
	}

	scope.MarkNode(node)

	if node.Type == content.NodeTypeSection {
		scope.MarkGlobal("sitemap")
		if ra.sectionHasPosts(node) {
			scope.MarkGlobal("feed")
		}
	}

	if node.Type == content.NodeTypePage {
		if node.Parent != nil {
			scope.MarkNode(node.Parent)
		}
		scope.MarkGlobal("sitemap")

		if ra.shouldIncludeInFeed(node) {
			scope.MarkGlobal("feed")
		}

		ra.analyzeTaxonomyChanges(node, relPath, oldTaxonomies, scope)
	}

	if relPath == "_index.md" || node.Path == "." {
		scope.MarkGlobal("home")
	}
}

// analyzeTaxonomyChanges detects taxonomy changes and marks affected pages
func (ra *RebuildAnalyzer) analyzeTaxonomyChanges(node *content.Node, relPath string, oldTaxonomies map[string]map[string]any, scope *RebuildScope) {
	oldValues, hasOld := oldTaxonomies[relPath]
	currentValues := ra.extractTaxonomyValues(node)

	for taxonomyName, taxonomy := range ra.parser.Taxonomies {
		oldTerms := ra.getTaxonomyTerms(oldValues, taxonomyName)
		newTerms := ra.getTaxonomyTerms(currentValues, taxonomyName)

		for _, term := range oldTerms {
			if !contains(newTerms, term) {
				slug := ra.getSlugForTerm(taxonomy, term)
				if slug != "" {
					scope.MarkTaxonomy(taxonomyName, slug)
				}
			}
		}

		for _, term := range newTerms {
			if !hasOld || !contains(oldTerms, term) {
				slug := ra.getSlugForTerm(taxonomy, term)
				if slug != "" {
					scope.MarkTaxonomy(taxonomyName, slug)
				}
			}
		}
	}
}

// getSlugForTerm finds the slug for a term in a taxonomy
func (ra *RebuildAnalyzer) getSlugForTerm(taxonomy *parser.Taxonomy, term string) string {
	for slug, entry := range taxonomy.Entries {
		if entry.Term == term {
			return slug
		}
	}
	return urlize(term)
}

func urlize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	slug := result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	return strings.Trim(slug, "-")
}

// findNode finds a content node by file path
func (ra *RebuildAnalyzer) findNode(relPath string) *content.Node {
	dir := filepath.Dir(relPath)
	if dir == "." {
		dir = ""
	}

	if filepath.Base(relPath) == "_index.md" {
		// ContentMap uses "." for root, but dir is "" after normalization
		if dir == "" {
			if node := ra.contentMap["."]; node != nil {
				return node
			}
		}
		return ra.contentMap[dir]
	}

	if filepath.Base(relPath) == "index.md" {
		parentDir := filepath.Dir(dir)
		if parentDir == "." {
			parentDir = ""
		}

		parentNode := ra.contentMap[parentDir]
		if parentNode == nil {
			return nil
		}

		slug := filepath.Base(dir)
		for _, child := range parentNode.Children {
			if child.Slug == slug || child.Slug == dir {
				return child
			}
		}
		return nil
	}

	parentNode := ra.contentMap[dir]
	if parentNode == nil {
		return nil
	}

	baseSlug := strings.TrimSuffix(filepath.Base(relPath), ".md")
	fullSlug := strings.TrimSuffix(relPath, ".md")

	for _, child := range parentNode.Children {
		if child.Slug == baseSlug || child.Slug == fullSlug {
			return child
		}
	}

	return nil
}

// extractTaxonomyValues extracts all taxonomy values from a node
func (ra *RebuildAnalyzer) extractTaxonomyValues(node *content.Node) map[string]any {
	values := make(map[string]any)

	for taxonomyName := range ra.parser.Taxonomies {
		if val, exists := node.Config[taxonomyName]; exists {
			values[taxonomyName] = val
		}
	}

	return values
}

// getTaxonomyTerms gets terms for a specific taxonomy
func (ra *RebuildAnalyzer) getTaxonomyTerms(values map[string]any, taxonomyName string) []string {
	if values == nil {
		return nil
	}

	val, exists := values[taxonomyName]
	if !exists {
		return nil
	}

	// Handle both string and []string
	switch v := val.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	default:
		return nil
	}
}

// sectionHasPosts checks if a section contains posts
func (ra *RebuildAnalyzer) sectionHasPosts(section *content.Node) bool {
	for _, child := range section.Children {
		if child.Type == content.NodeTypePage {
			return true
		}
	}
	return false
}

// shouldIncludeInFeed checks if a node should be in the feed
func (ra *RebuildAnalyzer) shouldIncludeInFeed(node *content.Node) bool {
	if generateFeed, ok := node.Config["generate_feed"].(bool); ok {
		return generateFeed
	}
	return false
}

// SnapshotTaxonomyValues captures current taxonomy values for changed files.
// This should be called BEFORE ParseFiles() to preserve old values for comparison.
func (ra *RebuildAnalyzer) SnapshotTaxonomyValues(changedFiles []string) map[string]map[string]any {
	snapshot := make(map[string]map[string]any)

	for _, file := range changedFiles {
		relPath, err := filepath.Rel(ra.contentDir, file)
		if err != nil {
			continue
		}

		node := ra.findNode(relPath)
		if node != nil {
			snapshot[relPath] = ra.extractTaxonomyValues(node)
		}
	}

	return snapshot
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
