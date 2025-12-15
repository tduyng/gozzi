+++
title = "Advanced Template Patterns"
date = 2025-12-15
template = "page.html"
+++

Advanced techniques for building sophisticated, data-driven templates.

## Conditional Template Loading

Load different partials based on context:

```html
<!-- Load section-specific sidebars -->
{{ if eq .Section.Name "blog" }}
    {{ template "partials/_blog_sidebar.html" . }}
{{ else if eq .Section.Name "notes" }}  
    {{ template "partials/_notes_sidebar.html" . }}
{{ else }}
    {{ template "partials/_default_sidebar.html" . }}
{{ end }}

<!-- Feature flags from config -->
{{ if .Site.Config.extra.enable_comments }}
    {{ template "partials/_comments.html" . }}
{{ end }}

{{ if .Site.Config.extra.enable_analytics }}
    {{ template "partials/_analytics.html" . }}
{{ end }}
```

## Complex Data Processing

### Filtering and Transforming Collections

```html
<!-- Get published, featured posts from last 30 days -->
{{ $posts := get_section "blog" }}
{{ $published := where $posts.Children "Config.draft" false }}
{{ $featured := where $published "Config.featured" true }}
{{ $recent := first 5 (reverse (where $published "Config.draft" false)) }}

<section class="featured">
    <h2>Featured Posts</h2>
    {{ range $featured }}
        <article>
            <h3><a href="{{ .Permalink }}">{{ .Config.title }}</a></h3>
            <p>{{ .Config.description }}</p>
        </article>
    {{ end }}
</section>

<section class="recent">
    <h2>Recent Posts</h2>
    {{ range $recent }}
        <article>
            <h3><a href="{{ .Permalink }}">{{ .Config.title }}</a></h3>
            <time>{{ date .Config.date "Jan 2" }}</time>
        </article>
    {{ end }}
</section>
```

### Grouping and Aggregation

```html
<!-- Group posts by year -->
{{ $posts := .Section.Children }}
{{ $groups := group_by $posts "Config.date.Year" }}

{{ range $groups }}
    <section class="year-{{ .Key }}">
        <h2>{{ .Key }}</h2>
        {{ range .Items }}
            <article>
                <time>{{ date .Config.date "Jan 2" }}</time>
                <a href="{{ .Permalink }}">{{ .Config.title }}</a>
            </article>
        {{ end }}
    </section>
{{ end }}
```

### Multi-Level Filtering

```html
<!-- Complex query: published blog posts tagged "golang", sorted by date -->
{{ $blog := get_section "blog" }}
{{ $published := where $blog.Children "Config.draft" false }}
{{ $goPosts := where $published "Config.tags" "golang" }}
{{ $sorted := sort $goPosts ".Config.date" "desc" }}

<div class="golang-posts">
    <h2>Go Programming Posts</h2>
    {{ range $sorted }}
        <article>
            <h3><a href="{{ .Permalink }}">{{ .Config.title }}</a></h3>
            <time>{{ date .Config.date "January 2, 2006" }}</time>
        </article>
    {{ end }}
</div>
```

## Dynamic Content Assembly

### Building Navigation Dynamically

```html
<!-- Generate navigation from sections -->
<nav class="main-nav">
    {{ $sections := slice "blog" "notes" "projects" }}
    {{ range $sections }}
        {{ $section := get_section . }}
        <a href="/{{ . }}/" class="{{ if eq $.Page.Section . }}active{{ end }}">
            {{ $section.Config.title }}
            <span class="count">({{ len $section.Children }})</span>
        </a>
    {{ end }}
</nav>
```

### Related Content

```html
<!-- Find related posts by tags -->
{{ $currentTags := .Page.Config.tags }}
{{ $blog := get_section "blog" }}
{{ $related := slice }}

{{ range $blog.Children }}
    {{ if ne .Permalink $.Page.Permalink }}
        {{ range .Config.tags }}
            {{ if contains $currentTags . }}
                {{ $related = $related | append $ }}
            {{ end }}
        {{ end }}
    {{ end }}
{{ end }}

{{ if gt (len $related) 0 }}
    <aside class="related-posts">
        <h3>Related Posts</h3>
        {{ range first 3 $related }}
            <article>
                <a href="{{ .Permalink }}">{{ .Config.title }}</a>
            </article>
        {{ end }}
    </aside>
{{ end }}
```

## Performance Optimization

### Caching Expensive Operations

```html
<!-- Cache section data in variable -->
{{ $blog := get_section "blog" }}
{{ $notes := get_section "notes" }}

<!-- Reuse without re-fetching -->
<p>{{ len $blog.Children }} blog posts</p>
<p>{{ len $notes.Children }} notes</p>

{{ range first 5 (reverse $blog.Children) }}
    <!-- Use cached data -->
{{ end }}
```

### Lazy Loading Content

```html
<!-- Load related content only if needed -->
{{ if .Page.Config.extra.show_related }}
    {{ $related := where .Section.Children "Config.tags" .Page.Config.tags }}
    {{ range first 3 $related }}
        <article>{{ .Config.title }}</article>
    {{ end }}
{{ end }}
```

## Advanced Templating Techniques

### Template Composition

