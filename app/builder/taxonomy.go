// Generic taxonomy page generation for tags, categories, series, and custom taxonomies.
// Builds index pages and term pages with support for series ordering and navigation.
package builder

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// generateTaxonomyPages generates all taxonomy pages (tags, categories, series, custom).
func (b *Builder) generateTaxonomyPages() error {
	for taxonomyName, taxonomy := range b.parser.Taxonomies {
		// Check if templates exist for this taxonomy
		indexTemplate := taxonomyName + ".html"             // e.g., "tags.html", "series.html"
		termTemplate := singularize(taxonomyName) + ".html" // e.g., "tag.html", "serie.html"

		indexExists := b.hasTemplate(indexTemplate)
		termExists := b.hasTemplate(termTemplate)

		if !indexExists && !termExists {
			continue
		}

		// Create taxonomy directory
		taxonomyDir := filepath.Join(b.site.OutputDir, taxonomyName)
		if err := os.MkdirAll(taxonomyDir, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", taxonomyName, err)
		}

		// Generate index page (e.g., /tags/)
		if indexExists {
			if err := b.generateTaxonomyIndex(taxonomyName, taxonomy, indexTemplate); err != nil {
				return fmt.Errorf("failed to generate %s index: %w", taxonomyName, err)
			}
		}

		// Generate term pages (e.g., /tags/go/, /series/git-mastery/)
		if termExists {
			for slug, entry := range taxonomy.Entries {
				if err := b.generateTaxonomyTerm(taxonomyName, slug, entry, termTemplate); err != nil {
					return fmt.Errorf("failed to generate %s term '%s': %w", taxonomyName, slug, err)
				}
			}
		}
	}

	return nil
}

// generateTaxonomyIndex generates the index page for a taxonomy (e.g., /tags/, /series/).
func (b *Builder) generateTaxonomyIndex(taxonomyName string, taxonomy *parser.Taxonomy, templateName string) error {
	// Build list of terms with metadata
	terms := make([]map[string]any, 0, len(taxonomy.Entries))
	for slug, entry := range taxonomy.Entries {
		permalink := b.buildTaxonomyPermalink(taxonomyName, slug)
		term := map[string]any{
			"Name":      entry.Term,
			"Slug":      slug,
			"Count":     entry.Count,
			"Permalink": permalink,
			"URL":       b.buildTaxonomyURL(permalink),
		}
		terms = append(terms, term)
	}

	// Sort alphabetically by term name
	slices.SortFunc(terms, func(a, b map[string]any) int {
		return strings.Compare(a["Name"].(string), b["Name"].(string))
	})

	// Build template data
	caser := cases.Title(language.English)
	data := map[string]any{
		"Site": map[string]any{
			"Config":     b.site.ToConfig(),
			"Taxonomies": b.buildTaxonomiesMap(),
		},
		"Page": map[string]any{
			"Title":    caser.String(taxonomyName),
			"Terms":    terms,
			"Path":     "/" + taxonomyName,
			"Taxonomy": taxonomyName,
		},
	}

	outputPath := filepath.Join(b.site.OutputDir, taxonomyName, "index.html")
	return b.renderTemplate(nil, outputPath, data, templateName)
}

// generateTaxonomyTerm generates a single term page (e.g., /tags/go/, /series/git-mastery/).
func (b *Builder) generateTaxonomyTerm(
	taxonomyName, slug string, entry *parser.TaxonomyEntry, templateName string,
) error {
	var pageList []map[string]any

	// Special handling for series - use ordered pages with prev/next navigation
	if taxonomyName == "series" {
		pageList = b.buildSeriesPageList(entry)
	} else {
		// For tags/categories - sort by date (newest first)
		pageList = b.buildStandardPageList(entry)
	}

	// Build template data
	caser := cases.Title(language.English)
	data := map[string]any{
		"Site": map[string]any{
			"Config":     b.site.ToConfig(),
			"Taxonomies": b.buildTaxonomiesMap(),
		},
		"Page": map[string]any{
			"Title":     fmt.Sprintf("%s: %s", caser.String(singularize(taxonomyName)), entry.Term),
			"Term":      entry.Term,
			"Slug":      slug,
			"Pages":     pageList,
			"Count":     entry.Count,
			"Permalink": b.buildTaxonomyPermalink(taxonomyName, slug),
			"Path":      b.buildTaxonomyPermalink(taxonomyName, slug),
			"Taxonomy":  taxonomyName,
		},
	}

	outputPath := filepath.Join(
		b.site.OutputDir,
		taxonomyName,
		slug,
		"index.html",
	)

	// Ensure term directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create term directory: %w", err)
	}

	return b.renderTemplate(nil, outputPath, data, templateName)
}

