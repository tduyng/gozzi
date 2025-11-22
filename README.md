# Gozzi 🍃

**A Go static site generator born from curiosity**

<p align="left">
  <a href="https://golang.org/">
    <img alt="Go Version" title="Go 1.25+" src="https://custom-icon-badges.demolab.com/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white&labelColor=302D41"/></a>
  <a href="https://github.com/tduyng/gozzi/actions/workflows/ci.yml">
    <img alt="Build Status" title="CI Status" src="https://img.shields.io/github/actions/workflow/status/tduyng/gozzi/ci.yml?branch=main&style=for-the-badge&logo=github-actions&logoColor=white&labelColor=302D41&color=89B4FA"/></a>
  <a href="LICENSE">
    <img alt="License" title="MIT License" src="https://custom-icon-badges.demolab.com/badge/License-MIT-A6E3A1?style=for-the-badge&logo=law&logoColor=white&labelColor=302D41"/></a>
  <a href="https://github.com/tduyng/gozzi/releases">
    <img alt="Release" title="Latest Release" src="https://custom-icon-badges.demolab.com/github/v/release/tduyng/gozzi?style=for-the-badge&logo=rocket&color=DDB6F2&logoColor=white&labelColor=302D41"/></a>
  <a href="https://github.com/tduyng/gozzi/stargazers">
    <img alt="Stars" title="Star on GitHub" src="https://custom-icon-badges.demolab.com/github/stars/tduyng/gozzi?style=for-the-badge&logo=star&color=F5E0DC&logoColor=white&labelColor=302D41"/></a>
</p>

I built **Gozzi** to learn how static site generators work. Now I use it for [tduyng.com](https://tduyng.com).

Gozzi is simple, fast, and helps you learn Go along the way.

## ✨ Features

- 🚀 **Fast builds** - Sub-second for most sites with live reload
- 📝 **Markdown first** - Write content, organize in folders
- 🎨 **Go templates** - Flexible HTML templates with 40+ helper functions
- 🌐 **SEO ready** - RSS, sitemap, Open Graph support
- 📊 **Rich content** - Native KaTeX math, Mermaid diagrams, syntax highlighting

## 🚀 Quick Start

```bash
# Install
go install github.com/tduyng/gozzi@latest

# Create a site
mkdir my-blog && cd my-blog
mkdir -p content/blog templates static

# Configure
cat > config.toml << 'EOF'
base_url = "https://example.com"
title = "My Blog"
EOF

# Create first post
cat > content/blog/hello.md << 'EOF'
+++
title = "Hello World"
date = 2024-01-15
+++

My first post!
EOF

# Serve with live reload
gozzi serve
```

Visit **http://localhost:1313** 🎉

## 📚 Documentation

📖 **[Read the full documentation →](https://tduyng.github.io/gozzi/)**

- 🚀 [Getting Started](https://tduyng.github.io/gozzi/guide/getting-started)
- 📖 [Installation](https://tduyng.github.io/gozzi/guide/installation)
- ⚙️ [Configuration](https://tduyng.github.io/gozzi/guide/configuration)
- 🎨 [Templates](https://tduyng.github.io/gozzi/guide/templates)
- 🏗️ [Content Structure](https://tduyng.github.io/gozzi/guide/content-structure)
- ✨ [Features](https://tduyng.github.io/gozzi/guide/features)
- 💻 [CLI Reference](https://tduyng.github.io/gozzi/reference/cli)
- 🔧 [Template Functions](https://tduyng.github.io/gozzi/reference/template-functions)
- 🌟 [Examples](https://tduyng.github.io/gozzi/examples/quick-start)

## 🌟 Real-World Example

My personal site [tduyng.com](https://tduyng.com) runs on Gozzi:

- 100+ blog posts and notes
- Native KaTeX math and Mermaid diagrams
- Full-text search
- ~100ms builds

[View the source →](https://github.com/tduyng/tduyng.github.io)

## 🛠️ Development

```bash
# Clone and build
git clone https://github.com/tduyng/gozzi.git
cd gozzi
go build -o gozzi .

# Run tests
go test ./...
```

**Contributing:** PRs welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

**Built for learning. Shared for learning.** 🍃

[Documentation](https://tduyng.github.io/gozzi/) · [Issues](https://github.com/tduyng/gozzi/issues) · [Discussions](https://github.com/tduyng/gozzi/discussions)
