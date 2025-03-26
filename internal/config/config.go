package config

import (
	"bytes"
	"maps"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type GlobalConfig struct {
	BaseURL     string `toml:"base_url"`
	Title       string `toml:"title"`
	OutputDir   string `toml:"output_dir"`
	CompileSass bool   `toml:"compile_sass"`
	MinifyHTML  bool   `toml:"minify_html"`

	Taxonomies []TaxonomyConfig `toml:"taxonomies"`
	Markdown   MarkdownConfig   `toml:"markdown"`

	Extra map[string]any `toml:"extra"`
}

type MarkdownConfig struct {
	HighlightCode  bool   `toml:"highlight_code"`
	HighlightTheme string `toml:"highlight_theme"`
}

type SectionConfig struct {
	Title      string         `toml:"title"`
	Template   string         `toml:"template"`
	PaginateBy int            `toml:"paginate_by"`
	SortBy     string         `toml:"sort_by"`
	Render     bool           `toml:"render"`
	Extra      map[string]any `toml:"extra"`
}

type PageConfig struct {
	Title      string              `toml:"title"`
	Date       string              `toml:"date"`
	Draft      bool                `toml:"draft"`
	Taxonomies map[string][]string `toml:"taxonomies"`
	Template   string              `toml:"template"`
	Extra      map[string]any      `toml:"extra"`
}

type TaxonomyConfig struct {
	Name         string `toml:"name"`
	PaginateBy   int    `toml:"paginate_by"`
	GenerateFeed bool   `toml:"generate_feed"`
}

type MergedConfig struct {
	BaseURL     string
	Title       string
	OutputDir   string
	CompileSass bool
	MinifyHTML  bool
	Template    string
	PaginateBy  int
	SortBy      string
	Render      bool
	Date        string
	Draft       bool
	Taxonomies  map[string][]string
	Markdown    Markdown
	Extra       map[string]any
}

type Markdown struct {
	HighlightCode  bool
	HighlightTheme string
}

func LoadConfig(path string) (*GlobalConfig, error) {
	var cfg GlobalConfig
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadSectionConfig(sectionPath string) (*SectionConfig, error) {
	cfgPath := filepath.Join(sectionPath, "_index.md")
	content, err := os.ReadFile(cfgPath)
	if os.IsNotExist(err) {
		return &SectionConfig{}, nil
	}

	if err != nil {
		return nil, err
	}

	return parseFrontMatter[SectionConfig](content)
}

func LoadPageConfig(content []byte) (*PageConfig, error) {
	return parseFrontMatter[PageConfig](content)
}

func MergeConfigs(global *GlobalConfig, section *SectionConfig, page *PageConfig) *MergedConfig {
	merged := &MergedConfig{
		BaseURL:     global.BaseURL,
		Title:       global.Title,
		OutputDir:   global.OutputDir,
		CompileSass: global.CompileSass,
		MinifyHTML:  global.MinifyHTML,
		Markdown: Markdown{
			HighlightCode:  global.Markdown.HighlightCode,
			HighlightTheme: global.Markdown.HighlightTheme,
		},
		Extra: mergeMaps(nil, global.Extra),
	}

	if section != nil {
		merged = mergeSection(merged, section)
	}

	if page != nil {
		merged = mergePage(merged, page)
	}

	if merged.Template == "" {
		merged.Template = "default.html"
	}

	return merged
}

func mergeSection(merged *MergedConfig, section *SectionConfig) *MergedConfig {
	if section.Title != "" {
		merged.Title = section.Title
	}

	if section.Template != "" {
		merged.Template = section.Template
	}

	if section.PaginateBy > 0 {
		merged.PaginateBy = section.PaginateBy
	}

	if section.SortBy != "" {
		merged.SortBy = section.SortBy
	}

	merged.Render = section.Render
	merged.Extra = mergeMaps(merged.Extra, section.Extra)
	return merged
}

func mergePage(merged *MergedConfig, page *PageConfig) *MergedConfig {
	if page.Title != "" {
		merged.Title = page.Title
	}

	if page.Date != "" {
		merged.Date = page.Date
	}

	if page.Template != "" {
		merged.Template = page.Template
	}

	merged.Draft = page.Draft
	merged.Taxonomies = page.Taxonomies
	merged.Extra = mergeMaps(merged.Extra, page.Extra)

	return merged
}

func mergeMaps(base, override map[string]any) map[string]any {
	result := make(map[string]any)
	maps.Copy(result, base)

	for k, overrideVal := range override {
		if baseVal, exist := result[k]; exist {
			baseMap, baseIsMap := baseVal.(map[string]any)
			overrideMap, overrideIsMap := overrideVal.(map[string]any)
			if baseIsMap && overrideIsMap {
				result[k] = mergeMaps(baseMap, overrideMap)
				continue
			}
		}
		result[k] = overrideVal
	}

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
