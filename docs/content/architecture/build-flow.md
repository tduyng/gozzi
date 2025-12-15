+++
title = "Build Flow"
date = 2025-12-15
template = "page.html"
+++

How `gozzi build` generates your static site.

## Build Steps

### 1. Load Configuration

```
gozzi build --config config.toml
  ↓
Read config.toml
  ↓
Parse TOML into Site struct
```

### 2. Parse Content

```
Walk content/ directory
  ↓
For each .md file:
  - Extract front matter (TOML between +++)
  - Parse markdown to HTML
  - Build content tree
  ↓
Create ContentMap
```

### 3. Load Templates

```
Walk templates/ directory
  ↓
For each .html file:
  - Parse template
  - Register template functions
  ↓
Template cache ready
```

### 4. Generate Site

```
For each node in content tree:
  - Find matching template
  - Render with data
  - Write HTML file
  ↓
Generate sitemap.xml
Generate atom.xml  
Generate robots.txt
Copy static/ assets
```

## Content Tree

Your directory structure becomes a tree:

```
content/
├── _index.md         → Section (root)
├── blog/
│   ├── _index.md     → Section (blog)
│   └── post.md       → Page
└── about/
    └── _index.md     → Section (about)
```

## Concurrent Processing

Gozzi processes files in parallel:

```
Worker Pool (NumCPU workers)
  ↓
Process multiple files simultaneously
  ↓
Faster builds
```

## Clean vs Incremental

**Incremental (default)**
- Preserves existing files
- Only overwrites Gozzi-generated files
- Fast

**Clean (`--clean` flag)**
- Removes ALL files first
- Fresh build
- Slower but guaranteed clean
