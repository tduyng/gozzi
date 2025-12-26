package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewI18n(t *testing.T) {
	i18n := NewI18n("en", "/tmp/data")

	if i18n.DefaultLang != "en" {
		t.Errorf("expected default lang 'en', got %s", i18n.DefaultLang)
	}

	if i18n.dataDir != "/tmp/data" {
		t.Errorf("expected dataDir '/tmp/data', got %s", i18n.dataDir)
	}

	if len(i18n.Languages) != 0 {
		t.Errorf("expected 0 languages, got %d", len(i18n.Languages))
	}
}

func TestNewI18n_EmptyDefaultLang(t *testing.T) {
	i18n := NewI18n("", "/tmp/data")

	if i18n.DefaultLang != "en" {
		t.Errorf("expected default lang to fall back to 'en', got %s", i18n.DefaultLang)
	}
}

func TestAddLanguage(t *testing.T) {
	i18n := NewI18n("en", "")

	i18n.AddLanguage("en", "English", 1, true)
	i18n.AddLanguage("fr", "Français", 2, false)

	if len(i18n.Languages) != 2 {
		t.Fatalf("expected 2 languages, got %d", len(i18n.Languages))
	}

	en := i18n.Languages["en"]
	if en.Code != "en" || en.Name != "English" || en.Weight != 1 || !en.IsDefault {
		t.Errorf("English language not configured correctly: %+v", en)
	}

	fr := i18n.Languages["fr"]
	if fr.Code != "fr" || fr.Name != "Français" || fr.Weight != 2 || fr.IsDefault {
		t.Errorf("French language not configured correctly: %+v", fr)
	}

	if i18n.DefaultLang != "en" {
		t.Errorf("expected default lang 'en', got %s", i18n.DefaultLang)
	}
}

func TestGetLanguage(t *testing.T) {
	i18n := NewI18n("en", "")
	i18n.AddLanguage("en", "English", 1, true)
	i18n.AddLanguage("fr", "Français", 2, false)

	en := i18n.GetLanguage("en")
	if en == nil || en.Code != "en" {
		t.Error("failed to get English language")
	}

	fr := i18n.GetLanguage("fr")
	if fr == nil || fr.Code != "fr" {
		t.Error("failed to get French language")
	}

	invalid := i18n.GetLanguage("de")
	if invalid != nil {
		t.Error("expected nil for non-existent language")
	}
}

func TestGetDefaultLanguage(t *testing.T) {
	i18n := NewI18n("en", "")
	i18n.AddLanguage("en", "English", 1, true)
	i18n.AddLanguage("fr", "Français", 2, false)

	defaultLang := i18n.GetDefaultLanguage()
	if defaultLang == nil || defaultLang.Code != "en" {
		t.Error("failed to get default language")
	}
}

func TestGetSortedLanguages(t *testing.T) {
	i18n := NewI18n("en", "")
	i18n.AddLanguage("fr", "Français", 3, false)
	i18n.AddLanguage("en", "English", 1, true)
	i18n.AddLanguage("es", "Español", 2, false)

	sorted := i18n.GetSortedLanguages()

	if len(sorted) != 3 {
		t.Fatalf("expected 3 languages, got %d", len(sorted))
	}

	// Should be sorted by weight: en (1), es (2), fr (3)
	if sorted[0].Code != "en" {
		t.Errorf("expected first language to be 'en', got %s", sorted[0].Code)
	}
	if sorted[1].Code != "es" {
		t.Errorf("expected second language to be 'es', got %s", sorted[1].Code)
	}
	if sorted[2].Code != "fr" {
		t.Errorf("expected third language to be 'fr', got %s", sorted[2].Code)
	}
}

func TestIsEnabled(t *testing.T) {
	i18n := NewI18n("en", "")

	if i18n.IsEnabled() {
		t.Error("expected i18n to be disabled with 0 languages")
	}

	i18n.AddLanguage("en", "English", 1, true)

	if i18n.IsEnabled() {
		t.Error("expected i18n to be disabled with 1 language")
	}

	i18n.AddLanguage("fr", "Français", 2, false)

	if !i18n.IsEnabled() {
		t.Error("expected i18n to be enabled with 2 languages")
	}
}

func TestSplitKey(t *testing.T) {
	tests := []struct {
		key      string
		expected []string
	}{
		{"simple", []string{"simple"}},
		{"section.contact", []string{"section", "contact"}},
		{"career.title", []string{"career", "title"}},
		{"a.b.c.d", []string{"a", "b", "c", "d"}},
		{"", []string{}},
	}

	for _, tt := range tests {
		result := splitKey(tt.key)
		if len(result) != len(tt.expected) {
			t.Errorf("splitKey(%q): expected %d parts, got %d", tt.key, len(tt.expected), len(result))
			continue
		}

		for i, part := range result {
			if part != tt.expected[i] {
				t.Errorf("splitKey(%q)[%d]: expected %q, got %q", tt.key, i, tt.expected[i], part)
			}
		}
	}
}

