package paginate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/content"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		sections map[string]*content.Node
		want     *Paginator
	}{
		{
			name:     "creates paginator with nil sections",
			sections: nil,
			want:     &Paginator{sections: nil},
		},
		{
			name:     "creates paginator with empty sections",
			sections: make(map[string]*content.Node),
			want:     &Paginator{sections: make(map[string]*content.Node)},
		},
		{
			name: "creates paginator with sections",
			sections: map[string]*content.Node{
				"blog": createTestSection("blog"),
			},
			want: &Paginator{
				sections: map[string]*content.Node{
					"blog": createTestSection("blog"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.sections)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPaginator_BuildLinks(t *testing.T) {
	tests := []struct {
		name     string
		sections map[string]*content.Node
		validate func(t *testing.T, sections map[string]*content.Node)
	}{
		{
			name:     "handles nil sections",
			sections: nil,
			validate: func(t *testing.T, sections map[string]*content.Node) {
				// Should not panic
			},
		},
		{
			name:     "handles empty sections",
			sections: make(map[string]*content.Node),
			validate: func(t *testing.T, sections map[string]*content.Node) {
				assert.Empty(t, sections)
			},
		},
		{
			name: "processes single section",
			sections: map[string]*content.Node{
				"blog": createSectionWithPages("blog", []string{"post1", "post2", "post3"}),
			},
			validate: func(t *testing.T, sections map[string]*content.Node) {
				section := sections["blog"]
				require.NotNil(t, section)
				assert.Len(t, section.Children, 3)

				// Verify pages are sorted by date (ascending)
				assertPageOrder(t, section.Children, []string{"post1", "post2", "post3"})

				// Verify pagination links
				assertPaginationLinks(t, section.Children)
			},
		},
		{
			name: "processes multiple sections",
			sections: map[string]*content.Node{
				"blog":  createSectionWithPages("blog", []string{"blog1", "blog2"}),
				"notes": createSectionWithPages("notes", []string{"note1", "note2", "note3"}),
			},
			validate: func(t *testing.T, sections map[string]*content.Node) {
				// Verify blog section
				blogSection := sections["blog"]
				require.NotNil(t, blogSection)
				assert.Len(t, blogSection.Children, 2)
				assertPageOrder(t, blogSection.Children, []string{"blog1", "blog2"})
				assertPaginationLinks(t, blogSection.Children)

				// Verify notes section
				notesSection := sections["notes"]
				require.NotNil(t, notesSection)
				assert.Len(t, notesSection.Children, 3)
				assertPageOrder(t, notesSection.Children, []string{"note1", "note2", "note3"})
				assertPaginationLinks(t, notesSection.Children)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.sections)
			p.BuildLinks()
			tt.validate(t, tt.sections)
		})
	}
}

func TestPaginator_processSection(t *testing.T) {
	tests := []struct {
		name     string
		section  *content.Node
		validate func(t *testing.T, section *content.Node)
	}{
		{
			name:    "handles section with no children",
			section: createTestSection("empty"),
			validate: func(t *testing.T, section *content.Node) {
				assert.Empty(t, section.Children)
			},
		},
		{
			name:    "handles section with single page",
			section: createSectionWithPages("single", []string{"page1"}),
			validate: func(t *testing.T, section *content.Node) {
				assert.Len(t, section.Children, 1)
				page := section.Children[0]
				assert.Nil(t, page.Higher, "single page should have no higher link")
				assert.Nil(t, page.Lower, "single page should have no lower link")
			},
		},
		{
			name: "sorts pages by date ascending",
			section: createSectionWithCustomDates("blog", map[string]time.Time{
				"newest": time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
				"oldest": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				"middle": time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			}),
			validate: func(t *testing.T, section *content.Node) {
				assert.Len(t, section.Children, 3)
				assertPageOrder(t, section.Children, []string{"oldest", "middle", "newest"})
				assertPaginationLinks(t, section.Children)
			},
		},
		{
			name:    "handles mixed node types",
			section: createSectionWithMixedNodes("mixed"),
			validate: func(t *testing.T, section *content.Node) {
				assert.Len(t, section.Children, 4) // 2 pages + 2 other nodes

				// Pages should come first, sorted by date
				pages := getPageNodes(section.Children)
				assert.Len(t, pages, 2)
				assertPageOrder(t, pages, []string{"page1", "page2"})
				assertPaginationLinks(t, pages)

				// Other nodes should be at the end
				otherNodes := getNonPageNodes(section.Children)
				assert.Len(t, otherNodes, 2)
			},
		},
		{
			name:    "handles pages with same date",
			section: createSectionWithSameDates("same-dates", []string{"page1", "page2", "page3"}),
			validate: func(t *testing.T, section *content.Node) {
				assert.Len(t, section.Children, 3)
				// With stable sort, original order should be preserved
				assertPaginationLinks(t, section.Children)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(nil)
			p.processSection(tt.section)
			tt.validate(t, tt.section)
		})
	}
}

func TestPaginator_processSection_PanicHandling(t *testing.T) {
	tests := []struct {
		name      string
		section   *content.Node
		wantPanic bool
	}{
		{
			name:      "panics with page missing date config",
			section:   createPageWithoutDate("no-date"),
			wantPanic: true,
		},
		{
			name:      "panics with page with wrong date type",
			section:   createPageWithWrongDateType("wrong-type"),
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(nil)
			if tt.wantPanic {
				assert.Panics(t, func() {
					p.processSection(tt.section)
				})
			} else {
				assert.NotPanics(t, func() {
					p.processSection(tt.section)
				})
			}
		})
	}
}

// Helper functions

func createTestSection(title string) *content.Node {
	return &content.Node{
		Type:     content.NodeTypeSection,
		Slug:     title,
		Children: []*content.Node{},
		Config:   map[string]any{"title": title},
	}
}

func createSectionWithPages(sectionTitle string, pageNames []string) *content.Node {
	section := createTestSection(sectionTitle)

	for i, name := range pageNames {
		page := &content.Node{
			Type: content.NodeTypePage,
			Slug: name,
			Config: map[string]any{
				"title": name,
				"date":  time.Date(2024, 1, i+1, 0, 0, 0, 0, time.UTC),
			},
		}
		section.Children = append(section.Children, page)
	}

	return section
}

func createSectionWithCustomDates(sectionTitle string, pageDates map[string]time.Time) *content.Node {
	section := createTestSection(sectionTitle)

	for name, date := range pageDates {
		page := &content.Node{
			Type: content.NodeTypePage,
			Slug: name,
			Config: map[string]any{
				"title": name,
				"date":  date,
			},
		}
		section.Children = append(section.Children, page)
	}

	return section
}

func createSectionWithMixedNodes(sectionTitle string) *content.Node {
	section := createTestSection(sectionTitle)

	// Add pages
	page1 := &content.Node{
		Type: content.NodeTypePage,
		Slug: "page1",
		Config: map[string]any{
			"title": "page1",
			"date":  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	page2 := &content.Node{
		Type: content.NodeTypePage,
		Slug: "page2",
		Config: map[string]any{
			"title": "page2",
			"date":  time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	// Add non-page nodes
	subsection := &content.Node{
		Type: content.NodeTypeSection,
		Slug: "subsection",
		Config: map[string]any{
			"title": "subsection",
		},
	}
	indexNode := &content.Node{
		Type: content.NodeTypeSection, // Using Section instead of Index
		Slug: "index",
		Config: map[string]any{
			"title": "index",
		},
	}

	section.Children = []*content.Node{subsection, page2, indexNode, page1}
	return section
}

func createSectionWithSameDates(sectionTitle string, pageNames []string) *content.Node {
	section := createTestSection(sectionTitle)
	sameDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, name := range pageNames {
		page := &content.Node{
			Type: content.NodeTypePage,
			Slug: name,
			Config: map[string]any{
				"title": name,
				"date":  sameDate,
			},
		}
		section.Children = append(section.Children, page)
	}

	return section
}

func createPageWithoutDate(title string) *content.Node {
	section := createTestSection("test")
	page1 := &content.Node{
		Type:   content.NodeTypePage,
		Slug:   title,
		Config: map[string]any{"title": title}, // Missing date field
	}
	page2 := &content.Node{
		Type: content.NodeTypePage,
		Slug: "other-page",
		Config: map[string]any{
			"title": "other-page",
			"date":  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	section.Children = []*content.Node{page1, page2}
	return section
}

func createPageWithWrongDateType(title string) *content.Node {
	section := createTestSection("test")
	page1 := &content.Node{
		Type: content.NodeTypePage,
		Slug: title,
		Config: map[string]any{
			"title": title,
			"date":  "not-a-time", // Wrong type
		},
	}
	page2 := &content.Node{
		Type: content.NodeTypePage,
		Slug: "other-page",
		Config: map[string]any{
			"title": "other-page",
			"date":  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	section.Children = []*content.Node{page1, page2}
	return section
}

func assertPageOrder(t *testing.T, nodes []*content.Node, expectedOrder []string) {
	pages := getPageNodes(nodes)
	require.Len(t, pages, len(expectedOrder), "unexpected number of pages")

	for i, expectedTitle := range expectedOrder {
		actualTitle := pages[i].Config["title"].(string)
		assert.Equal(t, expectedTitle, actualTitle, "page at position %d has wrong title", i)
	}
}

func assertPaginationLinks(t *testing.T, nodes []*content.Node) {
	pages := getPageNodes(nodes)

	for i, page := range pages {
		if i == 0 {
			assert.Nil(t, page.Higher, "first page should have no higher link")
		} else {
			assert.Equal(t, pages[i-1], page.Higher, "page %d should link to previous page", i)
		}

		if i == len(pages)-1 {
			assert.Nil(t, page.Lower, "last page should have no lower link")
		} else {
			assert.Equal(t, pages[i+1], page.Lower, "page %d should link to next page", i)
		}
	}
}

func getPageNodes(nodes []*content.Node) []*content.Node {
	var pages []*content.Node
	for _, node := range nodes {
		if node.Type == content.NodeTypePage {
			pages = append(pages, node)
		}
	}
	return pages
}

func getNonPageNodes(nodes []*content.Node) []*content.Node {
	var nonPages []*content.Node
	for _, node := range nodes {
		if node.Type != content.NodeTypePage {
			nonPages = append(nonPages, node)
		}
	}
	return nonPages
}
