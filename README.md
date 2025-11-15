# Gozzi 🍃

**A Go static site generator born from curiosity**

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/tduyng/gozzi)](https://github.com/tduyng/gozzi/releases)

I built **Gozzi** to learn how static site generators work. Now I use it for [tduyng.com](https://tduyng.com).

Gozzi is simple, fast, and helps you learn Go along the way.

## ✨ Features

- 🚀 **Fast builds** - Sub-second for most sites with live reload
- 📝 **Markdown first** - Write content, organize in folders
- 🎨 **Go templates** - Flexible HTML templates with 40+ helper functions
- 🌐 **SEO ready** - RSS, sitemap, Open Graph support
- 📊 **Rich content** - Math (KaTeX), diagrams (Mermaid), syntax highlighting

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

Visit **http://localhost:3000** 🎉

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
- Math expressions and diagrams
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
