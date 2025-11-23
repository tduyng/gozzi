# Template Inheritance

Create base templates for consistent structure across your site, reducing duplication and maintaining consistency.

## Base Template Pattern

### Creating a Base Template

**`templates/base.html`**:
```html
<!DOCTYPE html>
<html lang="{{ priority .Page.Config.lang .Section.Config.lang .Site.Config.lang }}">
{{ template "partials/_head.html" . }}
<body class="{{ block "body_class" . }}page{{ end }}">
    {{ template "partials/_header.html" . }}
    
    <main>
        {{ block "content" . }}
        <p>Default content</p>
        {{ end }}
    </main>
    
    {{ template "partials/_footer.html" . }}
    {{ template "partials/_scripts.html" . }}
</body>
</html>
```

### Extending Base Template

**`templates/post.html`** (inherits from base):
```html
{{ template "base.html" . }}

{{ define "body_class" }}post{{ end }}

{{ define "content" }}
<article class="prose">
    <h1>{{ .Page.Config.title }}</h1>
    
    <div class="meta">
        <time>{{ date .Page.Config.date "January 2, 2006" }}</time>
        {{ template "partials/_word_count.html" . }}
    </div>
    
    {{ .Page.Content }}
    
    {{ template "partials/_pagination.html" . }}
</article>
{{ end }}
```

## Blocks and Definitions

### `{{ block }}`

Define overridable sections:

```html
<!-- Base template -->
<body class="{{ block "body_class" . }}default{{ end }}">
    {{ block "header" . }}
        <header>Default Header</header>
    {{ end }}
    
    <main>
        {{ block "content" . }}
            <p>No content defined</p>
        {{ end }}
    </main>
    
    {{ block "sidebar" . }}
        <!-- Optional sidebar -->
    {{ end }}
</body>
```

### `{{ define }}`

Override blocks in child templates:

```html
<!-- Child template -->
{{ template "base.html" . }}

{{ define "body_class" }}custom-page{{ end }}

{{ define "header" }}
    <header class="custom">
        <h1>Custom Header</h1>
    </header>
{{ end }}

{{ define "content" }}
    <article>
        <h1>{{ .Page.Config.title }}</h1>
        {{ .Page.Content }}
    </article>
{{ end }}

<!-- sidebar block not defined = not rendered -->
```

## Complete Inheritance Example

### Base Layout

**`templates/layouts/base.html`**:
```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.language }}">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    
    {{ block "head_meta" . }}
        <title>{{ .Site.Config.title }}</title>
        <meta name="description" content="{{ .Site.Config.description }}">
    {{ end }}
    
    <link rel="stylesheet" href="/css/main.css">
    {{ block "head_extra" . }}{{ end }}
</head>
<body class="{{ block "body_class" . }}page{{ end }}">
    {{ template "partials/_header.html" . }}
    
    <div class="container">
        {{ block "before_content" . }}{{ end }}
        
        <main class="{{ block "main_class" . }}content{{ end }}">
            {{ block "content" . }}
                <p>Default content</p>
            {{ end }}
        </main>
        
        {{ block "after_content" . }}{{ end }}
    </div>
    
    {{ template "partials/_footer.html" . }}
    
    {{ block "scripts" . }}
        <script src="/js/main.js"></script>
    {{ end }}
</body>
</html>
```

### Blog Post Template

**`templates/post.html`**:
```html
{{ template "layouts/base.html" . }}

{{ define "body_class" }}post single{{ end }}

{{ define "head_meta" }}
    <title>{{ .Page.Config.title }} - {{ .Site.Config.title }}</title>
    <meta name="description" content="{{ .Page.Config.description }}">
    <meta property="og:type" content="article">
    <meta property="og:title" content="{{ .Page.Config.title }}">
{{ end }}

{{ define "head_extra" }}
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.css">
{{ end }}

{{ define "main_class" }}prose{{ end }}

{{ define "content" }}
<article>
    <header>
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
    </header>
    
    <div class="content">
        {{ .Page.Content }}
    </div>
</article>
{{ end }}

{{ define "after_content" }}
    {{ template "partials/_sharing.html" . }}
    {{ template "partials/_comments.html" . }}
{{ end }}
```