```html
<!-- Build complex layouts from smaller pieces -->
<div class="page-layout">
    <aside class="sidebar">
        {{ template "partials/_sidebar_about.html" . }}
        {{ template "partials/_sidebar_tags.html" . }}
        {{ template "partials/_sidebar_recent.html" . }}
    </aside>
    
    <main class="content">
        {{ block "content" . }}{{ end }}
    </main>
    
    <aside class="widgets">
        {{ template "partials/_widget_newsletter.html" . }}
        {{ template "partials/_widget_popular.html" . }}
    </aside>
</div>
```

### Recursive Templates

```html
<!-- Render nested comment threads -->
<!-- macros/comment_thread.html -->
<div class="comment" data-id="{{ .id }}">
    <div class="comment-body">
        <strong>{{ .author }}</strong>
        <p>{{ .text }}</p>
    </div>
    
    {{ if .replies }}
        <div class="comment-replies">
            {{ range .replies }}
                {{ template "macros/comment_thread.html" . }}
            {{ end }}
        </div>
    {{ end }}
</div>
```

### Template Middleware Pattern

```html
<!-- Wrap content with common elements -->
{{ define "with_sidebar" }}
    <div class="layout-with-sidebar">
        <main>{{ .content }}</main>
        <aside>{{ .sidebar }}</aside>
    </div>
{{ end }}

<!-- Usage -->
{{ template "with_sidebar" (dict 
    "content" .Page.Content 
    "sidebar" "Sidebar content here"
) }}
```

## State Management

### Maintaining State Across Templates

```html
<!-- Build state object -->
{{ $state := dict 
    "page" .Page 
    "site" .Site 
    "showTOC" .Page.Config.extra.toc
    "theme" (priority .Page.Config.extra.theme .Site.Config.extra.theme "light")
}}

<!-- Pass state to partials -->
{{ template "partials/_article.html" $state }}
```

### Context Preservation

```html
<!-- Preserve parent context in loops -->
{{ range .Section.Children }}
    <!-- $ refers to root context -->
    <article>
        <h2>{{ .Config.title }}</h2>
        <p>Site: {{ $.Site.Config.title }}</p>
        <p>Current section: {{ $.Section.Name }}</p>
    </article>
{{ end }}
```

## Error Handling

### Graceful Degradation

```html
<!-- Handle missing data gracefully -->
{{ with .Page.Config.extra.author }}
    <p>By {{ . }}</p>
{{ else }}
    <p>By {{ .Site.Config.extra.name }}</p>
{{ end }}

<!-- Provide fallback images -->
{{ $image := priority .Page.Config.img .Section.Config.img "/img/default.jpg" }}
<img src="{{ $image }}" alt="{{ .Page.Config.title }}">
```

### Validation

```html
<!-- Validate required data -->
{{ if and .Page.Config.title .Page.Content }}
    <article>
        <h1>{{ .Page.Config.title }}</h1>
        {{ .Page.Content }}
    </article>
{{ else }}
    <div class="error">
        <p>Missing required page data</p>
    </div>
{{ end }}
```

## Debugging Patterns

### Development Mode Features

```html
{{ if eq .Site.Environment "development" }}
    <div id="debug-toolbar">
        <h3>Debug Info</h3>
        <dl>
            <dt>Template:</dt>
            <dd>{{ .TemplateName }}</dd>
            
            <dt>Section:</dt>
            <dd>{{ .Section.Name }}</dd>
            
            <dt>Page Type:</dt>
            <dd>{{ .Page.Type }}</dd>
        </dl>
        
        <details>
            <summary>Full Page Data</summary>
            <pre>{{ .Page | json }}</pre>
        </details>
    </div>
{{ end }}
```

### Performance Monitoring

```html
<!-- Track template rendering time -->
{{ $start := now }}

<!-- Complex template logic -->
{{ range .Section.Children }}
    <!-- ... -->
{{ end }}

{{ $duration := sub (now) $start }}
{{ if eq .Site.Environment "development" }}
    <p>Rendered in {{ $duration }}ms</p>
{{ end }}
```

## Advanced Use Cases

### Multi-Language Support

```html
<!-- Language-aware content -->
{{ $lang := priority .Page.Config.lang .Site.Config.lang "en" }}

<html lang="{{ $lang }}">
    <head>
        {{ if eq $lang "en" }}
            <title>{{ .Page.Config.title }}</title>
        {{ else if eq $lang "fr" }}
            <title>{{ .Page.Config.title_fr }}</title>
        {{ end }}
    </head>
</html>
```

### A/B Testing

```html
<!-- Conditional variants -->
{{ $variant := .Site.Config.extra.ab_variant | default "A" }}

{{ if eq $variant "A" }}
    <div class="hero-variant-a">
        <h1>Welcome to Gozzi</h1>
    </div>
{{ else }}
    <div class="hero-variant-b">
        <h1>Build Fast Static Sites</h1>
    </div>
{{ end }}
```

### Dynamic Schema Generation

```html
<!-- Generate JSON-LD schema -->
<script type="application/ld+json">
{
    "@context": "https://schema.org",
    "@type": "BlogPosting",
    "headline": "{{ .Page.Config.title }}",
    "datePublished": "{{ .Page.Config.date.Format "2006-01-02" }}",
    "author": {
        "@type": "Person",
        "name": "{{ .Site.Config.extra.name }}"
    }
    {{ if .Page.Config.img }},
    "image": "{{ .Page.Config.img_url }}"
    {{ end }}
}
</script>
```

---

**Related:**
- [Template Functions](/templates/functions)
- [Template Variables](/templates/variables)
- [Examples](/templates/examples)
