package builder

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/tduyng/gozzi/app/cache"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/minify"
	"github.com/tduyng/gozzi/app/utils"
)

func (b *Builder) generateSection(node *content.Node) error {
	outputPath := filepath.Join(b.site.OutputDir, node.Slug, "index.html")

	// Templates render {{.Section.Content}} so sections need full content
	nodeMap := node.ToMap()

	data := map[string]any{
		"Site": map[string]any{
			"Config": b.site.ToConfig(),
		},
		"Config":  node.Config,
		"Page":    nodeMap,
		"Section": nodeMap,
	}

	// Generate the main section page
	if err := b.renderTemplate(node, outputPath, data); err != nil {
		return err
	}

	// Generate alias redirects
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
	// Parent section doesn't need full content
	parentMap := node.Parent.ToMapMinimal()

	// Add series navigation if this page is part of a series
	if seriesName, ok := node.Config["series"].(string); ok && seriesName != "" {
		if seriesNav := b.getSeriesNavigation(node, seriesName); seriesNav != nil {
			nodeMap["Series"] = seriesNav
		}
	}

	data := map[string]any{
		"Site": map[string]any{
			"Config": b.site.ToConfig(),
		},
		"Config": node.Config,
		"Page":   nodeMap, "Section": parentMap,
	}

	// Generate the main page
	if err := b.renderTemplate(node, outputPath, data); err != nil {
		return err
	}

	// Generate alias redirects
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
			"Config": b.site.ToConfig(),
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
	var tplName string

	if node != nil {
		for _, name := range node.TemplateChain() {
			tpl = b.templ.Lookup(name)
			if tpl != nil {
				tplName = name
				break
			}
		}
	} else if len(templateNames) > 0 {
		for _, name := range templateNames {
			tpl = b.templ.Lookup(name)
			if tpl != nil {
				tplName = name
				break
			}
		}
	} else {
		tpl = b.templ.Lookup("404.html")
		tplName = "404.html"
	}

	if tpl == nil {
		return utils.WrapWithContext(fmt.Errorf("no template found"), utils.ErrTemplate, utils.ErrorContext{
			Operation: "find_template",
			Component: "builder",
			Path:      outputPath,
		})
	}

	var cacheKey any
	if node != nil {
		// Convert time.Time to stable string for cache key
		key := map[string]any{
			"Path":      node.Path,
			"Content":   string(node.Content),
			"WordCount": node.WordCount,
			"ReadTime":  node.ReadTime,
		}

		if title, ok := node.Config["title"].(string); ok {
			key["Title"] = title
		}
		if date, ok := node.Config["date"].(time.Time); ok {
			key["Date"] = date.Format("2006-01-02")
		}
		if desc, ok := node.Config["description"].(string); ok {
			key["Description"] = desc
		}
		if template, ok := node.Config["template"].(string); ok {
			key["Template"] = template
		}
		if tags, ok := node.Config["tags"]; ok {
			key["Tags"] = fmt.Sprint(tags)
		}
		// Include series in cache key (affects navigation rendering)
		if series, ok := node.Config["series"]; ok {
			key["Series"] = fmt.Sprint(series)
		}
		if seriesOrder, ok := node.Config["series_order"]; ok {
			key["SeriesOrder"] = fmt.Sprint(seriesOrder)
		}
		// Include categories in cache key
		if categories, ok := node.Config["categories"]; ok {
			key["Categories"] = fmt.Sprint(categories)
		}
		// Extra config affects template rendering via partials
		if extra, ok := node.Config["extra"]; ok {
			key["Extra"] = extra
		}

		// Section pages include children metadata in cache key
		if node.Type == content.NodeTypeSection {
			isHomepage := node.Path == "." || node.Path == "" || node.Path == "_index.md" ||
				node.Slug == "" || node.Slug == "/"

			if isHomepage {
				// Homepage can reference ANY section's data, so include all sections in cache key
				// This ensures homepage cache invalidates when any child page changes
				allSections := make(map[string][]string)
				for sectionPath, sectionNode := range b.parser.ContentMap {
					if sectionNode.Type == content.NodeTypeSection && sectionPath != "." && len(sectionNode.Children) > 0 {
						childKeys := make([]string, len(sectionNode.Children))
						for i, child := range sectionNode.Children {
							parts := b.buildChildCacheKeyParts(child)
							childKeys[i] = strings.Join(parts, "|")
						}
						allSections[sectionPath] = childKeys
					}
				}
				if len(allSections) > 0 {
					key["AllSections"] = allSections
				}
			} else if len(node.Children) > 0 {
				// Regular sections only need their own children in cache key
				childKeys := make([]string, len(node.Children))
				for i, child := range node.Children {
					parts := b.buildChildCacheKeyParts(child)
					childKeys[i] = strings.Join(parts, "|")
				}
				key["Children"] = childKeys
			}
		}

		if b.site.Extra != nil {
			key["SiteExtra"] = b.site.Extra
		}

		cacheKey = key
	} else {
		cacheKey = b.createStableCacheKey(tplName, data)
	}

	dataHash, err := cache.ComputeDataHash(cacheKey)
	if err != nil {
		return utils.WrapWithContext(utils.ErrTemplate, err, utils.ErrorContext{
			Operation: "compute_cache_hash",
			Component: "builder",
			Path:      outputPath,
		})
	}

	content, cached, err := b.renderCache.GetOrCompute(tplName, dataHash, func() ([]byte, error) {
		return b.executeTemplate(tpl, data)
	})

	if err != nil {
		return utils.WrapWithContext(err, utils.ErrTemplate, utils.ErrorContext{
			Operation: "execute_template",
			Component: "builder",
			Path:      outputPath,
		})
	}

	if !cached && b.site.MinifyHTML {
		m := minify.New()
		if minified, err := m.MinifyHTML(content); err == nil {
			content = minified
			b.renderCache.Set(tplName, dataHash, content)
		}
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

// executeTemplate renders a template to bytes.
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

// renderTemplateDirect renders template without caching.
func (b *Builder) renderTemplateDirect(tpl *template.Template, outputPath string, data any) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_output_directory",
			Component: "builder",
			Path:      filepath.Dir(outputPath),
		})
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return utils.WrapWithContext(err, utils.ErrTemplate, utils.ErrorContext{
			Operation: "execute_template",
			Component: "builder",
			Path:      outputPath,
		})
	}

	content := buf.Bytes()

	if b.site.MinifyHTML {
		m := minify.New()
		minified, err := m.MinifyHTML(content)
		if err == nil {
			content = minified
		}
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

// createStableCacheKey creates a deterministic cache key for auxiliary pages.
func (b *Builder) createStableCacheKey(templateName string, data any) map[string]any {
	key := map[string]any{
		"Template": templateName,
	}

	dataMap, ok := data.(map[string]any)
	if !ok {
		return key
	}

	pageData, ok := dataMap["Page"].(map[string]any)
	if !ok {
		return key
	}

	if tag, ok := pageData["Tag"].(string); ok {
		key["Tag"] = tag

		if pages, ok := pageData["Pages"].([]map[string]any); ok {
			pageKeys := make([]string, len(pages))
			for i, page := range pages {
				parts := []string{}

				if permalink, ok := page["Permalink"].(string); ok {
					parts = append(parts, permalink)
				}

				if config, ok := page["Config"].(map[string]any); ok {
					if title, ok := config["title"].(string); ok {
						parts = append(parts, title)
					}
					if date, ok := config["date"].(time.Time); ok {
						parts = append(parts, date.Format("2006-01-02"))
					}
					if extra, ok := config["extra"].(map[string]any); ok {
						if featured, ok := extra["featured"].(bool); ok {
							parts = append(parts, fmt.Sprintf("featured:%v", featured))
						}
					}
				}

				pageKeys[i] = strings.Join(parts, "|")
			}
			key["Pages"] = pageKeys
		}
	}

	if tags, ok := pageData["Tags"].([]map[string]any); ok {
		tagKeys := make([]string, len(tags))
		for i, tag := range tags {
			name := fmt.Sprint(tag["Name"])
			count := fmt.Sprint(tag["Count"])
			permalink := fmt.Sprint(tag["Permalink"])
			tagKeys[i] = fmt.Sprintf("%s:%s:%s", name, count, permalink)
		}
		key["Tags"] = tagKeys
	}

	// Handle taxonomy index pages with Terms field (series, categories, etc.)
	if terms, ok := pageData["Terms"].([]map[string]any); ok {
		termKeys := make([]string, len(terms))
		for i, term := range terms {
			name := fmt.Sprint(term["Name"])
			slug := fmt.Sprint(term["Slug"])
			count := fmt.Sprint(term["Count"])
			permalink := fmt.Sprint(term["Permalink"])
			termKeys[i] = fmt.Sprintf("%s:%s:%s:%s", name, slug, count, permalink)
		}
		key["Terms"] = termKeys
	}

	if title, ok := pageData["Title"].(string); ok {
		key["Title"] = title
	}

	if path, ok := pageData["Path"].(string); ok {
		key["Path"] = path
	}

	if siteData, ok := dataMap["Site"].(map[string]any); ok {
		if config, ok := siteData["Config"].(map[string]any); ok {
			if baseURL, ok := config["base_url"]; ok {
				key["SiteBaseURL"] = fmt.Sprint(baseURL)
			}
			if siteTitle, ok := config["title"]; ok {
				key["SiteTitle"] = fmt.Sprint(siteTitle)
			}
		}
	}

	return key
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

	// Build series navigation data
	return map[string]any{
		"Name":         entry.Term,
		"Slug":         slug,
		"Permalink":    b.buildTaxonomyPermalink("series", slug),
		"TotalPosts":   len(seriesPages),
		"CurrentPart":  currentPosition,
		"PreviousPost": prevPage,
		"NextPost":     nextPage,
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

// buildChildCacheKeyParts extracts cache-relevant data from a child node.
// This ensures consistent cache key generation for child pages across different contexts.
func (b *Builder) buildChildCacheKeyParts(child *content.Node) []string {
	parts := []string{child.Path}

	if title, ok := child.Config["title"].(string); ok {
		parts = append(parts, title)
	}
	if date, ok := child.Config["date"].(time.Time); ok {
		parts = append(parts, date.Format("2006-01-02"))
	}
	if desc, ok := child.Config["description"].(string); ok {
		parts = append(parts, desc)
	}

	// Include extra config (affects template rendering)
	if extra, ok := child.Config["extra"]; ok {
		if extraMap, ok := extra.(map[string]any); ok {
			keys := make([]string, 0, len(extraMap))
			for k := range extraMap {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			extraParts := make([]string, 0, len(keys))
			for _, k := range keys {
				extraParts = append(extraParts, fmt.Sprintf("%s=%v", k, extraMap[k]))
			}
			parts = append(parts, strings.Join(extraParts, ","))
		} else {
			parts = append(parts, fmt.Sprintf("%v", extra))
		}
	}

	return parts
}
