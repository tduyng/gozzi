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
	oldTaxonomyValues  map[string]map[string]any  // Snapshot of old taxonomy values: relPath -> {field -> value}
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
		oldTaxonomyValues:  make(map[string]map[string]any),
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
	// Snapshot will either use provided values or try to snapshot from current contentMap
	// Note: If snapshotting from contentMap, it may be too late if ParseFiles already ran
	ct.snapshotOldTaxonomyValues(changedFiles, contentDir)

	for _, file := range changedFiles {
		ct.analyzeFile(file, contentDir)
	}
}

// SetOldTaxonomyValues allows pre-snapshotted values to be set.
// This should be used when values were captured before ParseFiles() updated the contentMap.
func (ct *ChangeTracker) SetOldTaxonomyValues(values map[string]map[string]any) {
	if values != nil {
		ct.oldTaxonomyValues = values
	}
}

// snapshotOldTaxonomyValues saves the current taxonomy values.
// NOTE: This will only work correctly if called BEFORE ParseFiles updates the contentMap.
func (ct *ChangeTracker) snapshotOldTaxonomyValues(changedFiles []string, contentDir string) {
	// Skip if values were already set via SetOldTaxonomyValues
	if len(ct.oldTaxonomyValues) > 0 {
		return
	}

	for _, file := range changedFiles {
		relPath := ct.normalizeFilePath(file, contentDir)
		if relPath == "" {
			continue
		}

		// Find the node for this file and snapshot its taxonomy values
		node := ct.findNodeByPath(relPath)
		if node != nil {
			ct.oldTaxonomyValues[relPath] = ct.extractTaxonomyValues(node)
		} else {
		}
	}
}

// findNodeByPath finds a node by its relative path.
func (ct *ChangeTracker) findNodeByPath(relPath string) *content.Node {
	dir := filepath.Dir(relPath)
	if dir == "." {
		dir = ""
	}

	if filepath.Base(relPath) == "_index.md" {
		if node, exists := ct.contentMap[dir]; exists {
			return node
		}
		return nil
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
					return child
				}
			}
		}
	}

	return nil
}

// normalizeFilePath converts an absolute file path to a relative path.
func (ct *ChangeTracker) normalizeFilePath(file, contentDir string) string {
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
		return ""
	}

	relPath, err := filepath.Rel(absContent, absFile)
	if err != nil {
		return ""
	}

	return filepath.ToSlash(relPath)
}

// extractTaxonomyValues extracts all taxonomy values from a node.
func (ct *ChangeTracker) extractTaxonomyValues(node *content.Node) map[string]any {
	values := make(map[string]any)

	// Extract standard taxonomy fields
	taxonomyFields := []string{"tags", "categories", "series"}
	for _, field := range taxonomyFields {
		if val, ok := node.Config[field]; ok {
			values[field] = cloneTaxonomyValue(val)
		}
	}

	// Extract custom taxonomies
	if extra, ok := node.Config["extra"].(map[string]any); ok {
		if taxonomies, ok := extra["taxonomies"].(map[string]any); ok {
			for taxName, val := range taxonomies {
				values["extra.taxonomies."+taxName] = cloneTaxonomyValue(val)
			}
		}
	}

	return values
}

