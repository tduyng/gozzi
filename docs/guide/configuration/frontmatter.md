# Front Matter Configuration

Front matter appears at the top of content files (markdown pages and section indexes) and controls individual page/section behavior.

## What is Front Matter?

Front matter is metadata at the beginning of your content files, written in TOML format and wrapped in `+++` delimiters:

```markdown
+++
title = "My Page Title"
date = 2025-01-15
tags = ["go", "tutorial"]
+++

# Content starts here

Your markdown content...
```

## Available Fields

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

## Date Formats

Gozzi supports TOML date format:

```toml
# Date only (preferred)
date = 2025-01-15

# Date with time
date = 2025-01-15T10:30:00Z

# Date with time and timezone
date = 2025-01-15T10:30:00-08:00
```

## Draft Content

Control content visibility with the `draft` field:

```toml
+++
title = "Work in Progress"
draft = true  # Won't appear in production builds
+++
```

**Behavior:**
- Draft pages are excluded from production builds
- Still visible in development mode (`gozzi serve`)
- Not included in sitemaps or feeds
- Useful for unpublished content

## Template Selection

Specify which template to use:

```toml
+++
title = "About Me"
template = "prose.html"  # Use custom template
+++
```

**Default template mapping:**
- Blog posts → `post.html`
- Section indexes → `section.html`
- Pages → `page.html`

See [Template Mapping](/guide/templates/mapping) for details.

## Tags

Categorize content with tags:

```toml
+++
title = "Advanced Go Patterns"
tags = ["golang", "design-patterns", "advanced"]
+++
```

**Features:**
- Automatic tag page generation (`/tags/golang/`)
- Tag cloud/listing page (`/tags/`)
- Filterable in templates
- RSS feed filtering

## Featured Content

Mark important content as featured:

```toml
+++
title = "Important Announcement"
featured = true
+++
```

Use in templates:

```html
{{ range .Site.Pages }}
  {{ if .Config.featured }}
    <article class="featured">
      <h2>{{ .Title }}</h2>
    </article>
  {{ end }}
{{ end }}
```

## Images

Set cover/sharing images:

```toml
+++
title = "My Post"
img = "/img/cover.jpg"  # Absolute path from static/
+++
```

**Usage:**
- Social media previews (Open Graph, Twitter Cards)
- Featured images in listings
- Default fallback to site-level `img`

## RSS/Feed Control

Control feed inclusion per page:

```toml
+++
title = "Privacy Policy"
generate_feed = false  # Exclude from RSS feed
+++
```

## Language Override

Override site language per page:

```toml
+++
title = "Bonjour"
language = "fr"  # French content
+++
```

## Extended Front Matter

Add custom fields with `[extra]`:

```toml
+++
title = "Tutorial"

[extra]
difficulty = "intermediate"
duration = "30 minutes"
prerequisites = ["Go basics", "HTML/CSS"]
github_url = "https://github.com/user/repo"
+++
```

Access in templates:

```html
<p>Difficulty: {{ .Page.Extra.difficulty }}</p>
<p>Duration: {{ .Page.Extra.duration }}</p>

<ul>
  {{ range .Page.Extra.prerequisites }}
    <li>{{ . }}</li>
  {{ end }}
</ul>
```

## Section-Specific Fields

Additional fields for `_index.md` (section pages):

| Field | Type | Description |
|-------|------|-------------|
| `paginate_by` | integer | Number of items per page |
| `sort_by` | string | Sort children by "date", "title", or "updated" |
| `slug` | string | Override URL segment (instead of folder name) |

### Section Example

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

## Template Access

### Page Data

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
  
  <div class="content">
    {{ .Page.Content | safe }}
  </div>
</article>
```

### Section Data

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

## Common Patterns

### Blog Post

```toml
+++
title = "Understanding Go Interfaces"
description = "Deep dive into Go's interface system"
date = 2025-01-20
tags = ["golang", "tutorial"]
template = "post.html"

[extra]
reading_time = "8 min"
toc = true
series = "Go Fundamentals"
+++
```

### Documentation Page

```toml
+++
title = "API Reference"
description = "Complete API documentation"
template = "docs.html"
draft = false

[extra]
sidebar = true
edit_url = "https://github.com/user/repo/edit/main/content/api.md"
+++
```

### Landing Page

```toml
+++
title = "Welcome"
template = "home.html"
generate_feed = false

[extra]
hero_title = "Build Fast Static Sites"
hero_subtitle = "With Gozzi"
cta_text = "Get Started"
cta_url = "/guide/getting-started"
+++
```

## Best Practices

1. **Required fields first** - `title`, `date`, `description`
2. **Use consistent dates** - YYYY-MM-DD format
3. **Write descriptions** - 150-160 characters for SEO
4. **Tag consistently** - Use lowercase, hyphenated tags
5. **Set updated dates** - When making significant changes
6. **Use draft mode** - For work-in-progress content
7. **Leverage [extra]** - For custom, template-specific data

## Validation

Common front matter errors:

```toml
# ❌ Wrong: Missing quotes for string
title = My Title

# ✅ Correct: Strings need quotes
title = "My Title"

# ❌ Wrong: Date format
date = "2025-1-5"

# ✅ Correct: Use leading zeros
date = 2025-01-05

# ❌ Wrong: Array without brackets
tags = "go", "tutorial"

# ✅ Correct: Arrays use brackets
tags = ["go", "tutorial"]
```

---

**Related:**
- [Configuration Inheritance](/guide/configuration/inheritance)
- [Template Variables](/guide/templates/variables)
- [Content Structure](/guide/concepts/content-structure)
