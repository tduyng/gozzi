// ChangeTracker tracks which content changed for incremental builds
package builder

import (
	"path/filepath"
	"strings"

	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
)

// ChangeTracker analyzes changed files and determines what needs regeneration
type ChangeTracker struct {
	changedPaths map[string]bool          // Content paths that changed
	changedNodes map[*content.Node]bool   // Nodes that need regeneration
	affectedTags map[string]bool          // Tags that need regeneration
	needsSitemap bool                     // Whether sitemap needs regeneration
	needsFeed    bool                     // Whether Atom feed needs regeneration
	needsRobots  bool                     // Whether robots.txt needs regeneration
	needs404     bool                     // Whether 404 page needs regeneration
	contentMap   map[string]*content.Node // Reference to content map
	parser       *parser.ContentParser    // Reference to parser for tag lookup
}

// NewChangeTracker creates a new change tracker
func NewChangeTracker(contentMap map[string]*content.Node, p *parser.ContentParser) *ChangeTracker {
	return &ChangeTracker{
		changedPaths: make(map[string]bool),
		changedNodes: make(map[*content.Node]bool),
		affectedTags: make(map[string]bool),
		contentMap:   contentMap,
		parser:       p,
		needsSitemap: false,
		needsFeed:    false,
		needsRobots:  false,
		needs404:     false,
	}
}

// AnalyzeChanges processes the list of changed files and determines impact
func (ct *ChangeTracker) AnalyzeChanges(changedFiles []string, contentDir string) {
	for _, file := range changedFiles {
		ct.analyzeFile(file, contentDir)
	}
}

// analyzeFile determines what needs regeneration for a single changed file
func (ct *ChangeTracker) analyzeFile(file, contentDir string) {
	// Normalize path and resolve symlinks (important on macOS with /var -> /private/var)
	absContent, _ := filepath.EvalSymlinks(contentDir)
	if absContent == "" {
		absContent, _ = filepath.Abs(contentDir)
	}

	absFile, _ := filepath.EvalSymlinks(file)
	if absFile == "" {
		absFile, _ = filepath.Abs(file)
	}

	// Skip if not in content directory
	if !strings.HasPrefix(absFile, absContent) {
		return
	}

	// Get relative path from content directory
	relPath, err := filepath.Rel(absContent, absFile)
	if err != nil {
		return
	}

	// Normalize path separators
	relPath = filepath.ToSlash(relPath)

	// Mark this path as changed
	ct.changedPaths[relPath] = true

	// Determine node impact
	ct.determineNodeImpact(relPath)

	// Sitemap always needs regeneration if any content changed
	ct.needsSitemap = true

	// Feed needs regeneration if blog posts changed
	if ct.isBlogPost(relPath) {
		ct.needsFeed = true
	}
}

// determineNodeImpact finds which nodes need regeneration for a changed file
func (ct *ChangeTracker) determineNodeImpact(relPath string) {
	// Try to find the node for this path
	dir := filepath.Dir(relPath)
	if dir == "." {
		dir = ""
	}

	// Check if this is an _index.md (section page)
	if filepath.Base(relPath) == "_index.md" {
		// Section page changed - regenerate the section
		if node, exists := ct.contentMap[dir]; exists {
			ct.changedNodes[node] = true
		}
		return
	}

	// Regular page - find it in the parent section's children
	parentDir := dir
	if parentSection, exists := ct.contentMap[parentDir]; exists {
		// Find the specific page node
		for _, child := range parentSection.Children {
			if child.Type == content.NodeTypePage {
				// Match by path - check if paths overlap
				if strings.Contains(child.Path, relPath) || strings.Contains(relPath, child.Path) {
					ct.changedNodes[child] = true

					// Track affected tags
					ct.trackAffectedTags(child)
					break
				}
			}
		}
	}
}

// trackAffectedTags identifies which tag pages need regeneration
func (ct *ChangeTracker) trackAffectedTags(node *content.Node) {
	// Get tags from node config
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

// isBlogPost checks if a path represents a blog post (not a section index)
func (ct *ChangeTracker) isBlogPost(relPath string) bool {
	// Exclude section indexes
	if filepath.Base(relPath) == "_index.md" {
		return false
	}

	// Common blog post patterns
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

// GetChangedNodesCount returns the number of nodes that need regeneration
func (ct *ChangeTracker) GetChangedNodesCount() int {
	return len(ct.changedNodes)
}

// GetAffectedTagsCount returns the number of tag pages that need regeneration
func (ct *ChangeTracker) GetAffectedTagsCount() int {
	return len(ct.affectedTags)
}

// GetChangedNodes returns a slice of all nodes that need regeneration
func (ct *ChangeTracker) GetChangedNodes() []*content.Node {
	nodes := make([]*content.Node, 0, len(ct.changedNodes))
	for node := range ct.changedNodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetAffectedTags returns a slice of all tags that need regeneration
func (ct *ChangeTracker) GetAffectedTags() []string {
	tags := make([]string, 0, len(ct.affectedTags))
	for tag := range ct.affectedTags {
		tags = append(tags, tag)
	}
	return tags
}
