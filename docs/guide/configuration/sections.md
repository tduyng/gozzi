# Section Configuration

Sections are directories with `_index.md` files that create list pages for collections of content.

## What are Sections?

Sections represent content collections in Gozzi:

```
content/
├── blog/
│   ├── _index.md      # Section configuration
│   ├── post-1.md      # Child page
│   └── post-2.md      # Child page
└── notes/
    ├── _index.md      # Section configuration
    └── note-1.md      # Child page
```

## Basic Section Example

```toml
+++
title = "Blog"
description = "Latest articles and thoughts"
template = "blog.html"
date = 2025-01-01

[extra]
hero_text = "Welcome to my blog!"
show_reading_time = true
paginate_by = 10
sort_by = "date"
+++

# Blog Section

This is the blog landing page content.
```

## Section-Specific Fields

Additional fields available for sections (`_index.md` files):

| Field | Type | Description |
|-------|------|-------------|
| `paginate_by` | integer | Number of items per page |
| `sort_by` | string | Sort children by "date", "title", or "updated" |
| `slug` | string | Override URL segment (instead of folder name) |

### Pagination

Control how many items appear per page:

```toml
+++
title = "Blog"
paginate_by = 10  # Show 10 posts per page
+++
```

Generates:
- `/blog/` - First 10 posts
- `/blog/page/2/` - Next 10 posts
- `/blog/page/3/` - And so on...

### Sorting

Control child page order:

```toml
+++
title = "Blog"
sort_by = "date"  # Options: "date", "title", "updated"
+++
```

**Sort options:**
- `"date"` - Newest first (default)
- `"title"` - Alphabetical by title
- `"updated"` - Most recently updated first

### Custom Slugs

Override the URL segment:

```toml
+++
title = "Web Development Blog"
slug = "blog"  # URL: /blog/ instead of /web-development-blog/
+++
```

## Template Access

### Section Data

In templates, section data is available as `.Section`:

```html
<h1>{{ .Section.Title }}</h1>
<p>{{ .Section.Description }}</p>
<p>{{ .Section.Extra.hero_text }}</p>

<!-- Section's own content (from _index.md body) -->
<div class="section-intro">
  {{ .Section.Content | safe }}
</div>
```

### Children Pages

Access child pages with `.Section.Children`:

```html
<!-- List all children -->
{{ range .Section.Children }}
  <article>
    <h2><a href="{{ .Permalink }}">{{ .Title }}</a></h2>
    <time>{{ .Date.Format "2006-01-02" }}</time>
    <p>{{ .Description }}</p>
  </article>
{{ end }}
```

### Pagination

With `paginate_by` set, use `.Paginate`:

```html
<!-- Current page items -->
{{ range .Paginate.Items }}
  <article>
    <h2>{{ .Title }}</h2>
  </article>
{{ end }}

<!-- Pagination nav -->
<nav class="pagination">
  {{ if .Paginate.HasPrev }}
    <a href="{{ .Paginate.PrevURL }}">← Previous</a>
  {{ end }}
  
  <span>Page {{ .Paginate.CurrentPage }} of {{ .Paginate.TotalPages }}</span>
  
  {{ if .Paginate.HasNext }}
    <a href="{{ .Paginate.NextURL }}">Next →</a>
  {{ end }}
</nav>
```

## Common Section Patterns

### Blog Section

```toml
+++
title = "Blog"
description = "Technical articles and tutorials"
template = "blog.html"
paginate_by = 10
sort_by = "date"

[extra]
show_tags = true
show_reading_time = true
+++

Welcome to my blog where I write about web development and Go programming.
```

### Notes Section

```toml
+++
title = "Notes"
description = "Quick thoughts and TILs"
template = "notes.html"
paginate_by = 20
sort_by = "updated"

[extra]
show_dates = true
compact_view = true
+++

A collection of quick notes and things I've learned.
```

### Projects Section

```toml
+++
title = "Projects"
description = "My open source projects and work"
template = "projects.html"
sort_by = "title"

[extra]
show_github_stats = true
+++

Check out my open source projects and contributions.
```

### Documentation Section

```toml
+++
title = "Documentation"
description = "Complete project documentation"
template = "docs.html"
sort_by = "title"

[extra]
sidebar = true
search_enabled = true
+++
```

## Extended Configuration

Use `[extra]` for section-specific data:

```toml
+++
title = "Blog"

[extra]
# Display options
show_reading_time = true
show_author = true
show_tags = true
compact_mode = false

# Section hero
hero_title = "Latest Articles"
hero_subtitle = "Thoughts on web development"
hero_image = "/img/blog-hero.jpg"

# Filtering
featured_tag = "tutorial"
exclude_tags = ["draft", "private"]

# Metadata
section_icon = "📝"
section_color = "#4a90e2"
+++
```

## Section Hierarchy

Sections can be nested:

```
content/
└── docs/
    ├── _index.md              # Docs section
    ├── getting-started/
    │   ├── _index.md          # Nested section
    │   ├── installation.md
    │   └── quickstart.md
    └── advanced/
        ├── _index.md          # Nested section
        └── configuration.md
```

Each nested section gets its own configuration and URL:
- `/docs/` - Main docs section
- `/docs/getting-started/` - Nested section
- `/docs/advanced/` - Nested section

## Template Selection

Sections use template mapping:

1. **Explicit**: `template = "custom.html"` in front matter
2. **Section-specific**: `blog.html` for `/blog/` section
3. **Default**: `section.html` as fallback

## Best Practices

1. **Always include _index.md** - Defines the section
2. **Set pagination** - For large content collections
3. **Use descriptive titles** - Clear section purpose
4. **Write section intro** - Content in `_index.md` body
5. **Configure sorting** - Match user expectations
6. **Use [extra] data** - For template customization
7. **Keep sections focused** - One topic per section

## Examples in Templates

### Featured Posts

```html
<section class="featured">
  <h2>Featured Posts</h2>
  {{ range .Section.Children }}
    {{ if .Config.featured }}
      <article>
        <h3>{{ .Title }}</h3>
        <p>{{ .Description }}</p>
      </article>
    {{ end }}
  {{ end }}
</section>
```

### Tag Filtering

```html
{{ if .Section.Extra.show_tags }}
  <div class="tags-filter">
    {{ range .Site.Tags }}
      <a href="/tags/{{ . | urlize }}">{{ . }}</a>
    {{ end }}
  </div>
{{ end }}
```

### Custom Sorting

```html
<!-- Recent updates -->
{{ range first 5 (sort .Section.Children ".UpdatedDate" "desc") }}
  <article class="recent-update">
    <h3>{{ .Title }}</h3>
    <time>Updated {{ .UpdatedDate.Format "Jan 2, 2006" }}</time>
  </article>
{{ end }}
```

---

**Related:**
- [Pages Configuration](/guide/configuration/pages)
- [Front Matter](/guide/configuration/frontmatter)
- [Content Structure](/guide/concepts/content-structure)
- [Template Mapping](/guide/templates/mapping)
