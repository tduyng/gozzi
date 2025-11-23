# Page Configuration

Individual content pages are markdown files (`.md`) or page bundles (`index.md` in a directory) that represent single pieces of content.

## What are Pages?

Pages are individual content files in Gozzi:

```
content/
├── about.md               # Simple page
├── blog/
│   ├── _index.md         # Section config
│   ├── post-1.md         # Simple page
│   └── post-2/           # Page bundle
│       ├── index.md      # Page content
│       ├── cover.jpg     # Page asset
│       └── diagram.png   # Page asset
```

## Page Types

### Simple Page

Single markdown file:

```markdown
+++
title = "About Me"
date = 2025-01-15
+++

# About

I'm a software engineer...
```

### Page Bundle

Directory with `index.md` and assets:

```
post-with-images/
├── index.md         # Page content
├── cover.jpg        # Local image
└── diagram.png      # Local image
```

```markdown
+++
title = "My Post with Images"
date = 2025-01-20
img = "cover.jpg"    # Relative to this directory
+++

# My Post

![Diagram](diagram.png)  <!-- Works! -->
```

## Basic Page Example

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

## Available Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `title` | string | "" | Page title (required) |
| `description` | string | "" | Meta description for SEO |
| `date` | date | current | Publication date |
| `updated` | date | "" | Last modified date |
| `template` | string | auto | Template file to use |
| `draft` | boolean | false | Hide from production |
| `generate_feed` | boolean | true | Include in RSS/Atom |
| `featured` | boolean | false | Mark as featured |
| `tags` | array | [] | List of tags |
| `language` | string | site | Override language |
| `img` | string | site | Cover/sharing image |
| `slug` | string | filename | Custom URL slug |

## Title and Description

Required for good SEO:

```toml
+++
title = "Understanding Go Interfaces"
description = "Deep dive into Go's interface system with practical examples"
+++
```

**Best practices:**
- Title: Clear, descriptive, under 60 characters
- Description: 150-160 characters, includes key terms

## Dates

### Publication Date

```toml
+++
date = 2025-01-15
+++
```

### Last Updated

```toml
+++
date = 2025-01-15
updated = 2025-01-20  # Show "updated" badge
+++
```

### Date Formats

```toml
# Date only (preferred)
date = 2025-01-15

# With time
date = 2025-01-15T10:30:00Z

# With timezone
date = 2025-01-15T10:30:00-08:00
```

## Template Selection

Specify which template renders the page:

```toml
+++
title = "About"
template = "prose.html"  # Use custom template
+++
```

Default mapping:
- `/blog/*.md` → `post.html`
- `/about.md` → `page.html`
- Custom: `template` field overrides

## Draft Mode

Hide pages from production:

```toml
+++
title = "Work in Progress"
draft = true
+++
```

**Behavior:**
- ✅ Visible in development (`gozzi serve`)
- ❌ Hidden in production builds
- ❌ Not in sitemaps or feeds
- ❌ Not in `.Site.Pages`

## Tags

Categorize and filter content:

```toml
+++
title = "Advanced Go Patterns"
tags = ["golang", "design-patterns", "advanced"]
+++
```

**Generated pages:**
- `/tags/golang/` - All posts with "golang" tag
- `/tags/` - All tags index

## Featured Content

Highlight important pages:

```toml
+++
title = "Important Post"
featured = true
+++
```

Use in templates:

```html
{{ range .Site.Pages }}
  {{ if .Config.featured }}
    <article class="featured">
      <span class="badge">Featured</span>
      <h2>{{ .Title }}</h2>
    </article>
  {{ end }}
{{ end }}
```

## Images

### Cover Images

```toml
+++
title = "My Post"
img = "/img/cover.jpg"  # Absolute from static/
+++
```

### Page Bundle Images

```toml
+++
title = "Tutorial"
img = "cover.jpg"  # Relative to page bundle
+++
```

**Usage:**
- Social media previews
- Featured image in listings
- Header backgrounds

## Custom Slugs

Override URL path:

```toml
+++
title = "Getting Started Guide"
slug = "quickstart"  # /blog/quickstart/ instead of /blog/getting-started-guide/
+++
```

## RSS/Feed Control

```toml
+++
title = "Privacy Policy"
generate_feed = false  # Don't include in RSS
+++
```

## Language Override

```toml
+++
title = "Bonjour le monde"
language = "fr"
+++
```

