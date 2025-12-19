// Package config handles TOML configuration loading and merging.
package config

import (
	"bytes"
	"maps"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/tduyng/gozzi/app/utils"
)

// Site represents the site-wide configuration loaded from config.toml.
type Site struct {
	BaseURL         string         `toml:"base_url"`
	FeedURL         string         `toml:"feed_url"`
	Description     string         `toml:"Description"`
	Extra           map[string]any `toml:"extra"`
	GenerateFeed    bool           `toml:"generate_feed"`
	Theme           string         `toml:"theme"`
	Img             string         `toml:"img"`
	Lang            string         `toml:"language"`
	OutputDir       string         `toml:"output_dir"`
	Title           string         `toml:"title"`
	StrictTemplates bool           `toml:"strict_templates"`
	SyntaxTheme     string         `toml:"syntax_theme"`
	MinifyCSS       bool           `toml:"minify_css"`
	MinifyHTML      bool           `toml:"minify_html"`
	MinifyJS        bool           `toml:"minify_js"`
	MinifyJSON      bool           `toml:"minify_json"`
	MinifySVG       bool           `toml:"minify_svg"`
	MinifyXML       bool           `toml:"minify_xml"`
	BuildTime       time.Time
	BuildDrafts     bool
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
	Template     string         `toml:"template"`
	Title        string         `toml:"title"`
	Featured     bool           `toml:"featured"`
	Updated      time.Time      `toml:"updated"`
}

// LoadSite loads site configuration from a TOML file.
func LoadSite(path string) (*Site, error) {
	var cfg Site
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, utils.WrapWithContext(utils.ErrConfig, err, utils.ErrorContext{
			Operation: "load_site_config",
			Component: "config",
			Path:      path,
		})
	}
	return &cfg, nil
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

		body = bytes.TrimSpace(parts[2])
	}

	return config, body, nil
}
