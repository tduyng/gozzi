+++
title = "Gozzi"
template = "home.html"
date = 2025-12-15
+++
## Features

### 📝 Markdown First
Write content in Markdown with TOML front matter. Organize files in folders. That's it.

### ⚡ Fast & Live
Sub-second builds. Live reload that actually works. See your changes instantly.

### 🎨 Go Templates
Use Go's HTML templates. Simple, powerful, and you learn Go templating along the way.

### 🔧 Clear Config
One TOML file. No magic. You control everything.

### 📦 Built-in Essentials
Tags, pagination, RSS, sitemap, TOC. Common features just work.

### 🚀 Deploy Anywhere
Generates static HTML. Deploy to GitHub Pages, Netlify, Vercel, or anywhere.

## Quick Example

Install and create your first site:

```bash
# Install
go install github.com/tduyng/gozzi@latest

# Create Site
mkdir my-blog && cd my-blog
mkdir -p content/blog templates static

# Create config
cat > config.toml << EOF
base_url = "https://example.com"
title = "My Blog"
EOF
```

```bash
# Create first post
mkdir -p content/blog/hello-world
cat > content/blog/hello-world/index.md << EOF
+++
title = "Hello World"
date = 2024-01-15
+++

# My First Post

Welcome to my blog built with Gozzi!
EOF
```

```bash
# Serve
gozzi serve
# Visit http://localhost:1313
```

## Real Usage

I built Gozzi to learn, and now I use it daily for [tduyng.com](https://tduyng.com):

- 100+ blog posts and notes
- KaTeX math expressions
- Mermaid diagrams
- Full-text search
- Builds in ~100ms

See the [real-world example](/examples/real-world) to learn from my setup.

## Why Gozzi?

Gozzi is **not trying to replace Hugo or Zola**. Those are excellent, mature tools. Use them if they fit your needs!

Gozzi might be a good fit if you:

- **Want to learn Go** by using a real tool
- **Like simplicity** over feature overload
- **Prefer clear configuration** in TOML
- **Enjoy understanding** how your tools work

## Open Source

Gozzi is MIT licensed. Built by learning, shared for learning.

[View source on GitHub](https://github.com/tduyng/gozzi) · [Report issues](https://github.com/tduyng/gozzi/issues) · [Discussions](https://github.com/tduyng/gozzi/discussions)
