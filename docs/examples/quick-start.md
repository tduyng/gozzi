# Quick Start Example

Build a complete blog with Gozzi in under 10 minutes.

## What We'll Build

A fully-functional blog with:

- ✅ Homepage with featured posts
- ✅ Blog listing page
- ✅ Individual blog posts
- ✅ Tag system
- ✅ RSS feed
- ✅ SEO optimization

## Project Setup

### 1. Create Project

```bash
mkdir gozzi-blog && cd gozzi-blog
mkdir -p content/blog static/css templates/{partials,macros}
```

### 2. Configuration

Create `config.toml`:

```toml
base_url = "https://myblog.com"
title = "My Awesome Blog"
description = "Thoughts on tech, life, and everything in between"
language = "en"
output_dir = "public"
generate_feed = true

[extra]
name = "Jane Developer"
bio = "Software engineer and writer"
avatar = "/img/avatar.jpg"

# Social links
links = [
  { name = "GitHub", icon = "github", url = "https://github.com/yourusername" },
  { name = "Twitter", icon = "twitter", url = "https://twitter.com/yourusername" },
  { name = "Email", icon = "email", url = "mailto:hello@myblog.com" },
]

# Navigation
sections = [
  { name = "blog", path = "/blog", is_external = false },
  { name = "about", path = "/about", is_external = false },
]
```

## Content Structure

### 3. Create Blog Posts

Create `content/blog/2024-01-15-first-post/index.md`:

```markdown
+++
title = "Getting Started with Gozzi"
date = 2024-01-15
description = "Learn how to build a blog with Gozzi in minutes"
tags = ["tutorial", "go", "blogging"]
featured = true
+++

# Getting Started with Gozzi

Welcome to my blog! Today I'll show you how easy it is to build a static site with Gozzi.

## Why Gozzi?

- **Fast**: Sub-second builds
- **Simple**: Clean configuration
- **Powerful**: Rich content support

## Quick Example

Here's a simple code example:

\`\`\`go
package main

import "fmt"

func main() {
fmt.Println("Hello from Gozzi!")
}
\`\`\`

Pretty cool, right? Stay tuned for more posts!
```

Create `content/blog/2024-01-20-second-post/index.md`:

```markdown
+++
title = "Building Modern Web Apps"
date = 2024-01-20
description = "Modern web development best practices"
tags = ["web", "javascript", "tutorial"]
featured = false
+++

# Building Modern Web Apps

Let's explore modern web development...

## Topics Covered

1. **Performance** - Make it fast
2. **Accessibility** - Make it usable
3. **SEO** - Make it discoverable
```

### 4. Create Blog Section

Create `content/blog/_index.md`:

```markdown
+++
title = "Blog"
description = "Latest articles and tutorials"
template = "blog.html"
+++

Welcome to my blog where I write about technology, programming, and life.
```

### 5. Create Homepage

Create `content/_index.md`:

```markdown
+++
title = "Home"
description = "Welcome to my blog"
template = "home.html"
+++

# Welcome!

I'm a developer who loves to write about technology.
```

## Templates

### 6. Create Base Layout

Create `templates/partials/_head.html`:

```html
<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{ priority .Page.Config.title .Section.Config.title .Site.Config.title }}</title>

    {{ $description := priority .Page.Config.description .Section.Config.description
    .Site.Config.description }}
    <meta name="description" content="{{ $description }}" />

    <!-- Styles -->
    <link rel="stylesheet" href="/css/main.css" />

    <!-- RSS Feed -->
    <link
        rel="alternate"
        type="application/atom+xml"
        title="{{ .Site.Config.title }}"
        href="/atom.xml"
    />
</head>
```

Create `templates/partials/_header.html`:

```html
<header class="site-header">
    <nav class="container">
        <a href="/" class="site-title">{{ .Site.Config.title }}</a>

        <ul class="nav-links">
            {{ range .Site.Config.extra.sections }}
            <li><a href="{{ .path }}">{{ .name }}</a></li>
            {{ end }}
        </ul>
    </nav>
</header>
```

Create `templates/partials/_footer.html`:

```html
<footer class="site-footer">
    <div class="container">
        <p>&copy; 2024 {{ .Site.Config.extra.name }}. Built with Gozzi.</p>

        <div class="social-links">
            {{ range .Site.Config.extra.links }}
            <a href="{{ .url }}" target="_blank" rel="noopener">{{ .name }}</a>
            {{ end }}
        </div>
    </div>
</footer>
```

### 7. Create Home Template

Create `templates/home.html`:

