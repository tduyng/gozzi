// ChangeTracker tracks which content changed for incremental builds
package builder

import (
	"path/filepath"
	"strings"

	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
)

// ChangeTracker tracks which content changed for incremental builds.
type ChangeTracker struct {
	changedPaths       map[string]bool            // Content paths that changed
	changedNodes       map[*content.Node]bool     // Nodes that need regeneration
	affectedTags       map[string]bool            // Tags that need regeneration (deprecated, use affectedTaxonomies)
	affectedTaxonomies map[string]map[string]bool // taxonomyName -> {term slugs}
	needsSitemap       bool                       // Whether sitemap needs regeneration
	needsFeed          bool                       // Whether Atom feed needs regeneration
	needsRobots        bool                       // Whether robots.txt needs regeneration
	needs404           bool                       // Whether 404 page needs regeneration
	needsHome          bool                       // Whether homepage needs regeneration
	needsBlogListing   bool                       // Whether blog listing page needs regeneration
	contentMap         map[string]*content.Node   // Reference to content map
	parser             *parser.ContentParser      // Reference to parser for taxonomy lookup
}

// NewChangeTracker creates a new change tracker
func NewChangeTracker(contentMap map[string]*content.Node, p *parser.ContentParser) *ChangeTracker {
	return &ChangeTracker{
		changedPaths:       make(map[string]bool),
		changedNodes:       make(map[*content.Node]bool),
		affectedTags:       make(map[string]bool),
		affectedTaxonomies: make(map[string]map[string]bool),
		contentMap:         contentMap,
		parser:             p,
		needsSitemap:       false,
		needsFeed:          false,
		needsRobots:        false,
		needs404:           false,
		needsHome:          false,
		needsBlogListing:   false,
	}
}

// AnalyzeChanges processes the list of changed files.
func (ct *ChangeTracker) AnalyzeChanges(changedFiles []string, contentDir string) {
	for _, file := range changedFiles {
		ct.analyzeFile(file, contentDir)
	}
}

// analyzeFile determines what needs regeneration for a single changed file.
func (ct *ChangeTracker) analyzeFile(file, contentDir string) {
	absContentRaw, _ := filepath.Abs(contentDir)
	absFileRaw, _ := filepath.Abs(file)

	absContent, err := filepath.EvalSymlinks(absContentRaw)
	if err != nil {
		absContent = absContentRaw
	}

	absFile, err := filepath.EvalSymlinks(absFileRaw)
	if err != nil {
		absFile = absFileRaw
	}

	if !strings.HasPrefix(absFile, absContent) {
		return
	}

	relPath, err := filepath.Rel(absContent, absFile)
	if err != nil {
		return
	}

	relPath = filepath.ToSlash(relPath)

	ct.changedPaths[relPath] = true

	ct.determineNodeImpact(relPath)

	ct.needsSitemap = true

	if ct.isBlogPost(relPath) {
		ct.needsFeed = true
		ct.needsHome = true
		ct.needsBlogListing = true
	}
}

// determineNodeImpact finds which nodes need regeneration.
func (ct *ChangeTracker) determineNodeImpact(relPath string) {
	dir := filepath.Dir(relPath)
	if dir == "." {
		dir = ""
	}

	if filepath.Base(relPath) == "_index.md" {
		if node, exists := ct.contentMap[dir]; exists {
			ct.changedNodes[node] = true
		}
		return
	}

	parentDir := dir
	if filepath.Base(relPath) == "index.md" {
		parentDir = filepath.Dir(dir)
		if parentDir == "." {
			parentDir = ""
		}
	}

	if parentSection, exists := ct.contentMap[parentDir]; exists {
		for _, child := range parentSection.Children {
			if child.Type == content.NodeTypePage {
				if strings.Contains(child.Path, relPath) || strings.Contains(relPath, child.Path) {
					ct.changedNodes[child] = true

					ct.trackAffectedTags(child)
					ct.trackAffectedTaxonomies(child)
					break
				}
			}
		}
	}
}

// trackAffectedTags identifies which tag pages need regeneration.
func (ct *ChangeTracker) trackAffectedTags(node *content.Node) {
	if tags, ok := node.Config["tags"]; ok {
		switch tagList := tags.(type) {
		case []string:
			for _, tag := range tagList {
				ct.affectedTags[strings.ToLower(strings.TrimSpace(tag))] = true
			}
		case []interface{}:
			for _, tag := range tagList {
				if tagStr, ok := tag.(string); ok {
					ct.affectedTags[strings.ToLower(strings.TrimSpace(tagStr))] = true
				}
			}
		}
	}
}

// trackAffectedTaxonomies identifies which taxonomy terms need regeneration.
func (ct *ChangeTracker) trackAffectedTaxonomies(node *content.Node) {
	// Track tags
	if tags, ok := node.Config["tags"]; ok {
		ct.trackTaxonomyTerms("tags", tags)
	}

	// Track categories
	if categories, ok := node.Config["categories"]; ok {
		ct.trackTaxonomyTerms("categories", categories)
	}

	// Track series
	if series, ok := node.Config["series"]; ok {
		ct.trackTaxonomyTerms("series", series)
	}

	// Track custom taxonomies from extra.taxonomies
	if extra, ok := node.Config["extra"].(map[string]any); ok {
		if taxonomies, ok := extra["taxonomies"].(map[string]any); ok {
			for taxName, terms := range taxonomies {
				ct.trackTaxonomyTerms(taxName, terms)
			}
		}
	}
}

