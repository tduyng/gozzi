+++
title = "Partials"
date = 2025-12-15
template = "page.html"
+++

Partials are reusable template snippets that promote DRY (Don't Repeat Yourself) principles and consistent design across your site.

## What are Partials?

Partials are shared template components stored in `templates/partials/` that can be included in multiple templates.

**Benefits:**
- Reduce code duplication
- Maintain consistency
- Easy updates (change once, apply everywhere)
- Cleaner, more maintainable templates

## Including Partials

### Basic Inclusion

```html
<!-- Include partial with current context -->
{{ template "partials/_head.html" . }}
{{ template "partials/_header.html" . }}
{{ template "partials/_footer.html" . }}
```

### With Custom Data

```html
<!-- Pass specific data to partial -->
{{ template "partials/_pagination.html" (dict "pages" .Section.Children "current" .Page) }}
```

## Common Partials

### `_head.html` - Document Head

```html
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{ priority .Page.Config.title .Section.Config.title .Site.Config.title }}</title>
    
    <!-- SEO Meta Tags -->
    {{ $description := priority .Page.Config.description .Section.Config.description .Site.Config.description }}
    <meta name="description" content="{{ $description }}">
    
    <!-- Open Graph -->
    <meta property="og:title" content="{{ .Page.Config.title }}">
    <meta property="og:description" content="{{ $description }}">
    {{ if .Page.Config.img }}
    <meta property="og:image" content="{{ .Page.Config.img_url }}">
    {{ end }}
    
    <!-- Canonical URL -->
    <link rel="canonical" href="{{ .Page.URL }}">
    
    <!-- Stylesheets -->
    <link rel="stylesheet" href="/css/main.css">
</head>
```

### `_header.html` - Site Header

```html
<header class="site-header">
    <div class="container">
        <a href="/" class="site-title">{{ .Site.Config.title }}</a>
        
        <nav class="main-nav">
            {{ range .Site.Config.extra.sections }}
                {{ $current := "" }}
                {{ if eq .path $.Page.Section }}
                    {{ $current = "current" }}
                {{ end }}
                
                <a href="{{ .path }}" class="nav-link {{ $current }}">
                    {{ .name }}
                </a>
            {{ end }}
        </nav>
        
        <div class="social-links">
            {{ range .Site.Config.extra.links }}
                <a href="{{ .url }}" class="social-link" title="{{ .name }}">
                    {{ $icon_path := printf "static/icon/%s.svg" .icon }}
                    {{ load $icon_path }}
                </a>
            {{ end }}
        </div>
    </div>
</header>
```

### `_footer.html` - Site Footer

```html
<footer class="site-footer">
    <div class="container">
        <p>{{ .Site.Config.extra.footer_copyright }}</p>
        
        <nav class="footer-nav">
            {{ range .Site.Config.extra.footer_menu }}
                <a href="{{ .path }}">{{ .name }}</a>
            {{ end }}
        </nav>
        
        <div class="social-links">
            {{ range .Site.Config.extra.links }}
                <a href="{{ .url }}">{{ .name }}</a>
            {{ end }}
        </div>
    </div>
</footer>
```

### `_toc.html` - Table of Contents

```html
{{ if .Page.Config.extra.toc }}
<nav class="toc">
    <h3>Table of Contents</h3>
    {{ .Toc }}
</nav>
{{ end }}
```

### `_word_count.html` - Reading Time

```html
{{ if .Page.WordCount }}
<div class="reading-stats">
    <span class="word-count">{{ .Page.WordCount }} words</span>
    <span class="reading-time">{{ .Page.ReadingTime }} min read</span>
</div>
{{ end }}
```

### `_sharing.html` - Social Sharing

```html
<div class="sharing">
    <h4>Share this post</h4>
    {{ $url := .Page.Permalink }}
    {{ $title := .Page.Config.title }}
    
    <a href="https://twitter.com/intent/tweet?url={{ $url }}&text={{ $title }}" 
       target="_blank" 
       rel="noopener">
        Twitter
    </a>
    
    <a href="https://www.linkedin.com/sharing/share-offsite/?url={{ $url }}" 
       target="_blank" 
       rel="noopener">
        LinkedIn  
    </a>
    
    <a href="https://www.facebook.com/sharer/sharer.php?u={{ $url }}" 
       target="_blank" 
       rel="noopener">
        Facebook
    </a>
</div>
```

### `_scripts.html` - JavaScript Includes

```html
<script src="/js/main.js"></script>

{{ if .Site.Config.extra.enable_analytics }}
    {{ template "partials/_analytics.html" . }}
{{ end }}

{{ if .Page.Config.extra.enable_comments }}
    {{ template "partials/_comments.html" . }}
{{ end }}
```

### `_comments.html` - Comments Section

```html
{{ if .Site.Config.extra.enable_comments }}
<div id="comments">
    <h3>Comments</h3>
    <!-- Your comment system integration -->
    <!-- Example: Giscus, Disqus, etc. -->
</div>
{{ end }}
```

### `_pagination.html` - Page Navigation

```html
{{ with .Page.Higher }}
<a href="{{ .Permalink }}" class="pagination-prev">
    ← {{ .Config.title }}
</a>
{{ end }}

{{ with .Page.Lower }}
<a href="{{ .Permalink }}" class="pagination-next">
    {{ .Config.title }} →
</a>
{{ end }}
```

## Naming Conventions

**Prefix with underscore:**
```
templates/partials/
├── _head.html       ✅ Good
├── _header.html     ✅ Good
├── head.html        ❌ Avoid (no underscore)
└── Header.html      ❌ Avoid (capitalized)
```

**Use descriptive names:**
```
_head.html           ✅ Clear purpose
_header.html         ✅ Clear purpose
_post_meta.html      ✅ Specific and clear
_thing.html          ❌ Vague
_1.html              ❌ Not descriptive
```

## Passing Data to Partials

### Pass Entire Context

```html
<!-- Pass full context (.) -->
{{ template "partials/_head.html" . }}

<!-- Partial can access all variables -->
<title>{{ .Page.Config.title }}</title>
```

### Pass Custom Dictionary

```html
<!-- Create custom data dict -->
{{ template "partials/_card.html" (dict 
    "title" .Page.Config.title 
    "description" .Page.Config.description
    "url" .Page.Permalink
) }}

<!-- In partial, access with . -->
<div class="card">
    <h3>{{ .title }}</h3>
    <p>{{ .description }}</p>
    <a href="{{ .url }}">Read more</a>
</div>
```

### Pass Specific Page

```html
<!-- Pass different page context -->
{{ range .Section.Children }}
    {{ template "partials/_post_card.html" . }}
{{ end }}

<!-- Partial receives individual page -->
<article class="post-card">
    <h2>{{ .Config.title }}</h2>
    <time>{{ date .Config.date "Jan 2, 2006" }}</time>
</article>
```

## Partial Organization

### By Function

```
templates/partials/
├── _head.html          # Document head
├── _header.html        # Site header
├── _footer.html        # Site footer
├── _scripts.html       # JavaScript
├── _analytics.html     # Analytics code
└── _comments.html      # Comments system
```

### By Component

```
templates/partials/
├── navigation/
│   ├── _main_nav.html
│   ├── _footer_nav.html
│   └── _breadcrumbs.html
├── content/
│   ├── _post_card.html
│   ├── _post_meta.html
│   └── _toc.html
└── social/
    ├── _sharing.html
    └── _author.html
```

## Best Practices

1. **Keep partials focused** - One responsibility per partial
2. **Use underscore prefix** - Clear identification as partial
3. **Document parameters** - Comment expected data structure
4. **Make reusable** - Design for multiple contexts
5. **Avoid logic in partials** - Keep complex logic in main templates
6. **Test with different data** - Ensure partials handle edge cases

## Advanced Patterns

### Conditional Partials

```html
<!-- Load different partials based on context -->
{{ if eq .Section.Name "blog" }}
    {{ template "partials/_blog_sidebar.html" . }}
{{ else if eq .Section.Name "notes" }}  
    {{ template "partials/_notes_sidebar.html" . }}
{{ else }}
    {{ template "partials/_default_sidebar.html" . }}
{{ end }}
```

### Nested Partials

```html
<!-- _head.html includes other partials -->
<head>
    {{ template "partials/_meta_tags.html" . }}
    {{ template "partials/_stylesheets.html" . }}
    {{ template "partials/_fonts.html" . }}
</head>
```

### Optional Partials

```html
<!-- Only include if feature enabled -->
{{ if .Site.Config.extra.enable_search }}
    {{ template "partials/_search.html" . }}
{{ end }}

{{ if .Page.Config.extra.show_author }}
    {{ template "partials/_author_bio.html" . }}
{{ end }}
```

---

**Related:**
- [Macros](/templates/macros)
- [Template Structure](/templates/structure)
- [Examples](/templates/examples)
