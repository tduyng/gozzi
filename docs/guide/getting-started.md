# Getting Started

This guide will help you create your first Gozzi site from scratch in less than 5 minutes.

## Prerequisites

Before you begin, make sure you have:

- **Go 1.25+** installed ([installation guide](https://go.dev/doc/install))
- **Basic command-line knowledge**
- **A text editor** (VS Code, Vim, Neovim, etc.)

::: tip
Check your Go version:

```bash
go version
# Should show go1.25 or higher
```

:::

## Step 1: Install Gozzi

Choose your preferred installation method:

::: code-group

```bash [go install]
go install github.com/tduyng/gozzi@latest
```

```bash [Homebrew (coming soon)]
brew install gozzi
```

```bash [From Release]
# Download from GitHub Releases
curl -LO https://github.com/tduyng/gozzi/releases/latest/download/gozzi_Darwin_arm64.tar.gz
tar -xzf gozzi_Darwin_arm64.tar.gz
sudo mv gozzi /usr/local/bin/
```

:::

Verify installation:

```bash
gozzi version
# Output: gozzi version 0.0.10
```

For detailed installation instructions, see the [Installation Guide](/guide/installation).

## Step 2: Create Project Structure

Create a new directory for your site:

```bash
mkdir my-blog
cd my-blog
```

Create the basic folder structure:

```bash
mkdir -p content/blog
mkdir templates
mkdir static
```

Your project should look like this:

```
my-blog/
├── content/      # Your markdown content
│   └── blog/
├── templates/    # HTML templates
└── static/       # CSS, JS, images
```

## Step 3: Create Configuration

Create a `config.toml` file in the root:

```toml
# config.toml
base_url = "https://example.com"
title = "My Awesome Blog"
description = "A blog built with Gozzi"
language = "en"
generate_feed = true

[extra]
name = "Your Name"
bio = "Developer and writer"
```

::: details What does this config do?

- `base_url`: Your site's URL (used for absolute links)
- `title`: Site title (appears in browser tabs and feeds)
- `description`: SEO meta description
- `generate_feed`: Create RSS/Atom feed automatically
- `[extra]`: Custom data accessible in templates
  :::

## Step 4: Create Your First Post

Create a new blog post:

```bash
mkdir -p content/blog/2024-01-15-hello-world
```

Create `content/blog/2024-01-15-hello-world/index.md`:

```markdown
+++
title = "Hello, World!"
date = 2024-01-15
description = "My first post with Gozzi"
tags = ["introduction", "blog"]
+++

# Welcome to My Blog!

This is my first post built with **Gozzi**, a fast static site generator written in Go.

## Why Gozzi?

- 🚀 Lightning fast builds
- 📊 Rich content support
- 🎨 Flexible templating
- 🌐 SEO ready

Let's build something amazing!
```

::: tip Date-based URLs
The folder name `2024-01-15-hello-world` becomes `/blog/hello-world/` in the URL. Gozzi automatically strips the date prefix.
:::

## Step 5: Create a Basic Template

Create `templates/post.html`:

```html
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>{{ .Page.Config.title }} - {{ .Site.Config.title }}</title>
        <meta name="description" content="{{ .Page.Config.description }}" />
        <style>
            body {
                max-width: 800px;
                margin: 0 auto;
                padding: 2rem;
                font-family:
                    system-ui,
                    -apple-system,
                    sans-serif;
                line-height: 1.6;
            }
            h1 {
                color: #10b981;
            }
            .meta {
                color: #666;
                margin-bottom: 2rem;
            }
            .tag {
                display: inline-block;
                background: #f0f0f0;
                padding: 0.25rem 0.5rem;
                border-radius: 0.25rem;
                margin-right: 0.5rem;
            }
        </style>
    </head>
    <body>
        <article>
            <h1>{{ .Page.Config.title }}</h1>

            <div class="meta">
                <time>{{ date .Page.Config.date "January 2, 2006" }}</time>

                {{ if .Page.Config.tags }}
                <div style="margin-top: 0.5rem;">
                    {{ range .Page.Config.tags }}
                    <span class="tag">#{{ . }}</span>
                    {{ end }}
                </div>
                {{ end }}
            </div>

            <div class="content">{{ .Page.Content }}</div>
        </article>
    </body>
</html>
```

## Step 6: Start Development Server

Start the development server with live reload:

```bash
gozzi serve
```

You should see:

```
Starting server on http://localhost:3000
Watching for changes...
```

Open your browser to [http://localhost:3000/blog/hello-world/](http://localhost:3000/blog/hello-world/) 🎉

::: tip Live Reload
Try editing your post or template - the browser will automatically refresh!
:::

## Step 7: Build for Production

When you're ready to deploy:

```bash
gozzi build
```

This creates a `public/` directory with your complete static site:

```
public/
├── index.html
├── blog/
│   └── hello-world/
│       └── index.html
├── atom.xml
├── sitemap.xml
└── robots.txt
```

You can now deploy the `public/` folder to any static hosting service!

## What's Next?

Now that you have a basic site running, explore more features:

### Learn the Basics

- [Content Structure](/guide/content-structure) - Organize pages and sections
- [Configuration](/guide/configuration) - Customize your site
- [Templates](/guide/templates) - Build custom layouts

### Add Features

- [Tags & Pagination](/guide/features) - Organize content
- [RSS & SEO](/guide/features) - Optimize for search engines
- [Built-in Features](/guide/features) - TOC, math, diagrams

### See Examples

- [Blog Setup](/examples/quick-start) - Complete blog example
- [Custom Templates](/guide/templates) - Advanced templating
- [Real-World Sites](/examples/real-world) - Production examples

## Common Next Steps

### Add More Posts

Create additional posts in the `content/blog/` directory:

```bash
mkdir -p content/blog/2024-01-16-second-post
```

### Create a Blog Listing Page

Create `content/blog/_index.md`:

```markdown
+++
title = "Blog"
description = "All my blog posts"
template = "blog.html"
+++

Welcome to my blog!
```

And create `templates/blog.html` to list all posts.

### Add CSS and Images

Put files in `static/`:

```
static/
├── css/
│   └── main.css
└── img/
    └── avatar.jpg
```

Reference in templates:

```html
<link rel="stylesheet" href="/css/main.css" /> <img src="/img/avatar.jpg" alt="Avatar" />
```

## Troubleshooting

### Port Already in Use

If port 1313 is busy:

```bash
gozzi serve --port 3000
```

### Command Not Found

Make sure Go's bin directory is in your PATH:

```bash
export PATH="$PATH:$HOME/go/bin"
```

### Template Errors

Check that your template file names match the content type:

- `post.html` for blog posts
- `blog.html` for blog section
- `home.html` for homepage

Need help? [Open an issue](https://github.com/tduyng/gozzi/issues) on GitHub!
