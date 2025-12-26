// ABOUTME: Generic taxonomy system for tags, categories, series, and custom taxonomies.
// ABOUTME: Provides unified parsing, tracking, and organization of content classifications.
package parser

import (
	"sort"
	"strings"

	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
)

// Taxonomy represents a classification system (tags, categories, series, etc.)
type Taxonomy struct {
	Name    string                    // Taxonomy name (e.g., "tags", "series")
	Entries map[string]*TaxonomyEntry // Slug -> entry mapping
	Config  config.TaxonomyConfig     // Configuration
}

// TaxonomyEntry tracks pages for a specific taxonomy term.
type TaxonomyEntry struct {
	Term     string              // Original term (e.g., "Git Mastery")
	Slug     string              // URL-friendly slug
	Pages    []*content.Node     // Associated pages
	Count    int                 // Number of pages
	Metadata map[string]any      // Extra data (e.g., series order)
	Seen     map[string]struct{} // Deduplication tracking
}

// SeriesPage wraps a page with its position in a series.
type SeriesPage struct {
	Node     *content.Node
	Order    int // series_order from frontmatter
	Position int // Calculated position (1-based)
}

// NewTaxonomy creates a new taxonomy with the given configuration.
func NewTaxonomy(name string, cfg config.TaxonomyConfig) *Taxonomy {
	return &Taxonomy{
		Name:    name,
		Entries: make(map[string]*TaxonomyEntry),
		Config:  cfg,
	}
}

// NewTaxonomyEntry creates a new taxonomy entry for a term.
func NewTaxonomyEntry(term, slug string) *TaxonomyEntry {
	return &TaxonomyEntry{
		Term:     term,
		Slug:     slug,
		Pages:    make([]*content.Node, 0),
		Metadata: make(map[string]any),
		Seen:     make(map[string]struct{}),
	}
}

// AddPage adds a page to this taxonomy entry.
func (te *TaxonomyEntry) AddPage(node *content.Node, metadata map[string]any) {
	// Check if page already exists
	if _, exists := te.Seen[node.Path]; exists {
		// Update existing page reference
		for i, page := range te.Pages {
			if page.Path == node.Path {
				te.Pages[i] = node
				// Update metadata if provided
				for k, v := range metadata {
					te.Metadata[k] = v
				}
				break
			}
		}
		return
	}

	// New page
	te.Pages = append(te.Pages, node)
	te.Seen[node.Path] = struct{}{}
	te.Count = len(te.Pages)

	// Store metadata
	for k, v := range metadata {
		te.Metadata[k] = v
	}
}

// GetOrCreateEntry gets or creates a taxonomy entry for a term.
func (t *Taxonomy) GetOrCreateEntry(term, slug string) *TaxonomyEntry {
	if entry, exists := t.Entries[slug]; exists {
		return entry
	}
	entry := NewTaxonomyEntry(term, slug)
	t.Entries[slug] = entry
	return entry
}

// GetSeriesPages returns pages ordered by series_order for a series taxonomy entry.
func (te *TaxonomyEntry) GetSeriesPages() []*SeriesPage {
	// Create series pages with order information
	seriesPages := make([]*SeriesPage, 0, len(te.Pages))

	for _, node := range te.Pages {
		order := 0
		if seriesOrder, ok := node.Config["series_order"]; ok {
			if orderInt, ok := seriesOrder.(int); ok {
				order = orderInt
			}
		}

		seriesPages = append(seriesPages, &SeriesPage{
			Node:  node,
			Order: order,
		})
	}

	// Sort by order (0 = unordered, goes last), then by date
	sort.Slice(seriesPages, func(i, j int) bool {
		// If both have order, compare orders
		if seriesPages[i].Order > 0 && seriesPages[j].Order > 0 {
			return seriesPages[i].Order < seriesPages[j].Order
		}
		// If only one has order, it comes first
		if seriesPages[i].Order > 0 {
			return true
		}
		if seriesPages[j].Order > 0 {
			return false
		}
		// Both unordered, sort by date (newest first for fallback)
		dateI := seriesPages[i].Node.Config["date"]
		dateJ := seriesPages[j].Node.Config["date"]
		return dateI.(interface{ Before(interface{}) bool }).Before(dateJ)
	})

	// Assign positions
	for i := range seriesPages {
		seriesPages[i].Position = i + 1
	}

	return seriesPages
}

