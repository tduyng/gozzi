package builder

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
)

func TestChangeTracker(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(oldWd) })
	require.NoError(t, os.Chdir(tempDir))

	// Create content directory
	contentDir := filepath.Join(tempDir, "content")
	require.NoError(t, os.MkdirAll(filepath.Join(contentDir, "blog"), 0755))

	// Create test content
	post1 := `+++
title = "First Post"
date = 2024-01-01T00:00:00Z
tags = ["go", "test"]
+++
Content here`
	require.NoError(t, os.WriteFile(filepath.Join(contentDir, "blog", "post1.md"), []byte(post1), 0644))

	post2 := `+++
title = "Second Post"
date = 2024-01-02T00:00:00Z
tags = ["go"]
+++
Content here`
	require.NoError(t, os.WriteFile(filepath.Join(contentDir, "blog", "post2.md"), []byte(post2), 0644))

	// Parse content
	site := &config.Site{
		Title:     "Test Site",
		BaseURL:   "https://example.com",
		OutputDir: filepath.Join(tempDir, "public"),
	}
	p := parser.NewParser(site)
	require.NoError(t, p.Parse(contentDir))

	t.Run("detects changed blog post", func(t *testing.T) {
		tracker := NewChangeTracker(p.ContentMap, p)

		changedFiles := []string{
			filepath.Join(contentDir, "blog", "post1.md"),
		}

		tracker.AnalyzeChanges(changedFiles, contentDir)

		assert.True(t, tracker.ShouldRegenerateFeed(), "Feed should regenerate for blog post change")
		assert.True(t, tracker.ShouldRegenerateSitemap(), "Sitemap should regenerate for any content change")
		assert.Equal(t, 1, tracker.GetChangedNodesCount(), "Should have 1 changed node")
		assert.Contains(t, tracker.GetAffectedTags(), "go", "Should affect 'go' tag")
		assert.Contains(t, tracker.GetAffectedTags(), "test", "Should affect 'test' tag")
	})

	t.Run("detects multiple changed files", func(t *testing.T) {
		tracker := NewChangeTracker(p.ContentMap, p)

		changedFiles := []string{
			filepath.Join(contentDir, "blog", "post1.md"),
			filepath.Join(contentDir, "blog", "post2.md"),
		}

		tracker.AnalyzeChanges(changedFiles, contentDir)

		assert.Equal(t, 2, tracker.GetChangedNodesCount(), "Should have 2 changed nodes")
		assert.True(t, tracker.ShouldRegenerateFeed(), "Feed should regenerate")
	})

	t.Run("tracks tag dependencies correctly", func(t *testing.T) {
		tracker := NewChangeTracker(p.ContentMap, p)

		// Only change post2 which has tag "go"
		changedFiles := []string{
			filepath.Join(contentDir, "blog", "post2.md"),
		}

		tracker.AnalyzeChanges(changedFiles, contentDir)

		affectedTags := tracker.GetAffectedTags()
		assert.Contains(t, affectedTags, "go", "Should affect 'go' tag")
		assert.NotContains(t, affectedTags, "test", "Should NOT affect 'test' tag (not in post2)")
	})

	t.Run("identifies blog posts correctly", func(t *testing.T) {
		tracker := NewChangeTracker(p.ContentMap, p)

		assert.True(t, tracker.isBlogPost("blog/my-post.md"))
		assert.True(t, tracker.isBlogPost("posts/another-post.md"))
		assert.False(t, tracker.isBlogPost("about/index.md"))
		assert.False(t, tracker.isBlogPost("_index.md"))
	})
}

func TestChangeTrackerWithSectionChanges(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(oldWd) })
	require.NoError(t, os.Chdir(tempDir))

	// Create content directory
	contentDir := filepath.Join(tempDir, "content")
	require.NoError(t, os.MkdirAll(filepath.Join(contentDir, "blog"), 0755))

	// Create section index
	sectionIndex := `+++
title = "Blog"
+++
Blog section content`
	require.NoError(t, os.WriteFile(filepath.Join(contentDir, "blog", "_index.md"), []byte(sectionIndex), 0644))

	// Parse content
	site := &config.Site{
		Title:     "Test Site",
		BaseURL:   "https://example.com",
		OutputDir: filepath.Join(tempDir, "public"),
	}
	p := parser.NewParser(site)
	require.NoError(t, p.Parse(contentDir))

	t.Run("detects section index change", func(t *testing.T) {
		tracker := NewChangeTracker(p.ContentMap, p)

		changedFiles := []string{
			filepath.Join(contentDir, "blog", "_index.md"),
		}

		tracker.AnalyzeChanges(changedFiles, contentDir)

		assert.True(t, tracker.ShouldRegenerateSitemap(), "Sitemap should regenerate")
		// Section changes should not affect feed (unless it's a blog post)
		assert.False(t, tracker.ShouldRegenerateFeed(), "Feed should not regenerate for section index")
		assert.Equal(t, 1, tracker.GetChangedNodesCount(), "Should have 1 changed node (section)")
	})
}

func TestChangeTrackerNodeLookup(t *testing.T) {
	// Create a simple content tree
	root := &content.Node{
		Type: content.NodeTypeSection,
		Path: "",
	}

	blogPost := &content.Node{
		Type:   content.NodeTypePage,
		Path:   "blog/my-post.md",
		Parent: root,
		Config: map[string]any{
			"title": "My Post",
			"tags":  []string{"go", "testing"},
			"date":  time.Now(),
		},
	}

	root.Children = []*content.Node{blogPost}

	contentMap := map[string]*content.Node{
		"":     root,
		"blog": root, // Simplified - in reality blog would be separate section
	}

	site := &config.Site{
		Title:   "Test",
		BaseURL: "https://example.com",
	}
	p := parser.NewParser(site)

	tracker := NewChangeTracker(contentMap, p)

	// Manually mark node as changed (simulating file detection)
	tracker.changedNodes[blogPost] = true
	tracker.trackAffectedTags(blogPost)

	assert.True(t, tracker.ShouldRegenerateNode(blogPost), "Should regenerate the changed node")
	assert.Contains(t, tracker.GetAffectedTags(), "go", "Should affect go tag")
	assert.Contains(t, tracker.GetAffectedTags(), "testing", "Should affect testing tag")
}
