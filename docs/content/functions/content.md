+++
title = "Content Functions"
date = 2025-12-15
template = "page.html"
+++

Functions for accessing and rendering content.

## Content Access

### `get_section`

Retrieves a section's pages.

```go
{{ $blog := get_section "blog" }}
{{ range $blog.Children }}
  <h2>{{ .Title }}</h2>
  <p>{{ .Description }}</p>
{{ end }}
```

### `markdown`

Renders Markdown to HTML.

```go
{{ markdown .Content }}
{{ markdown "# Hello\n\nWorld" }}
```

### `priority`

Returns first defined value.

```go
{{ priority .Page.Title .Site.Title }}
{{ priority .Description .Summary "No description" }}
```

### `related_posts`

Finds related posts using intelligent tag-based scoring.

**Key Features:**
- **O(k) Performance**: Tag index lookup, not O(n²) comparison
- **Smart Ranking**: Scores by tag overlap + recency + randomization
- **Variety**: Returns 6 candidates for client-side random selection

```go
{{ $section := get_section "blog" }}
{{ $related := related_posts .Page $section.Children }}

{{ if $related }}
<section id="related-posts">
  <h3>Related Posts</h3>
  {{ range $related }}
    <article>
      <a href="{{ .Permalink }}">{{ .Config.title }}</a>
    </article>
  {{ end }}
</section>
{{ end }}
```

**See [Collections](/functions/collections#related_posts) for detailed documentation.**

## Examples

### Section Listing

```go
{{ $blog := get_section "blog" }}
<section>
  <h2>{{ $blog.Title }}</h2>
  <ul>
    {{ range $blog.Children }}
      <li>
        <a href="{{ .Permalink }}">{{ .Title }}</a>
      </li>
    {{ end }}
  </ul>
</section>
```

### Fallback Values

```go
<title>{{ priority .Title .Site.Title }}</title>
<meta name="description" content="{{ priority .Description .Summary }}" />
```

### Inline Markdown

```go
<div class="note">
  {{ markdown .FrontMatter.note }}
</div>
```
