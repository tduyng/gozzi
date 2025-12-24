+++
title = "Template Examples"
date = 2025-12-15
template = "page.html"
+++

Complete, real-world template examples you can use as starting points for your Gozzi site.

## Complete Blog Template

**`templates/blog.html`** - Blog section listing:

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
            {{ template "partials/pagination.html" (dict "posts" $posts "per_page" 10) }}
        {{ end }}
    </main>
    
    {{ template "partials/_footer.html" . }}
</body>
</html>
```

## Advanced Post Template

**`templates/post.html`** - Individual post with all features:

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
                
                {{ template "partials/_pagination.html" . }}
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

## Homepage Template

**`templates/home.html`** - Dynamic homepage:

```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.lang }}">
{{ template "partials/_head.html" . }}
<body class="home">
    {{ template "partials/_header.html" . }}
    
    <main>
        <!-- Hero Section -->
        <section class="hero">
            <h1>{{ .Site.Config.title }}</h1>
            <p class="lead">{{ .Site.Config.description }}</p>
            
            {{ if .Site.Config.extra.hero_cta }}
                <a href="{{ .Site.Config.extra.hero_cta_url }}" class="cta-button">
                    {{ .Site.Config.extra.hero_cta_text }}
                </a>
            {{ end }}
        </section>
        
        <!-- Featured Posts -->
        {{ $blog := get_section "blog" }}
        {{ $featured := where $blog.Children "Config.featured" true }}
        {{ if gt (len $featured) 0 }}
        <section class="featured-posts">
            <h2>Featured Posts</h2>
            <div class="post-grid">
                {{ range first 3 $featured }}
                    {{ template "partials/post_card.html" (dict 
                        "title" .Config.title
                        "description" .Config.description
                        "url" .Permalink
                        "date" .Config.date
                        "tags" .Config.tags
                    ) }}
                {{ end }}
            </div>
        </section>
        {{ end }}
        
        <!-- Recent Posts -->
        <section class="recent-posts">
            <h2>Recent Posts</h2>
            {{ range first 5 (reverse $blog.Children) }}
                <article class="post-preview">
                    <h3><a href="{{ .Permalink }}">{{ .Config.title }}</a></h3>
                    <time>{{ date .Config.date "Jan 2, 2006" }}</time>
                    <p>{{ .Config.description }}</p>
                </article>
            {{ end }}
            <a href="/blog/" class="view-all">View all posts →</a>
        </section>
        
        <!-- About Section -->
        {{ if .Site.Config.extra.bio }}
        <section class="about">
            <h2>About</h2>
            <p>{{ .Site.Config.extra.bio }}</p>
        </section>
        {{ end }}
    </main>
    
    {{ template "partials/_footer.html" . }}
</body>
</html>
```

## Dynamic Navigation

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

## Tag Pages

### All Tags Index

**`templates/tags.html`**:

```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.lang }}">
{{ template "partials/_head.html" . }}
<body class="tags-index">
    {{ template "partials/_header.html" . }}
    
    <main>
        <h1>All Tags</h1>
        
        <div class="tag-cloud">
            {{ $allTags := .Site.Tags }}
            {{ range $allTags }}
                {{ $tagUrl := printf "/tags/%s" (. | urlize) }}
                <a href="{{ $tagUrl }}" class="tag">
                    {{ . }}
                </a>
            {{ end }}
        </div>
    </main>
    
    {{ template "partials/_footer.html" . }}
</body>
</html>
```

### Single Tag Page

**`templates/tag.html`**:

```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.lang }}">
{{ template "partials/_head.html" . }}
<body class="tag-page">
    {{ template "partials/_header.html" . }}
    
    <main>
        <header>
            <h1>Posts tagged: {{ .Tag.Name }}</h1>
            <p>{{ len .Tag.Pages }} posts</p>
        </header>
        
        <div class="posts">
            {{ range .Tag.Pages }}
                <article class="post-card">
                    <h2><a href="{{ .Permalink }}">{{ .Config.title }}</a></h2>
                    <time>{{ date .Config.date "January 2, 2006" }}</time>
                    <p>{{ .Config.description }}</p>
                </article>
            {{ end }}
        </div>
    </main>
    
    {{ template "partials/_footer.html" . }}
