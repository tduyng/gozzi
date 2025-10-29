// Package paginate provides pagination utilities for content sections in gozzi.
package paginate

import (
	"slices"
	"time"

	"github.com/tduyng/gozzi/app/content"
)

// Paginator builds pagination links for content sections.
type Paginator struct {
	sections map[string]*content.Node
}

// New creates a new Paginator for the given sections.
func New(sections map[string]*content.Node) *Paginator {
	return &Paginator{
		sections: sections,
	}
}

// BuildLinks generates next/previous links for all paginated sections.
func (p *Paginator) BuildLinks() {
	for _, section := range p.sections {
		p.processSection(section)
	}
}

func (p *Paginator) processSection(section *content.Node) {
	var pages []*content.Node
	var otherNodes []*content.Node

	for _, child := range section.Children {
		if child.Type == content.NodeTypePage {
			pages = append(pages, child)
		} else {
			otherNodes = append(otherNodes, child)
		}
	}

	slices.SortStableFunc(pages, func(a, b *content.Node) int {
		dateA := a.Config["date"].(time.Time)
		dateB := b.Config["date"].(time.Time)
		return dateA.Compare(dateB) // Ascending order (oldest first)
	})
	section.Children = append(pages, otherNodes...)

	for i := range pages {
		if i > 0 {
			pages[i].Higher = pages[i-1]
		}
		if i < len(pages)-1 {
			pages[i].Lower = pages[i+1]
		}
	}
}
