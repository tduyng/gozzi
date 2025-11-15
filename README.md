# Gozzi 🍃

**A Go static site generator born from curiosity**

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/tduyng/gozzi)](https://github.com/tduyng/gozzi/releases)

I built **Gozzi** to learn how static site generators work. Now I use it for [tduyng.com](https://tduyng.com).

Gozzi is simple, fast, and helps you learn Go along the way. It's not trying to replace Hugo or Zola - use those if you need them!

## ✨ Key Features

- **🚀 Fast & Live Reload**: Instant preview with `gozzi serve` and sub-second builds
- **📊 Rich Content**: Mermaid diagrams, KaTeX math, syntax highlighting, and auto-generated TOC
- **🏷️ Smart Organization**: Automatic tag management, pagination, and content grouping
- **🌐 SEO Ready**: Auto-generated RSS feeds, sitemaps, robots.txt, and Open Graph meta tags
- **🎨 Flexible Templates**: Go HTML templates with 40+ custom helper functions
- **📱 Developer Friendly**: File watching, live reload, and comprehensive CLI tools

## 🚀 Quick Start

### 1. Install Gozzi

```bash
# Using Go
go install github.com/tduyng/gozzi@latest

# Or download from releases
curl -L https://github.com/tduyng/gozzi/releases/latest/download/gozzi-linux-amd64 -o gozzi
chmod +x gozzi
```

### 2. Create Your Site

```bash
# Create directory structure
mkdir my-site && cd my-site
mkdir -p content/blog static templates

# Create basic config
cat > config.toml << EOF
base_url = "https://example.com"
title = "My Awesome Site"
description = "A site built with Gozzi"

[author]
name = "Your Name"
email = "you@example.com"
EOF
```

### 3. Add Content

```bash
# Create your first post
cat > content/blog/hello-world.md << EOF
---
title: "Hello, World!"
date: 2024-01-15T10:00:00Z
tags: ["introduction", "blog"]
---

# Welcome to my site!

This is my first post built with **Gozzi**.

## Features I love:

- Fast builds
- Live reload
- Rich content support
EOF
```

### 4. Create Templates

```bash
# Basic post template
cat > templates/post.html << EOF
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }} - {{ .Site.Config.title }}</title>
</head>
<body>
    <article>
        <h1>{{ .Title }}</h1>
        <time>{{ date .Date "January 2, 2006" }}</time>
        <div>{{ .Content }}</div>
    </article>
</body>
</html>
EOF
```

### 5. Build & Preview

```bash
# Start development server with live reload
gozzi serve

# Or build for production
gozzi build
```

Visit `http://localhost:3000` to see your site!

## 📚 Documentation

📖 **[Complete Documentation](https://tduyng.github.io/gozzi/)** - Visit our full documentation site

Quick links:
- 🚀 [Getting Started Guide](https://tduyng.github.io/gozzi/guide/getting-started)
- ⚙️ [Installation](https://tduyng.github.io/gozzi/guide/installation)
- 📝 [Configuration](https://tduyng.github.io/gozzi/guide/configuration)
- 🎨 [Templates](https://tduyng.github.io/gozzi/guide/templates)
- 💻 [CLI Reference](https://tduyng.github.io/gozzi/reference/cli)
- 📖 [Examples](https://tduyng.github.io/gozzi/examples/quick-start)


## 🔧 Advanced Usage

### Multi-Environment Configuration

```toml
# config.prod.toml
base_url = "https://mysite.com"
build_drafts = false

# config.dev.toml
base_url = "http://localhost:3000"
build_drafts = true
```

```bash
gozzi build --config config.prod.toml
gozzi serve --config config.dev.toml
```

### Custom Template Functions

Gozzi includes 40+ template functions:

```html
<!-- Date formatting -->
{{ date .Date "January 2, 2006" }}

<!-- Content filtering -->
{{ range where .Site.Pages "Featured" true | first 5 }}
  <article>{{ .Title }}</article>
{{ end }}

<!-- URL generation -->
<a href="{{ asset "css/main.css" }}">Stylesheet</a>

<!-- Math rendering -->
{{ if .HasMath }}
  <script src="https://cdn.jsdelivr.net/npm/katex@latest/dist/katex.min.js"></script>
{{ end }}
```

### Content Organization

```
content/
├── _index.md           # Home page
├── about/
│   └── _index.md       # About section
├── blog/
│   ├── _index.md       # Blog listing
│   ├── 2024-01-15-post1.md
│   └── 2024-01-20-post2.md
└── projects/
    ├── _index.md
    └── gozzi/
        └── _index.md
```

## 🌟 Real-World Example

See [tduyng.github.io](https://github.com/tduyng/tduyng.github.io) for a complete example featuring:

- 100+ blog posts and notes
- Multiple content sections
- Custom templates and styling
- Math expressions and diagrams
- Tag-based organization
- RSS feeds and SEO optimization

## 🛠️ Development

### Building from Source

```bash
git clone https://github.com/tduyng/gozzi.git
cd gozzi
go build -o gozzi .
```

### Running Tests

```bash
go test ./...
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Happy building with Gozzi!** 🍃

For questions, issues, or feature requests, please [open an issue](https://github.com/tduyng/gozzi/issues) on GitHub.
