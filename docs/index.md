---
layout: home

hero:
  name: "Gozzi"
  text: "A Go Static Site Generator"
  tagline: "Built for learning. Simple by design."
  image:
    src: /logo.svg
    alt: Gozzi
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/tduyng/gozzi

features:
  - icon: 📝
    title: Markdown First
    details: Write content in Markdown with TOML front matter. Organize files in folders. That's it.
    
  - icon: ⚡
    title: Fast & Live
    details: Sub-second builds. Live reload that actually works. See your changes instantly.
    
  - icon: 🎨
    title: Go Templates
    details: Use Go's HTML templates. Simple, powerful, and you learn Go templating along the way.
    
  - icon: 🔧
    title: Clear Config
    details: One TOML file. No magic. You control everything.
    
  - icon: 📦
    title: Built-in Essentials
    details: Tags, pagination, RSS, sitemap, TOC. Common features just work.
    
  - icon: 🚀
    title: Deploy Anywhere
    details: Generates static HTML. Deploy to GitHub Pages, Netlify, Vercel, or anywhere.

---

## Quick Example

Install and create your first site:

::: code-group

```bash [Install]
go install github.com/tduyng/gozzi@latest
```

```bash [Create Site]
mkdir my-blog && cd my-blog
mkdir -p content/blog templates static

# Create config
cat > config.toml << EOF
base_url = "https://example.com"
title = "My Blog"
EOF

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

```bash [Create Template]
cat > templates/post.html << EOF
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Page.Config.title }}</title>
</head>
<body>
    <article>
        <h1>{{ .Page.Config.title }}</h1>
        <time>{{ date .Page.Config.date "Jan 2, 2006" }}</time>
        {{ .Page.Content }}
    </article>
</body>
</html>
EOF
```

```bash [Serve]
gozzi serve
# Visit http://localhost:1313
```

:::

## Real Usage

I built Gozzi to learn, and now I use it daily for [tduyng.com](https://tduyng.com):

- 100+ blog posts and notes
- Math expressions with KaTeX
- Diagrams with Mermaid
- Full-text search
- Builds in ~500ms

See the [real-world example](/examples/real-world) to learn from my setup.

## Not Another Hugo

Gozzi isn't trying to replace Hugo, Zola, or Jekyll. Those are mature, battle-tested tools.

**Use Gozzi if you want to:**
- Learn how static site generators work
- Use a tool built with learning in mind
- Understand every part of your site generation
- Keep things simple and explicit

**Use Hugo/Zola if you need:**
- Production-critical stability
- Extensive plugin ecosystem
- Advanced theme marketplace
- Enterprise support

## Documentation

- **[Getting Started](/guide/getting-started)** - Build your first site
- **[Installation](/guide/installation)** - Platform-specific install
- **[Configuration](/guide/configuration)** - Configure your site
- **[Templates](/guide/templates)** - Create layouts
- **[CLI Reference](/reference/cli)** - Command-line tools

## Open Source

Gozzi is MIT licensed. Built by learning, shared for learning.

[View source on GitHub](https://github.com/tduyng/gozzi) · [Report issues](https://github.com/tduyng/gozzi/issues) · [Discussions](https://github.com/tduyng/gozzi/discussions)