// ParseTaxonomies extracts and indexes all taxonomies from frontmatter.
func (p *ContentParser) ParseTaxonomies(pageConfig *config.FrontMatter, pageNode *content.Node) {
	if p.Taxonomies == nil {
		p.Taxonomies = make(map[string]*Taxonomy)
	}

	// Parse tags
	if p.Site.Taxonomies.IsEnabled("tags") {
		p.parseTaxonomy("tags", pageConfig.Tags, pageNode, nil)
	}

	// Parse categories
	if p.Site.Taxonomies.IsEnabled("categories") {
		p.parseTaxonomy("categories", pageConfig.Categories, pageNode, nil)
	}

	// Parse series with order metadata
	if p.Site.Taxonomies.IsEnabled("series") && pageConfig.Series != "" {
		metadata := map[string]any{
			"series_order": pageConfig.SeriesOrder,
		}
		p.parseTaxonomy("series", []string{pageConfig.Series}, pageNode, metadata)
	}

	// Parse custom taxonomies from extra
	if customs, ok := pageConfig.Extra["taxonomies"].(map[string]any); ok {
		for name, terms := range customs {
			if !p.Site.Taxonomies.IsEnabled(name) {
				continue
			}

			// Convert terms to string slice
			var termList []string
			switch v := terms.(type) {
			case []any:
				for _, term := range v {
					if str, ok := term.(string); ok {
						termList = append(termList, str)
					}
				}
			case []string:
				termList = v
			case string:
				termList = []string{v}
			}

			p.parseTaxonomy(name, termList, pageNode, nil)
		}
	}
}

// RemovePageFromAllTaxonomies removes a page from all taxonomy entries.
// This should be called before re-parsing a page's taxonomies during incremental builds.
func (p *ContentParser) RemovePageFromAllTaxonomies(pageNode *content.Node) {
	for _, taxonomy := range p.Taxonomies {
		for _, entry := range taxonomy.Entries {
			// Remove this page from the entry
			newPages := make([]*content.Node, 0, len(entry.Pages))
			removed := false
			for _, page := range entry.Pages {
				if page.Path != pageNode.Path {
					newPages = append(newPages, page)
				} else {
					removed = true
				}
			}
			if removed {
				entry.Pages = newPages
				entry.Count = len(newPages)
				delete(entry.Seen, pageNode.Path)
			}
		}
	}
}

// parseTaxonomy processes a single taxonomy type.
func (p *ContentParser) parseTaxonomy(name string, terms []string, node *content.Node, metadata map[string]any) {
	// Get or create taxonomy
	taxonomy, exists := p.Taxonomies[name]
	if !exists {
		cfg := p.Site.Taxonomies.GetTaxonomyConfig(name)
		taxonomy = NewTaxonomy(name, cfg)
		p.Taxonomies[name] = taxonomy
	}

	// Process each term
	uniqueTerms := make(map[string]bool)
	for _, term := range terms {
		term = strings.TrimSpace(term)
		if term == "" {
			continue
		}

		// Normalize for uniqueness within this page
		normalizedTerm := strings.ToLower(term)
		if uniqueTerms[normalizedTerm] {
			continue
		}
		uniqueTerms[normalizedTerm] = true

		// Create slug
		slug := urlize(term)

		// Add to taxonomy
		entry := taxonomy.GetOrCreateEntry(term, slug)
		entry.AddPage(node, metadata)
	}
}

// urlize converts a string to a URL-friendly slug.
func urlize(s string) string {
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