```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.language }}">
    {{ template "partials/_head.html" . }}
    <body class="home">
        {{ template "partials/_header.html" . }}

        <main class="container">
            <section class="hero">
                <h1>{{ .Site.Config.title }}</h1>
                <p class="lead">{{ .Site.Config.description }}</p>
            </section>

            <section class="featured-posts">
                <h2>Featured Posts</h2>

                {{ $blog := get_section "blog" }} {{ $featured := where $blog.Children
                "Config.featured" true }}

                <div class="post-grid">
                    {{ range $featured }}
                    <article class="post-card">
                        <h3><a href="{{ .Permalink }}">{{ .Config.title }}</a></h3>
                        <time>{{ date .Config.date "January 2, 2006" }}</time>
                        <p>{{ .Config.description }}</p>

                        <div class="tags">
                            {{ range .Config.tags }}
                            <span class="tag">#{{ . }}</span>
                            {{ end }}
                        </div>
                    </article>
                    {{ end }}
                </div>
            </section>

            <section class="recent-posts">
                <h2>Recent Posts</h2>

                {{ $recent := first 5 (reverse $blog.Children) }}
                <ul class="post-list">
                    {{ range $recent }}
                    <li>
                        <a href="{{ .Permalink }}">{{ .Config.title }}</a>
                        <time>{{ date .Config.date "Jan 2, 2006" }}</time>
                    </li>
                    {{ end }}
                </ul>
            </section>
        </main>

        {{ template "partials/_footer.html" . }}
    </body>
</html>
```

### 8. Create Blog Listing Template

Create `templates/blog.html`:

```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.language }}">
    {{ template "partials/_head.html" . }}
    <body class="blog-list">
        {{ template "partials/_header.html" . }}

        <main class="container">
            <header class="section-header">
                <h1>{{ .Section.Config.title }}</h1>
                <p>{{ .Section.Config.description }}</p>
            </header>

            <div class="posts">
                {{ range .Section.Children }}
                <article class="post-card">
                    <h2><a href="{{ .Permalink }}">{{ .Config.title }}</a></h2>

                    <div class="meta">
                        <time>{{ date .Config.date "January 2, 2006" }}</time>
                        {{ if .Config.featured }}
                        <span class="badge featured">Featured</span>
                        {{ end }}
                    </div>

                    <p>{{ .Config.description }}</p>

                    {{ if .Config.tags }}
                    <div class="tags">
                        {{ range .Config.tags }}
                        <a href="/tags/{{ . | urlize }}" class="tag">#{{ . }}</a>
                        {{ end }}
                    </div>
                    {{ end }}
                </article>
                {{ end }}
            </div>
        </main>

        {{ template "partials/_footer.html" . }}
    </body>
</html>
```

### 9. Create Post Template

Create `templates/post.html`:

```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.language }}">
{{ template "partials/_head.html" . }}
<body class="post">
    {{ template "partials/_header.html" . }}

    <main class="container">
        <article class="prose">
            <header class="post-header">
                <h1>{{ .Page.Config.title }}</h1>

                <div class="meta">
                    <time datetime="{{ .Page.Config.date.Format "2006-01-02" }}">
                        {{ date .Page.Config.date "January 2, 2006" }}
                    </time>
                </div>

                {{ if .Page.Config.tags }}
                <div class="tags">
                    {{ range .Page.Config.tags }}
                    <a href="/tags/{{ . | urlize }}" class="tag">#{{ . }}</a>
                    {{ end }}
                </div>
                {{ end }}
            </header>

            <div class="content">
                {{ .Page.Content }}
            </div>

            <footer class="post-footer">
                <div class="post-nav">
                    {{ with .Page.Higher }}
                    <a href="{{ .Permalink }}" class="prev">
                        ← {{ .Config.title }}
                    </a>
                    {{ end }}

                    {{ with .Page.Lower }}
                    <a href="{{ .Permalink }}" class="next">
                        {{ .Config.title }} →
                    </a>
                    {{ end }}
                </div>
            </footer>
        </article>
    </main>

    {{ template "partials/_footer.html" . }}
</body>
</html>
```

## Styling

### 10. Add CSS

Create `static/css/main.css`:

