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
	return b.renderTemplate(node, outputPath, data)
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

	data := map[string]any{
		"Site": map[string]any{
			"Config": b.site.ToConfig(),
		},
		"Config": node.Config,
		"Page":   nodeMap, "Section": parentMap,
	}

	return b.renderTemplate(node, outputPath, data)
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
		if template, ok := node.Config["template"].(string); ok {
			key["Template"] = template
		}
		if tags, ok := node.Config["tags"]; ok {
			key["Tags"] = fmt.Sprint(tags)
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
				if blogSection, exists := b.parser.ContentMap["blog"]; exists && len(blogSection.Children) > 0 {
					blogChildKeys := make([]string, len(blogSection.Children))
					for i, post := range blogSection.Children {
						parts := []string{post.Path}

						if title, ok := post.Config["title"].(string); ok {
							parts = append(parts, title)
						}
						if date, ok := post.Config["date"].(time.Time); ok {
							parts = append(parts, date.Format("2006-01-02"))
						}
						if desc, ok := post.Config["description"].(string); ok {
							parts = append(parts, desc)
						}

						if extra, ok := post.Config["extra"]; ok {
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
							}
						}

						blogChildKeys[i] = strings.Join(parts, "|")
					}
					key["BlogPosts"] = blogChildKeys
				}
			} else if len(node.Children) > 0 {
				childKeys := make([]string, len(node.Children))
				for i, child := range node.Children {
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

					if node.Slug == "notes" {
						parts = append(parts, string(child.Content))
					}

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
		return b.renderTemplateDirect(tpl, outputPath, data)
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
