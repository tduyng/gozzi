// This file tests static asset handling and copying behavior.
package integration

import (
	"os"
	"path/filepath"
	"testing"
)

// TestStatic_Copying tests static file copying
func TestStatic_Copying(t *testing.T) {
	t.Run("AllStatic_CopiedToOutput", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		staticFiles := []string{"robots.txt", "style.css"}
		for _, file := range staticFiles {
			verifyFileExists(t, sitePath, file)
		}
	})

	t.Run("StaticContent_Verbatim", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileContent(t, sitePath, "style.css", "font-family")
		verifyFileContent(t, sitePath, "robots.txt", "User-agent")
	})
}

// TestStatic_NestedStructure tests nested static directories
func TestStatic_NestedStructure(t *testing.T) {
	t.Run("NestedDirs_PreserveStructure", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Create nested structure
		nestedDir := filepath.Join(sitePath, "static/assets/images")
		os.MkdirAll(nestedDir, 0755)
		os.WriteFile(filepath.Join(nestedDir, "logo.svg"), []byte("<svg></svg>"), 0644)

		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "assets/images/logo.svg")
		verifyFileContent(t, sitePath, "assets/images/logo.svg", "<svg>")
	})

	t.Run("DeepNesting_Works", func(t *testing.T) {
		sitePath := setupTestSite(t)

		deepDir := filepath.Join(sitePath, "static/a/b/c/d")
		os.MkdirAll(deepDir, 0755)
		os.WriteFile(filepath.Join(deepDir, "deep.txt"), []byte("deep"), 0644)

		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "a/b/c/d/deep.txt")
	})
}

// TestStatic_Changes tests static file modification
func TestStatic_Changes(t *testing.T) {
	t.Run("StaticChange_UpdatesFile", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Modify static file
		staticPath := filepath.Join(sitePath, "static/style.css")
		modifyFile(t, staticPath, "font-family", "font-family: monospace")

		// Rebuild
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileContent(t, sitePath, "style.css", "monospace")
	})

	t.Run("StaticChange_NoContentRebuild", func(t *testing.T) {
		t.Skip("TODO: Fix cache to preserve hit rate when static files change")
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Get baseline
		gen.ResetCacheStats()
		fullRebuild(t, gen, contentParser, sitePath)
		baselineHitRate := gen.GetCacheStats().HitRate

		// Modify static file
		staticPath := filepath.Join(sitePath, "static/style.css")
		modifyFile(t, staticPath, "font-family", "font-family: sans-serif")

		// Rebuild
		gen.ResetCacheStats()
		fullRebuild(t, gen, contentParser, sitePath)
		afterHitRate := gen.GetCacheStats().HitRate

		// Cache hit rate should be similar (static changes don't affect content)
		if afterHitRate < baselineHitRate-10 {
			t.Errorf("static change should not significantly affect cache (before: %.1f%%, after: %.1f%%)",
				baselineHitRate, afterHitRate)
		}
	})
}

// TestStatic_NewFiles tests adding new static files
func TestStatic_NewFiles(t *testing.T) {
	t.Run("NewStatic_Added", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Add new static file
		newFile := filepath.Join(sitePath, "static/new-asset.txt")
		os.WriteFile(newFile, []byte("new content"), 0644)

		// Rebuild
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "new-asset.txt")
		verifyFileContent(t, sitePath, "new-asset.txt", "new content")
	})

	t.Run("NewStatic_NoContentRebuild", func(t *testing.T) {
		t.Skip("TODO: Fix cache to preserve hit rate when new static files are added")
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Baseline
		gen.ResetCacheStats()
		fullRebuild(t, gen, contentParser, sitePath)
		baseline := gen.GetCacheStats().HitRate

		// Add static file
		newFile := filepath.Join(sitePath, "static/another.js")
		os.WriteFile(newFile, []byte("console.log('hi')"), 0644)

		// Rebuild
		gen.ResetCacheStats()
		fullRebuild(t, gen, contentParser, sitePath)
		after := gen.GetCacheStats().HitRate

		// Should not affect content cache
		if after < baseline-10 {
			t.Errorf("new static file should not affect content cache (before: %.1f%%, after: %.1f%%)",
				baseline, after)
		}
	})
}

// TestStatic_AuxiliaryGeneration tests generated auxiliary files
func TestStatic_AuxiliaryGeneration(t *testing.T) {
	t.Run("RobotsTxt_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "robots.txt")
	})

	t.Run("Sitemap_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "sitemap.xml")
		verifyFileContent(t, sitePath, "sitemap.xml", "<?xml")
	})

	t.Run("AtomFeed_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "atom.xml")
		verifyFileContent(t, sitePath, "atom.xml", "<?xml")
	})
}

// TestStatic_CustomStaticDir tests custom static directories
func TestStatic_CustomStaticDir(t *testing.T) {
	t.Run("MultipleStaticDirs_AllCopied", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Create files in different subdirs
		os.MkdirAll(filepath.Join(sitePath, "static/css"), 0755)
		os.MkdirAll(filepath.Join(sitePath, "static/js"), 0755)
		os.MkdirAll(filepath.Join(sitePath, "static/images"), 0755)

		os.WriteFile(filepath.Join(sitePath, "static/css/app.css"), []byte("body{}"), 0644)
		os.WriteFile(filepath.Join(sitePath, "static/js/app.js"), []byte("console.log()"), 0644)
		os.WriteFile(filepath.Join(sitePath, "static/images/logo.png"), []byte("PNG"), 0644)

		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "css/app.css")
		verifyFileExists(t, sitePath, "js/app.js")
		verifyFileExists(t, sitePath, "images/logo.png")
	})
}
