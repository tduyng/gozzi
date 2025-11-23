# Macros

Macros provide complex, reusable functionality with advanced features beyond simple partials.

## What are Macros?

Macros are advanced components in `templates/macros/` that:

- Accept parameters for customization
- Provide complex functionality
- Encapsulate behavior and styling
- Offer consistent UI components

## Common Macros

### Alert Macro

**`macros/alert.html`**:
```html
<!-- Usage: {{ template "macros/alert.html" (dict "type" "info" "message" "Important info!") }} -->
<div class="alert alert-{{ .type }}">
    {{ if eq .type "warning" }}⚠️{{ end }}
    {{ if eq .type "info" }}ℹ️{{ end }}
    {{ if eq .type "success" }}✅{{ end }}
    {{ if eq .type "error" }}❌{{ end }}
    {{ .message }}
</div>
```

**Usage:**
```html
{{ template "macros/alert.html" (dict "type" "warning" "message" "This is a warning!") }}
{{ template "macros/alert.html" (dict "type" "info" "message" "Helpful information") }}
{{ template "macros/alert.html" (dict "type" "success" "message" "Operation successful") }}
```

### Code Block Macro

**`macros/codeblock.html`**:
```html
<!-- Enhanced code blocks with copy functionality -->
<div class="codeblock">
    <div class="codeblock-header">
        <span class="language">{{ .lang }}</span>
        <button class="copy-btn" data-clipboard-text="{{ .code }}">Copy</button>
    </div>
    <pre><code class="language-{{ .lang }}">{{ .code }}</code></pre>
</div>
```

**Usage:**
```html
{{ $code := `package main
import "fmt"
func main() {
    fmt.Println("Hello!")
}` }}
{{ template "macros/codeblock.html" (dict "lang" "go" "code" $code) }}
```

### Pagination Macro

**`macros/pagination.html`**:
```html
<nav class="pagination">
    {{ $currentPage := .current | default 1 }}
    {{ $totalPages := div (len .items) .per_page }}
    
    {{ if gt $currentPage 1 }}
        <a href="?page={{ sub $currentPage 1 }}" class="prev">← Previous</a>
    {{ end }}
    
    <span class="page-info">
        Page {{ $currentPage }} of {{ $totalPages }}
    </span>
    
    {{ if lt $currentPage $totalPages }}
        <a href="?page={{ add $currentPage 1 }}" class="next">Next →</a>
    {{ end }}
</nav>
```

**Usage:**
```html
{{ template "macros/pagination.html" (dict 
    "items" .Section.Children 
    "per_page" 10 
    "current" 1
) }}
```

### Card Macro

**`macros/card.html`**:
```html
<div class="card {{ .class }}">
    {{ if .image }}
        <div class="card-image">
            <img src="{{ .image }}" alt="{{ .title }}">
        </div>
    {{ end }}
    
    <div class="card-content">
        <h3 class="card-title">{{ .title }}</h3>
        
        {{ if .subtitle }}
            <p class="card-subtitle">{{ .subtitle }}</p>
        {{ end }}
        
        <div class="card-body">
            {{ .body }}
        </div>
        
        {{ if .link }}
            <a href="{{ .link }}" class="card-link">{{ .linkText | default "Read more" }} →</a>
        {{ end }}
    </div>
</div>
```

**Usage:**
```html
{{ template "macros/card.html" (dict 
    "title" "Getting Started"
    "subtitle" "Quick introduction"
    "body" "Learn the basics..."
    "link" "/guide/getting-started"
    "linkText" "Start learning"
    "class" "featured"
) }}
```

### Badge Macro

**`macros/badge.html`**:
```html
<span class="badge badge-{{ .type }}">
    {{ if .icon }}
        <span class="badge-icon">{{ .icon }}</span>
    {{ end }}
    <span class="badge-text">{{ .text }}</span>
</span>
```

**Usage:**
```html
{{ template "macros/badge.html" (dict "type" "success" "text" "Published") }}
{{ template "macros/badge.html" (dict "type" "warning" "text" "Draft" "icon" "⚠️") }}
{{ template "macros/badge.html" (dict "type" "info" "text" "Featured" "icon" "⭐") }}
```

### Breadcrumbs Macro

**`macros/breadcrumbs.html`**:
```html
<nav class="breadcrumbs">
    <a href="/">Home</a>
    {{ range .path }}
        <span class="separator">/</span>
        <a href="{{ .url }}">{{ .name }}</a>
    {{ end }}
</nav>
```

