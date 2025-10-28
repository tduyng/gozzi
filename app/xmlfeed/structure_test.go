package xmlfeed

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAtomFeed_XMLMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		feed     AtomFeed
		validate func(t *testing.T, xmlData string)
	}{
		{
			name: "complete atom feed",
			feed: AtomFeed{
				Title:   "Test Blog",
				ID:      "https://example.com",
				Updated: "2024-01-01T12:00:00Z",
				Author: &AtomAuthor{
					Name:  "Test Author",
					Email: "test@example.com",
				},
				Links: []AtomLink{
					{Rel: "self", Href: "https://example.com/atom.xml"},
					{Href: "https://example.com"},
				},
				Entries: []AtomEntry{
					{
						Title:     "Test Post",
						ID:        "https://example.com/test-post",
						Updated:   "2024-01-01T12:00:00Z",
						Published: "2024-01-01T10:00:00Z",
						Summary:   "Test summary",
						Content: struct {
							Type string `xml:"type,attr"`
							Data string `xml:",cdata"`
						}{
							Type: "html",
							Data: "<p>Test content</p>",
						},
						Links: []AtomLink{
							{Href: "https://example.com/test-post"},
						},
						Categories: []string{"test", "blog"},
					},
				},
			},
			validate: func(t *testing.T, xmlData string) {
				assert.Contains(t, xmlData, `xmlns="http://www.w3.org/2005/Atom"`)
				assert.Contains(t, xmlData, `<title>Test Blog</title>`)
				assert.Contains(t, xmlData, `<id>https://example.com</id>`)
				assert.Contains(t, xmlData, `<updated>2024-01-01T12:00:00Z</updated>`)
				assert.Contains(t, xmlData, `<name>Test Author</name>`)
				assert.Contains(t, xmlData, `<email>test@example.com</email>`)
				assert.Contains(t, xmlData, `rel="self"`)
				assert.Contains(t, xmlData, `href="https://example.com/atom.xml"`)
				assert.Contains(t, xmlData, `<![CDATA[<p>Test content</p>]]>`)
				assert.Contains(t, xmlData, `<category>test</category>`)
				assert.Contains(t, xmlData, `<category>blog</category>`)
			},
		},
		{
			name: "minimal atom feed",
			feed: AtomFeed{
				Title:   "Minimal Blog",
				ID:      "https://minimal.com",
				Updated: "2024-01-01T12:00:00Z",
				Links: []AtomLink{
					{Href: "https://minimal.com"},
				},
				Entries: []AtomEntry{},
			},
			validate: func(t *testing.T, xmlData string) {
				assert.Contains(t, xmlData, `<title>Minimal Blog</title>`)
				assert.Contains(t, xmlData, `<id>https://minimal.com</id>`)
				assert.Contains(t, xmlData, `href="https://minimal.com"`)
				assert.NotContains(t, xmlData, `<author>`)
				assert.NotContains(t, xmlData, `<entry>`)
			},
		},
		{
			name: "feed with entry without optional fields",
			feed: AtomFeed{
				Title:   "Simple Feed",
				ID:      "https://simple.com",
				Updated: "2024-01-01T12:00:00Z",
				Links: []AtomLink{
					{Href: "https://simple.com"},
				},
				Entries: []AtomEntry{
					{
						Title:     "Simple Post",
						ID:        "https://simple.com/post",
						Updated:   "2024-01-01T12:00:00Z",
						Published: "2024-01-01T10:00:00Z",
						Content: struct {
							Type string `xml:"type,attr"`
							Data string `xml:",cdata"`
						}{
							Type: "html",
							Data: "Simple content",
						},
						Links: []AtomLink{
							{Href: "https://simple.com/post"},
						},
					},
				},
			},
			validate: func(t *testing.T, xmlData string) {
				assert.Contains(t, xmlData, `<title>Simple Post</title>`)
				assert.Contains(t, xmlData, `<![CDATA[Simple content]]>`)
				assert.NotContains(t, xmlData, `<summary>`)
				assert.NotContains(t, xmlData, `<category>`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlData, err := xml.MarshalIndent(tt.feed, "", "  ")
			require.NoError(t, err)

			xmlString := string(xmlData)
			tt.validate(t, xmlString)

			// Verify XML is valid by unmarshaling back
			var parsedFeed AtomFeed
			err = xml.Unmarshal(xmlData, &parsedFeed)
			assert.NoError(t, err)
		})
	}
}

