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
	changedPaths     map[string]bool          // Content paths that changed
	changedNodes     map[*content.Node]bool   // Nodes that need regeneration
	affectedTags     map[string]bool          // Tags that need regeneration
	needsSitemap     bool                     // Whether sitemap needs regeneration
	needsFeed        bool                     // Whether Atom feed needs regeneration
	needsRobots      bool                     // Whether robots.txt needs regeneration
	needs404         bool                     // Whether 404 page needs regeneration
	needsHome        bool                     // Whether homepage needs regeneration
	needsBlogListing bool                     // Whether blog listing page needs regeneration
	contentMap       map[string]*content.Node // Reference to content map
	parser           *parser.ContentParser    // Reference to parser for tag lookup
}

// NewChangeTracker creates a new change tracker
func NewChangeTracker(contentMap map[string]*content.Node, p *parser.ContentParser) *ChangeTracker {
	return &ChangeTracker{
		changedPaths:     make(map[string]bool),
		changedNodes:     make(map[*content.Node]bool),
		affectedTags:     make(map[string]bool),
		contentMap:       contentMap,
		parser:           p,
		needsSitemap:     false,
		needsFeed:        false,
		needsRobots:      false,
		needs404:         false,
		needsHome:        false,
		needsBlogListing: false,
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
