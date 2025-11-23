# Templates & Partials

Gozzi uses Go's standard `html/template` package for rendering HTML, providing a powerful and flexible templating system.

## Templates Guide

### Getting Started
- **[Template Structure](/guide/templates/structure)** - Directory layout and file organization
- **[Template Mapping](/guide/templates/mapping)** - How content maps to templates
- **[Development Server](/guide/templates/development)** - Live reload and hot reloading

### Template Data
- **[Template Variables](/guide/templates/variables)** - Available data: Site, Section, Page
- **[Template Functions](/guide/templates/functions)** - Built-in functions for templates
- **[Priority & Inheritance](/guide/templates/inheritance)** - Configuration cascading

### Components
- **[Partials](/guide/templates/partials)** - Reusable template snippets
- **[Macros](/guide/templates/macros)** - Advanced reusable components
- **[Advanced Patterns](/guide/templates/advanced)** - Complex template patterns

### Examples
- **[Complete Examples](/guide/templates/examples)** - Full template implementations

## Quick Example

**Basic post template:**
```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.lang }}">
{{ template "partials/_head.html" . }}
<body>
    {{ template "partials/_header.html" . }}
    
    <main>
        <article>
            <h1>{{ .Page.Config.title }}</h1>
            <time>{{ date .Page.Config.date "January 2, 2006" }}</time>
            {{ .Page.Content }}
        </article>
    </main>
    
    {{ template "partials/_footer.html" . }}
</body>
</html>
```

## Directory Structure

```
templates/
├── home.html            # Homepage
├── blog.html            # Blog list
├── post.html            # Blog post
├── partials/            # Shared snippets
│   ├── _head.html
│   ├── _header.html
│   └── _footer.html
└── macros/              # Reusable components
    ├── alert.html
    └── pagination.html
```

---

**Start with:** [Template Structure](/guide/templates/structure) to understand template organization.
