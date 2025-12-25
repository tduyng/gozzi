package parser

import (
	"testing"
	"time"

	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
)

func TestParseTaxonomies_Tags(t *testing.T) {
	cfg := &config.Site{
		Taxonomies: config.TaxonomiesConfig{
			"tags": {Enabled: true},
		},
	}
	p := NewParser(cfg)

	frontMatter := &config.FrontMatter{
		Tags: []string{"go", "programming", "tutorial"},
		Date: time.Now(),
	}

	node := &content.Node{
		Path:   "content/test.md",
		Config: map[string]any{"date": frontMatter.Date},
	}

	p.ParseTaxonomies(frontMatter, node)

	// Verify tags taxonomy was created
	tagsTax, exists := p.Taxonomies["tags"]
	if !exists {
		t.Fatal("Expected tags taxonomy to be created")
	}

	// Verify all tags were added
	expectedTags := []string{"go", "programming", "tutorial"}
	if len(tagsTax.Entries) != len(expectedTags) {
		t.Errorf("Expected %d tag entries, got %d", len(expectedTags), len(tagsTax.Entries))
	}

	// Verify each tag entry
	for _, tag := range expectedTags {
		entry, exists := tagsTax.Entries[tag]
		if !exists {
			t.Errorf("Expected tag entry for '%s'", tag)
			continue
		}

		if entry.Count != 1 {
			t.Errorf("Expected count 1 for tag '%s', got %d", tag, entry.Count)
		}

		if len(entry.Pages) != 1 {
			t.Errorf("Expected 1 page for tag '%s', got %d", tag, len(entry.Pages))
		}
	}
}

func TestParseTaxonomies_Categories(t *testing.T) {
	cfg := &config.Site{
		Taxonomies: config.TaxonomiesConfig{
			"categories": {Enabled: true},
		},
	}
	p := NewParser(cfg)

	frontMatter := &config.FrontMatter{
		Categories: []string{"Development", "Tools"},
		Date:       time.Now(),
	}

	node := &content.Node{
		Path:   "content/test.md",
		Config: map[string]any{"date": frontMatter.Date},
	}

	p.ParseTaxonomies(frontMatter, node)

	// Verify categories taxonomy was created
	categoriesTax, exists := p.Taxonomies["categories"]
	if !exists {
		t.Fatal("Expected categories taxonomy to be created")
	}

	// Verify entries
	if len(categoriesTax.Entries) != 2 {
		t.Errorf("Expected 2 category entries, got %d", len(categoriesTax.Entries))
	}

	// Check specific categories (slugified)
	devEntry, exists := categoriesTax.Entries["development"]
	if !exists {
		t.Error("Expected 'development' category entry")
	} else if devEntry.Term != "Development" {
		t.Errorf("Expected term 'Development', got '%s'", devEntry.Term)
	}
}

func TestParseTaxonomies_Series(t *testing.T) {
	cfg := &config.Site{
		Taxonomies: config.TaxonomiesConfig{
			"series": {Enabled: true},
		},
	}
	p := NewParser(cfg)

	// Create multiple pages in a series
	nodes := []*content.Node{
		{
			Path: "content/part1.md",
			Config: map[string]any{
				"date":         time.Now(),
				"series_order": 1,
			},
		},
		{
			Path: "content/part3.md",
			Config: map[string]any{
				"date":         time.Now().Add(2 * time.Hour),
				"series_order": 3,
			},
		},
		{
			Path: "content/part2.md",
			Config: map[string]any{
				"date":         time.Now().Add(1 * time.Hour),
				"series_order": 2,
			},
		},
	}

	frontMatters := []*config.FrontMatter{
		{Series: "Git Mastery", SeriesOrder: 1},
		{Series: "Git Mastery", SeriesOrder: 3},
		{Series: "Git Mastery", SeriesOrder: 2},
	}

	for i, fm := range frontMatters {
		p.ParseTaxonomies(fm, nodes[i])
	}

	// Verify series taxonomy
	seriesTax, exists := p.Taxonomies["series"]
	if !exists {
		t.Fatal("Expected series taxonomy to be created")
	}

	// Verify series entry
	gitSeries, exists := seriesTax.Entries["git-mastery"]
	if !exists {
		t.Fatal("Expected 'git-mastery' series entry")
	}

	if gitSeries.Count != 3 {
		t.Errorf("Expected 3 pages in series, got %d", gitSeries.Count)
	}

	// Get ordered pages
	seriesPages := gitSeries.GetSeriesPages()
	if len(seriesPages) != 3 {
		t.Errorf("Expected 3 series pages, got %d", len(seriesPages))
	}

	// Verify ordering by series_order
	expectedOrders := []int{1, 2, 3}
	for i, sp := range seriesPages {
		if sp.Order != expectedOrders[i] {
			t.Errorf("Expected order %d at position %d, got %d", expectedOrders[i], i, sp.Order)
		}
		if sp.Position != i+1 {
			t.Errorf("Expected position %d, got %d", i+1, sp.Position)
		}
	}
}

