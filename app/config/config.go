// Package config handles TOML configuration loading and merging.
package config

import (
	"bytes"
	"maps"
	"path/filepath"
	"reflect"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/tduyng/gozzi/app/i18n"
	"github.com/tduyng/gozzi/app/utils"
)

// Site represents the site-wide configuration loaded from config.toml.
type Site struct {
	BaseURL         string           `toml:"base_url"`
	FeedURL         string           `toml:"feed_url"`
	Description     string           `toml:"Description"`
	Extra           map[string]any   `toml:"extra"`
	GenerateFeed    bool             `toml:"generate_feed"`
	Theme           string           `toml:"theme"`
	Img             string           `toml:"img"`
	Lang            string           `toml:"language"`
	OutputDir       string           `toml:"output_dir"`
	Title           string           `toml:"title"`
	StrictTemplates bool             `toml:"strict_templates"`
	SyntaxTheme     string           `toml:"syntax_theme"`
	MinifyCSS       bool             `toml:"minify_css"`
	MinifyHTML      bool             `toml:"minify_html"`
	MinifyJS        bool             `toml:"minify_js"`
	MinifyJSON      bool             `toml:"minify_json"`
	MinifySVG       bool             `toml:"minify_svg"`
	MinifyXML       bool             `toml:"minify_xml"`
	Taxonomies      TaxonomiesConfig `toml:"taxonomies"`
	BuildTime       time.Time
	BuildDrafts     bool
	Data            map[string]any `toml:"-"` // Loaded from data/ directory, not from TOML
	I18n            *i18n.I18n     `toml:"-"` // Loaded from config languages section and data/ directory
}

// FrontMatter represents the TOML front matter in markdown content files.
type FrontMatter struct {
	Date         time.Time      `toml:"date"`
	Description  string         `toml:"description"`
	Draft        bool           `toml:"draft"`
	Extra        map[string]any `toml:"extra"`
	GenerateFeed bool           `toml:"generate_feed"`
	Img          string         `toml:"img"`
	Lang         string         `toml:"language"`
	Tags         []string       `toml:"tags"`
	Categories   []string       `toml:"categories"`
	Series       string         `toml:"series"`
	SeriesOrder  int            `toml:"series_order"`
	Template     string         `toml:"template"`
	Title        string         `toml:"title"`
	Featured     bool           `toml:"featured"`
	Updated      time.Time      `toml:"updated"`
}

// LoadSite loads site configuration from a TOML file.
// It also attempts to load data files from the data/ directory if it exists.
func LoadSite(path string) (*Site, error) {
	var cfg Site
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, utils.WrapWithContext(utils.ErrConfig, err, utils.ErrorContext{
			Operation: "load_site_config",
			Component: "config",
			Path:      path,
		})
	}

	// Initialize default taxonomies if not configured
	if cfg.Taxonomies == nil {
		cfg.Taxonomies = DefaultTaxonomiesConfig()
	} else {
		// Merge with defaults for any missing taxonomies
		defaults := DefaultTaxonomiesConfig()
		for name, defaultCfg := range defaults {
			if _, exists := cfg.Taxonomies[name]; !exists {
				cfg.Taxonomies[name] = defaultCfg
			}
		}
	}

	// Load i18n configuration from config file
	if err := cfg.loadI18nConfig(path); err != nil {
		return nil, utils.WrapWithContext(utils.ErrConfig, err, utils.ErrorContext{
			Operation: "load_i18n_config",
			Component: "config",
			Path:      path,
		})
	}

	return &cfg, nil
}

// LoadDataFiles loads data from the data/ directory.
// This is exposed to allow the builder to load data files without circular dependencies.
func (s *Site) LoadDataFiles(projectDir string) error {
	dataDir := filepath.Join(projectDir, "data")

	// Import data loader dynamically to avoid circular imports
	// We'll use a function variable pattern
	if dataLoaderFunc != nil {
		data, err := dataLoaderFunc(dataDir)
		if err != nil {
			return err
		}
		s.Data = data
	} else {
		// No data loader registered, initialize empty
		s.Data = make(map[string]any)
	}

	return nil
}

