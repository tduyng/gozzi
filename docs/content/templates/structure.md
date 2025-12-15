+++
title = "Template Structure"
date = 2025-12-15
template = "page.html"
+++

Gozzi uses Go's standard `html/template` package for rendering HTML. Your layouts and partials live in the `templates/` folder and define how each page or section of your site should look.

## Basic Template Organization

```
templates/
├── home.html            # Homepage layout
├── blog.html            # Blog section list layout  
├── post.html            # Blog post (single page)
├── notes.html           # Notes list
├── note.html            # Single note
├── prose.html           # General prose pages
├── tags.html            # All tags index
├── tag.html             # Posts under a specific tag
├── 404.html             # Custom 404 page
├── partials/            # Shared snippets
│   ├── _head.html
│   ├── _header.html
│   ├── _footer.html
│   ├── _scripts.html
│   ├── _toc.html        # Table of contents
│   ├── _word_count.html # Reading time
│   ├── _sharing.html    # Social sharing
│   ├── _comment.html    # Comments section
│   └── _reaction.html   # Reaction buttons
└── macros/              # Reusable components
    ├── alert.html
    ├── codeblock.html
    ├── important.html
    ├── mermaid.html
    ├── pagination.html
    └── warning.html
```

## Template Hierarchy

Templates follow an inheritance pattern:

1. **Base templates** - Define overall structure
2. **Section templates** - Handle list pages  
3. **Page templates** - Render individual content
4. **Partials** - Provide reusable components
5. **Macros** - Offer advanced functionality

## Template Types

### Base Templates

Foundation templates that define site-wide structure:

```html
<!-- base.html -->
<!DOCTYPE html>
<html lang="{{ .Site.Config.lang }}">
{{ template "partials/_head.html" . }}
<body>
    {{ template "partials/_header.html" . }}
    
    <main>
        {{ block "content" . }}
        <p>Default content</p>
        {{ end }}
    </main>
    
    {{ template "partials/_footer.html" . }}
</body>
</html>
```

### Section Templates

List pages for content collections (blogs, notes, etc.):

```html
<!-- blog.html -->
<!DOCTYPE html>
<html>
{{ template "partials/_head.html" . }}
<body>
    <h1>{{ .Section.Config.title }}</h1>
    
    {{ range .Section.Children }}
        <article>
            <h2><a href="{{ .Permalink }}">{{ .Config.title }}</a></h2>
            <time>{{ date .Config.date "Jan 2, 2006" }}</time>
        </article>
    {{ end }}
</body>
</html>
```

### Page Templates

Individual content rendering:

```html
<!-- post.html -->
<!DOCTYPE html>
<html>
{{ template "partials/_head.html" . }}
<body>
    <article>
        <h1>{{ .Page.Config.title }}</h1>
        <time>{{ date .Page.Config.date "January 2, 2006" }}</time>
        
        <div class="content">
            {{ .Page.Content }}
        </div>
    </article>
</body>
</html>
```

## Directory Organization

### Partials Directory

Store reusable components:

```
templates/partials/
├── _head.html       # <head> section
├── _header.html     # Site header/navigation
├── _footer.html     # Site footer
├── _scripts.html    # JavaScript includes
├── _toc.html        # Table of contents
├── _word_count.html # Reading time display
└── _sharing.html    # Social sharing buttons
```

**Naming convention:** Prefix with underscore (`_`) to indicate partial templates.

### Macros Directory

Complex, reusable components:

```
templates/macros/
├── alert.html       # Alert boxes (info, warning, error)
├── codeblock.html   # Enhanced code blocks
├── pagination.html  # Pagination controls
└── mermaid.html     # Mermaid diagram wrapper
```

## Output Directory Structure

After building, your site structure mirrors content organization:

```
public/
├── index.html           # Home page
├── blog/
│   ├── index.html       # Blog section list
│   └── my-post/
│       └── index.html   # Individual post
├── notes/
│   ├── index.html       # Notes section list  
│   └── my-note/
│       └── index.html   # Individual note
├── tags/
│   ├── index.html       # All tags
│   └── go/
│       └── index.html   # Posts tagged "go"
├── css/
│   └── main.css         # Stylesheets
├── js/
│   └── main.js          # JavaScript
└── img/
    └── avatar.webp      # Images and assets
```

## Template Context

Each template receives a context object (`.`) with data:

- **`.Site`** - Global site configuration
- **`.Section`** - Section data (for list pages)
- **`.Page`** - Page data (for single pages)
- **`.Toc`** - Table of contents
- **Custom data** from `[extra]` configuration

## Including Templates

### Basic Inclusion

```html
<!-- Include partial with current context -->
{{ template "partials/_head.html" . }}
{{ template "partials/_header.html" . }}
{{ template "partials/_footer.html" . }}
```

### With Custom Data

```html
<!-- Pass specific data to partial -->
{{ template "partials/_pagination.html" (dict "items" .Section.Children "per_page" 10) }}
```

## Best Practices

1. **Consistent naming** - Use clear, descriptive names
2. **Organize by purpose** - Group related templates
3. **Use partials** - Avoid code duplication
4. **Prefix partials** - Use underscore for reusable components
5. **Document complex logic** - Add comments for clarity
6. **Test all templates** - Ensure they work for all content types
7. **Keep templates focused** - One responsibility per template

## Common Template Pattern

```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.lang }}">
{{ template "partials/_head.html" . }}
<body class="{{ .Section.Name }}">
    {{ template "partials/_header.html" . }}
    
    <main>
        <!-- Template-specific content -->
    </main>
    
    {{ template "partials/_footer.html" . }}
    {{ template "partials/_scripts.html" . }}
</body>
</html>
```

---

**Related:**
- [Template Mapping](/templates/mapping)
- [Template Variables](/templates/variables)
- [Partials & Macros](/templates/partials)
- [Development Workflow](/templates/development)