// trackTaxonomyTerms adds terms to the affected taxonomies map.
func (ct *ChangeTracker) trackTaxonomyTerms(taxonomyName string, terms any) {
	if ct.affectedTaxonomies[taxonomyName] == nil {
		ct.affectedTaxonomies[taxonomyName] = make(map[string]bool)
	}

	// Convert terms to slugs and track them
	switch termList := terms.(type) {
	case string:
		// Single term (e.g., series)
		slug := urlizeTracker(termList)
		ct.affectedTaxonomies[taxonomyName][slug] = true
	case []string:
		for _, term := range termList {
			slug := urlizeTracker(term)
			ct.affectedTaxonomies[taxonomyName][slug] = true
		}
	case []interface{}:
		for _, term := range termList {
			if termStr, ok := term.(string); ok {
				slug := urlizeTracker(termStr)
				ct.affectedTaxonomies[taxonomyName][slug] = true
			}
		}
	}
}

// urlizeTracker converts a term to a URL-friendly slug (matches parser.urlize).
func urlizeTracker(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Remove non-alphanumeric characters except hyphens
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	// Remove consecutive hyphens
	slug := result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start/end
	slug = strings.Trim(slug, "-")

	return slug
}

// isBlogPost checks if a path represents a blog post.
func (ct *ChangeTracker) isBlogPost(relPath string) bool {
	if filepath.Base(relPath) == "_index.md" {
		return false
	}

	return strings.HasPrefix(relPath, "blog/") ||
		strings.HasPrefix(relPath, "posts/") ||
		strings.HasPrefix(relPath, "articles/")
}

// ShouldRegenerateNode checks if a specific node needs regeneration
func (ct *ChangeTracker) ShouldRegenerateNode(node *content.Node) bool {
	return ct.changedNodes[node]
}

// ShouldRegenerateTag checks if a specific tag page needs regeneration
func (ct *ChangeTracker) ShouldRegenerateTag(tag string) bool {
	return ct.affectedTags[tag]
}

// ShouldRegenerateSitemap returns true if sitemap needs regeneration
func (ct *ChangeTracker) ShouldRegenerateSitemap() bool {
	return ct.needsSitemap
}

// ShouldRegenerateFeed returns true if Atom feed needs regeneration
func (ct *ChangeTracker) ShouldRegenerateFeed() bool {
	return ct.needsFeed
}

// ShouldRegenerateRobots returns true if robots.txt needs regeneration
func (ct *ChangeTracker) ShouldRegenerateRobots() bool {
	return ct.needsRobots
}

// ShouldRegenerate404 returns true if 404 page needs regeneration
func (ct *ChangeTracker) ShouldRegenerate404() bool {
	return ct.needs404
}

// ShouldRegenerateHome returns true if homepage needs regeneration
func (ct *ChangeTracker) ShouldRegenerateHome() bool {
	return ct.needsHome
}

// ShouldRegenerateBlogListing returns true if blog listing page needs regeneration
func (ct *ChangeTracker) ShouldRegenerateBlogListing() bool {
	return ct.needsBlogListing
}

// GetChangedNodesCount returns the number of nodes that need regeneration
func (ct *ChangeTracker) GetChangedNodesCount() int {
	return len(ct.changedNodes)
}

// GetAffectedTagsCount returns the number of tag pages that need regeneration.
func (ct *ChangeTracker) GetAffectedTagsCount() int {
	return len(ct.affectedTags)
}

// GetChangedNodes returns a slice of all nodes that need regeneration.
func (ct *ChangeTracker) GetChangedNodes() []*content.Node {
	nodes := make([]*content.Node, 0, len(ct.changedNodes))
	for node := range ct.changedNodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetChangedNodesAfterParse re-resolves node pointers after parsing completes.
func (ct *ChangeTracker) GetChangedNodesAfterParse(contentMap map[string]*content.Node) []*content.Node {
	nodes := make([]*content.Node, 0, len(ct.changedNodes))

	for relPath := range ct.changedPaths {
		dir := filepath.Dir(relPath)
		if dir == "." {
			dir = ""
		}

		if filepath.Base(relPath) == "_index.md" {
			if node, exists := contentMap[dir]; exists {
				nodes = append(nodes, node)
			}
			continue
		}

		parentDir := dir
		if filepath.Base(relPath) == "index.md" {
			parentDir = filepath.Dir(dir)
			if parentDir == "." {
				parentDir = ""
			}
		}

		if parentSection, exists := contentMap[parentDir]; exists {
			for _, child := range parentSection.Children {
				if child.Type == content.NodeTypePage {
					if strings.Contains(child.Path, relPath) || strings.Contains(relPath, child.Path) {
						nodes = append(nodes, child)
						break
					}
				}
			}
		}
	}

	return nodes
}

// GetAffectedTags returns a slice of all tags that need regeneration.
func (ct *ChangeTracker) GetAffectedTags() []string {
	tags := make([]string, 0, len(ct.affectedTags))
	for tag := range ct.affectedTags {
		tags = append(tags, tag)
	}
	return tags
}

// GetAffectedTaxonomies returns a map of taxonomy names to affected term slugs.
func (ct *ChangeTracker) GetAffectedTaxonomies() map[string][]string {
	result := make(map[string][]string)
	for taxonomyName, terms := range ct.affectedTaxonomies {
		termList := make([]string, 0, len(terms))
		for slug := range terms {
			termList = append(termList, slug)
		}
		if len(termList) > 0 {
			result[taxonomyName] = termList
		}
	}
	return result
}

// GetAffectedTaxonomyCount returns the total number of affected taxonomy terms across all taxonomies.
func (ct *ChangeTracker) GetAffectedTaxonomyCount() int {
	count := 0
	for _, terms := range ct.affectedTaxonomies {
		count += len(terms)
	}
	return count
}
