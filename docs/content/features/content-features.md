+++
title = "Built-in Features"
date = 2025-12-15
template = "page.html"
+++

Gozzi includes essential features for static sites:

## 📝 Content Features

**Tags & Organization:**
- Automatic tag pages (`/tags/golang/`)
- Tag cloud generation
- Categorization and filtering

**Table of Contents:**
- Auto-generated from headings
- Nested structure with anchor links
- Available as `.Toc` in templates

**RSS/Atom Feeds:**
- Automatic feed generation
- Per-section feeds
- Control with `generate_feed` option

**Sitemap:**
- Auto-generated XML sitemap
- SEO-optimized
- Includes all pages and sections

## 🎨 Content Rendering

**[Shortcodes](/features/shortcodes):**
- Hugo-compatible syntax
- Reusable HTML components in markdown
- Self-closing and paired variants
- Access to template functions
- Clean, maintainable content

**[Syntax Highlighting](/features/syntax-highlighting):**
- 100+ languages supported
- Multiple themes (GitHub, Monokai, etc.)
- Server-side rendering for performance

**[Math Expressions](/features/math):**
- KaTeX rendering (server-side)
- Inline: `$E=mc^2$`
- Block: `$$...$$`
- Fast, beautiful math typography

**[Mermaid Diagrams](/features/diagrams):**
- Flowcharts, sequence diagrams
- Gantt charts, pie charts
- Server-side rendering
- Clean SVG output

## 🔍 SEO & Meta

**[SEO Automation](/features/seo):**
- Open Graph tags
- Twitter Cards
- Meta descriptions
- Structured data support

**Automatic Generation:**
- Canonical URLs
- Robots.txt support
- Sitemap.xml
- RSS feeds

## 🚀 Performance

**Fast Builds:**
- Concurrent processing
- Sub-second builds for small sites
- Template caching

**Optimized Output:**
- Clean HTML
- Minified when needed
- Static assets copied efficiently

## 📦 Built-in, No Plugins

All features work out of the box:
- No plugins to install
- No configuration needed (sensible defaults)
- Works immediately after installation

---

**Learn More:**
- [Syntax Highlighting](/features/syntax-highlighting)
- [Math Expressions](/features/math)
- [Mermaid Diagrams](/features/diagrams)
- [SEO](/features/seo)