// loadI18nConfig loads internationalization configuration from the config file.
func (s *Site) loadI18nConfig(configPath string) error {
	// Parse the raw TOML to extract languages section
	var rawConfig struct {
		Languages map[string]map[string]any `toml:"languages"`
	}

	if _, err := toml.DecodeFile(configPath, &rawConfig); err != nil {
		// If we can't decode, languages section probably doesn't exist, which is fine
		return nil
	}

	// If no languages configured, i18n is disabled
	if len(rawConfig.Languages) == 0 {
		return nil
	}

	// Find default language
	defaultLang := ""
	if defLangRaw, exists := rawConfig.Languages["default"]; exists {
		if defLangStr, ok := defLangRaw[""].(string); ok {
			defaultLang = defLangStr
		}
	}

	// If no default specified, check if site.Lang is set
	if defaultLang == "" && s.Lang != "" {
		defaultLang = s.Lang
	}

	// If still no default, use "en"
	if defaultLang == "" {
		defaultLang = "en"
	}

	// Get project directory from config path
	projectDir := filepath.Dir(configPath)
	dataDir := filepath.Join(projectDir, "data")

	// Create i18n manager
	s.I18n = i18n.NewI18n(defaultLang, dataDir)

	// Add each language
	for langCode, langConfig := range rawConfig.Languages {
		// Skip the "default" key
		if langCode == "default" {
			continue
		}

		name := langCode // Default name is the code
		if nameRaw, exists := langConfig["name"]; exists {
			if nameStr, ok := nameRaw.(string); ok {
				name = nameStr
			}
		}

		weight := 0
		if weightRaw, exists := langConfig["weight"]; exists {
			switch v := weightRaw.(type) {
			case int:
				weight = v
			case int64:
				weight = int(v)
			case float64:
				weight = int(v)
			}
		}

		isDefault := langCode == defaultLang
		s.I18n.AddLanguage(langCode, name, weight, isDefault)
	}

	return nil
}

// DataLoaderFunc is a function type for loading data files.
type DataLoaderFunc func(dataDir string) (map[string]any, error)

// dataLoaderFunc is the registered data loader (set by data package init).
var dataLoaderFunc DataLoaderFunc

// RegisterDataLoader registers the data loading function.
// This is called by the data package to avoid circular dependencies.
func RegisterDataLoader(loader DataLoaderFunc) {
	dataLoaderFunc = loader
}

// LoadFrontMatter parses TOML front matter from markdown content.
func LoadFrontMatter(content []byte) (*FrontMatter, []byte, error) {
	return parseFrontMatter[FrontMatter](content)
}

// ToConfig converts Site to a map for template rendering.
func (site *Site) ToConfig() map[string]any {
	siteConfig := make(map[string]any)
	if site != nil {
		siteConfig["base_url"] = site.BaseURL
		siteConfig["feed_url"] = site.FeedURL
		siteConfig["title"] = site.Title
		siteConfig["description"] = site.Description
		siteConfig["output_dir"] = site.OutputDir
		siteConfig["lang"] = site.Lang
		siteConfig["generate_feed"] = site.GenerateFeed
		siteConfig["strict_templates"] = site.StrictTemplates
		siteConfig["build_time"] = site.BuildTime
		siteConfig["extra"] = MergeExtra(make(map[string]any), site.Extra)
		siteConfig["data"] = site.Data
	}
	if site.OutputDir == "" {
		siteConfig["output_dir"] = "public"
	}

	if site.Lang == "" {
		siteConfig["lang"] = "en"
	}

	if site.FeedURL == "" {
		siteConfig["feed_url"] = site.BaseURL + "/atom.xml"
	}

	return siteConfig
}

// ToConfig converts FrontMatter to a map for template rendering.
func (frontMatter *FrontMatter) ToConfig() map[string]any {
	config := make(map[string]any)
	config["title"] = frontMatter.Title
	config["description"] = frontMatter.Description
	config["template"] = frontMatter.Template
	config["generate_feed"] = frontMatter.GenerateFeed
	config["lang"] = frontMatter.Lang
	config["featured"] = frontMatter.Featured
	config["draft"] = frontMatter.Draft
	config["tags"] = frontMatter.Tags
	config["categories"] = frontMatter.Categories
	config["series"] = frontMatter.Series
	config["series_order"] = frontMatter.SeriesOrder
	config["date"] = frontMatter.Date
	config["updated"] = frontMatter.Updated
	config["extra"] = MergeExtra(make(map[string]any), frontMatter.Extra)
	return config
}

// MergeConfigs merges site, section, and page configurations with proper precedence.
func MergeConfigs(site, section, page map[string]any) map[string]any {
	merged := maps.Clone(site)

	if section != nil {
		merged = mergeFrontMatter(merged, section)
	}

	if page != nil {
		merged = mergeFrontMatter(merged, page)
	}

	return merged
}

func mergeFrontMatter(merged, frontMatter map[string]any) map[string]any {
	if frontMatter["title"] != "" {
		merged["title"] = frontMatter["title"]
	}
	if frontMatter["description"] != "" {
		merged["description"] = frontMatter["description"]
	}
	if frontMatter["template"] != "" {
		merged["template"] = frontMatter["template"]
	}

	if frontMatter["img"] != "" {
		merged["img"] = frontMatter["img"]
	}

	if frontMatter["lang"] != "" {
		merged["lang"] = frontMatter["lang"]
	}

	merged["tags"] = frontMatter["tags"]
	merged["categories"] = frontMatter["categories"]
	merged["series"] = frontMatter["series"]
	merged["series_order"] = frontMatter["series_order"]

	if date, ok := frontMatter["date"].(time.Time); ok && !date.IsZero() {
		merged["date"] = date
	}

	if updated, ok := frontMatter["updated"].(time.Time); ok && !updated.IsZero() {
		merged["updated"] = updated
	} else {
		merged["updated"] = time.Time{}
	}

	boolFields := []string{"generate_feed", "draft", "featured"}
	for _, field := range boolFields {
		if val, exists := frontMatter[field]; exists {
			if b, ok := val.(bool); ok {
				merged[field] = b
			}
		}
	}

	if frontMatterExtra, exists := frontMatter["extra"]; exists {
		if extraMap, ok := frontMatterExtra.(map[string]any); ok {
			if existingExtra, hasExtra := merged["extra"]; hasExtra {
				if existingExtraMap, ok := existingExtra.(map[string]any); ok {
					merged["extra"] = MergeExtra(existingExtraMap, extraMap)
				} else {
					merged["extra"] = extraMap
				}
			} else {
				merged["extra"] = extraMap
			}
		}
	}
	return merged
}

