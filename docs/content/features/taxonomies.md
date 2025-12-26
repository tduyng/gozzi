+++
title = "Taxonomies"
date = 2025-12-26
template = "page.html"
+++

Taxonomies let you organize content into groups. Common examples: tags, series, categories.

## What are Taxonomies?

A taxonomy is a classification system. You add metadata to your content, and Gozzi automatically:

1. Creates index pages listing all terms (e.g., `/tags/`, `/series/`)
2. Creates pages for each term (e.g., `/tags/golang/`, `/series/neovim/`)
3. Lists all content with that term
4. Provides navigation between related content

## Quick Start

**1. Enable taxonomy in `config.toml`:**

```toml
[taxonomies]
tags = { enabled = true, paginate_by = 0 }
series = { enabled = true, paginate_by = 0 }
categories = { enabled = true, paginate_by = 10 }
```

**2. Add taxonomy to your content:**

```toml
+++
title = "My Post"
tags = ["golang", "web"]
series = "golang-tutorial"
series_order = 1
categories = ["programming"]
+++
```

**3. Access in templates:**

Gozzi generates pages automatically. You just need templates.

## Configuration

Enable taxonomies in `config.toml`:

```toml
[taxonomies]
# Taxonomy name = { enabled = true/false, paginate_by = number }
tags = { enabled = true, paginate_by = 0 }
series = { enabled = true, paginate_by = 0 }
categories = { enabled = true, paginate_by = 20 }
```

**Options:**

- `enabled`: Turn taxonomy on/off
- `paginate_by`: Items per page (0 = no pagination)

## Content Frontmatter

Add taxonomy values to your content:

**Simple values (single term):**

```toml
+++
series = "neovim-guide"
category = "tutorial"
+++
```

**Array values (multiple terms):**

```toml
+++
tags = ["golang", "web", "api"]
categories = ["programming", "backend"]
+++
```

## Series with Ordering

Series support special ordering with `series_order`:

```toml
+++
title = "Part 1: Introduction"
series = "golang-tutorial"
series_order = 1
+++
```

```toml
+++
title = "Part 2: Variables"
series = "golang-tutorial"
series_order = 2
+++
```

Posts are automatically ordered by `series_order` (or by date if not specified).

## Templates

Create templates for taxonomy pages:

### Index Page Template

Shows all terms in a taxonomy.

**File: `templates/tags.html`**

```html
<!DOCTYPE html>
<html>
<head>
    <title>All Tags</title>
</head>
<body>
    <h1>Tags</h1>
    <ul>
    {{ range .Terms }}
        <li>
            <a href="{{ .Permalink }}">{{ .Name }}</a>
            <span>({{ .Count }} posts)</span>
        </li>
    {{ end }}
    </ul>
</body>
</html>
```

**Available data:**

- `.Terms` - Array of all terms in this taxonomy
  - `.Name` - Term name (e.g., "golang")
  - `.Slug` - URL-safe slug (e.g., "golang")
  - `.Permalink` - Full URL (e.g., "/tags/golang/")
  - `.Count` - Number of pages with this term
  - `.Pages` - Array of pages with this term

### Term Page Template

Shows all content for a specific term.

**File: `templates/tag.html`**

```html
<!DOCTYPE html>
<html>
<head>
    <title>Tag: {{ .Term.Name }}</title>
</head>
<body>
    <h1>Posts tagged "{{ .Term.Name }}"</h1>
    <p>{{ .Term.Count }} posts</p>
    
    <ul>
    {{ range .Term.Pages }}
        <li>
            <a href="{{ .Permalink }}">{{ .Config.title }}</a>
            <time>{{ date .Config.date "Jan 2, 2006" }}</time>
        </li>
    {{ end }}
    </ul>
</body>
</html>
```

**Available data:**

- `.Term` - Current term information
  - `.Name` - Term name
  - `.Slug` - URL slug
  - `.Permalink` - Full URL
  - `.Count` - Number of pages
  - `.Pages` - Array of pages with this term

### Series Page Template

Series pages show ordered content.