func TestGetNestedValue(t *testing.T) {
	i18n := NewI18n("en", "")

	translations := map[string]any{
		"simple": "Simple Value",
		"section": map[string]any{
			"contact": "Contact",
			"about":   "About",
		},
		"career": map[string]any{
			"title": "Career Profile",
		},
	}

	tests := []struct {
		key      string
		expected string
	}{
		{"simple", "Simple Value"},
		{"section.contact", "Contact"},
		{"section.about", "About"},
		{"career.title", "Career Profile"},
		{"nonexistent", ""},
		{"section.nonexistent", ""},
		{"", ""},
	}

	for _, tt := range tests {
		result := i18n.getNestedValue(translations, tt.key)
		if result != tt.expected {
			t.Errorf("getNestedValue(%q): expected %q, got %q", tt.key, tt.expected, result)
		}
	}
}

func TestTranslate(t *testing.T) {
	i18n := NewI18n("en", "")
	i18n.AddLanguage("en", "English", 1, true)
	i18n.AddLanguage("fr", "Français", 2, false)

	// Set up translations
	i18n.Languages["en"].Translations = map[string]any{
		"hello": "Hello",
		"section": map[string]any{
			"contact": "Contact",
		},
	}

	i18n.Languages["fr"].Translations = map[string]any{
		"hello": "Bonjour",
		"section": map[string]any{
			"contact": "Contact",
		},
	}

	// Test English translations
	result, err := i18n.Translate("en", "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != "Hello" {
		t.Errorf("expected 'Hello', got %q", result)
	}

	result, err = i18n.Translate("en", "section.contact")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != "Contact" {
		t.Errorf("expected 'Contact', got %q", result)
	}

	// Test French translations
	result, err = i18n.Translate("fr", "hello")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != "Bonjour" {
		t.Errorf("expected 'Bonjour', got %q", result)
	}

	// Test non-existent key
	_, err = i18n.Translate("en", "nonexistent")
	if err == nil {
		t.Error("expected error for non-existent key")
	}

	// Test non-existent language
	_, err = i18n.Translate("de", "hello")
	if err == nil {
		t.Error("expected error for non-existent language")
	}
}

func TestLoadTranslations(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	i18nDir := filepath.Join(tmpDir, "i18n")
	os.MkdirAll(i18nDir, 0755)

	// Create translation files
	enFile := filepath.Join(i18nDir, "en.toml")
	os.WriteFile(enFile, []byte(`
hello = "Hello"

[section]
contact = "Contact"
`), 0644)

	frFile := filepath.Join(i18nDir, "fr.toml")
	os.WriteFile(frFile, []byte(`
hello = "Bonjour"

[section]
contact = "Contact"
`), 0644)

	// Create i18n instance
	i18n := NewI18n("en", tmpDir)
	i18n.AddLanguage("en", "English", 1, true)
	i18n.AddLanguage("fr", "Français", 2, false)

	// Load translations
	err := i18n.LoadTranslations()
	if err != nil {
		t.Fatalf("failed to load translations: %v", err)
	}

	// Verify English translations
	if val := i18n.getNestedValue(i18n.Languages["en"].Translations, "hello"); val != "Hello" {
		t.Errorf("expected English 'hello' to be 'Hello', got %q", val)
	}

	// Verify French translations
	if val := i18n.getNestedValue(i18n.Languages["fr"].Translations, "hello"); val != "Bonjour" {
		t.Errorf("expected French 'hello' to be 'Bonjour', got %q", val)
	}
}

func TestLoadTranslations_AlternativeDirectory(t *testing.T) {
	// Test backward compatibility with data/translations/ instead of data/i18n/
	tmpDir := t.TempDir()
	translationsDir := filepath.Join(tmpDir, "translations")
	os.MkdirAll(translationsDir, 0755)

	enFile := filepath.Join(translationsDir, "en.toml")
	os.WriteFile(enFile, []byte(`hello = "Hello"`), 0644)

	i18n := NewI18n("en", tmpDir)
	i18n.AddLanguage("en", "English", 1, true)

	err := i18n.LoadTranslations()
	if err != nil {
		t.Fatalf("failed to load translations from alternative directory: %v", err)
	}

	if val := i18n.getNestedValue(i18n.Languages["en"].Translations, "hello"); val != "Hello" {
		t.Errorf("expected 'Hello', got %q", val)
	}
}

func TestLoadTranslations_NoDirectory(t *testing.T) {
	// Should not error if no translation directory exists
	tmpDir := t.TempDir()

	i18n := NewI18n("en", tmpDir)
	i18n.AddLanguage("en", "English", 1, true)

	err := i18n.LoadTranslations()
	if err != nil {
		t.Errorf("unexpected error when no translation directory exists: %v", err)
	}
}
