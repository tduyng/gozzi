+++
title = "SEO Automation"
date = 2025-12-15
template = "page.html"
+++

Gozzi automatically handles essential SEO tasks to improve your site's search engine visibility.

## Sitemap Generation

Gozzi automatically generates XML sitemaps for all your content.

### Features

- **Location**: `/sitemap.xml` (automatically generated)
- **Format**: XML sitemap protocol
- **Content**: All pages with modification dates
- **Automatic**: No configuration needed

The sitemap includes:

- Homepage
- All content pages (blog posts, notes, etc.)
- Last modification date from frontmatter `updated` field
- Proper URL structure based on `base_url`

## Robots.txt

Gozzi automatically generates a basic robots.txt file.

### Generated File

- **Location**: `/robots.txt` (automatically generated)
- **Content**: Allows all crawlers and links to sitemap
- **Format**:

```text
User-agent: *
Allow: /

Sitemap: https://yoursite.com/sitemap.xml
```

### Customization

To customize robots.txt, create your own file at `static/robots.txt`. This will override the automatic generation:

```text
User-agent: *
Allow: /
Disallow: /admin/
Disallow: /drafts/

Sitemap: https://yoursite.com/sitemap.xml
```

## Meta Tags

### Open Graph

Essential for social media sharing:

```html
<meta property="og:title" content="{{ .Page.Config.title }}" />
<meta property="og:description" content="{{ .Page.Config.description }}" />
<meta property="og:type" content="article" />
<meta property="og:url" content="{{ .Site.base_url }}{{ .Page.Path }}" />
{{ with .Page.Config.img }}<meta property="og:image" content="{{ . }}" />{{ end }}
<meta property="og:site_name" content="{{ .Site.title }}" />
```

### Twitter Cards

Optimized display on Twitter:

```html
<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:title" content="{{ .Page.Config.title }}" />
<meta name="twitter:description" content="{{ .Page.Config.description }}" />
{{ with .Page.Config.img }}<meta name="twitter:image" content="{{ . }}" />{{ end }} {{ with
.Site.Extra.twitter }}<meta name="twitter:site" content="@{{ . }}" />{{ end }}
```

### Basic Meta Tags

```html
<meta name="description" content="{{ .Page.Config.description }}" />
<meta name="keywords" content="{{ range .Page.Config.tags }}{{ . }},{{ end }}" />
{{ with .Site.Extra.author }}<meta name="author" content="{{ . }}" />{{ end }}
<meta name="robots" content="index, follow" />
```

## Structured Data

### JSON-LD for Articles

Improve search result appearance:

```html
<script type="application/ld+json">
    {
        "@context": "https://schema.org",
        "@type": "BlogPosting",
        "headline": "{{ .Page.Config.title }}",
        "description": "{{ .Page.Config.description }}",
        {{ with .Site.Extra.author }}"author": {
            "@type": "Person",
            "name": "{{ . }}"
        },{{ end }}
        "datePublished": "{{ format_date .Page.Config.date "2006-01-02" }}",
        {{ if .Page.Config.updated }}"dateModified": "{{ format_date .Page.Config.updated "2006-01-02" }}",{{ end }}
        {{ with .Page.Config.img }}"image": "{{ . }}"{{ end }}
    }
</script>
```

### JSON-LD for Website

```html
<script type="application/ld+json">
    {
        "@context": "https://schema.org",
        "@type": "WebSite",
        "name": "{{ .Site.title }}",
        "url": "{{ .Site.base_url }}",
        "description": "{{ .Site.description }}"
    }
</script>
```

## Canonical URLs

Prevent duplicate content issues:

```html
<link rel="canonical" href="{{ .Site.base_url }}{{ .Page.Path }}" />
```

## Analytics Integration

Add analytics tracking in your template using custom `[extra]` config:

```html
<!-- Google Analytics -->
{{ with .Site.Extra.google_analytics }}
<script async src="https://www.googletagmanager.com/gtag/js?id={{ . }}"></script>
<script>
    window.dataLayer = window.dataLayer || []
    function gtag() {
        dataLayer.push(arguments)
    }
    gtag('js', new Date())
    gtag('config', '{{ . }}')
</script>
{{ end }}
```

Configuration:

```toml
[extra]
google_analytics = "G-XXXXXXXXXX"
```
