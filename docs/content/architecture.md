+++
title = "Architecture"
date = 2025-12-15
template = "page.html"
+++

How Gozzi works internally.

## Build Process

```
1. Read config.toml
2. Scan content/ directory
3. Parse markdown files
4. Process templates
5. Generate HTML
6. Copy static files
7. Write to public/
```

**Key steps:**
- Parse TOML front matter
- Render markdown to HTML
- Apply Go templates
- Generate RSS/sitemap
- Copy assets

## Development Server

```
1. Start HTTP server (default: localhost:1313)
2. Build site initially
3. Watch for file changes
4. Rebuild on change
5. Send live reload signal
```

**Live reload:**
- Server-Sent Events (SSE)
- Watches: content/, templates/, static/, config.toml
- Automatic browser refresh

## Content Processing

```
content/blog/post.md
    ↓
1. Parse front matter (TOML)
2. Extract markdown body
3. Render markdown to HTML
4. Apply syntax highlighting (if code blocks)
5. Render math (if KaTeX expressions)
6. Render diagrams (if Mermaid blocks)
7. Generate table of contents
    ↓
HTML output
```

## Template System

Gozzi uses Go's `html/template` engine:

**Template selection:**
1. Check front matter `template` field
2. Check content type (post, page, section)
3. Fall back to defaults

**Template inheritance:**
```
base.html (layout)
    ↓
post.html (extends base)
    ↓
content (rendered here)
```

**Partials:**
```html
<!-- templates/base.html -->
{{ template "partials/_header.html" . }}
{{ block "content" . }}{{ end }}
{{ template "partials/_footer.html" . }}
```

## Performance

**Fast builds:**
- Concurrent processing
- Minimal dependencies
- Efficient markdown parsing
- Template caching

**Typical build times:**
- 10 pages: ~10-20ms
- 100 pages: ~50-100ms
- 1000 pages: ~500ms-1s

## File Structure

```
project/
├── config.toml       # Site configuration
├── content/          # Markdown files
│   ├── _index.md     # Homepage
│   └── blog/
│       ├── _index.md # Section
│       └── post.md   # Page
├── templates/        # Go templates
│   ├── base.html
│   ├── post.html
│   └── partials/
├── static/           # Static assets
│   ├── css/
│   ├── js/
│   └── img/
└── public/           # Generated (git ignored)
    └── index.html
```

## Technical Details

- **Language:** Go 1.25+
- **Markdown:** goldmark parser
- **Templates:** Go html/template
- **Math:** KaTeX (server-side)
- **Diagrams:** Mermaid (server-side)
- **Syntax:** Chroma highlighter
