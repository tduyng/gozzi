# Templates & Partials

Gozzi uses Go's standard `html/template` package for rendering HTML, providing a powerful and flexible templating system. Your layouts and partials live in the `templates/` folder and define how each page or section of your site should look.

## Table of Contents

- [Template Structure](#template-structure)  
- [Development Server Features](#development-server-features)
- [Template Mapping](#template-mapping)
- [Template Variables](#template-variables)
- [Partials & Macros](#partials--macros)
- [Template Functions](#template-functions)
- [Advanced Patterns](#advanced-patterns)
- [Examples](#examples)

## Template Structure

### Basic Template Organization

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

### Template Hierarchy

Templates follow an inheritance pattern:
1. **Base templates** define overall structure
2. **Section templates** handle list pages  
3. **Page templates** render individual content
4. **Partials** provide reusable components
5. **Macros** offer advanced functionality

## Development Server Features

### 🔥 Live Template Reloading

When running `gozzi serve`, the development server provides instant template updates:

```sh
# Start development server
gozzi serve --port 1313

# Edit any template file:
# - templates/*.html
# - templates/partials/*.html  
# - templates/macros/*.html

# Browser automatically refreshes with changes!
```

#### How Live Reload Works

1. **File Watching**: Server monitors `templates/` directory for changes
2. **Template Reloading**: Automatically reloads template cache on file changes
3. **Smart Rebuilding**: Only regenerates affected pages
4. **Browser Refresh**: Sends refresh signal via Server-Sent Events (SSE)

#### Watched File Types

The development server watches these file extensions in `templates/`:
- `.html` - Template files
- `.css` - Stylesheets (in static folder)
- `.js` - JavaScript files (in static folder)
- `.md` - Content files
- `config.toml` - Configuration changes

### Development Context Variables

During development, templates have access to additional context:

```html
<!-- Development server info (only available during serve) -->
{{ if .IsDev }}
  <div id="dev-info">
    <p>Development Mode Active</p>
    <p>Port: {{ .DevPort }}</p>
    <p>Live Reload: Enabled</p>
  </div>
{{ end }}

<!-- Live reload script automatically injected -->
<!-- No need to manually add this during development -->
```

## Template Mapping

### Content-to-Template Mapping

Gozzi maps content to templates based on content type and location:

| Content | Template Used | Output Path | Description |
|---------|---------------|-------------|-------------|
| `content/_index.md` | `home.html` | `/index.html` | Site homepage |
| `content/blog/_index.md` | `blog.html` | `/blog/index.html` | Blog section list |
| `content/blog/post/index.md` | `post.html` | `/blog/post/index.html` | Individual blog post |
| `content/notes/_index.md` | `notes.html` | `/notes/index.html` | Notes section list |
| `content/notes/note/index.md` | `note.html` | `/notes/note/index.html` | Individual note |
| `content/about/index.md` | `prose.html` | `/about/index.html` | General page |

### Template Resolution Order

When selecting templates, Gozzi follows this priority:

1. **Front matter override**: `template = "custom.html"` in front matter
2. **Section-specific**: `{section}.html` (e.g., `blog.html`)
3. **Content type**: `post.html`, `note.html`, `prose.html`
4. **Default fallback**: `home.html` for root, `prose.html` for pages

## Template Variables

Each template receives a context object (`.`) with comprehensive site and content data.

### `.Site` - Global Configuration

Site-wide data from `config.toml`:

```html
<!-- Basic site info -->
<title>{{ .Site.Config.title }}</title>
<meta name="description" content="{{ .Site.Config.description }}">
<link rel="canonical" href="{{ .Site.Config.base_url }}">

<!-- Custom configuration -->
<p>{{ .Site.Config.extra.bio }}</p>
<img src="{{ .Site.Config.extra.avatar }}" alt="Profile">

<!-- Navigation data -->
{{ range .Site.Config.extra.sections }}
  <a href="{{ .path }}">{{ .name }}</a>
{{ end }}

<!-- Social links -->
{{ range .Site.Config.extra.links }}
  <a href="{{ .url }}">
    <img src="/icon/{{ .icon }}.svg" alt="{{ .name }}">
  </a>
{{ end }}
```

### `.Section` - Section Data

Available in section templates (list pages with `_index.md`):

```html
<!-- Section metadata -->
<h1>{{ .Section.Config.title }}</h1>
<p>{{ .Section.Config.description }}</p>

<!-- Section content -->
{{ .Section.Content }}

<!-- Custom section data -->
{{ if .Section.Config.extra.hero_text }}
  <div class="hero">{{ .Section.Config.extra.hero_text }}</div>
{{ end }}

<!-- List all children pages -->
{{ range .Section.Children }}
  <article>
    <h2><a href="{{ .Permalink }}">{{ .Config.title }}</a></h2>
    <time>{{ date .Config.date "January 2, 2006" }}</time>
    <p>{{ .Config.description }}</p>
    
    {{ if .Config.featured }}
      <span class="featured">Featured</span>
    {{ end }}
    
    <div class="tags">
      {{ range .Config.tags }}
        <span class="tag">#{{ . }}</span>
      {{ end }}
    </div>
  </article>
{{ end }}
```

### `.Page` - Individual Page Data

Used in page templates for individual content:

```html
<!-- Page metadata -->
<h1>{{ .Page.Config.title }}</h1>
<meta name="description" content="{{ .Page.Config.description }}">

<!-- Dates and timing -->
<time datetime="{{ .Page.Config.date.Format "2006-01-02" }}">
  {{ date .Page.Config.date "January 2, 2006" }}
</time>

{{ if not .Page.Config.updated.IsZero }}
  <p>Updated: {{ date .Page.Config.updated "January 2, 2006" }}</p>
{{ end }}

<!-- Content and metadata -->
{{ .Page.Content }}
<p>Permalink: {{ .Page.Permalink }}</p>

<!-- Navigation -->
{{ with .Page.Higher }}
  <a href="{{ .Permalink }}">← {{ .Config.title }}</a>
{{ end }}

{{ with .Page.Lower }}
  <a href="{{ .Permalink }}">{{ .Config.title }} →</a>
{{ end }}

<!-- Custom page data -->
{{ if .Page.Config.extra.toc }}
  {{ template "partials/_toc.html" . }}
{{ end }}

<!-- Tags and categorization -->
{{ range .Page.Config.tags }}
  <a href="/tags/{{ . | urlize }}">#{{ . }}</a>
{{ end }}
```

### Priority Function for Inheritance

Use the `priority` function to cascade values from page → section → site:

```html
<!-- Use page title, fall back to section, then site -->
<title>{{ priority .Page.Config.title .Section.Config.title .Site.Config.title }}</title>

<!-- Language inheritance -->
<html lang="{{ priority .Page.Config.lang .Section.Config.lang .Site.Config.lang }}">

<!-- Description with fallbacks -->
{{ $description := priority .Page.Config.description .Section.Config.description .Site.Config.description }}
<meta name="description" content="{{ $description }}">
```

## Partials & Macros

### Partials - Reusable Components

Partials are shared template snippets that promote DRY (Don't Repeat Yourself) principles:

#### Including Partials

```html
<!-- Basic inclusion with current context -->
{{ template "partials/_head.html" . }}
{{ template "partials/_header.html" . }}
{{ template "partials/_footer.html" . }}

<!-- Pass custom data to partial -->
{{ template "partials/_pagination.html" (dict "pages" .Section.Children "current" .Page) }}
```

#### Common Partials

**`_head.html` - Document Head**:
```html
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{ priority .Page.Config.title .Section.Config.title .Site.Config.title }}</title>
    
    <!-- SEO Meta Tags -->
    {{ $description := priority .Page.Config.description .Section.Config.description .Site.Config.description }}
    <meta name="description" content="{{ $description }}">
    
    <!-- Open Graph -->
    <meta property="og:title" content="{{ .Page.Config.title }}">
    <meta property="og:description" content="{{ $description }}">
    {{ if .Page.Config.img }}
    <meta property="og:image" content="{{ .Page.Config.img_url }}">
    {{ end }}
    
    <!-- Canonical URL -->
    <link rel="canonical" href="{{ .Page.URL }}">
    
    <!-- Stylesheets -->
    <link rel="stylesheet" href="/css/main.css">
</head>
```

**`_toc.html` - Table of Contents**:
```html
{{ if .Page.Config.extra.toc }}
<nav class="toc">
    <h3>Table of Contents</h3>
    <!-- TOC generated from page content -->
    {{ .Page.TOC }}
</nav>
{{ end }}
```

**`_sharing.html` - Social Sharing**:
```html
<div class="sharing">
    <h4>Share this post</h4>
    {{ $url := .Page.Permalink }}
    {{ $title := .Page.Config.title }}
    
    <a href="https://twitter.com/intent/tweet?url={{ $url }}&text={{ $title }}" target="_blank">
        Twitter
    </a>
    <a href="https://www.linkedin.com/sharing/share-offsite/?url={{ $url }}" target="_blank">
        LinkedIn  
    </a>
</div>
```

### Macros - Advanced Components

Macros provide complex, reusable functionality:

**Alert Macro** (`macros/alert.html`):
```html
<!-- Usage: {{ template "macros/alert.html" (dict "type" "info" "message" "Important info!") }} -->
<div class="alert alert-{{ .type }}">
    {{ if eq .type "warning" }}⚠️{{ end }}
    {{ if eq .type "info" }}ℹ️{{ end }}
    {{ if eq .type "success" }}✅{{ end }}
    {{ .message }}
</div>
```

**Code Block Macro** (`macros/codeblock.html`):
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

## Template Functions

Gozzi provides extensive template functions beyond Go's built-ins.

### Content Functions

```html
<!-- Get section data -->
{{ $blog := get_section "blog" }}
{{ range limit 5 (reverse $blog.Children) }}
  <a href="{{ .Permalink }}">{{ .Config.title }}</a>
{{ end }}

<!-- Render markdown inline -->
{{ $markdown := "**Bold text** and _italic_" }}
{{ markdown $markdown }}

<!-- Load external files -->
{{ $svg := load "static/icon/arrow.svg" }}
{{ $svg }}
```

### String Functions

```html
<!-- URL generation -->
{{ $slug := "My Great Post Title" | urlize }}
<!-- Output: my-great-post-title -->

<!-- String manipulation -->
{{ "hello world" | upper }}  <!-- HELLO WORLD -->
{{ "  trim me  " | trim }}   <!-- trim me -->
{{ "one,two,three" | split "," }}  <!-- ["one", "two", "three"] -->

<!-- Contains and prefix checking -->
{{ if contains .Page.Config.tags "go" }}
  <span class="go-post">Go Related</span>
{{ end }}

{{ if has_prefix .Page.URL "/blog/" }}
  <nav>Blog Navigation</nav>
{{ end }}
```

### Date Functions

```html
<!-- Format dates -->
{{ date .Page.Config.date "January 2, 2006" }}
{{ date .Page.Config.date "2006-01-02" }}
{{ date .Page.Config.date "Jan 2" }}

<!-- Current time -->
{{ $now := now }}
{{ date $now "2006-01-02 15:04:05" }}

<!-- Parse date strings -->
{{ $parsed := to_date "2025-01-15" }}
{{ date $parsed "January 2, 2006" }}
```

### Collection Functions

```html
<!-- Pagination -->
{{ $pages := limit 10 .Section.Children }}
{{ range $pages }}
  <article>{{ .Config.title }}</article>
{{ end }}

<!-- Grouping -->
{{ $groups := group_by .Section.Children "Config.tags" }}
{{ range $groups }}
  <h3>{{ .Key }}</h3>
  {{ range .Items }}
    <p>{{ .Config.title }}</p>
  {{ end }}
{{ end }}

<!-- Filtering -->
{{ $featured := where .Section.Children "Config.featured" true }}
{{ range $featured }}
  <div class="featured">{{ .Config.title }}</div>
{{ end }}

<!-- First and last -->
{{ $latest := first 3 (reverse .Section.Children) }}
{{ $oldest := last 3 .Section.Children }}
```

### Logic Functions

```html
<!-- Conditionals -->
{{ if and .Page.Config.featured (eq .Page.Config.draft false) }}
  <span class="featured-published">Featured & Published</span>
{{ end }}

{{ if or (eq .Section.Name "blog") (eq .Section.Name "notes") }}
  <div class="content-section">Content</div>
{{ end }}

<!-- Comparisons -->
{{ if ne .Page.Config.template "home.html" }}
  <nav>Secondary Navigation</nav>
{{ end }}

<!-- Default values -->
{{ $author := default "Anonymous" .Page.Config.extra.author }}
<p>By {{ $author }}</p>
```

### Asset Functions

```html
<!-- Asset path generation -->
{{ asset "/css/main.css" . }}  <!-- /css/main.css with cache busting -->
{{ asset "/js/app.js" . }}     <!-- /js/app.js?v=hash -->

<!-- Load and process files -->
{{ $styles := load "static/css/critical.css" }}
<style>{{ $styles }}</style>

<!-- Load as HTML attribute -->
{{ $icon := load_attribute "static/icon/menu.svg" }}
<div data-icon="{{ $icon }}">Menu</div>
```

## Advanced Patterns

### Template Inheritance

Create base templates for consistent structure:

**`base.html`**:
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

**`post.html`** (inherits from base):
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

### Conditional Template Loading

```html
<!-- Load different partials based on context -->
{{ if eq .Section.Name "blog" }}
  {{ template "partials/_blog_sidebar.html" . }}
{{ else if eq .Section.Name "notes" }}  
  {{ template "partials/_notes_sidebar.html" . }}
{{ else }}
  {{ template "partials/_default_sidebar.html" . }}
{{ end }}

<!-- Feature flags from config -->
{{ if .Site.Config.extra.enable_comments }}
  {{ template "partials/_comments.html" . }}
{{ end }}

{{ if .Site.Config.extra.enable_analytics }}
  {{ template "partials/_analytics.html" . }}
{{ end }}
```

### Data Processing

```html
<!-- Complex data manipulation -->
{{ $posts := get_section "blog" }}
{{ $featured := where $posts.Children "Config.featured" true }}
{{ $recent := first 5 (reverse (where $posts.Children "Config.draft" false)) }}

<section class="featured">
    <h2>Featured Posts</h2>
    {{ range $featured }}
        <article>
            <h3><a href="{{ .Permalink }}">{{ .Config.title }}</a></h3>
            <p>{{ .Config.description }}</p>
        </article>
    {{ end }}
</section>

<section class="recent">
    <h2>Recent Posts</h2>
    {{ range $recent }}
        <article>
            <h3><a href="{{ .Permalink }}">{{ .Config.title }}</a></h3>
            <time>{{ date .Config.date "Jan 2" }}</time>
        </article>
    {{ end }}
</section>
```

## Examples

### Complete Blog Template

**`blog.html`** - Blog section listing:
```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.lang }}">
{{ template "partials/_head.html" . }}
<body class="blog-list">
    {{ template "partials/_header.html" . }}
    
    <main>
        <header class="section-header">
            <h1>{{ .Section.Config.title }}</h1>
            <p>{{ .Section.Config.description }}</p>
            
            {{ if .Section.Config.extra.hero_text }}
            <div class="hero">{{ .Section.Config.extra.hero_text }}</div>
            {{ end }}
        </header>
        
        <div class="posts">
            {{ $posts := where .Section.Children "Config.draft" false }}
            {{ range $posts }}
            <article class="post-card">
                <h2><a href="{{ .Permalink }}">{{ .Config.title }}</a></h2>
                
                <div class="meta">
                    <time>{{ date .Config.date "January 2, 2006" }}</time>
                    {{ if .Config.featured }}
                        <span class="featured">Featured</span>
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
        
        {{ if gt (len $posts) 10 }}
            {{ template "macros/pagination.html" (dict "posts" $posts "per_page" 10) }}
        {{ end }}
    </main>
    
    {{ template "partials/_footer.html" . }}
</body>
</html>
```

### Advanced Post Template

**`post.html`** - Individual post with all features:
```html
<!DOCTYPE html>
<html lang="{{ priority .Page.Config.lang .Section.Config.lang .Site.Config.lang }}">
{{ template "partials/_head.html" . }}
<body class="post">
    {{ template "partials/_header.html" . }}
    
    <div id="wrapper">
        <aside class="sidebar">
            {{ if .Page.Config.extra.toc }}
                {{ template "partials/_toc.html" . }}
            {{ end }}
            
            <button id="back-to-top" aria-label="Back to top">
                {{ $icon := load "static/icon/arrow-up.svg" }}
                {{ $icon }}
            </button>
        </aside>
        
        <main>
            {{ template "partials/_copy_code.html" . }}
            
            <article class="prose">
                <header>
                    <h1>{{ .Page.Config.title }}</h1>
                    
                    <div class="post-meta">
                        <time datetime="{{ .Page.Config.date.Format "2006-01-02" }}">
                            {{ date .Page.Config.date "January 2, 2006" }}
                        </time>
                        
                        {{ if not .Page.Config.updated.IsZero }}
                        <span>Updated: 
                            <time datetime="{{ .Page.Config.updated.Format "2006-01-02" }}">
                                {{ date .Page.Config.updated "January 2, 2006" }}
                            </time>
                        </span>
                        {{ end }}
                        
                        {{ template "partials/_word_count.html" . }}
                    </div>
                    
                    {{ if .Page.Config.tags }}
                    <div class="tags">
                        {{ range .Page.Config.tags }}
                            {{ $tag_url := printf "/tags/%s" (. | urlize) }}
                            <a href="{{ $tag_url }}" class="tag">#{{ . }}</a>
                        {{ end }}
                    </div>
                    {{ end }}
                </header>
                
                {{ template "partials/_outdate.html" . }}
                
                <div class="content">
                    {{ .Page.Content }}
                </div>
                
                {{ pagination . }}
            </article>
            
            {{ template "partials/_sharing.html" . }}
            {{ template "partials/_reaction.html" . }}
            {{ template "partials/_comment.html" . }}
        </main>
    </div>
    
    {{ template "partials/_footer.html" . }}
    {{ template "partials/_scripts.html" . }}
</body>
</html>
```

### Dynamic Navigation

**`partials/_header.html`**:
```html
<header class="site-header">
    <div class="container">
        <a href="/" class="site-title">{{ .Site.Config.title }}</a>
        
        <nav class="main-nav">
            {{ range .Site.Config.extra.sections }}
                {{ $current := "" }}
                {{ if eq .path $.Page.Section }}
                    {{ $current = "current" }}
                {{ end }}
                
                <a href="{{ .path }}" class="nav-link {{ $current }}">
                    {{ .name }}
                </a>
            {{ end }}
        </nav>
        
        <div class="social-links">
            {{ range .Site.Config.extra.links }}
                <a href="{{ .url }}" class="social-link" title="{{ .name }}">
                    {{ $icon_path := printf "static/icon/%s.svg" .icon }}
                    {{ load $icon_path }}
                </a>
            {{ end }}
        </div>
    </div>
</header>
```

## Output Directory Structure

After building, your site structure mirrors the content organization:

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

## Best Practices

### Performance Optimization

- **Minimize template nesting** to reduce compilation time
- **Cache frequently used data** in variables
- **Use partials** to avoid code duplication
- **Preload critical assets** in the head section

### Development Workflow

1. **Start development server**: `gozzi serve --port 1313`
2. **Edit templates**: Changes are instantly reflected
3. **Test across sections**: Ensure templates work for all content types
4. **Validate HTML**: Use browser dev tools to check output
5. **Build for production**: `gozzi build` to generate final site

### Security Considerations

- **Use `safe` function carefully**: Only for trusted HTML content
- **Escape user content**: Let Go templates handle escaping automatically
- **Validate external data**: Check loaded files and user inputs

## Next Steps

- [Configuration](/guide/configuration) - Configure site settings and front matter
- [HTML Functions](/reference/template-functions) - Complete template function reference  
- [Content Structure](/guide/content-structure) - Organize your content effectively
- [CLI Reference](/reference/cli) - Development server and build commands
