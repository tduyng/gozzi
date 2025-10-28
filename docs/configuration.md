# Configuration

Gozzi uses TOML files for configuration, offering a clean and readable format for site settings. The main configuration file is `config.toml` by default, but you can specify alternative files using the `--config` flag.

## Table of Contents

- [Site Configuration](#site-configuration)
- [Environment-Specific Configs](#environment-specific-configs)
- [Front Matter Configuration](#front-matter-configuration)
- [Section Configuration](#section-configuration)
- [Page Configuration](#page-configuration)
- [Advanced Configuration](#advanced-configuration)

## Site Configuration

The main site configuration defines global settings for your entire website.

### Core Settings

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
```

### Extended Configuration

The `[extra]` section allows custom data accessible in templates:

```toml
[extra]
# Personal Information
name   = "Your Name"
id     = "username"
bio    = "Software Engineer"
avatar = "/img/avatar.webp"

# Social Links
links = [
  { name = "GitHub", icon = "github", url = "https://github.com/username" },
  { name = "Twitter", icon = "twitter", url = "https://x.com/username" },
  { name = "Email", icon = "email", url = "mailto:hello@example.com" },
]

# Navigation Configuration
sections = [
  { name = "blog", path = "/blog", is_external = false },
  { name = "notes", path = "/notes", is_external = false },
  { name = "about", path = "/about", is_external = false },
]

footer_menu = [
  { name = "tags", path = "/tags", is_external = false },
  { name = "subscribe", path = "/subscribe", is_external = false },
  { name = "contact", path = "/contact", is_external = false },
  { name = "privacy", path = "/privacy", is_external = false },
  { name = "rss", path = "/atom.xml", is_external = false },
]

# Navigation Styling
nav_separator = "☰"
nav_wrapper_left = "["
nav_wrapper_right = "]"
nav_wrapper_separator = ","

# Content Configuration
ended_words = "Thanks for reading! Feel free to reach out with questions or feedback."
reaction_align = "right"                          # "left" | "center" | "right"
reaction_endpoint = "https://your-reaction-api.com"
footer_copyright = "©2025 Your Name"

# Outdated Content Alerts
outdate_alert_text_before = "This article was last updated "
outdate_alert_text_after = " days ago and may be out of date."

# Analytics
enable_analytics = true
```

### Configuration Options Reference

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

## Environment-Specific Configs

Use different configuration files for different environments:

### Development Configuration
```toml
# config.dev.toml
base_url = "http://localhost:1313"
feed_url = "http://localhost:1313/atom.xml"
title = "My Site (Dev)"
description = "Development version"
language = "en"
output_dir = "public"

[extra]
enable_analytics = false  # Disable analytics in development
```

### Production Configuration
```toml
# config.prod.toml  
base_url = "https://mysite.com"
feed_url = "https://mysite.com/atom.xml"
title = "My Site"
description = "Production website"
language = "en"
output_dir = "public"

[extra]
enable_analytics = true   # Enable analytics in production
```

### Usage
```sh
# Development
gozzi serve --config config.dev.toml

# Production build
gozzi build --config config.prod.toml
```

## Front Matter Configuration

Front matter appears at the top of content files and controls individual page/section behavior.

### Available Fields

All front matter fields that can be used in both pages and sections:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `title` | string | "" | Page/section title |
| `description` | string | "" | Meta description for SEO |
| `date` | date | current | Publication date (YYYY-MM-DD) |
| `updated` | date | "" | Last modified date |
| `template` | string | auto | Template file to use (e.g., "post.html") |
| `draft` | boolean | false | Hide from production builds |
| `generate_feed` | boolean | true | Include in RSS/Atom feeds |
| `featured` | boolean | false | Mark as featured content |
| `tags` | array | [] | List of tags for categorization |
| `language` | string | site default | Override site language |
| `img` | string | site default | Custom cover/sharing image |

## Section Configuration

Sections are directories with `_index.md` files that create list pages.

### Basic Section Example
```toml
+++
title = "Blog"
description = "Latest articles and thoughts"
template = "blog.html"
date = 2025-01-01

[extra]
hero_text = "Welcome to my blog!"
show_reading_time = true
paginate_by = 10
sort_by = "date"
+++

# Blog Section

This is the blog landing page content.
```

### Section-Specific Front Matter

Additional fields available for sections:

| Field | Type | Description |
|-------|------|-------------|
| `paginate_by` | integer | Number of items per page |
| `sort_by` | string | Sort children by "date", "title", or "updated" |
| `slug` | string | Override URL segment (instead of folder name) |

### Template Access

In templates, section data is available as:
```html
<h1>{{ .Section.Title }}</h1>
<p>{{ .Section.Description }}</p>
<p>{{ .Section.Extra.hero_text }}</p>

<!-- Children pages -->
{{ range .Section.Children }}
  <article>
    <h2>{{ .Title }}</h2>
    <time>{{ .Date.Format "2006-01-02" }}</time>
  </article>
{{ end }}
```

## Page Configuration

Individual content pages (`.md` files or `index.md` in leaf bundles).

### Basic Page Example
```toml
+++
title = "Getting Started with Gozzi"
description = "A comprehensive guide to building static sites with Gozzi"
date = 2025-01-15
updated = 2025-01-20
tags = ["tutorial", "go", "static-site"]
template = "post.html"
draft = false
generate_feed = true
featured = true

[extra]
toc = true
reading_time = "5 min"
img = "/img/gozzi-tutorial.webp"
author = "Jane Doe"
+++

# Getting Started with Gozzi

Your markdown content here...
```

### Page Template Access

In templates, page data is available as:
```html
<article>
  <h1>{{ .Page.Title }}</h1>
  <p>{{ .Page.Description }}</p>
  <time>{{ .Page.Date.Format "January 2, 2006" }}</time>
  
  {{ if .Page.Config.updated }}
    <p>Updated: {{ .Page.Config.updated.Format "January 2, 2006" }}</p>
  {{ end }}
  
  {{ if .Page.Config.featured }}
    <span class="featured-badge">Featured</span>
  {{ end }}
  
  <div class="tags">
    {{ range .Page.Config.tags }}
      <span class="tag">{{ . }}</span>
    {{ end }}
  </div>
  
  {{ if .Page.Extra.toc }}
    <!-- Include table of contents -->
  {{ end }}
  
  <div class="content">
    {{ .Page.Content | safe }}
  </div>
</article>
```

## Advanced Configuration

### Custom Data Structures

The `[extra]` section supports complex nested data:

```toml
[extra]
# Simple values
author = "John Doe"
reading_time = "5 min"

# Arrays
tags = ["go", "static-site", "web"]
categories = ["Programming", "Web Development"]

# Objects
social = { twitter = "@username", github = "username" }

# Nested structures
navigation = [
  { name = "Home", url = "/", active = true },
  { name = "Blog", url = "/blog", active = false },
  { name = "About", url = "/about", active = false },
]

# Analytics configuration
analytics = { 
  google_id = "GA-XXXXXXXX", 
  plausible_domain = "example.com" 
}

# Theme customization
theme_config = {
  color_scheme = "dark",
  syntax_highlighting = true,
  show_reading_time = true,
  enable_comments = false
}
```

### Accessing Custom Data in Templates

```html
<!-- Simple values -->
<p>By {{ .Site.Extra.author }}</p>
<p>Reading time: {{ .Page.Extra.reading_time }}</p>

<!-- Arrays -->
{{ range .Page.Extra.categories }}
  <span class="category">{{ . }}</span>
{{ end }}

<!-- Objects -->
<a href="https://twitter.com/{{ .Site.Extra.social.twitter }}">Twitter</a>
<a href="https://github.com/{{ .Site.Extra.social.github }}">GitHub</a>

<!-- Nested structures -->
<nav>
  {{ range .Site.Extra.navigation }}
    <a href="{{ .url }}"{{ if .active }} class="active"{{ end }}>
      {{ .name }}
    </a>
  {{ end }}
</nav>
```

### Configuration Inheritance

Configuration follows a hierarchy where more specific settings override general ones:

1. **Site config** (global defaults)
2. **Section config** (overrides site for that section)
3. **Page config** (overrides both site and section)

```toml
# Site config.toml
[extra]
show_reading_time = true
theme = "light"

# Section _index.md
+++
[extra]
show_reading_time = false  # Overrides site setting for this section
+++

# Page front matter
+++
[extra]
theme = "dark"  # Overrides both site and section settings for this page
+++
```

### Multiple Configuration Files

Organize complex configurations across multiple files:

```sh
config/
├── config.toml           # Main configuration
├── config.dev.toml       # Development overrides
├── config.staging.toml   # Staging environment
└── config.prod.toml      # Production environment
```

## Configuration Examples

### Minimal Blog Configuration
```toml
base_url = "https://myblog.com"
title = "My Blog"
description = "Thoughts and tutorials"
language = "en"

[extra]
name = "Author Name"
bio = "Writer and developer"
```

### Complete Site Configuration
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

## Troubleshooting Configuration

### Common Issues

**TOML Syntax Errors**:
```sh
# Error: Key-value pairs must be on the same line
base_url = 
  "https://example.com"  # ❌ Wrong

base_url = "https://example.com"  # ✅ Correct
```

**Missing Required Fields**:
```toml
# Error: Missing required base_url
title = "My Site"  # ❌ Won't work without base_url

# Fix: Add required fields
base_url = "https://example.com"  # ✅ Required
title = "My Site"
```

**Date Format Issues**:
```toml
# Error: Wrong date format
date = "2025-1-15"  # ❌ Missing leading zeros

# Fix: Use YYYY-MM-DD format
date = 2025-01-15   # ✅ Correct
```

### Validation Commands

```sh
# Test configuration with build
gozzi build --config config.toml

# Test development server
gozzi serve --config config.dev.toml --port 3000
```

## Next Steps

- [CLI Reference](cli.md) - Learn command-line options
- [Content Structure](content_structure.md) - Organize your content
- [Templates](templates.md) - Customize your site's appearance
- [HTML Functions](htmlfunc.md) - Use template functions

## References

- [TOML Specification](https://toml.io/) - Official TOML documentation
- [Example Configuration](https://github.com/tduyng/tduyng.github.io/blob/main/config.toml) - Real-world example