**File: `templates/serie.html`**

```html
<!DOCTYPE html>
<html>
<head>
    <title>Series: {{ .Term.Name }}</title>
</head>
<body>
    <h1>{{ .Term.Name }} Series</h1>
    <p>{{ .Term.Count }} parts</p>
    
    <ol>
    {{ range .Term.Pages }}
        <li>
            <strong>Part {{ .Position }}</strong>
            <a href="{{ .Permalink }}">{{ .Config.title }}</a>
        </li>
    {{ end }}
    </ol>
</body>
</html>
```

**Series-specific data:**

- `.Position` - Part number in series (from `series_order`)

### Series Navigation in Posts

Show "Part X of Y" and prev/next links in individual posts.

**File: `templates/post.html`**

```html
<!DOCTYPE html>
<html>
<body>
    <article>
        <h1>{{ .Page.Config.title }}</h1>
        
        <!-- Series badge at top -->
        {{ with .Page.Series }}
        <div class="series-badge">
            Part <strong>{{ .CurrentPart }}</strong> of <strong>{{ .TotalPosts }}</strong> in
            <a href="{{ .Permalink }}">{{ .Name }}</a>
        </div>
        {{ end }}
        
        {{ .Page.Content }}
        
        <!-- Series navigation at bottom -->
        {{ with .Page.Series }}
        <nav class="series-nav">
            {{ if .PreviousPost }}
            <a href="{{ .PreviousPost.Permalink }}">
                ← {{ .PreviousPost.Title }}
            </a>
            {{ end }}
            
            {{ if .NextPost }}
            <a href="{{ .NextPost.Permalink }}">
                {{ .NextPost.Title }} →
            </a>
            {{ end }}
        </nav>
        {{ end }}
    </article>
</body>
</html>
```

**Series navigation data:**

- `.Page.Series` - Series information (only if page is in a series)
  - `.Name` - Series name
  - `.Slug` - URL slug
  - `.Permalink` - Series page URL
  - `.TotalPosts` - Total parts in series
  - `.CurrentPart` - This page's position
  - `.PreviousPost` - Previous part (or nil)
    - `.Title` - Post title
    - `.Permalink` - Post URL
  - `.NextPost` - Next part (or nil)

## Template Naming

Gozzi looks for templates using this naming pattern:

| Page Type | Template Files (in order) |
|-----------|---------------------------|
| All tags | `tags.html` → `page.html` |
| Single tag | `tag.html` → `page.html` |
| All series | `series.html` → `page.html` |
| Single series | `serie.html` → `page.html` |
| All categories | `categories.html` → `page.html` |
| Single category | `category.html` → `page.html` |

**Custom taxonomies:**

For a taxonomy named `authors`:
- All authors: `authors.html`
- Single author: `author.html`

(Singular form is the name minus 's', or you can use the same plural name)

## URL Structure

Taxonomy pages follow this pattern:

```
/tags/                  # All tags index
/tags/golang/           # Posts tagged "golang"
/series/                # All series index
/series/neovim-guide/   # Posts in "neovim-guide" series
/categories/            # All categories
/categories/tutorial/   # Posts in "tutorial" category
```

## Examples

### Example 1: Blog with Tags

**config.toml:**

```toml
[taxonomies]
tags = { enabled = true, paginate_by = 0 }
```

**content/blog/my-post.md:**

```toml
+++
title = "Building APIs in Go"
tags = ["golang", "api", "web"]
+++

Content here...
```

**Result:**

- `/tags/` - Lists all tags
- `/tags/golang/` - Shows all posts tagged "golang"
- `/tags/api/` - Shows all posts tagged "api"
- `/tags/web/` - Shows all posts tagged "web"

### Example 2: Tutorial Series

**config.toml:**

```toml
[taxonomies]
series = { enabled = true, paginate_by = 0 }
```

**content/tutorial/part-1.md:**

```toml
+++
title = "Introduction to Go"
series = "golang-tutorial"
series_order = 1
+++
```

**content/tutorial/part-2.md:**

```toml
+++
title = "Variables and Types"
series = "golang-tutorial"
series_order = 2
+++
```

