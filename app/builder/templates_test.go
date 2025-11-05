// Test file for templates.go
// Contains tests for template loading and reloading functionality
package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/parser"
)

func TestReloadTemplates(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Create templates directory
	templateDir := filepath.Join(tempDir, "templates")
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	// Create initial template
	templateContent := `<html><body>{{.Page.Title}}</body></html>`
	templatePath := filepath.Join(templateDir, "post.html")
	require.NoError(t, os.WriteFile(templatePath, []byte(templateContent), 0644))

	site := &config.Site{Title: "Test Site", BaseURL: "https://example.com", OutputDir: "output"}
	p := parser.NewParser(site)
	b, err := NewBuilder(site, p)
	require.NoError(t, err)

	// Reload templates to ensure all templates are loaded
	err = b.ReloadTemplates()
	require.NoError(t, err)
}