### Blog List Template

**`templates/blog.html`**:
```html
{{ template "layouts/base.html" . }}

{{ define "body_class" }}blog list{{ end }}

{{ define "head_meta" }}
    <title>{{ .Section.Config.title }} - {{ .Site.Config.title }}</title>
    <meta name="description" content="{{ .Section.Config.description }}">
{{ end }}

{{ define "before_content" }}
    <header class="section-header">
        <h1>{{ .Section.Config.title }}</h1>
        <p>{{ .Section.Config.description }}</p>
    </header>
{{ end }}

{{ define "content" }}
    {{ range .Section.Children }}
    <article class="post-card">
        <h2><a href="{{ .Permalink }}">{{ .Config.title }}</a></h2>
        <time>{{ date .Config.date "Jan 2, 2006" }}</time>
        <p>{{ .Config.description }}</p>
    </article>
    {{ end }}
{{ end }}
```

## Multi-Level Inheritance

### Three-Level Structure

```
base.html                    # Level 1: Foundation
  ├── content-base.html      # Level 2: Content layouts
  │   ├── post.html          # Level 3: Specific pages
  │   └── note.html
  └── list-base.html         # Level 2: List layouts
      ├── blog.html          # Level 3: Specific lists
      └── notes.html
```

**Level 1 - Foundation**:
```html
<!-- templates/base.html -->
<!DOCTYPE html>
<html>
{{ template "partials/_head.html" . }}
<body class="{{ block "body_class" . }}page{{ end }}">
    {{ block "layout" . }}
        <p>No layout defined</p>
    {{ end }}
</body>
</html>
```

**Level 2 - Content Layout**:
```html
<!-- templates/content-base.html -->
{{ template "base.html" . }}

{{ define "layout" }}
    {{ template "partials/_header.html" . }}
    <main>
        {{ block "article" . }}
            <article>Default article</article>
        {{ end }}
    </main>
    {{ template "partials/_footer.html" . }}
{{ end }}
```

**Level 3 - Specific Page**:
```html
<!-- templates/post.html -->
{{ template "content-base.html" . }}

{{ define "body_class" }}post{{ end }}

{{ define "article" }}
    <article class="prose">
        <h1>{{ .Page.Config.title }}</h1>
        {{ .Page.Content }}
    </article>
{{ end }}
```

## Best Practices

1. **Keep base templates simple** - Minimal structure
2. **Use semantic block names** - `content`, `header`, `sidebar`
3. **Provide sensible defaults** - Base template should render standalone
4. **Document blocks** - Comment what each block does
5. **Avoid deep nesting** - 2-3 levels maximum
6. **Test inheritance chain** - Ensure all templates work

## Common Patterns

### Optional Sidebar

```html
<!-- Base with optional sidebar -->
<div class="layout">
    <main>{{ block "content" . }}{{ end }}</main>
    
    {{ block "sidebar" . }}
        <!-- No sidebar by default -->
    {{ end }}
</div>

<!-- Add sidebar in child -->
{{ define "sidebar" }}
    <aside>
        <h3>Related Posts</h3>
        <!-- Sidebar content -->
    </aside>
{{ end }}
```

### Conditional Blocks

```html
<!-- Base template -->
{{ if block "has_toc" . }}false{{ end }}
    <aside class="toc">
        {{ block "toc_content" . }}{{ end }}
    </aside>
{{ end }}

<!-- Enable in child -->
{{ define "has_toc" }}true{{ end }}
{{ define "toc_content" }}
    {{ .Toc }}
{{ end }}
```

---

**Related:**
- [Template Structure](/guide/templates/structure)
- [Partials & Macros](/guide/templates/partials)
- [Examples](/guide/templates/examples)