**Usage:**
```html
{{ template "macros/breadcrumbs.html" (dict "path" (slice
    (dict "name" "Blog" "url" "/blog")
    (dict "name" "2025" "url" "/blog/2025")
    (dict "name" .Page.Config.title "url" .Page.Permalink)
)) }}
```

### Mermaid Diagram Macro

**`macros/mermaid.html`**:
```html
<div class="mermaid-wrapper">
    <pre class="mermaid">
{{ .diagram }}
    </pre>
</div>
```

**Usage:**
```html
{{ $diagram := `graph TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Action]
    B -->|No| D[End]` }}
{{ template "macros/mermaid.html" (dict "diagram" $diagram) }}
```

## Creating Custom Macros

### Step 1: Identify Reusable Pattern

Look for repeated code across templates:

```html
<!-- Repeated in multiple places -->
<div class="author-bio">
    <img src="/img/{{ .author }}.jpg" alt="{{ .author }}">
    <div>
        <h4>{{ .author }}</h4>
        <p>{{ .author_bio }}</p>
    </div>
</div>
```

### Step 2: Extract to Macro

**`macros/author_bio.html`**:
```html
<div class="author-bio">
    <img src="{{ .avatar | default "/img/default-avatar.jpg" }}" alt="{{ .name }}">
    <div class="author-info">
        <h4 class="author-name">{{ .name }}</h4>
        {{ if .bio }}
            <p class="author-bio">{{ .bio }}</p>
        {{ end }}
        {{ if .links }}
            <div class="author-links">
                {{ range .links }}
                    <a href="{{ .url }}">{{ .name }}</a>
                {{ end }}
            </div>
        {{ end }}
    </div>
</div>
```

### Step 3: Use Macro

```html
{{ template "macros/author_bio.html" (dict 
    "name" .Page.Config.extra.author
    "bio" .Page.Config.extra.author_bio
    "avatar" .Page.Config.extra.author_avatar
    "links" .Page.Config.extra.author_links
) }}
```

## Macro Best Practices

1. **Accept parameters via dict** - Flexible and explicit
2. **Provide defaults** - Graceful fallbacks for missing data
3. **Document parameters** - Comment expected inputs
4. **Keep focused** - One component per macro
5. **Style consistently** - Match site design system
6. **Test edge cases** - Handle missing/empty data
7. **Version carefully** - Changes affect all usages

## Advanced Patterns

### Macro with Conditional Features

```html
<!-- macros/post_card.html -->
<article class="post-card {{ .class }}">
    <h2><a href="{{ .url }}">{{ .title }}</a></h2>
    
    {{ if .showDate }}
        <time>{{ date .date "Jan 2, 2006" }}</time>
    {{ end }}
    
    {{ if .showTags }}
        <div class="tags">
            {{ range .tags }}
                <span>#{{ . }}</span>
            {{ end }}
        </div>
    {{ end }}
    
    {{ if .showExcerpt }}
        <p>{{ .excerpt }}</p>
    {{ end }}
    
    {{ if .showReadMore }}
        <a href="{{ .url }}" class="read-more">Read more →</a>
    {{ end }}
</article>
```

### Nested Macros

```html
<!-- macros/feature_grid.html -->
<div class="feature-grid">
    {{ range .features }}
        {{ template "macros/feature_card.html" . }}
    {{ end }}
</div>
```

### Macro with Validation

```html
<!-- macros/link_button.html -->
{{ if and .url .text }}
    <a href="{{ .url }}" class="btn btn-{{ .style | default "primary" }}">
        {{ .text }}
    </a>
{{ else }}
    <!-- Error: Missing required parameters -->
    <div class="error">Link button requires url and text</div>
{{ end }}
```

## Macros vs Partials

**Use Partials for:**
- Structural components (header, footer)
- Fixed layouts
- Simple inclusions
- Site-wide elements

**Use Macros for:**
- Reusable UI components
- Parameterized widgets
- Conditional rendering
- Complex functionality

## Example Library

Build a macro library for consistent components:

```
templates/macros/
├── alert.html          # Alert boxes
├── badge.html          # Status badges
├── button.html         # Styled buttons
├── card.html           # Content cards
├── breadcrumbs.html    # Navigation
├── pagination.html     # Page navigation
├── tabs.html           # Tab interface
├── accordion.html      # Collapsible sections
└── modal.html          # Modal dialogs
```

---

**Related:**
- [Partials](/guide/templates/partials)
- [Template Functions](/guide/templates/functions)
- [Examples](/guide/templates/examples)
