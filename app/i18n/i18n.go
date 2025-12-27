// Core i18n package for multilingual site support
// Handles language configs, translation loading, and locale management
package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
)

// Language represents a configured language
type Language struct {
	Code         string         // Language code (e.g., "en", "fr")
	Name         string         // Display name (e.g., "English", "Français")
	Weight       int            // Sort order (lower = higher priority)
	IsDefault    bool           // Whether this is the default language
	Translations map[string]any // Loaded translations for this language
}

// I18n manages multilingual configuration and translations
type I18n struct {
	Languages   map[string]*Language // All configured languages (key = language code)
	DefaultLang string               // Default language code
	dataDir     string               // Path to data directory
}

// Config represents the i18n configuration from config.toml
type Config struct {
	DefaultLang string                    `toml:"default_language"`
	Languages   map[string]LanguageConfig `toml:"languages"`
}

// LanguageConfig represents individual language configuration
type LanguageConfig struct {
	Name   string `toml:"name"`
	Weight int    `toml:"weight"`
}

// NewI18n creates a new i18n manager
func NewI18n(defaultLang string, dataDir string) *I18n {
	if defaultLang == "" {
		defaultLang = "en"
	}

	return &I18n{
		Languages:   make(map[string]*Language),
		DefaultLang: defaultLang,
		dataDir:     dataDir,
	}
}

// AddLanguage adds a new language to the i18n system
func (i *I18n) AddLanguage(code, name string, weight int, isDefault bool) {
	i.Languages[code] = &Language{
		Code:         code,
		Name:         name,
		Weight:       weight,
		IsDefault:    isDefault,
		Translations: make(map[string]any),
	}

	if isDefault {
		i.DefaultLang = code
	}
}

// LoadTranslations loads translation files for all configured languages
func (i *I18n) LoadTranslations() error {
	if i.dataDir == "" {
		return nil // No data directory configured
	}

	i18nDir := filepath.Join(i.dataDir, "i18n")
	if _, err := os.Stat(i18nDir); os.IsNotExist(err) {
		// Try alternative: data/translations/ (for backward compatibility)
		i18nDir = filepath.Join(i.dataDir, "translations")
		if _, err := os.Stat(i18nDir); os.IsNotExist(err) {
			return nil // No translation directory exists
		}
	}

	for code, lang := range i.Languages {
		translationFile := filepath.Join(i18nDir, code+".toml")
		if _, err := os.Stat(translationFile); os.IsNotExist(err) {
			continue // Translation file doesn't exist for this language
		}

		var translations map[string]any
		if _, err := toml.DecodeFile(translationFile, &translations); err != nil {
			return fmt.Errorf("failed to load translations for %s: %w", code, err)
		}

		lang.Translations = translations
	}

	return nil
}

// GetLanguage returns a language by code
func (i *I18n) GetLanguage(code string) *Language {
	return i.Languages[code]
}

// GetDefaultLanguage returns the default language
func (i *I18n) GetDefaultLanguage() *Language {
	return i.Languages[i.DefaultLang]
}

// GetSortedLanguages returns all languages sorted by weight
func (i *I18n) GetSortedLanguages() []*Language {
	langs := make([]*Language, 0, len(i.Languages))
	for _, lang := range i.Languages {
		langs = append(langs, lang)
	}

	sort.Slice(langs, func(i, j int) bool {
		return langs[i].Weight < langs[j].Weight
	})

	return langs
}

// IsEnabled returns whether i18n is enabled (more than one language)
func (i *I18n) IsEnabled() bool {
	return len(i.Languages) > 1
}

// Translate looks up a translation key for a specific language
// Key format: "section.contact" or "career.title"
// Falls back to default language if not found in requested language
func (i *I18n) Translate(langCode, key string) (string, error) {
	lang := i.GetLanguage(langCode)
	if lang == nil {
		return "", fmt.Errorf("language %s not found", langCode)
	}

	// Try to get translation from requested language
	if value := i.getNestedValue(lang.Translations, key); value != "" {
		return value, nil
	}

	// Fall back to default language
	if langCode != i.DefaultLang {
		defaultLang := i.GetDefaultLanguage()
		if defaultLang != nil {
			if value := i.getNestedValue(defaultLang.Translations, key); value != "" {
				return value, nil
			}
		}
	}

	return "", fmt.Errorf("translation key %s not found for language %s", key, langCode)
}

// getNestedValue retrieves a nested value from translation map
// Supports keys like "section.contact" -> translations["section"]["contact"]
func (i *I18n) getNestedValue(m map[string]any, key string) string {
	// Split key by dots
	keys := splitKey(key)

	current := m
	for i, k := range keys {
		if i == len(keys)-1 {
			// Last key - return the value
			if val, ok := current[k]; ok {
				if str, ok := val.(string); ok {
					return str
				}
			}
			return ""
		}

		// Intermediate key - navigate deeper
		if val, ok := current[k]; ok {
			if nested, ok := val.(map[string]any); ok {
				current = nested
			} else {
				return ""
			}
		} else {
			return ""
		}
	}

	return ""
}

// splitKey splits a translation key by dots
func splitKey(key string) []string {
	var parts []string
	var current string

	for _, ch := range key {
		if ch == '.' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}

// GetAllTranslations returns all translations for a language
func (i *I18n) GetAllTranslations(langCode string) map[string]any {
	lang := i.GetLanguage(langCode)
	if lang == nil {
		return make(map[string]any)
	}
	return lang.Translations
}
