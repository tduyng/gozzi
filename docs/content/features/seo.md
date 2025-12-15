+++
title = "SEO Automation"
date = 2025-12-15
template = "page.html"
+++

Gozzi automatically handles essential SEO tasks to improve your site's search engine visibility.

## Sitemap Generation

Automatically generates XML sitemaps following the sitemap protocol.

### Features

- **Location**: `/sitemap.xml`
- **Format**: XML sitemap protocol
- **Content**: All pages with modification dates
- **Priority**: Calculated based on content type
- **Change frequency**: Based on content updates

### Configuration

```toml
[sitemap]
enabled = true
change_freq = "weekly"
priority = 0.5
```

### Customization

Override defaults per page using front matter:

```yaml
---
title: 'Important Page'
sitemap:
    priority: 1.0
    change_freq: 'daily'
---
```

### Exclusion

Exclude pages from sitemap:

```yaml
---
title: 'Draft Page'
sitemap:
    exclude: true
---
```

## Robots.txt

Automatically generates robots.txt for crawler control.

### Default Configuration

- **Location**: `/robots.txt`
- **Content**: Default rules for common crawlers
- **Sitemap**: Automatically includes sitemap URL

### Generated File

```text
User-agent: *
Allow: /

Sitemap: https://yoursite.com/sitemap.xml
```

### Customization

Override by creating `static/robots.txt`:

```text
User-agent: *
Allow: /
Disallow: /admin/
Disallow: /private/

User-agent: Googlebot
Crawl-delay: 0

Sitemap: https://yoursite.com/sitemap.xml
```

## Meta Tags

### Open Graph

Essential for social media sharing:

```html
<meta property="og:title" content="{{ .Title }}" />
<meta property="og:description" content="{{ .Description }}" />
<meta property="og:type" content="article" />
<meta property="og:url" content="{{ .Permalink }}" />
<meta property="og:image" content="{{ .FeaturedImage }}" />
<meta property="og:site_name" content="{{ .Site.Title }}" />
```

### Twitter Cards

Optimized display on Twitter:

```html
<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:title" content="{{ .Title }}" />
<meta name="twitter:description" content="{{ .Description }}" />
<meta name="twitter:image" content="{{ .FeaturedImage }}" />
<meta name="twitter:site" content="@yourhandle" />
```

### Basic Meta Tags

```html
<meta name="description" content="{{ .Description }}" />
<meta name="keywords" content="{{ range .Tags }}{{ . }},{{ end }}" />
<meta name="author" content="{{ .Site.Author }}" />
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
        "headline": "{{ .Title }}",
        "description": "{{ .Description }}",
        "author": {
            "@type": "Person",
            "name": "{{ .Site.Author }}"
        },
        "datePublished": "{{ .Date.Format "2006-01-02" }}",
        "dateModified": "{{ .UpdatedDate.Format "2006-01-02" }}",
        "image": "{{ .FeaturedImage }}"
    }
</script>
```

### JSON-LD for Website

```html
<script type="application/ld+json">
    {
        "@context": "https://schema.org",
        "@type": "WebSite",
        "name": "{{ .Site.Title }}",
        "url": "{{ .Site.BaseURL }}",
        "description": "{{ .Site.Description }}"
    }
</script>
```

## Canonical URLs

Prevent duplicate content issues:

```html
<link rel="canonical" href="{{ .Permalink }}" />
```

## Performance SEO

### Build-Time Optimization

- **Server-side rendering** - All content pre-rendered
- **Zero JavaScript** - Core content accessible without JS
- **Fast page loads** - Critical for search rankings
- **Mobile-friendly** - Responsive by default

### Best Practices

1. **Use descriptive titles**: <code v-pre>{{ .Title }} | {{ .Site.Title }}</code>
2. **Write meta descriptions**: 150-160 characters optimal
3. **Optimize images**: Use appropriate formats and sizes
4. **Use semantic HTML**: Proper heading hierarchy
5. **Add alt text**: Images accessible and indexable

## Front Matter for SEO

Configure SEO per page:

```yaml
---
title: 'Complete Guide to Go Web Development'
description: 'Learn to build modern web applications with Go, covering HTTP servers, templating, databases, and deployment.'
keywords: ['golang', 'web development', 'tutorial', 'backend']
featured_image: '/images/go-web-guide.jpg'
date: 2024-01-15
updated: 2024-01-20
author: 'Your Name'
sitemap:
    priority: 0.9
    change_freq: 'monthly'
---
```

## Analytics Integration

Add analytics tracking (Google Analytics, Plausible, etc.) in your template:

```html
<!-- Google Analytics -->
{{ if .Site.GoogleAnalytics }}
<script async src="https://www.googletagmanager.com/gtag/js?id={{ .Site.GoogleAnalytics }}"></script>
<script>
    window.dataLayer = window.dataLayer || [];
    function gtag() {
        dataLayer.push(arguments);
    }
    gtag('js', new Date());
    gtag('config', '{{ .Site.GoogleAnalytics }}');
</script>
{{ end }}
```

Configuration:

```toml
[site]
google_analytics = "G-XXXXXXXXXX"
```

