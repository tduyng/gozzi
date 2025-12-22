package builder

import (
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
	"github.com/tduyng/gozzi/app/template/funcs"
	"os"
)

func (b *Builder) generateTagPages() error {
	tagsTemplateExists := b.hasTemplate("tags.html")
	tagTemplateExists := b.hasTemplate("tag.html")

	if !tagsTemplateExists && !tagTemplateExists {
		return nil
	}

	tagsDir := filepath.Join(b.site.OutputDir, "tags")
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return err
	}

	if tagsTemplateExists {
		if err := b.generateTagsIndex(); err != nil {
			return err
		}
	}

	if tagTemplateExists {
		for tag, pages := range b.parser.Tags {
			if err := b.generateTagPage(tag, pages); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Builder) generateTagsIndex() error {
	tags := make([]map[string]any, 0, len(b.parser.Tags))
	for tag, entry := range b.parser.Tags {
		permalink := b.buildTagPermalink(tag)
		tags = append(tags, map[string]any{
			"Name":      tag,
			"Count":     entry.Count,
			"Permalink": permalink,
			"URL":       b.buildTagURL(permalink),
		})
	}

	slices.SortFunc(tags, func(a, b map[string]any) int {
		return strings.Compare(a["Name"].(string), b["Name"].(string))
	})

	data := map[string]any{
		"Site": map[string]any{
			"Config": b.site.ToConfig(),
		},
		"Page": map[string]any{
			"Title": "Tags",
			"Tags":  tags,
			"Path":  "/tags",
		},
	}

	return b.renderTemplate(nil,
		filepath.Join(b.site.OutputDir, "tags", "index.html"),
		data, "tags.html")
}

func (b *Builder) generateTagPage(tag string, entry *parser.TagEntry) error {
	sortedPages := make([]*content.Node, len(entry.Pages))
	copy(sortedPages, entry.Pages)
	slices.SortFunc(sortedPages, func(a, b *content.Node) int {
		dateA := a.Config["date"].(time.Time)
		dateB := b.Config["date"].(time.Time)
		return dateB.Compare(dateA)
	})

	pageList := make([]map[string]any, len(sortedPages))
	for i, node := range sortedPages {
		pageList[i] = map[string]any{
			"Permalink": node.Permalink,
			"Config":    node.Config,
		}
	}

	data := map[string]any{
		"Site": map[string]any{
			"Config": b.site.ToConfig(),
		},
		"Page": map[string]any{
			"Title":     fmt.Sprintf("Tag: %s", tag),
			"Tag":       tag,
			"Pages":     pageList,
			"Permalink": b.buildTagPermalink(tag),
			"Path":      b.buildTagPermalink(tag),
		},
	}

	outputPath := filepath.Join(
		b.site.OutputDir,
		"tags",
		funcs.Urlize(tag),
		"index.html",
	)

	return b.renderTemplate(nil, outputPath, data, "tag.html")
}

func (b *Builder) buildTagPermalink(tag string) string {
	return path.Join("/tags", funcs.Urlize(tag)) + "/"
}

func (b *Builder) buildTagURL(tagLink string) string {
	return b.site.BaseURL + tagLink
}

// generateSelectiveTags generates only the specified tag pages.
func (b *Builder) generateSelectiveTags(affectedTags []string) error {
	tagTemplateExists := b.hasTemplate("tag.html")

	if !tagTemplateExists || len(affectedTags) == 0 {
		return nil
	}

	tagsDir := filepath.Join(b.site.OutputDir, "tags")
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return err
	}

	for _, tag := range affectedTags {
		if entry, exists := b.parser.Tags[tag]; exists {
			if err := b.generateTagPage(tag, entry); err != nil {
				return err
			}
		}
	}

	if b.hasTemplate("tags.html") {
		if err := b.generateTagsIndex(); err != nil {
			return err
		}
	}

	return nil
}
