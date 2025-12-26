+++
title = "Gozzi"
template = "home.html"
date = 2025-12-15
+++

A simple static site generator written in Go. Built for learning, designed for clarity.

## What is Gozzi?

Gozzi converts Markdown files into a static website. I built it to understand how static site generators work, and now I use it for [tduyng.com](https://tduyng.com) with 100+ posts building in ~100ms.

Gozzi might be right for you if you:

- Want to learn Go by using a real tool
- Prefer simplicity over feature-rich complexity
- Like understanding how your tools work

## Quick Start

```bash
# Install
go install github.com/tduyng/gozzi@latest

# Create site
mkdir my-blog && cd my-blog
mkdir -p content templates static

# Configure (config.toml)
echo 'base_url = "https://example.com"
title = "My Blog"' > config.toml

# Create first post (content/hello.md)
echo '+++
title = "Hello World"
date = 2024-01-15
+++

# Hello World

My first post!' > content/hello.md

# Serve
gozzi serve  # Visit http://localhost:1313
```

## Core Features

- **Fast** - Sub-second builds, instant live reload
- **Simple** - Markdown + TOML config + Go templates
- **Complete** - Taxonomies, series, pagination, RSS, sitemap, math, diagrams built-in
- **Portable** - Single binary, deploy anywhere

## Next Steps

- [Installation](/installation) - Install Gozzi
- [Getting Started](/getting-started) - Build your first complete site
- [Real Example](/examples/real-world) - How I use Gozzi for tduyng.com

## Resources

[GitHub](https://github.com/tduyng/gozzi) · [Issues](https://github.com/tduyng/gozzi/issues) · [Discussions](https://github.com/tduyng/gozzi/discussions)
