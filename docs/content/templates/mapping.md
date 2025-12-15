+++
title = "Template Mapping"
date = 2025-12-15
template = "page.html"
+++

Gozzi automatically maps content to templates based on content type and location. Understanding this mapping helps you create the right template files for your content.

## Content-to-Template Mapping

| Content | Template Used | Output Path | Description |
|---------|---------------|-------------|-------------|
| `content/_index.md` | `home.html` | `/index.html` | Site homepage |
| `content/blog/_index.md` | `blog.html` | `/blog/index.html` | Blog section list |
| `content/blog/post/index.md` | `post.html` | `/blog/post/index.html` | Individual blog post |
| `content/notes/_index.md` | `notes.html` | `/notes/index.html` | Notes section list |
| `content/notes/note/index.md` | `note.html` | `/notes/note/index.html` | Individual note |
| `content/about/index.md` | `prose.html` | `/about/index.html` | General page |

## Template Resolution Order

When selecting templates, Gozzi follows this priority:

1. **Front matter override** - `template = "custom.html"` in front matter
2. **Section-specific** - `{section}.html` (e.g., `blog.html`)
3. **Content type** - `post.html`, `note.html`, `prose.html`
4. **Default fallback** - `home.html` for root, `prose.html` for pages

## Front Matter Override

Override automatic mapping with front matter:

```toml
+++
title = "About Me"
template = "prose.html"  # Force use of prose.html
+++
```

```toml
+++
title = "Special Post"
template = "custom-layout.html"  # Use custom template
+++
```

## Section-Based Mapping

Sections automatically map to templates by name:

```
content/
├── blog/
│   └── _index.md        → templates/blog.html
├── notes/
│   └── _index.md        → templates/notes.html
└── projects/
    └── _index.md        → templates/projects.html
```

If section-specific template doesn't exist, falls back to generic templates.

## Page-Type Mapping

Individual pages within sections:

```
content/blog/
├── _index.md            → blog.html (section list)
├── tutorial.md          → post.html (single post)
└── announcement.md      → post.html (single post)
```

```
content/notes/
├── _index.md            → notes.html (section list)
├── quick-tip.md         → note.html (single note)
└── til.md               → note.html (single note)
```

## Special Templates

### Homepage

```
content/_index.md → templates/home.html
```

Always uses `home.html` for the site root.

### 404 Page

```
templates/404.html → public/404.html
```

Custom error page for missing content.

### Tag Pages

```
# All tags listing
/tags/ → templates/tags.html

# Single tag page
/tags/golang/ → templates/tag.html
```

## Fallback Templates

When specific templates don't exist:

```
Missing blog.html?     → Use section.html
Missing post.html?     → Use page.html
Missing note.html?     → Use page.html
Missing section.html?  → Use prose.html
Missing page.html?     → Error (required)
```

## Template Examples by Type

### Homepage Template

```html
<!-- templates/home.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Site.Config.title }}</title>
</head>
<body>
    <h1>Welcome to {{ .Site.Config.title }}</h1>
    <p>{{ .Site.Config.description }}</p>
    
    {{ $blog := get_section "blog" }}
    <h2>Recent Posts</h2>
    {{ range first 5 (reverse $blog.Children) }}
        <article>
            <h3><a href="{{ .Permalink }}">{{ .Config.title }}</a></h3>
        </article>
    {{ end }}
</body>
</html>
```

### Section List Template

```html
<!-- templates/blog.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Section.Config.title }} - {{ .Site.Config.title }}</title>
</head>
<body>
    <h1>{{ .Section.Config.title }}</h1>
    
    {{ range .Section.Children }}
        <article>
            <h2><a href="{{ .Permalink }}">{{ .Config.title }}</a></h2>
            <time>{{ date .Config.date "Jan 2, 2006" }}</time>
            <p>{{ .Config.description }}</p>
        </article>
    {{ end }}
</body>
</html>
```

### Single Page Template

```html
<!-- templates/post.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Page.Config.title }} - {{ .Site.Config.title }}</title>
</head>
<body>
    <article>
        <h1>{{ .Page.Config.title }}</h1>
        <time>{{ date .Page.Config.date "January 2, 2006" }}</time>
        
        <div class="content">
            {{ .Page.Content }}
        </div>
        
        {{ if .Page.Config.tags }}
            <div class="tags">
                {{ range .Page.Config.tags }}
                    <a href="/tags/{{ . | urlize }}">#{{ . }}</a>
                {{ end }}
            </div>
        {{ end }}
    </article>
</body>
</html>
```

## Custom Mapping Strategy

### By Content Type

Create different templates for content types:

```
templates/
├── blog.html       # Blog listings
├── post.html       # Blog posts
├── notes.html      # Notes listings
├── note.html       # Single notes
├── projects.html   # Projects listings
└── project.html    # Single projects
```

### By Section Name

Match template name to section folder:

```
content/tutorials/   → templates/tutorials.html
content/guides/      → templates/guides.html
content/docs/        → templates/docs.html
```

### By Front Matter

Explicit template selection:

```toml
# Standard blog post
+++
title = "My Post"
# Uses post.html (auto-detected)
+++

# Special layout
+++
title = "Landing Page"
template = "landing.html"
+++

# Documentation style
+++
title = "API Reference"
template = "docs.html"
+++
```

## Debugging Template Selection

### Check Which Template is Used

Add debug info to templates (development only):

```html
{{ if eq .Site.Environment "development" }}
    <div id="debug">
        Template: {{ .TemplateName }}
        Section: {{ .Section.Name }}
        Content Type: {{ .Page.Type }}
    </div>
{{ end }}
```

### Verify Template Existence

Build output shows which template was used:

```bash
gozzi build --verbose
# Output shows template resolution for each page
```

## Best Practices

1. **Name templates clearly** - Match content structure
2. **Use fallbacks** - Provide generic templates
3. **Document custom mappings** - Comment complex logic
4. **Test all content types** - Ensure templates work for all pages
5. **Avoid over-customization** - Too many templates = harder maintenance
6. **Use front matter sparingly** - Prefer automatic mapping

## Common Patterns

### Blog + Notes Site

```
templates/
├── home.html      # Homepage
├── blog.html      # Blog listing
├── post.html      # Blog posts
├── notes.html     # Notes listing
├── note.html      # Single notes
└── prose.html     # Other pages
```

### Documentation Site

```
templates/
├── home.html      # Homepage
├── docs.html      # Docs section
├── doc.html       # Single doc page
└── prose.html     # Other pages
```

### Portfolio Site

```
templates/
├── home.html       # Homepage
├── projects.html   # Projects listing
├── project.html    # Single project
├── blog.html       # Blog listing
├── post.html       # Blog posts
└── prose.html      # About, contact, etc.
```

---

**Related:**
- [Template Structure](/templates/structure)
- [Template Variables](/templates/variables)
- [Content Structure](/content-structure)
- [Front Matter Configuration](/config/frontmatter)