// buildSeriesPageList builds an ordered list of pages for a series with prev/next navigation.
func (b *Builder) buildSeriesPageList(entry *parser.TaxonomyEntry) []map[string]any {
	seriesPages := entry.GetSeriesPages()
	pageList := make([]map[string]any, len(seriesPages))

	for i, sp := range seriesPages {
		pageData := map[string]any{
			"Permalink": sp.Node.Permalink,
			"Config":    sp.Node.Config,
			"Position":  sp.Position,
			"Order":     sp.Order,
		}

		// Add prev/next navigation
		if i > 0 {
			prevPage := seriesPages[i-1]
			pageData["Prev"] = map[string]any{
				"Permalink": prevPage.Node.Permalink,
				"Config":    prevPage.Node.Config,
				"Position":  prevPage.Position,
			}
		}

		if i < len(seriesPages)-1 {
			nextPage := seriesPages[i+1]
			pageData["Next"] = map[string]any{
				"Permalink": nextPage.Node.Permalink,
				"Config":    nextPage.Node.Config,
				"Position":  nextPage.Position,
			}
		}

		pageList[i] = pageData
	}

	return pageList
}

// buildStandardPageList builds a date-sorted list of pages for tags/categories.
func (b *Builder) buildStandardPageList(entry *parser.TaxonomyEntry) []map[string]any {
	// Sort pages by date (newest first)
	sortedPages := make([]*content.Node, len(entry.Pages))
	copy(sortedPages, entry.Pages)
	slices.SortFunc(sortedPages, func(a, b *content.Node) int {
		dateA := a.Config["date"].(time.Time)
		dateB := b.Config["date"].(time.Time)
		return dateB.Compare(dateA)
	})

	// Build page list
	pageList := make([]map[string]any, len(sortedPages))
	for i, node := range sortedPages {
		pageList[i] = map[string]any{
			"Permalink": node.Permalink,
			"Config":    node.Config,
		}
	}

	return pageList
}

// buildTaxonomiesMap builds a map of all taxonomies for template access.
func (b *Builder) buildTaxonomiesMap() map[string]any {
	taxonomies := make(map[string]any)

	for name, taxonomy := range b.parser.Taxonomies {
		terms := make([]map[string]any, 0, len(taxonomy.Entries))
		for slug, entry := range taxonomy.Entries {
			terms = append(terms, map[string]any{
				"Name":      entry.Term,
				"Slug":      slug,
				"Count":     entry.Count,
				"Permalink": b.buildTaxonomyPermalink(name, slug),
			})
		}

		// Sort alphabetically
		slices.SortFunc(terms, func(a, b map[string]any) int {
			return strings.Compare(a["Name"].(string), b["Name"].(string))
		})

		taxonomies[name] = map[string]any{
			"Name":  name,
			"Terms": terms,
			"Count": len(terms),
		}
	}

	return taxonomies
}

// buildTaxonomyPermalink builds a permalink for a taxonomy term.
func (b *Builder) buildTaxonomyPermalink(taxonomyName, slug string) string {
	return path.Join("/", taxonomyName, slug) + "/"
}

// buildTaxonomyURL builds a full URL for a taxonomy term.
func (b *Builder) buildTaxonomyURL(permalink string) string {
	return b.site.BaseURL + permalink
}

// generateSelectiveTaxonomies regenerates only affected taxonomy terms.
func (b *Builder) generateSelectiveTaxonomies(taxonomyName string, affectedTerms []string) error {
	taxonomy, exists := b.parser.Taxonomies[taxonomyName]

	if !exists || len(affectedTerms) == 0 {
		return nil
	}

	termTemplate := singularize(taxonomyName) + ".html"
	if !b.hasTemplate(termTemplate) {
		return nil
	}

	// Create taxonomy directory
	taxonomyDir := filepath.Join(b.site.OutputDir, taxonomyName)
	if err := os.MkdirAll(taxonomyDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s directory: %w", taxonomyName, err)
	}

	// Regenerate affected term pages
	for _, slug := range affectedTerms {
		// Skip empty slugs (invalid)
		if slug == "" || strings.TrimSpace(slug) == "" {
			continue
		}

		if entry, exists := taxonomy.Entries[slug]; exists && entry.Count > 0 {
			if err := b.generateTaxonomyTerm(taxonomyName, slug, entry, termTemplate); err != nil {
				return err
			}
		} else {
			// Term no longer exists or has zero posts, delete the term page
			termDir := filepath.Join(b.site.OutputDir, taxonomyName, slug)
			if err := os.RemoveAll(termDir); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove empty term directory %s: %w", termDir, err)
			}
		}
	}

	// Regenerate index page
	indexTemplate := taxonomyName + ".html"
	if b.hasTemplate(indexTemplate) {
		if err := b.generateTaxonomyIndex(taxonomyName, taxonomy, indexTemplate); err != nil {
			return err
		}
	}

	return nil
}

// singularize converts plural taxonomy names to singular for template names.
// This is a simple implementation - add more rules as needed.
func singularize(plural string) string {
	switch plural {
	case "categories":
		return "category"
	case "series":
		return "serie" // Following Hugo convention
	case "tags":
		return "tag"
	default:
		// Simple rule: remove trailing 's'
		if strings.HasSuffix(plural, "s") {
			return plural[:len(plural)-1]
		}
		return plural
	}
}
