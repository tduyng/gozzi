+++
title = "Site Configuration"
date = 2025-12-15
template = "page.html"
+++

The main site configuration defines global settings for your entire website.

## Core Settings

```toml
# config.toml - Core site configuration
base_url     = "https://example.com"              # Site's base URL (required)
feed_url     = "https://example.com/atom.xml"     # RSS feed URL (auto-generated if not set)
title        = "My Site"                          # Site title (required)
description  = "Description for website"          # Site description for SEO
output_dir   = "public"                           # Output directory (default: "public")
language     = "en"                               # Default language code (default: "en")
theme        = "default"                          # Theme name (optional)
img          = "/img/default-cover.webp"          # Default sharing image
generate_feed = true                              # Generate RSS/Atom feed (default: false)
syntax_theme = "dracula"                          # Syntax highlighting theme (default: "dracula")
minify_css = true                                 # Minify CSS files (default: false)
minify_html = true                                # Minify HTML output (default: false)
minify_js = true                                  # Minify JavaScript files (default: false)
minify_json = true                                # Minify JSON files (default: false)
minify_svg = true                                 # Minify SVG files (default: false)
minify_xml = true                                 # Minify XML files (default: false)
homepage_cache_sections = ["blog", "notes"]       # Sections that trigger homepage rebuild (optional)
```

