+++
title = "Content Features"
date = 2025-12-15
template = "page.html"
+++

Gozzi provides powerful content management features built directly into the static site generator.

## Table of Contents (TOC)

Gozzi automatically generates a table of contents from your content headings.

### Usage in Templates

::: details Example Template Code

```html
<!-- Basic TOC display -->
{{ if .Toc }}
<nav class="toc">
    <h3>Table of Contents</h3>
    {{ .Toc }}
</nav>
{{ end }}
```

:::

### How It Works

- Automatically parses markdown headings (`#`, `##`, `###`, etc.)
- Generates nested HTML list structure
- Includes anchor links for navigation
- Available as `.Toc` variable in templates

### Customization

Style the TOC with CSS:

```css
.toc {
    background: var(--bg-secondary);
    padding: 1.5rem;
    border-radius: 8px;
    margin: 2rem 0;
}

.toc ul {
    list-style: none;
    padding-left: 1rem;
}

.toc a {
    color: var(--text-primary);
    text-decoration: none;
}

.toc a:hover {
    text-decoration: underline;
}
```

## Tag Management

Comprehensive tag management with automatic collection across content.

### Template Variables

::: details Template Examples

```html
<!-- Display tags for current page -->
{{ range .Tags }}
<span class="tag">{{ . }}</span>
{{ end }}

<!-- Access all site tags -->
{{ range .Site.Tags }}
<a href="/tags/{{ . | urlize }}">{{ . }}</a>
{{ end }}
```

:::

### Front Matter

Define tags in your content:

```yaml
---
title: 'My Post'
tags: ['golang', 'web-development', 'tutorial']
---
```

### Tag Pages

Gozzi automatically generates:

- Individual tag pages (`/tags/golang/`)
- Tag listing page (`/tags/`)
- Tag counts and post collections

## Pagination

Built-in pagination for large content collections.

### Configuration

```toml
[pagination]
per_page = 10
pagination_path = "page"
```

### Template Usage

::: details Pagination Template

```html
{{ range .Paginate.Items }}
<article>
    <h2>{{ .Title }}</h2>
</article>
{{ end }} {{ if .Paginate.HasPrev }}
<a href="{{ .Paginate.PrevURL }}">← Previous</a>
{{ end }} {{ if .Paginate.HasNext }}
<a href="{{ .Paginate.NextURL }}">Next →</a>
{{ end }}
```

:::

### Available Variables

- `.Paginate.Items` - Current page items
- `.Paginate.HasNext` - Has next page
- `.Paginate.HasPrev` - Has previous page
- `.Paginate.NextURL` - Next page URL
- `.Paginate.PrevURL` - Previous page URL
- `.Paginate.CurrentPage` - Current page number
- `.Paginate.TotalPages` - Total page count

## RSS/Atom Feeds

Automatically generates RSS feeds for content syndication.

### Features

- **Location**: `/atom.xml`
- **Format**: Atom 1.0
- **Content**: Latest posts with full content
- **Automatic**: No configuration required

### Configuration

```toml
[feed]
enabled = true
limit = 20
include_content = true
title = "My Blog Feed"
description = "Latest posts from my blog"
```

### Customization

Add feed link to your template:

```html
<link rel="alternate" type="application/atom+xml" title="{{ .Site.Title }}" href="/atom.xml" />
```

## Content Analytics

### Word Count & Reading Time

Automatically calculated metrics for each page.

::: details Template Examples

```html
<p>{{ .WordCount }} words</p>

{{ if gt .WordCount 1000 }}
<span class="long-read">Long read</span>
{{ end }}

<p>{{ .ReadingTime }} min read</p>
```

:::

### Available Variables

- `.WordCount` - Total words in content
- `.ReadingTime` - Estimated reading time (minutes)

### Calculation

Reading time assumes ~200 words per minute, a common average reading speed.

## Social Media Integration

Automatic image URL generation for social sharing.

### Front Matter

```yaml
---
title: 'My Post'
featured_image: 'cover.jpg'
image: '/images/post-cover.png'
description: 'Post description for social media'
---
```

### Template Usage

```html
<!-- Open Graph -->
<meta property="og:title" content="{{ .Title }}" />
<meta property="og:description" content="{{ .Description }}" />
<meta property="og:image" content="{{ .FeaturedImage }}" />

<!-- Twitter Card -->
<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:title" content="{{ .Title }}" />
<meta name="twitter:description" content="{{ .Description }}" />
<meta name="twitter:image" content="{{ .FeaturedImage }}" />
```

## Performance

All content features are optimized for speed:

- **Server-side generation** - All processing happens during build
- **Fast builds** - Minimal overhead even for large sites
- **Caching** - Generated content cached until source changes
- **No JavaScript** - All features work without client-side JS

---

**Related:**
- [Configuration](/config/)
- [Template Functions](/functions/)
- [Front Matter](/config/frontmatter)
