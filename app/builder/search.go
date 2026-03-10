package builder

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
)

// SearchEntry represents a single document in the search index
type SearchEntry struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
	URL         string   `json:"url"`
	Date        string   `json:"date"`
}

func (b *Builder) generateSearchIndex() error {
	var entries []SearchEntry

	// Reusable regex for stripping HTML
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	spaceRegex := regexp.MustCompile(`\s+`)

	// Walk all pages in the content map
	b.walkNodes(b.parser.ContentMap["."], func(n *content.Node) {
		// Only index individual pages, not section indexes
		if n.Type != content.NodeTypePage {
			return
		}

		// Skip drafts
		if config.GetBool(n.Config, "draft") {
			return
		}

		// Get standard fields safely
		title, _ := n.Config["title"].(string)
		if title == "" {
			return // Skip items without a title
		}

		description, _ := n.Config["description"].(string)

		var tags []string
		if t, ok := n.Config["tags"].([]string); ok {
			tags = t
		} else if tAny, ok := n.Config["tags"].([]any); ok {
			for _, tag := range tAny {
				if s, isStr := tag.(string); isStr {
					tags = append(tags, s)
				}
			}
		}

		dateStr := ""
		if d, ok := n.Config["date"].(time.Time); ok && !d.IsZero() {
			// Format to YYYY-MM-DD or full RFC3339 as needed
			dateStr = d.Format("2006-01-02")
		}

		// Clean the HTML content to plain text
		plainText := htmlRegex.ReplaceAllString(string(n.Content), " ")
		// Decode basic entities (similar to what you do in summary.go)
		plainText = strings.ReplaceAll(plainText, "&nbsp;", " ")
		plainText = strings.ReplaceAll(plainText, "&amp;", "&")
		plainText = strings.ReplaceAll(plainText, "&lt;", "<")
		plainText = strings.ReplaceAll(plainText, "&gt;", ">")
		plainText = strings.ReplaceAll(plainText, "&quot;", "\"")
		plainText = strings.ReplaceAll(plainText, "&#39;", "'")
		plainText = spaceRegex.ReplaceAllString(plainText, " ")
		plainText = strings.TrimSpace(plainText)

		// Limit content to 800 chars (same as your python script)
		if len(plainText) > 800 {
			plainText = plainText[:800]
		}

		entries = append(entries, SearchEntry{
			ID:          n.Permalink,
			Title:       title,
			Description: description,
			Content:     plainText,
			Tags:        tags,
			URL:         n.Permalink,
			Date:        dateStr,
		})
	})

	// Sort entries by date (newest first)
	slices.SortFunc(entries, func(a, b SearchEntry) int {
		return strings.Compare(b.Date, a.Date)
	})

	// Output as JSON
	outputPath := filepath.Join(b.site.OutputDir, "search-index.json")

	var jsonData []byte
	var err error

	if b.site.MinifyJSON {
		jsonData, err = json.Marshal(entries)
	} else {
		jsonData, err = json.MarshalIndent(entries, "", "  ")
	}

	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, jsonData, 0644)
}
