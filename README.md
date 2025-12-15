# Gozzi

A Go static site generator born from curiosity.

I built Gozzi to [learn how static site generators work](https://tduyng.com/blog/why-i-built-my-own-static-site-generator/). Now I use it for [tduyng.com](https://tduyng.com).

## Features

- Fast builds with live reload
- Markdown-first content
- Go templates with 40+ helper functions
- RSS, sitemap, Open Graph support
- Server-side KaTeX math & syntax highlighting
- Mermaid diagrams
- Built-in CSS and HTML minification (36-53% size reduction)

## Quick Start

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
minify_css = true
minify_html = true
minify_js = true
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

Visit http://localhost:1313

## Documentation

Read the full documentation at [tduyng.com/gozzi](https://tduyng.com/gozzi/)

- [Installation](https://tduyng.com/gozzi/installation)
- [Getting Started](https://tduyng.com/gozzi/getting-started)
- [Content Structure](https://tduyng.com/gozzi/content-structure)
- [Configuration](https://tduyng.com/gozzi/config/site)
- [Templates](https://tduyng.com/gozzi/templates/structure)
- [CLI Reference](https://tduyng.com/gozzi/cli/build)
- [Template Functions](https://tduyng.com/gozzi/functions/math-logic)
- [Examples](https://tduyng.com/gozzi/examples/quick-start)

## Example

My personal site [tduyng.com](https://tduyng.com) runs on Gozzi:

- 100+ blog posts and notes
- Server-side KaTeX math rendering (no JS needed)
- Syntax highlighting with Chroma
- Full-text search
- ~100ms builds

[View the source →](https://github.com/tduyng/tduyng.github.io)

## Development

```bash
# Clone and build
git clone https://github.com/tduyng/gozzi.git
cd gozzi
go build -o gozzi .

# Run tests
go test ./...
```

## License

MIT License - see [LICENSE](LICENSE) for details.
