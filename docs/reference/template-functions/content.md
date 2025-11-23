# Content Functions

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
