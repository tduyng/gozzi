// Test file for static.go
// Contains tests for static asset copying functionality
package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyStaticAssets(t *testing.T) {
	b, tempDir := createTestBuilder(t)

	// Create static directory with files
	staticDir := filepath.Join(tempDir, "static")
	require.NoError(t, os.MkdirAll(filepath.Join(staticDir, "css"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(staticDir, "style.css"), []byte("body { margin: 0; }"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(staticDir, "css", "main.css"), []byte(".header { color: blue; }"), 0644))

	err := b.copyStaticAssets()
	assert.NoError(t, err)

	// Verify files were copied to output
	content, err := os.ReadFile(filepath.Join(b.site.OutputDir, "style.css"))
	assert.NoError(t, err)
	assert.Equal(t, "body { margin: 0; }", string(content))

	content, err = os.ReadFile(filepath.Join(b.site.OutputDir, "css", "main.css"))
	assert.NoError(t, err)
	assert.Equal(t, ".header { color: blue; }", string(content))
}

func TestCopyCSSWithMinify(t *testing.T) {
	b, tempDir := createTestBuilder(t)
	b.site.MinifyCSS = true

	srcPath := filepath.Join(tempDir, "test.css")
	dstPath := filepath.Join(b.site.OutputDir, "test.css")

	cssContent := `
body {
    margin: 0;
    padding: 0;
}

/* Comment */
.container {
    max-width: 1200px;
}
`
	require.NoError(t, os.WriteFile(srcPath, []byte(cssContent), 0644))

	err := b.copyCSSWithMinify(srcPath, dstPath)
	assert.NoError(t, err)
	assert.FileExists(t, dstPath)

	content, err := os.ReadFile(dstPath)
	assert.NoError(t, err)
	assert.Less(t, len(content), len(cssContent), "Minified CSS should be smaller")
	assert.NotContains(t, string(content), "/* Comment */")
}

func TestCopyJSWithMinify(t *testing.T) {
	b, tempDir := createTestBuilder(t)
	b.site.MinifyJS = true

	srcPath := filepath.Join(tempDir, "test.js")
	dstPath := filepath.Join(b.site.OutputDir, "test.js")

	jsContent := `
function hello() {
    console.log('Hello World');
}

// This is a comment
const x = 1;
`
	require.NoError(t, os.WriteFile(srcPath, []byte(jsContent), 0644))

	err := b.copyJSWithMinify(srcPath, dstPath)
	assert.NoError(t, err)
	assert.FileExists(t, dstPath)

	content, err := os.ReadFile(dstPath)
	assert.NoError(t, err)
	assert.Less(t, len(content), len(jsContent), "Minified JS should be smaller")
	assert.NotContains(t, string(content), "// This is a comment")
}

func TestCopyJSONWithMinify(t *testing.T) {
	b, tempDir := createTestBuilder(t)
	b.site.MinifyJSON = true

	srcPath := filepath.Join(tempDir, "test.json")
	dstPath := filepath.Join(b.site.OutputDir, "test.json")

	jsonContent := `{
    "name": "Test",
    "value": 42,
    "active": true
}`
	require.NoError(t, os.WriteFile(srcPath, []byte(jsonContent), 0644))

	err := b.copyJSONWithMinify(srcPath, dstPath)
	assert.NoError(t, err)
	assert.FileExists(t, dstPath)

	content, err := os.ReadFile(dstPath)
	assert.NoError(t, err)
	assert.Less(t, len(content), len(jsonContent), "Minified JSON should be smaller")
	assert.NotContains(t, string(content), "\n    ")
}

func TestCopySVGWithMinify(t *testing.T) {
	b, tempDir := createTestBuilder(t)
	b.site.MinifySVG = true

	srcPath := filepath.Join(tempDir, "test.svg")
	dstPath := filepath.Join(b.site.OutputDir, "test.svg")

	svgContent := `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
    <!-- Icon -->
    <circle cx="50" cy="50" r="40" fill="red" />
</svg>`
	require.NoError(t, os.WriteFile(srcPath, []byte(svgContent), 0644))

	err := b.copySVGWithMinify(srcPath, dstPath)
	assert.NoError(t, err)
	assert.FileExists(t, dstPath)

	content, err := os.ReadFile(dstPath)
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(content), len(svgContent), "Minified SVG should be smaller or equal")
	assert.NotContains(t, string(content), "<!-- Icon -->")
}

func TestCopyXMLWithMinify(t *testing.T) {
	b, tempDir := createTestBuilder(t)
	b.site.MinifyXML = true

	srcPath := filepath.Join(tempDir, "test.xml")
	dstPath := filepath.Join(b.site.OutputDir, "test.xml")

	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<config>
    <setting name="timeout" value="30" />
    <setting name="retries" value="3" />
</config>`
	require.NoError(t, os.WriteFile(srcPath, []byte(xmlContent), 0644))

	err := b.copyXMLWithMinify(srcPath, dstPath)
	assert.NoError(t, err)
	assert.FileExists(t, dstPath)

	content, err := os.ReadFile(dstPath)
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(content), len(xmlContent), "Minified XML should be smaller or equal")
	assert.NotContains(t, string(content), "\n    ")
}

func TestMinifyWithInvalidContent(t *testing.T) {
	b, tempDir := createTestBuilder(t)

	t.Run("invalid_js_fallback", func(t *testing.T) {
		b.site.MinifyJS = true
		srcPath := filepath.Join(tempDir, "invalid.js")
		dstPath := filepath.Join(b.site.OutputDir, "invalid.js")

		invalidJS := `function broken( { console.log('broken'); }`
		require.NoError(t, os.WriteFile(srcPath, []byte(invalidJS), 0644))

		err := b.copyJSWithMinify(srcPath, dstPath)
		assert.NoError(t, err, "Should not error, should fallback to original")

		content, err := os.ReadFile(dstPath)
		assert.NoError(t, err)
		assert.Equal(t, invalidJS, string(content), "Should write original content on minify failure")
	})

	t.Run("invalid_json_fallback", func(t *testing.T) {
		b.site.MinifyJSON = true
		srcPath := filepath.Join(tempDir, "invalid.json")
		dstPath := filepath.Join(b.site.OutputDir, "invalid.json")

		invalidJSON := `{"name": "test", broken}`
		require.NoError(t, os.WriteFile(srcPath, []byte(invalidJSON), 0644))

		err := b.copyJSONWithMinify(srcPath, dstPath)
		assert.NoError(t, err, "Should not error, should fallback to original")

		content, err := os.ReadFile(dstPath)
		assert.NoError(t, err)
		assert.Equal(t, invalidJSON, string(content), "Should write original content on minify failure")
	})
}

func TestCopyFileWithoutMinify(t *testing.T) {
	b, tempDir := createTestBuilder(t)

	srcPath := filepath.Join(tempDir, "test.txt")
	dstPath := filepath.Join(b.site.OutputDir, "test.txt")

	content := "This is a text file\nwith multiple lines"
	require.NoError(t, os.WriteFile(srcPath, []byte(content), 0644))

	err := copyFile(srcPath, dstPath)
	assert.NoError(t, err)
	assert.FileExists(t, dstPath)

	readContent, err := os.ReadFile(dstPath)
	assert.NoError(t, err)
	assert.Equal(t, content, string(readContent), "Content should be unchanged")
}
