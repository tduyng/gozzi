// Atom feed and sitemap generation for the site builder.
// Creates XML feeds (Atom) and sitemaps for content discovery.
package builder

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
)

const (
	atomXSL    = `<?xml-stylesheet type="text/xsl" href="/atom.xsl"?>`
	sitemapXSL = `<?xml-stylesheet type="text/xsl" href="/sitemap.xsl"?>`
)

// AtomFeed represents an Atom XML feed structure.
type AtomFeed struct {
	XMLName xml.Name    `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string      `xml:"title"`
	ID      string      `xml:"id"`
	Updated string      `xml:"updated"`
	Author  *AtomAuthor `xml:"author,omitempty"`
	Links   []AtomLink  `xml:"link"`
	Entries []AtomEntry `xml:"entry"`
}

// AtomAuthor represents the author information in an Atom feed.
type AtomAuthor struct {
	Name  string `xml:"name"`
	Email string `xml:"email,omitempty"`
}

// AtomLink represents a link in an Atom feed entry.
type AtomLink struct {
	Rel  string `xml:"rel,attr,omitempty"`
	Href string `xml:"href,attr"`
}

// AtomEntry represents a single entry in an Atom feed.
type AtomEntry struct {
	Title     string `xml:"title"`
	ID        string `xml:"id"`
	Updated   string `xml:"updated"`
	Published string `xml:"published"`
	Summary   string `xml:"summary,omitempty"`
	Content   struct {
		Type string `xml:"type,attr"`
		Data string `xml:",cdata"`
	} `xml:"content"`
	Links      []AtomLink `xml:"link"`
	Categories []string   `xml:"category,omitempty"`
}

// Sitemap represents an XML sitemap for search engines.
type Sitemap struct {
	XMLName xml.Name     `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	URLs    []SitemapURL `xml:"url"`
}

// SitemapURL represents a single URL entry in a sitemap.
type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

func (b *Builder) generateAtomFeed() error {
	var entries []*content.Node
	b.walkNodes(b.parser.ContentMap["."], func(n *content.Node) {
		if n.Type != content.NodeTypePage {
			return
		}
		isDraft := config.GetBool(n.Config, "draft")
		generateFeed := config.GetBool(n.Config, "generate_feed")

		if !isDraft && generateFeed {
			entries = append(entries, n)
		}
	})

	slices.SortFunc(entries, func(a, b *content.Node) int {
		dateA := a.Config["date"].(time.Time)
		dateB := b.Config["date"].(time.Time)
		// Sort descending (newest first), so reverse comparison
		return dateB.Compare(dateA)
	})

	if len(entries) > 100 {
		entries = entries[:100]
	}

	feed := AtomFeed{
		Title:   b.site.Title,
		ID:      b.site.BaseURL,
		Updated: time.Now().UTC().Format(time.RFC3339),
		Links: []AtomLink{
			{Rel: "self", Href: b.site.BaseURL + "/atom.xml"},
			{Href: b.site.BaseURL},
		},
	}

	if b.site.Extra["author"] != nil {
		author := b.site.Extra["author"].(map[string]any)
		feed.Author = &AtomAuthor{
			Name:  author["name"].(string),
			Email: author["email"].(string),
		}
	}

	for _, entry := range entries {
		date := entry.Config["date"].(time.Time).UTC()
		updated := date
		if u, ok := entry.Config["updated"].(time.Time); ok {
			updated = u.UTC()
		}

		categories := make([]string, 0)
		if tags, ok := entry.Config["tags"].([]string); ok {
			categories = tags
		}

		feed.Entries = append(feed.Entries, AtomEntry{
			Title:     entry.Config["title"].(string),
			ID:        entry.Permalink,
			Published: date.Format(time.RFC3339),
			Updated:   updated.Format(time.RFC3339),
			Summary:   entry.Config["description"].(string),
			Content: struct {
				Type string `xml:"type,attr"`
				Data string `xml:",cdata"`
			}{
				Type: "html",
				Data: string(entry.Content),
			},
			Links: []AtomLink{
				{Href: entry.Permalink},
			},
			Categories: categories,
		})
	}

	return b.writeXMLFile("atom.xml", atomXSL, feed)
}

func (b *Builder) generateSitemap() error {
	var urls []SitemapURL

	// Content pages
	b.walkNodes(b.parser.ContentMap["."], func(n *content.Node) {
		if config.GetBool(n.Config, "draft") {
			return
		}

		lastMod := getLastMod(n)

		url := SitemapURL{
			Loc:     n.URL,
			LastMod: lastMod,
		}

		switch {
		case n.Type == content.NodeTypeSection && n.Path == ".":
			url.ChangeFreq = "daily"
			url.Priority = "1.0"
		case n.Type == content.NodeTypeSection:
			url.ChangeFreq = "weekly"
			url.Priority = "0.8"
		default:
			url.ChangeFreq = "monthly"
			url.Priority = "0.6"
		}

		urls = append(urls, url)
	})

	// Tag pages
	if b.hasTemplate("tags.html") {
		for tag := range b.parser.Tags {
			loc := b.buildTagURL(b.buildTagPermalink(tag))
			urls = append(urls, SitemapURL{
				Loc:        loc,
				ChangeFreq: "monthly",
				Priority:   "0.4",
			})
		}
	}

	return b.writeXMLFile("sitemap.xml", sitemapXSL, Sitemap{URLs: urls})
}

func (b *Builder) writeXMLFile(name string, xslHeader string, data any) error {
	path := filepath.Join(b.site.OutputDir, name)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	// Write XML header with XSL
	if _, err := file.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + xslHeader + "\n"); err != nil {
		return err
	}

	// Create XML encoder with proper settings
	enc := xml.NewEncoder(file)
	enc.Indent("", "  ")

	if err := enc.Encode(data); err != nil {
		return err
	}

	// Final newline
	_, err = file.WriteString("\n")
	return err
}

func getLastMod(n *content.Node) string {
	var lastMod time.Time

	// Check for valid updated date
	if updated, ok := n.Config["updated"].(time.Time); ok && !updated.IsZero() {
		lastMod = updated
	} else if date, ok := n.Config["date"].(time.Time); ok && !date.IsZero() {
		// Fall back to valid date
		lastMod = date
	} else {
		// Fallback to current time if no valid dates
		lastMod = time.Now()
	}

	return lastMod.UTC().Format("2006-01-02")
}
