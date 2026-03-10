package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/builder"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
)

func TestI18nMultilingualSite(t *testing.T) {
	t.Parallel()
	projectDir := filepath.Join("testdata", "i18n-site")

	// Load site config
	site, err := config.LoadSite(filepath.Join(projectDir, "config.toml"))
	require.NoError(t, err)

	site.ProjectDir = projectDir

	// Load data files (including translations)
	err = site.LoadDataFiles(projectDir)
	require.NoError(t, err)

	// Set output dir to a temp directory
	tempDir := t.TempDir()
	site.OutputDir = tempDir

	// Load i18n translations
	require.NotNil(t, site.I18n, "I18n should be initialized")
	err = site.I18n.LoadTranslations()
	require.NoError(t, err)

	// Verify i18n configuration
	assert.True(t, site.I18n.IsEnabled(), "I18n should be enabled with 2 languages")
	assert.Equal(t, "en", site.I18n.GetDefaultLanguage().Code)

	langs := site.I18n.GetSortedLanguages()
	require.Len(t, langs, 2)
	assert.Equal(t, "en", langs[0].Code)
	assert.Equal(t, "English", langs[0].Name)
	assert.Equal(t, "fr", langs[1].Code)
	assert.Equal(t, "Français", langs[1].Name)

	// Test translations loaded correctly
	enHello, err := site.I18n.Translate("en", "greeting.hello")
	require.NoError(t, err)
	assert.Equal(t, "Hello", enHello)

	frHello, err := site.I18n.Translate("fr", "greeting.hello")
	require.NoError(t, err)
	assert.Equal(t, "Bonjour", frHello)

	// Parse content
	contentParser := parser.NewParser(site)
	gen, err := builder.NewBuilder(site, contentParser)
	require.NoError(t, err)

	err = contentParser.Parse(filepath.Join(projectDir, "content"))
	require.NoError(t, err)

	// Verify English section
	enSection, exists := contentParser.ContentMap["en"]
	require.True(t, exists, "English section should exist")
	assert.Equal(t, "en", enSection.Config["lang"], "English section should have lang=en")
	assert.Equal(t, "Home", enSection.Config["title"])

	// Verify French section
	frSection, exists := contentParser.ContentMap["fr"]
	require.True(t, exists, "French section should exist")
	assert.Equal(t, "fr", frSection.Config["lang"], "French section should have lang=fr")
	assert.Equal(t, "Accueil", frSection.Config["title"])

	// Verify English about page
	require.NotEmpty(t, enSection.Children)
	var enAbout *content.Node
	for _, child := range enSection.Children {
		if strings.Contains(child.Path, "about") {
			enAbout = child
			break
		}
	}
	require.NotNil(t, enAbout, "English about page should exist")
	assert.Equal(t, "en", enAbout.Config["lang"], "English about page should have lang=en")
	assert.Equal(t, "About", enAbout.Config["title"])

	// Verify French about page
	require.NotEmpty(t, frSection.Children)
	var frAbout *content.Node
	for _, child := range frSection.Children {
		if strings.Contains(child.Path, "about") {
			frAbout = child
			break
		}
	}
	require.NotNil(t, frAbout, "French about page should exist")
	assert.Equal(t, "fr", frAbout.Config["lang"], "French about page should have lang=fr")
	assert.Equal(t, "À Propos", frAbout.Config["title"])

	// Build the site
	err = gen.Generate(contentParser.ContentMap["."])
	require.NoError(t, err)

	// Verify generated HTML files
	enIndexPath := filepath.Join(site.OutputDir, "en", "index.html")
	assert.FileExists(t, enIndexPath, "English index.html should be generated")

	frIndexPath := filepath.Join(site.OutputDir, "fr", "index.html")
	assert.FileExists(t, frIndexPath, "French index.html should be generated")

	// Read and verify English content
	enContent, err := os.ReadFile(enIndexPath)
	require.NoError(t, err)
	enHTML := string(enContent)

	assert.Contains(t, enHTML, `lang="en"`, "English page should have lang=en")
	assert.Contains(t, enHTML, "Welcome to our site", "English page should use English translation")
	assert.Contains(t, enHTML, "Welcome |", "English page should show English greeting")
	assert.Contains(t, enHTML, "Language: en", "English page should show language code")
	assert.Contains(t, enHTML, `<a href="/en/">English</a>`, "Should have English language link")
	assert.Contains(t, enHTML, `<a href="/fr/">Français</a>`, "Should have French language link")

	// Read and verify French content
	frContent, err := os.ReadFile(frIndexPath)
	require.NoError(t, err)
	frHTML := string(frContent)

	assert.Contains(t, frHTML, `lang="fr"`, "French page should have lang=fr")
	assert.Contains(t, frHTML, "Bienvenue sur notre site", "French page should use French translation")
	assert.Contains(t, frHTML, "Bienvenue |", "French page should show French greeting")
	assert.Contains(t, frHTML, "Language: fr", "French page should show language code")

	// Verify about pages
	enAboutPath := filepath.Join(site.OutputDir, "en", "about", "index.html")
	assert.FileExists(t, enAboutPath, "English about page should be generated")

	frAboutPath := filepath.Join(site.OutputDir, "fr", "about", "index.html")
	assert.FileExists(t, frAboutPath, "French about page should be generated")

	// Read and verify English about page
	enAboutContent, err := os.ReadFile(enAboutPath)
	require.NoError(t, err)
	enAboutHTML := string(enAboutContent)

	assert.Contains(t, enAboutHTML, "About Us", "English about should use English translation")
	assert.Contains(t, enAboutHTML, "Hello |", "English about should show English greeting")
	assert.Contains(t, enAboutHTML, "About us in English", "English about should have English content")

	// Read and verify French about page
	frAboutContent, err := os.ReadFile(frAboutPath)
	require.NoError(t, err)
	frAboutHTML := string(frAboutContent)

	assert.Contains(t, frAboutHTML, "À Propos", "French about should use French translation")
	assert.Contains(t, frAboutHTML, "Bonjour |", "French about should show French greeting")
	assert.Contains(t, frAboutHTML, "À propos de nous en français", "French about should have French content")

	// Clean up
	os.RemoveAll(site.OutputDir)
}