## Extended Page Data

Add custom fields with `[extra]`:

```toml
+++
title = "Tutorial"

[extra]
# Display options
toc = true
reading_time = "8 min"
difficulty = "intermediate"
show_author = true

# Content metadata
author = "Jane Doe"
author_url = "https://example.com"
series = "Go Fundamentals"
part = 3

# Links
github_url = "https://github.com/user/repo"
demo_url = "https://demo.example.com"

# Prerequisites
prerequisites = ["Go basics", "HTML/CSS"]
related_posts = ["/blog/post-1", "/blog/post-2"]
+++
```

## Template Access

```html
<article>
  <header>
    <h1>{{ .Page.Title }}</h1>
    <p class="description">{{ .Page.Description }}</p>
    
    <div class="meta">
      <time datetime="{{ .Page.Date.Format "2006-01-02" }}">
        {{ .Page.Date.Format "January 2, 2006" }}
      </time>
      
      {{ if .Page.Config.updated }}
        <span class="updated">
          Updated: {{ .Page.Config.updated.Format "Jan 2, 2006" }}
        </span>
      {{ end }}
      
      {{ if .Page.Config.featured }}
        <span class="badge featured">Featured</span>
      {{ end }}
      
      {{ if .Page.Extra.reading_time }}
        <span class="reading-time">{{ .Page.Extra.reading_time }}</span>
      {{ end }}
    </div>
    
    {{ if .Page.Config.tags }}
      <div class="tags">
        {{ range .Page.Config.tags }}
          <a href="/tags/{{ . | urlize }}" class="tag">{{ . }}</a>
        {{ end }}
      </div>
    {{ end }}
  </header>
  
  {{ if .Page.Extra.toc }}
    <aside class="toc">
      <h3>Table of Contents</h3>
      {{ .Toc }}
    </aside>
  {{ end }}
  
  <div class="content">
    {{ .Page.Content | safe }}
  </div>
  
  <footer>
    {{ if .Page.Extra.author }}
      <p>By {{ .Page.Extra.author }}</p>
    {{ end }}
  </footer>
</article>
```

## Common Page Patterns

### Blog Post

```toml
+++
title = "Understanding Goroutines"
description = "Learn how to use goroutines effectively in Go"
date = 2025-01-20
tags = ["golang", "concurrency"]
template = "post.html"

[extra]
reading_time = "6 min"
toc = true
author = "John Doe"
+++
```

### Documentation Page

```toml
+++
title = "Installation Guide"
description = "How to install and set up Gozzi"
template = "docs.html"
date = 2025-01-01

[extra]
sidebar = true
next_page = "/docs/quickstart"
prev_page = "/docs/intro"
+++
```

### Landing Page

```toml
+++
title = "Welcome to My Site"
description = "Personal website and blog"
template = "home.html"
generate_feed = false

[extra]
hero_title = "Hi, I'm Jane"
hero_subtitle = "Software Engineer & Writer"
show_recent_posts = true
+++
```

### Tutorial

```toml
+++
title = "Build a REST API in Go"
description = "Step-by-step tutorial for building REST APIs"
date = 2025-01-15
tags = ["golang", "api", "tutorial"]

[extra]
difficulty = "intermediate"
duration = "45 minutes"
prerequisites = ["Go installed", "Basic HTTP knowledge"]
github_url = "https://github.com/user/api-tutorial"
+++
```

## Best Practices

1. **Always set title** - Required for navigation and SEO
2. **Write descriptions** - 150-160 characters optimal
3. **Use tags consistently** - Lowercase, hyphenated
4. **Set dates** - Publication date required
5. **Update timestamps** - When making significant changes
6. **Use draft mode** - For work-in-progress
7. **Organize with bundles** - For pages with multiple assets
8. **Custom [extra] data** - For template flexibility

## Validation

Common errors:

```toml
# ❌ Missing quotes
title = My Title

# ✅ Correct
title = "My Title"

# ❌ Wrong date format
date = "Jan 15, 2025"

# ✅ Correct
date = 2025-01-15

# ❌ Invalid array
tags = "go", "tutorial"

# ✅ Correct
tags = ["go", "tutorial"]
```

---

**Related:**
- [Front Matter](/guide/configuration/frontmatter)
- [Sections Configuration](/guide/configuration/sections)
- [Content Structure](/guide/content-structure)
- [Template Variables](/guide/templates/variables)