**content/tutorial/part-3.md:**

```toml
+++
title = "Functions"
series = "golang-tutorial"
series_order = 3
+++
```

**Result:**

- `/series/` - Lists all series
- `/series/golang-tutorial/` - Shows 3 parts in order
- Each post shows "Part X of 3" badge
- Each post has prev/next navigation

### Example 3: Multi-Taxonomy

**config.toml:**

```toml
[taxonomies]
tags = { enabled = true, paginate_by = 0 }
series = { enabled = true, paginate_by = 0 }
categories = { enabled = true, paginate_by = 10 }
```

**content/posts/neovim-lsp.md:**

```toml
+++
title = "Neovim LSP Setup"
tags = ["neovim", "lsp", "editor"]
series = "neovim-guide"
series_order = 3
categories = ["tutorials"]
+++
```

**Result:**

- Tagged with 3 tags → appears on 3 tag pages
- Part of series → shows "Part 3 of X" badge
- In category → appears on category page
- Full navigation between all related content

## Incremental Builds

Gozzi is smart about rebuilding:

**When you modify a post:**
- Only affected taxonomy pages regenerate
- Other taxonomy pages stay cached
- Fast incremental rebuilds

**Example:**

```bash
# Modify a post tagged "golang"
vim content/posts/my-post.md

# Only these regenerate:
# - /tags/golang/ (tag changed)
# - /tags/ (index updated)
# - Your post
# Everything else uses cache
```

## Performance Tips

**1. Use pagination for large taxonomies:**

```toml
[taxonomies]
tags = { enabled = true, paginate_by = 20 }  # 20 posts per page
```

**2. Limit taxonomy usage:**

Only enable taxonomies you need:

```toml
[taxonomies]
tags = { enabled = true, paginate_by = 0 }
# Don't enable series if you don't use them
# series = { enabled = false }
```

**3. Keep term counts reasonable:**

Avoid hundreds of tags. Use 5-10 tags per post max.

## Troubleshooting

### Series navigation not showing

**Problem:** `.Page.Series` is nil in templates

**Solution:**

1. Check taxonomy is enabled:
   ```toml
   [taxonomies]
   series = { enabled = true, paginate_by = 0 }
   ```

2. Check frontmatter has `series`:
   ```toml
   series = "my-series"
   series_order = 1
   ```

3. Rebuild site (not just incremental):
   ```bash
   gozzi build
   ```

### Template not found

**Problem:** `template "tags.html" not found`

**Solution:**

Create the template file:
```bash
touch templates/tags.html
touch templates/tag.html
```

### Wrong order in series

**Problem:** Series posts in wrong order

**Solution:**

Add `series_order` to ALL posts in series:

```toml
series_order = 1  # Not optional for proper ordering
```

## Migration from Other Generators

### From Hugo

Hugo's taxonomy config:

```toml
[taxonomies]
  tag = "tags"
  category = "categories"
```

Gozzi equivalent:

```toml
[taxonomies]
tags = { enabled = true, paginate_by = 0 }
categories = { enabled = true, paginate_by = 0 }
```

Frontmatter stays the same!

### From Zola

Zola config:

```toml
[taxonomies]
[[taxonomies]]
name = "tags"
```

Gozzi equivalent:

```toml
[taxonomies]
tags = { enabled = true, paginate_by = 0 }
```

## Summary

**Key Points:**

1. ✅ Enable taxonomies in `config.toml`
2. ✅ Add taxonomy values to content frontmatter
3. ✅ Create templates for taxonomy pages
4. ✅ Use `.Page.Series` for series navigation
5. ✅ Gozzi handles URL generation automatically

**Common Use Cases:**

- **Tags**: Organize posts by topics
- **Series**: Multi-part tutorials with navigation
- **Categories**: High-level organization
- **Authors**: Multi-author blogs
- **Custom**: Any classification system you need

---

**Related:**
- [Frontmatter Configuration](/config/frontmatter)
- [Template Variables](/templates/variables)
- [Content Structure](/content-structure)
