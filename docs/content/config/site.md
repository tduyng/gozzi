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

## CSS and HTML Minification

Optimize your site's performance by enabling minification:

```toml
minify_css = true   # Reduces CSS file size by ~36%
minify_html = true  # Reduces HTML output by ~53%
```

Benefits:
- Smaller file sizes for faster page loads
- Removes whitespace, comments, and redundant code
- Graceful fallback - uses original if minification fails
- No impact on build time (~150ms)

## Best Practices

1. **Use HTTPS** for `base_url` in production
2. **Keep descriptions concise** - 150-160 characters optimal for SEO
3. **Set language code** - Helps with SEO and accessibility
4. **Use consistent paths** - End all URLs without trailing slashes
5. **Organize with `[extra]`** - Keep custom data separate from core config
6. **Choose appropriate syntax theme** - Match your site's design for consistency
7. **Enable minification in production** - Reduces bandwidth and improves load times

---

**Related:**
- [Syntax Highlighting](/features/syntax-highlighting)
- [Environment-Specific Configs](/config/environment)
- [Extended Configuration](/config/extended)
- [Front Matter](/config/frontmatter)