</body>
</html>
```

## 404 Error Page

**`templates/404.html`**:

```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.lang }}">
{{ template "partials/_head.html" . }}
<body class="error-404">
    {{ template "partials/_header.html" . }}
    
    <main>
        <div class="error-content">
            <h1>404</h1>
            <p>Page not found</p>
            <p>The page you're looking for doesn't exist.</p>
            
            <nav class="error-nav">
                <a href="/">← Back to home</a>
                <a href="/blog/">Browse posts</a>
            </nav>
            
            <!-- Suggest related content -->
            {{ $blog := get_section "blog" }}
            {{ if $blog }}
                <section class="suggested">
                    <h2>Popular Posts</h2>
                    {{ range first 5 (reverse $blog.Children) }}
                        <article>
                            <a href="{{ .Permalink }}">{{ .Config.title }}</a>
                        </article>
                    {{ end }}
                </section>
            {{ end }}
        </div>
    </main>
    
    {{ template "partials/_footer.html" . }}
</body>
</html>
```

## Documentation Template

**`templates/docs.html`** - Documentation with sidebar:

```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.lang }}">
{{ template "partials/_head.html" . }}
<body class="docs">
    {{ template "partials/_header.html" . }}
    
    <div class="docs-layout">
        <aside class="docs-sidebar">
            <nav>
                {{ $docs := get_section "docs" }}
                {{ range $docs.Children }}
                    <a href="{{ .Permalink }}" 
                       class="{{ if eq .Permalink $.Page.Permalink }}active{{ end }}">
                        {{ .Config.title }}
                    </a>
                {{ end }}
            </nav>
        </aside>
        
        <main class="docs-content">
            <article>
                <h1>{{ .Page.Config.title }}</h1>
                
                {{ if .Page.Config.extra.toc }}
                    <aside class="toc">
                        {{ .Toc }}
                    </aside>
                {{ end }}
                
                <div class="prose">
                    {{ .Page.Content }}
                </div>
                
                <nav class="docs-nav">
                    {{ with .Page.Higher }}
                        <a href="{{ .Permalink }}" class="prev">← {{ .Config.title }}</a>
                    {{ end }}
                    {{ with .Page.Lower }}
                        <a href="{{ .Permalink }}" class="next">{{ .Config.title }} →</a>
                    {{ end }}
                </nav>
            </article>
        </main>
    </div>
    
    {{ template "partials/_footer.html" . }}
</body>
</html>
```

## Portfolio/Projects Template

**`templates/projects.html`**:

```html
<!DOCTYPE html>
<html lang="{{ .Site.Config.lang }}">
{{ template "partials/_head.html" . }}
<body class="projects">
    {{ template "partials/_header.html" . }}
    
    <main>
        <h1>{{ .Section.Config.title }}</h1>
        
        <div class="projects-grid">
            {{ range .Section.Children }}
                <article class="project-card">
                    {{ if .Config.img }}
                        <img src="{{ .Config.img }}" alt="{{ .Config.title }}">
                    {{ end }}
                    
                    <h2>{{ .Config.title }}</h2>
                    <p>{{ .Config.description }}</p>
                    
                    <div class="project-meta">
                        {{ if .Config.extra.tech_stack }}
                            <div class="tech-stack">
                                {{ range .Config.extra.tech_stack }}
                                    <span class="tech">{{ . }}</span>
                                {{ end }}
                            </div>
                        {{ end }}
                        
                        <div class="project-links">
                            {{ if .Config.extra.github_url }}
                                <a href="{{ .Config.extra.github_url }}">GitHub</a>
                            {{ end }}
                            {{ if .Config.extra.demo_url }}
                                <a href="{{ .Config.extra.demo_url }}">Live Demo</a>
                            {{ end }}
                            <a href="{{ .Permalink }}">Details →</a>
                        </div>
                    </div>
                </article>
            {{ end }}
        </div>
    </main>
    
    {{ template "partials/_footer.html" . }}
</body>
</html>
```

---

**Related:**
- [Template Structure](/templates/structure)
- [Template Variables](/templates/variables)
- [Template Functions](/templates/functions)
- [Advanced Patterns](/templates/advanced)