func TestParseTaxonomies_Disabled(t *testing.T) {
	cfg := &config.Site{
		Taxonomies: config.TaxonomiesConfig{
			"tags":       {Enabled: false},
			"categories": {Enabled: false},
		},
	}
	p := NewParser(cfg)

	frontMatter := &config.FrontMatter{
		Tags:       []string{"go", "rust"},
		Categories: []string{"Development"},
		Date:       time.Now(),
	}

	node := &content.Node{
		Path:   "content/test.md",
		Config: map[string]any{"date": frontMatter.Date},
	}

	p.ParseTaxonomies(frontMatter, node)

	// Verify no taxonomies were created
	if len(p.Taxonomies) != 0 {
		t.Errorf("Expected no taxonomies when disabled, got %d", len(p.Taxonomies))
	}
}

func TestParseTaxonomies_DuplicateTags(t *testing.T) {
	cfg := &config.Site{
		Taxonomies: config.TaxonomiesConfig{
			"tags": {Enabled: true},
		},
	}
	p := NewParser(cfg)

	frontMatter := &config.FrontMatter{
		Tags: []string{"go", "Go", "GO", "programming"},
		Date: time.Now(),
	}

	node := &content.Node{
		Path:   "content/test.md",
		Config: map[string]any{"date": frontMatter.Date},
	}

	p.ParseTaxonomies(frontMatter, node)

	tagsTax := p.Taxonomies["tags"]

	// Should only have 2 unique tags (go and programming)
	if len(tagsTax.Entries) != 2 {
		t.Errorf("Expected 2 unique tags, got %d", len(tagsTax.Entries))
	}

	// Verify "go" entry exists and has only 1 page
	goEntry, exists := tagsTax.Entries["go"]
	if !exists {
		t.Fatal("Expected 'go' tag entry")
	}

	if goEntry.Count != 1 {
		t.Errorf("Expected count 1 for 'go' tag, got %d", goEntry.Count)
	}
}

func TestParseTaxonomies_MultiplePages(t *testing.T) {
	cfg := &config.Site{
		Taxonomies: config.TaxonomiesConfig{
			"tags": {Enabled: true},
		},
	}
	p := NewParser(cfg)

	// Add multiple pages with overlapping tags
	pages := []struct {
		path string
		tags []string
	}{
		{"content/post1.md", []string{"go", "tutorial"}},
		{"content/post2.md", []string{"go", "advanced"}},
		{"content/post3.md", []string{"rust", "tutorial"}},
	}

	for _, page := range pages {
		fm := &config.FrontMatter{
			Tags: page.tags,
			Date: time.Now(),
		}
		node := &content.Node{
			Path:   page.path,
			Config: map[string]any{"date": fm.Date},
		}
		p.ParseTaxonomies(fm, node)
	}

	tagsTax := p.Taxonomies["tags"]

	// Verify counts
	goEntry := tagsTax.Entries["go"]
	if goEntry.Count != 2 {
		t.Errorf("Expected 'go' tag to have 2 pages, got %d", goEntry.Count)
	}

	tutorialEntry := tagsTax.Entries["tutorial"]
	if tutorialEntry.Count != 2 {
		t.Errorf("Expected 'tutorial' tag to have 2 pages, got %d", tutorialEntry.Count)
	}

	advancedEntry := tagsTax.Entries["advanced"]
	if advancedEntry.Count != 1 {
		t.Errorf("Expected 'advanced' tag to have 1 page, got %d", advancedEntry.Count)
	}
}

func TestUrlize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Git Mastery", "git-mastery"},
		{"UPPER CASE", "upper-case"},
		{"snake_case", "snake-case"},
		{"Mixed_Case Words", "mixed-case-words"},
		{"Special@#$Chars!", "specialchars"},
		{"multiple---hyphens", "multiple-hyphens"},
		{"-leading-trailing-", "leading-trailing"},
		{"123 Numbers", "123-numbers"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := urlize(tt.input)
			if result != tt.expected {
				t.Errorf("urlize(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
