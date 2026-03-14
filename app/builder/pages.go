package builder

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/minify"
	"github.com/tduyng/gozzi/app/utils"
)

func (b *Builder) generateSection(node *content.Node) error {
	outputPath := filepath.Join(b.site.OutputDir, node.Slug, "index.html")

	nodeMap := node.ToMap()

	data := map[string]any{
		"Site": map[string]any{
			"Config":     b.site.ToConfig(),
			"Taxonomies": b.buildTaxonomiesMap(),
		},
		"Config":  node.Config,
		"Page":    nodeMap,
		"Section": nodeMap,
	}

	if node.Config["assets"] != "" {
		if err := b.copyPageAssets(node); err != nil {
			return err
		}
	}

	if err := b.renderTemplate(node, outputPath, data); err != nil {
		return err
	}

	return b.generateAliasRedirects(node)
}

func (b *Builder) generatePage(node *content.Node) error {
	outputPath := filepath.Join(
		b.site.OutputDir,
		node.Slug,
		"index.html",
	)

	if node.Config["assets"] != "" {
		if err := b.copyPageAssets(node); err != nil {
			return err
		}
	}
	nodeMap := node.ToMap()

	var parentMap map[string]any
	if node.Parent != nil {
		parentMap = node.Parent.ToMapMinimal()
	} else {
		parentMap = make(map[string]any)
	}

	if seriesName, ok := node.Config["series"].(string); ok && seriesName != "" {
		if seriesNav := b.getSeriesNavigation(node, seriesName); seriesNav != nil {
			nodeMap["Series"] = seriesNav
		}
	}

	data := map[string]any{
		"Site": map[string]any{
			"Config":     b.site.ToConfig(),
			"Taxonomies": b.buildTaxonomiesMap(),
		},
		"Config": node.Config,
		"Page":   nodeMap, "Section": parentMap,
	}

	if err := b.renderTemplate(node, outputPath, data); err != nil {
		return err
	}

	return b.generateAliasRedirects(node)
}

func (b *Builder) generate404Page() error {
	tpl := b.templ.Lookup("404.html")
	if tpl == nil {
		return nil
	}

	outputPath := filepath.Join(b.site.OutputDir, "404.html")
	data := map[string]any{
		"Site": map[string]any{
			"Config":     b.site.ToConfig(),
			"Taxonomies": b.buildTaxonomiesMap(),
		},
		"Page": map[string]any{
			"Title": "Page Not Found",
			"URL":   path.Join(b.site.BaseURL, "404.html"),
		},
	}

	return b.renderTemplate(nil, outputPath, data)
}

func (b *Builder) renderTemplate(node *content.Node, outputPath string, data any, templateNames ...string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	var tpl *template.Template

	if node != nil {
		for _, name := range node.TemplateChain() {
			tpl = b.templ.Lookup(name)
			if tpl != nil {
				break
			}
		}
	} else if len(templateNames) > 0 {
		for _, name := range templateNames {
			tpl = b.templ.Lookup(name)
			if tpl != nil {
				break
			}
		}
	} else {
		tpl = b.templ.Lookup("404.html")
	}

	if tpl == nil {
		return utils.WrapWithContext(fmt.Errorf("no template found"), utils.ErrTemplate, utils.ErrorContext{
			Operation: "find_template",
			Component: "builder",
			Path:      outputPath,
		})
	}

	content, err := b.executeTemplate(tpl, data)
	if err != nil {
		return utils.WrapWithContext(err, utils.ErrTemplate, utils.ErrorContext{
			Operation: "execute_template",
			Component: "builder",
			Path:      outputPath,
		})
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_output_directory",
			Component: "builder",
			Path:      filepath.Dir(outputPath),
		})
	}

	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "write_html_output",
			Component: "builder",
			Path:      outputPath,
		})
	}

	return nil
}

func (b *Builder) executeTemplate(tpl *template.Template, data any) ([]byte, error) {
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	content := buf.Bytes()

	// Apply minification if enabled
	if b.site.MinifyHTML {
		m := minify.New()
		if minified, err := m.MinifyHTML(content); err == nil {
			content = minified
		}
	}

	return content, nil
}

func (b *Builder) copyPageAssets(node *content.Node) error {
	assetsValue, exists := node.Config["assets"]
	if !exists {
		return nil
	}

	assets, ok := assetsValue.(string)
	if !ok {
		return nil
	}

	if b.site.ProjectDir != "" {
		assets = filepath.Join(b.site.ProjectDir, assets)
	}

	dest := filepath.Join(
		b.site.OutputDir,
		node.Slug,
		filepath.Base(assets),
	)

	if _, err := os.Stat(assets); os.IsNotExist(err) {
		return os.RemoveAll(dest)
	}

	if err := os.RemoveAll(dest); err != nil && !os.IsNotExist(err) {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "remove_old_assets",
			Component: "builder",
			Path:      dest,
		})
	}

	return copyDir(assets, dest)
}

// getSeriesNavigation builds series navigation data for a page.
func (b *Builder) getSeriesNavigation(node *content.Node, seriesName string) map[string]any {
	// Get the series taxonomy
	seriesTaxonomy, exists := b.parser.Taxonomies["series"]
	if !exists {
		return nil
	}

	// Find the series entry
	slug := urlizeHelper(seriesName)
	entry, exists := seriesTaxonomy.Entries[slug]
	if !exists {
		return nil
	}

	// Get ordered series pages
	seriesPages := entry.GetSeriesPages()
	if len(seriesPages) == 0 {
		return nil
	}

	// Find current page position
	var currentPosition int
	var prevPage, nextPage map[string]any

	for i, sp := range seriesPages {
		if sp.Node.Path == node.Path {
			currentPosition = sp.Position

			// Get previous page
			if i > 0 {
				prev := seriesPages[i-1]
				prevPage = map[string]any{
					"Title":     prev.Node.Config["title"],
					"Permalink": prev.Node.Permalink,
					"Position":  prev.Position,
				}
			}

			// Get next page
			if i < len(seriesPages)-1 {
				next := seriesPages[i+1]
				nextPage = map[string]any{
					"Title":     next.Node.Config["title"],
					"Permalink": next.Node.Permalink,
					"Position":  next.Position,
				}
			}
			break
		}
	}

	// Build all posts list for series navigation
	allPosts := make([]map[string]any, len(seriesPages))
	for i, sp := range seriesPages {
		allPosts[i] = map[string]any{
			"Title":     sp.Node.Config["title"],
			"Permalink": sp.Node.Permalink,
			"Position":  sp.Position,
			"Date":      sp.Node.Config["date"],
		}
	}

	// Build series navigation data
	return map[string]any{
		"Name":         entry.Term,
		"Slug":         slug,
		"Permalink":    b.buildTaxonomyPermalink("series", slug),
		"TotalPosts":   len(seriesPages),
		"CurrentPart":  currentPosition,
		"PreviousPost": prevPage,
		"NextPost":     nextPage,
		"AllPosts":     allPosts,
	}
}

// urlizeHelper converts a string to URL-friendly slug (helper for series nav).
func urlizeHelper(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	slug := result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	return strings.Trim(slug, "-")
}
