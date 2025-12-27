// This file handles generating redirect HTML pages for URL aliases.
// Aliases allow old URLs to redirect to new canonical page locations.
package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/utils"
)

// generateAliasRedirects creates redirect HTML files for all aliases defined in a node.
func (b *Builder) generateAliasRedirects(node *content.Node) error {
	if len(node.Aliases) == 0 {
		return nil
	}

	// Use full URL with permalink (includes trailing slash)
	canonicalURL := b.site.BaseURL + node.Permalink

	for _, alias := range node.Aliases {
		if err := b.generateSingleRedirect(alias, canonicalURL); err != nil {
			return err
		}
	}

	return nil
}

// generateSingleRedirect creates a redirect HTML file at the alias path.
func (b *Builder) generateSingleRedirect(aliasPath, targetURL string) error {
	// Normalize alias path - ensure it starts with /
	if !strings.HasPrefix(aliasPath, "/") {
		aliasPath = "/" + aliasPath
	}

	// Remove leading slash for filesystem path
	fsPath := strings.TrimPrefix(aliasPath, "/")

	// Create output path for the redirect HTML
	outputPath := filepath.Join(b.site.OutputDir, fsPath, "index.html")

	// Create redirect HTML content
	redirectHTML := b.createRedirectHTML(targetURL)

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_alias_directory",
			Component: "builder",
			Path:      filepath.Dir(outputPath),
		})
	}

	// Write redirect HTML file
	if err := os.WriteFile(outputPath, []byte(redirectHTML), 0644); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "write_alias_redirect",
			Component: "builder",
			Path:      outputPath,
		})
	}

	return nil
}

// createRedirectHTML generates the HTML content for a redirect page.
func (b *Builder) createRedirectHTML(targetURL string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="refresh" content="0; url=%s">
    <link rel="canonical" href="%s">
    <title>Redirecting...</title>
</head>
<body>
    <p>This page has moved to <a href="%s">%s</a>.</p>
</body>
</html>`, targetURL, targetURL, targetURL, targetURL)
}