// cloneTaxonomyValue creates a deep copy of a taxonomy value.
func cloneTaxonomyValue(val any) any {
	switch v := val.(type) {
	case string:
		return v
	case []string:
		clone := make([]string, len(v))
		copy(clone, v)
		return clone
	case []interface{}:
		clone := make([]interface{}, len(v))
		copy(clone, v)
		return clone
	default:
		return val
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
		if termList == "" {
			return
		}
		slug := urlizeTracker(termList)
		if slug != "" {
			ct.affectedTaxonomies[taxonomyName][slug] = true
		}
	case []string:
		for _, term := range termList {
			if term == "" {
				continue
			}
			slug := urlizeTracker(term)
			if slug != "" {
				ct.affectedTaxonomies[taxonomyName][slug] = true
			}
		}
	case []interface{}:
		for _, term := range termList {
			if termStr, ok := term.(string); ok && termStr != "" {
				slug := urlizeTracker(termStr)
				if slug != "" {
					ct.affectedTaxonomies[taxonomyName][slug] = true
				}
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

// markSeriesPostsForRegeneration marks all posts in a series for regeneration.
// This ensures that when a post is added/removed/reordered in a series,
// all posts get regenerated with updated prev/next navigation.
func (ct *ChangeTracker) markSeriesPostsForRegeneration(seriesValue any) {
	seriesName, ok := seriesValue.(string)
	if !ok || seriesName == "" {
		return
	}

	// Get the series taxonomy
	if ct.parser == nil || ct.parser.Taxonomies == nil {
		return
	}

	seriesTaxonomy, exists := ct.parser.Taxonomies["series"]
	if !exists {
		return
	}

	// Find the series entry
	slug := urlizeTracker(seriesName)
	entry, exists := seriesTaxonomy.Entries[slug]
	if !exists {
		return
	}

	// Mark all posts in this series for regeneration
	for _, node := range entry.Pages {
		ct.changedNodes[node] = true
	}
}

// compareTaxonomies compares old and new taxonomy values and tracks both.
// This ensures that when a post moves from one taxonomy term to another,
// both the old and new term pages are regenerated.
func (ct *ChangeTracker) compareTaxonomies(relPath string, newNode *content.Node) {
	// Track new taxonomy values (pages we're joining)
	ct.trackAffectedTaxonomies(newNode)

	// Get old taxonomy values from the snapshot
	oldValues, hasOldValues := ct.oldTaxonomyValues[relPath]

	if !hasOldValues {
		// No old values to compare, this might be a new file
		return
	}

	// Track old taxonomy values (pages we're leaving)
	// We need to regenerate these to remove the post
	taxonomyFields := []string{"tags", "categories", "series"}

	for _, field := range taxonomyFields {
		oldValue, oldExists := oldValues[field]
		newValue, newExists := newNode.Config[field]

		// Special handling for series: mark all posts in affected series for regeneration
		// This ensures prev/next navigation is updated for all posts in the series
		if field == "series" {
			// Series was added
			if !oldExists && newExists {
				ct.markSeriesPostsForRegeneration(newValue)
			}

			// Series was changed
			if oldExists && newExists && !equalTaxonomyValue(oldValue, newValue) {
				ct.markSeriesPostsForRegeneration(oldValue)
				ct.markSeriesPostsForRegeneration(newValue)
			}

			// Series was removed
			if oldExists && !newExists {
				ct.markSeriesPostsForRegeneration(oldValue)
			}

			// If series changed or removed, track old taxonomy terms
			if oldExists && (!newExists || !equalTaxonomyValue(oldValue, newValue)) {
				ct.trackTaxonomyTerms(field, oldValue)
			}
		} else {
			// For non-series taxonomies, use the original logic
			if oldExists && (!newExists || !equalTaxonomyValue(oldValue, newValue)) {
				ct.trackTaxonomyTerms(field, oldValue)
			}
		}
	}

	// Also check custom taxonomies in extra.taxonomies
	for key, oldValue := range oldValues {
		if strings.HasPrefix(key, "extra.taxonomies.") {
			taxName := strings.TrimPrefix(key, "extra.taxonomies.")

			// Check if this taxonomy still exists in new node with same value
			newHasThis := false
			if newExtra, ok := newNode.Config["extra"].(map[string]any); ok {
				if newTaxonomies, ok := newExtra["taxonomies"].(map[string]any); ok {
					if newValue, exists := newTaxonomies[taxName]; exists {
						if !equalTaxonomyValue(oldValue, newValue) {
							// Value changed, track old terms
							ct.trackTaxonomyTerms(taxName, oldValue)
						}
						newHasThis = true
					}
				}
			}

			// Taxonomy was removed, track old terms
			if !newHasThis {
				ct.trackTaxonomyTerms(taxName, oldValue)
			}
		}
	}
}

// equalTaxonomyValue checks if two taxonomy values are equal.
func equalTaxonomyValue(a, b any) bool {
	// Handle string values (series)
	if strA, okA := a.(string); okA {
		if strB, okB := b.(string); okB {
			return strA == strB
		}
		return false
	}

	// Handle []string values
	if sliceA, okA := a.([]string); okA {
		if sliceB, okB := b.([]string); okB {
			if len(sliceA) != len(sliceB) {
				return false
			}
			// Create maps for comparison
			mapA := make(map[string]bool)
			for _, v := range sliceA {
				mapA[v] = true
			}
			for _, v := range sliceB {
				if !mapA[v] {
					return false
				}
			}
			return true
		}
		return false
	}

	// Handle []interface{} values
	if sliceA, okA := a.([]interface{}); okA {
		if sliceB, okB := b.([]interface{}); okB {
			if len(sliceA) != len(sliceB) {
				return false
			}
			mapA := make(map[string]bool)
			for _, v := range sliceA {
				if strV, ok := v.(string); ok {
					mapA[strV] = true
				}
			}
			for _, v := range sliceB {
				if strV, ok := v.(string); ok {
					if !mapA[strV] {
						return false
					}
				}
			}
			return true
		}
		return false
	}

	return false
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
	// Use a map to avoid duplicates
	nodeSet := make(map[*content.Node]bool)

	for relPath := range ct.changedPaths {
		dir := filepath.Dir(relPath)
		if dir == "." {
			dir = ""
		}

		if filepath.Base(relPath) == "_index.md" {
			if node, exists := contentMap[dir]; exists {
				nodeSet[node] = true

				// Compare old vs new taxonomy values
				ct.compareTaxonomies(relPath, node)
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
						nodeSet[child] = true

						// Compare old vs new taxonomy values for this child
						ct.compareTaxonomies(relPath, child)
						break
					}
				}
			}
		}
	}

	// After comparing taxonomies, include all nodes marked for regeneration
	// This includes nodes added by markSeriesPostsForRegeneration
	for node := range ct.changedNodes {
		nodeSet[node] = true
	}

	// Convert set to slice
	nodes := make([]*content.Node, 0, len(nodeSet))
	for node := range nodeSet {
		nodes = append(nodes, node)
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