func TestAtomAuthor_XMLMarshaling(t *testing.T) {
	tests := []struct {
		name   string
		author AtomAuthor
		expect string
	}{
		{
			name: "author with email",
			author: AtomAuthor{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			expect: `<name>John Doe</name>` + "\n" + `  <email>john@example.com</email>`,
		},
		{
			name: "author without email",
			author: AtomAuthor{
				Name: "Jane Doe",
			},
			expect: `<name>Jane Doe</name>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlData, err := xml.MarshalIndent(tt.author, "", "  ")
			require.NoError(t, err)

			xmlString := string(xmlData)
			for _, expectPart := range strings.Split(tt.expect, "\n") {
				assert.Contains(t, xmlString, strings.TrimSpace(expectPart))
			}
		})
	}
}

func TestAtomLink_XMLMarshaling(t *testing.T) {
	tests := []struct {
		name string
		link AtomLink
		want string
	}{
		{
			name: "link with rel attribute",
			link: AtomLink{
				Rel:  "self",
				Href: "https://example.com/atom.xml",
			},
			want: `rel="self" href="https://example.com/atom.xml"`,
		},
		{
			name: "link without rel attribute",
			link: AtomLink{
				Href: "https://example.com",
			},
			want: `href="https://example.com"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlData, err := xml.MarshalIndent(tt.link, "", "  ")
			require.NoError(t, err)

			xmlString := string(xmlData)
			assert.Contains(t, xmlString, tt.want)
		})
	}
}

func TestAtomEntry_XMLMarshaling(t *testing.T) {
	entry := AtomEntry{
		Title:     "Test Entry",
		ID:        "https://example.com/entry",
		Updated:   "2024-01-01T12:00:00Z",
		Published: "2024-01-01T10:00:00Z",
		Summary:   "Entry summary",
		Content: struct {
			Type string `xml:"type,attr"`
			Data string `xml:",cdata"`
		}{
			Type: "html",
			Data: "<h1>Entry Content</h1>",
		},
		Links: []AtomLink{
			{Href: "https://example.com/entry"},
		},
		Categories: []string{"tech", "programming"},
	}

	xmlData, err := xml.MarshalIndent(entry, "", "  ")
	require.NoError(t, err)

	xmlString := string(xmlData)
	assert.Contains(t, xmlString, `<title>Test Entry</title>`)
	assert.Contains(t, xmlString, `<id>https://example.com/entry</id>`)
	assert.Contains(t, xmlString, `<updated>2024-01-01T12:00:00Z</updated>`)
	assert.Contains(t, xmlString, `<published>2024-01-01T10:00:00Z</published>`)
	assert.Contains(t, xmlString, `<summary>Entry summary</summary>`)
	assert.Contains(t, xmlString, `type="html"`)
	assert.Contains(t, xmlString, `<![CDATA[<h1>Entry Content</h1>]]>`)
	assert.Contains(t, xmlString, `href="https://example.com/entry"`)
	assert.Contains(t, xmlString, `<category>tech</category>`)
	assert.Contains(t, xmlString, `<category>programming</category>`)
}

func TestSitemap_XMLMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		sitemap  Sitemap
		validate func(t *testing.T, xmlData string)
	}{
		{
			name: "complete sitemap",
			sitemap: Sitemap{
				URLs: []SitemapURL{
					{
						Loc:        "https://example.com/",
						LastMod:    "2024-01-01",
						ChangeFreq: "daily",
						Priority:   "1.0",
					},
					{
						Loc:        "https://example.com/about",
						LastMod:    "2024-01-01",
						ChangeFreq: "monthly",
						Priority:   "0.8",
					},
				},
			},
			validate: func(t *testing.T, xmlData string) {
				assert.Contains(t, xmlData, `xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`)
				assert.Contains(t, xmlData, `<loc>https://example.com/</loc>`)
				assert.Contains(t, xmlData, `<lastmod>2024-01-01</lastmod>`)
				assert.Contains(t, xmlData, `<changefreq>daily</changefreq>`)
				assert.Contains(t, xmlData, `<priority>1.0</priority>`)
				assert.Contains(t, xmlData, `<loc>https://example.com/about</loc>`)
				assert.Contains(t, xmlData, `<changefreq>monthly</changefreq>`)
				assert.Contains(t, xmlData, `<priority>0.8</priority>`)
			},
		},
		{
			name: "minimal sitemap",
			sitemap: Sitemap{
				URLs: []SitemapURL{
					{
						Loc: "https://minimal.com/",
					},
				},
			},
			validate: func(t *testing.T, xmlData string) {
				assert.Contains(t, xmlData, `<loc>https://minimal.com/</loc>`)
				assert.NotContains(t, xmlData, `<lastmod>`)
				assert.NotContains(t, xmlData, `<changefreq>`)
				assert.NotContains(t, xmlData, `<priority>`)
			},
		},
		{
			name: "empty sitemap",
			sitemap: Sitemap{
				URLs: []SitemapURL{},
			},
			validate: func(t *testing.T, xmlData string) {
				assert.Contains(t, xmlData, `xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`)
				assert.NotContains(t, xmlData, `<url>`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlData, err := xml.MarshalIndent(tt.sitemap, "", "  ")
			require.NoError(t, err)

			xmlString := string(xmlData)
			tt.validate(t, xmlString)

			// Verify XML is valid by unmarshaling back
			var parsedSitemap Sitemap
			err = xml.Unmarshal(xmlData, &parsedSitemap)
			assert.NoError(t, err)
		})
	}
}

