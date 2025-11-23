# Template Variables

Each template receives a context object (`.`) with comprehensive site and content data. Understanding these variables is key to building dynamic templates.

## Context Object

The root context (`.`) contains:

- **`.Site`** - Global site configuration
- **`.Section`** - Section data (list pages with `_index.md`)
- **`.Page`** - Individual page data
- **`.Toc`** - Table of contents (if available)

## `.Site` - Global Configuration

Site-wide data from `config.toml`:

### Basic Site Info

```html
<!-- Core configuration -->
<title>{{ .Site.Config.title }}</title>
<meta name="description" content="{{ .Site.Config.description }}">
<link rel="canonical" href="{{ .Site.Config.base_url }}">
<html lang="{{ .Site.Config.language }}">

<!-- Feed URL -->
<link rel="alternate" type="application/atom+xml" href="{{ .Site.Config.feed_url }}">
```

### Extended Configuration

Access custom `[extra]` data:

```html
<!-- Personal info -->
<p>{{ .Site.Config.extra.name }}</p>
<p>{{ .Site.Config.extra.bio }}</p>
<img src="{{ .Site.Config.extra.avatar }}" alt="Profile">

<!-- Navigation sections -->
{{ range .Site.Config.extra.sections }}
  <a href="{{ .path }}">{{ .name }}</a>
{{ end }}

<!-- Social links -->
{{ range .Site.Config.extra.links }}
  <a href="{{ .url }}" title="{{ .name }}">
    <img src="/icon/{{ .icon }}.svg" alt="{{ .name }}">
  </a>
{{ end }}
```

## `.Section` - Section Data

Available in section templates (list pages):

### Section Metadata

```html
<!-- Section configuration -->
<h1>{{ .Section.Config.title }}</h1>
<p>{{ .Section.Config.description }}</p>
<time>{{ date .Section.Config.date "January 2, 2006" }}</time>

<!-- Section content (from _index.md body) -->
<div class="section-intro">
  {{ .Section.Content }}
</div>
```

### Section Custom Data

```html
<!-- Extended section data -->
{{ if .Section.Config.extra.hero_text }}
  <div class="hero">{{ .Section.Config.extra.hero_text }}</div>
{{ end }}

{{ if .Section.Config.extra.show_reading_time }}
  <!-- Show reading time for posts -->
{{ end }}
```

### Children Pages

```html
<!-- List all pages in section -->
{{ range .Section.Children }}
  <article>
    <h2><a href="{{ .Permalink }}">{{ .Config.title }}</a></h2>
    <time>{{ date .Config.date "January 2, 2006" }}</time>
    <p>{{ .Config.description }}</p>
    
    {{ if .Config.featured }}
      <span class="featured">Featured</span>
    {{ end }}
    
    <div class="tags">
      {{ range .Config.tags }}
        <span class="tag">#{{ . }}</span>
      {{ end }}
    </div>
  </article>
{{ end }}
```

### Section Name

```html
<!-- Check which section we're in -->
{{ if eq .Section.Name "blog" }}
  <div class="blog-specific">Blog content</div>
{{ else if eq .Section.Name "notes" }}
  <div class="notes-specific">Notes content</div>
{{ end }}
```

## `.Page` - Individual Page Data

Used in page templates:

### Page Metadata

```html
<!-- Title and description -->
<h1>{{ .Page.Config.title }}</h1>
<meta name="description" content="{{ .Page.Config.description }}">

<!-- Dates -->
<time datetime="{{ .Page.Config.date.Format "2006-01-02" }}">
  {{ date .Page.Config.date "January 2, 2006" }}
</time>

{{ if not .Page.Config.updated.IsZero }}
  <p>Updated: {{ date .Page.Config.updated "January 2, 2006" }}</p>
{{ end }}
```

### Page Content

```html
<!-- Rendered markdown content -->
<div class="content">
  {{ .Page.Content }}
</div>

<!-- Permalink -->
<link rel="canonical" href="{{ .Page.Permalink }}">
<a href="{{ .Page.Permalink }}">Permalink</a>
```

### Page URLs

```html
<!-- Full URL -->
<meta property="og:url" content="{{ .Page.URL }}">

<!-- Relative path -->
<a href="{{ .Page.Permalink }}">{{ .Page.Config.title }}</a>
```

### Page Navigation

```html
<!-- Previous/Next pages in section -->
{{ with .Page.Higher }}
  <a href="{{ .Permalink }}" class="prev">
    ← {{ .Config.title }}
  </a>
{{ end }}

{{ with .Page.Lower }}
  <a href="{{ .Permalink }}" class="next">
    {{ .Config.title }} →
  </a>
{{ end }}
```

### Page Configuration

```html
<!-- Template used -->
{{ if .Page.Config.template }}
  <!-- Using: {{ .Page.Config.template }} -->
{{ end }}

<!-- Draft status -->
{{ if .Page.Config.draft }}
  <div class="draft-warning">This is a draft</div>
{{ end }}

<!-- Featured status -->
{{ if .Page.Config.featured }}
  <span class="badge">Featured</span>
{{ end }}

<!-- Language -->
<html lang="{{ .Page.Config.language }}">
```

### Tags and Categorization