```css
/* Reset & Base */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family:
        -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
    line-height: 1.6;
    color: #333;
    background: #fff;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
}

/* Header */
.site-header {
    border-bottom: 1px solid #e5e7eb;
    padding: 1rem 0;
    position: sticky;
    top: 0;
    background: #fff;
    z-index: 100;
}

.site-header nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.site-title {
    font-size: 1.5rem;
    font-weight: bold;
    color: #10b981;
    text-decoration: none;
}

.nav-links {
    display: flex;
    gap: 2rem;
    list-style: none;
}

.nav-links a {
    color: #666;
    text-decoration: none;
}

.nav-links a:hover {
    color: #10b981;
}

/* Hero */
.hero {
    text-align: center;
    padding: 4rem 0;
}

.hero h1 {
    font-size: 3rem;
    margin-bottom: 1rem;
    color: #111;
}

.lead {
    font-size: 1.25rem;
    color: #666;
}

/* Post Grid */
.post-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 2rem;
    margin-top: 2rem;
}

.post-card {
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    padding: 1.5rem;
    transition: box-shadow 0.2s;
}

.post-card:hover {
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.post-card h2,
.post-card h3 {
    margin-bottom: 0.5rem;
}

.post-card a {
    color: #111;
    text-decoration: none;
}

.post-card a:hover {
    color: #10b981;
}

.post-card time {
    color: #666;
    font-size: 0.9rem;
}

.post-card p {
    margin-top: 0.5rem;
    color: #666;
}

/* Tags */
.tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-top: 1rem;
}

.tag {
    background: #f3f4f6;
    padding: 0.25rem 0.75rem;
    border-radius: 0.25rem;
    font-size: 0.875rem;
    color: #666;
    text-decoration: none;
}

.tag:hover {
    background: #e5e7eb;
}

/* Post Content */
.prose {
    max-width: 65ch;
    margin: 2rem auto;
}

.post-header {
    margin-bottom: 2rem;
}

.post-header h1 {
    font-size: 2.5rem;
    margin-bottom: 1rem;
}

.meta {
    color: #666;
    font-size: 0.9rem;
}

.badge {
    display: inline-block;
    padding: 0.25rem 0.5rem;
    border-radius: 0.25rem;
    font-size: 0.8rem;
    margin-left: 0.5rem;
}

.badge.featured {
    background: #10b981;
    color: white;
}

.content {
    margin-top: 2rem;
}

.content h2 {
    margin-top: 2rem;
    margin-bottom: 1rem;
}

.content p {
    margin-bottom: 1rem;
}

.content code {
    background: #f3f4f6;
    padding: 0.2rem 0.4rem;
    border-radius: 0.25rem;
    font-size: 0.9em;
}

.content pre {
    background: #1f2937;
    color: #f9fafb;
    padding: 1rem;
    border-radius: 0.5rem;
    overflow-x: auto;
    margin: 1rem 0;
}

.content pre code {
    background: none;
    padding: 0;
}

/* Post Navigation */
.post-nav {
    display: flex;
    justify-content: space-between;
    margin-top: 3rem;
    padding-top: 2rem;
    border-top: 1px solid #e5e7eb;
}

.post-nav a {
    color: #10b981;
    text-decoration: none;
}

/* Footer */
.site-footer {
    border-top: 1px solid #e5e7eb;
    padding: 2rem 0;
    margin-top: 4rem;
    text-align: center;
    color: #666;
}

.social-links {
    display: flex;
    gap: 1rem;
    justify-content: center;
    margin-top: 1rem;
}

.social-links a {
    color: #666;
    text-decoration: none;
}

.social-links a:hover {
    color: #10b981;
}
```

## Build & Deploy

### 11. Development

Start the development server:

```bash
gozzi serve
```

Visit [http://localhost:1313](http://localhost:1313)

### 12. Production Build

Build for production:

```bash
gozzi build
```

Your site is ready in the `public/` directory!

### 13. Deploy

Deploy to any static hosting:

::: code-group

```bash [Netlify]
# netlify.toml
[build]
  command = "gozzi build"
  publish = "public"
```

```bash [Vercel]
# vercel.json
{
  "buildCommand": "gozzi build",
  "outputDirectory": "public"
}
```

```bash [GitHub Pages]
# .github/workflows/deploy.yml
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
      - run: go install github.com/tduyng/gozzi@latest
      - run: gozzi build
      - uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

:::

## What You Built

Congratulations! You now have a complete blog with:

✅ Homepage with featured and recent posts  
✅ Blog listing page  
✅ Individual post pages with navigation  
✅ Tag system for categorization  
✅ RSS feed for subscribers  
✅ Responsive design  
✅ SEO optimization

## Next Steps

Enhance your blog:

- [Add pagination](/guide/features) for large post collections
- [Customize templates](/guide/templates) with advanced features
- [Use math expressions](/guide/features#mathematical-expressions-katex) with native KaTeX support
- [Create diagrams](/guide/features#mermaid-diagrams) with native Mermaid support
- [Enable comments](/guide/templates#adding-comments) with Giscus
- [Set up analytics](/guide/templates#analytics) with Plausible

## Complete Code

The complete example is available at:
[github.com/tduyng/gozzi-examples](https://github.com/tduyng/gozzi/tree/main/examples)

Happy blogging! 🚀
