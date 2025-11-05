// Package parser provides tag parsing and tracking functionality for content pages.
// Maintains tag index and associates pages with their tags.
package parser

import (
	"strings"

	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
)

// TagEntry tracks pages associated with a specific tag.
type TagEntry struct {
	Pages []*content.Node
	Count int
	Seen  map[string]struct{} // Track page paths
}

func (p *ContentParser) parseTags(pageConfig *config.FrontMatter, pageNode *content.Node) {
	uniqueTags := make(map[string]bool)
	for _, rawTag := range pageConfig.Tags {
		tag := strings.ToLower(strings.TrimSpace(rawTag))
		if tag == "" {
			continue
		}

		if _, exists := uniqueTags[tag]; exists {
			continue
		}
		uniqueTags[tag] = true

		// Get or create tag entry
		entry, exists := p.Tags[tag]
		if !exists {
			entry = &TagEntry{
				Seen: make(map[string]struct{}),
			}
			p.Tags[tag] = entry
		}

		// Add page if not already present
		if _, exists := entry.Seen[pageNode.Path]; !exists {
			entry.Pages = append(entry.Pages, pageNode)
			entry.Seen[pageNode.Path] = struct{}{}
			entry.Count = len(entry.Pages)
		}
	}
}