// MergeExtra recursively merges extra configuration maps.
func MergeExtra(base map[string]any, extra map[string]any) map[string]any {
	pageSpecificFields := map[string]bool{
		"img": true,
	}

	for k, v := range extra {
		if existing, exists := base[k]; exists {
			switch existingVal := existing.(type) {
			case map[string]any:
				if newMap, ok := v.(map[string]any); ok {
					base[k] = mergeMapsDeep(existingVal, newMap)
					continue
				}
			case []any:
				if newArr, ok := v.([]any); ok {
					base[k] = append(existingVal, newArr...)
					continue
				}
			}
		}
		base[k] = v
	}

	for field := range pageSpecificFields {
		if _, existsInExtra := extra[field]; !existsInExtra {
			delete(base, field)
		}
	}

	return base
}

// GetBool safely retrieves a boolean value from a configuration map.
func GetBool(m map[string]any, key string) bool {
	if val, exists := m[key]; exists {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

func mergeMapsDeep(base, override map[string]any) map[string]any {
	result := maps.Clone(base)

	for k, overrideVal := range override {
		if baseVal, exists := base[k]; exists {
			switch baseValTyped := baseVal.(type) {
			case map[string]any:
				if overrideMap, ok := overrideVal.(map[string]any); ok {
					result[k] = mergeMapsDeep(baseValTyped, overrideMap)
					continue
				}
			case []any:
				if overrideArr, ok := overrideVal.([]any); ok {
					result[k] = append(baseValTyped, overrideArr...)
					continue
				}
			}
		}
		result[k] = overrideVal
	}
	return result
}

func parseFrontMatter[T any](content []byte) (*T, []byte, error) {
	frontMatterDelim := []byte("+++")
	parts := bytes.SplitN(content, frontMatterDelim, 3)

	var (
		config *T
		body   []byte
		err    error
	)

	switch len(parts) {
	case 1:
		body = content
	case 2:
		body = bytes.Join(parts, frontMatterDelim)
	case 3:
		frontMatter := bytes.TrimSpace(parts[1])
		config = new(T)
		if err = toml.Unmarshal(frontMatter, config); err != nil {
			return nil, nil, utils.WrapWithContext(utils.ErrContent, err, utils.ErrorContext{
				Operation: "parse_frontmatter_toml",
				Component: "config",
			})
		}

		// Normalize dates to UTC for timezone-independent behavior
		normalizeDatesToUTC(config)

		body = bytes.TrimSpace(parts[2])
	}

	return config, body, nil
}

// normalizeDatesToUTC converts all time.Time values in a struct/map to UTC
// This ensures dates from TOML (which default to local timezone) are consistent
func normalizeDatesToUTC(v any) {
	if v == nil {
		return
	}

	// Handle pointers
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		// Iterate through struct fields
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if !field.CanSet() {
				continue
			}

			if field.Type() == reflect.TypeOf(time.Time{}) {
				// Convert time.Time field to UTC
				t := field.Interface().(time.Time)
				if !t.IsZero() {
					field.Set(reflect.ValueOf(t.UTC()))
				}
			} else if field.Kind() == reflect.Map || field.Kind() == reflect.Slice || field.Kind() == reflect.Struct {
				// Recursively normalize nested structures
				if field.CanInterface() {
					normalizeDatesToUTC(field.Addr().Interface())
				}
			}
		}

	case reflect.Map:
		// Iterate through map keys
		iter := val.MapRange()
		for iter.Next() {
			mapVal := iter.Value()
			if mapVal.Type() == reflect.TypeOf(time.Time{}) {
				t := mapVal.Interface().(time.Time)
				if !t.IsZero() {
					val.SetMapIndex(iter.Key(), reflect.ValueOf(t.UTC()))
				}
			} else if mapVal.Kind() == reflect.Map || mapVal.Kind() == reflect.Slice || mapVal.Kind() == reflect.Struct {
				// For nested structures, normalize recursively
				if mapVal.CanInterface() {
					normalizeDatesToUTC(mapVal.Interface())
				}
			}
		}

	case reflect.Slice:
		// Iterate through slice elements
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			if elem.CanInterface() {
				normalizeDatesToUTC(elem.Interface())
			}
		}
	}
}
