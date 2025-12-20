// Package parser provides tag parsing and tracking for content pages.
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
	Seen  map[string]struct{}
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

		entry, exists := p.Tags[tag]
		if !exists {
			entry = &TagEntry{
				Seen: make(map[string]struct{}),
			}
			p.Tags[tag] = entry
		}

		if _, exists := entry.Seen[pageNode.Path]; !exists {
			// New page with this tag
			entry.Pages = append(entry.Pages, pageNode)
			entry.Seen[pageNode.Path] = struct{}{}
			entry.Count = len(entry.Pages)
		} else {
			// Page already exists - replace the old pointer with new one
			for i, page := range entry.Pages {
				if page.Path == pageNode.Path {
					entry.Pages[i] = pageNode
					break
				}
			}
		}
	}
}
