// Shared test helpers for builder package tests
// Provides helper functions to create test builders and fixtures
package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/parser"
)

// createTestBuilderWithContent creates a fully initialized test builder with parsed content.
func createTestBuilderWithContent(t *testing.T, contentFiles map[string]string) (*Builder, string) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(oldWd) })
	os.Chdir(tempDir)

	// Create templates directory with basic template
	templateDir := filepath.Join(tempDir, "templates")
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	templateContent := `<html><body>{{.Page.Title}}</body></html>`
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"), []byte(templateContent), 0644))

	// Create essential fallback templates
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "default.html"), []byte(templateContent), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "404.html"), []byte(`<html><body>404 Not Found</body></html>`), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "home.html"), []byte(templateContent), 0644))

	// Create content directory and files
	contentDir := filepath.Join(tempDir, "content")
	require.NoError(t, os.MkdirAll(contentDir, 0755))

	for filename, content := range contentFiles {
		filePath := filepath.Join(contentDir, filename)
		fileDir := filepath.Dir(filePath)
		require.NoError(t, os.MkdirAll(fileDir, 0755))
		require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))
	}

	site := &config.Site{
		Title:     "Test Site",
		BaseURL:   "https://example.com",
		OutputDir: filepath.Join(tempDir, "output"),
	}

	// Create and initialize parser with content - this is the key!
	p := parser.NewParser(site)
	err := p.Parse(contentDir)
	require.NoError(t, err)

	b, err := NewBuilder(site, p)
	require.NoError(t, err)

	// Reload templates to ensure all templates are loaded
	err = b.ReloadTemplates()
	require.NoError(t, err)

	return b, tempDir
}

// createTestBuilder creates a minimal test builder without content files.
func createTestBuilder(t *testing.T) (*Builder, string) {
	return createTestBuilderWithContent(t, map[string]string{})
}
