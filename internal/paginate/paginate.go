package paginate

import (
	"sort"
	"time"

	"github.com/tduyng/gozzi/internal/content"
)

type Paginator struct {
	sections map[string]*content.Node
}

func New(sections map[string]*content.Node) *Paginator {
	return &Paginator{
		sections: sections,
	}
}

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

	sort.SliceStable(pages, func(i, j int) bool {
		dateI := pages[i].Config["date"].(time.Time)
		dateJ := pages[j].Config["date"].(time.Time)
		return dateI.Before(dateJ)
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
