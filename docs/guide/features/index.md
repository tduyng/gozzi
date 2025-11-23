# Built-in Features

Gozzi includes comprehensive features that enhance your static site without requiring additional plugins.

## Content Enhancement

- **[Mathematical Expressions](/guide/features/math)** - Server-side KaTeX rendering for beautiful math equations
- **[Mermaid Diagrams](/guide/features/diagrams)** - Client-side diagram rendering with automatic script injection
- **[Syntax Highlighting](/guide/features/syntax-highlighting)** - Server-side code highlighting with Chroma

## Content Organization

- **[Table of Contents](/guide/features/content-features#table-of-contents)** - Automatic TOC generation from headings
- **[Tag Management](/guide/features/content-features#tag-management)** - Comprehensive tag system across all content
- **[Pagination](/guide/features/content-features#pagination)** - Built-in pagination for large content collections
- **[Content Analytics](/guide/features/content-features#content-analytics)** - Word count and reading time calculation

## SEO & Discoverability

- **[RSS/Atom Feeds](/guide/features/seo#rss-atom-feeds)** - Automatic feed generation
- **[Sitemap](/guide/features/seo#sitemap-generation)** - XML sitemap for search engines
- **[Robots.txt](/guide/features/seo#robots-txt)** - Customizable robots file
- **[Social Media Integration](/guide/features/seo#social-media-integration)** - Open Graph meta tags

## Performance

All built-in features are optimized for speed:

- **Server-Side Rendering** - KaTeX and syntax highlighting rendered during build time
- **Client-Side Rendering** - Mermaid diagrams render in browser for interactivity
- **Fast Builds** - Native features add minimal overhead (still sub-second for most sites)
- **Caching** - Generated content cached until source changes
- **Minimal JS** - Only Mermaid requires JavaScript (~200KB); math and code work without JS

## Feature Integration

Built-in features work seamlessly together:

- Tag pages include pagination automatically
- RSS feeds respect tag filtering
- Social meta tags use calculated reading times
- TOC generation works with KaTeX math, Mermaid diagrams, and syntax highlighting
- Server-side features (KaTeX, syntax highlighting) have zero runtime overhead
- Client-side features (Mermaid) load asynchronously without blocking page render
- All features respect configuration inheritance

---

**Next Steps:**
- Learn about [Mathematical Expressions](/guide/features/math)
- Explore [Mermaid Diagrams](/guide/features/diagrams)
- Configure [SEO features](/guide/features/seo)
