# Gozzi 🍃

**A Go static site generator born from curiosity**

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**Gozzi** started as a learning project to understand modern static site generators. What began as "let's see how Hugo, Zola works" became my daily driver for [tduyng.com](https://tduyng.com)

> 📖 Documentation: Full guides on installation, configuration, content structure, templates, helper functions, features, and CLI usage live in the docs/ folder.

## What it does

Gozzi addresses practical needs for my website:

- Supports Mermaid diagrams, KaTeX for math notation, and syntax highlighting. Automatically generates a table of contents, manages tags, and handles pagination.
- Uses Go's HTML templates with custom helper functions (see [htmlfunc.go](github.com/tduyng/gozzi/blob/main/internal/generator/htmlfunc.go)) for tasks like date formatting and image processing.
- Generates RSS feeds, sitemap.xml, robots.txt, and Open Graph meta tags for social sharing.
- Offers rapid build times and live reloading during content authoring. Use gozzi serve to preview changes instantly.

## Try it out

### Installation

```bash
go install github.com/tduyng/gozzi@latest
```

### Basic setup

1. Create content folder

```bash
mkdir -p my-site/content/blog
```

2. Add config (full example [here](https://github.com/tduyng/tduyng.github.io/blob/main/config.toml))

```toml
base_url = "https://tduyng.com"
title = "My site"
```

See more in my website: [tduyng.com](https://github.com/tduyng/tduyng.github.io)

4. Build & preview

```bash
gozzi build
gozzi serve
```

## Current state

This is actively developed. While it's stable enough for my needs, some APIs might change as I keep improving it.
