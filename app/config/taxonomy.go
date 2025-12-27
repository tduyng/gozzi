// Taxonomy configuration structures for site-wide and frontmatter taxonomy settings.
// Supports tags, categories, series, and custom taxonomies with flexible configuration.
package config

// TaxonomyConfig represents configuration for a single taxonomy type.
type TaxonomyConfig struct {
	Enabled    bool   `toml:"enabled"`     // Whether this taxonomy is active
	PaginateBy int    `toml:"paginate_by"` // Items per page (0 = no pagination)
	Template   string `toml:"template"`    // Custom template override (optional)
}

// TaxonomiesConfig holds all taxonomy configurations.
type TaxonomiesConfig map[string]TaxonomyConfig

// DefaultTaxonomiesConfig returns the default taxonomy configuration.
func DefaultTaxonomiesConfig() TaxonomiesConfig {
	return TaxonomiesConfig{
		"tags": {
			Enabled:    true,
			PaginateBy: 0, // No pagination by default
		},
		"categories": {
			Enabled:    false,
			PaginateBy: 0,
		},
		"series": {
			Enabled:    false,
			PaginateBy: 0,
		},
	}
}

// GetTaxonomyConfig safely retrieves taxonomy config with fallback to defaults.
func (tc TaxonomiesConfig) GetTaxonomyConfig(name string) TaxonomyConfig {
	if cfg, exists := tc[name]; exists {
		return cfg
	}
	// Return default config for unknown taxonomies
	return TaxonomyConfig{Enabled: false}
}

// IsEnabled checks if a taxonomy is enabled.
func (tc TaxonomiesConfig) IsEnabled(name string) bool {
	return tc.GetTaxonomyConfig(name).Enabled
}