## Configuration Options Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `base_url` | string | **required** | Site's base URL for absolute links |
| `feed_url` | string | `{base_url}/atom.xml` | RSS/Atom feed URL |
| `title` | string | **required** | Site title for pages and feeds |
| `description` | string | "" | Site description for SEO meta tags |
| `language` | string | "en" | Default language code (ISO 639-1) |
| `output_dir` | string | "public" | Directory for generated static files |
| `theme` | string | "" | Theme name (if using themes) |
| `img` | string | "" | Default sharing image for social media |
| `generate_feed` | boolean | false | Whether to generate RSS/Atom feeds |
| `syntax_theme` | string | "dracula" | Syntax highlighting theme for code blocks |
| `minify_css` | boolean | false | Minify CSS files (36% avg reduction) |
| `minify_html` | boolean | false | Minify HTML output (53% avg reduction) |
| `minify_js` | boolean | false | Minify JavaScript files (53% avg reduction) |
| `minify_json` | boolean | false | Minify JSON files (51% avg reduction) |
| `minify_svg` | boolean | false | Minify SVG files (29% avg reduction) |
| `minify_xml` | boolean | false | Minify XML files (19% avg reduction) |
| `homepage_cache_sections` | array of strings | `[]` | Sections that trigger homepage rebuild (see [Performance](#performance-optimization)) |

## Minimal Blog Configuration

```toml
base_url = "https://myblog.com"
title = "My Blog"
description = "Thoughts and tutorials"
language = "en"

[extra]
name = "Author Name"
bio = "Writer and developer"
```

## Complete Site Configuration

```toml
base_url = "https://mysite.com"
feed_url = "https://mysite.com/feed.xml"
title = "Professional Portfolio"
description = "Software engineer and technical writer"
language = "en"
output_dir = "dist"
generate_feed = true
homepage_cache_sections = ["blog", "projects"]

[extra]
name = "Jane Developer"
id = "janedev"
bio = "Full-stack developer with 10+ years experience"
avatar = "/img/profile.webp"
img = "/img/social-card.webp"

links = [
  { name = "GitHub", icon = "github", url = "https://github.com/janedev" },
  { name = "LinkedIn", icon = "linkedin", url = "https://linkedin.com/in/janedev" },
  { name = "Email", icon = "email", url = "mailto:jane@example.com" },
]

sections = [
  { name = "blog", path = "/blog", is_external = false },
  { name = "projects", path = "/projects", is_external = false },
  { name = "speaking", path = "/speaking", is_external = false },
  { name = "about", path = "/about", is_external = false },
]

theme_config = {
  color_scheme = "auto",
  syntax_highlighting = true,
  show_reading_time = true,
  enable_toc = true
}

footer_copyright = "©2025 Jane Developer"
enable_analytics = true
```

## Required Fields

Gozzi requires two fields in your site configuration:

1. **`base_url`** - The full URL where your site will be hosted
2. **`title`** - The name of your site

Without these fields, Gozzi will fail to build.

## Syntax Highlighting

Configure syntax highlighting theme for code blocks:

```toml
syntax_theme = "github-dark"  # Default: "dracula"
```

See [Syntax Highlighting](/features/syntax-highlighting) for complete theme list and examples.

## Asset Minification

Optimize your site's performance by enabling minification for various file types:

```toml
minify_css = true   # Reduces CSS file size by ~36%
minify_html = true  # Reduces HTML output by ~53%
minify_js = true    # Reduces JavaScript file size by ~53%
minify_json = true  # Reduces JSON file size by ~51%
minify_svg = true   # Reduces SVG file size by ~29%
minify_xml = true   # Reduces XML file size by ~19%
```

Benefits:
- Smaller file sizes for faster page loads
- Removes whitespace, comments, and redundant code
- Graceful fallback - uses original if minification fails
- No significant impact on build time

## Performance Optimization

### Homepage Cache Control

Control when your homepage is rebuilt during incremental builds (serve mode):

```toml
# Only rebuild homepage when these sections change
homepage_cache_sections = ["blog", "notes"]
```

**How it works:**

Gozzi caches rendered pages to speed up incremental builds. When you modify a file in serve mode, Gozzi determines which pages need rebuilding. The `homepage_cache_sections` setting controls when your homepage is invalidated and rebuilt.

**Configuration options:**

```toml
# Option 1: Specify sections (recommended for performance)
homepage_cache_sections = ["blog", "notes"]
# Homepage rebuilds only when files in 'blog' or 'notes' sections change

# Option 2: Omit or leave empty (safe default)
# homepage_cache_sections = []
# Homepage rebuilds whenever ANY page changes (slower but always correct)
```

**When to use:**

| Site Structure | Recommended Config | Reason |
|----------------|-------------------|---------|
| Homepage shows blog posts only | `["blog"]` | Optimal performance |
| Homepage shows blog + notes | `["blog", "notes"]` | Only rebuild when relevant sections change |
| Homepage shows content from many sections | Omit config | Safe default ensures correctness |
| Unsure what homepage references | Omit config | Homepage always stays fresh |

**Example scenarios:**

```toml
# Blog-only site
homepage_cache_sections = ["blog"]
# Changing a blog post → homepage rebuilds ✓
# Changing about page → homepage stays cached ✓

# Multi-section site
homepage_cache_sections = ["blog", "notes", "projects"]
# Changing any listed section → homepage rebuilds ✓
# Changing contact page → homepage stays cached ✓

# Complex site (not configured)
# Changing ANY page → homepage rebuilds ✓ (safe but slower)
```

**Performance impact:**

- **With config:** Homepage only rebuilds when necessary (faster incremental builds)
- **Without config:** Homepage rebuilds on every change (safe but slower)
- For typical sites (<100 pages): Negligible impact
- For large sites (>500 pages): Can improve incremental build time by 20-50%

**Important notes:**

- This only affects **incremental builds** (serve mode)
- Full builds always regenerate everything
- If homepage looks stale, add the missing section to the config
- Root-level page changes never trigger homepage rebuild

## Best Practices

1. **Use HTTPS** for `base_url` in production
2. **Keep descriptions concise** - 150-160 characters optimal for SEO
3. **Set language code** - Helps with SEO and accessibility
4. **Use consistent paths** - End all URLs without trailing slashes
5. **Organize with `[extra]`** - Keep custom data separate from core config
6. **Choose appropriate syntax theme** - Match your site's design for consistency
7. **Enable minification in production** - Reduces bandwidth and improves load times
8. **Configure homepage cache sections** - Optimize incremental build performance

---

**Related:**
- [Syntax Highlighting](/features/syntax-highlighting)
- [Environment-Specific Configs](/config/environment)
- [Extended Configuration](/config/extended)
- [Front Matter](/config/frontmatter)
