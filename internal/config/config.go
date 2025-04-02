package config

import (
	"bytes"
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

func LoadFrontMatter(content []byte) (*FrontMatter, error) {
	return parseFrontMatter[FrontMatter](content)
}

func MergeConfigs(site *Site, section *FrontMatter, page *FrontMatter) map[string]any {
	merged := make(map[string]any)
	if site != nil {
		merged["base_url"] = site.BaseURL
		merged["feed_url"] = site.FeedURL
		merged["title"] = site.Title
		merged["description"] = site.Description
		merged["output_dir"] = site.OutputDir
		merged["lang"] = site.Lang
		merged["generate_feed"] = site.GenerateFeed
		merged["img"] = site.Img
		merged["extra"] = MergeExtra(make(map[string]any), site.Extra)
	}

	if section != nil {
		merged = mergeFrontMatter(merged, section)
	}

	if page != nil {
		merged = mergeFrontMatter(merged, page)
	}

	return merged
}

func mergeFrontMatter(merged map[string]any, frontMatter *FrontMatter) map[string]any {
	if frontMatter.Title != "" {
		merged["title"] = frontMatter.Title
	}
	if frontMatter.Description != "" {
		merged["description"] = frontMatter.Description
	}
	if frontMatter.Template != "" {
		merged["template"] = frontMatter.Template
	}

	if frontMatter.GenerateFeed {
		merged["generate_feed"] = frontMatter.GenerateFeed
	}

	if frontMatter.Img != "" {
		merged["img"] = frontMatter.Img
	}

	if frontMatter.Lang != "" {
		merged["lang"] = frontMatter.Lang
	}

	if frontMatter.Featured {
		merged["featured"] = frontMatter.Featured
	}

	if frontMatter.Draft {
		merged["draft"] = frontMatter.Draft
	}

	merged["tags"] = frontMatter.Tags
	merged["date"] = frontMatter.Date
	merged["update"] = frontMatter.Update

	merged["extra"] = MergeExtra(merged, frontMatter.Extra)
	return merged
}

func MergeExtra(base map[string]any, extra map[string]any) map[string]any {
	for k, v := range extra {
		if existing, exists := base[k]; exists {
			if existingMap, ok := existing.(map[string]any); ok {
				if newMap, ok := v.(map[string]any); ok {
					base[k] = mergeMaps(existingMap, newMap)
					continue
				}
			}
		}
		base[k] = v
	}
	return base
}

func mergeMaps(base, override map[string]any) map[string]any {
	result := make(map[string]any, len(base))
	maps.Copy(result, base)
	maps.Copy(result, override)
	return result
}

func parseFrontMatter[T any](content []byte) (*T, error) {
	frontMatterDelim := []byte("+++")
	parts := bytes.SplitN(content, frontMatterDelim, 3)

	if len(parts) < 3 {
		return nil, nil // No front matter found
	}

	frontMatter := bytes.TrimSpace(parts[1])

	var config T
	if err := toml.Unmarshal(frontMatter, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