```html
<!-- Page tags -->
{{ if .Page.Config.tags }}
  <div class="tags">
    {{ range .Page.Config.tags }}
      <a href="/tags/{{ . | urlize }}" class="tag">#{{ . }}</a>
    {{ end }}
  </div>
{{ end }}

<!-- Check for specific tag -->
{{ if contains .Page.Config.tags "golang" }}
  <span class="go-related">Go-related post</span>
{{ end }}
```

### Page Images

```html
<!-- Featured/cover image -->
{{ if .Page.Config.img }}
  <img src="{{ .Page.Config.img }}" alt="{{ .Page.Config.title }}">
  <meta property="og:image" content="{{ .Page.Config.img_url }}">
{{ end }}
```

### Custom Page Data

```html
<!-- Extended page data -->
{{ if .Page.Config.extra.toc }}
  <aside class="toc">{{ .Toc }}</aside>
{{ end }}

{{ if .Page.Config.extra.reading_time }}
  <span class="reading-time">{{ .Page.Config.extra.reading_time }}</span>
{{ end }}

{{ if .Page.Config.extra.author }}
  <p>By {{ .Page.Config.extra.author }}</p>
{{ end }}
```

## `.Toc` - Table of Contents

Automatically generated from page headings:

```html
<!-- Basic TOC display -->
{{ if .Toc }}
  <nav class="toc">
    <h3>Table of Contents</h3>
    {{ .Toc }}
  </nav>
{{ end }}

<!-- Conditional TOC from config -->
{{ if .Page.Config.extra.toc }}
  <aside class="toc">
    {{ .Toc }}
  </aside>
{{ end }}
```

## Priority Function for Inheritance

Cascade values from page → section → site:

```html
<!-- Use page title, fall back to section, then site -->
<title>{{ priority .Page.Config.title .Section.Config.title .Site.Config.title }}</title>

<!-- Language inheritance -->
<html lang="{{ priority .Page.Config.lang .Section.Config.lang .Site.Config.lang }}">

<!-- Description with fallbacks -->
{{ $description := priority .Page.Config.description .Section.Config.description .Site.Config.description }}
<meta name="description" content="{{ $description }}">

<!-- Image with fallbacks -->
{{ $image := priority .Page.Config.img .Section.Config.img .Site.Config.img }}
{{ if $image }}
  <meta property="og:image" content="{{ $image }}">
{{ end }}
```

## Common Variable Patterns

### Navigation State

```html
<!-- Highlight current section -->
{{ range .Site.Config.extra.sections }}
  {{ $active := "" }}
  {{ if eq .path $.Page.Section }}
    {{ $active = "active" }}
  {{ end }}
  <a href="{{ .path }}" class="{{ $active }}">{{ .name }}</a>
{{ end }}
```

### Conditional Content

```html
<!-- Show content based on data availability -->
{{ with .Page.Config.description }}
  <meta name="description" content="{{ . }}">
{{ end }}

{{ with .Page.Config.extra.author }}
  <p>Written by {{ . }}</p>
{{ end }}
```

### Date Formatting

```html
<!-- Various date formats -->
{{ date .Page.Config.date "January 2, 2006" }}
{{ date .Page.Config.date "2006-01-02" }}
{{ date .Page.Config.date "Jan 2" }}
{{ date .Page.Config.date "Monday, Jan 2, 2006" }}
```

### Array Iteration

```html
<!-- Iterate with index -->
{{ range $index, $item := .Section.Children }}
  <div class="item-{{ $index }}">
    <h2>{{ $item.Config.title }}</h2>
  </div>
{{ end }}

<!-- Check array length -->
{{ if gt (len .Page.Config.tags) 0 }}
  <p>Tagged with {{ len .Page.Config.tags }} tags</p>
{{ end }}
```

## Variable Scope

### Accessing Parent Context

```html
<!-- Inside range, access parent context with $ -->
{{ range .Section.Children }}
  <!-- .Site not available here, use $.Site -->
  <article>
    <h2>{{ .Config.title }}</h2>
    <p>Site: {{ $.Site.Config.title }}</p>
  </article>
{{ end }}
```

### Storing Variables

```html
<!-- Assign to variable for reuse -->
{{ $siteTitle := .Site.Config.title }}
{{ $pageTitle := .Page.Config.title }}
{{ $fullTitle := printf "%s | %s" $pageTitle $siteTitle }}
<title>{{ $fullTitle }}</title>

<!-- Complex data -->
{{ $posts := where .Section.Children "Config.draft" false }}
{{ $featured := where $posts "Config.featured" true }}
<p>{{ len $featured }} featured posts out of {{ len $posts }}</p>
```

## Debugging Variables

```html
<!-- Dump all page data (development only) -->
{{ if eq .Site.Environment "development" }}
  <details>
    <summary>Page Data (Debug)</summary>
    <pre>{{ .Page | json }}</pre>
  </details>
  
  <details>
    <summary>Section Data (Debug)</summary>
    <pre>{{ .Section | json }}</pre>
  </details>
{{ end }}
```

---

**Related:**
- [Template Functions](/guide/templates/functions)
- [Configuration](/guide/configuration/)
- [Template Structure](/guide/templates/structure)
