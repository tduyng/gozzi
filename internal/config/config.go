package config

import (
	"bytes"
	"fmt"
	"maps"
	"time"

	"github.com/BurntSushi/toml"
)

type Site struct {
	BaseURL      string         `toml:"base_url"`
	FeedURL      string         `toml:"feed_url"`
	Description  string         `toml:"Description"`
	Extra        map[string]any `toml:"extra"`
	GenerateFeed bool           `toml:"generate_feed"`
	Theme        string         `toml:"theme"`
	Img          string         `toml:"img"`
	Lang         string         `toml:"language"`
	OutputDir    string         `toml:"output_dir"`
	Title        string         `toml:"title"`
}

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
	Update       time.Time      `toml:"update"`
}

func LoadSite(path string) (*Site, error) {
	var cfg Site
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadFrontMatter(content []byte) (*FrontMatter, []byte, error) {
	return parseFrontMatter[FrontMatter](content)
}

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
		siteConfig["img"] = site.Img
		siteConfig["extra"] = MergeExtra(make(map[string]any), site.Extra)
	}
	return siteConfig
}

func (frontMatter *FrontMatter) ToConfig() map[string]any {
	config := make(map[string]any)
	config["title"] = frontMatter.Title
	config["description"] = frontMatter.Description
	config["template"] = frontMatter.Template
	config["generate_feed"] = frontMatter.GenerateFeed
	config["img"] = frontMatter.Img
	config["lang"] = frontMatter.Lang
	config["featured"] = frontMatter.Featured
	config["draft"] = frontMatter.Draft
	config["tags"] = frontMatter.Tags
	config["date"] = frontMatter.Date
	config["update"] = frontMatter.Update
	config["extra"] = MergeExtra(config, frontMatter.Extra)
	return config
}

func MergeConfigs(site, section, page map[string]any) map[string]any {
	merged := make(map[string]any)
	maps.Copy(merged, site)

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

	if frontMatter["generate_feed"].(bool) {
		merged["generate_feed"] = frontMatter["generate_feed"]
	}

	if frontMatter["img"] != "" {
		merged["img"] = frontMatter["img"]
	}

	if frontMatter["lang"] != "" {
		merged["lang"] = frontMatter["lang"]
	}

	if frontMatter["featured"].(bool) {
		merged["featured"] = frontMatter["featured"]
	}

	if frontMatter["draft"].(bool) {
		merged["draft"] = frontMatter["draft"].(bool)
	}

	merged["tags"] = frontMatter["tags"]
	merged["date"] = frontMatter["date"]
	merged["update"] = frontMatter["update"]

	merged["extra"] = MergeExtra(merged, frontMatter["extra"].(map[string]any))
	return merged
}

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

func mergeMapsDeep(base, override map[string]any) map[string]any {
	result := make(map[string]any, len(base))
	maps.Copy(result, base)

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
	case 1: // No front matter
		body = content
	case 2: // Only opening +++
		body = bytes.Join(parts, frontMatterDelim)
	case 3:
		frontMatter := bytes.TrimSpace(parts[1])
		config = new(T)
		if err = toml.Unmarshal(frontMatter, config); err != nil {
			return nil, nil, fmt.Errorf("front matter parsing failed: %w", err)
		}

		body = bytes.TrimSpace(parts[2])
	}

	return config, body, nil
}