func TestSitemapURL_XMLMarshaling(t *testing.T) {
	tests := []struct {
		name string
		url  SitemapURL
		want []string
	}{
		{
			name: "complete sitemap URL",
			url: SitemapURL{
				Loc:        "https://example.com/page",
				LastMod:    "2024-01-01",
				ChangeFreq: "weekly",
				Priority:   "0.6",
			},
			want: []string{
				`<loc>https://example.com/page</loc>`,
				`<lastmod>2024-01-01</lastmod>`,
				`<changefreq>weekly</changefreq>`,
				`<priority>0.6</priority>`,
			},
		},
		{
			name: "minimal sitemap URL",
			url: SitemapURL{
				Loc: "https://example.com/minimal",
			},
			want: []string{
				`<loc>https://example.com/minimal</loc>`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xmlData, err := xml.MarshalIndent(tt.url, "", "  ")
			require.NoError(t, err)

			xmlString := string(xmlData)
			for _, wantPart := range tt.want {
				assert.Contains(t, xmlString, wantPart)
			}
		})
	}
}

func TestXMLNamespaces(t *testing.T) {
	t.Run("atom feed namespace", func(t *testing.T) {
		feed := AtomFeed{
			Title:   "Test",
			ID:      "https://example.com",
			Updated: "2024-01-01T12:00:00Z",
		}

		xmlData, err := xml.MarshalIndent(feed, "", "  ")
		require.NoError(t, err)

		assert.Contains(t, string(xmlData), `xmlns="http://www.w3.org/2005/Atom"`)
	})

	t.Run("sitemap namespace", func(t *testing.T) {
		sitemap := Sitemap{
			URLs: []SitemapURL{
				{Loc: "https://example.com"},
			},
		}

		xmlData, err := xml.MarshalIndent(sitemap, "", "  ")
		require.NoError(t, err)

		assert.Contains(t, string(xmlData), `xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`)
	})
}

func TestStructureDefaults(t *testing.T) {
	t.Run("empty atom feed", func(t *testing.T) {
		feed := AtomFeed{}
		xmlData, err := xml.MarshalIndent(feed, "", "  ")
		require.NoError(t, err)

		// Should still have namespace and basic structure
		assert.Contains(t, string(xmlData), `xmlns="http://www.w3.org/2005/Atom"`)
	})

	t.Run("empty sitemap", func(t *testing.T) {
		sitemap := Sitemap{}
		xmlData, err := xml.MarshalIndent(sitemap, "", "  ")
		require.NoError(t, err)

		// Should still have namespace and basic structure
		assert.Contains(t, string(xmlData), `xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"`)
	})
}

func TestXMLUnmarshaling(t *testing.T) {
	t.Run("unmarshal atom feed", func(t *testing.T) {
		xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Test Feed</title>
  <id>https://example.com</id>
  <updated>2024-01-01T12:00:00Z</updated>
  <link href="https://example.com"/>
  <entry>
    <title>Test Entry</title>
    <id>https://example.com/entry</id>
    <updated>2024-01-01T12:00:00Z</updated>
    <published>2024-01-01T10:00:00Z</published>
    <content type="html"><![CDATA[<p>Content</p>]]></content>
    <link href="https://example.com/entry"/>
  </entry>
</feed>`

		var feed AtomFeed
		err := xml.Unmarshal([]byte(xmlData), &feed)
		require.NoError(t, err)

		assert.Equal(t, "Test Feed", feed.Title)
		assert.Equal(t, "https://example.com", feed.ID)
		assert.Equal(t, "2024-01-01T12:00:00Z", feed.Updated)
		assert.Len(t, feed.Links, 1)
		assert.Equal(t, "https://example.com", feed.Links[0].Href)
		assert.Len(t, feed.Entries, 1)
		assert.Equal(t, "Test Entry", feed.Entries[0].Title)
		assert.Equal(t, "<p>Content</p>", feed.Entries[0].Content.Data)
	})

	t.Run("unmarshal sitemap", func(t *testing.T) {
		xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com/</loc>
    <lastmod>2024-01-01</lastmod>
    <changefreq>daily</changefreq>
    <priority>1.0</priority>
  </url>
</urlset>`

		var sitemap Sitemap
		err := xml.Unmarshal([]byte(xmlData), &sitemap)
		require.NoError(t, err)

		assert.Len(t, sitemap.URLs, 1)
		assert.Equal(t, "https://example.com/", sitemap.URLs[0].Loc)
		assert.Equal(t, "2024-01-01", sitemap.URLs[0].LastMod)
		assert.Equal(t, "daily", sitemap.URLs[0].ChangeFreq)
		assert.Equal(t, "1.0", sitemap.URLs[0].Priority)
	})
}
