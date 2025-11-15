# Real-World Example: tduyng.com

This is how I actually use Gozzi for my personal website [tduyng.com](https://tduyng.com).

**Source code**: [github.com/tduyng/tduyng.github.io](https://github.com/tduyng/tduyng.github.io)

## What It Is

A personal blog and notes site with:
- 100+ blog posts (technical articles)
- 20+ quick notes (TILs, short learnings)
- Math expressions (KaTeX)
- Diagrams (Mermaid)
- Full-text search
- Comments (Giscus)

## Project Structure

```
tduyng.github.io/
├── config/
│   ├── config.toml         # Main config
│   ├── config.dev.toml     # Local development
│   └── config.prod.toml    # Production
│
├── content/
│   ├── _index.md           # Homepage
│   ├── about/_index.md     # About page
│   ├── blog/               # Blog posts (100+)
│   │   └── 2024-01-15-post-title/
│   │       └── index.md
│   ├── notes/              # Quick notes
│   │   └── 2024-01-20-note/
│   │       └── index.md
│   └── contact/_index.md   # Contact page
│
├── templates/
│   ├── home.html           # Homepage
│   ├── blog.html           # Blog listing
│   ├── post.html           # Single post
│   ├── notes.html          # Notes listing
│   ├── note.html           # Single note
│   ├── prose.html          # Generic pages
│   └── partials/           # Reusable pieces
│       ├── _head.html
│       ├── _header.html
│       ├── _footer.html
│       └── _toc.html
│
├── static/
│   ├── css/main.css
│   ├── js/main.js
│   └── img/
│
└── .github/workflows/
    └── deploy.yml          # Auto-deploy
```

## Key Configuration

```toml
# config.prod.toml
base_url = "https://tduyng.com"
title = "Tien Duy"
description = "Software Engineer & Technical Writer"
generate_feed = true

[extra]
name = "Tien Duy Nguyen"
bio = "Building things and writing about it"

# Navigation sections
sections = [
  { name = "blog", path = "/blog" },
  { name = "notes", path = "/notes" },
  { name = "about", path = "/about" },
]

# Social links
links = [
  { name = "GitHub", icon = "github", url = "https://github.com/tduyng" },
  { name = "Email", icon = "email", url = "mailto:hi@tduyng.com" },
]
```

## Content Example

### Blog Post

```markdown
+++
title = "Understanding Go Channels"
date = 2024-01-15
description = "Deep dive into Go channels and concurrency"
tags = ["go", "concurrency", "programming"]
featured = true

[extra]
toc = true
+++

# Understanding Go Channels

Channels are Go's way of enabling communication between goroutines...

## Basic Usage

```go
ch := make(chan int)
go func() {
    ch <- 42
}()
value := <-ch
```

## Key Concepts...
```

### Quick Note

```markdown
+++
title = "TIL: Git worktree"
date = 2024-01-20
tags = ["git", "til"]
+++

Learned about `git worktree` today. Super useful for working on multiple branches:

```bash
git worktree add ../feature-branch feature-branch
# Work in ../feature-branch directory
# No need to stash or commit
```
```

## Template Example

Simple post template with TOC:

```html
<!DOCTYPE html>
<html>
{{ template "partials/_head.html" . }}
<body>
  {{ template "partials/_header.html" . }}
  
  <main class="container">
    <div class="content-wrapper">
      <!-- Sidebar with TOC -->
      {{ if .Page.Config.extra.toc }}
      <aside class="toc">
        {{ .Page.TOC }}
      </aside>
      {{ end }}
      
      <!-- Main content -->
      <article>
        <h1>{{ .Page.Config.title }}</h1>
        
        <div class="meta">
          <time>{{ date .Page.Config.date "January 2, 2006" }}</time>
          
          {{ if .Page.Config.tags }}
          <div class="tags">
            {{ range .Page.Config.tags }}
            <a href="/tags/{{ . | urlize }}">#{{ . }}</a>
            {{ end }}
          </div>
          {{ end }}
        </div>
        
        <div class="prose">
          {{ .Page.Content }}
        </div>
      </article>
    </div>
  </main>
  
  {{ template "partials/_footer.html" . }}
</body>
</html>
```

## Build Performance

- **100+ posts**: ~500ms build time
- **Development**: <200ms incremental rebuilds
- **Live reload**: Instant browser refresh

## Deployment

GitHub Actions workflow:

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      
      - name: Install Gozzi
        run: go install github.com/tduyng/gozzi@latest
      
      - name: Build
        run: gozzi build --config config.prod.toml
      
      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

## What I Learned

### Good Decisions

1. **Simple folder structure** - Easy to find content
2. **Config per environment** - Different settings for dev/prod
3. **Partial templates** - Reuse header/footer everywhere
4. **KaTeX for math** - Works great for technical posts
5. **Tag-based organization** - Natural content discovery

### What Works Well

- **Fast builds** - 500ms even with 100+ posts
- **Live reload** - Saves so much time during writing
- **Clear config** - Easy to understand and modify
- **Go templates** - Simple enough, powerful enough

### If I Started Today

I'd do mostly the same, but:
- Maybe organize posts by year folders
- Add more examples in docs earlier
- Set up search from day one

## Key Takeaways

Gozzi works well for personal sites where you want:
- **Control** - You understand everything
- **Simplicity** - Clear configuration, no surprises
- **Speed** - Fast builds, good developer experience
- **Learning** - Great way to learn Go templates

The source code is all public, so you can:
- Copy the template structure
- See how features work together
- Use it as a starting point for your site

**Browse the code**: [github.com/tduyng/tduyng.github.io](https://github.com/tduyng/tduyng.github.io)
