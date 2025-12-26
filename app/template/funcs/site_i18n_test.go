// ABOUTME: Tests for i18n-related template functions
// ABOUTME: Validates language-aware URL generation and language detection
package funcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/i18n"
)

func TestSiteFuncs_LangURL(t *testing.T) {
	// Setup i18n
	i18nMgr := i18n.NewI18n("en", "")
	i18nMgr.AddLanguage("en", "English", 1, true)
	i18nMgr.AddLanguage("fr", "Français", 2, false)
	i18nMgr.AddLanguage("es", "Español", 3, false)

	ctx := &SiteContext{
		BaseURL: "https://example.com",
		I18n:    i18nMgr,
	}

	sf := NewSiteFuncs(ctx)

	tests := []struct {
		name       string
		langCode   string
		input      any
		want       string
		setupInput func() any
	}{
		{
			name:     "root page to french",
			langCode: "fr",
			input:    "/en/",
			want:     "/fr/",
		},
		{
			name:     "about page english to french",
			langCode: "fr",
			input:    "/en/about/",
			want:     "/fr/about/",
		},
		{
			name:     "about page french to spanish",
			langCode: "es",
			input:    "/fr/about/",
			want:     "/es/about/",
		},
		{
			name:     "blog post with slug",
			langCode: "fr",
			input:    "/en/blog/my-post/",
			want:     "/fr/blog/my-post/",
		},
		{
			name:     "node with permalink",
			langCode: "fr",
			setupInput: func() any {
				return &content.Node{
					Permalink: "/en/about/",
					Config: map[string]any{
						"lang": "en",
					},
				}
			},
			want: "/fr/about/",
		},
		{
			name:     "root node to french",
			langCode: "fr",
			setupInput: func() any {
				return &content.Node{
					Permalink: "/en/",
					Config: map[string]any{
						"lang": "en",
					},
				}
			},
			want: "/fr/",
		},
		{
			name:     "invalid language code",
			langCode: "invalid",
			input:    "/en/about/",
			want:     "/en/about/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.input
			if tt.setupInput != nil {
				input = tt.setupInput()
			}

			got := sf.LangURL(tt.langCode, input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSiteFuncs_LangURL_NoI18n(t *testing.T) {
	// No i18n configured
	ctx := &SiteContext{
		BaseURL: "https://example.com",
		I18n:    nil,
	}

	sf := NewSiteFuncs(ctx)

	// Should return original URL when i18n is not configured
	got := sf.LangURL("fr", "/about/")
	assert.Equal(t, "/about/", got)

	node := &content.Node{
		Permalink: "/blog/post/",
		Config:    map[string]any{},
	}
	got = sf.LangURL("fr", node)
	assert.Equal(t, "/blog/post/", got)
}

func TestSiteFuncs_CurrentLang(t *testing.T) {
	// Setup i18n
	i18nMgr := i18n.NewI18n("en", "")
	i18nMgr.AddLanguage("en", "English", 1, true)
	i18nMgr.AddLanguage("fr", "Français", 2, false)

	ctx := &SiteContext{
		I18n: i18nMgr,
	}

	sf := NewSiteFuncs(ctx)

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name: "from node config",
			input: &content.Node{
				Config: map[string]any{
					"lang": "fr",
				},
			},
			want: "fr",
		},
		{
			name: "from config map",
			input: map[string]any{
				"lang": "es",
			},
			want: "es",
		},
		{
			name: "fallback to default",
			input: &content.Node{
				Config: map[string]any{},
			},
			want: "en",
		},
		{
			name:  "no context",
			input: nil,
			want:  "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sf.CurrentLang(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSiteFuncs_CurrentLang_NoI18n(t *testing.T) {
	// No i18n configured
	ctx := &SiteContext{
		I18n: nil,
	}

	sf := NewSiteFuncs(ctx)

	// Should return "en" as ultimate fallback
	got := sf.CurrentLang(nil)
	assert.Equal(t, "en", got)
}
